package ports

import (
	"context"

	"github.com/ownerofglory/raspi-agent/internal/core/domain"
)

// VoiceAssistant defines the contract for a full voice interaction pipeline.
//
// Implementations of this interface handle the entire voice-assistant flow:
//  1. Accepting an audio input request (typically recorded speech from a user).
//  2. Transcribing the speech into text via an STT (speech-to-text) provider.
//  3. Generating a natural language response using an LLM or dialogue engine.
//  4. Converting the response text back into spoken audio via a TTS (text-to-speech) provider.
//
// The Assist method executes this pipeline asynchronously and returns a
// *receive-only* channel (<-chan *VoiceAssistantResult>) that streams
// intermediate or final results as they become available — such as transcribed
// text, partial responses, or synthesized audio.
//
// The provided context allows cancellation and timeout control across the
// entire interaction lifecycle (for example, to stop processing when the
// user cancels the request or a network error occurs).
type VoiceAssistant interface {
	// Assist begins processing the user's audio request and returns a channel
	// that streams VoiceAssistantResult messages containing transcribed text,
	// LLM responses, or audio output.
	//
	// Implementations may choose to:
	//   - Stream results incrementally (e.g., progressive transcription or TTS),
	//   - Or buffer and return a final, aggregated response.
	//
	// The context controls cancellation and timeout behavior.
	Assist(ctx context.Context, req *domain.VoiceAssistantRequest) (<-chan *domain.VoiceAssistantResult, error)
}
