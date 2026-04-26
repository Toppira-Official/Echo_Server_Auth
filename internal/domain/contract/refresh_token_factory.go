package contract

type RefreshTokenGenerator interface {
	Generate() (refreshToken string, err error)
}

type RefreshTokenHasher interface {
	Hash(token string) (hashed string, err error)
	Verify(token, hashed string) (ok bool, err error)
}
