package domain_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/zitadel/zitadel/backend/v3/domain"
	domainmock "github.com/zitadel/zitadel/backend/v3/domain/mock"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/dbmock"
	"github.com/zitadel/zitadel/backend/v3/storage/eventstore"
	legacy_es "github.com/zitadel/zitadel/internal/eventstore"
)

type testLegacyEventstore struct {
	commands []legacy_es.Command
}

func (es *testLegacyEventstore) PushWithNewClient(ctx context.Context, client database.QueryExecutor, commands ...legacy_es.Command) ([]legacy_es.Event, error) {
	es.commands = append(es.commands, commands...)
	return nil, nil
}

var _ eventstore.LegacyEventstore = (*testLegacyEventstore)(nil)

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
		command               func(ctrl *gomock.Controller) domain.Executor
		queryExecutor         func(ctrl *gomock.Controller) database.QueryExecutor
		expectedErr           error
		assertCollectedEvents func(t *testing.T, events []legacy_es.Command)
	}{
		{
			name:        "simple command without events",
			expectedErr: nil,
			command: func(ctrl *gomock.Controller) domain.Executor {
				cmd := domainmock.NewMockCommander(ctrl)
				gomock.InOrder(
					cmd.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil),
					cmd.EXPECT().Events(gomock.Any(), gomock.Any()).Return(nil, nil),
				)
				return cmd
			},
			assertCollectedEvents: func(t *testing.T, events []legacy_es.Command) {
				assert.Len(t, events, 0)
			},
		},
		{
			name:        "simple command with events",
			expectedErr: nil,
			command: func(ctrl *gomock.Controller) domain.Executor {
				cmd := domainmock.NewMockCommander(ctrl)
				gomock.InOrder(
					cmd.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil),
					cmd.EXPECT().Events(gomock.Any(), gomock.Any()).Return([]legacy_es.Command{&invokeTestEvent{id: "1"}}, nil),
				)
				return cmd
			},
			assertCollectedEvents: func(t *testing.T, events []legacy_es.Command) {
				require.Len(t, events, 1)
				assert.Equal(t, "1", events[0].Aggregate().ID)
			},
		},
		{
			name:        "command with sub commands",
			expectedErr: nil,
			command: func(ctrl *gomock.Controller) domain.Executor {
				mainCmd := domainmock.NewMockCommander(ctrl)
				subCmd := domainmock.NewMockCommander(ctrl)
				gomock.InOrder(
					mainCmd.EXPECT().Execute(gomock.Any(), gomock.Any()).DoAndReturn(
						func(ctx context.Context, opts *domain.InvokeOpts) error {
							return opts.Invoke(ctx, subCmd)
						},
					),
					subCmd.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil),
					subCmd.EXPECT().Events(gomock.Any(), gomock.Any()).Return([]legacy_es.Command{&invokeTestEvent{id: "2"}}, nil),
					mainCmd.EXPECT().Events(gomock.Any(), gomock.Any()).Return([]legacy_es.Command{&invokeTestEvent{id: "1"}}, nil),
				)
				return mainCmd
			},
			assertCollectedEvents: func(t *testing.T, events []legacy_es.Command) {
				require.Len(t, events, 2)
				assert.Equal(t, "1", events[0].Aggregate().ID)
				assert.Equal(t, "2", events[1].Aggregate().ID)
			},
		},
		{
			name:        "command with multiple sub commands",
			expectedErr: nil,
			command: func(ctrl *gomock.Controller) domain.Executor {
				mainCmd := domainmock.NewMockCommander(ctrl)
				firstSubCmd := domainmock.NewMockCommander(ctrl)
				secondSubCommand := domainmock.NewMockCommander(ctrl)
				gomock.InOrder(
					mainCmd.EXPECT().Execute(gomock.Any(), gomock.Any()).DoAndReturn(
						func(ctx context.Context, opts *domain.InvokeOpts) error {
							return opts.Invoke(ctx, firstSubCmd)
						},
					),
					firstSubCmd.EXPECT().Execute(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, io *domain.InvokeOpts) error {
						return io.Invoke(ctx, secondSubCommand)
					}),
					secondSubCommand.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil),
					secondSubCommand.EXPECT().Events(gomock.Any(), gomock.Any()).Return([]legacy_es.Command{&invokeTestEvent{id: "3"}}, nil),
					firstSubCmd.EXPECT().Events(gomock.Any(), gomock.Any()).Return([]legacy_es.Command{&invokeTestEvent{id: "2"}}, nil),
					mainCmd.EXPECT().Events(gomock.Any(), gomock.Any()).Return([]legacy_es.Command{&invokeTestEvent{id: "1"}}, nil),
				)
				return mainCmd
			},
			assertCollectedEvents: func(t *testing.T, events []legacy_es.Command) {
				require.Len(t, events, 3)
				assert.Equal(t, "1", events[0].Aggregate().ID)
				assert.Equal(t, "2", events[1].Aggregate().ID)
				assert.Equal(t, "3", events[2].Aggregate().ID)
			},
		},
		{
			name:        "only sub commands with events",
			expectedErr: nil,
			command: func(ctrl *gomock.Controller) domain.Executor {
				mainCmd := domainmock.NewMockCommander(ctrl)
				subCmd := domainmock.NewMockCommander(ctrl)
				gomock.InOrder(
					mainCmd.EXPECT().Execute(gomock.Any(), gomock.Any()).DoAndReturn(
						func(ctx context.Context, opts *domain.InvokeOpts) error {
							return opts.Invoke(ctx, subCmd)
						},
					),
					subCmd.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil),
					subCmd.EXPECT().Events(gomock.Any(), gomock.Any()).Return([]legacy_es.Command{&invokeTestEvent{id: "2"}}, nil),
					mainCmd.EXPECT().Events(gomock.Any(), gomock.Any()).Return(nil, nil),
				)
				return mainCmd
			},
			assertCollectedEvents: func(t *testing.T, events []legacy_es.Command) {
				require.Len(t, events, 1)
				assert.Equal(t, "2", events[0].Aggregate().ID)
			},
		},
		{
			name:        "executor batch all commands",
			expectedErr: nil,
			command: func(ctrl *gomock.Controller) domain.Executor {
				mainCommand := domainmock.NewMockCommander(ctrl)
				firstSubCommand := domainmock.NewMockCommander(ctrl)
				secondSubCommand := domainmock.NewMockCommander(ctrl)

				gomock.InOrder(
					mainCommand.EXPECT().Execute(gomock.Any(), gomock.Any()).DoAndReturn(
						func(ctx context.Context, opts *domain.InvokeOpts) error {
							return opts.Invoke(ctx,
								domain.BatchExecutors(
									firstSubCommand,
									secondSubCommand,
								),
							)
						},
					),
					firstSubCommand.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil),
					firstSubCommand.EXPECT().Events(gomock.Any(), gomock.Any()).Return([]legacy_es.Command{&invokeTestEvent{id: "2"}}, nil),
					secondSubCommand.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil),
					secondSubCommand.EXPECT().Events(gomock.Any(), gomock.Any()).Return([]legacy_es.Command{&invokeTestEvent{id: "3"}}, nil),
					mainCommand.EXPECT().Events(gomock.Any(), gomock.Any()).Return([]legacy_es.Command{&invokeTestEvent{id: "1"}}, nil),
				)

				return mainCommand
			},
			assertCollectedEvents: func(t *testing.T, events []legacy_es.Command) {
				require.Len(t, events, 3)
				assert.Equal(t, "1", events[0].Aggregate().ID)
				assert.Equal(t, "2", events[1].Aggregate().ID)
				assert.Equal(t, "3", events[2].Aggregate().ID)
			},
		},
		{
			name:        "executor batch with sub batch",
			expectedErr: nil,
			command: func(ctrl *gomock.Controller) domain.Executor {
				mainCommand := domainmock.NewMockCommander(ctrl)
				firstSubCommand := domainmock.NewMockCommander(ctrl)
				secondSubCommand := domainmock.NewMockCommander(ctrl)
				thirdSubCommand := domainmock.NewMockCommander(ctrl)
				fourthSubCommand := domainmock.NewMockCommander(ctrl)
				fifthSubCommand := domainmock.NewMockCommander(ctrl)

				gomock.InOrder(
					mainCommand.EXPECT().Execute(gomock.Any(), gomock.Any()).DoAndReturn(
						func(ctx context.Context, opts *domain.InvokeOpts) error {
							return opts.Invoke(ctx,
								domain.BatchExecutors(
									firstSubCommand,
									domain.BatchExecutors(
										secondSubCommand,
										thirdSubCommand,
									),
								),
							)
						},
					),
					firstSubCommand.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil),
					firstSubCommand.EXPECT().Events(gomock.Any(), gomock.Any()).Return([]legacy_es.Command{&invokeTestEvent{id: "2"}}, nil),
					secondSubCommand.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil),
					secondSubCommand.EXPECT().Events(gomock.Any(), gomock.Any()).Return([]legacy_es.Command{&invokeTestEvent{id: "3"}}, nil),
					thirdSubCommand.EXPECT().Execute(gomock.Any(), gomock.Any()).DoAndReturn(
						func(ctx context.Context, opts *domain.InvokeOpts) error {
							return opts.Invoke(ctx,
								domain.BatchExecutors(
									domain.BatchExecutors(
										fourthSubCommand,
									),
									domain.BatchExecutors(
										fifthSubCommand,
									),
								),
							)
						},
					),
					fourthSubCommand.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil),
					fourthSubCommand.EXPECT().Events(gomock.Any(), gomock.Any()).Return([]legacy_es.Command{&invokeTestEvent{id: "5"}}, nil),
					fifthSubCommand.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil),
					fifthSubCommand.EXPECT().Events(gomock.Any(), gomock.Any()).Return([]legacy_es.Command{&invokeTestEvent{id: "6"}}, nil),
					thirdSubCommand.EXPECT().Events(gomock.Any(), gomock.Any()).Return([]legacy_es.Command{&invokeTestEvent{id: "4"}}, nil),
					mainCommand.EXPECT().Events(gomock.Any(), gomock.Any()).Return([]legacy_es.Command{&invokeTestEvent{id: "1"}}, nil),
				)
				return mainCommand
			},
			assertCollectedEvents: func(t *testing.T, events []legacy_es.Command) {
				require.Len(t, events, 6)
				assert.Equal(t, "1", events[0].Aggregate().ID)
				assert.Equal(t, "2", events[1].Aggregate().ID)
				assert.Equal(t, "3", events[2].Aggregate().ID)
				assert.Equal(t, "4", events[3].Aggregate().ID)
				assert.Equal(t, "5", events[4].Aggregate().ID)
				assert.Equal(t, "6", events[5].Aggregate().ID)
			},
		},
		{
			name:        "command with sub query",
			expectedErr: nil,
			command: func(ctrl *gomock.Controller) domain.Executor {
				mainCmd := domainmock.NewMockCommander(ctrl)
				subQuery := domainmock.NewMockQuerier[string](ctrl)
				gomock.InOrder(
					mainCmd.EXPECT().Execute(gomock.Any(), gomock.Any()).DoAndReturn(
						func(ctx context.Context, opts *domain.InvokeOpts) error {
							return opts.Invoke(ctx, subQuery)
						},
					),
					subQuery.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil),
					mainCmd.EXPECT().Events(gomock.Any(), gomock.Any()).Return([]legacy_es.Command{&invokeTestEvent{id: "1"}}, nil),
				)
				return mainCmd
			},
			assertCollectedEvents: func(t *testing.T, events []legacy_es.Command) {
				require.Len(t, events, 1)
				assert.Equal(t, "1", events[0].Aggregate().ID)
			},
		},
		{
			name:        "begin fails",
			expectedErr: assert.AnError,
			command: func(ctrl *gomock.Controller) domain.Executor {
				return domainmock.NewMockCommander(ctrl)
			},
			queryExecutor: func(ctrl *gomock.Controller) database.QueryExecutor {
				pool := dbmock.NewMockPool(ctrl)
				pool.EXPECT().Begin(gomock.Any(), gomock.Any()).Return(nil, assert.AnError).Times(1)
				return pool
			},
			assertCollectedEvents: func(t *testing.T, events []legacy_es.Command) {
				require.Len(t, events, 0)
			},
		},
		{
			name:        "execute fails",
			expectedErr: assert.AnError,
			command: func(ctrl *gomock.Controller) domain.Executor {
				cmd := domainmock.NewMockCommander(ctrl)
				cmd.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(assert.AnError)
				return cmd
			},
			queryExecutor: func(ctrl *gomock.Controller) database.QueryExecutor {
				pool := dbmock.NewMockPool(ctrl)
				tx := dbmock.NewMockTransaction(ctrl)
				// gomock.InOrder(
				pool.EXPECT().Begin(gomock.Any(), gomock.Any()).Return(tx, nil).Times(1)
				tx.EXPECT().End(gomock.Any(), assert.AnError).Return(assert.AnError).Times(1)
				// )
				return pool
			},
			assertCollectedEvents: func(t *testing.T, events []legacy_es.Command) {
				require.Len(t, events, 0)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			i := domain.NewEventStoreInvoker(nil)
			opts := domain.InvokeOpts{
				Invoker: i,
			}
			if tt.queryExecutor != nil {
				domain.WithQueryExecutor(tt.queryExecutor(ctrl))(&opts)
			} else {
				pool := dbmock.NewMockPool(ctrl)
				tx := dbmock.NewMockTransaction(ctrl)
				tx.EXPECT().End(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
				pool.EXPECT().Begin(gomock.Any(), gomock.Any()).Return(tx, nil).Times(1)
				domain.WithQueryExecutor(pool)(&opts)
			}

			eventStore := new(testLegacyEventstore)
			domain.WithLegacyEventstore(eventStore)(&opts)

			err := opts.Invoke(t.Context(), tt.command(ctrl))
			require.ErrorIs(t, err, tt.expectedErr)
			tt.assertCollectedEvents(t, eventStore.commands)
		})
	}
}
