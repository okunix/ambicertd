package errutil

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

	pkgerr "github.com/pkg/errors"
)

type ErrorWithAttrs struct {
	err   error
	attrs []slog.Attr
}

func (e ErrorWithAttrs) Error() string {
	return e.err.Error()
}

func (e ErrorWithAttrs) Attrs() []slog.Attr {
	return e.attrs
}

func (e ErrorWithAttrs) Unwrap() error {
	return e.err
}

func WithAttrs(err error, attrs ...any) error {
	return ErrorWithAttrs{err: err, attrs: argsToAttrs(attrs...)}
}

func argsToAttrs(args ...any) []slog.Attr {
	attrs := make([]slog.Attr, 0, len(args))
	for i := 0; i < len(args); {
		switch key := args[i].(type) {
		case slog.Attr:
			attrs = append(attrs, key)
			i++
		case string:
			if i+1 == len(args) {
				attrs = append(attrs, slog.String("!BADKEY", key))
				i++
			} else {
				attrs = append(attrs, slog.Any(key, args[i+1]))
				i += 2
			}
		default:
			attrs = append(attrs, slog.Any("!BADKEY", key))
			i++
		}
	}
	return attrs
}

type attrError interface {
	Attrs() []slog.Attr
}

func Attrs(err error) []slog.Attr {
	var attrs []slog.Attr
	for err != nil {
		if ae, ok := err.(attrError); ok {
			attrs = append(attrs, ae.Attrs()...)
		}
		err = errors.Unwrap(err)
	}
	return attrs
}

type stackTracer interface {
	StackTrace() pkgerr.StackTrace
}

func StackTrace(err error) string {
	stackTrace := new(strings.Builder)
	for err != nil {
		if st, ok := err.(stackTracer); ok {
			fmt.Fprintf(stackTrace, "%+v", st)
		}
		err = errors.Unwrap(err)
	}
	return stackTrace.String()
}

func WithStack(err error) error {
	return pkgerr.WithStack(err)
}
