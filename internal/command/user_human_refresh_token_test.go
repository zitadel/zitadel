package command

import (
	"context"
	"encoding/base64"
	"testing"
	"time"

	"github.com/caos/oidc/pkg/oidc"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/id"
	id_mock "github.com/caos/zitadel/internal/id/mock"
	"github.com/caos/zitadel/internal/repository/user"
)

func TestCommands_AddUserAndRefreshToken(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		idGenerator  id.Generator
		iamDomain    string
		keyAlgorithm crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx                   context.Context
		orgID                 string
		agentID               string
		clientID              string
		userID                string
		refreshToken          string
		audience              []string
		scopes                []string
		authMethodsReferences []string
		lifetime              time.Duration
		authTime              time.Time
		refreshIdleExpiration time.Duration
		refreshExpiration     time.Duration
	}
	type res struct {
		token        *domain.Token
		refreshToken string
		err          func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "error access token, error",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
				),
			},
			args: args{},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "add refresh token, user inactive, error",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
				),
			},
			args: args{},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "renew refresh token, invalid token, error",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(user.NewHumanAddedEvent(
							context.Background(),
							&user.NewAggregate("userID", "orgID").Aggregate,
							"username",
							"firstname",
							"lastname",
							"nickname",
							"displayname",
							language.German,
							domain.GenderUnspecified,
							"email",
							true,
						)),
					),
					expectFilter(
						eventFromEventPusher(user.NewHumanAddedEvent(
							context.Background(),
							&user.NewAggregate("userID", "orgID").Aggregate,
							"username",
							"firstname",
							"lastname",
							"nickname",
							"displayname",
							language.German,
							domain.GenderUnspecified,
							"email",
							true,
						)),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "accessTokenID1"),
				keyAlgorithm: refreshTokenEncryptionAlgorithm(gomock.NewController(t)),
			},
			args: args{
				ctx:          context.Background(),
				refreshToken: "invalid",
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "renew refresh token, invalid token (invalid userID), error",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(user.NewHumanAddedEvent(
							context.Background(),
							&user.NewAggregate("userID", "orgID").Aggregate,
							"username",
							"firstname",
							"lastname",
							"nickname",
							"displayname",
							language.German,
							domain.GenderUnspecified,
							"email",
							true,
						)),
					),
					expectFilter(
						eventFromEventPusher(user.NewHumanAddedEvent(
							context.Background(),
							&user.NewAggregate("userID", "orgID").Aggregate,
							"username",
							"firstname",
							"lastname",
							"nickname",
							"displayname",
							language.German,
							domain.GenderUnspecified,
							"email",
							true,
						)),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "accessTokenID1"),
				keyAlgorithm: refreshTokenEncryptionAlgorithm(gomock.NewController(t)),
			},
			args: args{
				ctx:          context.Background(),
				userID:       "userID",
				orgID:        "orgID",
				refreshToken: base64.RawURLEncoding.EncodeToString([]byte("userID2:tokenID:token")),
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "renew refresh token, token inactive, error",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(user.NewHumanAddedEvent(
							context.Background(),
							&user.NewAggregate("userID", "orgID").Aggregate,
							"username",
							"firstname",
							"lastname",
							"nickname",
							"displayname",
							language.German,
							domain.GenderUnspecified,
							"email",
							true,
						)),
					),
					expectFilter(
						eventFromEventPusher(user.NewHumanAddedEvent(
							context.Background(),
							&user.NewAggregate("userID", "orgID").Aggregate,
							"username",
							"firstname",
							"lastname",
							"nickname",
							"displayname",
							language.German,
							domain.GenderUnspecified,
							"email",
							true,
						)),
					),
					expectFilter(
						eventFromEventPusher(user.NewHumanRefreshTokenAddedEvent(
							context.Background(),
							&user.NewAggregate("userID", "orgID").Aggregate,
							"tokenID",
							"applicationID",
							"userAgentID",
							"de",
							[]string{"clientID1"},
							[]string{oidc.ScopeOpenID, oidc.ScopeProfile, oidc.ScopeEmail, oidc.ScopeOfflineAccess},
							[]string{"password"},
							time.Now(),
							1*time.Hour,
							24*time.Hour,
						)),
						eventFromEventPusher(user.NewHumanRefreshTokenRemovedEvent(
							context.Background(),
							&user.NewAggregate("userID", "orgID").Aggregate,
							"tokenID",
						)),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "accessTokenID1"),
				keyAlgorithm: refreshTokenEncryptionAlgorithm(gomock.NewController(t)),
			},
			args: args{
				ctx:          context.Background(),
				userID:       "userID",
				orgID:        "orgID",
				refreshToken: base64.RawURLEncoding.EncodeToString([]byte("userID:tokenID:token")),
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "renew refresh token, token expired, error",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(user.NewHumanAddedEvent(
							context.Background(),
							&user.NewAggregate("userID", "orgID").Aggregate,
							"username",
							"firstname",
							"lastname",
							"nickname",
							"displayname",
							language.German,
							domain.GenderUnspecified,
							"email",
							true,
						)),
					),
					expectFilter(
						eventFromEventPusher(user.NewHumanAddedEvent(
							context.Background(),
							&user.NewAggregate("userID", "orgID").Aggregate,
							"username",
							"firstname",
							"lastname",
							"nickname",
							"displayname",
							language.German,
							domain.GenderUnspecified,
							"email",
							true,
						)),
					),
					expectFilter(
						eventFromEventPusher(user.NewHumanRefreshTokenAddedEvent(
							context.Background(),
							&user.NewAggregate("userID", "orgID").Aggregate,
							"tokenID",
							"applicationID",
							"userAgentID",
							"de",
							[]string{"clientID1"},
							[]string{oidc.ScopeOpenID, oidc.ScopeProfile, oidc.ScopeEmail, oidc.ScopeOfflineAccess},
							[]string{"password"},
							time.Now(),
							-1*time.Hour,
							24*time.Hour,
						)),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "accessTokenID1"),
				keyAlgorithm: refreshTokenEncryptionAlgorithm(gomock.NewController(t)),
			},
			args: args{
				ctx:          context.Background(),
				userID:       "userID",
				orgID:        "orgID",
				refreshToken: base64.RawURLEncoding.EncodeToString([]byte("userID:tokenID:tokenID")),
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		//fails because of timestamp equality
		//{
		//	name: "push failed, error",
		//	fields: fields{
		//		eventstore: eventstoreExpect(t,
		//			expectFilter(
		//				eventFromEventPusher(user.NewHumanAddedEvent(
		//					context.Background(),
		//					&user.NewAggregate("userID", "orgID").Aggregate,
		//					"username",
		//					"firstname",
		//					"lastname",
		//					"nickname",
		//					"displayname",
		//					language.German,
		//					domain.GenderUnspecified,
		//					"email",
		//					true,
		//				)),
		//			),
		//			expectFilter(
		//				eventFromEventPusherWithCreationDateNow(user.NewHumanAddedEvent(
		//					context.Background(),
		//					&user.NewAggregate("userID", "orgID").Aggregate,
		//					"username",
		//					"firstname",
		//					"lastname",
		//					"nickname",
		//					"displayname",
		//					language.German,
		//					domain.GenderUnspecified,
		//					"email",
		//					true,
		//				)),
		//			),
		//			expectFilter(
		//				eventFromEventPusherWithCreationDateNow(user.NewHumanRefreshTokenAddedEvent(
		//					context.Background(),
		//					&user.NewAggregate("userID", "orgID").Aggregate,
		//					"tokenID",
		//					"applicationID",
		//					"userAgentID",
		//					"de",
		//					[]string{"clientID1"},
		//					[]string{oidc.ScopeOpenID, oidc.ScopeProfile, oidc.ScopeEmail, oidc.ScopeOfflineAccess},
		//					[]string{"password"},
		//					time.Now(),
		//					1*time.Hour,
		//					24*time.Hour,
		//				)),
		//			),
		//			expectPushFailed(
		//				caos_errs.ThrowInternal(nil, "ERROR", "internal"),
		//				[]*repository.Event{
		//					eventFromEventPusher(user.NewUserTokenAddedEvent(
		//						context.Background(),
		//						&user.NewAggregate("userID", "orgID").Aggregate,
		//						"accessTokenID1",
		//						"clientID",
		//						"agentID",
		//						"de",
		//						[]string{"clientID1"},
		//						[]string{oidc.ScopeOpenID, oidc.ScopeProfile, oidc.ScopeEmail, oidc.ScopeOfflineAccess},
		//						time.Now().Add(5*time.Minute),
		//					)),
		//					eventFromEventPusher(user.NewHumanRefreshTokenRenewedEvent(
		//						context.Background(),
		//						&user.NewAggregate("userID", "orgID").Aggregate,
		//						"tokenID",
		//						"refreshToken1",
		//						1*time.Hour,
		//					)),
		//				},
		//			),
		//		),
		//		idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "accessTokenID1", "refreshToken1"),
		//		keyAlgorithm: refreshTokenEncryptionAlgorithm(gomock.NewController(t)),
		//	},
		//	args: args{
		//		ctx:                   context.Background(),
		//		orgID:                 "orgID",
		//		agentID:               "agentID",
		//		clientID:              "clientID",
		//		userID:                "userID",
		//		refreshToken:          base64.RawURLEncoding.EncodeToString([]byte("userID:tokenID:tokenID")),
		//		audience:              []string{"clientID1"},
		//		scopes:                []string{oidc.ScopeOpenID, oidc.ScopeProfile, oidc.ScopeEmail, oidc.ScopeOfflineAccess},
		//		authMethodsReferences: []string{"password"},
		//		lifetime:              5 * time.Minute,
		//		authTime:              time.Now(),
		//	},
		//	res: res{
		//		err: caos_errs.IsInternal,
		//	},
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:   tt.fields.eventstore,
				idGenerator:  tt.fields.idGenerator,
				iamDomain:    tt.fields.iamDomain,
				keyAlgorithm: tt.fields.keyAlgorithm,
			}
			got, gotRefresh, err := c.AddUserAndRefreshToken(tt.args.ctx, tt.args.orgID, tt.args.agentID, tt.args.clientID, tt.args.userID, tt.args.refreshToken,
				tt.args.audience, tt.args.scopes, tt.args.authMethodsReferences, tt.args.lifetime, tt.args.refreshIdleExpiration, tt.args.refreshExpiration, tt.args.authTime)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.token, got)
				assert.Equal(t, tt.res.refreshToken, gotRefresh)
			}
		})
	}
}

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
				err: caos_errs.IsErrorInvalidArgument,
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
				err: caos_errs.IsNotFound,
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
						)),
					),
					expectPushFailed(caos_errs.ThrowInternal(nil, "ERROR", "internal"),
						[]*repository.Event{
							eventFromEventPusher(user.NewHumanRefreshTokenRemovedEvent(
								context.Background(),
								&user.NewAggregate("userID", "orgID").Aggregate,
								"tokenID",
							)),
						},
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
				err: caos_errs.IsInternal,
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
						)),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(user.NewHumanRefreshTokenRemovedEvent(
								context.Background(),
								&user.NewAggregate("userID", "orgID").Aggregate,
								"tokenID",
							)),
						},
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
				err: caos_errs.IsErrorInvalidArgument,
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
				err: caos_errs.IsNotFound,
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
						)),
					),
					expectPushFailed(caos_errs.ThrowInternal(nil, "ERROR", "internal"),
						[]*repository.Event{
							eventFromEventPusher(user.NewHumanRefreshTokenRemovedEvent(
								context.Background(),
								&user.NewAggregate("userID", "orgID").Aggregate,
								"tokenID",
							)),
							eventFromEventPusher(user.NewHumanRefreshTokenRemovedEvent(
								context.Background(),
								&user.NewAggregate("userID", "orgID").Aggregate,
								"tokenID2",
							)),
						},
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
				err: caos_errs.IsInternal,
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
						)),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(user.NewHumanRefreshTokenRemovedEvent(
								context.Background(),
								&user.NewAggregate("userID", "orgID").Aggregate,
								"tokenID",
							)),
							eventFromEventPusher(user.NewHumanRefreshTokenRemovedEvent(
								context.Background(),
								&user.NewAggregate("userID", "orgID").Aggregate,
								"tokenID2",
							)),
						},
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

