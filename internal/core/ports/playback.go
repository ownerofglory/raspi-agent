package ports

import "context"

// Player defines the contract for any audio playback implementation.
//
// Implementations consume an incoming stream of audio data (as byte slices)
// and play it through the device’s audio output in real time.
//
// Typical use cases include:
//   - Playing raw PCM or decoded MP3/Opus chunks streamed from a backend.
//   - Acting as the final stage in a voice assistant pipeline after TTS.
//
// The PlaybackStream method must block until playback is complete,
// or return early if the context is canceled.
//
// Example usage:
//
//	player := audio.NewPlayer()
//	err := player.PlaybackStream(ctx, audioStream)
//	if err != nil {
//	    log.Fatal(err)
//	}
type Player interface {
	// PlaybackStream plays streamed audio chunks in real time.
	//
	// audioStream delivers consecutive chunks of audio data ([]byte),
	// typically representing small portions of PCM or decoded audio frames.
	// Implementations should handle variable chunk sizes gracefully.
	//
	// The context allows cancellation of playback (for example, if the
	// user stops the assistant or starts a new request).
	PlaybackStream(ctx context.Context, audioStream <-chan []byte) error
}
