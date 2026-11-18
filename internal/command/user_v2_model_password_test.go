package command

import (
	"context"
	"database/sql"
	"fmt"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/repository/user"
)

func TestUserV2WriteModel_Query_passwordChecks(t *testing.T) {
	t.Run("without password, no password checks", func(t *testing.T) {
		query := NewUserHumanWriteModel("user1", "org1", true, true, true, false, true, true, true).Query()

		require.Len(t, query.GetQueries(), 1)
		for _, eventType := range humanPasswordCheckEventTypes {
			assert.NotContains(t, query.GetQueries()[0].GetEventTypes(), eventType)
		}
	})
	for _, checksAfterSequence := range []uint64{0, 42} {
		t.Run(fmt.Sprintf("with password, only password checks after sequence %d", checksAfterSequence), func(t *testing.T) {
			wm := NewUserHumanWriteModel("user1", "org1", true, true, true, true, true, true, true)
			wm.passwordChecksAfterSequence = checksAfterSequence
			query := wm.Query()

			require.Len(t, query.GetQueries(), 2)
			state, checks := query.GetQueries()[0], query.GetQueries()[1]
			for _, eventType := range humanPasswordCheckEventTypes {
				assert.NotContains(t, state.GetEventTypes(), eventType)
			}
			assert.Zero(t, state.GetSequenceGreater())

			assert.Equal(t, []string{"user1"}, checks.GetAggregateIDs())
			assert.ElementsMatch(t, humanPasswordCheckEventTypes, checks.GetEventTypes())
			assert.Equal(t, checksAfterSequence, checks.GetSequenceGreater())
			assert.Equal(t, "org1", query.GetResourceOwner())
		})
	}
	t.Run("reload only loads new events", func(t *testing.T) {
		wm := NewUserHumanWriteModel("user1", "org1", false, false, false, true, false, false, false)
		assert.Zero(t, wm.Query().GetEventSequenceGreater())

		wm.ProcessedSequence = 7
		assert.Equal(t, uint64(7), wm.Query().GetEventSequenceGreater())

		// the checks up to the reset are not loaded, the reload doesn't need to scan them
		wm.passwordChecksAfterSequence = 42
		assert.Equal(t, uint64(42), wm.Query().GetEventSequenceGreater())

		wm.ProcessedSequence = 50
		assert.Equal(t, uint64(50), wm.Query().GetEventSequenceGreater())
	})
}

// Test_userV2WriteModel_passwordCheckResetEventTypes ensures that the reset event types used to bound the
// password checks are exactly the events which reset the failed password check count in [UserV2WriteModel.Reduce].
func Test_userV2WriteModel_passwordCheckResetEventTypes(t *testing.T) {
	ctx := context.Background()
	userAgg := &user.NewAggregate("user1", "org1").Aggregate
	failed := eventFromEventPusher(user.NewHumanPasswordCheckFailedEvent(ctx, userAgg, nil))
	// only the password option (plus the human and state events, which are always loaded),
	// as e.g. the metadata events can't be reduced without payload
	newWriteModel := func() *UserV2WriteModel {
		return NewUserHumanWriteModel("user1", "org1", false, false, false, true, false, false, false)
	}

	var queried []eventstore.EventType
	for _, query := range newWriteModel().Query().GetQueries() {
		queried = append(queried, query.GetEventTypes()...)
	}
	humanResetTypes := slices.DeleteFunc(slices.Clone(humanPasswordCheckResetEventTypes), func(eventType eventstore.EventType) bool {
		return eventType == user.MachineAddedEventType
	})
	assert.Subset(t, queried, humanResetTypes, "reset events must be loaded by the write model")

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
			wm := newWriteModel()
			require.NoError(t, es.FilterToQueryReducer(ctx, wm))

			isReset := slices.Contains(humanPasswordCheckResetEventTypes, eventType)
			assert.Equal(t, isReset, wm.PasswordCheckFailedCount == 0, "failed count after event: %d", wm.PasswordCheckFailedCount)
		})
	}
}

func TestCommands_UserHumanWriteModel_onlyLoadsPasswordChecksSinceLastReset(t *testing.T) {
	ctx := context.Background()
	userAgg := &user.NewAggregate("user1", "org1").Aggregate

	resetEvent := eventFromEventPusher(user.NewHumanPasswordCheckSucceededEvent(ctx, userAgg, nil))
	resetEvent.Seq = 42
	failedEvent := eventFromEventPusher(user.NewHumanPasswordCheckFailedEvent(ctx, userAgg, nil))
	failedEvent.Seq = 43

	c := &Commands{
		eventstore: expectEventstore(
			// last reset of the failed password check count
			expectFilterQuery(
				func(query *eventstore.SearchQueryBuilder) {
					assert.True(t, query.GetDesc())
					assert.Equal(t, uint64(1), query.GetLimit())
					require.Len(t, query.GetQueries(), 1)
					assert.ElementsMatch(t, humanPasswordCheckResetEventTypes, query.GetQueries()[0].GetEventTypes())
				},
				resetEvent,
			),
			// user and the password checks after the reset
			expectFilterQuery(
				func(query *eventstore.SearchQueryBuilder) {
					require.Len(t, query.GetQueries(), 2)
					assert.Equal(t, uint64(42), query.GetQueries()[1].GetSequenceGreater())
				},
				eventFromEventPusher(newAddHumanEvent("$plain$x$password", false, true, "", language.English)),
				failedEvent,
			),
		)(t),
	}

	wm, err := c.UserHumanWriteModel(ctx, "user1", "org1", false, false, false, true, false, false, false)
	require.NoError(t, err)
	assert.Equal(t, uint64(1), wm.PasswordCheckFailedCount)
}

func TestCommands_UserHumanWriteModel_withoutPasswordNoResetLookup(t *testing.T) {
	ctx := context.Background()
	c := &Commands{
		eventstore: expectEventstore(
			expectFilterQuery(
				func(query *eventstore.SearchQueryBuilder) {
					assert.False(t, query.GetDesc())
					require.Len(t, query.GetQueries(), 1)
				},
				eventFromEventPusher(newAddHumanEvent("$plain$x$password", false, true, "", language.English)),
			),
		)(t),
	}

	_, err := c.UserHumanWriteModel(ctx, "user1", "org1", true, false, false, false, false, false, false)
	require.NoError(t, err)
}