func refreshTokenEncryptionAlgorithm(ctrl *gomock.Controller) crypto.EncryptionAlgorithm {
	mCrypto := crypto.NewMockEncryptionAlgorithm(ctrl)
	mCrypto.EXPECT().Algorithm().AnyTimes().Return("enc")
	mCrypto.EXPECT().EncryptionKeyID().AnyTimes().Return("id")
	mCrypto.EXPECT().Encrypt(gomock.Any()).AnyTimes().DoAndReturn(
		func(refrehToken []byte) ([]byte, error) {
			return refrehToken, nil
		},
	)
	mCrypto.EXPECT().Decrypt(gomock.Any(), gomock.Any()).AnyTimes().DoAndReturn(
		func(refrehToken []byte, keyID string) ([]byte, error) {
			if keyID != "id" {
				return nil, caos_errs.ThrowInternal(nil, "id", "invalid key id")
			}
			return refrehToken, nil
		},
	)
	return mCrypto
}

func TestCommands_addRefreshToken(t *testing.T) {
	authTime := time.Now().Add(-1 * time.Hour)
	type fields struct {
		eventstore   *eventstore.Eventstore
		idGenerator  id.Generator
		keyAlgorithm crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx                   context.Context
		accessToken           *domain.Token
		authMethodsReferences []string
		authTime              time.Time
		idleExpiration        time.Duration
		expiration            time.Duration
	}
	type res struct {
		event        *user.HumanRefreshTokenAddedEvent
		refreshToken string
		err          func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{

		{
			name: "add refresh Token",
			fields: fields{
				eventstore:   eventstoreExpect(t),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "refreshTokenID"),
				keyAlgorithm: refreshTokenEncryptionAlgorithm(gomock.NewController(t)),
			},
			args: args{
				ctx: context.Background(),
				accessToken: &domain.Token{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "userID",
						ResourceOwner: "org1",
					},
					TokenID:           "accessTokenID1",
					ApplicationID:     "clientID",
					UserAgentID:       "agentID",
					Audience:          []string{"clientID1"},
					Expiration:        time.Now().Add(5 * time.Minute),
					Scopes:            []string{oidc.ScopeOpenID, oidc.ScopeProfile, oidc.ScopeEmail, oidc.ScopeOfflineAccess},
					PreferredLanguage: "de",
				},
				authMethodsReferences: []string{"password"},
				authTime:              authTime,
				idleExpiration:        1 * time.Hour,
				expiration:            10 * time.Hour,
			},
			res: res{
				event: user.NewHumanRefreshTokenAddedEvent(
					context.Background(),
					&user.NewAggregate("userID", "org1").Aggregate,
					"refreshTokenID",
					"clientID",
					"agentID",
					"de",
					[]string{"clientID1"},
					[]string{oidc.ScopeOpenID, oidc.ScopeProfile, oidc.ScopeEmail, oidc.ScopeOfflineAccess},
					[]string{"password"},
					authTime,
					1*time.Hour,
					10*time.Hour,
				),
				refreshToken: base64.RawURLEncoding.EncodeToString([]byte("userID:refreshTokenID:refreshTokenID")),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:   tt.fields.eventstore,
				idGenerator:  tt.fields.idGenerator,
				keyAlgorithm: tt.fields.keyAlgorithm,
			}
			gotEvent, gotRefreshToken, err := c.addRefreshToken(tt.args.ctx, tt.args.accessToken, tt.args.authMethodsReferences, tt.args.authTime, tt.args.idleExpiration, tt.args.expiration)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.event, gotEvent)
				assert.Equal(t, tt.res.refreshToken, gotRefreshToken)
			}
		})
	}
}

