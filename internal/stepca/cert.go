package stepca

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/ownerofglory/raspi-agent/internal/core/domain"
	"go.step.sm/crypto/jose"
	"go.step.sm/crypto/randutil"
)

const idLength = 64

type enrollProvider struct {
	stepCAURL        string
	provisionerToken string
	provisionerName  string
	pem              []byte
	jwk              []byte
	client           *http.Client
}

func NewProvider(stepCAURL, provisionerName, provisionerToken string, pem, jwk []byte) *enrollProvider {
	rootCAs := x509.NewCertPool()
	rootCAs.AppendCertsFromPEM(pem)

	return &enrollProvider{
		stepCAURL:        stepCAURL,
		provisionerToken: provisionerToken,
		provisionerName:  provisionerName,
		jwk:              jwk,
		pem:              pem,
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					RootCAs: rootCAs,
				},
			},
		},
	}
}

type stepCASignResponse struct {
	Crt        string   `json:"crt"`
	CA         string   `json:"ca"`
	CertChain  []string `json:"certChain"`
	TLSOptions *struct {
		CipherSuites  []string `json:"cipherSuite,omitempty"`
		MinVersion    float64  `json:"minVersion"`
		MaxVersion    float64  `json:"maxVersion"`
		Renegotiation bool     `json:"renegotiation"`
	} `json:"tlsOptions,omitempty"`
}

func (e *enrollProvider) Sign(ctx context.Context, req *domain.CertEnrollRequest) (*domain.CertSignResult, error) {
	signReq := map[string]string{
		"csr": req.CSR,
		"ott": e.provisionerToken,
	}

	body, err := json.Marshal(signReq)
	if err != nil {
		slog.Error("enroll enroll: failed to marshal enrollment request", "err", err)
		return nil, fmt.Errorf("enroll enroll: failed to marshal enrollment request")
	}

	url := fmt.Sprintf("%s/1.0/sign", e.stepCAURL)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		slog.Error("enroll enroll: failed to create http request", "err", err)
		return nil, fmt.Errorf("enroll enroll: failed to create http request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+e.provisionerToken)

	resp, err := e.client.Do(httpReq)
	if err != nil {
		slog.Error("enroll enroll: failed to send enrollment request", "err", err)
		return nil, fmt.Errorf("enroll enroll: failed to send enrollment request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			slog.Error("enroll enroll: failed to read enrollment response body", "err", err)
			return nil, fmt.Errorf("enroll enroll: failed to read enrollment response body: %w", err)
		}

		slog.Error("Step CA error", "status", resp.Status, "body", string(b))
		return nil, fmt.Errorf("step CA error: %s", string(b))
	}

	var signResp stepCASignResponse
	if err := json.NewDecoder(resp.Body).Decode(&signResp); err != nil {
		slog.Error("Unable to decode step CA sign response", "err", err)
		return nil, fmt.Errorf("enroll enroll: failed to decode step CA sign response: %w", err)
	}

	return &domain.CertSignResult{
		Crt:       signResp.Crt,
		Ca:        signResp.CA,
		CertChain: signResp.CertChain,
	}, nil

}

// generateOTT generates a one-time token for a CSR request.
func (e *enrollProvider) generateOTT(deviceCN string) (string, error) {
	now := time.Now()

	// Load the JWK key
	jwk, err := jose.ParseKey(e.jwk, jose.WithPassword([]byte(e.provisionerToken)))
	if err != nil {
		return "", err
	}

	// Create JWT signer
	opts := new(jose.SignerOptions).WithType("JWT").WithHeader("kid", jwk.KeyID)
	sig, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.ES256, Key: jwk.Key}, opts)
	if err != nil {
		return "", err
	}

	id, err := randutil.ASCII(idLength)
	if err != nil {
		return "", err
	}

	// Claims expected by Step CA
	claims := struct {
		jose.Claims
		SANS []string `json:"sans"`
	}{
		Claims: jose.Claims{
			ID:        id,
			Subject:   deviceCN,
			Issuer:    e.provisionerName,
			NotBefore: jose.NewNumericDate(now),
			Expiry:    jose.NewNumericDate(now.Add(time.Minute)),
			Audience:  []string{e.stepCAURL},
		},
		SANS: []string{deviceCN},
	}

	return jose.Signed(sig).Claims(claims).CompactSerialize()
}
