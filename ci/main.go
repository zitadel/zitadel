package main

import (
	"context"
)

type Ci struct{}

// Build the whole stack
func (c *Ci) Build(ctx context.Context, directory *Directory) error {
	vote := dag.Docs().Build(directory.Directory("."))
	_, err := vote.Sync(ctx)
	if err != nil {
		return err
	}

	return nil
}
