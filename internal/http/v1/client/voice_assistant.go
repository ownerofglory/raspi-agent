package client

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
)

const PostReceiveAssistanceURL string = backendBasePath + "/v1/voice-assistance"

type voiceAssistant struct {
	client  *http.Client
	baseURL string
}

func NewVoiceAssistant(baseURL string) *voiceAssistant {
	return &voiceAssistant{
		client: &http.Client{
			Timeout: 0,
		},
	}
}

func (v *voiceAssistant) ReceiveVoiceAssistance(ctx context.Context, filePath string) (<-chan []byte, error) {
	url := fmt.Sprintf("%s%s", v.baseURL, PostReceiveAssistanceURL)

	file, err := os.Open(filePath)
	if err != nil {
		slog.Error("Error opening file", "err", err)
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	// prepare multipart/form-data request
	bodyReader, bodyWriter := io.Pipe()
	multipartWriter := multipart.NewWriter(bodyWriter)

	go func() {
		defer bodyWriter.Close()
		defer file.Close()
		part, err := multipartWriter.CreateFormFile("audio", filePath)
		if err != nil {
			slog.Error("failed to create multipart form", "err", err)
			return
		}
		if _, err := io.Copy(part, file); err != nil {
			slog.Error("failed to copy file to multipart", "err", err)
			return
		}
		multipartWriter.Close()
	}()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", multipartWriter.FormDataContentType())

	// send request
	resp, err := v.client.Do(req)
	if err != nil {
		slog.Error("failed to send request to voice assistant", "err", err)
		return nil, fmt.Errorf("failed to send audio: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		slog.Error("failed to send request to voice assistant", "status", resp.Status)
		defer resp.Body.Close()
		return nil, fmt.Errorf("backend returned status: %s", resp.Status)
	}

	ch := make(chan []byte)

	// read streaming audio from response
	go func() {
		defer resp.Body.Close()
		defer close(ch)

		audioBuf := make([]byte, 4096)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				n, err := resp.Body.Read(audioBuf)
				if n > 0 {
					// copy the slice to avoid overwriting by next read
					copyBuf := make([]byte, n)
					copy(copyBuf, audioBuf[:n])
					ch <- copyBuf
				}
				if err != nil {
					if err == io.EOF {
						slog.Info("Audio stream finished")
						return
					}
					slog.Error("failed to read audio stream", "err", err)
					return
				}
			}
		}
	}()

	return ch, nil
}
