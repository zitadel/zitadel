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
	"github.com/zitadel/zitadel/internal/feature"
	"github.com/zitadel/zitadel/internal/id"
	"github.com/zitadel/zitadel/internal/id/mock"
	"github.com/zitadel/zitadel/internal/repository/deviceauth"
	"github.com/zitadel/zitadel/internal/repository/oidcsession"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/repository/sessionlogout"
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
			assertObjectDetails(t, tt.wantDetails, gotDetails)
		})
	}
}

func TestCommands_ApproveDeviceAuth(t *testing.T) {
	ctx := authz.WithInstanceID(context.Background(), "instance1")
	now := time.Now()
	pushErr := errors.New("pushErr")

	type fields struct {
		eventstore func(*testing.T) *eventstore.Eventstore
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
		sessionID         string
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
				eventstore: expectEventstore(
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
				"sessionID",
			},
			wantErr: zerrors.ThrowNotFound(nil, "COMMAND-Hief9", "Errors.DeviceAuth.NotFound"),
		},
		{
			name: "push error",
			fields: fields{
				eventstore: expectEventstore(
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
							"sessionID",
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
				"sessionID",
			},
			wantErr: pushErr,
		},
		{
			name: "success",
			fields: fields{
				eventstore: expectEventstore(
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
							"sessionID",
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
				"sessionID",
			},
			wantDetails: &domain.ObjectDetails{
				ResourceOwner: "instance1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.fields.eventstore(t),
			}
			gotDetails, err := c.ApproveDeviceAuth(tt.args.ctx, tt.args.id, tt.args.userID, tt.args.userOrgID, tt.args.authMethods, tt.args.authTime, tt.args.preferredLanguage, tt.args.userAgent, tt.args.sessionID)
			require.ErrorIs(t, err, tt.wantErr)
			assertObjectDetails(t, tt.wantDetails, gotDetails)
		})
	}
}

