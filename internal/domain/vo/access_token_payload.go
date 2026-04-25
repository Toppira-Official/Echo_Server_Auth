package vo

import (
	"errors"
	"time"
)

var (
	ErrAccessTokenPayloadExpiresAtMustBeAfterIssuedAt = errors.New("expires at must be after issued at")
)

type AccessTokenPayload struct {
	credentialID CredentialID
	issuedAtUTC  time.Time
	expiredAtUTC time.Time
}

func NewAccessTokenPayload(
	credentialID CredentialID,
	issuedAtUTC, expiredAtUTC time.Time,
) (AccessTokenPayload, error) {
	if !expiredAtUTC.After(issuedAtUTC) {
		return AccessTokenPayload{}, ErrAccessTokenPayloadExpiresAtMustBeAfterIssuedAt
	}

	return AccessTokenPayload{
		credentialID: credentialID,
		issuedAtUTC:  issuedAtUTC,
		expiredAtUTC: expiredAtUTC,
	}, nil
}

func (a AccessTokenPayload) CredentialID() CredentialID     { return a.credentialID }
func (a AccessTokenPayload) IssuedAtUTC() time.Time         { return a.issuedAtUTC }
func (a AccessTokenPayload) ExpiredAtUTC() time.Time        { return a.expiredAtUTC }
func (a AccessTokenPayload) IsExpired(atUTC time.Time) bool { return atUTC.After(a.expiredAtUTC) }
func (a AccessTokenPayload) IsValid(atUTC time.Time) bool   { return !a.IsExpired(atUTC) }
func (a AccessTokenPayload) Lifetime() time.Duration        { return a.expiredAtUTC.Sub(a.issuedAtUTC) }
