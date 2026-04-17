package contract

type PasswordEncoder interface {
	Hash(rawPassword string) (hashedPassword string, err error)
	Compare(rawPassword string, hashedPassword string) (ok bool, err error)
}
