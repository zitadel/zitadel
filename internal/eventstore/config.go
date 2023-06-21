package eventstore

import (
	"time"
)

type Config struct {
	PushTimeout time.Duration

	Pusher  Pusher
	Querier Querier
}
