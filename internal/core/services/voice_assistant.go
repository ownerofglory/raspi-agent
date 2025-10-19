package services

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ownerofglory/raspi-agent/internal/core/domain"
	"github.com/ownerofglory/raspi-agent/internal/core/ports"
)

// voiceAssistant implements the domain.VoiceAssistant interface.
//
// It orchestrates the three primary stages of a conversational AI pipeline:
//  1. Transcription (speech-to-text)
//  2. Completion (natural language response generation)
//  3. Speech synthesis (text-to-speech)
//
// Each stage delegates to an injected provider (STT, LLM, TTS) via ports,
// making the component easily testable and replaceable.
type voiceAssistant struct {
	transcription ports.TranscriptionProvider
	speech        ports.SpeechProvider
	completion    ports.CompletionProvider
}

// NewVoiceAssistant constructs a new voiceAssistant instance.
//
// The caller must provide concrete implementations of the following providers:
//   - transcription: converts user audio to text (STT engine)
//   - completion: generates responses based on the text (LLM engine)
//   - speech: converts response text back into speech (TTS engine)
//
// This composition allows modular configuration — for example, combining
// Whisper STT with GPT-based completion and Piper or OpenAI TTS.
func NewVoiceAssistant(stt ports.TranscriptionProvider, tts ports.SpeechProvider, cmpl ports.CompletionProvider) *voiceAssistant {
	return &voiceAssistant{
		transcription: stt,
		speech:        tts,
		completion:    cmpl,
	}
}

// Assist executes a full voice interaction flow.
//
// It performs the following steps sequentially:
//  1. Transcribes the input audio using the STT provider.
//  2. Sends the transcribed text to the LLM for completion.
//  3. Streams the synthesized speech of the LLM response.
//
// Returns a *receive-only* channel (<-chan *VoiceAssistantResult) that yields
// progressively generated audio chunks, allowing immediate playback while
// speech synthesis continues in the background.
//
// The context controls cancellation and timeout across all stages.
// If any stage fails, the function logs the error, cleans up, and closes
// the result channel gracefully.
func (v *voiceAssistant) Assist(ctx context.Context, req *domain.VoiceAssistantRequest) (<-chan *domain.VoiceAssistantResult, error) {
	tr := domain.TranscribeRequest{
		Audio: req.Audio,
	}
	transcribe, err := v.transcription.Transcribe(ctx, tr)
	if err != nil {
		slog.Error("Failed to transcribe", "error", err)
		return nil, fmt.Errorf("failed to transcribe: %w", err)
	}

	cr := domain.CompletionRequest{
		Prompt: transcribe.Text,
	}
	completion, err := v.completion.CreateCompletion(ctx, &cr)
	if err != nil {
		slog.Error("Failed to create completion", "error", err)
		return nil, fmt.Errorf("failed to create completion: %w", err)
	}

	sr := domain.SpeechRequest{
		Text: completion.Text,
	}
	speechCh, err := v.speech.ProduceSpeechAudio(ctx, &sr)
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
