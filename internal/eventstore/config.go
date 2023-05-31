package eventstore

import (
	"time"
)

type Config struct {
	PushTimeout              time.Duration
	AllowOrderByCreationDate bool
	UseV2                    bool

	Pusher  Pusher
	Querier Querier
}
