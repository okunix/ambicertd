package logutil

import (
	"log/slog"
	"os"

	"github.com/okunix/ambicertd/errutil"
)

const errorKey = "error"

func replaceAttr(groups []string, a slog.Attr) slog.Attr {
	if a.Key == errorKey {
		err, ok := a.Value.Any().(error)
		if !ok {
			return a
		}
		attrs := []slog.Attr{slog.String("msg", err.Error())}
		attrs = append(attrs, errutil.Attrs(err)...)
		if stackTrace := errutil.StackTrace(err); stackTrace != "" {
			attrs = append(attrs, slog.String("stack_trace", stackTrace))
		}
		return slog.GroupAttrs(errorKey, attrs...)
	}
	return a
}

func Logger(level slog.Level) (*slog.Logger, error) {
	stdoutHandler := slog.NewJSONHandler(os.Stdout,
		&slog.HandlerOptions{Level: level, ReplaceAttr: replaceAttr})
	return slog.New(stdoutHandler), nil
}
