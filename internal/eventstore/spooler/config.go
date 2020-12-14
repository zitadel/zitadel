package spooler

import (
	"os"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/query"
	"github.com/caos/zitadel/internal/id"
)

type Config struct {
	Eventstore        eventstore.Eventstore
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

	return &Spooler{
		handlers:   c.ViewHandlers,
		lockID:     lockID,
		eventstore: c.Eventstore,
		locker:     c.Locker,
		queue:      make(chan *spooledHandler, len(c.ViewHandlers)),
		workers:    c.ConcurrentWorkers,
	}
}
