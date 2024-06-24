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
			if pusherKey != "v3(postgres)" {
				continue
			}
			b.Run(fmt.Sprintf("Benchmark_Push_SameAggregate-%s-%s", pusherKey, cmdsKey), func(b *testing.B) {
				b.StopTimer()
				cleanupEventstore(clients[pusherKey])
				b.StartTimer()

				var errorCount int

				for n := 0; n < b.N; n++ {
					_, err := store.Push(ctx, cmds...)
					if err != nil {
						errorCount++
						// b.Error(err)
					}
				}
				b.ReportMetric(float64(errorCount), "error_count")
				b.ReportMetric(float64(b.Elapsed().Nanoseconds()), "elapsed_ns")
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
			if pusherKey != "v3(postgres)" {
				continue
			}
			b.Run(fmt.Sprintf("Benchmark_Push_DifferentAggregate-%s-%s", cmdsKey, pusherKey), func(b *testing.B) {
				b.StopTimer()
				cleanupEventstore(clients[pusherKey])

				ctx, cancel := context.WithCancel(context.Background())
				b.StartTimer()

				i := 0
				var errorCount int
				var asdf int

				b.SetParallelism(8)
				b.RunParallel(func(p *testing.PB) {
					for p.Next() {
						asdf, i = i, i+1
						_, err := store.Push(ctx, commandCreator(strconv.Itoa(asdf))...)
						if err != nil {
							errorCount++
							// b.Error(err)
						}
					}
				})

				b.ReportMetric(float64(b.Elapsed().Nanoseconds()), "elapsed_ns")
				b.ReportMetric(float64(errorCount), "error_count")
				cancel()
			})
		}
	}
}
