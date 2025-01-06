package service

type IDGenerator interface {
	Generate() (id string, err error)
}
