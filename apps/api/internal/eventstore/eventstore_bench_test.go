package eventstore_test

import (
	"context"
	_ "embed"
	"fmt"
	"strconv"
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
				cleanupEventstore(clients[pusherKey])
				b.StartTimer()

				for n := 0; n < b.N; n++ {
					_, err := store.Push(ctx, store.Client().DB, cmds...)
					if err != nil {
						b.Error(err)
					}
				}
			})
		}
	}
}

func Benchmark_Push_MultipleAggregate_Parallel(b *testing.B) {
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

	commandCreators := map[string]func(id string) []eventstore.Command{
		"no payload one command": func(id string) []eventstore.Command {
			return []eventstore.Command{
				generateCommand(eventstore.AggregateType(b.Name()), id),
			}
		},
		"small payload one command": func(id string) []eventstore.Command {
			return []eventstore.Command{
				generateCommand(eventstore.AggregateType(b.Name()), id, withTestData(smallPayload)),
			}
		},
		"big payload one command": func(id string) []eventstore.Command {
			return []eventstore.Command{
				generateCommand(eventstore.AggregateType(b.Name()), id, withTestData(bigPayload)),
			}
		},
		"no payload multiple commands": func(id string) []eventstore.Command {
			return []eventstore.Command{
				generateCommand(eventstore.AggregateType(b.Name()), id),
				generateCommand(eventstore.AggregateType(b.Name()), id),
				generateCommand(eventstore.AggregateType(b.Name()), id),
			}
		},
		"mixed payload multiple command": func(id string) []eventstore.Command {
			return []eventstore.Command{
				generateCommand(eventstore.AggregateType(b.Name()), id, withTestData(smallPayload)),
				generateCommand(eventstore.AggregateType(b.Name()), id, withTestData(bigPayload)),
				generateCommand(eventstore.AggregateType(b.Name()), id, withTestData(smallPayload)),
				generateCommand(eventstore.AggregateType(b.Name()), id, withTestData(bigPayload)),
			}
		},
	}

	for cmdsKey, commandCreator := range commandCreators {
		for pusherKey, store := range pushers {
			b.Run(fmt.Sprintf("Benchmark_Push_DifferentAggregate-%s-%s", cmdsKey, pusherKey), func(b *testing.B) {
				b.StopTimer()
				cleanupEventstore(clients[pusherKey])

				ctx, cancel := context.WithCancel(context.Background())
				b.StartTimer()

				i := 0

				b.RunParallel(func(p *testing.PB) {
					for p.Next() {
						i++
						_, err := store.Push(ctx, store.Client().DB, commandCreator(strconv.Itoa(i))...)
						if err != nil {
							b.Error(err)
						}
					}
				})
				cancel()
			})
		}
	}
}
