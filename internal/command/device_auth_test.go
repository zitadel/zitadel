package command

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id"
	id_mock "github.com/zitadel/zitadel/internal/id/mock"
	"github.com/zitadel/zitadel/internal/repository/deviceauth"
)

func TestCommands_AddDeviceAuth(t *testing.T) {
	ctx := authz.WithInstanceID(context.Background(), "instance1")
	idErr := errors.New("idErr")
	pushErr := errors.New("pushErr")
	now := time.Now()

	unique := deviceauth.NewAddUniqueConstraints("client_id", "123", "456")
	require.Len(t, unique, 2)

	type fields struct {
		eventstore  *eventstore.Eventstore
		idGenerator id.Generator
	}
	type args struct {
		ctx        context.Context
		clientID   string
		deviceCode string
		userCode   string
		expires    time.Time
		scopes     []string
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantID      string
		wantDetails *domain.ObjectDetails
		wantErr     error
	}{
		{
			name: "idGenerator error",
			fields: fields{
				eventstore: eventstoreExpect(t),
				idGenerator: func() id.Generator {
					m := id_mock.NewMockGenerator(gomock.NewController(t))
					m.EXPECT().Next().Return("", idErr)
					return m
				}(),
			},
			args: args{
				ctx:        ctx,
				clientID:   "client_id",
				deviceCode: "123",
				userCode:   "456",
				expires:    now,
				scopes:     []string{"a", "b", "c"},
			},
			wantErr: idErr,
		},
		{
			name: "success",
			fields: fields{
				eventstore: eventstoreExpect(t, expectPush(
					deviceauth.NewAddedEvent(
						ctx,
						deviceauth.NewAggregate("1999", "instance1"),
						"client_id", "123", "456", now,
						[]string{"a", "b", "c"},
					),
				)),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "1999"),
			},
			args: args{
				ctx:        authz.WithInstanceID(context.Background(), "instance1"),
				clientID:   "client_id",
				deviceCode: "123",
				userCode:   "456",
				expires:    now,
				scopes:     []string{"a", "b", "c"},
			},
			wantID: "1999",
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
						deviceauth.NewAggregate("1999", "instance1"),
						"client_id", "123", "456", now,
						[]string{"a", "b", "c"},
					)),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "1999"),
			},
			args: args{
				ctx:        authz.WithInstanceID(context.Background(), "instance1"),
				clientID:   "client_id",
				deviceCode: "123",
				userCode:   "456",
				expires:    now,
				scopes:     []string{"a", "b", "c"},
			},
			wantErr: pushErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:  tt.fields.eventstore,
				idGenerator: tt.fields.idGenerator,
			}
			gotID, gotDetails, err := c.AddDeviceAuth(tt.args.ctx, tt.args.clientID, tt.args.deviceCode, tt.args.userCode, tt.args.expires, tt.args.scopes)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.wantID, gotID)
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
		ctx     context.Context
		id      string
		subject string
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
					expectFilter(
						eventFromEventPusherWithInstanceID("instance1",
							deviceauth.NewAddedEvent(
								ctx,
								deviceauth.NewAggregate("1999", "instance1"),
								"client_id", "123", "456", now,
								[]string{"a", "b", "c"},
							),
						),
						eventFromEventPusherWithInstanceID("instance1",
							deviceauth.NewRemovedEvent(
								ctx,
								deviceauth.NewAggregate("1999", "instance1"),
								"client_id", "123", "456",
							),
						),
					),
				),
			},
			args:    args{ctx, "1999", "subj"},
			wantErr: caos_errs.ThrowNotFound(nil, "COMMAND-Hief9", "Errors.DeviceAuth.NotFound"),
		},
		{
			name: "push error",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(eventFromEventPusherWithInstanceID(
						"instance1",
						deviceauth.NewAddedEvent(
							ctx,
							deviceauth.NewAggregate("1999", "instance1"),
							"client_id", "123", "456", now,
							[]string{"a", "b", "c"},
						),
					)),
					expectPushFailed(pushErr,
						deviceauth.NewApprovedEvent(
							ctx, deviceauth.NewAggregate("1999", "instance1"), "subj",
						),
					),
				),
			},
			args:    args{ctx, "1999", "subj"},
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
							deviceauth.NewAggregate("1999", "instance1"),
							"client_id", "123", "456", now,
							[]string{"a", "b", "c"},
						),
					)),
					expectPush(
						deviceauth.NewApprovedEvent(
							ctx, deviceauth.NewAggregate("1999", "instance1"), "subj",
						),
					),
				),
			},
			args: args{ctx, "1999", "subj"},
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
			gotDetails, err := c.ApproveDeviceAuth(tt.args.ctx, tt.args.id, tt.args.subject)
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
					expectFilter(
						eventFromEventPusherWithInstanceID("instance1",
							deviceauth.NewAddedEvent(
								ctx,
								deviceauth.NewAggregate("1999", "instance1"),
								"client_id", "123", "456", now,
								[]string{"a", "b", "c"},
							),
						),
						eventFromEventPusherWithInstanceID("instance1",
							deviceauth.NewRemovedEvent(
								ctx,
								deviceauth.NewAggregate("1999", "instance1"),
								"client_id", "123", "456",
							),
						),
					),
				),
			},
			args:    args{ctx, "1999", domain.DeviceAuthCanceledDenied},
			wantErr: caos_errs.ThrowNotFound(nil, "COMMAND-gee5A", "Errors.DeviceAuth.NotFound"),
		},
		{
			name: "push error",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(eventFromEventPusherWithInstanceID(
						"instance1",
						deviceauth.NewAddedEvent(
							ctx,
							deviceauth.NewAggregate("1999", "instance1"),
							"client_id", "123", "456", now,
							[]string{"a", "b", "c"},
						),
					)),
					expectPushFailed(pushErr,
						deviceauth.NewCanceledEvent(
							ctx, deviceauth.NewAggregate("1999", "instance1"),
							domain.DeviceAuthCanceledDenied,
						),
					),
				),
			},
			args:    args{ctx, "1999", domain.DeviceAuthCanceledDenied},
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
							deviceauth.NewAggregate("1999", "instance1"),
							"client_id", "123", "456", now,
							[]string{"a", "b", "c"},
						),
					)),
					expectPush(
						deviceauth.NewCanceledEvent(
							ctx, deviceauth.NewAggregate("1999", "instance1"),
							domain.DeviceAuthCanceledDenied,
						),
					),
				),
			},
			args: args{ctx, "1999", domain.DeviceAuthCanceledDenied},
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
							deviceauth.NewAggregate("1999", "instance1"),
							"client_id", "123", "456", now,
							[]string{"a", "b", "c"},
						),
					)),
					expectPush(
						deviceauth.NewCanceledEvent(
							ctx, deviceauth.NewAggregate("1999", "instance1"),
							domain.DeviceAuthCanceledExpired,
						),
					),
				),
			},
			args: args{ctx, "1999", domain.DeviceAuthCanceledExpired},
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

