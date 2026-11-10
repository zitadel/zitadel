package handler

import "context"

// Init initializes the projection with the given check
type Init func(context.Context, *Check) error

type Check struct {
	Executes []func(ctx context.Context, executer Executer, projectionName string) (bool, error)
}

func (c *Check) IsNoop() bool {
	return len(c.Executes) == 0
}
