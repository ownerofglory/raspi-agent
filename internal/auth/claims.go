package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// UserClaims represents the JWT claims used by this application.
//
// It embeds jwt.RegisteredClaims to remain compatible with the standard
// fields defined in RFC 7519 (such as exp, iat, iss), while also adding
// custom user-related fields like ID and Email.
//
// This type implements jwt.Claims by providing methods such as
// GetExpirationTime, GetIssuedAt, etc., which the JWT library uses
// during parsing and validation.
type UserClaims struct {
	// ID is the application's internal user identifier (e.g., UUID).
	ID string `json:"id"`

	// Email is the user's primary email address.
	Email string `json:"email"`

	// Issuer identifies the service that issued the token.
	Issuer string `json:"Issuer"`

	// Expires is the time at which the token should no longer be accepted.
	Expires time.Time `json:"expiresAt"`

	// Issued is the time at which the token was generated.
	Issued time.Time `json:"issuedAt"`

	// RegisteredClaims ensures compatibility with standard JWT validation,
	// e.g. checking exp, nbf, iss, etc.
	jwt.RegisteredClaims
}

// NewUserClaims creates a new UserClaims instance with the given user ID,
// email, and Issuer. It automatically sets the issued-at time to now and
// the expiration time based on the default tokenValidity duration.
//
// This helper is typically used when generating new JWTs for authenticated
// users.
func NewUserClaims(id string, email string, issuer string) *UserClaims {
	return &UserClaims{
		ID:      id,
		Email:   email,
		Issuer:  issuer,
		Issued:  time.Now(),
		Expires: time.Now().Add(tokenValidity * time.Second),
	}
}

// GetExpirationTime returns the token's expiration time in NumericDate format.
// Required by the jwt.Claims interface.
func (u UserClaims) GetExpirationTime() (*jwt.NumericDate, error) {
	return u.ExpiresAt, nil
}

// GetIssuedAt returns the token's issued-at time in NumericDate format.
// Required by the jwt.Claims interface.
func (u UserClaims) GetIssuedAt() (*jwt.NumericDate, error) {
	return u.IssuedAt, nil
}

// GetNotBefore returns the "not before" claim, which is not used in this application.
// Always returns nil.
func (u UserClaims) GetNotBefore() (*jwt.NumericDate, error) {
	return nil, nil
}

// GetIssuer returns the Issuer of the token.
// Required by the jwt.Claims interface.
func (u UserClaims) GetIssuer() (string, error) {
	return u.Issuer, nil
}

// GetSubject returns the subject of the token, which in this application
// corresponds to the user's email.
// Required by the jwt.Claims interface.
func (u UserClaims) GetSubject() (string, error) {
	return u.Email, nil
}

// GetAudience returns the token's audience claim, which is not used here.
// Always returns an empty list.
// Required by the jwt.Claims interface.
func (u UserClaims) GetAudience() (jwt.ClaimStrings, error) {
	return []string{}, nil
}
