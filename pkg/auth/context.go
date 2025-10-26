package auth

import "context"

// authContext is a concrete implementation of Context.
// It wraps a standard context.Context and provides access
// to the authenticated UserPrincipal.
type authContext struct {
	context.Context
}

// contextKey is a private type used to avoid collisions
// when storing values in context.
type contextKey string

const UserPrincipalKey contextKey = "user-principal"

// GetUserPrincipal retrieves the UserPrincipal from the context.
// If no principal is stored or if the type assertion fails,
// it returns nil.
func (a authContext) GetUserPrincipal() UserPrincipal {
	up, ok := a.Value(UserPrincipalKey).(UserPrincipal)
	if !ok {
		return nil
	}
	return up
}

// NewAuthContext wraps an existing context.Context into an authContext.
// This allows user identity information to be propagated alongside
// deadlines, cancellation signals, and other context values.
func NewAuthContext(ctx context.Context) Context {
	return &authContext{ctx}
}

// WithUserPrincipal returns a new Context containing the given UserPrincipal.
// It can be used in authentication middleware to attach the authenticated
// user to the request context.
func WithUserPrincipal(ctx context.Context, up UserPrincipal) Context {
	return &authContext{
		context.WithValue(ctx, UserPrincipalKey, up),
	}
}
