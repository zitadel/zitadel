package main

import (
	"context"
)

type Ci struct{}

// Build the whole stack
func (c *Ci) Build(ctx context.Context, directory *Directory) error {
	console := dag.Console().Build(directory)
	_, err := console.Sync(ctx)
	if err != nil {
		return err
	}

	core := dag.core().Build(directory)
	_, err := core.Sync(ctx)
	if err != nil {
		return err
	}
	return nil
}
