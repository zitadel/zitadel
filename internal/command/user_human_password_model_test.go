package command

import (
	"context"
	"database/sql"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/eventstore/repository/mock"
	"github.com/zitadel/zitadel/internal/repository/user"
)

func TestHumanPasswordWriteModel_Query(t *testing.T) {
	tests := []struct {
		name                string
		checksAfterSequence uint64
	}{
		{name: "without reset sequence, all password checks", checksAfterSequence: 0},
		{name: "with reset sequence, only password checks after it", checksAfterSequence: 42},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wm := NewHumanPasswordWriteModel("user1", "org1")
			wm.checksAfterSequence = tt.checksAfterSequence
			query := wm.Query()

			require.Len(t, query.GetQueries(), 2)
			state, checks := query.GetQueries()[0], query.GetQueries()[1]

			assert.Equal(t, []string{"user1"}, state.GetAggregateIDs())
			assert.ElementsMatch(t, humanPasswordStateEventTypes, state.GetEventTypes())
			assert.Zero(t, state.GetSequenceGreater())

			assert.Equal(t, []string{"user1"}, checks.GetAggregateIDs())
			assert.ElementsMatch(t, humanPasswordCheckEventTypes, checks.GetEventTypes())
			assert.Equal(t, tt.checksAfterSequence, checks.GetSequenceGreater())
			assert.True(t, checks.GetPositionAfter().IsZero(), "positions can disagree with the sequence order of historical events")

			assert.Equal(t, "org1", query.GetResourceOwner())
		})
	}
}

func TestHumanPasswordWriteModel_Query_reload(t *testing.T) {
	wm := NewHumanPasswordWriteModel("user1", "org1")
	wm.checksAfterSequence = 42
	assert.Zero(t, wm.Query().GetEventSequenceGreater(), "the first load must not filter by sequence")

	// the last reduced event is a state event before the last reset:
	// the reload doesn't need to scan the checks up to the reset
	wm.ProcessedSequence = 7
	assert.Equal(t, uint64(42), wm.Query().GetEventSequenceGreater())

	wm.ProcessedSequence = 50
	assert.Equal(t, uint64(50), wm.Query().GetEventSequenceGreater())
}

// Test_humanPasswordCheckResetEventTypes ensures that the reset event types are exactly the events
// which reset the failed password check count in [HumanPasswordWriteModel.Reduce].
// A reset type which doesn't reset the count would hide failed checks and weaken the lockout.
func Test_humanPasswordCheckResetEventTypes(t *testing.T) {
	ctx := context.Background()
	userAgg := &user.NewAggregate("user1", "org1").Aggregate
	failed := eventFromEventPusher(user.NewHumanPasswordCheckFailedEvent(ctx, userAgg, nil))

	queried := append(slices.Clone(humanPasswordStateEventTypes), humanPasswordCheckEventTypes...)
	assert.Subset(t, queried, humanPasswordCheckResetEventTypes, "reset events must be loaded by the write model")

	for _, eventType := range queried {
		t.Run(string(eventType), func(t *testing.T) {
			event := &repository.Event{
				AggregateID:   "user1",
				AggregateType: user.AggregateType,
				ResourceOwner: sql.NullString{String: "org1", Valid: true},
				Typ:           eventType,
				Data:          []byte("{}"),
			}
			es := expectEventstore(expectFilter(failed, failed, event))(t)
			wm := NewHumanPasswordWriteModel("user1", "org1")
			require.NoError(t, es.FilterToQueryReducer(ctx, wm))

			isReset := slices.Contains(humanPasswordCheckResetEventTypes, eventType)
			assert.Equal(t, isReset, wm.PasswordCheckFailedCount == 0, "failed count after event: %d", wm.PasswordCheckFailedCount)
		})
	}
}

func Test_lastPasswordCheckResetSequence(t *testing.T) {
	ctx := context.Background()
	assertQuery := func(query *eventstore.SearchQueryBuilder) {
		assert.True(t, query.GetDesc())
		assert.Equal(t, uint64(1), query.GetLimit())
		require.Len(t, query.GetQueries(), 1)
		assert.Equal(t, []string{"user1"}, query.GetQueries()[0].GetAggregateIDs())
		assert.ElementsMatch(t, humanPasswordCheckResetEventTypes, query.GetQueries()[0].GetEventTypes())
	}
	t.Run("last reset", func(t *testing.T) {
		event := eventFromEventPusher(user.NewHumanPasswordCheckSucceededEvent(ctx, &user.NewAggregate("user1", "org1").Aggregate, nil))
		event.Seq = 42
		es := expectEventstore(expectFilterQuery(assertQuery, event))(t)

		sequence, err := lastPasswordCheckResetSequence(ctx, es, "user1")
		require.NoError(t, err)
		assert.Equal(t, uint64(42), sequence)
	})
	t.Run("no reset", func(t *testing.T) {
		es := expectEventstore(expectFilterQuery(assertQuery))(t)

		sequence, err := lastPasswordCheckResetSequence(ctx, es, "user1")
		require.NoError(t, err)
		assert.Zero(t, sequence)
	})
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
	resetEvent.Seq = 42

	failedEvent := eventFromEventPusher(user.NewHumanPasswordCheckFailedEvent(ctx, userAgg, nil))
	failedEvent.Seq = 43

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
				assert.Equal(t, uint64(42), query.GetQueries()[1].GetSequenceGreater())
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
