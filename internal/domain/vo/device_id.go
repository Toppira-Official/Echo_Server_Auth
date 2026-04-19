package vo

import "errors"

var (
	ErrDeviceIDRequired = errors.New("device id required")
)

type DeviceID struct {
	value string
}

func NewDeviceID(id string) (DeviceID, error) {
	if id == "" {
		return DeviceID{}, ErrDeviceIDRequired
	}

	return DeviceID{value: id}, nil
}

func (u DeviceID) Value() string { return u.value }
