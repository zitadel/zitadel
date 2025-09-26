package domain

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/dbmock"
	"github.com/zitadel/zitadel/backend/v3/storage/eventstore"
	legacy_es "github.com/zitadel/zitadel/internal/eventstore"
)

type testLegacyEventstore struct{}

func (es *testLegacyEventstore) PushWithNewClient(ctx context.Context, client database.QueryExecutor, commands ...legacy_es.Command) ([]legacy_es.Event, error) {
	return nil, nil
}

var _ eventstore.LegacyEventstore = (*testLegacyEventstore)(nil)

type invokeTestCommand struct {
	events  []legacy_es.Command
	execute func(ctx context.Context, opts *CommandOpts) error
}

// Events implements Commander.
func (i *invokeTestCommand) Events(ctx context.Context, opts *CommandOpts) []legacy_es.Command {
	return i.events
}

// Execute implements Commander.
func (i *invokeTestCommand) Execute(ctx context.Context, opts *CommandOpts) (err error) {
	if i.execute == nil {
		return nil
	}
	return i.execute(ctx, opts)
}

// String implements Commander.
func (i *invokeTestCommand) String() string {
	return "invokeTestCommand"
}

// Validate implements Commander.
func (i *invokeTestCommand) Validate(ctx context.Context, opts *CommandOpts) (err error) {
	return nil
}

var _ Commander = (*invokeTestCommand)(nil)

type invokeTestEvent struct {
	id string
}

// Aggregate implements [legacy_es.Command].
func (i *invokeTestEvent) Aggregate() *legacy_es.Aggregate {
	return &legacy_es.Aggregate{
		ID:            i.id,
		Type:          "test",
		ResourceOwner: "test",
		InstanceID:    "test",
		Version:       legacy_es.Version("v1"),
	}
}

// Creator implements [legacy_es.Command].
func (i *invokeTestEvent) Creator() string {
	return "test"
}

// Fields implements [legacy_es.Command].
func (i *invokeTestEvent) Fields() []*legacy_es.FieldOperation {
	return nil
}

// Payload implements [legacy_es.Command].
func (i *invokeTestEvent) Payload() any {
	return nil
}

// Revision implements [legacy_es.Command].
func (i *invokeTestEvent) Revision() uint16 {
	return 1
}

// Type implements [legacy_es.Command].
func (i *invokeTestEvent) Type() legacy_es.EventType {
	return "test"
}

// UniqueConstraints implements [legacy_es.Command].
func (i *invokeTestEvent) UniqueConstraints() []*legacy_es.UniqueConstraint {
	return nil
}

var _ legacy_es.Command = (*invokeTestEvent)(nil)

func Test_eventCollector_Invoke(t *testing.T) {
	tests := []struct {
		name                  string
		command               Commander
		expectedErr           error
		assertCollectedEvents func(t *testing.T, events []legacy_es.Command)
	}{
		{
			name:        "simple command with events",
			expectedErr: nil,
			command: &invokeTestCommand{
				events: []legacy_es.Command{&invokeTestEvent{id: "1"}},
			},
			assertCollectedEvents: func(t *testing.T, events []legacy_es.Command) {
				require.Len(t, events, 1)
				assert.Equal(t, "1", events[0].Aggregate().ID)
			},
		},
		{
			name:        "simple command without events",
			expectedErr: nil,
			command:     &invokeTestCommand{},
			assertCollectedEvents: func(t *testing.T, events []legacy_es.Command) {
				assert.Len(t, events, 0)
			},
		},
		{
			name:        "command with sub commands",
			expectedErr: nil,
			command: &invokeTestCommand{
				events: []legacy_es.Command{&invokeTestEvent{id: "1"}},
				execute: func(ctx context.Context, opts *CommandOpts) error {
					return opts.Invoke(ctx, &invokeTestCommand{
						events: []legacy_es.Command{&invokeTestEvent{id: "2"}},
					})
				},
			},
			assertCollectedEvents: func(t *testing.T, events []legacy_es.Command) {
				require.Len(t, events, 2)
				assert.Equal(t, "1", events[0].Aggregate().ID)
				assert.Equal(t, "2", events[1].Aggregate().ID)
			},
		},
		{
			name:        "only sub commands with events",
			expectedErr: nil,
			command: &invokeTestCommand{
				execute: func(ctx context.Context, opts *CommandOpts) error {
					return opts.Invoke(ctx, &invokeTestCommand{
						events: []legacy_es.Command{&invokeTestEvent{id: "2"}},
					})
				},
			},
			assertCollectedEvents: func(t *testing.T, events []legacy_es.Command) {
				require.Len(t, events, 1)
				assert.Equal(t, "2", events[0].Aggregate().ID)
			},
		},
	}

	SetLegacyEventstore(new(testLegacyEventstore))

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			pool := dbmock.NewMockPool(ctrl)

			tx := dbmock.NewMockTransaction(ctrl)
			tx.EXPECT().End(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			pool.EXPECT().Begin(gomock.Any(), gomock.Any()).Return(tx, nil).AnyTimes()

			i := newEventStoreInvoker(nil)
			opts := CommandOpts{
				Invoker: i,
				DB:      pool,
			}
			err := opts.Invoke(t.Context(), tt.command)
			require.ErrorIs(t, err, tt.expectedErr)
			tt.assertCollectedEvents(t, i.collector.events)
		})
	}
}
