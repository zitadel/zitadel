package eventstore

import (
	"context"
	"database/sql"
)

type Pusher interface {
	healthier
	// Push writes the intents to the storage
	// if an intent implements [PushReducerIntent] [PushReducerIntent.Reduce] is called after
	// the intent was stored
	Push(ctx context.Context, intent *PushIntent) error
}

func NewPushIntent(instance string, opts ...PushOpt) *PushIntent {
	intent := &PushIntent{
		instance: instance,
	}

	for _, opt := range opts {
		opt(intent)
	}

	return intent
}

type PushIntent struct {
	instance   string
	reducer    Reducer
	tx         *sql.Tx
	aggregates []*PushAggregate
}

func (pi *PushIntent) Instance() string {
	return pi.instance
}

func (pi *PushIntent) Reduce(events ...*StorageEvent) error {
	if pi.reducer == nil {
		return nil
	}
	return pi.reducer.Reduce(events...)
}

func (pi *PushIntent) Tx() *sql.Tx {
	return pi.tx
}

func (pi *PushIntent) Aggregates() []*PushAggregate {
	return pi.aggregates
}

type PushOpt func(pi *PushIntent)

func PushReducer(reducer Reducer) PushOpt {
	return func(pi *PushIntent) {
		pi.reducer = reducer
	}
}

func PushTx(tx *sql.Tx) PushOpt {
	return func(pi *PushIntent) {
		pi.tx = tx
	}
}

func AppendAggregate(owner, typ, id string, opts ...PushAggregateOpt) PushOpt {
	return AppendAggregates(NewPushAggregate(owner, typ, id, opts...))
}

func AppendAggregates(aggregates ...*PushAggregate) PushOpt {
	return func(pi *PushIntent) {
		for _, aggregate := range aggregates {
			aggregate.parent = pi
		}
		pi.aggregates = append(pi.aggregates, aggregates...)
	}
}

type PushAggregate struct {
	parent *PushIntent
	// typ of the aggregate
	typ string
	// id of the aggregate
	id string
	// owner of the aggregate
	owner string
	// Commands is an ordered list of changes on the aggregate
	commands []*Command
	// CurrentSequence checks the current state of the aggregate.
	// The following types match the current sequence of the aggregate as described:
	// * nil or [SequenceIgnore]: Not relevant to add the commands
	// * [SequenceMatches]: Must exactly match
	// * [SequenceAtLeast]: Must be >= the given sequence
	currentSequence CurrentSequence
}

func NewPushAggregate(owner, typ, id string, opts ...PushAggregateOpt) *PushAggregate {
	pa := &PushAggregate{
		typ:   typ,
		id:    id,
		owner: owner,
	}

	for _, opt := range opts {
		opt(pa)
	}

	return pa
}

func (pa *PushAggregate) Type() string {
	return pa.typ
}

func (pa *PushAggregate) ID() string {
	return pa.id
}

func (pa *PushAggregate) Owner() string {
	return pa.owner
}

func (pa *PushAggregate) Commands() []*Command {
	return pa.commands
}

func (pa *PushAggregate) Aggregate() *Aggregate {
	return &Aggregate{
		ID:       pa.id,
		Type:     pa.typ,
		Owner:    pa.owner,
		Instance: pa.parent.instance,
	}
}

func (pa *PushAggregate) CurrentSequence() CurrentSequence {
	return pa.currentSequence
}

type PushAggregateOpt func(pa *PushAggregate)

func SetCurrentSequence(currentSequence CurrentSequence) PushAggregateOpt {
	return func(pa *PushAggregate) {
		pa.currentSequence = currentSequence
	}
}

func IgnoreCurrentSequence() PushAggregateOpt {
	return func(pa *PushAggregate) {
		pa.currentSequence = SequenceIgnore()
	}
}

func CurrentSequenceMatches(sequence uint32) PushAggregateOpt {
	return func(pa *PushAggregate) {
		pa.currentSequence = SequenceMatches(sequence)
	}
}

func CurrentSequenceAtLeast(sequence uint32) PushAggregateOpt {
	return func(pa *PushAggregate) {
		pa.currentSequence = SequenceAtLeast(sequence)
	}
}

func AppendCommands(commands ...*Command) PushAggregateOpt {
	return func(pa *PushAggregate) {
		pa.commands = append(pa.commands, commands...)
	}
}
