package spooler

import (
	"math/rand"
	"os"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/query"
	"github.com/caos/zitadel/internal/id"
)

type Config struct {
	Eventstore   eventstore.Eventstore
	Locker       Locker
	ViewHandlers []query.Handler
}

func (c *Config) New() *Spooler {
	lockID, err := os.Hostname()
	if err != nil || lockID == "" {
		lockID, err = id.SonyFlakeGenerator.Next()
		logging.Log("SPOOL-bdO56").OnError(err).Panic("unable to generate lockID")
	}

	rand.Shuffle(len(c.ViewHandlers), func(i, j int) {
		c.ViewHandlers[i], c.ViewHandlers[j] = c.ViewHandlers[j], c.ViewHandlers[i]
	})

	return &Spooler{
		handlers:   c.ViewHandlers,
		lockID:     lockID,
		eventstore: c.Eventstore,
		locker:     c.Locker,
		workers:    len(c.ViewHandlers),
	}
}
