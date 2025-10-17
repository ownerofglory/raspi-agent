package domain

import "io"

// TranscribeRequest represents a request payload sent to the transcription service.
// It contains the audio data (e.g., a WAV or Opus stream) that should be converted to text.
//
// The `Audio` field is an io.Reader rather than a []byte, allowing the handler to stream
// large audio files without loading the entire content into memory.
type TranscribeRequest struct {
	Audio io.Reader `json:"audio"`
}

// TranscribeResult represents the response returned by the transcription service.
// The `Text` field contains the recognized speech from the provided audio input.
//
// Example JSON response:
//
//	{
//	  "text": "hello, how is the weather today?"
//	}
type TranscribeResult struct {
	Text string `json:"text"`
}

// SpeechRequest represents the input payload for speech synthesis (TTS).
//
// The Text field contains the text that should be converted into spoken audio.
// Example:
//
//	req := &domain.SpeechRequest{
//	    Text: "Hello, I am Rhaspy!",
//	}
type SpeechRequest struct {
	Text string `json:"text"`
}

// SpeechResult represents a single chunk or complete piece of generated audio.
//
// The Audio field provides the audio data as an io.Reader, enabling the consumer
// to stream it directly to a player, encoder, or network socket without needing
// to buffer the entire file in memory.
//
// Example usage:
//
//	for result := range speechChan {
//	    io.Copy(outputWriter, result.Audio)
//	}
type SpeechResult struct {
	Audio io.Reader
}
