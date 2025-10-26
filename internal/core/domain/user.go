package domain

const (
	// LocalProvider identifies users that are managed locally
	// (e.g. registered directly with email + password).
	LocalProvider = "local"

	// GoogleProvider identifies users that are authenticated
	// through Google as an external identity provider.
	GoogleProvider = "google"
)

// User defines the minimal interface for representing a user model in the system.
//
// This abstraction allows your domain logic and services to depend on
// a stable contract, rather than being tied to a specific persistence
// model (e.g. database struct) or transport layer (e.g. JWT claims).
//
// Typical implementations might come from:
//   - A database model (e.g. ORM struct implementing this interface)
//   - An identity provider (mapping claims into this contract)
//   - A mock object for testing
//
// By using an interface, the core application logic remains decoupled
// from infrastructure concerns, which follows the principles of hexagonal
// architecture (ports and adapters).
type User interface {
	// ID returns the unique identifier of the user
	// (e.g. UUID or database ID).
	ID() string

	// Email returns the user's primary email address.
	Email() string

	// Firstname returns the user's given (first) name.
	Firstname() string

	// Lastname returns the user's family (last) name.
	Lastname() string

	// Password returns the hashed password if the user has one,
	// or nil if the user authenticates via an external provider.
	Password() *string

	// Provider identifies the source of this user
	// (e.g. "local", "google").
	Provider() string
}

// localUser is a User implementation backed by local storage.
// Typically created during direct user registration.
type localUser struct {
	id        string  `json:"id"`
	email     string  `json:"email"`
	password  *string `json:"password"`
	firstname string  `json:"firstname"`
	lastname  string  `json:"lastname"`
}

// googleUser is a User implementation backed by Google as
// an external identity provider.
type googleUser struct {
	id        string `json:"id"`
	email     string `json:"email"`
	firstname string `json:"firstname"`
	lastname  string `json:"lastname"`
}

// NewLocalUser creates a new localUser instance with the given fields.
//
// The password is expected to already be hashed before being passed here.
func NewLocalUser(id string, email, password, firstname, lastname string) User {
	return &localUser{
		id:        id,
		email:     email,
		password:  &password,
		firstname: firstname,
		lastname:  lastname,
	}
}

func (l *localUser) ID() string {
	return l.id
}

func (l *localUser) Email() string {
	return l.email
}

func (l *localUser) Firstname() string {
	return l.firstname
}

func (l *localUser) Lastname() string {
	return l.lastname
}

func (l *localUser) Password() *string {
	return l.password
}

func (l *localUser) Provider() string {
	return LocalProvider
}

// NewGoogleUser creates a new googleUser instance with the given fields.
//
// Since Google users authenticate externally, Password() always returns nil.
func NewGoogleUser(id, email, firstname, lastname string) User {
	return &googleUser{
		id:        id,
		email:     email,
		firstname: firstname,
		lastname:  lastname,
	}
}

func (l *googleUser) ID() string {
	return l.id
}

func (l *googleUser) Email() string {
	return l.email
}

func (l *googleUser) Firstname() string {
	return l.firstname
}

func (l *googleUser) Lastname() string {
	return l.lastname
}

func (l *googleUser) Password() *string {
	return nil
}

func (l *googleUser) Provider() string {
	return GoogleProvider
}
