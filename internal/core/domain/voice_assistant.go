package domain

import "io"

// VoiceAssistantRequest represents a single input to the voice assistant.
//
// It encapsulates an audio stream (typically recorded user speech) to be processed
// by the assistant pipeline — which may include transcription (speech-to-text),
// natural language understanding, and text-to-speech synthesis.
//
// The Audio field is an io.Reader, allowing flexible data sources such as:
//   - A microphone stream (e.g., PortAudio or ALSA)
//   - A temporary audio file
//   - A network or pipe stream .
type VoiceAssistantRequest struct {
	Audio io.Reader
}

// VoiceAssistantResult represents a single output message from the assistant.
//
// It encapsulates the generated audio stream that corresponds to the assistant's
// spoken reply. Depending on the pipeline design, multiple results may be streamed
// progressively through a channel to enable low-latency playback.
//
// The Audio field is an io.Reader that provides raw audio data — typically
// in a playable format (e.g., MPEG or WAV) suitable for immediate playback.
//
// The consumer (e.g., onboard agent or client) is responsible for reading and
// playing or saving the audio data as it arrives.
type VoiceAssistantResult struct {
	Audio io.Reader
}
