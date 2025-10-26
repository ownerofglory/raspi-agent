package device

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/ownerofglory/raspi-agent/internal/core/domain"
)

const backendBasePath = "/raspi-agent/management/api"

// PostEnrollDeviceURL defines the backend API path for device enrollment.
// The placeholder {deviceId} is replaced with the actual device identifier.
const PostEnrollDeviceURL string = backendBasePath + "/v1/device/{deviceId}/enroll"

// Default subject information for generated CSRs.
// These values are embedded into the certificate subject (DN).
const (
	defaultCountry  = "DE"
	defaultState    = "BW"
	defaultLocation = "Stuttgart"
	defaultOrg      = "ownerofglory"
	defaultOrgUnit  = "Devices"
)

// Filenames for generated private key and CSR artifacts.
const (
	deviceKeyFile = "device.key"
	deviceCSRFile = "device.csr"
)

// deviceEnrollClient is a client for enrolling devices with the backend.
// It handles generating a private key and CSR, sending it to the backend
// enrollment endpoint, and receiving a signed certificate.
type deviceEnrollClient struct {
	client  *http.Client
	baseURL string
}

// deviceEnrollRequest represents the JSON payload sent to the backend during enrollment.
type deviceEnrollRequest struct {
	CSR      string `json:"csr"`
	DeviceID string `json:"deviceId"`
	OTP      string `json:"otp"`
	UserID   string `json:"userId"`
}

// deviceEnrollResponse represents the expected JSON response
// from the backend, containing the signed certificate chain.
type deviceEnrollResponse struct {
	CertSign struct {
		Crt       string   `json:"crt"`
		Ca        string   `json:"ca"`
		CertChain []string `json:"certChain"`
	} `json:"certSign"`
}

// NewDeviceEnrollClient creates a new enrollment client that communicates
// with the backend device management API.
func NewDeviceEnrollClient(baseURL string) *deviceEnrollClient {
	return &deviceEnrollClient{
		client:  &http.Client{},
		baseURL: baseURL,
	}
}

// Enroll generates a CSR for the given device, sends it to the backend
// along with user and OTP information, and returns the signed certificate chain.
//
// The resulting certificate can then be stored locally and used for mTLS authentication
func (d *deviceEnrollClient) Enroll(deviceID, userID, otp string) (*domain.DeviceEnrollmentResult, error) {
	err := d.generateCSR(deviceID)
	if err != nil {
		return nil, fmt.Errorf("enroll device csr: %w", err)
	}

	deviceCSR, err := os.ReadFile(deviceCSRFile)
	if err != nil {
		return nil, fmt.Errorf("enroll device csr: %w", err)
	}

	reqBody := deviceEnrollRequest{
		DeviceID: deviceID,
		OTP:      otp,
		UserID:   userID,
		CSR:      string(deviceCSR),
	}

	reqPayload, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("enroll device csr: %w", err)
	}

	path := strings.Replace(PostEnrollDeviceURL, "{deviceId}", deviceID, 1)
	url := fmt.Sprintf("%s%s", d.baseURL, path)

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqPayload))
	if err != nil {
		return nil, fmt.Errorf("enroll device csr: %w", err)
	}

	resp, err := d.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("enroll device csr: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("enroll device csr: bad status code: %d", resp.StatusCode)
	}

	respPayload, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("enroll device csr: %w", err)
	}

	dRes := deviceEnrollResponse{}
	err = json.Unmarshal(respPayload, &dRes)
	if err != nil {
		return nil, fmt.Errorf("enroll device csr: %w", err)
	}

	result := domain.DeviceEnrollmentResult{
		CertSign: &domain.CertSignResult{
			Crt:       dRes.CertSign.Crt,
			Ca:        dRes.CertSign.Ca,
			CertChain: dRes.CertSign.CertChain,
		},
	}

	return &result, nil
}

// generateCSR creates a private key and CSR for a given device ID,
// following the same fields as the OpenSSL-based script.
//
// The CSR includes the device ID as the Common Name (CN),
// and SAN (Subject Alternative Name) with an Extended Key Usage for clientAuth.
func (d *deviceEnrollClient) generateCSR(deviceID string) error {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("failed to generate private key: %w", err)
	}

	keyOut, err := os.Create(deviceKeyFile)
	if err != nil {
		return fmt.Errorf("failed to create device key: %w", err)
	}
	defer keyOut.Close()
	err = pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	if err != nil {
		return fmt.Errorf("failed to encode rsa private key: %w", err)
	}

	csrTemplate := x509.CertificateRequest{
		Subject: pkix.Name{
			Country:            []string{defaultCountry},
			Province:           []string{defaultState},
			Locality:           []string{defaultLocation},
			Organization:       []string{defaultOrg},
			OrganizationalUnit: []string{defaultOrgUnit},
			CommonName:         deviceID,
		},
		DNSNames: []string{deviceID},
		ExtraExtensions: []pkix.Extension{
			{
				Id:       []int{2, 5, 29, 37}, // OID for extendedKeyUsage
				Critical: false,
				Value: mustMarshalExtKeyUsage([]x509.ExtKeyUsage{
					x509.ExtKeyUsageClientAuth,
				}),
			},
		},
	}

	csrDER, err := x509.CreateCertificateRequest(rand.Reader, &csrTemplate, key)
	if err != nil {
		return fmt.Errorf("failed to create certificate request: %w", err)
	}

	csrOut, err := os.Create(deviceCSRFile)
	if err != nil {
		return fmt.Errorf("failed to create device csr: %w", err)
	}
	defer csrOut.Close()
	err = pem.Encode(csrOut, &pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrDER})
	if err != nil {
		return fmt.Errorf("failed to encode device csr: %w", err)
	}

	return nil
}

// mustMarshalExtKeyUsage encodes an ExtKeyUsage extension as ASN.1 DER.
func mustMarshalExtKeyUsage(usages []x509.ExtKeyUsage) []byte {
	ekuOIDs := []asn1.ObjectIdentifier{}
	for _, usage := range usages {
		switch usage {
		case x509.ExtKeyUsageClientAuth:
			ekuOIDs = append(ekuOIDs, []int{1, 3, 6, 1, 5, 5, 7, 3, 2})
		case x509.ExtKeyUsageServerAuth:
			ekuOIDs = append(ekuOIDs, []int{1, 3, 6, 1, 5, 5, 7, 3, 1})
		default:
			continue
		}
	}
	b, err := asn1.Marshal(ekuOIDs)
	if err != nil {
		panic(err)
	}
	return b
}
