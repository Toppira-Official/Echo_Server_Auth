package contract

import "auth/internal/domain/vo"

type AccessTokenSigner interface {
	Generate(payload vo.AccessTokenPayload) (token string, err error)
	Verify(token string) (payload vo.AccessTokenPayload, err error)
}
