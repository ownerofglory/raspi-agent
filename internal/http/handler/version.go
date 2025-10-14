package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// go build -ldflags "-X github.com/ownerofglory/raspi-agent/internal/http/handler.AppVersion=v1.2.3"
var AppVersion string

const GetVersionEndpoint = basePath + "/version"

type VersionResponse struct {
	Version string `json:"version"`
}

// HandleGetVersion http handler function that returns application version
//
// Response:
//   - 200 OK with "application/json" version response VersionResponse
func HandleGetVersion(rw http.ResponseWriter, _ *http.Request) {
	vr := VersionResponse{
		Version: AppVersion,
	}

	payload, err := json.Marshal(vr)
	if err != nil {
		slog.Error("failed to marshal version response")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Header().Add("Content-Type", "application/json")
	_, _ = rw.Write(payload)
}
