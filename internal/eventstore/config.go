package eventstore

import (
	"time"
)

type Config struct {
	PushTimeout              time.Duration
	AllowOrderByCreationDate bool

	Pusher  Pusher
	Querier Querier
}
