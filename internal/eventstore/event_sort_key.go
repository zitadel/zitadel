package eventstore

import "github.com/shopspring/decimal"

// EventSortKey is the events2 cursor: position, in_tx_order, instance_id, aggregate_type, aggregate_id, sequence.
type EventSortKey struct {
	Position      decimal.Decimal
	InTxOrder     uint32
	InstanceID    string
	AggregateType AggregateType
	AggregateID   string
	Sequence      uint64
}

func EventSortKeyFromEvent(e Event) EventSortKey {
	return EventSortKey{
		Position:      e.Position(),
		InTxOrder:     e.InTxOrder(),
		InstanceID:    e.Aggregate().InstanceID,
		AggregateType: e.Aggregate().Type,
		AggregateID:   e.Aggregate().ID,
		Sequence:      e.Sequence(),
	}
}

func (k EventSortKey) IsZero() bool {
	return k.Position.IsZero()
}
