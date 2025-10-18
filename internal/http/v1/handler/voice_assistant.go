package handler

import (
	"context"
	"encoding/base64"
	"encoding/json"
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
)

type voiceAssistantResponse struct {
	AudioChunk string `json:"audioChunk"`
}

type voiceAssistantHandler struct {
	assistant ports.VoiceAssistant
}

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

	written, err := io.Copy(tmpFile, audioFile)
	if err != nil {
		slog.Error("Unable to copy audio to file", "err", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	if written != fh.Size {
		slog.Warn("Audio file copy size missmatch", "size", fh.Size, "written", written)
	}

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()
	resCh, err := v.assistant.Assist(ctx, &domain.VoiceAssistantRequest{
		Audio: audioFile,
	})
	if err != nil {
		slog.Error("Unable to assist audio", "err", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "text/event-stream")
	rw.Header().Set("Cache-Control", "no-cache")
	rw.Header().Set("Connection", "keep-alive")
	rw.Header().Set("Transfer-Encoding", "chunked")

	for {
		res, ok := <-resCh
		if !ok {
			slog.Debug("Audio channel closed")
			break
		}

		data, err := io.ReadAll(res.Audio)
		if err != nil {
			slog.Error("Unable to read audio", "err", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		dst := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
		base64.StdEncoding.Encode(dst, data)

		voiceChunk := voiceAssistantResponse{
			AudioChunk: string(dst),
		}

		chunk, err := json.Marshal(voiceChunk)
		if err != nil {
			slog.Error("Error marshalling chunk", "err", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err = fmt.Fprintf(rw, "data: %v\n\n", string(chunk))
		if err != nil {
			slog.Error("Error writing data", "err", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		rw.(http.Flusher).Flush()
	}
}
