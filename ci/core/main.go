package main

type Core struct{}

func (c *Core) Build(directory *Directory) *Container {
	return dag.Container().From("alpine").
		WithWorkdir("/src").
		WithMountedDirectory("./console/dist/console", directory).
		WithExec([]string{"ls", "-la", "./console/dist/console/"})
}
