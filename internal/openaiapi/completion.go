package openaiapi

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/openai/openai-go/v3"
	"github.com/ownerofglory/raspi-agent/internal/core/domain"
)

// completionClient implements the CompletionProvider interface using the
// OpenAI API as a backend.
//
// It wraps an initialized *openai.Client and exposes a high-level method
// to create chat-style completions based on a user prompt.
type completionClient struct {
	client *openai.Client
}

// NewCompletionClient creates a new instance of completionClient.
// The provided OpenAI client must be pre-configured with a valid API key.
func NewCompletionClient(client *openai.Client) *completionClient {
	return &completionClient{
		client: client,
	}
}

// CreateCompletion generates a text completion for the given prompt using
// the OpenAI Chat Completions API.
//
// It currently uses the GPT-4o-mini model for efficiency but can be
// parameterized or made configurable in future versions.
func (c *completionClient) CreateCompletion(ctx context.Context, req *domain.CompletionRequest) (*domain.CompletionResult, error) {
	messages := make([]openai.ChatCompletionMessageParamUnion, 0)

	systemMessage := openai.SystemMessage("You are an AI agent named Vicky")
	messages = append(messages, systemMessage)

	userContent := openai.TextContentPart(req.Prompt)
	userMessage := openai.UserMessage([]openai.ChatCompletionContentPartUnionParam{
		userContent,
	})

	messages = append(messages, userMessage)
	completion, err := c.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: messages,
		Model:    openai.ChatModelGPT4oMini,
	})
	if err != nil {
		slog.Error("completion create completion request fail", "err", err)
		return nil, fmt.Errorf("completion create completion request fail: %w", err)
	}

	m := completion.Choices[0].Message
	return &domain.CompletionResult{
		Text: m.Content,
	}, nil
}
