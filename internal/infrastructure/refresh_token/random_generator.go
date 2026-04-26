package refreshtoken

import (
	"crypto/rand"
	"encoding/base64"
)

type RandomRefreshTokenFactory struct {
	length int
}

func NewRandomRefreshTokenFactory() *RandomRefreshTokenFactory {
	return &RandomRefreshTokenFactory{length: 60}
}

func (f *RandomRefreshTokenFactory) Generate() (string, error) {
	bytes := make([]byte, f.length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(bytes), nil
}
