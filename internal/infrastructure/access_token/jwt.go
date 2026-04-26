package accesstoken

import (
	"auth/internal/domain/contract"
	"auth/internal/domain/vo"
	"errors"
	"time"

	"github.com/Ali127Dev/xerr"
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

	signed, err := token.SignedString(j.secretKey)
	if err != nil {
		return "", xerr.Wrap(
			err,
			xerr.CodeInternalError,
			xerr.WithMessage("failed to sign access token"),
		)
	}

	return signed, nil
}

func (s *JwtAccessTokenSigner) Verify(tokenString string) (vo.AccessTokenPayload, error) {
	var claims jwt.RegisteredClaims

	token, err := jwt.ParseWithClaims(
		tokenString, &claims, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, xerr.New(
					xerr.CodeInvalidToken,
					xerr.WithMessage("unexpected signing method"),
				)
			}
			return s.secretKey, nil
		},
		jwt.WithExpirationRequired(),
		jwt.WithIssuedAt(),
		jwt.WithTimeFunc(s.clock.NowUTC),
		jwt.WithLeeway(2*time.Second),
	)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return vo.AccessTokenPayload{}, xerr.New(
				xerr.CodeExpiredToken,
				xerr.WithMessage("access token has expired"),
			)
		}

		return vo.AccessTokenPayload{}, xerr.Wrap(
			err,
			xerr.CodeInvalidToken,
			xerr.WithMessage("failed to parse access token"),
		)
	}

	if !token.Valid {
		return vo.AccessTokenPayload{}, xerr.New(
			xerr.CodeInvalidToken,
			xerr.WithMessage("invalid access token"),
		)
	}

	credentialID, err := vo.NewCredentialID(claims.ID)
	if err != nil {
		return vo.AccessTokenPayload{}, xerr.Wrap(
			err,
			xerr.CodeInvalidToken,
			xerr.WithMessage("invalid credential id in token"),
		)
	}

	issuedAt := claims.IssuedAt.Time.UTC()
	expiresAt := claims.ExpiresAt.Time.UTC()

	payload, err := vo.NewAccessTokenPayload(credentialID, issuedAt, expiresAt)
	if err != nil {
		return vo.AccessTokenPayload{}, xerr.Wrap(
			err,
			xerr.CodeInvalidToken,
			xerr.WithMessage("invalid token payload"),
		)
	}

	return payload, nil
}
