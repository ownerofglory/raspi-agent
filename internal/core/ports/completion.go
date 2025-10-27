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
	CreateCompletion(ctx context.Context, req *domain.CompletionRequest) (*domain.CompletionResult, error)
}

// SummaryProvider defines the interface for generating summaries of text or conversations.
// Implementations of this interface produce concise summaries based on a series of input messages.
type SummaryProvider interface {
	// CreateSummary generates a summary for the provided conversation or text input.
	// It takes a context for cancellation and timeout control, along with a summary request
	// containing the messages to summarize. It returns the generated summary result or an error.
	CreateSummary(ctx context.Context, req *domain.SummaryRequest) (*domain.SummaryResult, error)
}
