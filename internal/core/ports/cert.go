package ports

import (
	"context"

	"github.com/ownerofglory/raspi-agent/internal/core/domain"
)

//go:generate go tool go.uber.org/mock/mockgen -source=cert.go -package=ports -destination=cert_mock.go EnrollmentHandler

// EnrollmentHandler defines the contract for any component capable
// of handling certificate enrollment for devices.
//
// The Enroll method receives a certificate signing request (CSR) from a device,
// sends it to a CA (e.g., Step CA) for signing, and returns the resulting
// signed certificate bundle.
//
// Implementations should validate and authenticate enrollment requests
// to ensure that only authorized devices are able to obtain certificates.
type EnrollmentHandler interface {
	// Sign handles the certificate enrollment process for a device.
	//
	// req contains the CSR to be signed.
	// The method should return a CertSignResult containing the signed
	// certificate, CA certificate, and any intermediate certificates.
	//
	// It should return an error if the signing process fails or if the
	// enrollment request is invalid or unauthorized.
	Sign(ctx context.Context, req *domain.CertSignRequest) (*domain.CertSignResult, error)
}
