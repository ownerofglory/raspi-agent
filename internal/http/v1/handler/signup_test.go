package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ownerofglory/raspi-agent/internal/core/domain"
	"github.com/ownerofglory/raspi-agent/internal/core/ports"
	"go.uber.org/mock/gomock"
)

func TestHandleSignup(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockUserSrv := ports.NewMockUserService(ctrl)
	h := NewSignupHandler(mockUserSrv)

	existingEmail := "existing@example.com"
	existingUser := domain.NewLocalUser("local", existingEmail, "password", "name", "surname")
	strongPassword, _ := GenerateSecurePassword(12)
	differentPassword, _ := GenerateSecurePassword(12)
	weakPassword := GenerateWeakPassword(6)

	testCases := []struct {
		name           string
		email          string
		firstname      string
		lastname       string
		password       string
		passwordRepeat string
		statusCode     int
	}{
		{
			name:           "valid email",
			email:          "email@example.com",
			firstname:      "firstname",
			lastname:       "lastname",
			password:       strongPassword,
			passwordRepeat: strongPassword,
			statusCode:     http.StatusCreated,
		},
		//{
		//	name:           "user already exists",
		//	email:          existingEmail,
		//	firstname:      "firstname",
		//	lastname:       "lastname",
		//	password:       "Password123!",
		//	passwordRepeat: "Password123!",
		//	statusCode:     http.StatusBadRequest,
		//},
		{
			name:           "invalid email",
			email:          "wrong email",
			firstname:      "firstname",
			lastname:       "lastname",
			password:       strongPassword,
			passwordRepeat: strongPassword,
			statusCode:     http.StatusBadRequest,
		},
		{
			name:           "password too short",
			email:          "email@example.com",
			firstname:      "firstname",
			lastname:       "lastname",
			password:       weakPassword,
			passwordRepeat: weakPassword,
			statusCode:     http.StatusBadRequest,
		},
		{
			name:           "password does not match",
			email:          "email@example.com",
			firstname:      "firstname",
			lastname:       "lastname",
			password:       strongPassword,
			passwordRepeat: differentPassword,
			statusCode:     http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.email == existingEmail {
				mockUserSrv.EXPECT().GetUserByEmail(gomock.Any(), gomock.Eq(tc.email)).Return(existingUser, nil)
			} else {
				mockUserSrv.EXPECT().GetUserByEmail(gomock.Any(), gomock.Any()).Return(nil, errors.New("user not found")).AnyTimes()
				mockUserSrv.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(existingUser, nil).AnyTimes()
			}

			reqBody := localSignup{
				Email:          tc.email,
				Firstname:      tc.firstname,
				Lastname:       tc.lastname,
				Password:       tc.password,
				PasswordRepeat: tc.passwordRepeat,
			}

			payload, _ := json.Marshal(reqBody)
			req := httptest.NewRequest(http.MethodPost, PostSignupPath, bytes.NewBuffer(payload))
			rec := httptest.NewRecorder()

			h.HandleSignup(rec, req)

			res := rec.Result()
			if res.StatusCode != tc.statusCode {
				t.Fatalf("expected status %d, got %d", tc.statusCode, res.StatusCode)
			}
		})
	}
}
