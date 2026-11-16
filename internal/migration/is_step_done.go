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

// IsStepDone reports whether a setup step has a Done event. A missing events
// table is not done. Started without Done or Failed is an error. Failed (e.g.
// after `zitadel setup cleanup`) is not done: cursor mode requires a Done event.
func IsStepDone(ctx context.Context, es eventQuerier, name string) (bool, error) {
	if es == nil {
		return false, nil
	}
	states := new(StepStates)
	err := es.FilterToQueryReducer(ctx, states)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "42P01" {
			return false, nil
		}
		return false, err
	}
	step := states.byName(name)
	if step == nil {
		return false, nil
	}
	switch step.state {
	case StepDone:
		return true, nil
	case StepStarted:
		return false, fmt.Errorf("setup %s started without a done or failed event; filter_offset may already hold in_tx_order", name)
	case StepFailed:
		return false, nil
	default:
		return false, nil
	}
}
