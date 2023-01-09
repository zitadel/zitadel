package eventstore_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"testing"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore"
)

func Benchmark_Eventstore_PushOneAggregate(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	ctx = authz.SetCtxData(ctx, authz.CtxData{UserID: "adlerhurst", OrgID: "myorg"})
	defer cancel()

	agg := eventstore.NewAggregate(ctx, "S7boD", "test", "v1")

	tests := []struct {
		name   string
		client *sql.DB
		cmds   []eventstore.Command
	}{
		{
			name:   "1 event - no payload - sequential",
			client: localClient,
			cmds: []eventstore.Command{
				commandWithoutPayload(ctx, agg),
			},
		},
		{
			name:   "1 event - payload - sequential",
			client: localClient,
			cmds: []eventstore.Command{
				cmdWithPayload(ctx, b, agg),
				cmdWithPayload(ctx, b, agg),
			},
		},
		{
			name:   "5 event - no payload - sequential",
			client: localClient,
			cmds: []eventstore.Command{
				commandWithoutPayload(ctx, agg),
				commandWithoutPayload(ctx, agg),
				commandWithoutPayload(ctx, agg),
				commandWithoutPayload(ctx, agg),
				commandWithoutPayload(ctx, agg),
			},
		},
		{
			name:   "5 event - payload - sequential",
			client: localClient,
			cmds: []eventstore.Command{
				cmdWithPayload(ctx, b, agg),
				cmdWithPayload(ctx, b, agg),
				cmdWithPayload(ctx, b, agg),
				cmdWithPayload(ctx, b, agg),
				cmdWithPayload(ctx, b, agg),
			},
		},

		{
			name:   "2 events - no payload - parallel",
			client: localClient,
			cmds: []eventstore.Command{
				commandWithoutPayload(ctx, agg),
			},
		},
		{
			name:   "2 events - payload - parallel",
			client: localClient,
			cmds: []eventstore.Command{
				cmdWithPayload(ctx, b, agg),
				cmdWithPayload(ctx, b, agg),
			},
		},
		{
			name:   "5 event - no payload - parallel",
			client: localClient,
			cmds: []eventstore.Command{
				commandWithoutPayload(ctx, agg),
				commandWithoutPayload(ctx, agg),
				commandWithoutPayload(ctx, agg),
				commandWithoutPayload(ctx, agg),
				commandWithoutPayload(ctx, agg),
			},
		},
		{
			name:   "5 event - payload - parallel",
			client: localClient,
			cmds: []eventstore.Command{
				cmdWithPayload(ctx, b, agg),
				cmdWithPayload(ctx, b, agg),
				cmdWithPayload(ctx, b, agg),
				cmdWithPayload(ctx, b, agg),
				cmdWithPayload(ctx, b, agg),
			},
		},
	}
	for _, tt := range tests {
		execTest(b, tt.client, tt.name, tt.cmds)
	}
}

