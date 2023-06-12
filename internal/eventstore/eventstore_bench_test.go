package eventstore

import (
	"context"
	"testing"
)

func Benchmark_Push_SameAggregate(b *testing.B) {
	// TODO: before
	var es *Eventstore
	ctx := context.Background()()
	b.Run("Benchmark_Push_SameAggregate", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			newTestEvent()
			_, err := es.Push(ctx)
			if err != nil {
				b.Error(err)
			}
		}
	})
}
