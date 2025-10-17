package ports

import (
	"context"

	"github.com/ownerofglory/raspi-agent/internal/core/domain"
)

// TranscriptionProvider defines the contract for any service capable of
// converting spoken audio into text (speech-to-text).
//
// Implementations of this interface may call external APIs (e.g., OpenAI Whisper),
// or run a local transcription engine (e.g., Whisper.cpp, Vosk, DeepSpeech, etc.).
//
// The Transcribe method takes a context for cancellation or timeout control,
// and a TranscribeRequest containing the audio stream or file to process.
//
// It returns a TranscribeResult with the recognized text, or an error if
// transcription fails.
type TranscriptionProvider interface {
	Transcribe(ctx context.Context, req domain.TranscribeRequest) (*domain.TranscribeResult, error)
}
