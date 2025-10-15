package audio

import (
	"encoding/binary"
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/gordonklaus/portaudio"
)

type recorder struct{}

type RecordingResult struct {
	data []int16
	err  error
}

func (r RecordingResult) GetError() error {
	return r.err
}

func (r RecordingResult) SaveTo(f io.Writer) error {
	channels := 2
	sampleRate := 4096
	byteRate := sampleRate * channels * 2
	blockAlign := channels * 2
	dataLen := len(r.data) * 2
	riffLen := 36 + dataLen

	f.Write([]byte("RIFF"))
	binary.Write(f, binary.LittleEndian, uint32(riffLen))
	f.Write([]byte("WAVE"))

	f.Write([]byte("fmt "))
	binary.Write(f, binary.LittleEndian, uint32(16))         // chunk size
	binary.Write(f, binary.LittleEndian, uint16(1))          // PCM format
	binary.Write(f, binary.LittleEndian, uint16(channels))   // channels
	binary.Write(f, binary.LittleEndian, uint32(sampleRate)) // sample rate
	binary.Write(f, binary.LittleEndian, uint32(byteRate))   // byte rate
	binary.Write(f, binary.LittleEndian, uint16(blockAlign)) // block align
	binary.Write(f, binary.LittleEndian, uint16(16))         // bits per sample

	f.Write([]byte("data"))
	binary.Write(f, binary.LittleEndian, uint32(dataLen))
	binary.Write(f, binary.LittleEndian, r.data)

	return nil
}

func NewRecorder() *recorder {
	return &recorder{}
}

func (r *recorder) RecordAudio(duration time.Duration) <-chan RecordingResult {
	ch := make(chan RecordingResult)

	go func() {
		if err := portaudio.Initialize(); err != nil {
			slog.Error("Unable to initialize portaudio", "err", err)
			ch <- RecordingResult{data: nil, err: fmt.Errorf("unable to initialize portaudio: %v", err)}
			return
		}
		defer portaudio.Terminate()

		slog.Debug("Portaudio initialized")

		seconds := duration.Seconds()
		channels := 1
		chunkSize := 4096

		devices, err := portaudio.Devices()
		if err != nil {
			slog.Error("Unable to get devices", "err", err)
			ch <- RecordingResult{data: nil, err: fmt.Errorf("unable to get devices: %v", err)}
			return
		}
		slog.Debug("Available devices:")
		for i, d := range devices {
			slog.Debug("%d: %s (inputs=%d, outputs=%d, default SR=%.0f)\n",
				i, d.Name, d.MaxInputChannels, d.MaxOutputChannels, d.DefaultSampleRate)
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
			ch <- RecordingResult{data: nil, err: fmt.Errorf("no input device found")}
			return
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
			ch <- RecordingResult{data: nil, err: fmt.Errorf("unable to open stream: %v", err)}
			return
		}
		defer stream.Close()
		slog.Debug("Stream opened. Recording audio...")

		if err := stream.Start(); err != nil {
			slog.Error("Unable to start stream", "err", err)
			ch <- RecordingResult{data: nil, err: fmt.Errorf("unable to start stream: %v", err)}
			return
		}

		for framesRecorded := 0; framesRecorded < int(totalFrames); framesRecorded += chunkSize {
			if err := stream.Read(); err != nil {
				if paErr, ok := err.(portaudio.Error); ok && paErr == portaudio.InputOverflowed {
					slog.Debug("Warning: input overflow (skipping some samples)")
					continue
				}
				panic(err)
			}
			buffer = append(buffer, chunk...)
		}

		if err := stream.Stop(); err != nil {
			slog.Error("Unable to stop stream", "err", err)
			ch <- RecordingResult{data: nil, err: fmt.Errorf("unable to stop stream: %v", err)}
			return
		}

		ch <- RecordingResult{data: buffer, err: nil}
	}()

	return ch
}
