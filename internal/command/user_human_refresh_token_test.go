package command

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zitadel/oidc/v3/pkg/oidc"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommands_RevokeRefreshToken(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx     context.Context
		userID  string
		orgID   string
		tokenID string
	}
	type res struct {
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"missing param, error",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{},
			res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			"token not active, error",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
				),
			},
			args{
				context.Background(),
				"userID",
				"orgID",
				"tokenID",
			},
			res{
				err: zerrors.IsNotFound,
			},
		},
		{
			"push failed, error",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(user.NewHumanRefreshTokenAddedEvent(
							context.Background(),
							&user.NewAggregate("userID", "orgID").Aggregate,
							"tokenID",
							"clientID",
							"agentID",
							"de",
							[]string{"clientID1"},
							[]string{oidc.ScopeOpenID, oidc.ScopeProfile, oidc.ScopeEmail, oidc.ScopeOfflineAccess},
							[]string{"password"},
							time.Now(),
							1*time.Hour,
							10*time.Hour,
							nil,
						)),
					),
					expectPushFailed(zerrors.ThrowInternal(nil, "ERROR", "internal"),
						user.NewHumanRefreshTokenRemovedEvent(
							context.Background(),
							&user.NewAggregate("userID", "orgID").Aggregate,
							"tokenID",
						),
					),
				),
			},
			args{
				context.Background(),
				"userID",
				"orgID",
				"tokenID",
			},
			res{
				err: zerrors.IsInternal,
			},
		},
		{
			"revoke, ok",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(user.NewHumanRefreshTokenAddedEvent(
							context.Background(),
							&user.NewAggregate("userID", "orgID").Aggregate,
							"tokenID",
							"clientID",
							"agentID",
							"de",
							[]string{"clientID1"},
							[]string{oidc.ScopeOpenID, oidc.ScopeProfile, oidc.ScopeEmail, oidc.ScopeOfflineAccess},
							[]string{"password"},
							time.Now(),
							1*time.Hour,
							10*time.Hour,
							nil,
						)),
					),
					expectPush(
						user.NewHumanRefreshTokenRemovedEvent(
							context.Background(),
							&user.NewAggregate("userID", "orgID").Aggregate,
							"tokenID",
						),
					),
				),
			},
			args{
				context.Background(),
				"userID",
				"orgID",
				"tokenID",
			},
			res{
				want: &domain.ObjectDetails{
					ResourceOwner: "orgID",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := c.RevokeRefreshToken(tt.args.ctx, tt.args.userID, tt.args.orgID, tt.args.tokenID)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}

func TestCommands_RevokeRefreshTokens(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx      context.Context
		userID   string
		orgID    string
		tokenIDs []string
	}
	type res struct {
		err func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"missing tokenIDs, error",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				context.Background(),
				"userID",
				"orgID",
				nil,
			},
			res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			"one token not active, error",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(user.NewHumanRefreshTokenAddedEvent(
							context.Background(),
							&user.NewAggregate("userID", "orgID").Aggregate,
							"tokenID",
							"clientID",
							"agentID",
							"de",
							[]string{"clientID1"},
							[]string{oidc.ScopeOpenID, oidc.ScopeProfile, oidc.ScopeEmail, oidc.ScopeOfflineAccess},
							[]string{"password"},
							time.Now(),
							1*time.Hour,
							10*time.Hour,
							nil,
						)),
					),
					expectFilter(),
				),
			},
			args{
				context.Background(),
				"userID",
				"orgID",
				[]string{"tokenID", "tokenID2"},
			},
			res{
				err: zerrors.IsNotFound,
			},
		},
		{
			"push failed, error",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(user.NewHumanRefreshTokenAddedEvent(
							context.Background(),
							&user.NewAggregate("userID", "orgID").Aggregate,
							"tokenID",
							"clientID",
							"agentID",
							"de",
							[]string{"clientID"},
							[]string{oidc.ScopeOpenID, oidc.ScopeProfile, oidc.ScopeEmail, oidc.ScopeOfflineAccess},
							[]string{"password"},
							time.Now(),
							1*time.Hour,
							10*time.Hour,
							nil,
						)),
					),
					expectFilter(
						eventFromEventPusher(user.NewHumanRefreshTokenAddedEvent(
							context.Background(),
							&user.NewAggregate("userID", "orgID").Aggregate,
							"tokenID2",
							"clientID2",
							"agentID",
							"de",
							[]string{"clientID2"},
							[]string{oidc.ScopeOpenID, oidc.ScopeProfile, oidc.ScopeEmail, oidc.ScopeOfflineAccess},
							[]string{"password"},
							time.Now(),
							1*time.Hour,
							10*time.Hour,
							nil,
						)),
					),
					expectPushFailed(zerrors.ThrowInternal(nil, "ERROR", "internal"),
						user.NewHumanRefreshTokenRemovedEvent(
							context.Background(),
							&user.NewAggregate("userID", "orgID").Aggregate,
							"tokenID",
						),
						user.NewHumanRefreshTokenRemovedEvent(
							context.Background(),
							&user.NewAggregate("userID", "orgID").Aggregate,
							"tokenID2",
						),
					),
				),
			},
			args{
				context.Background(),
				"userID",
				"orgID",
				[]string{"tokenID", "tokenID2"},
			},
			res{
				err: zerrors.IsInternal,
			},
		},
		{
			"revoke, ok",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(user.NewHumanRefreshTokenAddedEvent(
							context.Background(),
							&user.NewAggregate("userID", "orgID").Aggregate,
							"tokenID",
							"clientID",
							"agentID",
							"de",
							[]string{"clientID1"},
							[]string{oidc.ScopeOpenID, oidc.ScopeProfile, oidc.ScopeEmail, oidc.ScopeOfflineAccess},
							[]string{"password"},
							time.Now(),
							1*time.Hour,
							10*time.Hour,
							nil,
						)),
					),
					expectFilter(
						eventFromEventPusher(user.NewHumanRefreshTokenAddedEvent(
							context.Background(),
							&user.NewAggregate("userID", "orgID").Aggregate,
							"tokenID2",
							"clientID2",
							"agentID",
							"de",
							[]string{"clientID2"},
							[]string{oidc.ScopeOpenID, oidc.ScopeProfile, oidc.ScopeEmail, oidc.ScopeOfflineAccess},
							[]string{"password"},
							time.Now(),
							1*time.Hour,
							10*time.Hour,
							nil,
						)),
					),
					expectPush(
						user.NewHumanRefreshTokenRemovedEvent(
							context.Background(),
							&user.NewAggregate("userID", "orgID").Aggregate,
							"tokenID",
						),
						user.NewHumanRefreshTokenRemovedEvent(
							context.Background(),
							&user.NewAggregate("userID", "orgID").Aggregate,
							"tokenID2",
						),
					),
				),
			},
			args{
				context.Background(),
				"userID",
				"orgID",
				[]string{"tokenID", "tokenID2"},
			},
			res{
				err: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.fields.eventstore,
			}
			err := c.RevokeRefreshTokens(tt.args.ctx, tt.args.userID, tt.args.orgID, tt.args.tokenIDs)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}
