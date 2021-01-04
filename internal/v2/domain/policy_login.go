package domain

type PasswordlessType int32

const (
	PasswordlessTypeNotAllowed PasswordlessType = iota
	PasswordlessTypeAllowed

	passwordlessCount
)

func (f PasswordlessType) Valid() bool {
	return f >= 0 && f < passwordlessCount
}
