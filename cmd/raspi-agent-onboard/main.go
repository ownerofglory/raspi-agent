package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/ownerofglory/raspi-agent/internal/audio"
	"github.com/ownerofglory/raspi-agent/internal/http/v1/client"
	"github.com/ownerofglory/raspi-agent/internal/orchestrator"
	"github.com/ownerofglory/raspi-agent/internal/wakeword"
)

var (
	porcupineAccessKey   = flag.String("porcupineAccessKey", "", "porcupine SDK access key from 'https://console.picovoice.ai/'")
	porcupineLibPath     = flag.String("porcupineLibPath", "", "porcupine library path e.g. 'lib/libpv_porcupine.so'")
	porcupineModelPath   = flag.String("porcupineModelPath", "", "porcupine model parameters path, e.g. 'resources/porcupine_params.pv'")
	porcupineKeywordPath = flag.String("porcupineKeywordPath", "", "porcupine keyword path, e.g. 'resources/Hey-Rhaspy_en_raspberry-pi_v3_0_0.ppn'")

	backendBaseURL = flag.String("backendBaseURL", "", "Backend base URL")
)

func main() {
	flag.Parse()

	slog.SetLogLoggerLevel(slog.LevelDebug)
	slog.Debug("Starting application")

	listener := wakeword.NewPorcupineListener(*porcupineAccessKey, *porcupineModelPath, *porcupineLibPath, *porcupineKeywordPath)
	recorder := audio.NewRecorder()
	player := audio.NewPortAudioPlayer()
	assistant := client.NewVoiceAssistant(*backendBaseURL)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	orch := orchestrator.NewOnboardOrchestrator(listener, recorder, player, assistant)
	go func() {
		err := orch.Run(ctx)
		if err != nil {
			slog.Error("Error running onboardorchestrator", "error", err)
			return
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	slog.Debug("Received signal, shutting down")
}
