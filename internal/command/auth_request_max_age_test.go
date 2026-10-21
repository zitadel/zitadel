package command

import (
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/repository/authrequest"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommands_LinkSessionToAuthRequest_maxAge(t *testing.T) {
	requestCreatedAt := time.Date(2026, time.January, 1, 12, 0, 0, 0, time.UTC)
	tests := []struct {
		name     string
		maxAge   *time.Duration
		authTime time.Time
		wantErr  bool
	}{
		{name: "no limit", authTime: requestCreatedAt.Add(-time.Hour)},
		{name: "stale authentication", maxAge: gu.Ptr(10 * time.Second), authTime: requestCreatedAt.Add(-time.Hour), wantErr: true},
		{name: "recent authentication", maxAge: gu.Ptr(10 * time.Second), authTime: requestCreatedAt.Add(-5 * time.Second)},
		{name: "exact age limit", maxAge: gu.Ptr(10 * time.Second), authTime: requestCreatedAt.Add(-10 * time.Second)},
		{name: "zero requires reauthentication", maxAge: gu.Ptr(time.Duration(0)), authTime: requestCreatedAt.Add(-time.Second), wantErr: true},
		{name: "zero accepts reauthentication", maxAge: gu.Ptr(time.Duration(0)), authTime: requestCreatedAt.Add(time.Second)},
		{name: "missing authentication", maxAge: gu.Ptr(10 * time.Second), wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := authz.NewMockContext("instanceID", "orgID", "loginClient")
			authAggregate := &authrequest.NewAggregate("V2_max_age", "instanceID").Aggregate
			sessionAggregate := &session.NewAggregate("sessionID", "instanceID").Aggregate
			added := authrequest.NewAddedEvent(ctx, authAggregate,
				"loginClient", "clientID", "https://example.com/callback", "state", "nonce",
				[]string{"openid"}, []string{"clientID"}, domain.OIDCResponseTypeCode,
				domain.OIDCResponseModeQuery, nil, nil, nil, tt.maxAge, nil, nil, false,
				"https://issuer.example.com", "")
			expects := []expect{
				expectFilter(eventFromEventPusherWithCreationDate(added, requestCreatedAt)),
				expectFilter(
					eventFromEventPusher(session.NewAddedEvent(ctx, sessionAggregate, nil)),
					eventFromEventPusher(session.NewUserCheckedEvent(ctx, sessionAggregate, "userID", "orgID", tt.authTime, nil)),
					eventFromEventPusher(session.NewPasswordCheckedEvent(ctx, sessionAggregate, tt.authTime)),
				),
			}
			if !tt.wantErr {
				expects = append(expects, expectPush(authrequest.NewSessionLinkedEvent(ctx, authAggregate,
					"sessionID", "userID", tt.authTime, []domain.UserAuthMethodType{domain.UserAuthMethodTypePassword})))
			}
			c := &Commands{
				eventstore:           expectEventstore(expects...)(t),
				sessionTokenVerifier: newMockTokenVerifierValid(),
			}
			_, linked, err := c.LinkSessionToAuthRequest(ctx, "V2_max_age", "sessionID", "valid-session-token", true, nil)
			if tt.wantErr {
				require.True(t, zerrors.IsPreconditionFailed(err), "want reauthentication required, got %v", err)
				require.Nil(t, linked)
				return
			}
			require.NoError(t, err)
			require.Equal(t, "sessionID", linked.SessionID)
			require.WithinDuration(t, tt.authTime, linked.AuthTime, 0)
		})
	}
}
