package main

import (
	"context"
	"errors"
	"github.com/caarlos0/env/v11"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
	"github.com/ownerofglory/raspi-agent/config"
	"github.com/ownerofglory/raspi-agent/internal/http/v1/handler"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	slog.Info("Starting app")

	// Config parsing
	var cfg config.RaspiAgentConfig
	err := env.Parse(&cfg)
	if err != nil {
		slog.Error("Failed to parse config", "error", err)
		os.Exit(1)
	}

	// Logger setup
	logLevel := slog.LevelInfo
	if err := logLevel.UnmarshalText([]byte(cfg.LogLevel)); err != nil {
		logLevel = slog.LevelInfo
	}
	slog.SetLogLoggerLevel(logLevel)

	logger := httplog.NewLogger("cvthis", httplog.Options{
		LogLevel: logLevel,
	})

	slog.SetDefault(logger.Logger)

	// Chi setup
	r := chi.NewRouter()

	// HTTP handler setup
	r.Get(handler.GetVersionEndpoint, handler.HandleGetVersion)

	httpServer := http.Server{
		Addr:    cfg.ServerAddr,
		Handler: r,
	}

	go func() {
		slog.Info("Starting HTTP Server")

		err := httpServer.ListenAndServe()

		if !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Server shutdown unexpected", "err", err)
		}

		slog.Info("HTTP Server finished")
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		slog.Error("HTTP shutdown error:", "err", err)
	}

	slog.Info("App finished")
}
