package ports

import (
	"context"

	"github.com/ownerofglory/raspi-agent/internal/core/domain"
)

// DeviceService defines the high-level contract for managing
// the lifecycle of physical IoT or edge devices — from initial
// registration (by a user) to secure enrollment (by the device).
//
// It’s typically implemented by a backend service that integrates
// with:
//   - A user management or identity provider (for user linkage).
//   - A certificate authority (e.g., Step CA) for issuing mTLS certs.
type DeviceService interface {
	// RegisterDevice registers a new device under a specific user.
	// The backend generates:
	//   - A unique DeviceID (e.g. "device-raspi-001")
	//   - A one-time password (OTP)
	//
	// These values are returned in DeviceRegistrationResult and
	// must later be provided by the device during enrollment.
	//
	// Example:
	//   req := domain.DeviceRegistration{
	//     UserID: "usr_12345",
	//     Name:   "Living Room Pi",
	//   }
	//   res, err := svc.RegisterDevice(ctx, req)
	//
	// Returns:
	//   - DeviceRegistrationResult with DeviceID, UserID, Name, OTP
	//   - Error if the operation fails
	RegisterDevice(ctx context.Context, reg domain.DeviceRegistration) (*domain.DeviceRegistrationResult, error)

	// EnrollDevice handles the secure enrollment of a registered device.
	//
	// The device provides:
	//   - Its DeviceID (issued during registration)
	//   - A valid OTP
	//   - A PEM-encoded CSR (Certificate Signing Request)
	//
	// The backend validates the OTP, then forwards the CSR to the
	// configured Certificate Authority (e.g., Step CA) for signing.
	//
	// Returns:
	//   - DeviceEnrollmentResult containing the signed certificate
	//   - Error if OTP validation or signing fails
	EnrollDevice(ctx context.Context, enr domain.DeviceEnrollment) (*domain.DeviceEnrollmentResult, error)
}

// DeviceRepo defines the data access layer for device records.
// It abstracts persistence logic (e.g., database, in-memory, or remote store)
// from higher-level business logic in the service layer.
//
// Implementations of this interface should be responsible for creating,
// updating, retrieving, and deleting device entities in the storage backend.
type DeviceRepo interface {
	// Save persists a new device record in the repository.
	// It should return the created device (including its generated ID, if any).
	Save(ctx context.Context, device domain.Device) (*domain.Device, error)

	// Update modifies an existing device record.
	// Returns the updated device instance or an error if not found.
	Update(ctx context.Context, device domain.Device) (*domain.Device, error)

	// Find retrieves a single device by its unique identifier.
	// Returns nil if no matching device is found.
	Find(ctx context.Context, id string) (*domain.Device, error)

	// Remove deletes a device from the repository by its ID.
	// It should be idempotent — calling Remove on a non-existent ID
	// should not cause an error.
	Remove(ctx context.Context, id string) error

	// FindByUserID returns all devices associated with a given user.
	// The slice may be empty if the user has no registered devices.
	FindByUserID(ctx context.Context, userId string) ([]domain.Device, error)
}
