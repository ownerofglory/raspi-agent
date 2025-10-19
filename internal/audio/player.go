package audio

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"

	"github.com/gordonklaus/portaudio"
	"github.com/tosone/minimp3"
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
		slog.Error("Error initializing portaudio")
		return fmt.Errorf("portaudio initialize failed: %w", err)
	}
	defer portaudio.Terminate()

	// Pick output device (pulse/pipewire preferred)
	devices, err := portaudio.Devices()
	if err != nil {
		slog.Error("Error getting devices")
		return fmt.Errorf("list devices: %w", err)
	}

	var output *portaudio.DeviceInfo
	for _, d := range devices {
		name := strings.ToLower(d.Name)
		if strings.Contains(name, "pulse") || strings.Contains(name, "pipewire") {
			output = d
			break
		}
	}
	if output == nil {
		slog.Error("no pulse/pipewire output device found")
		return fmt.Errorf("no pulse/pipewire output device found")
	}

	// Pipe for streaming
	pr, pw := io.Pipe()

	// Write chunks from audioStream into pipe
	go func() {
		defer pw.Close()
		for {
			select {
			case <-ctx.Done():
				slog.Debug("Stopping portaudio playback")
				return
			case chunk, ok := <-audioStream:
				if !ok {
					slog.Debug("Audio stream channel closed")
					return
				}
				if len(chunk) > 0 {
					_, _ = pw.Write(chunk)
				}
			}
		}
	}()

	// Initialize minimp3 decoder
	decoder, err := minimp3.NewDecoder(pr)
	if err != nil {
		return fmt.Errorf("failed to create mp3 decoder: %w", err)
	}
	defer decoder.Close()

	// Wait until decoding actually starts
	started := decoder.Started()
	<-started

	slog.Info("MP3 decoder started", "rate", decoder.SampleRate, "channels", decoder.Channels)

	// Open PortAudio stream
	channels := decoder.Channels
	if channels == 0 {
		channels = 2
	}
	buf := make([]int16, 512*channels)
	stream, err := portaudio.OpenStream(portaudio.StreamParameters{
		Output: portaudio.StreamDeviceParameters{
			Device:   output,
			Channels: channels,
			Latency:  output.DefaultLowOutputLatency,
		},
		SampleRate:      float64(decoder.SampleRate),
		FramesPerBuffer: len(buf) / channels,
	}, &buf)
	if err != nil {
		return fmt.Errorf("failed to open audio stream: %w", err)
	}
	defer stream.Close()

	if err := stream.Start(); err != nil {
		return fmt.Errorf("failed to start stream: %w", err)
	}
	defer stream.Stop()

	slog.Info("🎧 Playing decoded audio via PortAudio", "device", output.Name)

	// Playback loop
	pcmBuf := make([]byte, 4096)
	for {
		select {
		case <-ctx.Done():
			slog.Info("Playback canceled")
			return nil
		default:
			n, err := decoder.Read(pcmBuf)
			if n > 0 {
				// Convert []byte -> []int16 without allocation
				samples := bytesToInt16(pcmBuf[:n])
				for len(samples) > 0 {
					copied := copy(buf, samples)
					samples = samples[copied:]
					if err := stream.Write(); err != nil {
						if strings.Contains(err.Error(), "underflow") {
							slog.Warn("Output underflow")
							continue
						}
						slog.Error("Stream write failed", "err", err)
						return err
					}
				}
			}
			if err != nil {
				if err == io.EOF {
					slog.Info("Playback finished")
					return nil
				}
				return fmt.Errorf("decoder read failed: %w", err)
			}
		}
	}
}

// helper to interpret little-endian bytes as int16 samples
func bytesToInt16(data []byte) []int16 {
	n := len(data) / 2
	samples := make([]int16, n)
	for i := 0; i < n; i++ {
		samples[i] = int16(uint16(data[2*i]) | uint16(data[2*i+1])<<8)
	}
	return samples
}
