package auth

import (
	"errors"
	"net/http"
)

// WithRoles creates an AuthorizationFunc that checks whether the
// authenticated user has at least one of the required roles.
//
// If the user is not authenticated or does not have any of the
// specified roles, the request is rejected with an error.
//
// Example:
//
//	authz := WithRoles("admin", "moderator")
//	if err := authz(rw, r); err != nil {
//		http.Error(rw, err.Error(), http.StatusForbidden)
//		return
//	}
func WithRoles(roles ...string) AuthorizationFunc {
	return func(rw http.ResponseWriter, r *http.Request) error {
		aCtx := NewAuthContext(r.Context())
		up := aCtx.GetUserPrincipal()
		if up == nil {
			return errors.New("no authenticated user in context")
		}

		userRoles := up.Roles()
		for _, required := range roles {
			for _, actual := range userRoles {
				if required == actual {
					return nil // authorized
				}
			}
		}

		return errors.New("forbidden: missing required role")
	}
}

// WithUserId creates an AuthorizationFunc that ensures the user ID
// in the request path matches the authenticated user's ID.
//
// The `param` argument specifies which path parameter contains the user ID.
// If the user is not authenticated, the parameter is missing, or the IDs
// do not match, the request is rejected.
//
// Example:
//
//	// Ensures the path /users/{id} matches the authenticated user
//	authz := WithUserId("id")
//	if err := authz(rw, r); err != nil {
//		http.Error(rw, err.Error(), http.StatusForbidden)
//		return
//	}
func WithUserId(param string) AuthorizationFunc {
	return func(rw http.ResponseWriter, r *http.Request) error {
		aCtx := NewAuthContext(r.Context())
		up := aCtx.GetUserPrincipal()
		if up == nil {
			return errors.New("no authenticated user in context")
		}

		userId := r.PathValue(param)
		if userId == "" {
			return errors.New("missing user id in path")
		}

		if userId != up.ID() {
			return errors.New("forbidden: user id mismatch")
		}
		return nil
	}
}

// WithPathParam returns an AuthorizationFunc that validates a path parameter
// using custom logic provided by the caller.
//
// The `param` argument specifies the name of the path parameter to extract
// from the request. The `valueResolver` function receives both the actual
// parameter value and the authenticated UserPrincipal. It should return nil
// if the user is authorized to access the resource, or an error if not.
//
// If the user is not authenticated, the path parameter is missing, or
// the resolver returns an error, the request will be rejected.
//
// Example:
//
//	WithPathParam("chatId", func(value string, principal UserPrincipal) error {
//		chat, err := chatService.getChatById(value)
//		if err != nil {
//			return fmt.Errorf("chat not found: %w", err)
//		}
//		if chat.UserId != principal.ID() {
//			return errors.New("forbidden: not your chat")
//		}
//		return nil
//	})
func WithPathParam(param string, valueResolver func(value string, principal UserPrincipal) error) AuthorizationFunc {
	return func(rw http.ResponseWriter, r *http.Request) error {
		aCtx := NewAuthContext(r.Context())
		up := aCtx.GetUserPrincipal()
		if up == nil {
			return errors.New("no authenticated user in context")
		}

		paramVal := r.PathValue(param)
		if paramVal == "" {
			return errors.New("missing param id in path: " + param)
		}

		err := valueResolver(paramVal, up)
		if err != nil {
			return errors.New("forbidden: " + err.Error())
		}
		return nil
	}
}

// WithPrincipal creates an AuthorizationFunc that applies custom
// authorization logic based on the authenticated UserPrincipal.
//
// The principalResolver is given the UserPrincipal from the request
// context. It should return nil if the user is authorized, or an error
// if access should be denied.
//
// If no authenticated user is present in the context, or the resolver
// returns an error, the request is rejected.
//
// Example:
//
//	// Allow only if the user has a valid tariff for weather API usage
//	authz := WithPrincipal(func(principal UserPrincipal) error {
//		if !principal.HasTariff("premium") {
//			return errors.New("tariff does not allow weather API")
//		}
//		return nil
//	})
func WithPrincipal(principalResolver func(principal UserPrincipal) error) AuthorizationFunc {
	return func(rw http.ResponseWriter, r *http.Request) error {
		aCtx := NewAuthContext(r.Context())
		up := aCtx.GetUserPrincipal()
		if up == nil {
			return errors.New("no authenticated user in context")
		}

		err := principalResolver(up)
		if err != nil {
			return errors.New("forbidden: " + err.Error())
		}
		return nil
	}
}
