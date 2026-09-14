package eventstore_test

import (
	"context"
	"database/sql"
	"sync"
	"testing"
	"time"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestEventstore_Push_OneAggregate(t *testing.T) {
	type args struct {
		ctx                  context.Context
		commands             []eventstore.Command
		uniqueDataType       string
		uniqueDataField      string
		uniqueDataInstanceID string
	}
	type eventsRes struct {
		pushedEventsCount int
		uniqueCount       int
		assetCount        int
		aggType           eventstore.AggregateType
		aggIDs            database.TextArray[string]
	}
	type res struct {
		wantErr   bool
		eventsRes eventsRes
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "push 1 event",
			args: args{
				ctx: context.Background(),
				commands: []eventstore.Command{
					generateCommand(eventstore.AggregateType(t.Name()), "1"),
				},
			},
			res: res{
				wantErr: false,
				eventsRes: eventsRes{
					pushedEventsCount: 1,
					aggIDs:            []string{"1"},
					aggType:           eventstore.AggregateType(t.Name()),
				}},
		},
		{
			name: "push two events on agg",
			args: args{
				ctx: context.Background(),
				commands: []eventstore.Command{
					generateCommand(eventstore.AggregateType(t.Name()), "6"),
					generateCommand(eventstore.AggregateType(t.Name()), "6"),
				},
			},
			res: res{
				wantErr: false,
				eventsRes: eventsRes{
					pushedEventsCount: 2,
					aggIDs:            []string{"6"},
					aggType:           eventstore.AggregateType(t.Name()),
				},
			},
		},
		{
			name: "failed push because context canceled",
			args: args{
				ctx: canceledCtx(),
				commands: []eventstore.Command{
					generateCommand(eventstore.AggregateType(t.Name()), "9"),
				},
			},
			res: res{
				wantErr: true,
				eventsRes: eventsRes{
					pushedEventsCount: 0,
					aggIDs:            []string{"9"},
					aggType:           eventstore.AggregateType(t.Name()),
				},
			},
		},
		{
			name: "push 1 event and add unique constraint",
			args: args{
				ctx: context.Background(),
				commands: []eventstore.Command{
					generateCommand(eventstore.AggregateType(t.Name()), "10",
						generateAddUniqueConstraint("usernames", "field"),
					),
				},
			},
			res: res{
				wantErr: false,
				eventsRes: eventsRes{
					pushedEventsCount: 1,
					uniqueCount:       1,
					aggIDs:            []string{"10"},
					aggType:           eventstore.AggregateType(t.Name()),
				}},
		},
		{
			name: "push 1 event and remove unique constraint",
			args: args{
				ctx: context.Background(),
				commands: []eventstore.Command{
					generateCommand(eventstore.AggregateType(t.Name()), "11",
						generateRemoveUniqueConstraint("usernames", "testremove"),
					),
				},
				uniqueDataType:  "usernames",
				uniqueDataField: "testremove",
			},
			res: res{
				wantErr: false,
				eventsRes: eventsRes{
					pushedEventsCount: 1,
					uniqueCount:       0,
					aggIDs:            []string{"11"},
					aggType:           eventstore.AggregateType(t.Name()),
				}},
		},
		{
			name: "push 1 event and remove instance unique constraints",
			args: args{
				ctx: context.Background(),
				commands: []eventstore.Command{
					generateCommand(eventstore.AggregateType(t.Name()), "12",
						generateRemoveUniqueConstraint("instance", "instanceID"),
					),
				},
				uniqueDataType:       "usernames",
				uniqueDataField:      "testremove",
				uniqueDataInstanceID: "instanceID",
			},
			res: res{
				wantErr: false,
				eventsRes: eventsRes{
					pushedEventsCount: 1,
					uniqueCount:       0,
					aggIDs:            []string{"12"},
					aggType:           eventstore.AggregateType(t.Name()),
				}},
		},
		{
			name: "push 1 event and add asset",
			args: args{
				ctx: context.Background(),
				commands: []eventstore.Command{
					generateCommand(eventstore.AggregateType(t.Name()), "13"),
				},
			},
			res: res{
				wantErr: false,
				eventsRes: eventsRes{
					pushedEventsCount: 1,
					assetCount:        1,
					aggIDs:            []string{"13"},
					aggType:           eventstore.AggregateType(t.Name()),
				}},
		},
		{
			name: "push 1 event and remove asset",
			args: args{
				ctx: context.Background(),
				commands: []eventstore.Command{
					generateCommand(eventstore.AggregateType(t.Name()), "14"),
				},
			},
			res: res{
				wantErr: false,
				eventsRes: eventsRes{
					pushedEventsCount: 1,
					assetCount:        0,
					aggIDs:            []string{"14"},
					aggType:           eventstore.AggregateType(t.Name()),
				}},
		},
	}
	for _, tt := range tests {
		for pusherName, pusher := range pushers {
			t.Run(pusherName+"/"+tt.name, func(t *testing.T) {
				t.Cleanup(cleanupEventstore(clients[pusherName]))
				db := eventstore.NewEventstore(
					&eventstore.Config{
						Querier: queriers["v2(inmemory)"],
						Pusher:  pusher,
					},
				)

				if tt.args.uniqueDataType != "" && tt.args.uniqueDataField != "" {
					err := fillUniqueData(tt.args.uniqueDataType, tt.args.uniqueDataField, tt.args.uniqueDataInstanceID)
					if err != nil {
						t.Error("unable to prefill insert unique data: ", err)
						return
					}
				}
				if _, err := db.Push(tt.args.ctx, tt.args.commands...); (err != nil) != tt.res.wantErr {
					t.Errorf("eventstore.Push() error = %v, wantErr %v", err, tt.res.wantErr)
				}

				assertEventCount(t,
					clients[pusherName],
					database.TextArray[eventstore.AggregateType]{tt.res.eventsRes.aggType},
					tt.res.eventsRes.aggIDs,
					tt.res.eventsRes.pushedEventsCount,
				)

				assertUniqueConstraint(t, clients[pusherName], tt.args.commands, tt.res.eventsRes.uniqueCount)
			})
		}
	}
}

