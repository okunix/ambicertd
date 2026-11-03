package router

import (
	"log/slog"
	"net/http"
	"time"
)

type middleware func(http.Handler) http.Handler

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func requestLogger(logger *slog.Logger) middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now().UTC()
			writer := &wrappedWriter{ResponseWriter: w, statusCode: http.StatusOK}
			h.ServeHTTP(writer, r)
			latency := time.Since(startTime)
			logger.Info(
				"incoming request",
				slog.Int64("latency", latency.Microseconds()),
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("code", writer.statusCode),
			)
		})
	}
}
