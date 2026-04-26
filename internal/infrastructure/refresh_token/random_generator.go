package refreshtoken

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/Ali127Dev/xerr"
)

type RandomRefreshTokenFactory struct {
	length int
}

func NewRandomRefreshTokenFactory() *RandomRefreshTokenFactory {
	return &RandomRefreshTokenFactory{length: 60}
}

func (f *RandomRefreshTokenFactory) Generate() (string, error) {
	bytes := make([]byte, f.length)
	if _, err := rand.Read(bytes); err != nil {
		return "", xerr.Wrap(
			err,
			xerr.CodeInternalError,
			xerr.WithMessage("failed to generate refresh token"),
			xerr.WithMeta("length", f.length),
		)
	}

	return base64.RawURLEncoding.EncodeToString(bytes), nil
}
