package middleware

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	appAuth "github.com/ownerofglory/raspi-agent/internal/auth"
	"github.com/ownerofglory/raspi-agent/pkg/auth"
)

// CertHeaderName Name of the HTTP header that contains the certificate
const CertHeaderName = "X-Forwarded-Tls-Client-Cert"

// Authenticated wraps an HTTP handler with an authentication check.
//
// It accepts an AuthenticationFunc, which inspects the incoming request
// (e.g. Authorization header, cookies, etc.), validates credentials, and
// returns an auth.Context containing a UserPrincipal.
//
// If authentication fails, the middleware writes a 401 Unauthorized
// response and does not call the next handler.
func Authenticated(af auth.AuthenticationFunc) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authCtx, err := af(w, r)
			if err != nil {
				slog.Error("Authentication error: ", "error", err)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			rAuth := r.WithContext(authCtx)

			h.ServeHTTP(w, rAuth)
		})
	}
}

// WithJWT returns an AuthenticationFunc that validates a JWT from the
// Authorization header in the form "Bearer <token>".
//
// If the token is valid, it is parsed into claims and converted into a
// UserPrincipal, which is stored in the request context.
// If the header is missing or the token is invalid, an error is returned.
func WithJWT(key string) auth.AuthenticationFunc {
	return func(rw http.ResponseWriter, r *http.Request) (auth.Context, error) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			slog.Error("Authorization header is empty")
			return nil, errors.New("Authorization header is empty")
		}

		bearerToken, found := strings.CutPrefix(authHeader, "Bearer ")
		if !found {
			slog.Error("Authorization header is invalid")
			return nil, errors.New("Authorization header is invalid")
		}

		claims, err := appAuth.ParseJWT(bearerToken, []byte(key))
		if err != nil {
			return nil, err
		}

		up := appAuth.NewUserPrincipal(claims)
		ctxWithAuth := context.WithValue(r.Context(), auth.UserPrincipalKey, up)

		authContext := auth.NewAuthContext(ctxWithAuth)

		return authContext, nil
	}
}

// WithDeviceCertHeader creates an AuthenticationFunc that extracts and parses
// an X.509 device certificate provided in a specific HTTP header.
//
// The header is expected to contain a PEM-encoded certificate (e.g. the
// `-----BEGIN CERTIFICATE-----` block). The middleware parses the certificate,
// retrieves the device identifier from the certificate's Subject Common Name (CN),
// and stores it in the request context using the `appAuth.DeviceKey`.
//
// If the header is missing, malformed, or the certificate cannot be parsed,
// the authentication will fail and an error will be returned.
//
// Example usage:
//
//	mux.Handle("/devices/{id}/data",
//		middleware.WrapFunc(
//			handleDeviceData,
//			middleware.Authenticated(middleware.WithDeviceCertHeader("X-Device-Cert")),
//			middleware.Authorized(middleware.HavingDeviceID("id")),
//		),
//	)
//
// Example header:
//
//	X-Device-Cert: -----BEGIN CERTIFICATE-----\nMIIB...==\n-----END CERTIFICATE-----
//
// The resulting context value can be retrieved later with:
//
//	deviceID, _ := r.Context().Value(appAuth.DeviceKey).(string)
func WithDeviceCertHeader(headerName string) auth.AuthenticationFunc {
	return func(rw http.ResponseWriter, r *http.Request) (auth.Context, error) {
		certHeader := r.Header.Get(headerName)
		if certHeader == "" {
			slog.Error("Certificate header is empty")
			return nil, errors.New("certificate header is empty")
		}

		block, _ := pem.Decode([]byte(certHeader))
		if block == nil {
			slog.Error("Certificate header is invalid")
			return nil, errors.New("certificate header is invalid")
		}

		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			slog.Error("Unable to parse certificate", "error", err)
			return nil, fmt.Errorf("unable to parse certificate: %w", err)
		}

		deviceID := cert.Subject.CommonName
		ctx := context.WithValue(r.Context(), appAuth.DeviceKey, deviceID)

		return auth.NewAuthContext(ctx), nil
	}
}

// Authorized wraps an HTTP handler with one or more authorization checks.
//
// Each AuthorizationFunc is invoked in order. If any of them return an error,
// the middleware writes a 403 Forbidden response and does not call the next
// handler.
//
// Typical AuthorizationFuncs include checks for user roles or matching user IDs.
func Authorized(authFuncs ...auth.AuthorizationFunc) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for _, fn := range authFuncs {
				if err := fn(w, r); err != nil {
					slog.Error("Authorization error", "error", err)
					w.WriteHeader(http.StatusForbidden)
					return
				}
			}

			h.ServeHTTP(w, r)
		})
	}
}

// HavingDeviceID creates an AuthorizationFunc that ensures the device ID
// present in the request path matches the device identifier extracted
// from a previously authenticated certificate.
//
// The `param` argument specifies the name of the path parameter that
// contains the expected device ID (for example, "id" in `/devices/{id}`).
//
// If the device ID is missing from the context, the path parameter is not
// provided, or the values do not match, the request is rejected with an error.
//
// Example:
//
//	authz := middleware.HavingDeviceID("id")
//	if err := authz(w, r); err != nil {
//		http.Error(w, err.Error(), http.StatusForbidden)
//		return
//	}
//
// Combined with WithDeviceCertHeader, this ensures that a device can only
// access its own resources, as identified by the certificate Common Name (CN).
func HavingDeviceID(param string) auth.AuthorizationFunc {
	return func(rw http.ResponseWriter, r *http.Request) error {
		aCtx := auth.NewAuthContext(r.Context())
		deviceID, ok := aCtx.Value(appAuth.DeviceKey).(string)
		if !ok {
			slog.Error("Unable to find device id in context")
			return errors.New("unable to find device id in context")
		}

		deviceIDParam := r.PathValue(param)
		if deviceIDParam == "" {
			return errors.New("missing device id in path")
		}

		if deviceIDParam != deviceID {
			return errors.New("forbidden: device id mismatch")
		}
		return nil
	}
}
