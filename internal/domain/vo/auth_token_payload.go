package vo

import (
	"errors"
	"time"
)

var (
	ErrAuthTokenPayloadExpiresAtMustBeAfterIssuerAt = errors.New("expires at must be after issuer at")
)

type AuthTokenPayload struct {
	userID       UserID
	issuedAtUTC  time.Time
	expiredAtUTC time.Time
}

func NewAuthTokenPayload(
	userID UserID,
	issuedAt, expiredAt time.Time,
) (AuthTokenPayload, error) {
	if expiredAt.Before(issuedAt) {
		return AuthTokenPayload{}, ErrAuthTokenPayloadExpiresAtMustBeAfterIssuerAt
	}

	return AuthTokenPayload{
		userID:       userID,
		issuedAtUTC:  issuedAt,
		expiredAtUTC: expiredAt,
	}, nil
}

func (a AuthTokenPayload) UserID() UserID                 { return a.userID }
func (a AuthTokenPayload) IssuedAtUTC() time.Time         { return a.issuedAtUTC }
func (a AuthTokenPayload) ExpiredAtUTC() time.Time        { return a.expiredAtUTC }
func (a AuthTokenPayload) IsExpired(atUTC time.Time) bool { return atUTC.After(a.expiredAtUTC) }
func (a AuthTokenPayload) IsValid(atUTC time.Time) bool   { return !a.IsExpired(atUTC) }
func (a AuthTokenPayload) Lifetime() time.Duration        { return a.expiredAtUTC.Sub(a.issuedAtUTC) }
