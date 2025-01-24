package execution

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/query"
)

func Test_reduceEventExecution(t *testing.T) {
	type fields struct {
		config WorkerConfig
	}
	type args struct {
		ctx     context.Context
		event   *EventExecution
		targets []*query.ExecutionTarget
	}
	type res struct {
		wantErr bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"max TTL",
			fields{
				config: WorkerConfig{MaxTtl: time.Minute},
			},
			args{
				ctx:     context.Background(),
				event:   newMockEventExecution(time.Now().Add(-time.Minute * 2)),
				targets: mockTargets(0),
			},
			res{
				wantErr: false,
			},
		},
		{
			"none",
			fields{
				config: WorkerConfig{MaxTtl: time.Minute},
			},
			args{
				ctx:     context.Background(),
				event:   newMockEventExecution(time.Now()),
				targets: mockTargets(0),
			},
			res{
				wantErr: false,
			},
		},
		{
			"single",
			fields{
				config: WorkerConfig{MaxTtl: time.Minute},
			},
			args{
				ctx:     context.Background(),
				event:   newMockEventExecution(time.Now()),
				targets: mockTargets(1),
			},
			res{
				wantErr: false,
			},
		},
		{
			"multiple",
			fields{
				config: WorkerConfig{MaxTtl: time.Minute},
			},
			args{
				ctx:     context.Background(),
				event:   newMockEventExecution(time.Now()),
				targets: mockTargets(3),
			},
			res{
				wantErr: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			worker := &Worker{
				config: tt.fields.config,
				now:    time.Now,
			}

			closeFuncs := make([]func(), len(tt.args.targets))
			calledFuncs := make([]func() bool, len(tt.args.targets))
			for i := range tt.args.targets {
				url, closeF, calledF := TestServerCall(
					nil,
					time.Second,
					http.StatusOK,
					nil,
				)
				tt.args.targets[i].Endpoint = url
				closeFuncs[i] = closeF
				calledFuncs[i] = calledF
			}

			data, err := json.Marshal(tt.args.targets)
			require.NoError(t, err)
			tt.args.event.TargetsData = data

			err = worker.reduceEventExecution(
				tt.args.ctx,
				tt.args.event,
			)

			if tt.res.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			for _, closeF := range closeFuncs {
				closeF()
			}
			for _, calledF := range calledFuncs {
				assert.True(t, calledF())
			}
		})
	}
}

func newMockEventExecution(createdAt time.Time) *EventExecution {
	return &EventExecution{
		Aggregate: &eventstore.Aggregate{
			ID:            "aggID",
			Type:          "aggType",
			ResourceOwner: "resourceOwner",
			InstanceID:    "instanceID",
			Version:       "version",
		},
		Sequence:  15,
		EventType: "session.added",
		CreatedAt: createdAt,
		UserID:    "userID",
		EventData: []byte("eventData"),
	}
}
