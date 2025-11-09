package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/ownerofglory/raspi-agent/internal/auth"
	"github.com/ownerofglory/raspi-agent/internal/core/domain"
	"github.com/ownerofglory/raspi-agent/internal/core/ports"
	"golang.org/x/oauth2"
)

const googleAPIsURL = "https://www.googleapis.com/oauth2/v2"

// googleUserInfo represents the JSON response from Google's
// OAuth2 v2 `userinfo` endpoint.
type googleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
	Link          string `json:"link"`
	HD            string `json:"hd"`
}

// googleOAuth2Handler handles Google OAuth2 login and callback flows.
//
// It is responsible for redirecting the user to Google's login page,
// exchanging the auth code for a token, fetching user info from Google,
// persisting the user in the system, and returning a signed JWT.
type googleOAuth2Handler struct {
	cfg         *oauth2.Config
	jwtKey      []byte
	userService ports.UserService
}

// NewGoogleOAuth2Handler creates a new Google OAuth2 handler with the given
// OAuth2 configuration, JWT signing key, and UserService for persistence.
func NewGoogleOAuth2Handler(cfg *oauth2.Config, jwtKey []byte, userService ports.UserService) *googleOAuth2Handler {
	return &googleOAuth2Handler{
		cfg:         cfg,
		jwtKey:      jwtKey,
		userService: userService,
	}
}

// HandleLogin initiates the OAuth2 flow by redirecting the user to Google's
// authorization page with a state parameter.
//
// The state value is currently hardcoded and should be randomized and stored
// (e.g. in a cookie or session) to protect against CSRF attacks.
func (h *googleOAuth2Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	url := h.cfg.AuthCodeURL(state)

	http.SetCookie(w, &http.Cookie{
		Name:     "google-auth",
		Value:    state,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	})
	http.Redirect(w, r, url, http.StatusSeeOther)
}

// HandleOAuth2Callback completes the OAuth2 flow after Google redirects back.
//
// Steps:
//  1. Exchange the authorization code for an access token.
//  2. Fetch the user's profile from Google's API.
//  3. Upsert the user in the system (create if new, fetch if existing).
//  4. Generate a signed JWT for the user.
//  5. Respond with a login success payload (JSON).
func (h *googleOAuth2Handler) HandleOAuth2Callback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 1. Exchange auth code for token
	authCode := r.URL.Query().Get("code")
	token, err := h.cfg.Exchange(ctx, authCode)
	if err != nil {
		slog.Error("google oauth2 exchange failed", "error", err)
		http.Error(w, "failed to exchange auth code", http.StatusInternalServerError)
		return
	}

	// 2. Fetch user profile from Google
	userInfo, err := fetchGoogleUserInfo(token.AccessToken)
	if err != nil {
		slog.Error("failed to fetch google user info", "error", err)
		http.Error(w, "failed to fetch user info", http.StatusInternalServerError)
		return
	}

	// 3. Upsert user
	user, err := h.findOrCreateUser(ctx, userInfo)
	if err != nil {
		slog.Error("failed to persist user", "error", err)
		http.Error(w, "failed to persist user", http.StatusInternalServerError)
		return
	}

	// 4. Generate JWT
	claims := auth.NewUserClaims(user.ID(), user.Email(), auth.Issuer)
	jwt, err := auth.GenerateJWT(h.jwtKey, claims)
	if err != nil {
		slog.Error("failed to generate jwt", "error", err)
		http.Error(w, "failed to generate jwt", http.StatusInternalServerError)
		return
	}

	// 5. Respond with login success
	payload, err := json.Marshal(loginSuccess{ID: user.ID(), Token: jwt})
	if err != nil {
		slog.Error("failed to marshal login response", "error", err)
		http.Error(w, "failed to serialize response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(payload)
}

// fetchGoogleUserInfo retrieves the user profile from Google's UserInfo API.
func fetchGoogleUserInfo(accessToken string) (*googleUserInfo, error) {
	resp, err := http.Get(googleAPIsURL + "/userinfo?access_token=" + accessToken)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var info googleUserInfo
	if err := json.Unmarshal(data, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

// findOrCreateUser retrieves an existing user by email, or creates one if not found.
func (h *googleOAuth2Handler) findOrCreateUser(ctx context.Context, info *googleUserInfo) (domain.User, error) {
	user, err := h.userService.GetUserByEmail(ctx, info.Email)
	if err != nil {
		if errors.Is(err, domain.UserNotFound) {
			userUUID, err := uuid.NewV7()
			if err != nil {
				return nil, err
			}
			newUser := domain.NewGoogleUser(userUUID.String(), info.Email, info.GivenName, info.FamilyName)
			return h.userService.CreateUser(ctx, newUser)
		}
		return nil, err
	}
	return user, nil
}
