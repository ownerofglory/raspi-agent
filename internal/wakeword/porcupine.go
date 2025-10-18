package wakeword

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/gordonklaus/portaudio"
	porcupine "github.com/sigidagi/porcupine/binding/go/v2"
)

const modelPath = "model"
const libraryPath = "model"

type porcupineListener struct {
	p *porcupine.Porcupine
}

func NewListener() *porcupineListener {
	var p = &porcupine.Porcupine{
		AccessKey:   "",
		LibraryPath: libraryPath,
		ModelPath:   modelPath,
	}
	return &porcupineListener{
		p: p,
	}
}

type WakeListenerResult struct {
	err error
}

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
