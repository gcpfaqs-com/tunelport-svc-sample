package types

import "errors"

var (
	ErrDeviceNotFound      = errors.New("device not found")
	ErrDeviceAlreadyExists = errors.New("device already exists")
	ErrInvalidStatus       = errors.New("invalid device status")
)
