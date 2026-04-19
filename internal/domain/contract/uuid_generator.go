package contract

type UuidGenerator interface {
	Generate() (uuid string, err error)
}
