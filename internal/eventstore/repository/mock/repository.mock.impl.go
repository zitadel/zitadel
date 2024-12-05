package mock

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
)

type MockRepository struct {
	*MockPusher
	*MockQuerier
}

func NewRepo(t *testing.T) *MockRepository {
	controller := gomock.NewController(t)
	return &MockRepository{
		MockPusher:  NewMockPusher(controller),
		MockQuerier: NewMockQuerier(controller),
	}
}

func (m *MockRepository) ExpectFilterNoEventsNoError() *MockRepository {
	m.MockQuerier.ctrl.T.Helper()

	m.MockQuerier.EXPECT().FilterToReducer(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	return m
}

func (m *MockRepository) ExpectFilterEvents(events ...eventstore.Event) *MockRepository {
	m.MockQuerier.ctrl.T.Helper()

	m.MockQuerier.EXPECT().FilterToReducer(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, _ *eventstore.SearchQueryBuilder, reduce eventstore.Reducer) error {
			for _, event := range events {
				if err := reduce(event); err != nil {
					return err
				}
			}
			return nil
		},
	)
	return m
}

func (m *MockRepository) ExpectFilterEventsError(err error) *MockRepository {
	m.MockQuerier.ctrl.T.Helper()

	m.MockQuerier.EXPECT().FilterToReducer(gomock.Any(), gomock.Any(), gomock.Any()).Return(err)
	return m
}

func (m *MockRepository) ExpectInstanceIDs(hasFilters []*repository.Filter, instanceIDs ...string) *MockRepository {
	m.MockQuerier.ctrl.T.Helper()

	matcher := gomock.Any()
	if len(hasFilters) > 0 {
		matcher = &filterQueryMatcher{SubQueries: [][]*repository.Filter{hasFilters}}
	}
	m.MockQuerier.EXPECT().InstanceIDs(gomock.Any(), matcher).Return(instanceIDs, nil)
	return m
}

func (m *MockRepository) ExpectInstanceIDsError(err error) *MockRepository {
	m.MockQuerier.ctrl.T.Helper()

	m.MockQuerier.EXPECT().InstanceIDs(gomock.Any(), gomock.Any()).Return(nil, err)
	return m
}

// ExpectPush checks if the expectedCommands are send to the Push method.
// The call will sleep at least the amount of passed duration.
func (m *MockRepository) ExpectPush(expectedCommands []eventstore.Command, sleep time.Duration) *MockRepository {
	m.MockPusher.EXPECT().Push(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, _ database.ContextQueryExecuter, commands ...eventstore.Command) ([]eventstore.Event, error) {
			m.MockPusher.ctrl.T.Helper()

			time.Sleep(sleep)

			if len(expectedCommands) != len(commands) {
				return nil, fmt.Errorf("unexpected amount of commands: want %d, got %d", len(expectedCommands), len(commands))
			}
			for i, expectedCommand := range expectedCommands {
				if !assert.Equal(m.MockPusher.ctrl.T, expectedCommand.Aggregate(), commands[i].Aggregate()) {
					m.MockPusher.ctrl.T.Errorf("invalid command.Aggregate [%d]: expected: %#v got: %#v", i, expectedCommand.Aggregate(), commands[i].Aggregate())
				}
				if !assert.Equal(m.MockPusher.ctrl.T, expectedCommand.Creator(), commands[i].Creator()) {
					m.MockPusher.ctrl.T.Errorf("invalid command.Creator [%d]: expected: %#v got: %#v", i, expectedCommand.Creator(), commands[i].Creator())
				}
				if !assert.Equal(m.MockPusher.ctrl.T, expectedCommand.Type(), commands[i].Type()) {
					m.MockPusher.ctrl.T.Errorf("invalid command.Type [%d]: expected: %#v got: %#v", i, expectedCommand.Type(), commands[i].Type())
				}
				if !assert.Equal(m.MockPusher.ctrl.T, expectedCommand.Revision(), commands[i].Revision()) {
					m.MockPusher.ctrl.T.Errorf("invalid command.Revision [%d]: expected: %#v got: %#v", i, expectedCommand.Revision(), commands[i].Revision())
				}

				var expectedPayload []byte
				expectedPayload, ok := expectedCommand.Payload().([]byte)
				if !ok {
					expectedPayload, _ = json.Marshal(expectedCommand.Payload())
				}
				if string(expectedPayload) == "" {
					expectedPayload = []byte("null")
				}
				gotPayload, _ := json.Marshal(commands[i].Payload())

				if !assert.Equal(m.MockPusher.ctrl.T, expectedPayload, gotPayload) {
					m.MockPusher.ctrl.T.Errorf("invalid command.Payload [%d]: expected: %#v got: %#v", i, expectedCommand.Payload(), commands[i].Payload())
				}
				if !assert.ElementsMatch(m.MockPusher.ctrl.T, expectedCommand.UniqueConstraints(), commands[i].UniqueConstraints()) {
					m.MockPusher.ctrl.T.Errorf("invalid command.UniqueConstraints [%d]: expected: %#v got: %#v", i, expectedCommand.UniqueConstraints(), commands[i].UniqueConstraints())
				}
			}
			events := make([]eventstore.Event, len(commands))
			for i, command := range commands {
				events[i] = &mockEvent{
					Command: command,
				}
			}
			return events, nil
		},
	)
	return m
}

