package openaiapi

import (
	"context"
	"fmt"

	"github.com/openai/openai-go/v3"
	"github.com/ownerofglory/raspi-agent/internal/core/domain"
)

type embeddingClient struct {
	client *openai.Client
}

func NewEmbeddingClient(client *openai.Client) *embeddingClient {
	return &embeddingClient{
		client: client,
	}
}

func (e *embeddingClient) CreateEmbedding(ctx context.Context, req domain.EmbeddingRequest) (*domain.EmbeddingResult, error) {
	input := openai.EmbeddingNewParamsInputUnion{
		OfString: openai.Opt(req.Data),
	}
	res, err := e.client.Embeddings.New(ctx, openai.EmbeddingNewParams{
		Model: openai.EmbeddingModelTextEmbeddingAda002,
		Input: input,
	})
	if err != nil {
		return nil, fmt.Errorf("create embedding failed: %w", err)
	}

	result := domain.EmbeddingResult{
		Embedding: res.Data[0].Embedding,
	}

	return &result, nil
}