func TestEventstore_Push_MultipleAggregate(t *testing.T) {
	type args struct {
		commands []eventstore.Command
	}
	type eventsRes struct {
		pushedEventsCount int
		aggType           database.TextArray[eventstore.AggregateType]
		aggID             database.TextArray[string]
	}
	type res struct {
		wantErr   bool
		eventsRes eventsRes
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "push two aggregates",
			args: args{
				commands: []eventstore.Command{
					generateCommand(eventstore.AggregateType(t.Name()), "100"),
					generateCommand(eventstore.AggregateType(t.Name()), "101"),
				},
			},
			res: res{
				wantErr: false,
				eventsRes: eventsRes{
					pushedEventsCount: 2,
					aggID:             []string{"100", "101"},
					aggType:           database.TextArray[eventstore.AggregateType]{eventstore.AggregateType(t.Name())},
				},
			},
		},
		{
			name: "push two aggregates both multiple events",
			args: args{
				commands: []eventstore.Command{
					generateCommand(eventstore.AggregateType(t.Name()), "102"),
					generateCommand(eventstore.AggregateType(t.Name()), "102"),
					generateCommand(eventstore.AggregateType(t.Name()), "103"),
					generateCommand(eventstore.AggregateType(t.Name()), "103"),
				},
			},
			res: res{
				wantErr: false,
				eventsRes: eventsRes{
					pushedEventsCount: 4,
					aggID:             []string{"102", "103"},
					aggType:           database.TextArray[eventstore.AggregateType]{eventstore.AggregateType(t.Name())},
				},
			},
		},
		{
			name: "push two aggregates mixed multiple events",
			args: args{
				commands: []eventstore.Command{
					generateCommand(eventstore.AggregateType(t.Name()), "106"),
					generateCommand(eventstore.AggregateType(t.Name()), "106"),
					generateCommand(eventstore.AggregateType(t.Name()), "106"),
					generateCommand(eventstore.AggregateType(t.Name()), "106"),
					generateCommand(eventstore.AggregateType(t.Name()), "107"),
					generateCommand(eventstore.AggregateType(t.Name()), "107"),
					generateCommand(eventstore.AggregateType(t.Name()), "107"),
					generateCommand(eventstore.AggregateType(t.Name()), "107"),
					generateCommand(eventstore.AggregateType(t.Name()), "108"),
					generateCommand(eventstore.AggregateType(t.Name()), "108"),
					generateCommand(eventstore.AggregateType(t.Name()), "108"),
					generateCommand(eventstore.AggregateType(t.Name()), "108"),
				},
			},
			res: res{
				wantErr: false,
				eventsRes: eventsRes{
					pushedEventsCount: 12,
					aggID:             []string{"106", "107", "108"},
					aggType:           database.TextArray[eventstore.AggregateType]{eventstore.AggregateType(t.Name())},
				},
			},
		},
	}
	for _, tt := range tests {
		for pusherName, pusher := range pushers {
			t.Run(pusherName+"/"+tt.name, func(t *testing.T) {
				t.Cleanup(cleanupEventstore(clients[pusherName]))

				db := eventstore.NewEventstore(
					&eventstore.Config{
						Querier: queriers["v2(inmemory)"],
						Pusher:  pusher,
					},
				)
				if _, err := db.Push(context.Background(), tt.args.commands...); (err != nil) != tt.res.wantErr {
					t.Errorf("eventstore.Push() error = %v, wantErr %v", err, tt.res.wantErr)
				}

				assertEventCount(t, clients[pusherName], tt.res.eventsRes.aggType, tt.res.eventsRes.aggID, tt.res.eventsRes.pushedEventsCount)
			})
		}
	}
}

