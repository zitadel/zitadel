package repository

type aggregateRoot interface {
	ID() string
	Type() string
}
