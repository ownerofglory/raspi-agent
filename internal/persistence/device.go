package persistence

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/ownerofglory/raspi-agent/internal/core/domain"
	"github.com/ownerofglory/raspi-agent/internal/persistence/entity"
	"gorm.io/gorm"
)

// deviceRepo is a GORM-based implementation of ports.DeviceRepo.
//
// It belongs to the **infrastructure / adapter** layer in the hexagonal architecture
// and persists domain.Device entities using PostgreSQL (or any SQL driver supported by GORM).
type deviceRepo struct {
	db *gorm.DB
}

// NewDeviceRepo creates a new GORM-backed Device repository.

func NewDeviceRepo(db *gorm.DB) *deviceRepo {
	return &deviceRepo{db: db}
}

// Save inserts a new device into the database.
//
// If the device already exists, consider using Update instead.
// Returns the saved domain.Device.
func (r *deviceRepo) Save(ctx context.Context, device domain.Device) (*domain.Device, error) {
	entityDevice, err := toDeviceEntity(device)
	if err != nil {
		return nil, fmt.Errorf("save device: %w", err)
	}

	if err := r.db.WithContext(ctx).Create(&entityDevice).Error; err != nil {
		slog.Error("failed to save device", "err", err, "device", device)
		return nil, fmt.Errorf("save device: %w", err)
	}

	return toDomainDevice(&entityDevice), nil
}

// Update modifies an existing device entry.
//
// Returns the updated domain.Device.
func (r *deviceRepo) Update(ctx context.Context, device domain.Device) (*domain.Device, error) {
	entityDevice, err := toDeviceEntity(device)
	if err != nil {
		return nil, fmt.Errorf("update device: %w", err)
	}

	if err := r.db.WithContext(ctx).Save(&entityDevice).Error; err != nil {
		slog.Error("failed to update device", "err", err, "device_id", device.ID)
		return nil, fmt.Errorf("update device: %w", err)
	}

	return toDomainDevice(&entityDevice), nil
}

// Find retrieves a single device by its ID.
func (r *deviceRepo) Find(ctx context.Context, id string) (*domain.Device, error) {
	var e entity.Device
	if err := r.db.WithContext(ctx).
		Preload("User").
		First(&e, "id = ?", id).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("device %s not found: %w", id, domain.ErrDeviceNotFound)
		}

		slog.Error("failed to find device", "err", err, "id", id)
		return nil, err
	}

	return toDomainDevice(&e), nil
}

// Remove deletes a device by ID.
func (r *deviceRepo) Remove(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Delete(&entity.Device{}, "id = ?", id).Error; err != nil {
		slog.Error("failed to delete device", "err", err, "id", id)
		return fmt.Errorf("remove device: %w", err)
	}
	return nil
}

// FindByUserID retrieves all devices belonging to a specific user.
func (r *deviceRepo) FindByUserID(ctx context.Context, userID string) ([]domain.Device, error) {
	var entities []entity.Device

	if err := r.db.WithContext(ctx).
		Joins("User").
		Where("users.id = ?", userID).
		Find(&entities).Error; err != nil {

		slog.Error("failed to find devices by user id", "err", err, "user_id", userID)
		return nil, fmt.Errorf("find devices by user id: %w", err)
	}

	devices := make([]domain.Device, 0, len(entities))
	for _, e := range entities {
		devices = append(devices, *toDomainDevice(&e))
	}

	return devices, nil
}

// toDeviceEntity converts a domain.Device to a persistence entity.Device.
func toDeviceEntity(d domain.Device) (entity.Device, error) {
	var e entity.Device

	// Parse device ID
	if d.ID != nil {
		id, err := uuid.Parse(*d.ID)
		if err != nil {
			return e, fmt.Errorf("invalid device ID: %w", err)
		}
		e.ID = id
	}

	e.Name = d.Name
	if d.OTP != nil {
		e.OTP = *d.OTP
	}
	e.EnrollmentStatus = string(d.EnrollmentStatus)

	// Attach user if present
	if d.UserID != nil {
		userUUID, err := uuid.Parse(*d.UserID)
		if err != nil {
			return e, fmt.Errorf("invalid user ID: %w", err)
		}
		e.User = &entity.User{ID: userUUID}
	}

	return e, nil
}

// toDomainDevice converts a persistence entity.Device to a domain.Device.
func toDomainDevice(e *entity.Device) *domain.Device {
	idStr := e.ID.String()
	var userID *string

	if e.User != nil {
		uid := e.User.ID.String()
		userID = &uid
	}

	otp := e.OTP

	return &domain.Device{
		ID:               &idStr,
		UserID:           userID,
		OTP:              &otp,
		Name:             e.Name,
		EnrollmentStatus: domain.DeviceEnrollmentState(e.EnrollmentStatus),
	}
}
