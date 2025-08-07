package preparation

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

// Validation of the input values of the command and if correct returns
// the function to create commands or if not valid an error
type Validation func() (CreateCommands, error)

// CreateCommands builds the commands
// the filter param is an extended version of the eventstore filter method
// it filters for events including the commands on the current context
type CreateCommands func(context.Context, FilterToQueryReducer) ([]eventstore.Command, error)

// FilterToQueryReducer is an abstraction of the eventstore method
type FilterToQueryReducer func(ctx context.Context, queryFactory *eventstore.SearchQueryBuilder) ([]eventstore.Event, error)

var (
	//ErrNotExecutable is thrown if no command creator was created
	ErrNotExecutable = zerrors.ThrowInvalidArgument(nil, "PREPA-pH70n", "Errors.Internal")
)

// PrepareCommands checks the passed validations and if ok creates the commands
//
// Deprecated: filter causes unneeded allocation. Use [eventstore.FilterToQueryReducer] instead.
func PrepareCommands(ctx context.Context, filter FilterToQueryReducer, validations ...Validation) (cmds []eventstore.Command, err error) {
	commanders, err := validate(validations)
	if err != nil {
		return nil, err
	}
	return create(ctx, filter, commanders)
}

func validate(validations []Validation) ([]CreateCommands, error) {
	creators := make([]CreateCommands, 0, len(validations))

	for _, validate := range validations {
		cmds, err := validate()
		if err != nil {
			return nil, err
		}
		creators = append(creators, cmds)
	}

	if len(creators) == 0 {
		return nil, ErrNotExecutable
	}
	return creators, nil
}

func create(ctx context.Context, filter FilterToQueryReducer, commanders []CreateCommands) (cmds []eventstore.Command, err error) {
	for _, command := range commanders {
		cmd, err := command(ctx, transactionFilter(filter, cmds))
		if err != nil {
			return nil, err
		}
		cmds = append(cmds, cmd...)
	}

	return cmds, nil
}

func transactionFilter(filter FilterToQueryReducer, commands []eventstore.Command) FilterToQueryReducer {
	return func(ctx context.Context, query *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
		events, err := filter(ctx, query)
		if err != nil {
			return nil, err
		}
		matches := query.Matches(commands...)
		for _, command := range matches {
			events = append(events, command.(eventstore.Event))
		}
		return events, nil
	}
}
