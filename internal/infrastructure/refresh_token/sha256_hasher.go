package refreshtoken

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
)

type Sha256RefreshTokenHasher struct{}

func NewSha256RefreshTokenHasher() *Sha256RefreshTokenHasher {
	return &Sha256RefreshTokenHasher{}
}

func (*Sha256RefreshTokenHasher) Hash(token string) (string, error) {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:]), nil
}

func (*Sha256RefreshTokenHasher) Verify(token, hashed string) (bool, error) {
	hash := sha256.Sum256([]byte(token))
	computed := hex.EncodeToString(hash[:])

	ok := subtle.ConstantTimeCompare(
		[]byte(computed),
		[]byte(hashed),
	) == 1

	return ok, nil
}
