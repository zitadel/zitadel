package eventsourcing

import (
	"github.com/sony/sonyflake"
)

var idGenerator = sonyflake.NewSonyflake(sonyflake.Settings{})
