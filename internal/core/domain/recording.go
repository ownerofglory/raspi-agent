package domain

import "io"

// RecordingResult represents the result of an audio recording operation.
//
// Implementations of this interface encapsulate the recorded audio data
// and provide a method to persist it to an `io.Writer`, such as a file,
// network connection, or in-memory buffer.
//
// The `SaveTo` method writes the audio data (for example, in WAV format)
// to the provided writer and returns an error if any I/O operation fails.
//
// Example:
//
//	file, _ := os.Create("/tmp/sample.wav")
//	defer file.Close()
//
//	if err := recording.SaveTo(file); err != nil {
//	    log.Fatalf("failed to save recording: %v", err)
//	}
type RecordingResult interface {
	// SaveTo writes the recorded audio data to the given writer.
	// Implementations must ensure that the audio is serialized in a
	// standard, readable format (e.g. PCM-encoded WAV).
	//
	// Returns an error if the write operation fails.
	SaveTo(writer io.Writer) error
}
