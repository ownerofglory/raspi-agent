package tools

import (
	"context"
	"fmt"
	"sync"
)

// Schema defines an open-typed metadata description for a tool.
// Example: parameters, input validation, expected outputs.
type Schema map[string]any

// Tool defines the contract for any executable tool or command.
// Tools can represent system utilities, API calls, or internal functions.
type Tool interface {
	// Name returns the unique tool identifier.
	Name() string

	// Description provides a short human-readable summary.
	Description() string

	// UserMessage is a friendly prompt for a human user.
	UserMessage() string

	// Schema returns a structured definition of accepted arguments.
	Schema() Schema

	// Execute runs the tool with the given args (usually JSON or CLI-style string).
	// Returns arbitrary data or an error.
	Execute(ctx context.Context, args string) (any, error)
}

// Tools is a thread-safe registry of all available Tool instances.
type Tools struct {
	mu    sync.RWMutex
	store map[string]Tool
}

// New creates a new empty tool registry.
func New() *Tools {
	return &Tools{
		store: make(map[string]Tool),
	}
}

// Add registers a new tool by its Name.
// It panics if a duplicate name is added — ensuring global uniqueness.
func (t *Tools) Add(tool Tool) {
	t.mu.Lock()
	defer t.mu.Unlock()

	name := tool.Name()
	if _, exists := t.store[name]; exists {
		panic(fmt.Sprintf("tool with name '%s' already registered", name))
	}
	t.store[name] = tool
}

// Get retrieves a tool by name.
func (t *Tools) Get(name string) (Tool, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	tool, ok := t.store[name]
	if !ok {
		return nil, fmt.Errorf("tool '%s' not found", name)
	}
	return tool, nil
}

// ExecuteByName runs a tool by name with the given args.
func (t *Tools) ExecuteByName(ctx context.Context, name, args string) (any, error) {
	tool, err := t.Get(name)
	if err != nil {
		return nil, err
	}
	return tool.Execute(ctx, args)
}

// List returns all registered tools with metadata.
func (t *Tools) List() []map[string]string {
	t.mu.RLock()
	defer t.mu.RUnlock()

	list := make([]map[string]string, 0, len(t.store))
	for _, tool := range t.store {
		list = append(list, map[string]string{
			"name":        tool.Name(),
			"description": tool.Description(),
			"userMessage": tool.UserMessage(),
		})
	}
	return list
}
