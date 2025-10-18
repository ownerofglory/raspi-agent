package ports

import (
	"context"
	"time"

	"github.com/ownerofglory/raspi-agent/internal/core/domain"
)

// Recorder defines the contract for capturing audio input from a microphone
// or other input device.
//
// Implementations of this interface are responsible for managing the
// audio stream lifecycle — initializing the input device, capturing samples,
// and returning the recorded data as a `domain.RecordingResult`.
//
// The captured audio may later be saved (e.g. as a WAV file) or streamed
// directly to another service, depending on the implementation.
type Recorder interface {
	// RecordAudio starts recording audio for the specified duration.
	//
	// It must block until the recording completes or the context is cancelled.
	// Implementations should handle proper device cleanup and return
	// any errors encountered during capture.
	//
	// Returns:
	//   - A `domain.RecordingResult` containing the recorded audio samples.
	//   - An `error` if initialization or recording fails.
	RecordAudio(ctx context.Context, duration time.Duration) (domain.RecordingResult, error)
}
