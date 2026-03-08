package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"backend/internal/config"
	"backend/internal/httpapi"
)

func main() {
	handler := slog.NewTextHandler(os.Stdout, nil)
	slog.SetDefault(slog.New(handler))

	if err := run(); err != nil {
		slog.Error("server exited", "error", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	app, err := httpapi.New()
	if err != nil {
		return err
	}

	address := fmt.Sprintf("0.0.0.0:%d", cfg.Port)
	server := &http.Server{
		Addr:              address,
		Handler:           app.Router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	slog.Info("starting backend", "environment", cfg.Environment, "address", address)

	err = server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}

	return err
}
