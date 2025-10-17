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
