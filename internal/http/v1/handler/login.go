package handler

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strings"

	z "github.com/Oudwins/zog"
	"github.com/ownerofglory/raspi-agent/internal/auth"
	"github.com/ownerofglory/raspi-agent/internal/core/domain"
	"github.com/ownerofglory/raspi-agent/internal/core/ports"
	"golang.org/x/crypto/bcrypt"
)

// PostLoginPath is the URL path for login requests.
const PostLoginPath = basePath + "/login"

// localLogin represents the expected JSON payload for login requests.
//
// Example:
//
//	{
//	  "email": "user@example.com",
//	  "password": "secret"
//	}
type localLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

var localLoginSchema = z.Struct(z.Shape{
	"email": z.String().
		Email(z.Message("invalid email format")).
		Min(5, z.Message("email must be at least 5 characters")).
		Max(254, z.Message("email must be at most 254 characters")),

	"password": z.String().
		Min(8, z.Message("password must be at least 8 characters")).
		Max(128, z.Message("password must be at most 128 characters")),
})

// loginSuccess is the JSON response returned on successful login.
//
// Example:
//
//	{
//	  "id": "user-uuid",
//	  "token": "jwt-token"
//	}
type loginSuccess struct {
	ID    string `json:"id"`
	Token string `json:"token"`
}

// loginHandler implements http.Handler for the login endpoint.
// It validates the request, creates a JWT, and returns it to the client.
type loginHandler struct {
	jwtKey      string
	userService ports.UserService
}

// NewLoginHandler constructs a new loginHandler with the given JWT signing key.
func NewLoginHandler(jwtKey string, userService ports.UserService) *loginHandler {
	return &loginHandler{
		jwtKey:      jwtKey,
		userService: userService,
	}
}

// HandleLogin handles login requests.
// It reads the JSON payload, validates credentials (placeholder),
// generates a new JWT for the user, and writes the token as JSON response.
func (h *loginHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	// Read request body
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("Authentication error", "error", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Parse JSON credentials
	var creds localLogin
	if err := json.Unmarshal(payload, &creds); err != nil {
		slog.Error("Authentication error", "error", err)
		http.Error(w, "invalid JSON payload", http.StatusBadRequest)
		return
	}

	issues := localLoginSchema.Validate(&creds)
	if issues != nil {
		slog.Error("Authentication error", "error", issues)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	email := strings.ToLower(creds.Email)
	user, err := h.userService.GetUserByEmail(r.Context(), email)
	if err != nil {
		slog.Error("error when getting user by email", "email", email, "err", err)
		http.Error(w, "wrong username or password", http.StatusNotFound)
		return
	}

	if user.Provider() != domain.LocalProvider {
		slog.Error("User created using different provider", "email", email)
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(*user.Password()), []byte(creds.Password))
	if err != nil {
		slog.Error("Authentication error. Wrong password: ", "error", err)
		http.Error(w, "wrong username or password", http.StatusNotFound)
		return
	}

	// Create claims and sign a JWT
	claims := auth.NewUserClaims(user.ID(), user.Provider(), auth.Issuer)
	jwtToken, err := auth.GenerateJWT([]byte(h.jwtKey), claims)
	if err != nil {
		slog.Error("Authentication error", "error", err)
		http.Error(w, "could not generate token", http.StatusInternalServerError)
		return
	}

	// Build success response
	ls := loginSuccess{
		ID:    user.ID(),
		Token: jwtToken,
	}

	respPayload, err := json.Marshal(ls)
	if err != nil {
		slog.Error("Authentication error", "error", err)
		http.Error(w, "could not marshal response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(respPayload)
}