func (m *MockRepository) ExpectPushFailed(err error, expectedCommands []eventstore.Command) *MockRepository {
	m.MockPusher.ctrl.T.Helper()

	m.MockPusher.EXPECT().Push(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, _ database.ContextQueryExecuter, commands ...eventstore.Command) ([]eventstore.Event, error) {
			if len(expectedCommands) != len(commands) {
				return nil, fmt.Errorf("unexpected amount of commands: want %d, got %d", len(expectedCommands), len(commands))
			}
			for i, expectedCommand := range expectedCommands {
				assert.Equal(m.MockPusher.ctrl.T, expectedCommand.Aggregate(), commands[i].Aggregate())
				assert.Equal(m.MockPusher.ctrl.T, expectedCommand.Creator(), commands[i].Creator())
				assert.Equal(m.MockPusher.ctrl.T, expectedCommand.Type(), commands[i].Type())
				assert.Equal(m.MockPusher.ctrl.T, expectedCommand.Revision(), commands[i].Revision())
				assert.Equal(m.MockPusher.ctrl.T, expectedCommand.Payload(), commands[i].Payload())
				assert.ElementsMatch(m.MockPusher.ctrl.T, expectedCommand.UniqueConstraints(), commands[i].UniqueConstraints())
			}

			return nil, err
		},
	)
	return m
}

type mockEvent struct {
	eventstore.Command
	sequence  uint64
	createdAt time.Time
}

// DataAsBytes implements eventstore.Event
func (e *mockEvent) DataAsBytes() []byte {
	if e.Payload() == nil {
		return nil
	}
	payload, err := json.Marshal(e.Payload())
	if err != nil {
		panic(err)
	}
	return payload
}

func (e *mockEvent) Unmarshal(ptr any) error {
	if e.Payload() == nil {
		return nil
	}
	payload, err := json.Marshal(e.Payload())
	if err != nil {
		return err
	}
	return json.Unmarshal(payload, ptr)
}

func (e *mockEvent) Sequence() uint64 {
	return e.sequence
}

func (e *mockEvent) Position() float64 {
	return 0
}

func (e *mockEvent) CreatedAt() time.Time {
	return e.createdAt
}

func (m *MockRepository) ExpectRandomPush(expectedCommands []eventstore.Command) *MockRepository {
	m.MockPusher.EXPECT().Push(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, _ database.ContextQueryExecuter, commands ...eventstore.Command) ([]eventstore.Event, error) {
			assert.Len(m.MockPusher.ctrl.T, commands, len(expectedCommands))

			events := make([]eventstore.Event, len(commands))
			for i, command := range commands {
				events[i] = &mockEvent{
					Command: command,
				}
			}

			return events, nil
		},
	)
	return m
}

func (m *MockRepository) ExpectRandomPushFailed(err error, expectedEvents []eventstore.Command) *MockRepository {
	m.MockPusher.EXPECT().Push(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, _ database.ContextQueryExecuter, events ...eventstore.Command) ([]eventstore.Event, error) {
			assert.Len(m.MockPusher.ctrl.T, events, len(expectedEvents))
			return nil, err
		},
	)
	return m
}