func TestCommands_ApproveDeviceAuthFromSession(t *testing.T) {
	ctx := authz.WithInstanceID(context.Background(), "instance1")
	now := time.Now()
	pushErr := errors.New("pushErr")

	type fields struct {
		eventstore      func(*testing.T) *eventstore.Eventstore
		tokenVerifier   func(ctx context.Context, sessionToken, sessionID, tokenID string) (err error)
		checkPermission domain.PermissionCheck
	}
	type args struct {
		ctx          context.Context
		deviceCode   string
		sessionID    string
		sessionToken string
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
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx,
				"notfound",
				"sessionID",
				"sessionToken",
			},
			wantErr: zerrors.ThrowNotFound(nil, "COMMAND-D2hf2", "Errors.DeviceAuth.NotFound"),
		},
		{
			name: "not initialized, error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusherWithInstanceID(
							"instance1",
							deviceauth.NewAddedEvent(
								ctx,
								deviceauth.NewAggregate("deviceCode", "instance1"),
								"client_id", "deviceCode", "456", now,
								[]string{"a", "b", "c"},
								[]string{"projectID", "clientID"}, true,
							),
						),
						eventFromEventPusherWithInstanceID(
							"instance1",
							deviceauth.NewCanceledEvent(
								ctx,
								deviceauth.NewAggregate("deviceCode", "instance1"),
								domain.DeviceAuthCanceledDenied,
							)),
					),
				),
			},
			args: args{
				ctx,
				"deviceCode",
				"sessionID",
				"sessionToken",
			},
			wantErr: zerrors.ThrowPreconditionFailed(nil, "COMMAND-D30Jf", "Errors.DeviceAuth.AlreadyHandled"),
		},
		{
			name: "missing permission, error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(eventFromEventPusherWithInstanceID(
						"instance1",
						deviceauth.NewAddedEvent(
							ctx,
							deviceauth.NewAggregate("deviceCode", "instance1"),
							"client_id", "deviceCode", "456", now,
							[]string{"a", "b", "c"},
							[]string{"projectID", "clientID"}, true,
						),
					)),
				),
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args: args{
				ctx,
				"deviceCode",
				"sessionID",
				"sessionToken",
			},
			wantErr: zerrors.ThrowPermissionDenied(nil, "AUTHZ-HKJD33", "Errors.PermissionDenied"),
		},
		{
			name: "session not active, error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(eventFromEventPusherWithInstanceID(
						"instance1",
						deviceauth.NewAddedEvent(
							ctx,
							deviceauth.NewAggregate("deviceCode", "instance1"),
							"client_id", "deviceCode", "456", now,
							[]string{"a", "b", "c"},
							[]string{"projectID", "clientID"}, true,
						),
					)),
					expectFilter(),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx,
				"deviceCode",
				"sessionID",
				"sessionToken",
			},
			wantErr: zerrors.ThrowPreconditionFailed(nil, "COMMAND-Flk38", "Errors.Session.NotExisting"),
		},
		{
			name: "invalid session token, error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(eventFromEventPusherWithInstanceID(
						"instance1",
						deviceauth.NewAddedEvent(
							ctx,
							deviceauth.NewAggregate("deviceCode", "instance1"),
							"client_id", "deviceCode", "456", now,
							[]string{"a", "b", "c"},
							[]string{"projectID", "clientID"}, true,
						),
					)),
					expectFilter(eventFromEventPusherWithInstanceID(
						"instance1",
						session.NewAddedEvent(ctx,
							&session.NewAggregate("sessionID", "instance1").Aggregate,
							&domain.UserAgent{
								FingerprintID: gu.Ptr("fp1"),
								IP:            net.ParseIP("1.2.3.4"),
								Description:   gu.Ptr("firefox"),
								Header:        http.Header{"foo": []string{"bar"}},
							},
						)),
					)),
				tokenVerifier:   newMockTokenVerifierInvalid(),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx,
				"deviceCode",
				"sessionID",
				"invalidToken",
			},
			wantErr: zerrors.ThrowPermissionDenied(nil, "COMMAND-sGr42", "Errors.Session.Token.Invalid"),
		},
		{
			name: "push error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(eventFromEventPusherWithInstanceID(
						"instance1",
						deviceauth.NewAddedEvent(
							ctx,
							deviceauth.NewAggregate("deviceCode", "instance1"),
							"client_id", "deviceCode", "456", now,
							[]string{"a", "b", "c"},
							[]string{"projectID", "clientID"}, true,
						),
					)),
					expectFilter(
						eventFromEventPusher(
							session.NewAddedEvent(ctx,
								&session.NewAggregate("sessionID", "instance1").Aggregate,
								&domain.UserAgent{
									FingerprintID: gu.Ptr("fp1"),
									IP:            net.ParseIP("1.2.3.4"),
									Description:   gu.Ptr("firefox"),
									Header:        http.Header{"foo": []string{"bar"}},
								},
							),
						),
						eventFromEventPusher(
							session.NewUserCheckedEvent(ctx, &session.NewAggregate("sessionID", "instance1").Aggregate,
								"userID", "orgID", testNow, &language.Afrikaans),
						),
						eventFromEventPusher(
							session.NewPasswordCheckedEvent(ctx, &session.NewAggregate("sessionID", "instance1").Aggregate,
								testNow),
						),
						eventFromEventPusherWithCreationDateNow(
							session.NewLifetimeSetEvent(ctx, &session.NewAggregate("sessionID", "instance1").Aggregate,
								2*time.Minute),
						),
					),
					expectPushFailed(pushErr,
						deviceauth.NewApprovedEvent(
							ctx, deviceauth.NewAggregate("deviceCode", "instance1"), "userID", "orgID",
							[]domain.UserAuthMethodType{domain.UserAuthMethodTypePassword},
							testNow, &language.Afrikaans, &domain.UserAgent{
								FingerprintID: gu.Ptr("fp1"),
								IP:            net.ParseIP("1.2.3.4"),
								Description:   gu.Ptr("firefox"),
								Header:        http.Header{"foo": []string{"bar"}},
							},
							"sessionID",
						),
					),
				),
				tokenVerifier:   newMockTokenVerifierValid(),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx,
				"deviceCode",
				"sessionID",
				"sessionToken",
			},
			wantErr: pushErr,
		},
		{
			name: "authorized",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(eventFromEventPusherWithInstanceID(
						"instance1",
						deviceauth.NewAddedEvent(
							ctx,
							deviceauth.NewAggregate("deviceCode", "instance1"),
							"client_id", "deviceCode", "456", now,
							[]string{"a", "b", "c"},
							[]string{"projectID", "clientID"}, true,
						),
					)),
					expectFilter(
						eventFromEventPusher(
							session.NewAddedEvent(ctx,
								&session.NewAggregate("sessionID", "instance1").Aggregate,
								&domain.UserAgent{
									FingerprintID: gu.Ptr("fp1"),
									IP:            net.ParseIP("1.2.3.4"),
									Description:   gu.Ptr("firefox"),
									Header:        http.Header{"foo": []string{"bar"}},
								},
							),
						),
						eventFromEventPusher(
							session.NewUserCheckedEvent(ctx, &session.NewAggregate("sessionID", "instance1").Aggregate,
								"userID", "orgID", testNow, &language.Afrikaans),
						),
						eventFromEventPusher(
							session.NewPasswordCheckedEvent(ctx, &session.NewAggregate("sessionID", "instance1").Aggregate,
								testNow),
						),
						eventFromEventPusherWithCreationDateNow(
							session.NewLifetimeSetEvent(ctx, &session.NewAggregate("sessionID", "instance1").Aggregate,
								2*time.Minute),
						),
					),
					expectPush(
						deviceauth.NewApprovedEvent(
							ctx, deviceauth.NewAggregate("deviceCode", "instance1"), "userID", "orgID",
							[]domain.UserAuthMethodType{domain.UserAuthMethodTypePassword},
							testNow, &language.Afrikaans, &domain.UserAgent{
								FingerprintID: gu.Ptr("fp1"),
								IP:            net.ParseIP("1.2.3.4"),
								Description:   gu.Ptr("firefox"),
								Header:        http.Header{"foo": []string{"bar"}},
							},
							"sessionID",
						),
					),
				),
				tokenVerifier:   newMockTokenVerifierValid(),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx,
				"deviceCode",
				"sessionID",
				"sessionToken",
			},
			wantDetails: &domain.ObjectDetails{
				ResourceOwner: "instance1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:           tt.fields.eventstore(t),
				sessionTokenVerifier: tt.fields.tokenVerifier,
				checkPermission:      tt.fields.checkPermission,
			}
			gotDetails, err := c.ApproveDeviceAuthWithSession(tt.args.ctx, tt.args.deviceCode, tt.args.sessionID, tt.args.sessionToken)
			require.ErrorIs(t, err, tt.wantErr)
			assertObjectDetails(t, tt.wantDetails, gotDetails)
		})
	}
}

