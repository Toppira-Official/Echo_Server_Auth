package contract

import "auth/internal/domain/vo"

type PasswordEncoder interface {
	Hash(rawPassword string) (hashedPassword vo.HashedPassword, err error)
	Compare(rawPassword string, hashedPassword vo.HashedPassword) (ok bool, err error)
}
