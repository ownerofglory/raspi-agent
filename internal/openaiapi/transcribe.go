package openaiapi

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/openai/openai-go/v3"
	"github.com/ownerofglory/raspi-agent/internal/core/domain"
)

type speechToText struct {
	client *openai.Client
}

func NewSpeechToTextClient(client *openai.Client) *speechToText {
	return &speechToText{client: client}
}

func (s *speechToText) Transcribe(ctx context.Context, req domain.TranscribeRequest) (*domain.TranscribeResult, error) {
	params := openai.AudioTranscriptionNewParams{
		Model: openai.AudioModelGPT4oMiniTranscribe,
		File:  req.Audio,
	}
	res, err := s.client.Audio.Transcriptions.New(ctx, params)
	if err != nil {
		slog.Error("Failed to transcribe audio", "err", err)
		return nil, fmt.Errorf("failed to transcribe audio: %w", err)
	}

	return &domain.TranscribeResult{Text: res.Text}, nil
}
