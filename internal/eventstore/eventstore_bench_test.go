package eventstore_test

import (
	"context"
	_ "embed"
	"fmt"
	"testing"

	"github.com/zitadel/zitadel/internal/eventstore"
)

//go:embed bench_payload.txt
var text string

func Benchmark_Push_SameAggregate(b *testing.B) {
	ctx := context.Background()

	smallPayload := struct {
		Username  string
		Firstname string
		Lastname  string
	}{
		Username:  "username",
		Firstname: "firstname",
		Lastname:  "lastname",
	}

	bigPayload := struct {
		Username  string
		Firstname string
		Lastname  string
		Text      string
	}{
		Username:  "username",
		Firstname: "firstname",
		Lastname:  "lastname",
		Text:      text,
	}

	commands := map[string][]eventstore.Command{
		"no payload one command": {
			generateCommand(eventstore.AggregateType(b.Name()), "id"),
		},
		"small payload one command": {
			generateCommand(eventstore.AggregateType(b.Name()), "id", withTestData(smallPayload)),
		},
		"big payload one command": {
			generateCommand(eventstore.AggregateType(b.Name()), "id", withTestData(bigPayload)),
		},
		"no payload multiple commands": {
			generateCommand(eventstore.AggregateType(b.Name()), "id"),
			generateCommand(eventstore.AggregateType(b.Name()), "id"),
			generateCommand(eventstore.AggregateType(b.Name()), "id"),
		},
		"mixed payload multiple command": {
			generateCommand(eventstore.AggregateType(b.Name()), "id", withTestData(smallPayload)),
			generateCommand(eventstore.AggregateType(b.Name()), "id", withTestData(bigPayload)),
			generateCommand(eventstore.AggregateType(b.Name()), "id", withTestData(smallPayload)),
			generateCommand(eventstore.AggregateType(b.Name()), "id", withTestData(bigPayload)),
		},
	}

	for cmdsKey, cmds := range commands {
		for pusherKey, store := range pushers {
			b.Run(fmt.Sprintf("Benchmark_Push_SameAggregate-%s-%s", pusherKey, cmdsKey), func(b *testing.B) {
				b.StopTimer()
				cleanupEventstore()
				b.StartTimer()

				for n := 0; n < b.N; n++ {
					_, err := store.Push(ctx, cmds...)
					if err != nil {
						b.Error(err)
					}
				}
			})
		}
	}
}
