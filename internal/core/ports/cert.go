package ports

import (
	"context"

	"github.com/ownerofglory/raspi-agent/internal/core/domain"
)

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
	// Enroll handles the certificate enrollment process for a device.
	//
	// req contains the CSR to be signed.
	// The method should return a CertSignResult containing the signed
	// certificate, CA certificate, and any intermediate certificates.
	//
	// It should return an error if the signing process fails or if the
	// enrollment request is invalid or unauthorized.
	Enroll(ctx context.Context, req *domain.CertEnrollRequest) (*domain.CertSignResult, error)
}
