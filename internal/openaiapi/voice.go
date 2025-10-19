package openaiapi

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"

	"github.com/openai/openai-go/v3"
	"github.com/ownerofglory/raspi-agent/internal/core/domain"
	"github.com/tmaxmax/go-sse"
)

type textToSpeech struct {
	client *openai.Client
}

func NewTextToSpeechClient(client *openai.Client) *textToSpeech {
	return &textToSpeech{client: client}
}

const (
	textToSpeechEventTypeDelta = "speech.audio.delta"
	textToSpeechEventTypeDone  = "speech.audio.done"
)

type textToSpeechEvent struct {
	Type        string `json:"type"`
	Usage       string `json:"usage"`
	AudioBase64 string `json:"audio"`
}

func (c *textToSpeech) ProduceSpeechSSE(ctx context.Context, req *domain.SpeechRequest) (<-chan *domain.SpeechResult, error) {
	params := openai.AudioSpeechNewParams{
		Input:        req.Text,
		Model:        openai.SpeechModelGPT4oMiniTTS,
		Voice:        openai.AudioSpeechNewParamsVoiceShimmer,
		StreamFormat: openai.AudioSpeechNewParamsStreamFormatSSE,
	}
	response, err := c.client.Audio.Speech.New(ctx, params)
	if err != nil {
		slog.Error("Failed to Text to Speech", "err", err)
		return nil, fmt.Errorf("failed to Text to Speech: %w", err)
	}

	ch := make(chan *domain.SpeechResult)

	go func() {
		defer response.Body.Close()
		defer close(ch)

		for ev, err := range sse.Read(response.Body, nil) {
			select {
			case <-ctx.Done():
				slog.Warn("Speech generation canceled")
				return
			default:
				if err != nil {
					slog.Error("error processing an event", "err", err)
					return
				}
				var event textToSpeechEvent
				err := json.Unmarshal([]byte(ev.Data), &event)
				if err != nil {
					slog.Error("error unmarshalling an event", "err", err)
					return
				}

				if event.Type == textToSpeechEventTypeDone {
					slog.Debug("Speech event is done")
					return
				}

				audioData, err := base64.StdEncoding.DecodeString(event.AudioBase64)
				if err != nil {
					slog.Error("error decoding audio", "err", err)
					return
				}

				ch <- &domain.SpeechResult{
					Audio: bytes.NewReader(audioData),
				}
			}
		}
	}()

	return ch, nil
}

func (c *textToSpeech) ProduceSpeechAudio(ctx context.Context, req *domain.SpeechRequest) (<-chan *domain.SpeechResult, error) {
	params := openai.AudioSpeechNewParams{
		Input:        req.Text,
		Model:        openai.SpeechModelGPT4oMiniTTS,
		Voice:        openai.AudioSpeechNewParamsVoiceShimmer,
		StreamFormat: openai.AudioSpeechNewParamsStreamFormatAudio,
	}
	response, err := c.client.Audio.Speech.New(ctx, params)
	if err != nil {
		slog.Error("Failed to Text to Speech", "err", err)
		return nil, fmt.Errorf("failed to Text to Speech: %w", err)
	}

	ch := make(chan *domain.SpeechResult)

	go func() {
		defer response.Body.Close()
		defer close(ch)

		buf := make([]byte, 4096) // 4KB buffer for streaming chunks
		for {
			n, err := response.Body.Read(buf)
			if err != nil {
				if err == io.EOF {
					slog.Debug("TTS stream completed")
					break
				}
				slog.Error("TTS stream read error", "err", err)
				break
			}

			if n > 0 {
				// Send the current audio chunk downstream.
				chunk := make([]byte, n)
				copy(chunk, buf[:n])

				select {
				case ch <- &domain.SpeechResult{Audio: io.NopCloser(io.NewSectionReader(
					bytes.NewReader(chunk), 0, int64(len(chunk)),
				))}:
					// successfully sent
				case <-ctx.Done():
					slog.Warn("TTS stream canceled by context")
					return
				}
			}
		}
	}()

	return ch, nil
}
