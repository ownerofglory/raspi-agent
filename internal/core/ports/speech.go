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

// SpeechProvider defines the contract for any text-to-speech (TTS) service.
//
// Implementations of this interface take text input (a SpeechRequest) and produce
// corresponding audio output as a stream of SpeechResult messages.
//
// The ProduceSpeech method returns a *receive-only* channel (<-chan *SpeechResult),
// allowing the caller to consume audio chunks as they are generated — for example,
// when using a streaming TTS engine like OpenAI Realtime API, Piper, or Coqui TTS.
//
// The context allows cancellation and timeout control during generation.
type SpeechProvider interface {
	ProduceSpeech(ctx context.Context, req *domain.SpeechRequest) (<-chan *domain.SpeechResult, error)
}
