package ports

import "context"

// WakeListener defines the contract for components capable of detecting
// a predefined wake word or activation phrase (e.g. “Hey Vicky”, “Raspi”)
// from live microphone input.
//
// Implementations are typically long-running processes that continuously
// listen to audio input and block until a wake word is detected or the
// provided context is cancelled.
type WakeListener interface {
	// Listen starts listening for the configured wake word using the
	// system’s default microphone or other input source.
	//
	// The method should block until:
	//   - The wake word is detected, in which case it returns nil.
	//   - The provided context is cancelled, in which case it returns ctx.Err().
	//   - An unrecoverable error occurs during initialization or processing.
	//
	// Implementations should ensure that any audio resources (e.g. PortAudio streams)
	// are properly initialized and cleaned up.
	Listen(ctx context.Context) error
}
