package main

import (
	"context"
)

type Ci struct{}

// Build the whole stack
func (c *Ci) Build(ctx context.Context, directory *Directory) (*Directory, error) {
	console := dag.Console().Build(directory.Directory("."))
	_, err := console.Sync(ctx)
	if err != nil {
		return nil, err
	}

	return console, nil
}

func (c *Ci) Test(ctx context.Context, directory *Directory) (*Directory, error) {
	console := dag.Console().Build(directory)
	return console, nil
}
