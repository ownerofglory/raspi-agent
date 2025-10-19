package handler

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/ownerofglory/raspi-agent/internal/core/domain"
	"github.com/ownerofglory/raspi-agent/internal/core/ports"
)

const (
	filePathPattern = "/tmp/voice%s.wav"

	PostReceiveVoiceAssistance = basePath + "/v1/voice-assistance"
)

type voiceAssistantResponse struct {
	AudioChunk string `json:"audioChunk"`
}

// voiceAssistantHandler handles HTTP requests for voice assistant operations.
//
// It bridges the HTTP layer with the core domain layer by translating incoming
// multipart audio uploads into domain.VoiceAssistantRequest objects and
// streaming the resulting synthesized audio directly to the client.
type voiceAssistantHandler struct {
	assistant ports.VoiceAssistant
}

// NewVoiceAssistantHandler constructs a new HTTP handler for voice assistant requests.
//
// The handler requires a concrete implementation of the ports.VoiceAssistant interface,
// which orchestrates transcription, completion, and TTS in the backend.
func NewVoiceAssistantHandler(va ports.VoiceAssistant) *voiceAssistantHandler {
	return &voiceAssistantHandler{
		assistant: va,
	}
}

type voiceAssistantRequest struct {
	Audio []byte `json:"audio"`
}

func NewVoiceAssistant() *voiceAssistantHandler {
	return &voiceAssistantHandler{}
}

// HandleAssist processes an uploaded audio file, runs it through the voice assistant pipeline,
// and streams the resulting synthesized audio back to the client.
//
// Request:
//   - Method: POST
//   - Content-Type: multipart/form-data
//   - Form field: "audio" (audio file)
//
// Response:
//   - Content-Type: audio/mpeg
//   - Transfer-Encoding: chunked
//   - The connection is kept alive to stream generated audio progressively.
//
// Flow:
//  1. The uploaded audio file is saved temporarily.
//  2. It’s passed into the assistant pipeline for processing.
//  3. The handler streams the resulting audio chunks as they become available.
//
// Example client usage (curl):
//
//	curl -X POST -F "audio=@sample.wav" http://<host>/v1/voice-assistance --output reply.mp3
func (v *voiceAssistantHandler) HandleAssist(rw http.ResponseWriter, r *http.Request) {
	audioFile, fh, err := r.FormFile("audio")
	if err != nil {
		slog.Error("Unable to get form data file", "err", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	defer audioFile.Close()
	defer r.Body.Close()

	filePath := fmt.Sprintf(filePathPattern, time.Now().Format("20060102150405"))
	tmpFile, err := os.Create(filePath)
	if err != nil {
		slog.Error("Unable to create temporary file", "err", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer tmpFile.Close()
	defer os.Remove(filePath)

	written, err := io.Copy(tmpFile, audioFile)
	if err != nil {
		slog.Error("Unable to copy audio to file", "err", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	if written != fh.Size {
		slog.Warn("Audio file copy size missmatch", "size", fh.Size, "written", written)
	}
	// Reopen saved file for reading
	tmpFile.Seek(0, io.SeekStart)

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()
	resCh, err := v.assistant.Assist(ctx, &domain.VoiceAssistantRequest{
		Audio: tmpFile,
	})
	if err != nil {
		slog.Error("Unable to assist audio", "err", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "audio/mpeg")
	rw.Header().Set("Cache-Control", "no-cache")
	rw.Header().Set("Connection", "keep-alive")
	rw.Header().Set("Transfer-Encoding", "chunked")

	flusher, ok := rw.(http.Flusher)
	if !ok {
		http.Error(rw, "streaming not supported", http.StatusInternalServerError)
		return
	}

	slog.Info("Streaming audio response to client...")

	for {
		select {
		case <-ctx.Done():
			slog.Info("Client disconnected or request canceled")
			return

		case res, ok := <-resCh:
			if !ok {
				slog.Info("Assistant stream completed")
				return
			}

			// Stream the raw audio bytes directly to the client
			_, err := io.Copy(rw, res.Audio)
			if err != nil {
				slog.Error("Error writing audio chunk", "err", err)
				return
			}
			flusher.Flush()
		}
	}
}