func Benchmark_Eventstore_PushMultipleAggregatesSequential(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	ctx = authz.SetCtxData(ctx, authz.CtxData{UserID: "adlerhurst", OrgID: "myorg"})
	defer cancel()

	agg1 := eventstore.NewAggregate(ctx, "ng5PD", "test", "v1")
	agg2 := eventstore.NewAggregate(ctx, "e4epE", "test", "v1")
	agg3 := eventstore.NewAggregate(ctx, "vE0uJ", "test", "v1")

	tests := []struct {
		name    string
		client  *sql.DB
		cmds    []eventstore.Command
		workers int
	}{
		{
			name:   "1 event - no payload - sequential",
			client: localClient,
			cmds: []eventstore.Command{
				commandWithoutPayload(ctx, agg1),
				commandWithoutPayload(ctx, agg2),
			},
			workers: 1,
		},
		{
			name:   "1 event - payload - sequential",
			client: localClient,
			cmds: []eventstore.Command{
				cmdWithPayload(ctx, b, agg1),
				cmdWithPayload(ctx, b, agg2),
			},
			workers: 1,
		},
		{
			name:   "5 event - no payload - sequential",
			client: localClient,
			cmds: []eventstore.Command{
				commandWithoutPayload(ctx, agg1),
				commandWithoutPayload(ctx, agg2),
				commandWithoutPayload(ctx, agg3),
				commandWithoutPayload(ctx, agg1),
				commandWithoutPayload(ctx, agg2),
				commandWithoutPayload(ctx, agg3),
			},
			workers: 1,
		},
		{
			name:   "5 event - payload - sequential",
			client: localClient,
			cmds: []eventstore.Command{
				commandWithoutPayload(ctx, agg1),
				commandWithoutPayload(ctx, agg2),
				commandWithoutPayload(ctx, agg3),
				commandWithoutPayload(ctx, agg1),
				commandWithoutPayload(ctx, agg2),
				commandWithoutPayload(ctx, agg3),
			},
			workers: 1,
		},
		{
			name:   "1 event - payload - parallel 2",
			client: localClient,
			cmds: []eventstore.Command{
				cmdWithPayload(ctx, b, agg1),
				cmdWithPayload(ctx, b, agg2),
			},
			workers: 2,
		},
		{
			name:   "5 event - no payload - parallel 5",
			client: localClient,
			cmds: []eventstore.Command{
				commandWithoutPayload(ctx, agg1),
				commandWithoutPayload(ctx, agg2),
				commandWithoutPayload(ctx, agg3),
				commandWithoutPayload(ctx, agg1),
				commandWithoutPayload(ctx, agg2),
				commandWithoutPayload(ctx, agg3),
			},
			workers: 5,
		},
		{
			name:   "5 event - payload - parallel 5",
			client: localClient,
			cmds: []eventstore.Command{
				commandWithoutPayload(ctx, agg1),
				commandWithoutPayload(ctx, agg2),
				commandWithoutPayload(ctx, agg3),
				commandWithoutPayload(ctx, agg1),
				commandWithoutPayload(ctx, agg2),
				commandWithoutPayload(ctx, agg3),
			},
			workers: 5,
		},

		{
			name:   "1 event - no payload - parallel 10",
			client: localClient,
			cmds: []eventstore.Command{
				commandWithoutPayload(ctx, agg1),
				commandWithoutPayload(ctx, agg2),
			},
			workers: 10,
		},
		{
			name:   "1 event - payload - parallel 10",
			client: localClient,
			cmds: []eventstore.Command{
				cmdWithPayload(ctx, b, agg1),
				cmdWithPayload(ctx, b, agg2),
			},
			workers: 10,
		},
		{
			name:   "6 event - no payload - parallel 10",
			client: localClient,
			cmds: []eventstore.Command{
				commandWithoutPayload(ctx, agg1),
				commandWithoutPayload(ctx, agg2),
				commandWithoutPayload(ctx, agg3),
				commandWithoutPayload(ctx, agg1),
				commandWithoutPayload(ctx, agg2),
				commandWithoutPayload(ctx, agg3),
			},
			workers: 10,
		},
		{
			name:   "6 event - payload - parallel 10",
			client: localClient,
			cmds: []eventstore.Command{
				commandWithoutPayload(ctx, agg1),
				commandWithoutPayload(ctx, agg2),
				commandWithoutPayload(ctx, agg3),
				commandWithoutPayload(ctx, agg1),
				commandWithoutPayload(ctx, agg2),
				commandWithoutPayload(ctx, agg3),
			},
			workers: 10,
		},
	}
	for _, tt := range tests {
		execTest(b, tt.client, tt.name, tt.cmds)
	}
}

func Benchmark_Eventstore_Push2AggregatesParallelWithoutPayload(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	ctx = authz.SetCtxData(ctx, authz.CtxData{UserID: "adlerhurst", OrgID: "myorg"})
	defer cancel()

	agg1 := eventstore.NewAggregate(ctx, "ng5PD", "test", "v1")
	agg2 := eventstore.NewAggregate(ctx, "e4epE", "test", "v1")

	t := struct {
		name   string
		client *sql.DB
		cmds   []eventstore.Command
	}{
		name:   "1 event - no payload - parallel 2",
		client: localClient,
		cmds: []eventstore.Command{
			commandWithoutPayload(ctx, agg1),
			commandWithoutPayload(ctx, agg2),
		},
	}

	execTestParallel(b, t.client, t.name, t.cmds)
}

func Benchmark_Eventstore_Push3AggregatesParallelWithoutPayload(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	ctx = authz.SetCtxData(ctx, authz.CtxData{UserID: "adlerhurst", OrgID: "myorg"})
	defer cancel()

	agg1 := eventstore.NewAggregate(ctx, "ng5PD", "test", "v1")
	agg2 := eventstore.NewAggregate(ctx, "e4epE", "test", "v1")
	agg3 := eventstore.NewAggregate(ctx, "jVGaS", "test", "v1")

	t := struct {
		name   string
		client *sql.DB
		cmds   []eventstore.Command
	}{
		name:   "5 event - no payload - parallel 5",
		client: localClient,
		cmds: []eventstore.Command{
			commandWithoutPayload(ctx, agg1),
			commandWithoutPayload(ctx, agg2),
			commandWithoutPayload(ctx, agg3),
			commandWithoutPayload(ctx, agg1),
			commandWithoutPayload(ctx, agg2),
			commandWithoutPayload(ctx, agg3),
		},
	}

	execTestParallel(b, t.client, t.name, t.cmds)
}

