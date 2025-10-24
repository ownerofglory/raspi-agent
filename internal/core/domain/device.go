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
