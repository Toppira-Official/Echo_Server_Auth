package vo

import "errors"

var (
	ErrCredentialIDRequired = errors.New("user id required")
)

type CredentialID struct {
	value string
}

func NewCredentialID(id string) (CredentialID, error) {
	if id == "" {
		return CredentialID{}, ErrCredentialIDRequired
	}

	return CredentialID{value: id}, nil
}

func (u CredentialID) Value() string { return u.value }
