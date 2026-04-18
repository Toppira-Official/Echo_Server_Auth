package vo

import (
	"errors"
	"strings"
)

var (
	ErrHashedPasswordRequired           = errors.New("hashed password required")
	ErrHashedPasswordMustNotBePlainText = errors.New("hashed password must not be plain text")
	ErrHashedPasswordInvalid            = errors.New("hashed password format invalid")
)

type HashedPassword struct {
	value string
}

func NewHashedPassword(hashedPassword string) (HashedPassword, error) {
	if hashedPassword == "" {
		return HashedPassword{}, ErrHashedPasswordRequired
	}
	if strings.TrimSpace(hashedPassword) != hashedPassword {
		return HashedPassword{}, ErrHashedPasswordInvalid
	}
	if len(hashedPassword) < 20 {
		return HashedPassword{}, ErrHashedPasswordMustNotBePlainText
	}

	return HashedPassword{value: hashedPassword}, nil
}

func (p HashedPassword) Value() string { return p.value }
