package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/ProRocketeers/url-shortener/api"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func createRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Use(
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

	r.Get("/swagger*", httpSwagger.WrapHandler)
	r.Get("/swagger", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/index.html", http.StatusMovedPermanently)
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello"))
	})

	r.Post("/shorten", api.ShortenUrl)
	r.Get("/{slug:[a-zA-Z0-9]+}", api.RedirectSlug)

	return r
}

func RunServerGracefully(config ServerConfig) error {
	router := createRouter()

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
		log.Info().Str("version", config.Version).Msgf("Starting server on port %d", config.Port)
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
