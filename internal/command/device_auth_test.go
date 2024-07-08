package command

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id_generator"
	"github.com/zitadel/zitadel/internal/id_generator/mock"
	"github.com/zitadel/zitadel/internal/repository/deviceauth"
	"github.com/zitadel/zitadel/internal/repository/oidcsession"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommands_AddDeviceAuth(t *testing.T) {
	ctx := authz.WithInstanceID(context.Background(), "instance1")
	pushErr := errors.New("pushErr")
	now := time.Now()

	unique := deviceauth.NewAddUniqueConstraints("123", "456")
	require.Len(t, unique, 2)

	type fields struct {
		eventstore func(*testing.T) *eventstore.Eventstore
	}
	type args struct {
		ctx              context.Context
		clientID         string
		deviceCode       string
		userCode         string
		expires          time.Time
		scopes           []string
		audience         []string
		needRefreshToken bool
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantDetails *domain.ObjectDetails
		wantErr     error
	}{
		{
			name: "success",
			fields: fields{
				eventstore: expectEventstore(expectPush(
					deviceauth.NewAddedEvent(
						ctx,
						deviceauth.NewAggregate("123", "instance1"),
						"client_id", "123", "456", now,
						[]string{"a", "b", "c"},
						[]string{"projectID", "clientID"}, true,
					),
				)),
			},
			args: args{
				ctx:              authz.WithInstanceID(context.Background(), "instance1"),
				clientID:         "client_id",
				deviceCode:       "123",
				userCode:         "456",
				expires:          now,
				scopes:           []string{"a", "b", "c"},
				audience:         []string{"projectID", "clientID"},
				needRefreshToken: true,
			},
			wantDetails: &domain.ObjectDetails{
				ResourceOwner: "instance1",
			},
		},
		{
			name: "push error",
			fields: fields{
				eventstore: expectEventstore(expectPushFailed(pushErr,
					deviceauth.NewAddedEvent(
						ctx,
						deviceauth.NewAggregate("123", "instance1"),
						"client_id", "123", "456", now,
						[]string{"a", "b", "c"},
						[]string{"projectID", "clientID"}, false,
					)),
				),
			},
			args: args{
				ctx:              authz.WithInstanceID(context.Background(), "instance1"),
				clientID:         "client_id",
				deviceCode:       "123",
				userCode:         "456",
				expires:          now,
				scopes:           []string{"a", "b", "c"},
				audience:         []string{"projectID", "clientID"},
				needRefreshToken: false,
			},
			wantErr: pushErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.fields.eventstore(t),
			}
			gotDetails, err := c.AddDeviceAuth(tt.args.ctx, tt.args.clientID, tt.args.deviceCode, tt.args.userCode, tt.args.expires, tt.args.scopes, tt.args.audience, tt.args.needRefreshToken)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.wantDetails, gotDetails)
		})
	}
}

