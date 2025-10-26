package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/ownerofglory/raspi-agent/internal/core/domain"
	"github.com/ownerofglory/raspi-agent/internal/core/ports"
)

const (
	// PostRegisterDeviceURL is the backend API path for device registration.
	// A new device record is created for the specified user.
	PostRegisterDeviceURL = baseManagementPath + "/v1/users/{userId}/devices"

	// PostEnrollDeviceURL is the backend API path for device enrollment.
	// The device submits its CSR and OTP to obtain a signed certificate.
	PostEnrollDeviceURL = baseManagementPath + "/v1/users/{userId}/devices/{deviceId}/enroll"
)

// deviceRegistrationReq defines the JSON payload for registering a new device.
type deviceRegistrationReq struct {
	// Name is the human-readable name for the device (e.g., “Raspberry Pi Living Room”).
	Name string `json:"name"`
}

// deviceRegisterResp defines the JSON response returned after a device is registered.
type deviceRegisterResp struct {
	DeviceID string `json:"deviceId"`
	UserID   string `json:"userId"`
	Name     string `json:"name"`
	OTP      string `json:"otp"`
}

// deviceEnrollmentReq defines the JSON payload for device enrollment.
// The CSR should be the PEM-encoded CSR, and OTP should match the one issued during registration.
type deviceEnrollmentReq struct {
	CSR string `json:"csr"`
	OTP string `json:"otp"`
}

// deviceEnrollmentResp defines the JSON response returned after a successful enrollment.
type deviceEnrollmentResp struct {
	CertSign struct {
		Crt       string   `json:"crt"`
		Ca        string   `json:"ca"`
		CertChain []string `json:"certChain"`
	} `json:"certSign"`
}

// deviceHandler handles device registration and enrollment HTTP requests.
// It delegates business logic to the injected DeviceService.
type deviceHandler struct {
	service ports.DeviceService
}

// NewDeviceHandler returns a new instance of deviceHandler.
func NewDeviceHandler(service ports.DeviceService) *deviceHandler {
	return &deviceHandler{service: service}
}

// HandlePostRegisterDevice registers a new device for a given user.
//
// Endpoint: POST /v1/users/{userId}/devices
//
// Expected JSON body:
//
//	{
//	  "name": "Raspberry Pi 5"
//	}
//
// Response 200 OK:
//
//	{
//	  "deviceId": "1234-abcd",
//	  "userId": "user-5678",
//	  "name": "Raspberry Pi 5",
//	  "otp": "XYZA12"
//	}
func (d *deviceHandler) HandlePostRegisterDevice(rw http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userId")

	defer r.Body.Close()
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	var req deviceRegistrationReq
	err = json.Unmarshal(reqBody, &req)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	device, err := d.service.RegisterDevice(r.Context(), domain.DeviceRegistration{
		UserID: userID,
		Name:   req.Name,
	})
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	reg := deviceRegisterResp{
		OTP:      device.OTP,
		Name:     device.Name,
		DeviceID: device.DeviceID,
		UserID:   userID,
	}
	respBody, err := json.Marshal(reg)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	_, _ = rw.Write(respBody)
}

// HandlePostEnrollDevice enrolls a device by verifying its OTP and signing its CSR.
//
// Endpoint: POST /v1/users/{userId}/devices/{deviceId}/enroll
//
// Expected JSON body:
//
//	{
//	  "csr": "-----BEGIN CERTIFICATE REQUEST-----...",
//	  "otp": "XYZA12"
//	}
//
// Response 200 OK:
//
//	{
//	  "certSign": {
//	    "crt": "-----BEGIN CERTIFICATE-----...",
//	    "ca": "-----BEGIN CERTIFICATE-----...",
//	    "certChain": ["..."]
//	  }
//	}
func (d *deviceHandler) HandlePostEnrollDevice(rw http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userId")
	deviceID := r.PathValue("deviceId")

	defer r.Body.Close()
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	var req deviceEnrollmentReq
	err = json.Unmarshal(reqBody, &req)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	enrollment, err := d.service.EnrollDevice(r.Context(), domain.DeviceEnrollment{
		UserID:   userID,
		DeviceID: deviceID,
		CSR:      req.CSR,
		OTP:      req.OTP,
	})
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	enr := deviceEnrollmentResp{
		CertSign: struct {
			Crt       string   `json:"crt"`
			Ca        string   `json:"ca"`
			CertChain []string `json:"certChain"`
		}{
			Ca:        enrollment.CertSign.Ca,
			Crt:       enrollment.CertSign.Crt,
			CertChain: enrollment.CertSign.CertChain,
		},
	}

	respBody, err := json.Marshal(enr)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	_, _ = rw.Write(respBody)
}
