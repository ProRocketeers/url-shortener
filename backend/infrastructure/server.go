package infrastructure

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path"
	"time"

	v1 "github.com/ProRocketeers/url-shortener/api/v1"
	"github.com/ProRocketeers/url-shortener/docs"
	"github.com/ProRocketeers/url-shortener/domain"
	"github.com/ProRocketeers/url-shortener/domain/query"
	"github.com/ProRocketeers/url-shortener/domain/services"
	"github.com/ProRocketeers/url-shortener/domain/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"github.com/swaggo/swag"
	"gorm.io/gorm"
)

type dependencies struct {
	db                  *gorm.DB
	shortLinkRepository *storage.ShortLinkRepository
	shortLinkService    *services.ShortLinkService
	apiHandler          *v1.ApiHandler
	adminApiHandler     *v1.AdminApiHandler
}

func createDependencies(config Config) (dependencies, error) {
	db, err := ConnectToDatabase(config)
	if err != nil {
		return dependencies{}, fmt.Errorf("could not connect to database: %v", err)
	}
	shortLinkRepository := &storage.ShortLinkRepository{
		Repository: storage.Repository{
			DB: query.Use(db),
		},
	}
	shortLinkService := &services.ShortLinkService{
		Repository: shortLinkRepository,
		BaseUrl:    config.Domain.BaseUrl,
	}
	requestInfoRepository := &storage.RequestInfoRepository{
		Repository: storage.Repository{
			DB: query.Use(db),
		},
	}
	requestInfoService := &services.RequestInfoService{
		Repository: requestInfoRepository,
	}
	apiHandler := v1.NewApiHandler(shortLinkService, requestInfoService)
	adminApiHandler := v1.NewAdminApiHandler(shortLinkService, requestInfoService)

	return dependencies{db, shortLinkRepository, shortLinkService, apiHandler, adminApiHandler}, nil
}

func createRouter(dependencies *dependencies, config Config) *chi.Mux {
	r := chi.NewRouter()

	r.Use(
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Body != nil {
					rawBody, err := io.ReadAll(r.Body)
					if err != nil {
						log.Error().Err(err).Msg("failed to read request body in middleware")
						http.Error(w, "failed to read body", http.StatusBadRequest)
						return
					}

					ctx := context.WithValue(r.Context(), "body", rawBody)

					r.Body = io.NopCloser(bytes.NewReader(rawBody))

					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}
				next.ServeHTTP(w, r)
			})
		},
		// injects random generated `X-Request-ID` header
		middleware.RequestID,
		// sets the `request.RemoteAddr`
		middleware.RealIP,
		// writes a log line that contains:
		// - timestamp
		// - server hostname, random request ID
		// - route, HTTP version
		// - remote address
		// - status code, response size, response time
		middleware.RequestLogger(&ZerologChiFormatter{
			Logger: log.Logger,
		}),
		// error handling - recovers from panics, logs it and returns 500
		middleware.Recoverer,
		middleware.Heartbeat("/health"),
		// automatically redirects trailing slashes to the path without it
		middleware.RedirectSlashes,
		middleware.Timeout(60*time.Second),
	)

	r.Handle("/metrics", promhttp.Handler())

	basePath := func() string {
		if config.Domain.BaseUrl.Path != "" {
			return config.Domain.BaseUrl.Path
		}
		return "/"
	}()

	docs.SetupSwaggerParams(swag.Spec{
		Title:    "URL Shortener API",
		Version:  config.Metadata.Version,
		Host:     config.Domain.BaseUrl.Host,
		BasePath: basePath,
		Schemes: func() []string {
			if config.Environment == DevelopmentEnvironment {
				return []string{"http"}
			}
			return []string{"https"}
		}(),
	})

	r.Route(basePath, func(r chi.Router) {
		r.Get("/swagger*", httpSwagger.WrapHandler)
		r.Get("/swagger", func(w http.ResponseWriter, r *http.Request) {
			// properly resolve the redirect with respect to the base path
			redirectUrl := config.Domain.BaseUrl.ResolveReference(
				&url.URL{
					Path: path.Join(basePath, "swagger", "index.html"),
				},
			).String()
			http.Redirect(w, r, redirectUrl, http.StatusMovedPermanently)
		})

		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("hello"))
		})

		r.Route("/v1", func(r chi.Router) {
			r.Post("/shorten", dependencies.apiHandler.ShortenUrl)
			r.Get("/{slug:[a-zA-Z0-9]+}", dependencies.apiHandler.RedirectSlug)

			r.Route("/admin", func(r chi.Router) {
				r.Route("/link", func(r chi.Router) {
					r.Post("/", dependencies.adminApiHandler.CreateShortLink)
					r.Get("/list", dependencies.adminApiHandler.ListShortLinks)
					r.Get("/id/{id:\\d+}", dependencies.adminApiHandler.GetShortLinkById)
					r.Put("/id/{id:\\d+}", dependencies.adminApiHandler.UpdateShortLinkById)
					r.Delete("/id/{id:\\d+}", dependencies.adminApiHandler.DeleteShortLinkById)
					r.Get("/slug/{slug:[a-zA-Z0-9]+}", dependencies.adminApiHandler.GetShortLinkBySlug)
				})
				r.Route("/info", func(r chi.Router) {
					r.Get("/", dependencies.adminApiHandler.FindSingleRequestInfo)
					r.Get("/list", dependencies.adminApiHandler.ListRequestInfos)
				})
			})
		})
	})
	return r
}

func RunServerGracefully(config Config) error {
	dependencies, err := createDependencies(config)
	if err != nil {
		return fmt.Errorf("could not create server dependencies: %v", err)
	}

	router := createRouter(&dependencies, config)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port),
		Handler: router,
	}

	// Handle SIGTERM and SIGINT
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	serverError := make(chan error, 1)
	// Start the server in goroutine to not block main thread
	go func() {
		defer close(serverError)
		log.Info().
			Str("version", config.Metadata.Version).
			Str("commit", config.Metadata.CommitHash).
			Str("buildTime", config.Metadata.BuildTime).
			Msgf("Starting server on port %d", config.Port)

		cleanupTask := domain.CleanupTask{
			Context:  ctx,
			DB:       dependencies.db,
			Interval: config.Domain.ExpiredLinkCleanupInterval,
		}
		log.Info().Msgf("Starting background job - cleaning up expired links - every %v", cleanupTask.Interval.String())
		cleanupTask.Run()

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			// encountered an error, gracefully shutdown
			serverError <- err
		}
	}()

	// Wait for either context cancellation (interrupt), or a server error
	select {
	case err := <-serverError:
		log.Error().Err(err).Msg("Server error")
	case <-ctx.Done():
		log.Info().Msg("Shutting down gracefully...")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("graceful shutdown failed: %v", err)
	}
	log.Info().Msg("Server stopped")
	return nil
}
