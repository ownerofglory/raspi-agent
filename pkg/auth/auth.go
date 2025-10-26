package auth

import (
	"context"
	"net/http"
)

// UserPrincipal represents the identity of an authenticated user.
//
// Implementations of this interface should encapsulate the essential
// information about a user that is relevant for authentication and
// authorization. Typically, this includes a stable identifier (e.g., user ID),
// their email address, and the roles or permissions assigned to them.
//
// The values returned here are usually extracted from a verified JWT, a
// database lookup, or another identity provider.
type UserPrincipal interface {
	// ID returns the unique identifier of the user (e.g., UUID, database ID).
	ID() string

	// Email returns the user's primary email address.
	Email() string

	// Roles returns a list of roles assigned to the user.
	// Roles are commonly used for coarse-grained authorization checks.
	Roles() []string
}

// Context extends the standard context.Context with user identity information.
//
// It allows handlers and middleware deeper in the request pipeline to access
// the authenticated user's principal alongside the usual context values
// (cancellation, deadlines, etc.).
type Context interface {
	context.Context

	// GetUserPrincipal returns the authenticated user principal stored in
	// the context. If the request is not authenticated, implementations
	// may return nil.
	GetUserPrincipal() UserPrincipal
}

type (
	// AuthenticationFunc represents a function responsible for authenticating
	// an HTTP request.
	//
	// Implementations should inspect the incoming request (e.g., headers,
	// cookies, tokens), validate the authentication information, and return
	// a Context containing the associated UserPrincipal if successful.
	//
	// If authentication fails, an error should be returned, and the
	// implementation may also write an appropriate HTTP response.
	AuthenticationFunc func(rw http.ResponseWriter, r *http.Request) (Context, error)

	// AuthorizationFunc represents a function responsible for authorizing
	// an HTTP request.
	//
	// Implementations should enforce access control based on the current
	// user's identity (obtained from the Context) and the requested resource.
	// If authorization fails, an error should be returned, and the
	// implementation may also write an appropriate HTTP response.
	AuthorizationFunc func(rw http.ResponseWriter, r *http.Request) error
)