func TestEventstore_Push_Parallel(t *testing.T) {
	type args struct {
		commands [][]eventstore.Command
	}
	type eventsRes struct {
		pushedEventsCount int
		aggTypes          database.TextArray[eventstore.AggregateType]
		aggIDs            database.TextArray[string]
	}
	type res struct {
		minErrCount int
		eventsRes   eventsRes
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "clients push different aggregates",
			args: args{
				commands: [][]eventstore.Command{
					{
						generateCommand(eventstore.AggregateType(t.Name()), "200"),
						generateCommand(eventstore.AggregateType(t.Name()), "200"),
						generateCommand(eventstore.AggregateType(t.Name()), "200"),
						generateCommand(eventstore.AggregateType(t.Name()), "201"),
						generateCommand(eventstore.AggregateType(t.Name()), "201"),
						generateCommand(eventstore.AggregateType(t.Name()), "201"),
					},
					{
						generateCommand(eventstore.AggregateType(t.Name()), "202"),
						generateCommand(eventstore.AggregateType(t.Name()), "203"),
						generateCommand(eventstore.AggregateType(t.Name()), "203"),
					},
				},
			},
			res: res{
				minErrCount: 0,
				eventsRes: eventsRes{
					aggIDs:            []string{"200", "201", "202", "203"},
					pushedEventsCount: 9,
					aggTypes:          database.TextArray[eventstore.AggregateType]{eventstore.AggregateType(t.Name())},
				},
			},
		},
		{
			name: "clients push same aggregates",
			args: args{
				commands: [][]eventstore.Command{
					{
						generateCommand(eventstore.AggregateType(t.Name()), "204"),
						generateCommand(eventstore.AggregateType(t.Name()), "204"),
					},
					{
						generateCommand(eventstore.AggregateType(t.Name()), "204"),
						generateCommand(eventstore.AggregateType(t.Name()), "204"),
					},
					{
						generateCommand(eventstore.AggregateType(t.Name()), "204"),
						generateCommand(eventstore.AggregateType(t.Name()), "204"),
					},
					{
						generateCommand(eventstore.AggregateType(t.Name()), "204"),
						generateCommand(eventstore.AggregateType(t.Name()), "204"),
					},
				},
			},
			res: res{
				minErrCount: 0,
				eventsRes: eventsRes{
					aggIDs:            []string{"204"},
					pushedEventsCount: 8,
					aggTypes:          database.TextArray[eventstore.AggregateType]{eventstore.AggregateType(t.Name())},
				},
			},
		},
		{
			name: "clients push different aggregates",
			args: args{
				commands: [][]eventstore.Command{
					{
						generateCommand(eventstore.AggregateType(t.Name()), "207"),
						generateCommand(eventstore.AggregateType(t.Name()), "207"),
						generateCommand(eventstore.AggregateType(t.Name()), "207"),
						generateCommand(eventstore.AggregateType(t.Name()), "207"),
						generateCommand(eventstore.AggregateType(t.Name()), "207"),
						generateCommand(eventstore.AggregateType(t.Name()), "207"),
					},
					{
						generateCommand(eventstore.AggregateType(t.Name()), "208"),
						generateCommand(eventstore.AggregateType(t.Name()), "208"),
						generateCommand(eventstore.AggregateType(t.Name()), "208"),
						generateCommand(eventstore.AggregateType(t.Name()), "208"),
						generateCommand(eventstore.AggregateType(t.Name()), "208"),
					},
				},
			},
			res: res{
				minErrCount: 0,
				eventsRes: eventsRes{
					aggIDs:            []string{"207", "208"},
					pushedEventsCount: 11,
					aggTypes:          database.TextArray[eventstore.AggregateType]{eventstore.AggregateType(t.Name())},
				},
			},
		},
	}
	for _, tt := range tests {
		for pusherName, pusher := range pushers {
			t.Run(pusherName+"/"+tt.name, func(t *testing.T) {
				t.Cleanup(cleanupEventstore(clients[pusherName]))

				db := eventstore.NewEventstore(
					&eventstore.Config{
						Querier: queriers["v2(inmemory)"],
						Pusher:  pusher,
					},
				)

				errs := pushAggregates(db, tt.args.commands)

				if len(errs) < tt.res.minErrCount {
					t.Errorf("eventstore.Push() error count = %d, wanted err count %d, errs: %v", len(errs), tt.res.minErrCount, errs)
				}

				assertEventCount(t, clients[pusherName], tt.res.eventsRes.aggTypes, tt.res.eventsRes.aggIDs, tt.res.eventsRes.pushedEventsCount)
			})
		}
	}
}

