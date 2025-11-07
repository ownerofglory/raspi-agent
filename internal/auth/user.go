package auth

// userPrincipal is a concrete implementation of the UserPrincipal interface
// that wraps parsed JWT claims (UserClaims).
//
// It adapts the application-specific UserClaims into the generic
// UserPrincipal interface expected by the rest of the authentication/authorization
// layer.
type userPrincipal struct {
	*UserClaims
}

// ID returns the unique identifier of the user, taken from the JWT claims.
func (u userPrincipal) ID() string {
	return u.UserClaims.ID
}

// Email returns the user's email address, taken from the JWT claims.
func (u userPrincipal) Email() string {
	return u.UserClaims.Email
}

// Roles returns the roles assigned to the user.
// Currently this is hardcoded to ["ROLE_USER"], but in a real-world scenario
// you would likely read roles from the JWT claims or from a user store.
func (u userPrincipal) Roles() []string {
	return []string{"ROLE_USER"}
}

// NewUserPrincipal creates a new UserPrincipal implementation from the given
// JWT claims. This function is typically called after successfully parsing
// a JWT in WithJWT authentication middleware.
func NewUserPrincipal(uc *UserClaims) *userPrincipal {
	return &userPrincipal{uc}
}
