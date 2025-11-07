package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"log/slog"
	"time"
)

const (
	// tokenValidity defines how long a generated token is valid (24h).
	tokenValidity = 60 * 60 * 24
	// Issuer identifies this service as the JWT Issuer.
	Issuer = "raspi-agent"
)

// GenerateJWT creates and signs a new JWT for the given user claims.
//
// The token is signed using HS256 with the provided key.
// Standard claims (iss, sub, exp, iat) are included in the payload.
// The returned string is the compact serialized JWT.
func GenerateJWT(key []byte, claims *UserClaims) (string, error) {
	if claims.Expires.IsZero() {
		claims.Expires = time.Now().Add(tokenValidity * time.Second)
	}
	if claims.Issued.IsZero() {
		claims.Issued = time.Now()
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    claims.ID,
		"email": claims.Email,
		"iss":   Issuer,
		"sub":   claims.ID,
		"exp":   claims.Expires.Unix(),
		"iat":   claims.Issued.Unix(),
	})

	signed, err := token.SignedString(key)
	if err != nil {
		slog.Error("Failed to sign JWT", "error", err)
		return "", err
	}
	return signed, nil
}

// ParseJWT parses and validates a JWT string using the given signing key.
//
// It expects the token to be signed with HS256 and to contain UserClaims.
// If the token is valid and signature matches, the claims are returned.
// Otherwise, an error is logged and returned.
func ParseJWT(ts string, key []byte) (*UserClaims, error) {
	uc := &UserClaims{}
	token, err := jwt.ParseWithClaims(ts, uc, func(t *jwt.Token) (any, error) {
		return key, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		slog.Error("Failed to parse JWT", "error", err)
		return nil, err
	}

	if token.Valid {
		return uc, nil
	}

	slog.Error("Invalid JWT", "token", ts)
	return nil, err
}
