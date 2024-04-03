package command

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/deviceauth"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommands_AddDeviceAuth(t *testing.T) {
	ctx := authz.WithInstanceID(context.Background(), "instance1")
	pushErr := errors.New("pushErr")
	now := time.Now()

	unique := deviceauth.NewAddUniqueConstraints("123", "456")
	require.Len(t, unique, 2)

	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx        context.Context
		clientID   string
		deviceCode string
		userCode   string
		expires    time.Time
		scopes     []string
		audience   []string
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
				eventstore: eventstoreExpect(t, expectPush(
					deviceauth.NewAddedEvent(
						ctx,
						deviceauth.NewAggregate("123", "instance1"),
						"client_id", "123", "456", now,
						[]string{"a", "b", "c"},
						[]string{"projectID", "clientID"},
					),
				)),
			},
			args: args{
				ctx:        authz.WithInstanceID(context.Background(), "instance1"),
				clientID:   "client_id",
				deviceCode: "123",
				userCode:   "456",
				expires:    now,
				scopes:     []string{"a", "b", "c"},
				audience:   []string{"projectID", "clientID"},
			},
			wantDetails: &domain.ObjectDetails{
				ResourceOwner: "instance1",
			},
		},
		{
			name: "push error",
			fields: fields{
				eventstore: eventstoreExpect(t, expectPushFailed(pushErr,
					deviceauth.NewAddedEvent(
						ctx,
						deviceauth.NewAggregate("123", "instance1"),
						"client_id", "123", "456", now,
						[]string{"a", "b", "c"},
						[]string{"projectID", "clientID"},
					)),
				),
			},
			args: args{
				ctx:        authz.WithInstanceID(context.Background(), "instance1"),
				clientID:   "client_id",
				deviceCode: "123",
				userCode:   "456",
				expires:    now,
				scopes:     []string{"a", "b", "c"},
				audience:   []string{"projectID", "clientID"},
			},
			wantErr: pushErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.fields.eventstore,
			}
			gotDetails, err := c.AddDeviceAuth(tt.args.ctx, tt.args.clientID, tt.args.deviceCode, tt.args.userCode, tt.args.expires, tt.args.scopes, tt.args.audience)
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
		ctx         context.Context
		id          string
		subject     string
		authMethods []domain.UserAuthMethodType
		authTime    time.Time
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
				ctx, "123", "subj",
				[]domain.UserAuthMethodType{domain.UserAuthMethodTypePassword},
				time.Unix(123, 456),
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
							[]string{"projectID", "clientID"},
						),
					)),
					expectPushFailed(pushErr,
						deviceauth.NewApprovedEvent(
							ctx, deviceauth.NewAggregate("123", "instance1"), "subj",
							[]domain.UserAuthMethodType{domain.UserAuthMethodTypePassword},
							time.Unix(123, 456),
						),
					),
				),
			},
			args: args{
				ctx, "123", "subj",
				[]domain.UserAuthMethodType{domain.UserAuthMethodTypePassword},
				time.Unix(123, 456),
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
							[]string{"projectID", "clientID"},
						),
					)),
					expectPush(
						deviceauth.NewApprovedEvent(
							ctx, deviceauth.NewAggregate("123", "instance1"), "subj",
							[]domain.UserAuthMethodType{domain.UserAuthMethodTypePassword},
							time.Unix(123, 456),
						),
					),
				),
			},
			args: args{
				ctx, "123", "subj",
				[]domain.UserAuthMethodType{domain.UserAuthMethodTypePassword},
				time.Unix(123, 456),
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
			gotDetails, err := c.ApproveDeviceAuth(tt.args.ctx, tt.args.id, tt.args.subject, tt.args.authMethods, tt.args.authTime)
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
							[]string{"projectID", "clientID"},
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
							[]string{"projectID", "clientID"},
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
							[]string{"projectID", "clientID"},
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