func TestCommands_CancelDeviceAuth(t *testing.T) {
	ctx := authz.WithInstanceID(context.Background(), "instance1")
	now := time.Now()
	pushErr := errors.New("pushErr")

	type fields struct {
		eventstore      func(*testing.T) *eventstore.Eventstore
		checkPermission domain.PermissionCheck
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
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args:    args{ctx, "123", domain.DeviceAuthCanceledDenied},
			wantErr: zerrors.ThrowNotFound(nil, "COMMAND-gee5A", "Errors.DeviceAuth.NotFound"),
		},
		{
			name: "missing permission, error",
			fields: fields{
				eventstore: expectEventstore(
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
				),
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args:    args{ctx, "123", domain.DeviceAuthCanceledDenied},
			wantErr: zerrors.ThrowPermissionDenied(nil, "AUTHZ-HKJD33", "Errors.PermissionDenied"),
		},
		{
			name: "push error",
			fields: fields{
				eventstore: expectEventstore(
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
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args:    args{ctx, "123", domain.DeviceAuthCanceledDenied},
			wantErr: pushErr,
		},
		{
			name: "success/denied",
			fields: fields{
				eventstore: expectEventstore(
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
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{ctx, "123", domain.DeviceAuthCanceledDenied},
			wantDetails: &domain.ObjectDetails{
				ResourceOwner: "instance1",
			},
		},
		{
			name: "success/expired",
			fields: fields{
				eventstore: expectEventstore(
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
				checkPermission: newMockPermissionCheckAllowed(),
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
				eventstore:      tt.fields.eventstore(t),
				checkPermission: tt.fields.checkPermission,
			}
			gotDetails, err := c.CancelDeviceAuth(tt.args.ctx, tt.args.id, tt.args.reason)
			require.ErrorIs(t, err, tt.wantErr)
			assertObjectDetails(t, tt.wantDetails, gotDetails)
		})
	}
}

func TestCommands_CreateOIDCSessionFromDeviceAuth(t *testing.T) {
	ctx := authz.WithInstanceID(context.Background(), "instance1")

	type fields struct {
		eventstore                      func(*testing.T) *eventstore.Eventstore
		idGenerator                     id.Generator
		defaultAccessTokenLifetime      time.Duration
		defaultRefreshTokenLifetime     time.Duration
		defaultRefreshTokenIdleLifetime time.Duration
		keyAlgorithm                    crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx                  context.Context
		deviceCode           string
		backChannelLogoutURI string
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
				"",
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
				"",
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
				"",
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
				"",
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
				"",
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
				"",
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
				"",
			},
			wantErr: DeviceAuthStateError(domain.DeviceAuthStateDone),
		},
		{
			name: "user not active",
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
								"sessionID",
							),
						),
					),
					expectFilter(
						user.NewHumanAddedEvent(
							ctx,
							&user.NewAggregate("userID", "org1").Aggregate,
							"username",
							"firstname",
							"lastname",
							"nickname",
							"displayname",
							language.English,
							domain.GenderUnspecified,
							"email",
							false,
						),
						user.NewUserDeactivatedEvent(
							ctx,
							&user.NewAggregate("userID", "org1").Aggregate,
						),
					),
				),
				idGenerator:                     mock.NewIDGeneratorExpectIDs(t),
				defaultAccessTokenLifetime:      time.Hour,
				defaultRefreshTokenLifetime:     7 * 24 * time.Hour,
				defaultRefreshTokenIdleLifetime: 24 * time.Hour,
				keyAlgorithm:                    crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx,
				"123",
				"",
			},
			wantErr: zerrors.ThrowPreconditionFailed(nil, "OIDCS-kj3g2", "Errors.User.NotActive"),
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
								"sessionID",
							),
						),
					),
					expectFilter(
						user.NewHumanAddedEvent(
							ctx,
							&user.NewAggregate("userID", "org1").Aggregate,
							"username",
							"firstname",
							"lastname",
							"nickname",
							"displayname",
							language.English,
							domain.GenderUnspecified,
							"email",
							false,
						),
					),
					expectFilter(), // token lifetime
					expectPush(
						oidcsession.NewAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
							"userID", "org1", "sessionID", "clientID", []string{"audience"}, []string{"openid", "offline_access"},
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
				"",
			},
			want: &OIDCSession{
				TokenID:           "V2_oidcSessionID-at_accessTokenID",
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
				Reason:    domain.TokenReasonAuthRequest,
				SessionID: "sessionID",
			},
		},
		{
			name: "approved with backChannelLogout (feature enabled), success",
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
								"sessionID",
							),
						),
					),
					expectFilter(
						user.NewHumanAddedEvent(
							ctx,
							&user.NewAggregate("userID", "org1").Aggregate,
							"username",
							"firstname",
							"lastname",
							"nickname",
							"displayname",
							language.English,
							domain.GenderUnspecified,
							"email",
							false,
						),
					),
					expectFilter(), // token lifetime
					expectPush(
						oidcsession.NewAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
							"userID", "org1", "sessionID", "clientID", []string{"audience"}, []string{"openid", "offline_access"},
							[]domain.UserAuthMethodType{domain.UserAuthMethodTypePassword}, testNow, "", &language.Afrikaans, &domain.UserAgent{
								FingerprintID: gu.Ptr("fp1"),
								IP:            net.ParseIP("1.2.3.4"),
								Description:   gu.Ptr("firefox"),
								Header:        http.Header{"foo": []string{"bar"}},
							},
						),
						sessionlogout.NewBackChannelLogoutRegisteredEvent(context.Background(),
							&sessionlogout.NewAggregate("sessionID", "instance1").Aggregate,
							"V2_oidcSessionID",
							"userID",
							"clientID",
							"backChannelLogoutURI",
						),
						oidcsession.NewAccessTokenAddedEvent(context.Background(),
							&oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
							"at_accessTokenID", []string{"openid", "offline_access"}, time.Hour, domain.TokenReasonAuthRequest, nil,
						),
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
				authz.WithFeatures(ctx, feature.Features{
					EnableBackChannelLogout: true,
				}),
				"123",
				"backChannelLogoutURI",
			},
			want: &OIDCSession{
				TokenID:           "V2_oidcSessionID-at_accessTokenID",
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
				Reason:    domain.TokenReasonAuthRequest,
				SessionID: "sessionID",
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
								"sessionID",
							),
						),
					),
					expectFilter(
						user.NewHumanAddedEvent(
							ctx,
							&user.NewAggregate("userID", "org1").Aggregate,
							"username",
							"firstname",
							"lastname",
							"nickname",
							"displayname",
							language.English,
							domain.GenderUnspecified,
							"email",
							false,
						),
					),
					expectFilter(), // token lifetime
					expectPush(
						oidcsession.NewAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
							"userID", "org1", "sessionID", "clientID", []string{"audience"}, []string{"openid", "offline_access"},
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
				"",
			},
			want: &OIDCSession{
				TokenID:           "V2_oidcSessionID-at_accessTokenID",
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
				RefreshToken: "VjJfb2lkY1Nlc3Npb25JRC1ydF9yZWZyZXNoVG9rZW5JRDp1c2VySUQ", //V2_oidcSessionID-rt_refreshTokenID:userID
				SessionID:    "sessionID",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:                      tt.fields.eventstore(t),
				idGenerator:                     tt.fields.idGenerator,
				defaultAccessTokenLifetime:      tt.fields.defaultAccessTokenLifetime,
				defaultRefreshTokenLifetime:     tt.fields.defaultRefreshTokenLifetime,
				defaultRefreshTokenIdleLifetime: tt.fields.defaultRefreshTokenIdleLifetime,
				keyAlgorithm:                    tt.fields.keyAlgorithm,
			}
			got, err := c.CreateOIDCSessionFromDeviceAuth(tt.args.ctx, tt.args.deviceCode, tt.args.backChannelLogoutURI)
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
