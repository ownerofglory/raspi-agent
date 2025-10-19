package audio

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"log/slog"

	"github.com/gordonklaus/portaudio"
)

const (
	sampleRate           = 44100
	channels             = 2
	playerBytesPerSample = 2 // 16-bit audio
	framesPerBuffer      = 512
)

// portAudioPlayer handles playback of streamed audio using PortAudio.
// It consumes chunks of []byte audio data (typically raw PCM or decoded MP3).
type portAudioPlayer struct{}

// NewPortAudioPlayer constructs a new audio player instance.
func NewPortAudioPlayer() *portAudioPlayer {
	return &portAudioPlayer{}
}

// PlaybackStream plays streamed audio chunks in real time using PortAudio.
//
// The audioStream channel should deliver small PCM audio chunks (e.g., 4KB each).
// It automatically stops playback when the channel closes or the context is canceled.
//
// If your backend returns MP3, decode to PCM (int16 samples) before calling this.
func (p *portAudioPlayer) PlaybackStream(ctx context.Context, audioStream <-chan []byte) error {
	if err := portaudio.Initialize(); err != nil {
		slog.Error("portaudio initialize failed", "error", err)
		return fmt.Errorf("portaudio initialize failed: %w", err)
	}
	defer portaudio.Terminate()

	// Create output buffer and stream
	buf := make([]int16, framesPerBuffer*channels)
	stream, err := portaudio.OpenDefaultStream(0, channels, float64(sampleRate), len(buf), &buf)
	if err != nil {
		return fmt.Errorf("failed to open audio stream: %w", err)
	}
	defer stream.Close()

	if err := stream.Start(); err != nil {
		return fmt.Errorf("failed to start audio stream: %w", err)
	}
	defer stream.Stop()

	slog.Info("Playing streamed...")

	for {
		select {
		case <-ctx.Done():
			slog.Error("Audio stream stopped")
			return fmt.Errorf("playback stream stopped: %w", ctx.Err())
		case audio, ok := <-audioStream:
			if !ok {
				slog.Debug("Audio stream closed")
				return nil
			}
			// Convert bytes → int16 samples
			samples := bytesToInt16(audio)

			// Play in chunks sized to framesPerBuffer
			for len(samples) > 0 {
				n := copy(buf, samples)
				samples = samples[n:]

				if err := stream.Write(); err != nil {
					slog.Warn("stream write error", "err", err)
					break
				}
			}
		}
	}
}

// bytesToInt16 converts a byte slice (little-endian PCM) to int16 samples.
func bytesToInt16(data []byte) []int16 {
	if len(data)%2 != 0 {
		data = data[:len(data)-1] // trim stray byte
	}
	samples := make([]int16, len(data)/2)
	buf := bytes.NewReader(data)
	_ = binary.Read(buf, binary.LittleEndian, &samples)
	return samples
}
