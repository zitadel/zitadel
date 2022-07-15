package handler

import (
	"context"
	"time"
)

//Lock is used for mutex handling if needed on the projection
type Lock func(context.Context, time.Duration, ...string) <-chan error

//Unlock releases the mutex of the projection
type Unlock func(...string) error

type Scheduler struct {
}
