package main

type HelloDagger struct{}

func (m *HelloDagger) Build(source *Directory) *Container {
	return dag.Container().From("alpine").
		WithDirectory("/app", dag.console.build).
		WithExec([]string{"ls", "-la", "/app"})
}
