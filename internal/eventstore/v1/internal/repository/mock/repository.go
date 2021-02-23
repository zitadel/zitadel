package mock

import (
	"context"
	"testing"

	"github.com/caos/zitadel/internal/eventstore/v1/models"
	gomock "github.com/golang/mock/gomock"
)

func NewMock(t *testing.T) *MockRepository {
	return NewMockRepository(gomock.NewController(t))
}

func (m *MockRepository) ExpectFilter(query *models.SearchQuery, eventAmount int) *MockRepository {
	events := make([]*models.Event, eventAmount)
	m.EXPECT().Filter(context.Background(), query).Return(events, nil).MaxTimes(1)
	return m
}

func (m *MockRepository) ExpectFilterFail(query *models.SearchQuery, err error) *MockRepository {
	m.EXPECT().Filter(context.Background(), query).Return(nil, err).MaxTimes(1)
	return m
}

func (m *MockRepository) ExpectPush(aggregates ...*models.Aggregate) *MockRepository {
	m.EXPECT().PushAggregates(context.Background(), aggregates).Return(nil).MaxTimes(1)
	return m
}

func (m *MockRepository) ExpectPushError(err error, aggregates ...*models.Aggregate) *MockRepository {
	m.EXPECT().PushAggregates(context.Background(), aggregates).Return(err).MaxTimes(1)
	return m
}