func TestEventstore_Push_ResourceOwner(t *testing.T) {
	type args struct {
		commands []eventstore.Command
	}
	type res struct {
		resourceOwners database.TextArray[string]
	}
	type fields struct {
		aggregateIDs  database.TextArray[string]
		aggregateType string
	}
	tests := []struct {
		name   string
		args   args
		res    res
		fields fields
	}{
		{
			name: "two events of same aggregate same resource owner",
			args: args{
				commands: []eventstore.Command{
					generateCommand(eventstore.AggregateType(t.Name()), "500", func(e *testEvent) { e.BaseEvent.Agg.ResourceOwner = "caos" }),
					generateCommand(eventstore.AggregateType(t.Name()), "500", func(e *testEvent) { e.BaseEvent.Agg.ResourceOwner = "caos" }),
				},
			},
			fields: fields{
				aggregateIDs:  []string{"500"},
				aggregateType: t.Name(),
			},
			res: res{
				resourceOwners: []string{"caos", "caos"},
			},
		},
		{
			name: "two events of different aggregate same resource owner",
			args: args{
				commands: []eventstore.Command{
					generateCommand(eventstore.AggregateType(t.Name()), "501", func(e *testEvent) { e.BaseEvent.Agg.ResourceOwner = "caos" }),
					generateCommand(eventstore.AggregateType(t.Name()), "502", func(e *testEvent) { e.BaseEvent.Agg.ResourceOwner = "caos" }),
				},
			},
			fields: fields{
				aggregateIDs:  []string{"501", "502"},
				aggregateType: t.Name(),
			},
			res: res{
				resourceOwners: []string{"caos", "caos"},
			},
		},
		{
			name: "two events of different aggregate different resource owner",
			args: args{
				commands: []eventstore.Command{
					generateCommand(eventstore.AggregateType(t.Name()), "503", func(e *testEvent) { e.BaseEvent.Agg.ResourceOwner = "caos" }),
					generateCommand(eventstore.AggregateType(t.Name()), "504", func(e *testEvent) { e.BaseEvent.Agg.ResourceOwner = "zitadel" }),
				},
			},
			fields: fields{
				aggregateIDs:  []string{"503", "504"},
				aggregateType: t.Name(),
			},
			res: res{
				resourceOwners: []string{"caos", "zitadel"},
			},
		},
		{
			name: "events of different aggregate different resource owner",
			args: args{
				commands: []eventstore.Command{
					generateCommand(eventstore.AggregateType(t.Name()), "505", func(e *testEvent) { e.BaseEvent.Agg.ResourceOwner = "caos" }),
					generateCommand(eventstore.AggregateType(t.Name()), "505", func(e *testEvent) { e.BaseEvent.Agg.ResourceOwner = "caos" }),
					generateCommand(eventstore.AggregateType(t.Name()), "506", func(e *testEvent) { e.BaseEvent.Agg.ResourceOwner = "zitadel" }),
					generateCommand(eventstore.AggregateType(t.Name()), "506", func(e *testEvent) { e.BaseEvent.Agg.ResourceOwner = "zitadel" }),
				},
			},
			fields: fields{
				aggregateIDs:  []string{"505", "506"},
				aggregateType: t.Name(),
			},
			res: res{
				resourceOwners: []string{"caos", "caos", "zitadel", "zitadel"},
			},
		},
		{
			name: "events of different aggregate different resource owner per event",
			args: args{
				commands: []eventstore.Command{
					generateCommand(eventstore.AggregateType(t.Name()), "507", func(e *testEvent) { e.BaseEvent.Agg.ResourceOwner = "caos" }),
					generateCommand(eventstore.AggregateType(t.Name()), "507", func(e *testEvent) { e.BaseEvent.Agg.ResourceOwner = "ignored" }),
					generateCommand(eventstore.AggregateType(t.Name()), "508", func(e *testEvent) { e.BaseEvent.Agg.ResourceOwner = "zitadel" }),
					generateCommand(eventstore.AggregateType(t.Name()), "508", func(e *testEvent) { e.BaseEvent.Agg.ResourceOwner = "ignored" }),
				},
			},
			fields: fields{
				aggregateIDs:  []string{"507", "508"},
				aggregateType: t.Name(),
			},
			res: res{
				resourceOwners: []string{"caos", "caos", "zitadel", "zitadel"},
			},
		},
		{
			name: "events of one aggregate different resource owner per event",
			args: args{
				commands: []eventstore.Command{
					generateCommand(eventstore.AggregateType(t.Name()), "509", func(e *testEvent) { e.BaseEvent.Agg.ResourceOwner = "caos" }),
					generateCommand(eventstore.AggregateType(t.Name()), "509", func(e *testEvent) { e.BaseEvent.Agg.ResourceOwner = "ignored" }),
					generateCommand(eventstore.AggregateType(t.Name()), "509", func(e *testEvent) { e.BaseEvent.Agg.ResourceOwner = "ignored" }),
					generateCommand(eventstore.AggregateType(t.Name()), "509", func(e *testEvent) { e.BaseEvent.Agg.ResourceOwner = "ignored" }),
				},
			},
			fields: fields{
				aggregateIDs:  []string{"509"},
				aggregateType: t.Name(),
			},
			res: res{
				resourceOwners: []string{"caos", "caos", "caos", "caos"},
			},
		},
	}
	for _, tt := range tests {
		for pusherName, pusher := range pushers {
			t.Run(pusherName+"/"+tt.name, func(t *testing.T) {
				t.Cleanup(cleanupEventstore(clients[pusherName]))

				db := eventstore.NewEventstore(
					&eventstore.Config{
						Querier: queriers["v2(inmemory)"],
						Pusher:  pusher,
					},
				)

				events, err := db.Push(context.Background(), tt.args.commands...)
				if err != nil {
					t.Errorf("eventstore.Push() error = %v", err)
				}

				if len(events) != len(tt.res.resourceOwners) {
					t.Errorf("length of events (%d) and resource owners (%d) must be equal", len(events), len(tt.res.resourceOwners))
					return
				}

				for i, event := range events {
					if event.Aggregate().ResourceOwner != tt.res.resourceOwners[i] {
						t.Errorf("resource owner not expected want: %q got: %q", tt.res.resourceOwners[i], event.Aggregate().ResourceOwner)
					}
				}

				assertResourceOwners(t, clients[pusherName], tt.res.resourceOwners, tt.fields.aggregateIDs, tt.fields.aggregateType)
			})
		}
	}
}

