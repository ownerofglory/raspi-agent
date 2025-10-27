package ports

import (
	"context"

	"github.com/ownerofglory/raspi-agent/internal/core/domain"
)

// EmbeddingProvider defines the interface for generating text embeddings.
// Implementations of this interface convert input text into numerical
// vector representations suitable for similarity search, clustering,
// or other machine learning tasks.
type EmbeddingProvider interface {
	// CreateEmbedding generates an embedding vector for the provided input text.
	// The method takes a context for cancellation and a request specifying
	// the input data, and returns the computed embedding result or an error.
	CreateEmbedding(ctx context.Context, req domain.EmbeddingRequest) (domain.EmbeddingResult, error)
}
