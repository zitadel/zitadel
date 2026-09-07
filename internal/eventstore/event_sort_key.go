package eventstore

import "github.com/shopspring/decimal"

// EventSortKey is the events2 cursor: position, in_tx_order, aggregate_type, aggregate_id, sequence.
type EventSortKey struct {
	Position      decimal.Decimal
	InTxOrder     uint32
	AggregateType AggregateType
	AggregateID   string
	Sequence      uint64
}

func EventSortKeyFromEvent(e Event) EventSortKey {
	return EventSortKey{
		Position:      e.Position(),
		InTxOrder:     e.InTxOrder(),
		AggregateType: e.Aggregate().Type,
		AggregateID:   e.Aggregate().ID,
		Sequence:      e.Sequence(),
	}
}

func (k EventSortKey) IsZero() bool {
	return k.Position.IsZero()
}

// IdentityEquals reports whether k and o are the same events2 row.
// InTxOrder is ignored so a pre-migration filter_offset of 0 still matches.
func (k EventSortKey) IdentityEquals(o EventSortKey) bool {
	return k.Position.Equal(o.Position) &&
		k.AggregateType == o.AggregateType &&
		k.AggregateID == o.AggregateID &&
		k.Sequence == o.Sequence
}
