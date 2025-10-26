package domain

import "errors"

// User domain errors
var (
	UserNotFound      = errors.New("user not found")
	UserAlreadyExists = errors.New("user already exists")
)

// Device domain errors
var (
	DeviceNotFound = errors.New("device not found")
)
