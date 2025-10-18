package ports

import (
	"context"

	"github.com/ownerofglory/raspi-agent/internal/core/domain"
)

// CompletionProvider defines the contract for any service capable of generating
// text completions in response to a given prompt.
//
// Implementations may use large language models (LLMs) — such as OpenAI GPT,
// Anthropic Claude, or local inference engines — to produce contextually
// relevant text responses based on the provided input.
type CompletionProvider interface {
	// CreateCompletion generates a text response for the given request.
	//
	// The context controls cancellation and timeout. Implementations should
	// handle retries, streaming, and rate limiting as appropriate.
	//
	// Returns:
	//   - A pointer to CompletionResult containing the generated text.
	//   - An error if the generation process fails or the request is invalid.
	CreateCompletion(ctx context.Context, req domain.CompletionRequest) (*domain.CompletionResult, error)
}
