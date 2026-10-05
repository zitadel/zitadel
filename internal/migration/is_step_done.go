package migration

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/zitadel/zitadel/internal/eventstore"
)

type eventQuerier interface {
	FilterToQueryReducer(ctx context.Context, reducer eventstore.QueryReducer) error
}

type stepDoneCheck struct {
	name string
	done bool
}

func (s *stepDoneCheck) AppendEvents(events ...eventstore.Event) {
	if len(events) > 0 {
		s.done = true
	}
}

func (s *stepDoneCheck) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		Limit(1).
		InstanceID("").
		AddQuery().
		AggregateTypes(SystemAggregate).
		AggregateIDs(SystemAggregateID).
		EventTypes(DoneType).
		EventData(map[string]interface{}{
			"name": s.name,
		}).
		Builder()
}

func (*stepDoneCheck) Reduce() error {
	return nil
}

var _ eventstore.QueryReducer = (*stepDoneCheck)(nil)

// IsStepDone reports whether a setup step with the given name has a Done event.
// A missing events table (fresh setup before 14) is treated as not done.
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
	return check.done, nil
}
