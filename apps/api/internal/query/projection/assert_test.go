package projection

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

func Test_assertEvent(t *testing.T) {
	type args struct {
		event      eventstore.Event
		assertFunc func(eventstore.Event) (eventstore.Event, error)
	}
	type testCase struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}
	tests := []testCase{
		{
			name: "correct event type",
			args: args{
				event: instance.NewInstanceAddedEvent(context.Background(), &instance.NewAggregate("instance-id").Aggregate, "instance-name"),
				assertFunc: func(event eventstore.Event) (eventstore.Event, error) {
					return assertEvent[*instance.InstanceAddedEvent](event)
				},
			},
			wantErr: assert.NoError,
		}, {
			name: "wrong event type",
			args: args{
				event: instance.NewInstanceRemovedEvent(context.Background(), &instance.NewAggregate("instance-id").Aggregate, "instance-name", nil),
				assertFunc: func(event eventstore.Event) (eventstore.Event, error) {
					return assertEvent[*instance.InstanceAddedEvent](event)
				},
			},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.args.assertFunc(tt.args.event)
			if !tt.wantErr(t, err) {
				return
			}
		})
	}
}
