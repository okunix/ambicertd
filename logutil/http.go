package logutil

import (
	"context"
	"errors"
	"log/slog"
)

type HttpErrorLogger struct {
	slog.Handler
}

func (h HttpErrorLogger) Handle(ctx context.Context, r slog.Record) error {
	newRecord := slog.NewRecord(r.Time, r.Level, "http server exception", r.PC)
	newRecord.AddAttrs(slog.Any("error", errors.New(r.Message)))
	r.Attrs(func(a slog.Attr) bool {
		newRecord.AddAttrs(a)
		return true
	})
	return h.Handler.Handle(ctx, newRecord)
}
