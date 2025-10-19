package audio

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/gordonklaus/portaudio"
	"github.com/ownerofglory/raspi-agent/internal/core/domain"
)

// WAV format constants define the structure and encoding
// of PCM (Pulse-Code Modulation) audio data inside a .wav file.
const (
	// pcmFormat identifies linear PCM (uncompressed) audio data.
	// Value 1 means standard PCM encoding.
	pcmFormat = 1

	// bitsPerSample defines the resolution of each sample.
	// 16 bits per sample = CD-quality audio.
	bitsPerSample = 16

	// bytesPerSample is derived from bitsPerSample (16 bits = 2 bytes).
	// It’s used when calculating byte rates and data lengths.
	bytesPerSample = bitsPerSample / 8

	// fmtChunkSize is always 16 for PCM format WAV files.
	// This chunk describes format details such as sample rate and bit depth.
	fmtChunkSize = 16
)

// WAV header identifiers are the ASCII tags that mark sections of a .wav file.
// Each is exactly 4 bytes and identifies a logical block of the file.
const (
	// riffHeader identifies the overall file as a RIFF container.
	// It’s followed by the total file size and the "WAVE" format specifier.
	riffHeader = "RIFF"

	// waveHeader specifies that the RIFF container stores audio data in WAVE format.
	waveHeader = "WAVE"

	// fmtHeader marks the beginning of the format chunk,
	// which stores sample rate, channel count, and encoding type.
	fmtHeader = "fmt "

	// dataHeader marks the beginning of the actual audio sample data.
	dataHeader = "data"
)

type recorder struct{}

// NewRecorder creates a new recorder instance
func NewRecorder() *recorder {
	return &recorder{}
}

// recordingResult holds raw PCM data and metadata
type recordingResult struct {
	data []int16

	channels   int
	chunkSize  int
	sampleRate int
}

// SaveTo writes the recorded PCM data as a WAV file to an io.Writer
func (r *recordingResult) SaveTo(f io.Writer) error {
	channels := r.channels
	sampleRate := r.sampleRate
	byteRate := sampleRate * channels * bytesPerSample
	blockAlign := channels * bytesPerSample
	dataLen := len(r.data) * bytesPerSample
	riffLen := 36 + dataLen

	// Write WAV headers
	f.Write([]byte(riffHeader))
	binary.Write(f, binary.LittleEndian, uint32(riffLen))
	f.Write([]byte(waveHeader))

	f.Write([]byte(fmtHeader))
	binary.Write(f, binary.LittleEndian, uint32(fmtChunkSize))
	binary.Write(f, binary.LittleEndian, uint16(pcmFormat))
	binary.Write(f, binary.LittleEndian, uint16(channels))
	binary.Write(f, binary.LittleEndian, uint32(sampleRate))
	binary.Write(f, binary.LittleEndian, uint32(byteRate))
	binary.Write(f, binary.LittleEndian, uint16(blockAlign))
	binary.Write(f, binary.LittleEndian, uint16(bitsPerSample))

	f.Write([]byte(dataHeader))
	binary.Write(f, binary.LittleEndian, uint32(dataLen))
	binary.Write(f, binary.LittleEndian, r.data)

	return nil
}

// RecordAudio captures audio for a specified duration
func (r *recorder) RecordAudio(ctx context.Context, duration time.Duration) (domain.RecordingResult, error) {
	if err := portaudio.Initialize(); err != nil {
		slog.Error("Unable to initialize portaudio", "err", err)
		return nil, fmt.Errorf("unable to initialize portaudio: %v", err)
	}
	defer portaudio.Terminate()

	slog.Debug("Portaudio initialized")

	seconds := duration.Seconds()
	channels := 1
	chunkSize := 4096

	devices, err := portaudio.Devices()
	if err != nil {
		slog.Error("Unable to get devices", "err", err)
		return nil, fmt.Errorf("unable to get devices: %v", err)
	}
	slog.Debug("Available devices:")
	for i, d := range devices {
		slog.Debug("", "idx", i,
			"name", d.Name, "inputs", d.MaxInputChannels, "outputs", d.MaxOutputChannels, "SR", d.DefaultSampleRate)
	}

	var inputDevice *portaudio.DeviceInfo
	for _, d := range devices {
		if d.MaxInputChannels > 0 {
			inputDevice = d
			break
		}
	}
	if inputDevice == nil {
		slog.Error("no input device found")
		return nil, fmt.Errorf("no input device found")
	}
	slog.Debug("Using input device:", "device", inputDevice.Name)

	// use device's default sample rate
	sampleRate := inputDevice.DefaultSampleRate
	totalFrames := sampleRate * seconds

	// prepare buffers
	buffer := make([]int16, 0, int(totalFrames)*channels)
	chunk := make([]int16, chunkSize*channels)

	stream, err := portaudio.OpenStream(portaudio.StreamParameters{
		Input: portaudio.StreamDeviceParameters{
			Device:   inputDevice,
			Channels: channels,
			Latency:  inputDevice.DefaultLowInputLatency,
		},
		SampleRate:      sampleRate,
		FramesPerBuffer: chunkSize,
	}, chunk)

	if err != nil {
		slog.Error("Unable to open stream", "err", err)
		return nil, fmt.Errorf("unable to open stream: %v", err)
	}
	defer stream.Close()
	slog.Debug("Stream opened. Recording audio...")

	if err := stream.Start(); err != nil {
		slog.Error("Unable to start stream", "err", err)
		return nil, fmt.Errorf("unable to start stream: %v", err)
	}

	for framesRecorded := 0; framesRecorded < int(totalFrames); framesRecorded += chunkSize {
		select {
		case <-ctx.Done():
			slog.Debug("Audio recording cancelled")
			return nil, fmt.Errorf("Audio recording cancelled")
		default:
			if err := stream.Read(); err != nil {
				var paErr portaudio.Error
				if errors.As(err, &paErr) && errors.Is(paErr, portaudio.InputOverflowed) {
					slog.Debug("Warning: input overflow (skipping some samples)")
					continue
				}
				slog.Error("Streaming error", "err", err)
				return nil, fmt.Errorf("streaming error: %v", err)
			}
			buffer = append(buffer, chunk...)
		}
	}

	if err := stream.Stop(); err != nil {
		slog.Error("Unable to stop stream", "err", err)
		return nil, fmt.Errorf("unable to stop stream: %v", err)
	}

	slog.Debug("Recording finished successfully")
	return &recordingResult{
		data:       buffer,
		channels:   channels,
		chunkSize:  chunkSize,
		sampleRate: int(sampleRate),
	}, nil
}
