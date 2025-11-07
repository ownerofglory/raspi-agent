package services

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ownerofglory/raspi-agent/internal/core/domain"
	"github.com/ownerofglory/raspi-agent/internal/core/ports"
)

// userService provides the business logic for managing users.
//
// It depends on a UserRepository for persistence and encapsulates
// higher-level concerns such as logging and error wrapping.
type userService struct {
	userRepo ports.UserRepo
}

// NewUserService creates a new instance of userService.
//
// The returned service uses the given UserRepository for persistence.
func NewUserService(userRepo ports.UserRepo) *userService {
	return &userService{userRepo: userRepo}
}

// CreateUser creates a new user in the system.
//
// It logs the operation, persists the user using the repository,
// and returns the saved entity. If persistence fails, the error is
// logged and wrapped before being returned.
func (us *userService) CreateUser(ctx context.Context, user domain.User) (domain.User, error) {
	slog.Info("create user")
	slog.Debug("User", "val", user)

	savedUser, err := us.userRepo.Save(ctx, user)
	if err != nil {
		slog.Error("error when saving user", "err", err)
		return nil, fmt.Errorf("error when saving user: %w", err)
	}

	return *savedUser, nil
}

// GetUser retrieves a user by their unique identifier.
//
// It delegates to the repository and wraps any returned error with
// additional context.
func (us *userService) GetUser(ctx context.Context, id string) (domain.User, error) {
	slog.Debug("get user by id", "id", id)

	user, err := us.userRepo.Find(ctx, id)
	if err != nil {
		slog.Error("error when getting user", "err", err)
		return nil, fmt.Errorf("error when getting user: %w", err)
	}
	return user, nil
}

func (us *userService) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	slog.Debug("get user by email", "email", email)

	user, err := us.userRepo.FindByEmail(ctx, email)
	if err != nil {
		slog.Error("error when getting user", "err", err)
		return nil, fmt.Errorf("error when getting user: %w", err)
	}
	return user, nil
}
