package logutil

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/okunix/ambicertd/errutil"
)

const LoggerContextKey = "loggerContext"

type LoggerContext struct {
	UserID string
	Error  error
}

const errorAttrKey = "error"

func expandError(err error) slog.Attr {
	if multierr, ok := err.(errutil.MultiError); ok {
		errs := multierr.Unwrap()
		attrs := make([]slog.Attr, 0, len(errs))
		for k, v := range errs {
			attr := expandError(v)
			attr.Key = fmt.Sprintf("error_%d", k)
			attrs = append(attrs, attr)
		}
		return slog.GroupAttrs("errors", attrs...)
	}

	attrs := []slog.Attr{slog.String("msg", err.Error())}

	attrs = append(attrs, errutil.Attrs(err)...)

	if stackTrace := errutil.StackTrace(err); stackTrace != "" {
		attrs = append(attrs, slog.String("stack_trace", stackTrace))
	}

	return slog.GroupAttrs(errorAttrKey, attrs...)
}

func replaceAttr(groups []string, a slog.Attr) slog.Attr {
	if a.Key == errorAttrKey {
		err, ok := a.Value.Any().(error)
		if !ok {
			return a
		}
		return expandError(err)
	}
	return a
}

func Logger(level slog.Level) (*slog.Logger, error) {
	stdoutHandler := slog.NewJSONHandler(os.Stdout,
		&slog.HandlerOptions{Level: level, ReplaceAttr: replaceAttr})
	return slog.New(stdoutHandler), nil
}
