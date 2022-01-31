package types

import (
	"time"
)

type Duration struct {
	time.Duration
}

func (d *Duration) UnmarshalText(data []byte) (err error) {
	d.Duration, err = time.ParseDuration(string(data))
	return err
}
