package domain

// CertEnrollRequest represents a certificate enrollment request
// coming from a device or client. The CSR field should contain
// a base64-encoded PEM CSR (Certificate Signing Request).
type CertEnrollRequest struct {
	// CSR is the base64-encoded PEM-encoded CSR to be signed.
	CSR string
}

// CertSignResult represents the result of a successful certificate
// signing operation from the CA. It includes the issued certificate,
// the issuing CA certificate, and the full certificate chain if applicable.
type CertSignResult struct {
	// Crt is the PEM-encoded X.509 certificate issued to the device.
	Crt string

	// Ca is the PEM-encoded certificate of the issuing Certificate Authority.
	Ca string

	// CertChain optionally includes the entire PEM-encoded certificate chain,
	// starting from the issued certificate up to the root CA.
	CertChain []string
}
