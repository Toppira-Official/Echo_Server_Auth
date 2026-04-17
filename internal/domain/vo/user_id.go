package vo

import "errors"

var (
	ErrUserIDUserIdRequired = errors.New("user id required")
)

type UserID struct {
	value string
}

func NewUserID(id string) (UserID, error) {
	if id == "" {
		return UserID{}, ErrUserIDUserIdRequired
	}

	return UserID{value: id}, nil
}

func (u UserID) Value() string { return u.value }
