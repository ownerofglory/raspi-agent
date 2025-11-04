package ports

import (
	"context"

	"github.com/ownerofglory/raspi-agent/internal/core/domain"
)

// AgentMemory defines the interface for a component that can recall stored
// information or "memories" for a given user or context.
//
// Implementations of AgentMemory may use different backends—such as
// databases, vector stores, or in-memory caches—to retrieve relevant data
// based on the supplied RecallRequest.
//
// The Recall method should return a RecallResult containing any matching
// memories, or an error if the recall operation fails.
type AgentMemory interface {
	// Recall retrieves memories for the given user or query context.
	//
	// The ctx argument allows cancellation or timeout of the recall operation.
	// The req parameter specifies the user and optional query text.
	// The method returns a RecallResult with the retrieved memories or an error.
	Recall(ctx context.Context, req domain.RecallRequest) (*domain.RecallResult, error)
}
