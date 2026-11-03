package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/okunix/ambicertd/config"
	"github.com/okunix/ambicertd/logutil"
	"github.com/okunix/ambicertd/router"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "/etc/ambicertd/ambicertd.yaml", "ambicertd config file path")
}

func main() {
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	exitCode := run(ctx)
	cancel()
	os.Exit(exitCode)
}

func run(ctx context.Context) int {
	// initializing logger
	logger, err := logutil.Logger(slog.LevelDebug)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize logger: %s\n", err.Error())
		return 1
	}

	// reading config file
	cfg, err := config.ReadFile(configPath)
	if err != nil {
		logger.Error("failed to read config file", "error", err)
		return 1
	}

	// starting server
	httpErrorHandler := logutil.HttpErrorLogger{Handler: logger.Handler()}
	server := http.Server{
		Addr:           cfg.Server.Addr,
		Handler:        router.Router(logger),
		MaxHeaderBytes: 1 << 20,
		ErrorLog:       slog.NewLogLogger(httpErrorHandler, slog.LevelError),
	}

	logger.Debug("Server is running", "addr", server.Addr)

	var serverErr error
	doneCh := make(chan struct{})
	go func() {
		defer close(doneCh)
		if cfg.Server.TLS != nil {
			logger.Debug("TLS is enabled", "cert", cfg.Server.TLS.Cert)
			serverErr = server.ListenAndServeTLS(
				cfg.Server.TLS.Cert,
				cfg.Server.TLS.Key,
			)
			return
		}
		serverErr = server.ListenAndServe()
	}()

	// waiting for termination
	select {
	case <-ctx.Done():
	case <-doneCh:
	}

	// shutting down
	logger.Debug("Server is shutting down")

	if serverErr != nil {
		logger.Error("server error", "error", serverErr)
		return 1
	}

	timeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(timeout); err != nil {
		logger.Error("failed to shutdown the server", "error", err)
		return 1
	}

	return 0
}
