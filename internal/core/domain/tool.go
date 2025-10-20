package domain

import "context"

type AgentTool interface {
	Execute(ctx context.Context, args string) (any, error)
	Schema() ToolSchema
	Name() string
	Description() string
	UserMessage() string
}

type ToolSchema map[string]any

type ToolDef struct {
	Schema ToolSchema
	Tool   func(context.Context, string) (any, error)
}

func (t ToolDef) Execute(ctx context.Context, args string) (any, error) {
	return t.Tool(ctx, args)
}