func TestEventstore_Push_ResourceOwner_ReusedAggregateID(t *testing.T) {
	aggregateType := eventstore.AggregateType(t.Name())
	aggregateID := "510"

	for pusherName, pusher := range pushers {
		t.Run(pusherName, func(t *testing.T) {
			t.Cleanup(cleanupEventstore(clients[pusherName]))

			db := eventstore.NewEventstore(
				&eventstore.Config{
					Querier: queriers["v2(inmemory)"],
					Pusher:  pusher,
				},
			)

			initialEvents, err := db.Push(context.Background(),
				generateCommand(
					aggregateType,
					aggregateID,
					func(e *testEvent) { e.Agg.ResourceOwner = "old" },
					withEventType("test.initial"),
				),
			)
			if err != nil {
				t.Fatalf("initial eventstore.Push() error = %v", err)
			}
			if len(initialEvents) != 1 {
				t.Fatalf("initial events length = %d, want 1", len(initialEvents))
			}
			if initialEvents[0].Aggregate().ResourceOwner != "old" {
				t.Fatalf("initial resource owner = %q, want %q", initialEvents[0].Aggregate().ResourceOwner, "old")
			}

			reusedEvents, err := db.Push(context.Background(),
				generateEnforcedCommand(
					aggregateType,
					aggregateID,
					func(e *testEvent) { e.Agg.ResourceOwner = "new" },
					withEventType("test.recreated"),
				),
				generateCommand(
					aggregateType,
					aggregateID,
					func(e *testEvent) { e.Agg.ResourceOwner = "ignored" },
					withEventType("test.followup"),
				),
			)
			if err != nil {
				t.Fatalf("reused eventstore.Push() error = %v", err)
			}
			if len(reusedEvents) != 2 {
				t.Fatalf("reused events length = %d, want 2", len(reusedEvents))
			}
			for i, event := range reusedEvents {
				if event.Aggregate().ResourceOwner != "new" {
					t.Fatalf("reused event %d resource owner = %q, want %q", i, event.Aggregate().ResourceOwner, "new")
				}
			}

			assertResourceOwnersBySequence(t, clients[pusherName], []string{"old", "new", "new"}, string(aggregateType), aggregateID)
		})
	}
}

func pushAggregates(es *eventstore.Eventstore, aggregateCommands [][]eventstore.Command) []error {
	wg := sync.WaitGroup{}
	errs := make([]error, 0)
	errsMu := sync.Mutex{}
	wg.Add(len(aggregateCommands))

	ctx, cancel := context.WithCancel(context.Background())

	for _, commands := range aggregateCommands {
		go func(events []eventstore.Command) {
			<-ctx.Done()

			_, err := es.Push(context.Background(), events...) //nolint:contextcheck
			if err != nil {
				errsMu.Lock()
				errs = append(errs, err)
				errsMu.Unlock()
			}

			wg.Done()
		}(commands)
	}

	// wait till all routines are started
	time.Sleep(100 * time.Millisecond)
	cancel()
	wg.Wait()

	return errs
}

func assertResourceOwners(t *testing.T, db *database.DB, resourceOwners, aggregateIDs database.TextArray[string], aggregateType string) {
	t.Helper()

	eventCount := 0
	err := db.Query(func(rows *sql.Rows) error {
		for i := 0; rows.Next(); i++ {
			var resourceOwner string
			err := rows.Scan(&resourceOwner)
			if err != nil {
				return err
			}
			if resourceOwner != resourceOwners[i] {
				t.Errorf("unexpected resource owner in queried event. want %q, got: %q", resourceOwners[i], resourceOwner)
			}
			eventCount++
		}
		return nil
	}, "SELECT owner FROM eventstore.events2 WHERE aggregate_type = $1 AND aggregate_id = ANY($2) ORDER BY position, in_tx_order", aggregateType, aggregateIDs)
	if err != nil {
		t.Error("query failed: ", err)
		return
	}

	if eventCount != len(resourceOwners) {
		t.Errorf("wrong queried event count: want %d, got %d", len(resourceOwners), eventCount)
	}
}

