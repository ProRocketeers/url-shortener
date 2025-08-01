package infrastructure

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
)

type ZerologChiFormatter struct {
	Logger zerolog.Logger
}

// https://pkg.go.dev/github.com/go-chi/chi/v5@v5.2.2/middleware#LogFormatter
func (f *ZerologChiFormatter) NewLogEntry(r *http.Request) middleware.LogEntry {
	ctx := r.Context()
	requestId := middleware.GetReqID(ctx)
	return &ZeroLogEntry{
		Logger: f.Logger.With().
			Str("method", r.Method).
			Str("url", r.URL.String()).
			Str("requestId", requestId).
			Str("remoteAddr", r.RemoteAddr).
			Logger(),
	}
}

type ZeroLogEntry struct {
	Logger zerolog.Logger
}

// https://pkg.go.dev/github.com/go-chi/chi/v5@v5.2.2/middleware#LogEntry
func (e *ZeroLogEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra any) {
	event := e.Logger.Info().
		Int("status", status).
		Int("bytes", bytes).
		Dur("duration", elapsed)

	if extra != nil {
		event = event.Interface("extra", extra)
	}

	event.Send()
}

func (e *ZeroLogEntry) Panic(v any, stack []byte) {
	e.Logger.Error().Interface("panic", v).Bytes("stack", stack).Msg("panic recovered")
}
