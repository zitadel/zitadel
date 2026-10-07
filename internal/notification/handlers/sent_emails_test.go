package handlers

import (
	"encoding/json"
	"io"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore"
	es_repo_mock "github.com/zitadel/zitadel/internal/eventstore/repository/mock"
	"github.com/zitadel/zitadel/internal/notification/senders"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/repository/user"
)

// Test_sentEmails_Query ensures that the sent events of all email notifications are counted
// and that the filter on the payload matches what the sent events store.
func Test_sentEmails_Query(t *testing.T) {
	ctx := t.Context()
	userAgg := &user.NewAggregate(userID, orgID).Aggregate
	sessionAgg := &session.NewAggregate("session1", instanceID).Aggregate
	delivered := senders.DeliveryInfo{ProviderID: "provider1"}

	query := (&sentEmails{
		instanceID: instanceID,
		providerID: "provider1",
		since:      time.Now().Add(-time.Hour),
		limit:      10,
	}).Query()

	sent := []eventstore.Command{
		user.NewHumanInitialCodeSentEvent(ctx, userAgg, delivered),
		user.NewHumanEmailCodeSentEvent(ctx, userAgg, delivered),
		user.NewHumanInviteCodeSentEvent(ctx, userAgg, delivered),
		user.NewHumanPasswordCodeSentEvent(ctx, userAgg, nil, delivered),
		user.NewHumanPasswordChangeSentEvent(ctx, userAgg, delivered),
		user.NewHumanOTPEmailCodeSentEvent(ctx, userAgg, delivered),
		user.NewHumanPasswordlessInitCodeSentEvent(ctx, userAgg, "code1", delivered),
		user.NewDomainClaimedSentEvent(ctx, userAgg, delivered),
		session.NewOTPEmailSentEvent(ctx, sessionAgg, delivered),
	}
	assert.Len(t, query.Matches(sent...), len(sent), "every email sent event is counted")

	notCounted := []eventstore.Command{
		// SMS are not sent through the email provider
		user.NewHumanOTPSMSCodeSentEvent(ctx, userAgg, nil),
		session.NewOTPSMSSentEvent(ctx, sessionAgg, nil),
		// a rejected notification was not sent
		user.NewHumanInviteCodeAddedEvent(ctx, userAgg, nil, time.Hour, "", false, "", ""),
	}
	assert.Empty(t, query.Matches(notCounted...))

	// the filter on the payload must match what the sent events of delivered emails store,
	// and must not match suppressed emails, which do not state a provider
	for _, tt := range []struct {
		name string
		info senders.DeliveryInfo
		want bool
	}{
		{name: "delivered through the provider", info: delivered, want: true},
		{name: "delivered through another provider", info: senders.DeliveryInfo{ProviderID: "provider2"}, want: false},
		{name: "suppressed", info: senders.DeliveryInfo{DeliverySuppressed: true}, want: false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			payload, err := json.Marshal(user.NewHumanInviteCodeSentEvent(ctx, userAgg, tt.info).Payload())
			require.NoError(t, err)
			stored := make(map[string]any)
			require.NoError(t, json.Unmarshal(payload, &stored))
			matches := true
			for key, value := range senders.DeliveredBy("provider1") {
				if stored[key] != value {
					matches = false
				}
			}
			assert.Equal(t, tt.want, matches)
		})
	}
}

func TestNotificationQueries_SentEmails(t *testing.T) {
	ctx := authz.WithInstanceID(t.Context(), instanceID)
	userAgg := &user.NewAggregate(userID, orgID).Aggregate
	sent := func(n int) []eventstore.Event {
		events := make([]eventstore.Event, n)
		for i := range events {
			events[i] = user.NewHumanInviteCodeSentEvent(ctx, userAgg, senders.DeliveryInfo{ProviderID: "provider1"})
		}
		return events
	}

	tests := []struct {
		name      string
		querier   func(t *testing.T) *es_repo_mock.MockRepository
		wantCount uint64
		wantErr   error
	}{
		{
			name:      "nothing sent",
			querier:   func(t *testing.T) *es_repo_mock.MockRepository { return es_repo_mock.NewRepo(t).ExpectFilterEvents() },
			wantCount: 0,
		},
		{
			name: "sent",
			querier: func(t *testing.T) *es_repo_mock.MockRepository {
				return es_repo_mock.NewRepo(t).ExpectFilterEvents(sent(3)...)
			},
			wantCount: 3,
		},
		{
			name: "query fails",
			querier: func(t *testing.T) *es_repo_mock.MockRepository {
				return es_repo_mock.NewRepo(t).ExpectFilterEventsError(io.ErrClosedPipe)
			},
			wantErr: io.ErrClosedPipe,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NewNotificationQueries(nil,
				eventstore.NewEventstore(&eventstore.Config{Querier: tt.querier(t).MockQuerier}),
				externalDomain, externalPort, externalSecure, "", nil, nil, nil, nil,
			)
			count, err := n.SentEmails(ctx, "provider1", time.Now().Add(-time.Hour), 10)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.wantCount, count)
		})
	}
}

// Test_sentEmails_Query_bounds ensures that only the events of the instance within the window are read,
// and not more of them than the limit.
func Test_sentEmails_Query_bounds(t *testing.T) {
	since := time.Date(2026, time.October, 5, 12, 0, 0, 0, time.UTC)
	query := (&sentEmails{instanceID: instanceID, providerID: "provider1", since: since, limit: 10}).Query()

	require.NotNil(t, query.GetInstanceID())
	assert.Equal(t, instanceID, *query.GetInstanceID())
	assert.Equal(t, uint64(10), query.GetLimit())
	// the position of an event is the time it was created at in seconds
	assert.True(t, query.GetPositionAtLeast().Equal(decimal.NewFromInt(since.Unix())), query.GetPositionAtLeast().String())
	require.Len(t, query.GetQueries(), 1)
	assert.Equal(t, senders.DeliveredBy("provider1"), query.GetQueries()[0].GetEventData())
}
