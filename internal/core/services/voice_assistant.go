package services

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ownerofglory/raspi-agent/internal/core/domain"
	"github.com/ownerofglory/raspi-agent/internal/core/ports"
)

type voiceAssistant struct {
	stt  ports.TranscriptionProvider
	tts  ports.SpeechProvider
	cmpl ports.CompletionProvider
}

func NewVoiceAssistant() *voiceAssistant {
	return &voiceAssistant{}
}

func (v *voiceAssistant) Assist(ctx context.Context, req *domain.VoiceAssistantRequest) (<-chan *domain.VoiceAssistantResult, error) {
	tr := domain.TranscribeRequest{
		Audio: req.Audio,
	}
	transcribe, err := v.stt.Transcribe(ctx, tr)
	if err != nil {
		slog.Error("Failed to transcribe", "error", err)
		return nil, fmt.Errorf("failed to transcribe: %w", err)
	}

	cr := domain.CompletionRequest{
		Prompt: transcribe.Text,
	}
	completion, err := v.cmpl.CreateCompletion(ctx, &cr)
	if err != nil {
		slog.Error("Failed to create completion", "error", err)
		return nil, fmt.Errorf("failed to create completion: %w", err)
	}

	sr := domain.SpeechRequest{
		Text: completion.Text,
	}
	speechCh, err := v.tts.ProduceSpeech(ctx, &sr)
	if err != nil {
		slog.Error("Failed to produce speech", "error", err)
		return nil, fmt.Errorf("failed to produce speech: %w", err)
	}

	resCh := make(chan *domain.VoiceAssistantResult)

	go func() {
		defer close(resCh)
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-speechCh:
				if !ok {
					slog.Warn("Speech channel closed")
					return
				}

				r := domain.VoiceAssistantResult{
					Audio: msg.Audio,
				}

				resCh <- &r
			}
		}
	}()

	return resCh, nil
}