func Benchmark_Eventstore_Push2AggregatesParallelWithPayload(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	ctx = authz.SetCtxData(ctx, authz.CtxData{UserID: "adlerhurst", OrgID: "myorg"})
	defer cancel()

	agg1 := eventstore.NewAggregate(ctx, "ng5PD", "test", "v1")
	agg2 := eventstore.NewAggregate(ctx, "e4epE", "test", "v1")

	t := struct {
		name   string
		client *sql.DB
		cmds   []eventstore.Command
	}{
		name:   "1 event - no payload - parallel 2",
		client: localClient,
		cmds: []eventstore.Command{
			cmdWithPayload(ctx, b, agg1),
			cmdWithPayload(ctx, b, agg2),
		},
	}

	execTestParallel(b, t.client, t.name, t.cmds)
}

func Benchmark_Eventstore_Push3AggregatesParallelWithPayload(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	ctx = authz.SetCtxData(ctx, authz.CtxData{UserID: "adlerhurst", OrgID: "myorg"})
	defer cancel()

	agg1 := eventstore.NewAggregate(ctx, "ng5PD", "test", "v1")
	agg2 := eventstore.NewAggregate(ctx, "e4epE", "test", "v1")
	agg3 := eventstore.NewAggregate(ctx, "jVGaS", "test", "v1")

	t := struct {
		name   string
		client *sql.DB
		cmds   []eventstore.Command
	}{
		name:   "5 event - no payload - parallel 5",
		client: localClient,
		cmds: []eventstore.Command{
			cmdWithPayload(ctx, b, agg1),
			cmdWithPayload(ctx, b, agg2),
			cmdWithPayload(ctx, b, agg3),
			cmdWithPayload(ctx, b, agg1),
			cmdWithPayload(ctx, b, agg2),
			cmdWithPayload(ctx, b, agg3),
		},
	}

	execTestParallel(b, t.client, t.name, t.cmds)
}

func execTest(b *testing.B, client *sql.DB, name string, commands []eventstore.Command) {
	b.Helper()

	es, err := eventstore.Start(&eventstore.Config{Client: client})
	if err != nil {
		b.Fatal("unable to init eventstore: ", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	ctx = authz.SetCtxData(ctx, authz.CtxData{UserID: "adlerhurst", OrgID: "myorg"})
	defer cancel()

	b.Run(name, func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			if _, err := es.Push(ctx, commands...); err != nil {
				b.Error("push failed: ", err)
			}
		}
	})

	if _, err = localClient.Exec("TRUNCATE eventstore.events;"); err != nil {
		b.Fatal("unable to truncate table: ", err)
	}
}

func execTestParallel(b *testing.B, client *sql.DB, name string, commands []eventstore.Command) {
	b.Helper()

	es, err := eventstore.Start(&eventstore.Config{Client: client})
	if err != nil {
		b.Fatal("unable to init eventstore: ", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	ctx = authz.SetCtxData(ctx, authz.CtxData{UserID: "adlerhurst", OrgID: "myorg"})
	defer cancel()

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			_, err := es.Push(ctx, commands...)
			if err != nil {
				b.Error(err)
			}
		}
	})

	if _, err = localClient.Exec("TRUNCATE eventstore.events;"); err != nil {
		b.Fatal("unable to truncate table: ", err)
	}
}

type benchCommand struct {
	eventstore.BaseEvent
	payload []byte
}

func commandWithoutPayload(ctx context.Context, agg *eventstore.Aggregate) *benchCommand {
	typ := eventstore.EventType("test")
	return &benchCommand{
		BaseEvent: *eventstore.NewBaseEventForPush(ctx, agg, typ),
	}
}

func cmdWithPayload(ctx context.Context, b *testing.B, agg *eventstore.Aggregate) *benchCommand {
	b.Helper()

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
		Username:    "peterfile",
		Firstname:   "Peter",
		Lastname:    "File",
		Email:       "peter.file@somemail.com",
		DisplayName: "Peter File",
		Gender:      10,
	})
	if err != nil {
		b.Fatal("unable to create payload: ", err)
	}

	return cmd
}

func (cmd *benchCommand) Data() interface{} {
	if len(cmd.payload) == 0 {
		return nil
	}
	return cmd.payload
}

func (cmd *benchCommand) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}
