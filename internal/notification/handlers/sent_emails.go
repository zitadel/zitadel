package handlers

import (
	"context"
	"time"

	"github.com/shopspring/decimal"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/notification/senders"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/repository/user"
)

// emailSentTypes are the events stating that an email notification was sent.
// They carry the [senders.DeliveryInfo] of the email.
var emailSentTypes = []eventstore.EventType{
	user.UserV1InitialCodeSentType,
	user.HumanInitialCodeSentType,
	user.UserV1EmailCodeSentType,
	user.HumanEmailCodeSentType,
	user.HumanInviteCodeSentType,
	user.UserV1PasswordCodeSentType,
	user.HumanPasswordCodeSentType,
	user.HumanPasswordChangeSentType,
	user.HumanOTPEmailCodeSentType,
	user.HumanPasswordlessInitCodeSentType,
	user.UserDomainClaimedSentType,
	session.OTPEmailSentType,
}

// SentEmails returns the amount of emails the instance sent through the provider since the given time.
// At most limit emails are counted. Suppressed emails do not state a provider and are therefore not counted.
func (n *NotificationQueries) SentEmails(ctx context.Context, providerID string, since time.Time, limit uint64) (uint64, error) {
	sent := &sentEmails{
		instanceID: authz.GetInstance(ctx).InstanceID(),
		providerID: providerID,
		since:      since,
		limit:      limit,
	}
	if err := n.es.FilterToQueryReducer(ctx, sent); err != nil {
		return 0, err
	}
	return sent.count, nil
}

type sentEmails struct {
	instanceID string
	providerID string
	since      time.Time
	limit      uint64

	count uint64
}

func (s *sentEmails) Reduce() error {
	return nil
}

func (s *sentEmails) AppendEvents(events ...eventstore.Event) {
	s.count += uint64(len(events))
}

func (s *sentEmails) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		InstanceID(s.instanceID).
		// the position of an event is the time it was created at
		PositionAtLeast(decimal.NewFromInt(s.since.Unix())).
		Limit(s.limit).
		AddQuery().
		AggregateTypes(user.AggregateType, session.AggregateType).
		EventTypes(emailSentTypes...).
		EventData(senders.DeliveredBy(s.providerID)).
		Builder()
}
