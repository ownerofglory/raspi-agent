package onboard

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/ownerofglory/raspi-agent/internal/core/domain"
	"github.com/ownerofglory/raspi-agent/internal/core/ports"
)

type onboardOrchestrator struct {
	listener       ports.WakeListener
	recorder       ports.Recorder
	player         ports.Player
	voiceAssistant ports.VoiceAssistant
}

func NewOrchestrator(listener ports.WakeListener, recorder ports.Recorder, player ports.Player, voiceAssistant ports.VoiceAssistant) *onboardOrchestrator {
	return &onboardOrchestrator{
		listener:       listener,
		recorder:       recorder,
		player:         player,
		voiceAssistant: voiceAssistant,
	}
}

func (o *onboardOrchestrator) Run(ctx context.Context) error {
	ctxWithCancel, cancel := context.WithCancel(ctx)

	recordCh := make(chan domain.RecordingResult)
	defer close(recordCh)

	// wake-record routine
	go func() {
		err := o.recordUponWake(ctxWithCancel, recordCh)
		if err != nil {
			cancel()
			return
		}
	}()

	fpCh := make(chan string)
	defer close(fpCh)

	// record processing routine
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case recordResult, ok := <-recordCh:
				if !ok {
					return
				}
				filePath := "recording" + time.Now().Format(time.RFC3339) + ".wav"
				func() {
					f, err := os.Create(filePath)
					if err != nil {
						slog.Error("Failed to create recording file", "err", err)
						cancel()
						return
					}
					defer f.Close()

					err = recordResult.SaveTo(f)
					if err != nil {
						slog.Error("Unable to save recording", "error", err)
						cancel()
						return
					}
				}()
				fpCh <- filePath
			}
		}
	}()

	// request backend and play
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case filePath, ok := <-fpCh:
				if !ok {
					return
				}

				file, err := os.Open(filePath)
				if err != nil {
					slog.Error("Failed to open file", "err", err)
					return
				}
				defer file.Close()

				req := domain.VoiceAssistantRequest{
					Audio: file,
				}
				assistance, err := o.voiceAssistant.Assist(ctx, &req)
				if err != nil {
					slog.Error("Unable to receive voice", "error", err)
					return
				}

				streamCh := make(chan []byte)
				go func() {
					err = o.player.PlaybackStream(ctx, streamCh)
					if err != nil {
						slog.Error("Unable to playback", "error", err)
						return
					}
				}()

				for {
					select {
					case <-ctx.Done():
						return
					case res, ok := <-assistance:
						if !ok {
							return
						}

						data, err := io.ReadAll(res.Audio)
						if err != nil {
							slog.Error("Unable to read audio", "error", err)
							return
						}

						streamCh <- data
					}
				}
			}
		}
	}()

	<-ctxWithCancel.Done()
	return nil
}

func (o *onboardOrchestrator) recordUponWake(ctx context.Context, resCh chan<- domain.RecordingResult) error {
	wakeCh := make(chan error)
	defer close(wakeCh)

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context canceled: %w", ctx.Err())
		default:
			// wake listener routine
			go func() {
				err := o.listener.Listen(ctx)
				if err != nil {
					slog.Error("Failed to listen for a wake word", "error", err)
					wakeCh <- fmt.Errorf("failed to listen for a wake word: %w", err)
					return
				}

				slog.Debug("Wake word detected")
				wakeCh <- nil
			}()

			err := <-wakeCh
			if err != nil {
				slog.Error("Failed to send a wake word", "error", err)
				return fmt.Errorf("failed to send a wake word: %w", err)
			}

			audio, err := o.recorder.RecordAudio(ctx, 8*time.Second)
			if err != nil {
				slog.Error("Failed to record audio input", "error", err)
				return fmt.Errorf("failed to record audio input: %w", err)
			}

			resCh <- audio
		}
	}
}
