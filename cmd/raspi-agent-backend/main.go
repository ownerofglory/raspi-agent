package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/ownerofglory/raspi-agent/config"
	"github.com/ownerofglory/raspi-agent/internal/core/services"
	"github.com/ownerofglory/raspi-agent/internal/http/v1/handler"
	"github.com/ownerofglory/raspi-agent/internal/openaiapi"
	"github.com/ownerofglory/raspi-agent/internal/persistence"
	"github.com/ownerofglory/raspi-agent/internal/stepca"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
	logger := httplog.NewLogger("raspi-agent-backend", httplog.Options{
		LogLevel: logLevel,
	})
	slog.SetDefault(logger.Logger)

	// Certificate provider setup
	certProvider := stepca.NewProvider(cfg.StepCAURL,
		cfg.StepCAProvisionerName,
		cfg.StepCAProvisionerToken,
		[]byte(cfg.StepCAPEM),
		[]byte(cfg.StepCAJWK))

	// ORM setup
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", cfg.PostgresHost, cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresDB, cfg.PostgresPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	// Repo setup
	deviceRepo := persistence.NewDeviceRepo(db)
	userRepo := persistence.NewUserRepository(db)

	// AI client setup
	openAIClient := openai.NewClient(option.WithAPIKey(cfg.OpenAIAPIKey), option.WithBaseURL(cfg.OpenAIAPIURL))
	tts := openaiapi.NewTextToSpeechClient(&openAIClient)
	stt := openaiapi.NewSpeechToTextClient(&openAIClient)
	cmpl := openaiapi.NewCompletionClient(&openAIClient)

	// service setup
	deviceService := services.NewDeviceService(userRepo, deviceRepo, certProvider)
	deviceHandler := handler.NewDeviceHandler(deviceService)

	// voice assistant setup
	va := services.NewVoiceAssistant(stt, tts, cmpl)
	vh := handler.NewVoiceAssistantHandler(va)

	fs := http.FileServer(http.Dir("ui/dist"))

	// Chi setup
	r := chi.NewRouter()

	// HTTP handler registration
	r.Post(handler.PostReceiveVoiceAssistance, vh.HandleAssist)
	r.Post(handler.PostRegisterDeviceURL, deviceHandler.HandlePostRegisterDevice)
	r.Post(handler.PostEnrollDeviceURL, deviceHandler.HandlePostEnrollDevice)
	r.Get(handler.GetVersionEndpoint, handler.HandleGetVersion)
	// UI
	r.Get(handler.BaseUIPath+"*", http.StripPrefix(handler.BaseUIPath, fs).ServeHTTP)

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
