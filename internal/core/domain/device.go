package domain

// DeviceRegistration represents a user-initiated request to
// register a new physical device in the system.
//
// This request typically comes from a web or mobile client
// authenticated via a user JWT. It associates the device with
// the user and generates an OTP (one-time password) that will
// later be used by the device during certificate enrollment.
//
// Fields:
//   - UserID: Identifier of the user who owns this device.
//   - Name:   Human-friendly name for the device (for UI display).
type DeviceRegistration struct {
	UserID string
	Name   string
}

// DeviceRegistrationResult is returned by the backend when a
// new device registration succeeds.
//
// It includes the newly assigned device ID, the user association,
// and the OTP that the device must use in the subsequent enrollment
// step to obtain its mTLS certificate from Step CA.
//
// Fields:
//   - DeviceID: Unique backend-generated identifier for the device.
//   - UserID:   ID of the user who owns the device.
//   - Name:     Friendly name of the device.
//   - OTP:      One-time password for device enrollment authentication.
type DeviceRegistrationResult struct {
	DeviceID string
	UserID   string
	Name     string
	OTP      string
}

// DeviceEnrollment represents a request from a device to
// enroll and obtain an X.509 certificate signed by the
// backend (via Step CA).
//
// The device generates its own private key and CSR, then
// sends the CSR along with its ID and OTP for verification.
// Upon success, the backend signs the CSR and returns a
// signed certificate chain.
//
// Fields:
//   - CSR: Base64-encoded PEM CSR (Certificate Signing Request).
//   - DeviceID:  The device’s unique identifier.
//   - OTP: One-time password for verification before signing.
//   - UserID: ID of the user the device belongs to.
type DeviceEnrollment struct {
	CSR      string
	DeviceID string
	OTP      string
	UserID   string
}

// DeviceEnrollmentResult represents the response returned
// to a successfully enrolled device.
//
// It contains the signed certificate, intermediate CA,
// and optional certificate chain details returned from
// Step CA or your enrollment provider.
//
// Fields:
//   - CertSign: The signed certificate and chain details.
type DeviceEnrollmentResult struct {
	CertSign *CertSignResult
}

// DeviceEnrollmentState represents the current status of a device
// in the enrollment lifecycle. It indicates whether a device is newly
// created, successfully enrolled with a certificate, or disabled.
type DeviceEnrollmentState string

const (
	// DeviceEnrollmentStateEnrolled means the device has successfully
	// completed enrollment and holds a valid client certificate.
	DeviceEnrollmentStateEnrolled DeviceEnrollmentState = "enrolled"

	// DeviceEnrollmentStateCreated means the device record exists
	// (e.g., registered by a user) but has not yet completed the
	// certificate enrollment process.
	DeviceEnrollmentStateCreated DeviceEnrollmentState = "created"

	// DeviceEnrollmentStateDisabled means the device has been explicitly
	// disabled or revoked and can no longer authenticate with the system.
	DeviceEnrollmentStateDisabled DeviceEnrollmentState = "disabled"
)

// Device represents a registered hardware or software client (e.g., a Raspberry Pi).
// It is uniquely identified and associated with a user account. Devices are enrolled
// by generating a CSR and receiving a signed certificate from the backend.
type Device struct {
	// ID is the unique identifier of the device (e.g., UUID or serial number).
	ID *string

	// UserID associates the device with a user or account that owns it.
	UserID *string

	// OTP is an optional one-time passcode issued during device registration.
	// It is used to authenticate the device during its first enrollment.
	OTP *string

	// EnrollmentStatus represents the device’s current lifecycle state.
	// It indicates whether the device is registered, enrolled, or disabled.
	EnrollmentStatus DeviceEnrollmentState
}
