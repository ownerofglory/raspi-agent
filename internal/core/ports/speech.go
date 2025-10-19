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
	// Transcribe converts spoken audio into text based on the parameters
	// defined in the TranscribeRequest.
	//
	// The context allows for cancellation and timeout control.
	// Returns a pointer to TranscribeResult on success or an error if
	// the transcription process fails.
	Transcribe(ctx context.Context, req domain.TranscribeRequest) (*domain.TranscribeResult, error)
}

// SpeechProvider defines the contract for any Text-to-Speech (TTS) service.
//
// Implementations of this interface convert text input (via a SpeechRequest)
// into an audio stream represented by a sequence of SpeechResult messages.
//
// The provider may support multiple output formats or transport methods —
// for example, streaming audio over Server-Sent Events (SSE) or as raw
// audio data (MPEG/WAV) suitable for direct playback.
//
// Both methods return a *receive-only* channel (<-chan *SpeechResult) that
// yields audio chunks progressively as they are generated. The channel must
// be closed when the generation process completes or fails.
//
// The provided context allows for cancellation and timeout control during
// the streaming lifecycle (for example, when a user cancels playback).
type SpeechProvider interface {
	// ProduceSpeechSSE generates spoken audio from text and streams it as
	// Server-Sent Events (SSE) — suitable for clients that expect JSON-based
	// event messages (e.g., browsers or frontend dashboards).
	//
	// Each SpeechResult contains a base64-encoded audio chunk or other event
	// payload. The channel is closed when the stream completes or fails.
	//
	// Use this method when your transport layer requires SSE (text/event-stream).
	ProduceSpeechSSE(ctx context.Context, req *domain.SpeechRequest) (<-chan *domain.SpeechResult, error)

	// ProduceSpeechAudio generates spoken audio from text and streams it
	// as raw binary audio chunks — for example, MPEG or WAV data that can
	// be piped directly into a speaker or audio playback process.
	//
	// Each SpeechResult contains an io.Reader providing a portion of the
	// audio output. The caller is responsible for sequentially reading and
	// combining the chunks for continuous playback.
	//
	// Use this method when working with direct audio playback pipelines,
	// such as Raspberry Pi speakers or hardware audio devices.
	ProduceSpeechAudio(ctx context.Context, req *domain.SpeechRequest) (<-chan *domain.SpeechResult, error)
}
