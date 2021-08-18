package handler

type Handler interface {
}

type HandlerOption func(Handler)
