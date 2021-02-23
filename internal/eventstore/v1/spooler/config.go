package spooler

import (
	"github.com/caos/zitadel/internal/eventstore/v1"
	"math/rand"
	"os"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore/v1/query"
	"github.com/caos/zitadel/internal/id"
)

type Config struct {
	Eventstore        v1.Eventstore
	Locker            Locker
	ViewHandlers      []query.Handler
	ConcurrentWorkers int
}

func (c *Config) New() *Spooler {
	lockID, err := os.Hostname()
	if err != nil || lockID == "" {
		lockID, err = id.SonyFlakeGenerator.Next()
		logging.Log("SPOOL-bdO56").OnError(err).Panic("unable to generate lockID")
	}

	//shuffle the handlers for better balance when running multiple pods
	rand.Shuffle(len(c.ViewHandlers), func(i, j int) {
		c.ViewHandlers[i], c.ViewHandlers[j] = c.ViewHandlers[j], c.ViewHandlers[i]
	})

	return &Spooler{
		handlers:   c.ViewHandlers,
		lockID:     lockID,
		eventstore: c.Eventstore,
		locker:     c.Locker,
		queue:      make(chan *spooledHandler, len(c.ViewHandlers)),
		workers:    c.ConcurrentWorkers,
	}
}