func TestCommands_RemoveDeviceAuth(t *testing.T) {
	ctx := authz.WithInstanceID(context.Background(), "instance1")
	now := time.Now()
	pushErr := errors.New("pushErr")

	unique := deviceauth.NewRemoveUniqueConstraints("client_id", "123", "456")
	require.Len(t, unique, 2)

	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantDetails *domain.ObjectDetails
		wantErr     error
	}{
		{
			name: "push error",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(eventFromEventPusherWithInstanceID(
						"instance1",
						deviceauth.NewAddedEvent(
							ctx,
							deviceauth.NewAggregate("1999", "instance1"),
							"client_id", "123", "456", now,
							[]string{"a", "b", "c"},
						),
					)),
					expectPushFailed(pushErr,
						deviceauth.NewRemovedEvent(
							ctx, deviceauth.NewAggregate("1999", "instance1"),
							"client_id", "123", "456",
						),
					),
				),
			},
			args:    args{ctx, "1999"},
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
							deviceauth.NewAggregate("1999", "instance1"),
							"client_id", "123", "456", now,
							[]string{"a", "b", "c"},
						),
					)),
					expectPush(
						deviceauth.NewRemovedEvent(
							ctx, deviceauth.NewAggregate("1999", "instance1"),
							"client_id", "123", "456",
						),
					),
				),
			},
			args: args{ctx, "1999"},
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
			gotDetails, err := c.RemoveDeviceAuth(tt.args.ctx, tt.args.id)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, gotDetails, tt.wantDetails)
		})
	}
}