func TestCommands_ApproveDeviceAuth(t *testing.T) {
	ctx := authz.WithInstanceID(context.Background(), "instance1")
	now := time.Now()
	pushErr := errors.New("pushErr")

	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx               context.Context
		id                string
		userID            string
		userOrgID         string
		authMethods       []domain.UserAuthMethodType
		authTime          time.Time
		preferredLanguage *language.Tag
		userAgent         *domain.UserAgent
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantDetails *domain.ObjectDetails
		wantErr     error
	}{
		{
			name: "not found error",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
				),
			},
			args: args{
				ctx, "123", "subj", "orgID",
				[]domain.UserAuthMethodType{domain.UserAuthMethodTypePassword},
				time.Unix(123, 456), &language.Afrikaans, &domain.UserAgent{
					FingerprintID: gu.Ptr("fp1"),
					IP:            net.ParseIP("1.2.3.4"),
					Description:   gu.Ptr("firefox"),
					Header:        http.Header{"foo": []string{"bar"}},
				},
			},
			wantErr: zerrors.ThrowNotFound(nil, "COMMAND-Hief9", "Errors.DeviceAuth.NotFound"),
		},
		{
			name: "push error",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(eventFromEventPusherWithInstanceID(
						"instance1",
						deviceauth.NewAddedEvent(
							ctx,
							deviceauth.NewAggregate("123", "instance1"),
							"client_id", "123", "456", now,
							[]string{"a", "b", "c"},
							[]string{"projectID", "clientID"}, true,
						),
					)),
					expectPushFailed(pushErr,
						deviceauth.NewApprovedEvent(
							ctx, deviceauth.NewAggregate("123", "instance1"), "subj", "orgID",
							[]domain.UserAuthMethodType{domain.UserAuthMethodTypePassword},
							time.Unix(123, 456), &language.Afrikaans, &domain.UserAgent{
								FingerprintID: gu.Ptr("fp1"),
								IP:            net.ParseIP("1.2.3.4"),
								Description:   gu.Ptr("firefox"),
								Header:        http.Header{"foo": []string{"bar"}},
							},
						),
					),
				),
			},
			args: args{
				ctx, "123", "subj", "orgID",
				[]domain.UserAuthMethodType{domain.UserAuthMethodTypePassword},
				time.Unix(123, 456), &language.Afrikaans, &domain.UserAgent{
					FingerprintID: gu.Ptr("fp1"),
					IP:            net.ParseIP("1.2.3.4"),
					Description:   gu.Ptr("firefox"),
					Header:        http.Header{"foo": []string{"bar"}},
				},
			},
			wantErr: pushErr,
		},
		{
			name: "success",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(eventFromEventPusherWithInstanceID(
						"instance1",
						deviceauth.NewAddedEvent(
							ctx,
							deviceauth.NewAggregate("123", "instance1"),
							"client_id", "123", "456", now,
							[]string{"a", "b", "c"},
							[]string{"projectID", "clientID"}, true,
						),
					)),
					expectPush(
						deviceauth.NewApprovedEvent(
							ctx, deviceauth.NewAggregate("123", "instance1"), "subj", "orgID",
							[]domain.UserAuthMethodType{domain.UserAuthMethodTypePassword},
							time.Unix(123, 456), &language.Afrikaans, &domain.UserAgent{
								FingerprintID: gu.Ptr("fp1"),
								IP:            net.ParseIP("1.2.3.4"),
								Description:   gu.Ptr("firefox"),
								Header:        http.Header{"foo": []string{"bar"}},
							},
						),
					),
				),
			},
			args: args{
				ctx, "123", "subj", "orgID",
				[]domain.UserAuthMethodType{domain.UserAuthMethodTypePassword},
				time.Unix(123, 456), &language.Afrikaans, &domain.UserAgent{
					FingerprintID: gu.Ptr("fp1"),
					IP:            net.ParseIP("1.2.3.4"),
					Description:   gu.Ptr("firefox"),
					Header:        http.Header{"foo": []string{"bar"}},
				},
			},
			wantDetails: &domain.ObjectDetails{
				ResourceOwner: "instance1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.fields.eventstore,
			}
			gotDetails, err := c.ApproveDeviceAuth(tt.args.ctx, tt.args.id, tt.args.userID, tt.args.userOrgID, tt.args.authMethods, tt.args.authTime, tt.args.preferredLanguage, tt.args.userAgent)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, gotDetails, tt.wantDetails)
		})
	}
}

func TestCommands_CancelDeviceAuth(t *testing.T) {
	ctx := authz.WithInstanceID(context.Background(), "instance1")
	now := time.Now()
	pushErr := errors.New("pushErr")

	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx    context.Context
		id     string
		reason domain.DeviceAuthCanceled
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantDetails *domain.ObjectDetails
		wantErr     error
	}{
		{
			name: "not found error",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
				),
			},
			args:    args{ctx, "123", domain.DeviceAuthCanceledDenied},
			wantErr: zerrors.ThrowNotFound(nil, "COMMAND-gee5A", "Errors.DeviceAuth.NotFound"),
		},
		{
			name: "push error",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(eventFromEventPusherWithInstanceID(
						"instance1",
						deviceauth.NewAddedEvent(
							ctx,
							deviceauth.NewAggregate("123", "instance1"),
							"client_id", "123", "456", now,
							[]string{"a", "b", "c"},
							[]string{"projectID", "clientID"}, true,
						),
					)),
					expectPushFailed(pushErr,
						deviceauth.NewCanceledEvent(
							ctx, deviceauth.NewAggregate("123", "instance1"),
							domain.DeviceAuthCanceledDenied,
						),
					),
				),
			},
			args:    args{ctx, "123", domain.DeviceAuthCanceledDenied},
			wantErr: pushErr,
		},
		{
			name: "success/denied",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(eventFromEventPusherWithInstanceID(
						"instance1",
						deviceauth.NewAddedEvent(
							ctx,
							deviceauth.NewAggregate("123", "instance1"),
							"client_id", "123", "456", now,
							[]string{"a", "b", "c"},
							[]string{"projectID", "clientID"}, true,
						),
					)),
					expectPush(
						deviceauth.NewCanceledEvent(
							ctx, deviceauth.NewAggregate("123", "instance1"),
							domain.DeviceAuthCanceledDenied,
						),
					),
				),
			},
			args: args{ctx, "123", domain.DeviceAuthCanceledDenied},
			wantDetails: &domain.ObjectDetails{
				ResourceOwner: "instance1",
			},
		},
		{
			name: "success/expired",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(eventFromEventPusherWithInstanceID(
						"instance1",
						deviceauth.NewAddedEvent(
							ctx,
							deviceauth.NewAggregate("123", "instance1"),
							"client_id", "123", "456", now,
							[]string{"a", "b", "c"},
							[]string{"projectID", "clientID"}, true,
						),
					)),
					expectPush(
						deviceauth.NewCanceledEvent(
							ctx, deviceauth.NewAggregate("123", "instance1"),
							domain.DeviceAuthCanceledExpired,
						),
					),
				),
			},
			args: args{ctx, "123", domain.DeviceAuthCanceledExpired},
			wantDetails: &domain.ObjectDetails{
				ResourceOwner: "instance1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.fields.eventstore,
			}
			gotDetails, err := c.CancelDeviceAuth(tt.args.ctx, tt.args.id, tt.args.reason)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, gotDetails, tt.wantDetails)
		})
	}
}

