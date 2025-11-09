package middleware

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	appAuth "github.com/ownerofglory/raspi-agent/internal/auth"
	"github.com/ownerofglory/raspi-agent/pkg/auth"
)

// Authenticated wraps an HTTP handler with an authentication check.
//
// It accepts an AuthenticationFunc, which inspects the incoming request
// (e.g. Authorization header, cookies, etc.), validates credentials, and
// returns an auth.Context containing a UserPrincipal.
//
// If authentication fails, the middleware writes a 401 Unauthorized
// response and does not call the next handler.
func Authenticated(af auth.AuthenticationFunc) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authCtx, err := af(w, r)
			if err != nil {
				slog.Error("Authentication error: ", "error", err)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			rAuth := r.WithContext(authCtx)

			h.ServeHTTP(w, rAuth)
		})
	}
}

// WithJWT returns an AuthenticationFunc that validates a JWT from the
// Authorization header in the form "Bearer <token>".
//
// If the token is valid, it is parsed into claims and converted into a
// UserPrincipal, which is stored in the request context.
// If the header is missing or the token is invalid, an error is returned.
func WithJWT(key string) auth.AuthenticationFunc {
	return func(rw http.ResponseWriter, r *http.Request) (auth.Context, error) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			slog.Error("Authorization header is empty")
			return nil, errors.New("Authorization header is empty")
		}

		bearerToken, found := strings.CutPrefix(authHeader, "Bearer ")
		if !found {
			slog.Error("Authorization header is invalid")
			return nil, errors.New("Authorization header is invalid")
		}

		claims, err := appAuth.ParseJWT(bearerToken, []byte(key))
		if err != nil {
			return nil, err
		}

		up := appAuth.NewUserPrincipal(claims)
		ctxWithAuth := context.WithValue(r.Context(), auth.UserPrincipalKey, up)

		authContext := auth.NewAuthContext(ctxWithAuth)

		return authContext, nil
	}
}

// Authorized wraps an HTTP handler with one or more authorization checks.
//
// Each AuthorizationFunc is invoked in order. If any of them return an error,
// the middleware writes a 403 Forbidden response and does not call the next
// handler.
//
// Typical AuthorizationFuncs include checks for user roles or matching user IDs.
func Authorized(authFuncs ...auth.AuthorizationFunc) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for _, fn := range authFuncs {
				if err := fn(w, r); err != nil {
					slog.Error("Authorization error", "error", err)
					w.WriteHeader(http.StatusForbidden)
					return
				}
			}

			h.ServeHTTP(w, r)
		})
	}
}
