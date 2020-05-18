package id

type Generator interface {
	Next() (string, error)
}
