package handler

import "github.com/zitadel/zitadel/internal/eventstore"

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

type mockEventReducer struct {
	wasCalled bool
	statement *Statement
}

func (m *mockEventReducer) reduce(event eventstore.Event) (*Statement, error) {
	m.wasCalled = true
	return m.statement, nil
}
