package domain

// RecallRequest represents a request to recall stored information or memories
// associated with a specific user.
//
// UserID identifies the user whose memories are being requested.
// Text optionally contains a query or context string that narrows down
// which memories should be recalled.
type RecallRequest struct {
	UserID string // unique identifier of the user
	Text   string // optional query text for contextual recall
}

// RecallResult holds the response to a RecallRequest.
//
// Memories is a list of memory strings that match or relate to the recall query.
type RecallResult struct {
	Memories []string // recalled memory entries
}
