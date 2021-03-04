package mock

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/caos/zitadel/internal/eventstore/repository"
)

func NewRepo(t *testing.T) *MockRepository {
	return NewMockRepository(gomock.NewController(t))
}

func ExpectFilterNoEventsNoError() {

}

func (m *MockRepository) ExpectFilterNoEventsNoError() *MockRepository {
	m.EXPECT().Filter(gomock.Any(), gomock.Any()).Return(nil, nil)
	return m
}

func (m *MockRepository) ExpectFilterEvents2(expectedSearchQuery *repository.SearchQuery, events ...*repository.Event) *MockRepository {
	m.EXPECT().Filter(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, searchQuery *repository.SearchQuery) ([]*repository.Event, error) {
			assert.Equal(m.ctrl.T, expectedSearchQuery, searchQuery)
			return events, nil
		})
	return m
}

func (m *MockRepository) ExpectFilterEvents(events ...*repository.Event) *MockRepository {
	m.EXPECT().Filter(gomock.Any(), gomock.Any()).Return(events, nil)
	return m
}

func (m *MockRepository) ExpectPush(expectedEvents []*repository.Event, expectedUniqueConstraints ...*repository.UniqueConstraint) *MockRepository {
	m.EXPECT().Push(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, events []*repository.Event, uniqueConstraints ...*repository.UniqueConstraint) error {
			assert.Equal(m.ctrl.T, expectedEvents, events)
			assert.Equal(m.ctrl.T, expectedUniqueConstraints, uniqueConstraints)
			return nil
		},
	)
	return m
}

func (m *MockRepository) ExpectPushFailed(err error, expectedEvents []*repository.Event, expectedUniqueConstraints ...*repository.UniqueConstraint) *MockRepository {
	m.EXPECT().Push(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, events []*repository.Event, uniqueConstraints ...*repository.UniqueConstraint) error {
			assert.Equal(m.ctrl.T, expectedEvents, events)
			assert.Equal(m.ctrl.T, expectedUniqueConstraints, uniqueConstraints)
			return err
		},
	)
	return m
}
