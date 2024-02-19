package eventstore

import "github.com/zitadel/zitadel/internal/v2/database"

type EventQuery struct {
	Exts []EventQueryExt
}

func getExt[T EventQueryExt](b *EventQuery) *T {
	for _, ext := range b.Exts {
		if t, ok := ext.(*T); ok {
			return t
		}
	}

	ext := new(T)
	b.Exts = append(b.Exts, ext)
	return ext
}

type eventQueryOpt func(b *EventQuery)

type EventQueryExt interface {
	// queryExt()
}

type AggregateTypesFilter struct {
	types []string
}

func (at *AggregateTypesFilter) Types() []string {
	return at.types
}

func FilterAggregateTypes(types ...string) eventQueryOpt {
	return func(b *EventQuery) {
		ext := getExt[AggregateTypesFilter](b)
		ext.types = append(ext.types, types...)
	}
}

type AggregateIDsFilter struct {
	ids []string
}

func (at *AggregateIDsFilter) IDs() []string {
	return at.ids
}

func FilterAggregateIDs(types ...string) eventQueryOpt {
	return func(b *EventQuery) {
		ext := getExt[AggregateIDsFilter](b)
		ext.ids = append(ext.ids, types...)
	}
}

type EventTypesFilter struct {
	types []string
}

func (at *EventTypesFilter) Types() []string {
	return at.types
}

func FilterEventTypes(types ...string) eventQueryOpt {
	return func(b *EventQuery) {
		ext := getExt[EventTypesFilter](b)
		ext.types = append(ext.types, types...)
	}
}

type SequenceFilter[T SequenceFilterType] struct {
	filter *T
}

func (s *SequenceFilter[T]) Filter() *T {
	return s.filter
}

func FilterSequence[T SequenceFilterType](filter *T) eventQueryOpt {
	return func(b *EventQuery) {
		ext := getExt[SequenceFilter[T]](b)
		ext.filter = filter
	}
}

type SequenceFilterType interface {
	SequenceEqualsType | SequenceAtLeastType | SequenceGreaterType | SequenceAtMostType | SequenceLessType | SequenceBetweenType
}

type SequenceEqualsType database.NumberFilter[uint32]

type SequenceAtLeastType database.NumberFilter[uint32]

type SequenceGreaterType database.NumberFilter[uint32]

type SequenceAtMostType database.NumberFilter[uint32]

type SequenceLessType database.NumberFilter[uint32]

type SequenceBetweenType struct {
	database.NumberBetweenFilter[uint32]
}
