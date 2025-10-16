package ports

import (
	"context"
	"time"

	"github.com/ownerofglory/raspi-agent/internal/core/domain"
)

type Recorder interface {
	RecordAudio(ctx context.Context, duration time.Duration) (domain.RecordingResult, error)
}
