package sql

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
)

func Test_subQueriesOfSingleAggregate(t *testing.T) {
	aggregate := func(typ eventstore.AggregateType, id string, more ...*repository.Filter) []*repository.Filter {
		return append([]*repository.Filter{
			repository.NewFilter(repository.FieldAggregateType, typ, repository.OperationEquals),
			repository.NewFilter(repository.FieldAggregateID, id, repository.OperationEquals),
		}, more...)
	}
	tests := []struct {
		name       string
		subQueries [][]*repository.Filter
		want       bool
	}{
		{
			name: "same aggregate",
			subQueries: [][]*repository.Filter{
				aggregate("user", "1"),
				aggregate("user", "1", repository.NewFilter(repository.FieldSequence, uint64(4), repository.OperationGreater)),
			},
			want: true,
		},
		{
			name:       "different aggregate id",
			subQueries: [][]*repository.Filter{aggregate("user", "1"), aggregate("user", "2")},
			want:       false,
		},
		{
			name:       "different aggregate type",
			subQueries: [][]*repository.Filter{aggregate("user", "1"), aggregate("org", "1")},
			want:       false,
		},
		{
			name: "missing aggregate id",
			subQueries: [][]*repository.Filter{
				aggregate("user", "1"),
				{repository.NewFilter(repository.FieldAggregateType, eventstore.AggregateType("user"), repository.OperationEquals)},
			},
			want: false,
		},
		{
			name: "multiple aggregate ids",
			subQueries: [][]*repository.Filter{
				aggregate("user", "1"),
				{
					repository.NewFilter(repository.FieldAggregateType, eventstore.AggregateType("user"), repository.OperationEquals),
					repository.NewFilter(repository.FieldAggregateID, []string{"1", "2"}, repository.OperationIn),
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, subQueriesOfSingleAggregate(tt.subQueries))
		})
	}
}
