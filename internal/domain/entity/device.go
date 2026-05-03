package entity

import (
	"auth/internal/domain/vo"
	"errors"
	"time"
)

var (
	ErrDeviceRefreshTokenRequired = errors.New("refresh token is required")
	ErrDeviceUserAgentRequired    = errors.New("user agent is required")
	ErrDeviceIPAddressRequired    = errors.New("ip address is required")
	ErrDeviceExpiresAtInvalid     = errors.New("expires at must be in the future")
)

type Device struct {
	id           vo.DeviceID
	credentialID vo.CredentialID
	refreshToken string
	userAgent    string
	ipAddress    string
	expiresAt    time.Time
	lastUsedAt   time.Time
}

func NewDevice(
	id vo.DeviceID,
	refreshToken string,
	expiresAt, nowUTC time.Time,
	userAgent string, ip string,
) (*Device, error) {
	if refreshToken == "" {
		return nil, ErrDeviceRefreshTokenRequired
	}

	if userAgent == "" {
		return nil, ErrDeviceUserAgentRequired
	}

	if ip == "" {
		return nil, ErrDeviceIPAddressRequired
	}

	if !expiresAt.After(nowUTC) {
		return nil, ErrDeviceExpiresAtInvalid
	}

	return &Device{
		id:           id,
		refreshToken: refreshToken,
		expiresAt:    expiresAt,
		lastUsedAt:   nowUTC,
		userAgent:    userAgent,
		ipAddress:    ip,
	}, nil
}

func (d *Device) ID() vo.DeviceID               { return d.id }
func (d *Device) CredentialID() vo.CredentialID { return d.credentialID }
func (d *Device) RefreshToken() string          { return d.refreshToken }
func (d *Device) UserAgent() string             { return d.userAgent }
func (d *Device) IPAddress() string             { return d.ipAddress }
func (d *Device) ExpiresAt() time.Time          { return d.expiresAt }
func (d *Device) LastUsedAt() time.Time         { return d.lastUsedAt }

func (d *Device) UpdateRefreshToken(token string, expiresAt time.Time, nowUTC time.Time) {
	d.refreshToken = token
	d.expiresAt = expiresAt
	d.lastUsedAt = nowUTC
}

func (d *Device) Revoke(nowUTC time.Time) {
	d.refreshToken = ""
}
