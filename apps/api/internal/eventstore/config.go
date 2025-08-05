package eventstore

import (
	"time"
)

type Config struct {
	PushTimeout time.Duration
	MaxRetries  uint32

	Pusher   Pusher
	Querier  Querier
	Searcher Searcher
}
