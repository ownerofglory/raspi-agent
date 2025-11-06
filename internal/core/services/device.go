package services

import (
	"context"
	"crypto/rand"
	"fmt"
	"log/slog"
	"math/big"

	"github.com/ownerofglory/raspi-agent/internal/core/domain"
	"github.com/ownerofglory/raspi-agent/internal/core/ports"
)

const (
	passwordCharset = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"0123456789" +
		"!@#$%^&*()-_=+[]{}<>?/|"
)

type deviceService struct {
	userRepo    ports.UserRepo
	deviceRepo  ports.DeviceRepo
	certHandler ports.EnrollmentHandler
}

func NewDeviceService(userRepo ports.UserRepo, deviceRepo ports.DeviceRepo, certHandler ports.EnrollmentHandler) *deviceService {
	return &deviceService{
		userRepo:    userRepo,
		deviceRepo:  deviceRepo,
		certHandler: certHandler,
	}
}

func (s *deviceService) RegisterDevice(ctx context.Context, reg domain.DeviceRegistration) (*domain.DeviceRegistrationResult, error) {
	user, err := s.userRepo.Find(ctx, reg.UserID)
	if err != nil {
		slog.Error("failed to find user", "userId", reg.UserID)
		return nil, err
	}

	userID := reg.UserID
	otp, err := generatePassword(16)
	if err != nil {
		slog.Error("failed to generate password for user", "userId", reg.UserID)
		return nil, fmt.Errorf("failed to generate otp: %w", err)
	}
	device := domain.Device{
		UserID:           &userID,
		OTP:              &otp,
		EnrollmentStatus: domain.DeviceEnrollmentStateCreated,
	}

	saved, err := s.deviceRepo.Save(ctx, device)
	if err != nil {
		slog.Error("failed to save device", "userId", reg.UserID)
		return nil, fmt.Errorf("failed to save device: %w", err)
	}

	res := domain.DeviceRegistrationResult{
		DeviceID: *saved.ID,
		UserID:   user.ID(),
		Name:     saved.Name,
	}

	return &res, nil
}

func (s *deviceService) EnrollDevice(ctx context.Context, enr domain.DeviceEnrollment) (*domain.DeviceEnrollmentResult, error) {
	device, err := s.deviceRepo.Find(ctx, enr.DeviceID)
	if err != nil {
		slog.Error("failed to find device", "deviceID", enr.DeviceID)
		return nil, fmt.Errorf("failed to find device: %w", err)
	}

	if device.UserID == nil {
		slog.Error("failed to find device user", "deviceID", enr.DeviceID)
		return nil, fmt.Errorf("failed to find device user: %w", err)
	}

	if *device.UserID != enr.UserID {
		slog.Error("Device is not register to the user", "deviceID", enr.DeviceID, "userID", enr.UserID)
		return nil, fmt.Errorf("device is not register to the user device id :%s, user id: %s", enr.DeviceID, enr.UserID)
	}

	if device.OTP == nil || enr.OTP == *device.OTP {
		slog.Error("Device OTP is not enrolled", "deviceID", enr.DeviceID)
		return nil, fmt.Errorf("device OTP is not enrolled")
	}

	certSignResult, err := s.certHandler.Sign(ctx, &domain.CertSignRequest{
		CSR:      enr.CSR,
		DeviceID: enr.DeviceID,
	})
	if err != nil {
		slog.Error("failed to sign certificate", "deviceID", enr.DeviceID, "error", err)
		return nil, fmt.Errorf("failed to sign certificate: %w", err)
	}

	res := domain.DeviceEnrollmentResult{
		CertSign: certSignResult,
	}

	return &res, nil
}

// generatePassword creates a cryptographically secure random password
// of the given length. It uses only Go's standard library (crypto/rand),
// so it’s safe for device OTPs, API keys, or temporary credentials.
//
// Example:
//
//	pass, _ := GeneratePassword(16)
//	fmt.Println(pass)
//
// Returns an error only if the system random source fails.
func generatePassword(length int) (string, error) {
	if length <= 0 {
		length = 16
	}

	var password []byte
	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(passwordCharset))))
		if err != nil {
			return "", err
		}
		password = append(password, passwordCharset[n.Int64()])
	}
	return string(password), nil
}
