package entity

import (
	"auth/internal/domain/vo"
	"errors"
	"time"
)

var (
	ErrCredentialUsernameRequired = errors.New("username is required")
)

type Credential struct {
	id             vo.UserID
	username       string
	hashedPassword vo.HashedPassword
	createdAt      time.Time
	updatedAt      time.Time
}

func NewCredential(
	id vo.UserID,
	username string,
	nowUTC time.Time,
	hashedPassword vo.HashedPassword,
) (*Credential, error) {
	if username == "" {
		return nil, ErrCredentialUsernameRequired
	}

	return &Credential{
		id:             id,
		username:       username,
		hashedPassword: hashedPassword,
		createdAt:      nowUTC,
		updatedAt:      nowUTC,
	}, nil
}

func (c *Credential) ID() vo.UserID                     { return c.id }
func (c *Credential) Username() string                  { return c.username }
func (c *Credential) HashedPassword() vo.HashedPassword { return c.hashedPassword }
func (c *Credential) CreatedAt() time.Time              { return c.createdAt }
func (c *Credential) UpdatedAt() time.Time              { return c.updatedAt }

func (c *Credential) ChangePassword(newHashedPassword vo.HashedPassword, nowUTC time.Time) {
	c.hashedPassword = newHashedPassword
	c.updatedAt = nowUTC
}
