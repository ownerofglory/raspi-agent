package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strings"

	z "github.com/Oudwins/zog"
	"github.com/Oudwins/zog/internals"
	"github.com/google/uuid"
	"github.com/ownerofglory/raspi-agent/internal/core/domain"
	"github.com/ownerofglory/raspi-agent/internal/core/ports"
	"golang.org/x/crypto/bcrypt"
)

// PostSignupPath is the URL path for login requests.
const PostSignupPath = basePath + "/auth/signup"

type localSignup struct {
	Email          string `json:"email"`
	Firstname      string `json:"firstname"`
	Lastname       string `json:"lastname"`
	Password       string `json:"password"`
	PasswordRepeat string `json:"passwordRepeat"`
}

var localSignupSchema = z.Struct(z.Shape{
	"email": z.String().
		Email(z.Message("invalid email format")).
		Min(5, z.Message("email must be at least 5 characters")).
		Max(254, z.Message("email must be at most 254 characters")),

	"firstname": z.String().
		Min(1, z.Message("firstname is required")).
		Max(100, z.Message("firstname must be at most 100 characters")),

	"lastname": z.String().
		Min(1, z.Message("lastname is required")).
		Max(100, z.Message("lastname must be at most 100 characters")),

	"password": z.String().
		Min(8, z.Message("password must be at least 8 characters")).
		Max(128, z.Message("password must be at most 128 characters")),

	"passwordRepeat": z.String().
		Min(8, z.Message("passwordRepeat must be at least 8 characters")).
		Max(128, z.Message("passwordRepeat must be at most 128 characters")),
}).TestFunc(func(val any, ctx internals.Ctx) bool {
	schema, ok := val.(*localSignup)
	if !ok {
		return false
	}
	return schema.Password == schema.PasswordRepeat
})

// signupHandler handles user signup requests.
//
// It validates the request payload, hashes the password, creates a new user,
// and persists it using the UserService.
type signupHandler struct {
	userService ports.UserService
}

// NewSignupHandler creates a new signupHandler with the given UserService.
func NewSignupHandler(userService ports.UserService) *signupHandler {
	return &signupHandler{
		userService: userService,
	}
}

// HandleSignup processes an HTTP signup request.
//
// It performs the following steps:
//  1. Reads and unmarshals the request body into a localSignup struct.
//  2. Validates the input using localSignupSchema.
//  3. Hashes the user's password with bcrypt.
//  4. Creates a new user with a generated UUID.
//  5. Persists the user via UserService.
//  6. Responds with HTTP 201 Created if successful.
//
// Errors during validation or persistence result in appropriate
// HTTP 400 or 500 responses.
func (h *signupHandler) HandleSignup(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("error reading request body", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var ls localSignup
	err = json.Unmarshal(reqBody, &ls)
	if err != nil {
		slog.Error("error unmarshalling request body", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	issues := localSignupSchema.Validate(&ls)
	if issues != nil {
		slog.Error("error validating request body", "err", issues)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	email := strings.ToLower(ls.Email)
	existingUser, err := h.userService.GetUserByEmail(r.Context(), email)
	if !errors.Is(err, domain.UserNotFound) && existingUser != nil {
		slog.Error("user already exists", "email", email)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(ls.Password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("error generating password hash", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userUUID, err := uuid.NewUUID()
	if err != nil {
		slog.Error("error creating user uuid", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user := domain.NewLocalUser(userUUID.String(), email, string(passwordHash), ls.Firstname, ls.Lastname)

	_, err = h.userService.CreateUser(r.Context(), user)
	if err != nil {
		slog.Error("error creating user", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
