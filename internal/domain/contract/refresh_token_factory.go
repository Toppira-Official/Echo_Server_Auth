package contract

type RefreshTokenFactory interface {
	Generate() (refreshToken string, err error)
}
