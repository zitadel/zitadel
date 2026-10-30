package migration

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/zitadel/zitadel/internal/eventstore"
)

type eventQuerier interface {
	FilterToQueryReducer(ctx context.Context, reducer eventstore.QueryReducer) error
}

type stepDoneCheck struct {
	name    string
	started bool
	done    bool
	failed  bool
}

func (s *stepDoneCheck) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch event.Type() {
		case StartedType:
			s.started = true
			s.done = false
			s.failed = false
		case DoneType:
			s.done = true
			s.failed = false
		case failedType:
			s.failed = true
			s.done = false
		}
	}
}

func (s *stepDoneCheck) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		InstanceID("").
		AddQuery().
		AggregateTypes(SystemAggregate).
		AggregateIDs(SystemAggregateID).
		EventTypes(StartedType, DoneType, failedType).
		EventData(map[string]interface{}{
			"name": s.name,
		}).
		Builder()
}

func (*stepDoneCheck) Reduce() error {
	return nil
}

var _ eventstore.QueryReducer = (*stepDoneCheck)(nil)

// IsStepDone reports whether a setup step has a Done event. A missing events
// table is not done. Started without Done or Failed is an error.
func IsStepDone(ctx context.Context, es eventQuerier, name string) (bool, error) {
	if es == nil {
		return false, nil
	}
	check := &stepDoneCheck{name: name}
	err := es.FilterToQueryReducer(ctx, check)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "42P01" {
			return false, nil
		}
		return false, err
	}
	if check.done {
		return true, nil
	}
	if check.started && !check.failed {
		return false, fmt.Errorf("setup %s started without a done or failed event; filter_offset may already hold in_tx_order", name)
	}
	return false, nil
}
