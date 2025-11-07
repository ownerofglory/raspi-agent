package auth

import (
	"testing"
	"time"
)

func TestGenerateJWT(t *testing.T) {
	key := []byte("a-test-secret-at-least-256-bits-long-")
	now := time.Now()

	tests := []struct {
		name    string
		claims  UserClaims
		wantErr bool
	}{
		{
			name: "minimal claims",
			claims: UserClaims{
				ID:    "123",
				Email: "user@example.com",
			},
			wantErr: false,
		},
		{
			name: "with explicit times",
			claims: UserClaims{
				ID:      "456",
				Email:   "withtimes@example.com",
				Issued:  now,
				Expires: now.Add(1 * time.Hour),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GenerateJWT(key, &tt.claims)
			if (err != nil) != tt.wantErr {
				t.Fatalf("GenerateJWT() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && token == "" {
				t.Error("GenerateJWT() returned empty token")
			}
		})
	}
}

func TestParseJWT(t *testing.T) {
	key := []byte("a-test-secret-at-least-256-bits-long")
	now := time.Now()

	validClaims := UserClaims{ID: "abc", Email: "valid@example.com"}
	validToken, _ := GenerateJWT(key, &validClaims)

	expiredClaims := UserClaims{
		ID:      "expired",
		Email:   "expired@example.com",
		Issued:  now.Add(-2 * time.Hour),
		Expires: now.Add(-1 * time.Hour),
	}
	expiredToken, _ := GenerateJWT(key, &expiredClaims)

	tests := []struct {
		name     string
		token    string
		parseKey []byte
		wantErr  bool
		wantID   string
	}{
		{
			name:     "valid token",
			token:    validToken,
			parseKey: key,
			wantErr:  false,
			wantID:   "abc",
		},
		{
			name:     "invalid signature",
			token:    validToken,
			parseKey: []byte("wrong-key"),
			wantErr:  true,
		},
		{
			name:     "expired token",
			token:    expiredToken,
			parseKey: key,
			wantErr:  true,
		},
		{
			name:     "malformed token",
			token:    "not-a-jwt",
			parseKey: key,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsed, err := ParseJWT(tt.token, tt.parseKey)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseJWT() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && parsed.ID != tt.wantID {
				t.Errorf("ParseJWT() ID = %v, want %v", parsed.ID, tt.wantID)
			}
		})
	}
}