func TestCommands_renewRefreshToken(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		idGenerator  id.Generator
		keyAlgorithm crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx            context.Context
		userID         string
		orgID          string
		refreshToken   string
		idleExpiration time.Duration
	}
	type res struct {
		event           *user.HumanRefreshTokenRenewedEvent
		newRefreshToken string
		err             func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "empty token, error",
			fields: fields{
				eventstore: eventstoreExpect(t),
			},
			args: args{
				ctx: context.Background(),
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "invalid token, error",
			fields: fields{
				eventstore:   eventstoreExpect(t),
				keyAlgorithm: refreshTokenEncryptionAlgorithm(gomock.NewController(t)),
			},
			args: args{
				ctx:          context.Background(),
				refreshToken: "invalid",
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "invalid token (invalid userID), error",
			fields: fields{
				eventstore:   eventstoreExpect(t),
				keyAlgorithm: refreshTokenEncryptionAlgorithm(gomock.NewController(t)),
			},
			args: args{
				ctx:          context.Background(),
				userID:       "userID",
				orgID:        "orgID",
				refreshToken: base64.RawURLEncoding.EncodeToString([]byte("userID2:tokenID:token")),
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "token inactive, error",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(user.NewHumanRefreshTokenAddedEvent(
							context.Background(),
							&user.NewAggregate("userID", "orgID").Aggregate,
							"tokenID",
							"applicationID",
							"userAgentID",
							"de",
							[]string{"clientID1"},
							[]string{oidc.ScopeOpenID, oidc.ScopeProfile, oidc.ScopeEmail, oidc.ScopeOfflineAccess},
							[]string{"password"},
							time.Now(),
							1*time.Hour,
							24*time.Hour,
						)),
						eventFromEventPusher(user.NewHumanRefreshTokenRemovedEvent(
							context.Background(),
							&user.NewAggregate("userID", "orgID").Aggregate,
							"tokenID",
						)),
					),
				),
				keyAlgorithm: refreshTokenEncryptionAlgorithm(gomock.NewController(t)),
			},
			args: args{
				ctx:          context.Background(),
				userID:       "userID",
				orgID:        "orgID",
				refreshToken: base64.RawURLEncoding.EncodeToString([]byte("userID:tokenID:token")),
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "token expired, error",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(user.NewHumanRefreshTokenAddedEvent(
							context.Background(),
							&user.NewAggregate("userID", "orgID").Aggregate,
							"tokenID",
							"applicationID",
							"userAgentID",
							"de",
							[]string{"clientID1"},
							[]string{oidc.ScopeOpenID, oidc.ScopeProfile, oidc.ScopeEmail, oidc.ScopeOfflineAccess},
							[]string{"password"},
							time.Now(),
							1*time.Hour,
							24*time.Hour,
						)),
					),
				),
				keyAlgorithm: refreshTokenEncryptionAlgorithm(gomock.NewController(t)),
			},
			args: args{
				ctx:          context.Background(),
				userID:       "userID",
				orgID:        "orgID",
				refreshToken: base64.RawURLEncoding.EncodeToString([]byte("userID:tokenID:tokenID")),
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "token renewed, ok",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusherWithCreationDateNow(user.NewHumanRefreshTokenAddedEvent(
							context.Background(),
							&user.NewAggregate("userID", "orgID").Aggregate,
							"tokenID",
							"applicationID",
							"userAgentID",
							"de",
							[]string{"clientID1"},
							[]string{oidc.ScopeOpenID, oidc.ScopeProfile, oidc.ScopeEmail, oidc.ScopeOfflineAccess},
							[]string{"password"},
							time.Now(),
							1*time.Hour,
							24*time.Hour,
						)),
					),
				),
				keyAlgorithm: refreshTokenEncryptionAlgorithm(gomock.NewController(t)),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "refreshToken1"),
			},
			args: args{
				ctx:            context.Background(),
				userID:         "userID",
				orgID:          "orgID",
				refreshToken:   base64.RawURLEncoding.EncodeToString([]byte("userID:tokenID:tokenID")),
				idleExpiration: 1 * time.Hour,
			},
			res: res{
				event: user.NewHumanRefreshTokenRenewedEvent(
					context.Background(),
					&user.NewAggregate("userID", "orgID").Aggregate,
					"tokenID",
					"refreshToken1",
					1*time.Hour,
				),
				newRefreshToken: base64.RawURLEncoding.EncodeToString([]byte("userID:tokenID:refreshToken1")),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:   tt.fields.eventstore,
				idGenerator:  tt.fields.idGenerator,
				keyAlgorithm: tt.fields.keyAlgorithm,
			}
			gotEvent, gotNewRefreshToken, err := c.renewRefreshToken(tt.args.ctx, tt.args.userID, tt.args.orgID, tt.args.refreshToken, tt.args.idleExpiration)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.event, gotEvent)
				assert.Equal(t, tt.res.newRefreshToken, gotNewRefreshToken)
			}
		})
	}
}
