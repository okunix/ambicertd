package router

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/okunix/ambicertd/errutil"
)

func Router(logger *slog.Logger) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "ok\n")
	})
	mux.HandleFunc("GET /notok", func(w http.ResponseWriter, r *http.Request) {
		logger.Error("hell nah", "error", errors.Join(errors.New("test"), errutil.WithStack(errors.New("test with stack"))))
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "notok\n")
	})

	handler := requestLogger(logger)(mux)
	return handler
}
