package accesstoken

import (
	"auth/internal/domain/contract"
	"auth/internal/domain/vo"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtAccessTokenSignerConfig struct {
	SecretKey []byte
}

type JwtAccessTokenSigner struct {
	secretKey []byte
	clock     contract.Clock
}

func NewJwtAccessTokenSigner(
	cfg JwtAccessTokenSignerConfig,
	clock contract.Clock,
) *JwtAccessTokenSigner {
	return &JwtAccessTokenSigner{
		secretKey: cfg.SecretKey,
		clock:     clock,
	}
}

func (j *JwtAccessTokenSigner) Generate(payload vo.AccessTokenPayload) (string, error) {
	claims := jwt.RegisteredClaims{
		ID:        payload.CredentialID().Value(),
		Subject:   payload.CredentialID().Value(),
		IssuedAt:  jwt.NewNumericDate(payload.IssuedAtUTC()),
		ExpiresAt: jwt.NewNumericDate(payload.ExpiredAtUTC()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}

func (s *JwtAccessTokenSigner) Verify(tokenString string) (vo.AccessTokenPayload, error) {
	var claims jwt.RegisteredClaims

	token, err := jwt.ParseWithClaims(
		tokenString, &claims, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return s.secretKey, nil
		},
		jwt.WithExpirationRequired(),
		jwt.WithIssuedAt(),
		jwt.WithTimeFunc(s.clock.NowUTC),
		jwt.WithLeeway(2*time.Second),
	)
	if err != nil {
		return vo.AccessTokenPayload{}, err
	}

	if !token.Valid {
		return vo.AccessTokenPayload{}, errors.New("invalid token")
	}

	credentialID, err := vo.NewCredentialID(claims.ID)
	if err != nil {
		return vo.AccessTokenPayload{}, err
	}

	issuedAt := claims.IssuedAt.Time.UTC()
	expiresAt := claims.ExpiresAt.Time.UTC()

	payload, err := vo.NewAccessTokenPayload(credentialID, issuedAt, expiresAt)
	if err != nil {
		return vo.AccessTokenPayload{}, err
	}

	return payload, nil
}
