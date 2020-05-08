package spooler

import (
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/sony/sonyflake"
	"strconv"
)

type Config struct {
	Eventstore      eventstore.Eventstore
	Locker          Locker
	ViewHandlers    []Handler
	ConcurrentTasks int
}

func (c *Config) New() *Spooler {
	lockID, err := sonyflake.NewSonyflake(sonyflake.Settings{}).NextID()
	logging.Log("SPOOL-bdO56").OnError(err).Panic("unable to generate lockID")

	return &Spooler{
		handlers:        c.ViewHandlers,
		lockID:          strconv.FormatUint(lockID, 10),
		eventstore:      c.Eventstore,
		locker:          c.Locker,
		queue:           make(chan *spooledHandler),
		concurrentTasks: c.ConcurrentTasks,
	}
}