func assertResourceOwnersBySequence(t *testing.T, db *database.DB, resourceOwners []string, aggregateType, aggregateID string) {
	t.Helper()

	eventCount := 0
	err := db.Query(func(rows *sql.Rows) error {
		for i := 0; rows.Next(); i++ {
			var resourceOwner string
			if err := rows.Scan(&resourceOwner); err != nil {
				return err
			}
			if resourceOwner != resourceOwners[i] {
				t.Errorf("unexpected resource owner in queried event. want %q, got: %q", resourceOwners[i], resourceOwner)
			}
			eventCount++
		}
		return nil
	}, `SELECT owner FROM eventstore.events2 WHERE aggregate_type = $1 AND aggregate_id = $2 ORDER BY "sequence"`, aggregateType, aggregateID)
	if err != nil {
		t.Error("query failed: ", err)
		return
	}

	if eventCount != len(resourceOwners) {
		t.Errorf("wrong queried event count: want %d, got %d", len(resourceOwners), eventCount)
	}
}

func assertEventCount(t *testing.T, db *database.DB, aggTypes database.TextArray[eventstore.AggregateType], aggIDs database.TextArray[string], maxPushedEventsCount int) {
	t.Helper()

	var count int
	err := db.QueryRow(func(row *sql.Row) error {
		return row.Scan(&count)
	}, "SELECT count(*) FROM eventstore.events2 where aggregate_type = ANY($1) AND aggregate_id = ANY($2)", aggTypes, aggIDs)

	if err != nil {
		t.Errorf("unexpected err in row.Scan: %v", err)
		return
	}

	if count > maxPushedEventsCount {
		t.Errorf("expected push count %d got %d", maxPushedEventsCount, count)
	}
}

func assertUniqueConstraint(t *testing.T, db *database.DB, commands []eventstore.Command, expectedCount int) {
	t.Helper()

	var uniqueConstraint *eventstore.UniqueConstraint
	for _, command := range commands {
		if e := command.(*testEvent); len(e.uniqueConstraints) > 0 {
			uniqueConstraint = e.uniqueConstraints[0]
			break
		}
	}
	if uniqueConstraint == nil {
		return
	}

	var uniqueCount int
	err := db.QueryRow(func(row *sql.Row) error {
		return row.Scan(&uniqueCount)
	}, "SELECT COUNT(*) FROM eventstore.unique_constraints where unique_type = $1 AND unique_field = $2", uniqueConstraint.UniqueType, uniqueConstraint.UniqueField)
	if err != nil {
		t.Error("unable to query inserted rows: ", err)
		return
	}
	if uniqueCount != expectedCount {
		t.Errorf("expected unique count %d got %d", expectedCount, uniqueCount)
	}
}

func TestEventstore_Push_UniqueConstraintOwners(t *testing.T) {
	for pusherName, pusher := range pushers {
		t.Run(pusherName, func(t *testing.T) {
			runUniqueConstraintOwnerCases(t, pusher, clients[pusherName])
		})
	}
}

func runUniqueConstraintOwnerCases(t *testing.T, pusher eventstore.Pusher, client *database.DB) {
	t.Helper()
	db := eventstore.NewEventstore(&eventstore.Config{
		Querier: queriers["v2(inmemory)"],
		Pusher:  pusher,
	})
	const instanceID = "owners-instance"
	for _, tc := range uniqueConstraintOwnerCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Cleanup(cleanupEventstore(client))
			tc.run(t, db, client, instanceID)
		})
	}
}

type uniqueConstraintOwnerCase struct {
	name string
	run  func(t *testing.T, db *eventstore.Eventstore, client *database.DB, instanceID string)
}

