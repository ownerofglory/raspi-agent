package handler

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"errors"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"

	mathRand "math/rand"

	"github.com/ownerofglory/raspi-agent/internal/core/domain"
	"github.com/ownerofglory/raspi-agent/internal/core/ports"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func TestHandleLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockUserSrv := ports.NewMockUserService(ctrl)
	key := "a-test-secret-at-least-256-bits-long-"
	h := NewLoginHandler(key, mockUserSrv)

	existingEmail := "test@test.com"
	nonExistingEmail := "non-existing@test.com"
	externalEmail := "external@google.com"
	strongPassword, _ := GenerateSecurePassword(12)
	differentPassword, _ := GenerateSecurePassword(12)
	password, _ := bcrypt.GenerateFromPassword([]byte(strongPassword), bcrypt.DefaultCost)

	existingUser := domain.NewLocalUser("local", existingEmail, string(password), "test", "text")
	nonExistingUser := domain.NewLocalUser("local", nonExistingEmail, string(password), "non", "non")
	externalUser := domain.NewGoogleUser("google", externalEmail, "first", "last")

	testCases := []struct {
		name       string
		email      string
		password   string
		statusCode int
	}{
		{
			name:       "success email",
			email:      existingEmail,
			password:   strongPassword,
			statusCode: http.StatusOK,
		},
		{
			name:       "email correct, password wrong",
			email:      existingEmail,
			password:   differentPassword,
			statusCode: http.StatusNotFound,
		},
		{
			name:       "non-existing email",
			email:      nonExistingEmail,
			password:   strongPassword,
			statusCode: http.StatusNotFound,
		},
		{
			name:       "external user",
			email:      externalEmail,
			password:   string(password),
			statusCode: http.StatusNotFound,
		},
		{
			name:       "wrong email format",
			email:      "wrong-email",
			password:   string(password),
			statusCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.email == existingEmail {
				mockUserSrv.EXPECT().GetUserByEmail(gomock.Any(), existingUser.Email()).Return(existingUser, nil)
			} else if tc.email == nonExistingEmail {
				mockUserSrv.EXPECT().GetUserByEmail(gomock.Any(), nonExistingUser.Email()).Return(nil, errors.New("not found"))
			} else if tc.email == externalEmail {
				mockUserSrv.EXPECT().GetUserByEmail(gomock.Any(), externalEmail).Return(externalUser, errors.New("not found"))
			}

			reqBody := localLogin{
				Email:    tc.email,
				Password: tc.password,
			}

			payload, _ := json.Marshal(reqBody)
			req := httptest.NewRequest(http.MethodPost, PostLoginPath, bytes.NewBuffer(payload))
			rec := httptest.NewRecorder()

			h.HandleLogin(rec, req)

			res := rec.Result()
			if res.StatusCode != tc.statusCode {
				t.Fatalf("expected status %d, got %d", tc.statusCode, res.StatusCode)
			}

			if tc.statusCode == http.StatusOK {
				defer res.Body.Close()

				resData, err := io.ReadAll(res.Body)
				if err != nil {
					t.Errorf("expected error to be nil got %v", err)
				}

				var ls loginSuccess
				err = json.Unmarshal(resData, &ls)
				if err != nil {
					t.Errorf("expected error to be nil got %v", err)
				}

				if ls.Token == "" {
					t.Error("expected token to be non-empty")
				}
			}

		})
	}
}

const (
	lowerChars   = "abcdefghijklmnopqrstuvwxyz"
	upperChars   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digitChars   = "0123456789"
	specialChars = "!@#$%^&*()-_=+[]{}<>?/|"
	allChars     = lowerChars + upperChars + digitChars + specialChars
)

// GenerateSecurePassword generates a random secure password.
// It ensures at least one lowercase, one uppercase, one digit, and one special character.
func GenerateSecurePassword(length int) (string, error) {
	if length < 8 {
		length = 8 // enforce minimum length
	}

	// Ensure each category is present at least once
	categories := []string{lowerChars, upperChars, digitChars, specialChars}
	password := make([]byte, length)

	// Fill first slots with required categories
	for i, chars := range categories {
		c, err := randomChar(chars)
		if err != nil {
			return "", err
		}
		password[i] = c
	}

	// Fill remaining slots with random chars from all categories
	for i := len(categories); i < length; i++ {
		c, err := randomChar(allChars)
		if err != nil {
			return "", err
		}
		password[i] = c
	}

	// Shuffle the result so the first 4 chars are not predictable
	shuffle(password)

	return string(password), nil
}

func shuffle(data []byte) {
	for i := len(data) - 1; i > 0; i-- {
		jRand, _ := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		j := int(jRand.Int64())
		data[i], data[j] = data[j], data[i]
	}
}

func randomChar(chars string) (byte, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
	if err != nil {
		return 0, err
	}
	return chars[n.Int64()], nil
}

const weakCharset = "abcdefghijklmnopqrstuvwxyz0123456789"

// GenerateWeakPassword generates a weak, non-secure password
// using only lowercase letters and digits.
func GenerateWeakPassword(length int) string {
	if length < 4 {
		length = 4
	}
	password := make([]byte, length)
	for i := range password {
		password[i] = weakCharset[mathRand.Intn(len(weakCharset))]
	}
	return string(password)
}
