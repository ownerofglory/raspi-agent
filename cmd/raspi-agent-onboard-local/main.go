package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/ownerofglory/raspi-agent/internal/audio"
	"github.com/ownerofglory/raspi-agent/internal/core/services"
	"github.com/ownerofglory/raspi-agent/internal/onboard"
	"github.com/ownerofglory/raspi-agent/internal/openaiapi"
	"github.com/ownerofglory/raspi-agent/internal/wakeword"
)

var (
	porcupineAccessKey   = flag.String("porcupineAccessKey", "", "porcupine SDK access key from 'https://console.picovoice.ai/'")
	porcupineLibPath     = flag.String("porcupineLibPath", "", "porcupine library path e.g. 'lib/libpv_porcupine.so'")
	porcupineModelPath   = flag.String("porcupineModelPath", "", "porcupine model parameters path, e.g. 'resources/porcupine_params.pv'")
	porcupineKeywordPath = flag.String("porcupineKeywordPath", "", "porcupine keyword path, e.g. 'resources/Hey-Rhaspy_en_raspberry-pi_v3_0_0.ppn'")

	openAIURL    = flag.String("openAIURL", "", "OpenAI base URL")
	openAIAPIKey = flag.String("openAIAPIKey", "", "OpenAI API token")
)

func main() {
	flag.Parse()

	slog.SetLogLoggerLevel(slog.LevelDebug)
	slog.Debug("Starting application")

	listener := wakeword.NewPorcupineListener(*porcupineAccessKey, *porcupineModelPath, *porcupineLibPath, *porcupineKeywordPath)
	recorder := audio.NewRecorder()
	player := audio.NewPortAudioPlayer()

	c := openai.NewClient(option.WithAPIKey(*openAIAPIKey),
		option.WithBaseURL(*openAIURL))

	tts := openaiapi.NewTextToSpeechClient(&c)
	stt := openaiapi.NewSpeechToTextClient(&c)
	cmpl := openaiapi.NewCompletionClient(&c)

	assistant := services.NewVoiceAssistant(stt, tts, cmpl)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	orch := onboard.NewOrchestrator(listener, recorder, player, assistant)
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
