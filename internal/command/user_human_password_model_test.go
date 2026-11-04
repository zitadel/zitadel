package command

import (
	"context"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository/mock"
	"github.com/zitadel/zitadel/internal/repository/user"
)

func TestHumanPasswordWriteModel_Query(t *testing.T) {
	t.Run("without reset position, all password checks", func(t *testing.T) {
		wm := NewHumanPasswordWriteModel("user1", "org1")
		query := wm.Query()

		require.Len(t, query.GetQueries(), 1)
		sub := query.GetQueries()[0]
		assert.Equal(t, []string{"user1"}, sub.GetAggregateIDs())
		assert.Subset(t, sub.GetEventTypes(), humanPasswordStateEventTypes)
		assert.Subset(t, sub.GetEventTypes(), humanPasswordCheckEventTypes)
		assert.True(t, sub.GetPositionAfter().IsZero())
		assert.Equal(t, "org1", query.GetResourceOwner())
	})
	t.Run("with reset position, only password checks after it", func(t *testing.T) {
		wm := NewHumanPasswordWriteModel("user1", "org1")
		wm.checksAfter = decimal.NewFromInt(42)
		query := wm.Query()

		require.Len(t, query.GetQueries(), 2)
		state, checks := query.GetQueries()[0], query.GetQueries()[1]

		assert.Equal(t, []string{"user1"}, state.GetAggregateIDs())
		assert.ElementsMatch(t, humanPasswordStateEventTypes, state.GetEventTypes())
		assert.True(t, state.GetPositionAfter().IsZero())

		assert.Equal(t, []string{"user1"}, checks.GetAggregateIDs())
		assert.ElementsMatch(t, humanPasswordCheckEventTypes, checks.GetEventTypes())
		assert.True(t, checks.GetPositionAfter().Equal(decimal.NewFromInt(42)))

		assert.Equal(t, "org1", query.GetResourceOwner())
	})
}

func TestPasswordCheckResetWriteModel(t *testing.T) {
	wm := newPasswordCheckResetWriteModel("user1")

	query := wm.Query()
	assert.True(t, query.GetDesc())
	assert.Equal(t, uint64(1), query.GetLimit())
	require.Len(t, query.GetQueries(), 1)
	assert.Equal(t, []string{"user1"}, query.GetQueries()[0].GetAggregateIDs())
	assert.ElementsMatch(t, humanPasswordCheckResetEventTypes, query.GetQueries()[0].GetEventTypes())

	event := eventFromEventPusher(user.NewHumanPasswordCheckSucceededEvent(context.Background(), &user.NewAggregate("user1", "org1").Aggregate, nil))
	event.Pos = decimal.NewFromInt(42)
	wm.AppendEvents(event)
	require.NoError(t, wm.Reduce())
	assert.True(t, wm.position.Equal(decimal.NewFromInt(42)))
}

// expectFilterQuery expects a filter call, lets assertQuery check the query and returns the events.
func expectFilterQuery(assertQuery func(*eventstore.SearchQueryBuilder), events ...eventstore.Event) expect {
	return func(m *mock.MockRepository) {
		m.MockQuerier.EXPECT().FilterToReducer(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
			func(_ context.Context, query *eventstore.SearchQueryBuilder, reduce eventstore.Reducer) error {
				assertQuery(query)
				for _, event := range events {
					if err := reduce(event); err != nil {
						return err
					}
				}
				return nil
			},
		)
	}
}

func Test_checkPassword_onlyLoadsChecksSinceLastReset(t *testing.T) {
	ctx := context.Background()
	userAgg := &user.NewAggregate("user1", "org1").Aggregate

	resetEvent := eventFromEventPusher(user.NewHumanPasswordCheckSucceededEvent(ctx, userAgg, nil))
	resetEvent.Pos = decimal.NewFromInt(42)

	failedEvent := eventFromEventPusher(user.NewHumanPasswordCheckFailedEvent(ctx, userAgg, nil))
	failedEvent.Pos = decimal.NewFromInt(43)

	es := expectEventstore(
		// last reset of the failed password check count
		expectFilterQuery(
			func(query *eventstore.SearchQueryBuilder) {
				assert.True(t, query.GetDesc())
				assert.Equal(t, uint64(1), query.GetLimit())
			},
			resetEvent,
		),
		// password state and the password checks after the reset
		expectFilterQuery(
			func(query *eventstore.SearchQueryBuilder) {
				require.Len(t, query.GetQueries(), 2)
				assert.True(t, query.GetQueries()[1].GetPositionAfter().Equal(decimal.NewFromInt(42)))
			},
			eventFromEventPusher(user.NewHumanAddedEvent(ctx, userAgg,
				"username", "firstname", "lastname", "nickname", "displayname",
				language.German, domain.GenderUnspecified, "email@test.ch", false,
			)),
			eventFromEventPusher(user.NewHumanPasswordChangedEvent(ctx, userAgg, "$plain$x$password", false, "")),
			failedEvent,
		),
		expectFilter(), // recheck
	)(t)

	commands, err := checkPassword(ctx, "user1", "password", es, mockPasswordHasher("x"), nil, nil)
	require.NoError(t, err)
	require.Len(t, commands, 1)
	assert.Equal(t, user.HumanPasswordCheckSucceededType, commands[0].Type())
}
