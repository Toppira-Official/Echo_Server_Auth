package password

import (
	"auth/internal/domain/vo"
	"errors"

	"github.com/Ali127Dev/xerr"
	"golang.org/x/crypto/bcrypt"
)

type BcryptPasswordEncoder struct{}

func NewBcryptPasswordEncoder() *BcryptPasswordEncoder { return &BcryptPasswordEncoder{} }

func (b *BcryptPasswordEncoder) Hash(rawPassword string) (vo.HashedPassword, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		return vo.HashedPassword{}, xerr.Wrap(
			err,
			xerr.CodeInternalError,
			xerr.WithMessage("failed to hash password"),
			xerr.WithMeta("algo", "bcrypt"),
		)
	}

	return vo.NewHashedPassword(string(hash))
}

func (b *BcryptPasswordEncoder) Compare(rawPassword string, hashedPassword vo.HashedPassword) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword.Value()), []byte(rawPassword))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}

		return false, xerr.Wrap(
			err,
			xerr.CodeInternalError,
			xerr.WithMessage("failed to compare password hash"),
			xerr.WithMeta("algo", "bcrypt"),
		)
	}

	return true, nil
}
