package eventstore_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"testing"

	"github.com/zitadel/logging"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore/v4"
)

type benchmarkTest struct {
	name    string
	cmds    []eventstore.Cmd
	workers int
}

var (
	instanceID1 = "1"
	instanceID2 = "2"

	ctx = authz.SetCtxData(context.Background(), authz.CtxData{UserID: "adlerhurst", OrgID: "myorg"})

	ctxInstance1 = authz.WithInstanceID(ctx, instanceID1)
	ctxInstance2 = authz.WithInstanceID(ctx, instanceID2)

	agg1 = eventstore.NewAggregate(ctx, "ng5PD", "test")
	agg2 = eventstore.NewAggregate(ctx, "e4epE", "test")

	testCases = []benchmarkTest{
		{
			name: "without",
			cmds: []eventstore.Cmd{
				commandWithoutPayload(ctx, agg1),
				commandWithoutPayload(ctx, agg2),
			},
			workers: 1,
		},
		{
			name: "with",
			cmds: []eventstore.Cmd{
				cmdWithPayload(ctx, agg1),
				cmdWithPayload(ctx, agg2),
			},
			workers: 1,
		},
		{
			name: "without",
			cmds: []eventstore.Cmd{
				commandWithoutPayload(ctx, agg1),
				commandWithoutPayload(ctx, agg2),
			},
			workers: 3,
		},
		{
			name: "with",
			cmds: []eventstore.Cmd{
				cmdWithPayload(ctx, agg1),
				cmdWithPayload(ctx, agg2),
			},
			workers: 3,
		},
		{
			name: "without",
			cmds: []eventstore.Cmd{
				commandWithoutPayload(ctx, agg1),
				commandWithoutPayload(ctx, agg2),
			},
			workers: 5,
		},
		{
			name: "with",
			cmds: []eventstore.Cmd{
				cmdWithPayload(ctx, agg1),
				cmdWithPayload(ctx, agg2),
			},
			workers: 5,
		},
		{
			name: "without",
			cmds: []eventstore.Cmd{
				commandWithoutPayload(ctx, agg1),
				commandWithoutPayload(ctx, agg2),
			},
			workers: 7,
		},
		{
			name: "with",
			cmds: []eventstore.Cmd{
				cmdWithPayload(ctx, agg1),
				cmdWithPayload(ctx, agg2),
			},
			workers: 7,
		},
		{
			name: "without",
			cmds: []eventstore.Cmd{
				commandWithoutPayload(ctx, agg1),
				commandWithoutPayload(ctx, agg2),
			},
			workers: 9,
		},
		{
			name: "with",
			cmds: []eventstore.Cmd{
				cmdWithPayload(ctx, agg1),
				cmdWithPayload(ctx, agg2),
			},
			workers: 9,
		},
	}
)

func Benchmark_Eventstore_Push(b *testing.B) {
	for _, tt := range testCases {
		execTest(b, localClient, tt.workers, fmt.Sprintf("%d-workers_%s-payload", tt.workers, tt.name), tt.cmds, false)
	}
}

func Benchmark_Eventstore_Push_2_instances(b *testing.B) {
	for _, tt := range testCases {
		execTest(b, localClient, tt.workers, fmt.Sprintf("%d-workers_%s-payload", tt.workers, tt.name), tt.cmds, true)
	}
}

func execTest(b *testing.B, client *sql.DB, workers int, name string, commands []eventstore.Cmd, twoInstances bool) {
	b.Helper()

	// warmup sql connections
	var warmupWg sync.WaitGroup
	warmupWg.Add(40)
	for i := 0; i < 40; i++ {
		go func() {
			client.Ping()
			warmupWg.Done()
		}()
	}
	warmupWg.Wait()

	es, err := eventstore.Start(&eventstore.Config{Client: client})
	if err != nil {
		b.Fatal("unable to init eventstore: ", err)
	}
	if _, err = localClient.Exec("TRUNCATE zitadel.eventstore_v4.events;"); err != nil {
		b.Fatal("unable to truncate table: ", err)
	}
	if _, err = localClient.Exec("TRUNCATE zitadel.eventstore_v4.sequences;"); err != nil {
		b.Fatal("unable to truncate table: ", err)
	}

	b.Run(name, func(b *testing.B) {
		var wg sync.WaitGroup
		errCount := float64(0)
		for i := 0; i < workers; i++ {
			wg.Add(1)
			go func(i int) {
				pushCtx := authz.WithInstanceID(ctx, strconv.Itoa(i))
				for n := 0; n < b.N; n++ {
					if _, err := es.Push(pushCtx, commands...); err != nil {
						errCount++
					}
				}
				wg.Done()
			}(i)
		}
		wg.Wait()
		b.ReportMetric(errCount, "errors")
	})
}

type benchCommand struct {
	eventstore.Command `json:"-"`
	payload            []byte
}

func commandWithoutPayload(ctx context.Context, agg *eventstore.Aggregate) *benchCommand {
	return &benchCommand{
		Command: *eventstore.NewCommand("test", 1, &eventstore.Editor{UserID: "user-1", Service: "svc-1"}, agg),
	}
}

func cmdWithPayload(ctx context.Context, agg *eventstore.Aggregate) *benchCommand {
	cmd := commandWithoutPayload(ctx, agg)
	var err error

	cmd.payload, err = json.Marshal(struct {
		Username    string
		Firstname   string
		Lastname    string
		Email       string
		DisplayName string
		Gender      int8
	}{
		// The IT crowd S2.E4
		Username:    "peterfile",
		Firstname:   "Peter",
		Lastname:    "File",
		Email:       "peter.file@somemail.com",
		DisplayName: "Peter File",
		Gender:      10,
	})
	if err != nil {
		logging.Fatal("unable to create payload: ", err)
	}

	return cmd
}

func (cmd *benchCommand) Payload() interface{} {
	if len(cmd.payload) == 0 {
		return nil
	}
	return cmd.payload
}

// func (cmd *benchCommand) UniqueConstraints() []*eventstore.EventUniqueConstraint {
// 	return nil
// }
