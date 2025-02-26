package repository

type Handler interface {
	SetNext(next Handler) Handler
}
