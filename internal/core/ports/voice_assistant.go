package ports

import (
	"context"

	"github.com/ownerofglory/raspi-agent/internal/core/domain"
)

type VoiceAssistant interface {
	Assist(ctx context.Context, req *domain.VoiceAssistantRequest) (<-chan *domain.VoiceAssistantResult, error)
}
