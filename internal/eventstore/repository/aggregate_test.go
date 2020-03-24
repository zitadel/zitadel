package repository

import "github.com/caos/eventstore-lib/pkg/models"

type mockAggregateRoot struct {
	typ string
	id  string
}

func (agg *mockAggregateRoot) ID() string {
	return agg.id
}

func (agg *mockAggregateRoot) Type() string {
	return agg.typ
}

type mockDefaultAggregate struct {
	mockAggregateRoot
}

func (agg *mockDefaultAggregate) Events() models.Events {
	return nil
}

func (agg *mockDefaultAggregate) LatestSequence() uint64 {
	return 0
}

type mockValidateSequenceAggregate struct {
	mockAggregateRoot
	latestSequence uint64
}

func (agg mockValidateSequenceAggregate) Events() models.Events {
	return nil
}
func (agg mockValidateSequenceAggregate) LatestSequence() uint64 {
	return agg.latestSequence
}
