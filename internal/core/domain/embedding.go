package domain

// EmbeddingRequest represents a request to generate a vector embedding
// from a given input string.
type EmbeddingRequest struct {
	// Data is the raw text input to be converted into an embedding vector.
	Data string
}

// EmbeddingResult represents the output of an embedding generation operation.
// It contains the numerical embedding vector and information about the model used.
type EmbeddingResult struct {
	// Embedding is the numerical vector representation of the input text.
	Embedding []float64

	// Model is the name or identifier of the model used to generate the embedding.
	Model string
}