var uniqueConstraintOwnerCases = []uniqueConstraintOwnerCase{
	{
		name: "add stores owners",
		run: func(t *testing.T, db *eventstore.Eventstore, client *database.DB, instanceID string) {
			_, err := db.Push(context.Background(), generateCommand("owners-add", "1",
				withInstanceID(instanceID),
				generateAddUniqueConstraint("usernames", "alice", "org:org-1", "user:user-1"),
			))
			if err != nil {
				t.Fatal(err)
			}
			assertOwners(t, client, instanceID, "usernames", "alice", []string{"org:org-1", "user:user-1"})
		},
	},
	{
		name: "omit owners stores empty",
		run: func(t *testing.T, db *eventstore.Eventstore, client *database.DB, instanceID string) {
			_, err := db.Push(context.Background(), generateCommand("owners-empty", "1",
				withInstanceID(instanceID),
				generateAddUniqueConstraint("usernames", "bob"),
			))
			if err != nil {
				t.Fatal(err)
			}
			assertOwners(t, client, instanceID, "usernames", "bob", nil)
		},
	},
	{
		name: "add duplicate already exists",
		run: func(t *testing.T, db *eventstore.Eventstore, client *database.DB, instanceID string) {
			cmd := func() eventstore.Command {
				return generateCommand("owners-dup", "1",
					withInstanceID(instanceID),
					generateAddUniqueConstraint("usernames", "alice", "org:org-1"),
				)
			}
			if _, err := db.Push(context.Background(), cmd()); err != nil {
				t.Fatal(err)
			}
			_, err := db.Push(context.Background(), cmd())
			if !zerrors.IsErrorAlreadyExists(err) {
				t.Fatalf("expected already exists, got %v", err)
			}
		},
	},
	{
		name: "remove by field case insensitive",
		run: func(t *testing.T, db *eventstore.Eventstore, client *database.DB, instanceID string) {
			if err := insertUniqueConstraint(client, instanceID, "usernames", "mixedcase", "org:org-1"); err != nil {
				t.Fatal(err)
			}
			_, err := db.Push(context.Background(), generateCommand("owners-field-remove", "1",
				withInstanceID(instanceID),
				generateRemoveUniqueConstraint("usernames", "MixedCase"),
			))
			if err != nil {
				t.Fatal(err)
			}
			assertUniqueCount(t, client, instanceID, 0)
		},
	},
	{
		name: "remove by owner org",
		run: func(t *testing.T, db *eventstore.Eventstore, client *database.DB, instanceID string) {
			if err := insertUniqueConstraint(client, instanceID, "usernames", "alice", "org:org-x", "user:u1"); err != nil {
				t.Fatal(err)
			}
			if err := insertUniqueConstraint(client, instanceID, "usernames", "bob", "org:org-y", "user:u2"); err != nil {
				t.Fatal(err)
			}
			_, err := db.Push(context.Background(), generateCommand("owners-org-remove", "1",
				withInstanceID(instanceID),
				generateRemoveUniqueConstraintsByOwner(eventstore.UniqueConstraintOwnerOrg, "org-x"),
			))
			if err != nil {
				t.Fatal(err)
			}
			assertUniqueCount(t, client, instanceID, 1)
			assertOwners(t, client, instanceID, "usernames", "bob", []string{"org:org-y", "user:u2"})
		},
	},
	{
		name: "remove by owner user",
		run: func(t *testing.T, db *eventstore.Eventstore, client *database.DB, instanceID string) {
			if err := insertUniqueConstraint(client, instanceID, "usernames", "alice", "org:org-1", "user:u1"); err != nil {
				t.Fatal(err)
			}
			if err := insertUniqueConstraint(client, instanceID, "external_idps", "idp1ext1", "org:org-1", "idp:idp1", "user:u1"); err != nil {
				t.Fatal(err)
			}
			if err := insertUniqueConstraint(client, instanceID, "usernames", "bob", "org:org-1", "user:u2"); err != nil {
				t.Fatal(err)
			}
			_, err := db.Push(context.Background(), generateCommand("owners-user-remove", "1",
				withInstanceID(instanceID),
				generateRemoveUniqueConstraintsByOwner(eventstore.UniqueConstraintOwnerUser, "u1"),
			))
			if err != nil {
				t.Fatal(err)
			}
			assertUniqueCount(t, client, instanceID, 1)
			assertOwners(t, client, instanceID, "usernames", "bob", []string{"org:org-1", "user:u2"})
		},
	},
	{
		name: "remove by owner idp",
		run: func(t *testing.T, db *eventstore.Eventstore, client *database.DB, instanceID string) {
			if err := insertUniqueConstraint(client, instanceID, "external_idps", "idp1ext1", "org:org-1", "idp:idp-1", "user:u1"); err != nil {
				t.Fatal(err)
			}
			if err := insertUniqueConstraint(client, instanceID, "external_idps", "idp2ext2", "org:org-1", "idp:idp-2", "user:u2"); err != nil {
				t.Fatal(err)
			}
			_, err := db.Push(context.Background(), generateCommand("owners-idp-remove", "1",
				withInstanceID(instanceID),
				generateRemoveUniqueConstraintsByOwner(eventstore.UniqueConstraintOwnerIDP, "idp-1"),
			))
			if err != nil {
				t.Fatal(err)
			}
			assertUniqueCount(t, client, instanceID, 1)
			assertOwners(t, client, instanceID, "external_idps", "idp2ext2", []string{"org:org-1", "idp:idp-2", "user:u2"})
		},
	},
	{
		name: "two org tags either org deletes grant",
		run: func(t *testing.T, db *eventstore.Eventstore, client *database.DB, instanceID string) {
			if err := insertUniqueConstraint(client, instanceID, "project_grant", "granted:project", "org:granting", "org:granted", "project:project"); err != nil {
				t.Fatal(err)
			}
			_, err := db.Push(context.Background(), generateCommand("owners-grant-org", "1",
				withInstanceID(instanceID),
				generateRemoveUniqueConstraintsByOwner(eventstore.UniqueConstraintOwnerOrg, "granted"),
			))
			if err != nil {
				t.Fatal(err)
			}
			assertUniqueCount(t, client, instanceID, 0)
		},
	},
	{
		name: "instance remove deletes tagged rows",
		run: func(t *testing.T, db *eventstore.Eventstore, client *database.DB, instanceID string) {
			if err := insertUniqueConstraint(client, instanceID, "usernames", "alice", "org:org-1"); err != nil {
				t.Fatal(err)
			}
			if err := insertUniqueConstraint(client, "other", "usernames", "alice", "org:org-1"); err != nil {
				t.Fatal(err)
			}
			_, err := db.Push(context.Background(), generateCommand("owners-instance-remove", "1",
				withInstanceID(instanceID),
				generateRemoveInstanceUniqueConstraints(),
			))
			if err != nil {
				t.Fatal(err)
			}
			assertUniqueCount(t, client, instanceID, 0)
			assertUniqueCount(t, client, "other", 1)
		},
	},
	{
		name: "global and empty owners not deleted by owner",
		run: func(t *testing.T, db *eventstore.Eventstore, client *database.DB, instanceID string) {
			if err := insertUniqueConstraint(client, "", "instance_domain", "example.com"); err != nil {
				t.Fatal(err)
			}
			if err := insertUniqueConstraint(client, instanceID, "usernames", "untagged"); err != nil {
				t.Fatal(err)
			}
			if err := insertUniqueConstraint(client, instanceID, "usernames", "tagged", "org:org-1"); err != nil {
				t.Fatal(err)
			}
			_, err := db.Push(context.Background(), generateCommand("owners-skip-empty", "1",
				withInstanceID(instanceID),
				generateRemoveUniqueConstraintsByOwner(eventstore.UniqueConstraintOwnerOrg, "org-1"),
			))
			if err != nil {
				t.Fatal(err)
			}
			assertUniqueCount(t, client, "", 1)
			assertUniqueCount(t, client, instanceID, 1)
			assertOwners(t, client, instanceID, "usernames", "untagged", nil)
		},
	},
	{
		name: "mixed batch deletes then inserts",
		run: func(t *testing.T, db *eventstore.Eventstore, client *database.DB, instanceID string) {
			if err := insertUniqueConstraint(client, instanceID, "usernames", "old", "org:org-1"); err != nil {
				t.Fatal(err)
			}
			if err := insertUniqueConstraint(client, instanceID, "org_name", "acme", "org:org-1"); err != nil {
				t.Fatal(err)
			}
			_, err := db.Push(context.Background(), generateCommand("owners-mixed", "1",
				withInstanceID(instanceID),
				generateAddUniqueConstraint("usernames", "new", "org:org-1", "user:u1"),
				generateRemoveUniqueConstraint("usernames", "old"),
				generateRemoveUniqueConstraintsByOwner(eventstore.UniqueConstraintOwnerOrg, "org-1"),
			))
			if err != nil {
				t.Fatal(err)
			}
			assertUniqueCount(t, client, instanceID, 1)
			assertOwners(t, client, instanceID, "usernames", "new", []string{"org:org-1", "user:u1"})
		},
	},
}

