package ports

import (
	"context"

	"github.com/ownerofglory/raspi-agent/internal/core/domain"
)

//go:generate go tool go.uber.org/mock/mockgen -source=voice_assistant.go -package=ports -destination=voice_assistant_mock.go VoiceAssistant,VoiceAssistantClient

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

// VoiceAssistantClient defines the contract for a client capable of sending
// recorded user audio to a backend voice assistant service and receiving
// a streamed audio reply.
//
// The backend is typically an HTTP or WebSocket endpoint that accepts
// an audio file (e.g., WAV or MP3) and returns a streaming audio response
// (for example, an `audio/mpeg` chunked transfer).
//
// The ReceiveVoiceAssistance method uploads the user's voice recording
// to the backend and returns a *receive-only* channel of []byte chunks,
// each containing a portion of the assistant's spoken response.
//
// Example:
//
//	ch, err := client.ReceiveVoiceAssistance(ctx, "/tmp/request.wav")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Stream response to speaker
//	player.PlaybackStream(ctx, ch)
type VoiceAssistantClient interface {
	// ReceiveVoiceAssistance sends a voice request to the backend and
	// returns a stream of audio chunks representing the assistant's reply.
	//
	// The returned channel will be closed automatically when the stream ends
	// or if the context is canceled.
	ReceiveVoiceAssistance(ctx context.Context, filePath string) (<-chan []byte, error)
}
