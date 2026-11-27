package command

import (
	"context"
	"database/sql"
	"slices"
	"testing"
	"time"

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
		name                    string
		checksAfterSequence     uint64
		passwordChangedSequence uint64
		wantHashSequenceGreater uint64
	}{
		{name: "without bounds, all password checks and hash events", checksAfterSequence: 0, passwordChangedSequence: 0, wantHashSequenceGreater: 0},
		{name: "with reset sequence, only password checks after it", checksAfterSequence: 42, passwordChangedSequence: 0, wantHashSequenceGreater: 0},
		{name: "with password changed sequence, only hash events from the change on", checksAfterSequence: 42, passwordChangedSequence: 40, wantHashSequenceGreater: 39},
		{name: "password changed by the first event, all hash events", checksAfterSequence: 1, passwordChangedSequence: 1, wantHashSequenceGreater: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wm := NewHumanPasswordWriteModel("user1", "org1")
			wm.checksAfterSequence = tt.checksAfterSequence
			wm.passwordChangedSequence = tt.passwordChangedSequence
			query := wm.Query()

			require.Len(t, query.GetQueries(), 3)
			state, hashes, checks := query.GetQueries()[0], query.GetQueries()[1], query.GetQueries()[2]

			assert.Equal(t, []string{"user1"}, state.GetAggregateIDs())
			assert.ElementsMatch(t, humanPasswordStateEventTypes, state.GetEventTypes())
			assert.Zero(t, state.GetSequenceGreater())

			assert.Equal(t, []string{"user1"}, hashes.GetAggregateIDs())
			assert.ElementsMatch(t, humanPasswordHashEventTypes, hashes.GetEventTypes())
			assert.Equal(t, tt.wantHashSequenceGreater, hashes.GetSequenceGreater())
			assert.True(t, hashes.GetPositionAfter().IsZero(), "positions can disagree with the sequence order of historical events")

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

	queried := slices.Concat(humanPasswordStateEventTypes, humanPasswordHashEventTypes, humanPasswordCheckEventTypes)
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

// Test_humanPasswordHashEventTypes ensures that skipping the hash events before the last password change
// doesn't change the reduced state: a change must overwrite everything the hash events set.
func Test_humanPasswordHashEventTypes(t *testing.T) {
	ctx := context.Background()
	userAgg := &user.NewAggregate("user1", "org1").Aggregate

	queried := slices.Concat(humanPasswordStateEventTypes, humanPasswordHashEventTypes, humanPasswordCheckEventTypes)
	assert.Len(t, queried, len(slices.Compact(slices.Sorted(slices.Values(queried)))), "event types must only be loaded by one sub query")
	assert.Subset(t, humanPasswordCheckResetEventTypes, humanPasswordChangedEventTypes, "a password change must be found by the reset lookup")

	added := eventFromEventPusher(user.NewHumanAddedEvent(ctx, userAgg,
		"username", "firstname", "lastname", "nickname", "displayname",
		language.German, domain.GenderUnspecified, "email@test.ch", false,
	))
	codeAdded := eventFromEventPusher(user.NewHumanPasswordCodeAddedEvent(ctx, userAgg, nil, time.Hour, domain.NotificationTypeEmail, "", ""))
	older := []eventstore.Event{
		eventFromEventPusher(user.NewHumanPasswordChangedEvent(ctx, userAgg, "$plain$x$first", true, "")),
		eventFromEventPusher(user.NewHumanPasswordHashUpdatedEvent(ctx, userAgg, "$plain$x$updated")),
		eventFromEventPusher(user.NewHumanPasswordCheckFailedEvent(ctx, userAgg, nil)),
	}
	for _, changed := range []eventstore.Event{
		eventFromEventPusher(user.NewHumanPasswordChangedEvent(ctx, userAgg, "$plain$x$last", false, "")),
		&repository.Event{
			AggregateID:   "user1",
			AggregateType: user.AggregateType,
			ResourceOwner: sql.NullString{String: "org1", Valid: true},
			Typ:           user.UserV1PasswordChangedType,
			Data:          []byte(`{"encodedHash":"$plain$x$last"}`),
		},
	} {
		t.Run(string(changed.Type()), func(t *testing.T) {
			all := NewHumanPasswordWriteModel("user1", "org1")
			require.NoError(t, expectEventstore(expectFilter(slices.Concat([]eventstore.Event{added}, older, []eventstore.Event{codeAdded, changed})...))(t).FilterToQueryReducer(ctx, all))

			bounded := NewHumanPasswordWriteModel("user1", "org1")
			require.NoError(t, expectEventstore(expectFilter(added, codeAdded, changed))(t).FilterToQueryReducer(ctx, bounded))

			assert.Equal(t, "$plain$x$last", bounded.EncodedHash)
			assert.Equal(t, all.EncodedHash, bounded.EncodedHash)
			assert.Equal(t, all.SecretChangeRequired, bounded.SecretChangeRequired)
			assert.Equal(t, all.Code, bounded.Code)
			assert.Equal(t, all.PasswordCheckFailedCount, bounded.PasswordCheckFailedCount)
			assert.Equal(t, all.UserState, bounded.UserState)
		})
	}
}

func Test_lastHumanPasswordLoadBounds(t *testing.T) {
	ctx := context.Background()
	userAgg := &user.NewAggregate("user1", "org1").Aggregate
	assertQuery := func(query *eventstore.SearchQueryBuilder) {
		assert.True(t, query.GetDesc())
		assert.Equal(t, uint64(humanPasswordBoundsLookupLimit), query.GetLimit())
		require.Len(t, query.GetQueries(), 1)
		assert.Equal(t, []string{"user1"}, query.GetQueries()[0].GetAggregateIDs())
		assert.ElementsMatch(t, humanPasswordCheckResetEventTypes, query.GetQueries()[0].GetEventTypes())
	}
	withSequence := func(event *repository.Event, sequence uint64) *repository.Event {
		event.Seq = sequence
		return event
	}
	tests := []struct {
		name                            string
		events                          []eventstore.Event
		wantLastResetSequence           uint64
		wantLastPasswordChangedSequence uint64
	}{
		{
			name: "no reset",
		},
		{
			name: "last reset is the password change",
			events: []eventstore.Event{
				withSequence(eventFromEventPusher(user.NewHumanPasswordChangedEvent(ctx, userAgg, "$plain$x$password", false, "")), 42),
				withSequence(eventFromEventPusher(user.NewHumanPasswordCheckSucceededEvent(ctx, userAgg, nil)), 40),
			},
			wantLastResetSequence:           42,
			wantLastPasswordChangedSequence: 42,
		},
		{
			name: "password change before the last reset",
			events: []eventstore.Event{
				withSequence(eventFromEventPusher(user.NewHumanPasswordCheckSucceededEvent(ctx, userAgg, nil)), 45),
				withSequence(eventFromEventPusher(user.NewHumanPasswordChangedEvent(ctx, userAgg, "$plain$x$password", false, "")), 42),
				withSequence(eventFromEventPusher(user.NewHumanPasswordChangedEvent(ctx, userAgg, "$plain$x$older", false, "")), 30),
			},
			wantLastResetSequence:           45,
			wantLastPasswordChangedSequence: 42,
		},
		{
			name: "no password change among the last resets",
			events: []eventstore.Event{
				withSequence(eventFromEventPusher(user.NewHumanPasswordCheckSucceededEvent(ctx, userAgg, nil)), 45),
				withSequence(eventFromEventPusher(user.NewHumanPasswordCheckSucceededEvent(ctx, userAgg, nil)), 43),
			},
			wantLastResetSequence:           45,
			wantLastPasswordChangedSequence: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			es := expectEventstore(expectFilterQuery(assertQuery, tt.events...))(t)

			bounds, err := lastHumanPasswordLoadBounds(ctx, es, "user1")
			require.NoError(t, err)
			assert.Equal(t, tt.wantLastResetSequence, bounds.lastResetSequence)
			assert.Equal(t, tt.wantLastPasswordChangedSequence, bounds.lastPasswordChangedSequence)
		})
	}
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

	changedEvent := eventFromEventPusher(user.NewHumanPasswordChangedEvent(ctx, userAgg, "$plain$x$password", false, ""))
	changedEvent.Seq = 30

	failedEvent := eventFromEventPusher(user.NewHumanPasswordCheckFailedEvent(ctx, userAgg, nil))
	failedEvent.Seq = 43

	es := expectEventstore(
		// last resets of the failed password check count, including the last password change
		expectFilterQuery(
			func(query *eventstore.SearchQueryBuilder) {
				assert.True(t, query.GetDesc())
				assert.Equal(t, uint64(humanPasswordBoundsLookupLimit), query.GetLimit())
			},
			resetEvent,
			changedEvent,
		),
		// password state, the hash events from the last change on and the password checks after the reset
		expectFilterQuery(
			func(query *eventstore.SearchQueryBuilder) {
				require.Len(t, query.GetQueries(), 3)
				assert.Zero(t, query.GetQueries()[0].GetSequenceGreater())
				assert.Equal(t, uint64(29), query.GetQueries()[1].GetSequenceGreater())
				assert.Equal(t, uint64(42), query.GetQueries()[2].GetSequenceGreater())
			},
			eventFromEventPusher(user.NewHumanAddedEvent(ctx, userAgg,
				"username", "firstname", "lastname", "nickname", "displayname",
				language.German, domain.GenderUnspecified, "email@test.ch", false,
			)),
			changedEvent,
			failedEvent,
		),
		expectFilter(), // recheck
	)(t)

	commands, err := checkPassword(ctx, "user1", "password", es, mockPasswordHasher("x"), nil, nil)
	require.NoError(t, err)
	require.Len(t, commands, 1)
	assert.Equal(t, user.HumanPasswordCheckSucceededType, commands[0].Type())
}
