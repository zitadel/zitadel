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

func (m *MockRepository) ExpectFilterNoEventsNoError() *MockRepository {
	m.EXPECT().Filter(gomock.Any(), gomock.Any()).Return(nil, nil)
	return m
}

func (m *MockRepository) ExpectFilterEvents(events ...*repository.Event) *MockRepository {
	m.EXPECT().Filter(gomock.Any(), gomock.Any()).Return(events, nil)
	return m
}

func (m *MockRepository) ExpectPush(expectedEvents []*repository.Event, expectedAssets []*repository.Asset, expectedUniqueConstraints ...*repository.UniqueConstraint) *MockRepository {
	m.EXPECT().Push(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, events []*repository.Event, assets []*repository.Asset, uniqueConstraints ...*repository.UniqueConstraint) error {
			assert.Equal(m.ctrl.T, expectedEvents, events)
			if expectedUniqueConstraints == nil {
				expectedUniqueConstraints = []*repository.UniqueConstraint{}
			}
			assert.Equal(m.ctrl.T, expectedUniqueConstraints, uniqueConstraints)
			if expectedAssets == nil {
				expectedAssets = []*repository.Asset{}
			}
			assert.Equal(m.ctrl.T, expectedAssets, assets)
			return nil
		},
	)
	return m
}

func (m *MockRepository) ExpectPushFailed(err error, expectedEvents []*repository.Event, expectedAssets []*repository.Asset, expectedUniqueConstraints ...*repository.UniqueConstraint) *MockRepository {
	m.EXPECT().Push(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, events []*repository.Event, assets []*repository.Asset, uniqueConstraints ...*repository.UniqueConstraint) error {
			assert.Equal(m.ctrl.T, expectedEvents, events)
			if expectedUniqueConstraints == nil {
				expectedUniqueConstraints = []*repository.UniqueConstraint{}
			}
			assert.Equal(m.ctrl.T, expectedUniqueConstraints, uniqueConstraints)

			if expectedAssets == nil {
				expectedAssets = []*repository.Asset{}
			}
			assert.Equal(m.ctrl.T, expectedAssets, assets)
			return err
		},
	)
	return m
}
