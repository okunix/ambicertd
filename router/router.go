package router

import (
	"fmt"
	"log/slog"
	"net/http"
)

func Router(logger *slog.Logger) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "ok\n")
	})
	mux.HandleFunc("GET /notok", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "notok\n")
	})

	handler := requestLogger(logger)(mux)
	return handler
}
