package router

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/okunix/ambicertd/logutil"
)

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
			ctx := r.Context()
			loggerCtx := &logutil.LoggerContext{}
			r = r.WithContext(context.WithValue(ctx, logutil.LoggerContextKey, loggerCtx))

			startTime := time.Now().UTC()
			writer := &wrappedWriter{ResponseWriter: w, statusCode: http.StatusOK}
			h.ServeHTTP(writer, r)

			lg := logger
			if loggerCtx.UserID != "" {
				lg = lg.With("user_id", loggerCtx.UserID)
			}

			if loggerCtx.Error != nil {
				lg = lg.With("error", loggerCtx.UserID)
			}

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
