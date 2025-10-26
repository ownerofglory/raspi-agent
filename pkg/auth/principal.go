package auth

// anonymousPrincipal is a UserPrincipal implementation that represents
// an unauthenticated (guest) user. It is useful in cases where you want
// to avoid nil checks but still differentiate between authenticated and
// unauthenticated requests.
type anonymousPrincipal struct{}

// ID always returns an empty string for anonymous users.
func (a anonymousPrincipal) ID() string {
	return ""
}

// Email always returns an empty string for anonymous users.
func (a anonymousPrincipal) Email() string {
	return ""
}

// Roles always returns an empty slice, since anonymous users have no roles.
func (a anonymousPrincipal) Roles() []string {
	return []string{}
}

// NewAnonymousPrincipal constructs a new anonymousPrincipal as a UserPrincipal.
func NewAnonymousPrincipal() UserPrincipal {
	return anonymousPrincipal{}
}
