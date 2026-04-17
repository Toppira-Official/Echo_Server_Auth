package contract

import "auth/internal/domain/vo"

type AuthTokenSigner interface {
	Generate(payload vo.AuthTokenPayload) (token string, err error)
	Verify(token string) (payload vo.AuthTokenPayload, err error)
}
