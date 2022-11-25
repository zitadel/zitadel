package spooler

import (
	"math/rand"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/eventstore"
	v1 "github.com/zitadel/zitadel/internal/eventstore/v1"
	"github.com/zitadel/zitadel/internal/eventstore/v1/query"
	"github.com/zitadel/zitadel/internal/id"
)

type Config struct {
	Eventstore          v1.Eventstore
	EventstoreV2        *eventstore.Eventstore
	Locker              Locker
	ViewHandlers        []query.Handler
	ConcurrentWorkers   int
	ConcurrentInstances int
}

func (c *Config) New() *Spooler {
	lockID, err := id.SonyFlakeGenerator().Next()
	logging.OnError(err).Panic("unable to generate lockID")

	//shuffle the handlers for better balance when running multiple pods
	rand.Shuffle(len(c.ViewHandlers), func(i, j int) {
		c.ViewHandlers[i], c.ViewHandlers[j] = c.ViewHandlers[j], c.ViewHandlers[i]
	})

	return &Spooler{
		handlers:            c.ViewHandlers,
		lockID:              lockID,
		eventstore:          c.Eventstore,
		esV2:                c.EventstoreV2,
		locker:              c.Locker,
		queue:               make(chan *spooledHandler, len(c.ViewHandlers)),
		workers:             c.ConcurrentWorkers,
		concurrentInstances: c.ConcurrentInstances,
	}
}
