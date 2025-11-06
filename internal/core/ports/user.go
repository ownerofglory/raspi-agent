package ports

import (
	"context"

	"github.com/ownerofglory/raspi-agent/internal/core/domain"
)

// UserService defines the core business logic operations related to users.
//
// This interface lives in the **core** (domain/application) layer and expresses
// what the system can *do* regarding users — independent of delivery or
// persistence mechanisms.
type UserService interface {
	// GetUser retrieves a single user by ID.
	// Implementations should return a domain.NotFoundError if no user exists.
	GetUser(ctx context.Context, id string) (domain.User, error)
}

// UserRepo defines the persistence interface for user records.
//
// Implementations of this interface handle actual data storage and retrieval,
// whether via SQL, NoSQL, in-memory, or external APIs.
//
// Typical usage:
//   - `Find` is used when you know a user’s ID (e.g., for profile lookup).
//   - `FindByEmail` supports authentication or email-based identification.
//   - `Save` handles both create and update operations.
//   - `Delete` removes the user (or marks them as inactive) in the datastore.
type UserRepo interface {
	// Find retrieves a user by their unique identifier.
	// Returns (domain.User, nil) if found, or an error otherwise.
	Find(ctx context.Context, id string) (domain.User, error)

	// FindByEmail retrieves a user by their email address.
	// Returns (domain.User, nil) if found, or an error otherwise.
	FindByEmail(ctx context.Context, email string) (domain.User, error)

	// Save persists a new or existing user to the repository.
	// Implementations should handle insert vs. update logic internally.
	Save(ctx context.Context, user domain.User) (*domain.User, error)

	// Delete permanently removes a user by ID.
	// Depending on business rules, this may perform a soft-delete instead.
	Delete(ctx context.Context, id string) error
}
