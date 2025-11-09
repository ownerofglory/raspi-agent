package handler

import (
	"log/slog"
	"net/http"
)

// PostAuthOAuth2LoginPath is the URL path for initiating an OAuth2 login.
//
// Example: /api/v1/auth/oauth2/google/login
const PostAuthOAuth2LoginPath = basePath + "/auth/oauth2/{provider}/login"

// PostAuthOAuth2CallbackPath is the URL path for handling an OAuth2 callback.
//
// Example: /api/v1/auth/oauth2/google/callback
const PostAuthOAuth2CallbackPath = basePath + "/auth/oauth2/{provider}/callback"

// oauth2Handler routes OAuth2 login and callback requests to the correct
// provider-specific implementation.
//
// Currently only Google is supported, but this struct provides a simple
// extension point for adding more providers (GitHub, Apple, etc.).
type oauth2Handler struct {
	google *googleOAuth2Handler
}

// NewOAuth2Handler creates a new oauth2Handler with the given Google handler.
//
// Other providers can be added to this constructor in the future.
func NewOAuth2Handler(google *googleOAuth2Handler) *oauth2Handler {
	return &oauth2Handler{
		google: google,
	}
}

// HandleLogin starts the OAuth2 login flow for the given provider.
//
// It inspects the {provider} path parameter and delegates to the correct
// handler. If the provider is not recognized, it returns a 400 Bad Request.
func (h *oauth2Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	providerName := r.PathValue("provider")
	if providerName == "google" {
		h.google.HandleLogin(w, r)
		return
	}

	slog.Error("oauth2 handler error: invalid provider", "provider", providerName)
	w.WriteHeader(http.StatusBadRequest)
}

// HandleCallback completes the OAuth2 login flow for the given provider.
//
// It inspects the {provider} path parameter and delegates to the correct
// handler. If the provider is not recognized, it returns a 400 Bad Request.
func (h *oauth2Handler) HandleCallback(w http.ResponseWriter, r *http.Request) {
	providerName := r.PathValue("provider")
	if providerName == "google" {
		h.google.HandleOAuth2Callback(w, r)
		return
	}
	slog.Error("oauth2 handler error: invalid provider", "provider", providerName)
	w.WriteHeader(http.StatusBadRequest)
}