func insertUniqueConstraint(db *database.DB, instanceID, uniqueType, uniqueField string, owners ...string) error {
	if len(owners) == 0 {
		_, err := db.Exec("INSERT INTO eventstore.unique_constraints (instance_id, unique_type, unique_field) VALUES ($1, $2, $3)", instanceID, uniqueType, uniqueField)
		return err
	}
	_, err := db.Exec("INSERT INTO eventstore.unique_constraints (instance_id, unique_type, unique_field, owners) VALUES ($1, $2, $3, $4)", instanceID, uniqueType, uniqueField, database.TextArray[string](owners))
	return err
}

func assertUniqueCount(t *testing.T, db *database.DB, instanceID string, want int) {
	t.Helper()
	var count int
	err := db.QueryRow(func(row *sql.Row) error {
		return row.Scan(&count)
	}, "SELECT COUNT(*) FROM eventstore.unique_constraints WHERE instance_id = $1", instanceID)
	if err != nil {
		t.Fatal(err)
	}
	if count != want {
		t.Errorf("unique count = %d, want %d", count, want)
	}
}

func assertOwners(t *testing.T, db *database.DB, instanceID, uniqueType, uniqueField string, want []string) {
	t.Helper()
	var owners database.TextArray[string]
	err := db.QueryRow(func(row *sql.Row) error {
		return row.Scan(&owners)
	}, "SELECT owners FROM eventstore.unique_constraints WHERE instance_id = $1 AND unique_type = $2 AND unique_field = $3", instanceID, uniqueType, uniqueField)
	if err != nil {
		t.Fatal(err)
	}
	if len(owners) != len(want) {
		t.Fatalf("owners = %v, want %v", []string(owners), want)
	}
	got := map[string]int{}
	for _, owner := range owners {
		got[owner]++
	}
	for _, owner := range want {
		if got[owner] == 0 {
			t.Fatalf("owners = %v, want %v", []string(owners), want)
		}
		got[owner]--
	}
}
