package handler

var _ Projection = (*projection)(nil)

type projection struct {
	name     string
	reducers []AggregateReducer
}

// Name implements Projection
func (p *projection) Name() string {
	return p.name
}

// Reducers implements Projection
func (p *projection) Reducers() []AggregateReducer {
	return p.reducers
}
