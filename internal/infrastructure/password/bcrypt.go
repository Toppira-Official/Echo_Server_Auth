package password

import (
	"auth/internal/domain/vo"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type BcryptPasswordEncoder struct{}

func NewBcryptPasswordEncoder() *BcryptPasswordEncoder { return &BcryptPasswordEncoder{} }

func (b *BcryptPasswordEncoder) Hash(rawPassword string) (vo.HashedPassword, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		return vo.HashedPassword{}, err
	}

	return vo.NewHashedPassword(string(hash))
}

func (b *BcryptPasswordEncoder) Compare(rawPassword string, hashedPassword vo.HashedPassword) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword.Value()), []byte(rawPassword))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
