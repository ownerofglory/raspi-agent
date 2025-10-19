package wakeword

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/gordonklaus/portaudio"
	porcupine "github.com/sigidagi/porcupine/binding/go/v2"
)

// porcupineListener provides real-time wake-word detection using the
// Picovoice Porcupine engine.
//
// It encapsulates the lifecycle of the Porcupine instance and manages
// an audio input stream from the system microphone via PortAudio.
type porcupineListener struct {
	p *porcupine.Porcupine
}

// NewPorcupineListener creates and configures a new Porcupine listener.
//
// Parameters:
//   - accessKey:    Picovoice access key used to authenticate the SDK.
//   - modelPath:    Path to the Porcupine model file (e.g. "porcupine_params.pv").
//   - libraryPath:  Path to the shared library (e.g. "libpv_porcupine.so").
//   - keywordPath:  Path to the keyword file (.ppn) for the desired wake word.
//
// Returns:
//
//	A ready-to-use *porcupineListener instance.
func NewPorcupineListener(accessKey, modelPath, libraryPath, keywordPath string) *porcupineListener {
	var p = &porcupine.Porcupine{
		AccessKey:   accessKey,
		ModelPath:   modelPath,
		LibraryPath: libraryPath,
		KeywordPaths: []string{
			keywordPath,
		},
	}
	return &porcupineListener{
		p: p,
	}
}

// Listen continuously captures audio input and processes it using Porcupine
// until a wake word is detected or the provided context is cancelled.
//
// When the wake word is detected, Listen returns nil. If an error occurs
// during initialization or audio processing, it returns a descriptive error.
//
// Typical usage:
//
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//	listener := NewPorcupineListener(...)
//	err := listener.Listen(ctx)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println("Wake word detected!")
func (l *porcupineListener) Listen(ctx context.Context) error {
	err := l.p.Init()
	if err != nil {
		slog.Error("porcupine: failed to initialize porcupine", "err", err)
		return fmt.Errorf("porcupine: failed to initialize porcupine: %w", err)
	}
	defer l.p.Delete()

	slog.Debug("Porcupine ready", "SampleRate", porcupine.SampleRate, "FrameLength", porcupine.FrameLength)
	if err := portaudio.Initialize(); err != nil {
		slog.Error("portaudio init:", "err", err)
		return fmt.Errorf("portaudio init:: %w", err)
	}
	defer portaudio.Terminate()

	buf := make([]int16, porcupine.FrameLength)
	stream, err := portaudio.OpenDefaultStream(1, 0, float64(porcupine.SampleRate), len(buf), buf)
	if err != nil {
		slog.Error("OpenDefaultStream:", "err", err)
		return fmt.Errorf("porcupine: OpenDefaultStream:: %w", err)
	}
	defer stream.Close()

	if err := stream.Start(); err != nil {
		slog.Error("stream start:", "err", err)
		return fmt.Errorf("porcupine: stream start:: %w", err)
	}

	slog.Debug("Listening for wake word...")

	for {
		select {
		case <-ctx.Done():
		default:
			if err := stream.Read(); err != nil {
				slog.Warn("stream read error:", "err", err)
				continue
			}

			res, err := l.p.Process(buf)
			if err != nil {
				slog.Warn("porcupine process error:", "err", err)
				continue
			}
			if res >= 0 {
				slog.Debug("Wake word detected!")
				return nil
			}
		}
	}
}
