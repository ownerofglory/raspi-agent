package persistence

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/ownerofglory/raspi-agent/internal/core/domain"
	"github.com/ownerofglory/raspi-agent/internal/persistence/entity"
	"gorm.io/gorm"
)

// userRepo is a GORM-based implementation of port.UserRepository.
//
// It is responsible for persisting and retrieving user entities from the database,
// and converting between domain.User and the persistence entity.
type userRepo struct {
	db *gorm.DB
}

// NewUserRepository creates a new UserRepository backed by GORM.
func NewUserRepository(db *gorm.DB) *userRepo {
	return &userRepo{db: db}
}

// Find retrieves a user by their unique identifier.
//
// It queries the database for the user entity, converts it to a domain.User,
// and returns it. If not found or conversion fails, an error is returned.
func (u *userRepo) Find(ctx context.Context, id string) (domain.User, error) {
	userEntity, err := gorm.G[entity.User](u.db).
		Where("id = ?", id).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user with id %s not found: %w", id, domain.UserNotFound)
		}
		slog.Error("unable to find user", "id", id)
		return nil, err
	}

	user, err := createDomainUser(&userEntity)
	if err != nil {
		slog.Error("unable to convert persisted user", "err", err)
		return nil, err
	}

	return user, nil
}

// FindByEmail retrieves a user by their email address.
//
// It queries the database for the user entity, converts it to a domain.User,
// and returns it. If not found or conversion fails, an error is returned.
func (u *userRepo) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	userEntity, err := gorm.G[entity.User](u.db).
		Where("email = ?", email).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user with email %s not found: %w", email, domain.UserNotFound)
		}
		slog.Error("unable to find user", "email", email)
		return nil, err
	}

	user, err := createDomainUser(&userEntity)
	if err != nil {
		slog.Error("unable to convert persisted user", "err", err)
		return nil, err
	}

	return user, nil
}

// createUserEntity converts a domain.User into a persistence entity.User.
func createUserEntity(user domain.User) *entity.User {
	return &entity.User{
		FirstName:    user.Firstname(),
		LastName:     user.Lastname(),
		Email:        user.Email(),
		Provider:     user.Provider(),
		PasswordHash: user.Password(),
	}
}

// createDomainUser converts a persistence entity.User into a domain.User.
//
// It dispatches based on the Provider field and constructs the appropriate
// domain-specific user type. If the provider is unsupported, an error is returned.
func createDomainUser(entity *entity.User) (domain.User, error) {
	var user domain.User
	switch entity.Provider {
	case domain.LocalProvider:
		user = domain.NewLocalUser(entity.ID.String(), entity.Email, *entity.PasswordHash, entity.FirstName, entity.LastName)
	case domain.GoogleProvider:
		user = domain.NewGoogleUser(entity.ID.String(), entity.Email, entity.FirstName, entity.LastName)
	default:
		return nil, fmt.Errorf("provider %s not supported", entity.Provider)
	}

	return user, nil
}