func TestCommands_CreateOIDCSessionFromDeviceAuth(t *testing.T) {
	ctx := authz.WithInstanceID(context.Background(), "instance1")

	type fields struct {
		eventstore                      func(*testing.T) *eventstore.Eventstore
		idGenerator                     id_generator.Generator
		defaultAccessTokenLifetime      time.Duration
		defaultRefreshTokenLifetime     time.Duration
		defaultRefreshTokenIdleLifetime time.Duration
		keyAlgorithm                    crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx        context.Context
		deviceCode string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *OIDCSession
		wantErr error
	}{
		{
			name: "device auth filter error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilterError(io.ErrClosedPipe),
				),
			},
			args: args{
				ctx,
				"device1",
			},
			wantErr: io.ErrClosedPipe,
		},
		{
			name: "not yet approved",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusherWithInstanceID(
							"instance1",
							deviceauth.NewAddedEvent(
								ctx,
								deviceauth.NewAggregate("123", "instance1"),
								"clientID", "123", "456", time.Now().Add(time.Minute),
								[]string{"openid", "offline_access"},
								[]string{"audience"}, false,
							),
						),
					),
				),
			},
			args: args{
				ctx,
				"123",
			},
			wantErr: DeviceAuthStateError(domain.DeviceAuthStateInitiated),
		},
		{
			name: "not found",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx,
				"123",
			},
			wantErr: zerrors.ThrowNotFound(nil, "COMMAND-ua1Vo", "Errors.DeviceAuth.NotFound"),
		},
		{
			name: "expired",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusherWithInstanceID(
							"instance1",
							deviceauth.NewAddedEvent(
								ctx,
								deviceauth.NewAggregate("123", "instance1"),
								"clientID", "123", "456", time.Now().Add(-time.Minute),
								[]string{"openid", "offline_access"},
								[]string{"audience"}, false,
							),
						),
					),
					expectPushSlow(time.Second, deviceauth.NewCanceledEvent(ctx,
						deviceauth.NewAggregate("123", "instance1"),
						domain.DeviceAuthCanceledExpired,
					)),
				),
			},
			args: args{
				ctx,
				"123",
			},
			wantErr: DeviceAuthStateError(domain.DeviceAuthStateExpired),
		},
		{
			name: "already expired",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusherWithInstanceID(
							"instance1",
							deviceauth.NewAddedEvent(
								ctx,
								deviceauth.NewAggregate("123", "instance1"),
								"clientID", "123", "456", time.Now().Add(-time.Minute),
								[]string{"openid", "offline_access"},
								[]string{"audience"}, false,
							),
						),
						eventFromEventPusherWithInstanceID(
							"instance1",
							deviceauth.NewCanceledEvent(ctx,
								deviceauth.NewAggregate("123", "instance1"),
								domain.DeviceAuthCanceledExpired,
							),
						),
					),
				),
			},
			args: args{
				ctx,
				"123",
			},
			wantErr: DeviceAuthStateError(domain.DeviceAuthStateExpired),
		},
		{
			name: "denied",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusherWithInstanceID(
							"instance1",
							deviceauth.NewAddedEvent(
								ctx,
								deviceauth.NewAggregate("123", "instance1"),
								"clientID", "123", "456", time.Now().Add(-time.Minute),
								[]string{"openid", "offline_access"},
								[]string{"audience"}, false,
							),
						),
						eventFromEventPusherWithInstanceID(
							"instance1",
							deviceauth.NewCanceledEvent(ctx,
								deviceauth.NewAggregate("123", "instance1"),
								domain.DeviceAuthCanceledDenied,
							),
						),
					),
				),
			},
			args: args{
				ctx,
				"123",
			},
			wantErr: DeviceAuthStateError(domain.DeviceAuthStateDenied),
		},
		{
			name: "already done",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusherWithInstanceID(
							"instance1",
							deviceauth.NewAddedEvent(
								ctx,
								deviceauth.NewAggregate("123", "instance1"),
								"clientID", "123", "456", time.Now().Add(-time.Minute),
								[]string{"openid", "offline_access"},
								[]string{"audience"}, false,
							),
						),
						eventFromEventPusherWithInstanceID(
							"instance1",
							deviceauth.NewCanceledEvent(ctx,
								deviceauth.NewAggregate("123", "instance1"),
								domain.DeviceAuthCanceledDenied,
							),
						),
						eventFromEventPusherWithInstanceID(
							"instance1",
							deviceauth.NewDoneEvent(ctx,
								deviceauth.NewAggregate("123", "instance1"),
							),
						),
					),
				),
			},
			args: args{
				ctx,
				"123",
			},
			wantErr: DeviceAuthStateError(domain.DeviceAuthStateDone),
		},
		{
			name: "approved, success",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusherWithInstanceID(
							"instance1",
							deviceauth.NewAddedEvent(
								ctx,
								deviceauth.NewAggregate("123", "instance1"),
								"clientID", "123", "456", time.Now().Add(-time.Minute),
								[]string{"openid", "offline_access"},
								[]string{"audience"}, false,
							),
						),
						eventFromEventPusherWithInstanceID(
							"instance1",
							deviceauth.NewApprovedEvent(ctx,
								deviceauth.NewAggregate("123", "instance1"),
								"userID", "org1",
								[]domain.UserAuthMethodType{domain.UserAuthMethodTypePassword},
								testNow, &language.Afrikaans, &domain.UserAgent{
									FingerprintID: gu.Ptr("fp1"),
									IP:            net.ParseIP("1.2.3.4"),
									Description:   gu.Ptr("firefox"),
									Header:        http.Header{"foo": []string{"bar"}},
								},
							),
						),
					),
					expectFilter(), // token lifetime
					expectPush(
						oidcsession.NewAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
							"userID", "org1", "", "clientID", []string{"audience"}, []string{"openid", "offline_access"},
							[]domain.UserAuthMethodType{domain.UserAuthMethodTypePassword}, testNow, "", &language.Afrikaans, &domain.UserAgent{
								FingerprintID: gu.Ptr("fp1"),
								IP:            net.ParseIP("1.2.3.4"),
								Description:   gu.Ptr("firefox"),
								Header:        http.Header{"foo": []string{"bar"}},
							},
						),
						oidcsession.NewAccessTokenAddedEvent(context.Background(),
							&oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
							"at_accessTokenID", []string{"openid", "offline_access"}, time.Hour, domain.TokenReasonAuthRequest, nil,
						),
						user.NewUserTokenV2AddedEvent(context.Background(), &user.NewAggregate("userID", "org1").Aggregate, "at_accessTokenID"),
						deviceauth.NewDoneEvent(ctx,
							deviceauth.NewAggregate("123", "instance1"),
						),
					),
				),
				idGenerator:                     mock.NewIDGeneratorExpectIDs(t, "oidcSessionID", "accessTokenID"),
				defaultAccessTokenLifetime:      time.Hour,
				defaultRefreshTokenLifetime:     7 * 24 * time.Hour,
				defaultRefreshTokenIdleLifetime: 24 * time.Hour,
				keyAlgorithm:                    crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx,
				"123",
			},
			want: &OIDCSession{
				TokenID:           "V2_oidcSessionID.at_accessTokenID",
				ClientID:          "clientID",
				UserID:            "userID",
				Audience:          []string{"audience"},
				Expiration:        time.Time{}.Add(time.Hour),
				Scope:             []string{"openid", "offline_access"},
				AuthMethods:       []domain.UserAuthMethodType{domain.UserAuthMethodTypePassword},
				AuthTime:          testNow,
				PreferredLanguage: &language.Afrikaans,
				UserAgent: &domain.UserAgent{
					FingerprintID: gu.Ptr("fp1"),
					IP:            net.ParseIP("1.2.3.4"),
					Description:   gu.Ptr("firefox"),
					Header:        http.Header{"foo": []string{"bar"}},
				},
				Reason: domain.TokenReasonAuthRequest,
			},
		},
		{
			name: "approved, with refresh token",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusherWithInstanceID(
							"instance1",
							deviceauth.NewAddedEvent(
								ctx,
								deviceauth.NewAggregate("123", "instance1"),
								"clientID", "123", "456", time.Now().Add(-time.Minute),
								[]string{"openid", "offline_access"},
								[]string{"audience"}, true,
							),
						),
						eventFromEventPusherWithInstanceID(
							"instance1",
							deviceauth.NewApprovedEvent(ctx,
								deviceauth.NewAggregate("123", "instance1"),
								"userID", "org1",
								[]domain.UserAuthMethodType{domain.UserAuthMethodTypePassword},
								testNow, &language.Afrikaans, &domain.UserAgent{
									FingerprintID: gu.Ptr("fp1"),
									IP:            net.ParseIP("1.2.3.4"),
									Description:   gu.Ptr("firefox"),
									Header:        http.Header{"foo": []string{"bar"}},
								},
							),
						),
					),
					expectFilter(), // token lifetime
					expectPush(
						oidcsession.NewAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
							"userID", "org1", "", "clientID", []string{"audience"}, []string{"openid", "offline_access"},
							[]domain.UserAuthMethodType{domain.UserAuthMethodTypePassword}, testNow, "", &language.Afrikaans, &domain.UserAgent{
								FingerprintID: gu.Ptr("fp1"),
								IP:            net.ParseIP("1.2.3.4"),
								Description:   gu.Ptr("firefox"),
								Header:        http.Header{"foo": []string{"bar"}},
							},
						),
						oidcsession.NewAccessTokenAddedEvent(context.Background(),
							&oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
							"at_accessTokenID", []string{"openid", "offline_access"}, time.Hour, domain.TokenReasonAuthRequest, nil,
						),
						user.NewUserTokenV2AddedEvent(context.Background(), &user.NewAggregate("userID", "org1").Aggregate, "at_accessTokenID"),
						oidcsession.NewRefreshTokenAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
							"rt_refreshTokenID", 7*24*time.Hour, 24*time.Hour,
						),
						deviceauth.NewDoneEvent(ctx,
							deviceauth.NewAggregate("123", "instance1"),
						),
					),
				),
				idGenerator:                     mock.NewIDGeneratorExpectIDs(t, "oidcSessionID", "accessTokenID", "refreshTokenID"),
				defaultAccessTokenLifetime:      time.Hour,
				defaultRefreshTokenLifetime:     7 * 24 * time.Hour,
				defaultRefreshTokenIdleLifetime: 24 * time.Hour,
				keyAlgorithm:                    crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx,
				"123",
			},
			want: &OIDCSession{
				TokenID:           "V2_oidcSessionID.at_accessTokenID",
				ClientID:          "clientID",
				UserID:            "userID",
				Audience:          []string{"audience"},
				Expiration:        time.Time{}.Add(time.Hour),
				Scope:             []string{"openid", "offline_access"},
				AuthMethods:       []domain.UserAuthMethodType{domain.UserAuthMethodTypePassword},
				AuthTime:          testNow,
				PreferredLanguage: &language.Afrikaans,
				UserAgent: &domain.UserAgent{
					FingerprintID: gu.Ptr("fp1"),
					IP:            net.ParseIP("1.2.3.4"),
					Description:   gu.Ptr("firefox"),
					Header:        http.Header{"foo": []string{"bar"}},
				},
				Reason:       domain.TokenReasonAuthRequest,
				RefreshToken: "VjJfb2lkY1Nlc3Npb25JRC5ydF9yZWZyZXNoVG9rZW5JRDp1c2VySUQ", //V2_oidcSessionID.rt_refreshTokenID:userID
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:                      tt.fields.eventstore(t),
				defaultAccessTokenLifetime:      tt.fields.defaultAccessTokenLifetime,
				defaultRefreshTokenLifetime:     tt.fields.defaultRefreshTokenLifetime,
				defaultRefreshTokenIdleLifetime: tt.fields.defaultRefreshTokenIdleLifetime,
				keyAlgorithm:                    tt.fields.keyAlgorithm,
			}
			id_generator.SetGenerator(tt.fields.idGenerator)
			got, err := c.CreateOIDCSessionFromDeviceAuth(tt.args.ctx, tt.args.deviceCode)
			c.jobs.Wait()

			require.ErrorIs(t, err, tt.wantErr)

			if got != nil {
				assert.WithinRange(t, got.AuthTime, tt.want.AuthTime.Add(-time.Second), tt.want.AuthTime.Add(time.Second))
				got.AuthTime = time.Time{}
				tt.want.AuthTime = time.Time{}
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
