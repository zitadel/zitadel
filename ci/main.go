package main

import (
	"context"
)

type Ci struct{}

// Build the whole stack
func (c *Ci) Build(ctx context.Context, directory *Directory) (*Directory, error) {
	output := dag.Console().Build(directory)
	_, err := output.Sync(ctx)
	if err != nil {
		return nil, err
	}

	// core := dag.Core().Build(console)
	// _, err = core.Sync(ctx)
	// if err != nil {
	// 	return err
	// }
	return output, err

}
