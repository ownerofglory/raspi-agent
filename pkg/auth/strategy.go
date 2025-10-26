package auth

import "net/http"

// WithAnonymous creates an AuthenticationFunc that always attaches
// an anonymousPrincipal to the request context. Useful when you want
// to allow unauthenticated access but still keep the Context interface
// consistent for downstream code.
func WithAnonymous() AuthenticationFunc {
	return func(rw http.ResponseWriter, r *http.Request) (Context, error) {
		// Start with the base context
		authCtx := WithUserPrincipal(r.Context(), NewAnonymousPrincipal())
		return authCtx, nil
	}
}
