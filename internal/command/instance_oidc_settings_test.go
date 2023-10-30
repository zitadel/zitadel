package command

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

func TestCommandSide_AddOIDCConfig(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx        context.Context
		oidcConfig *domain.OIDCSettings
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
			name: "oidc settings, error already exists",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewOIDCSettingsAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								time.Hour*1,
								time.Hour*1,
								time.Hour*1,
								time.Hour*1,
							),
						),
					),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "INSTANCE"),
				oidcConfig: &domain.OIDCSettings{
					AccessTokenLifetime:        1 * time.Hour,
					IdTokenLifetime:            1 * time.Hour,
					RefreshTokenIdleExpiration: 1 * time.Hour,
					RefreshTokenExpiration:     1 * time.Hour,
				},
			},
			res: res{
				err: caos_errs.IsErrorAlreadyExists,
			},
		},
		{
			name: "add oidc settings, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
					expectPush(
						instance.NewOIDCSettingsAddedEvent(
							context.Background(),
							&instance.NewAggregate("INSTANCE").Aggregate,
							time.Hour*1,
							time.Hour*1,
							time.Hour*1,
							time.Hour*1,
						),
					),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "INSTANCE"),
				oidcConfig: &domain.OIDCSettings{
					AccessTokenLifetime:        1 * time.Hour,
					IdTokenLifetime:            1 * time.Hour,
					RefreshTokenIdleExpiration: 1 * time.Hour,
					RefreshTokenExpiration:     1 * time.Hour,
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "INSTANCE",
				},
			},
		},
		{
			name: "add oidc settings, invalid argument 1",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "INSTANCE"),
				oidcConfig: &domain.OIDCSettings{
					AccessTokenLifetime:        0 * time.Hour,
					IdTokenLifetime:            1 * time.Hour,
					RefreshTokenIdleExpiration: 1 * time.Hour,
					RefreshTokenExpiration:     1 * time.Hour,
				},
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "add oidc settings, invalid argument 2",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "INSTANCE"),
				oidcConfig: &domain.OIDCSettings{
					AccessTokenLifetime:        1 * time.Hour,
					IdTokenLifetime:            0 * time.Hour,
					RefreshTokenIdleExpiration: 1 * time.Hour,
					RefreshTokenExpiration:     1 * time.Hour,
				},
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "add oidc settings, invalid argument 3",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "INSTANCE"),
				oidcConfig: &domain.OIDCSettings{
					AccessTokenLifetime:        1 * time.Hour,
					IdTokenLifetime:            1 * time.Hour,
					RefreshTokenIdleExpiration: 0 * time.Hour,
					RefreshTokenExpiration:     1 * time.Hour,
				},
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "add oidc settings, invalid argument 4",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "INSTANCE"),
				oidcConfig: &domain.OIDCSettings{
					AccessTokenLifetime:        1 * time.Hour,
					IdTokenLifetime:            1 * time.Hour,
					RefreshTokenIdleExpiration: 1 * time.Hour,
					RefreshTokenExpiration:     0 * time.Hour,
				},
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.AddOIDCSettings(tt.args.ctx, tt.args.oidcConfig)
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

func TestCommandSide_ChangeOIDCConfig(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx        context.Context
		oidcConfig *domain.OIDCSettings
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
			name: "oidc settings not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx: context.Background(),
				oidcConfig: &domain.OIDCSettings{
					AccessTokenLifetime:        1 * time.Hour,
					IdTokenLifetime:            1 * time.Hour,
					RefreshTokenIdleExpiration: 1 * time.Hour,
					RefreshTokenExpiration:     1 * time.Hour,
				},
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "no changes, invalid argument error 1",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "INSTANCE"),
				oidcConfig: &domain.OIDCSettings{
					AccessTokenLifetime:        0 * time.Hour,
					IdTokenLifetime:            1 * time.Hour,
					RefreshTokenIdleExpiration: 1 * time.Hour,
					RefreshTokenExpiration:     1 * time.Hour,
				},
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "no changes, invalid argument error 2",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "INSTANCE"),
				oidcConfig: &domain.OIDCSettings{
					AccessTokenLifetime:        1 * time.Hour,
					IdTokenLifetime:            0 * time.Hour,
					RefreshTokenIdleExpiration: 1 * time.Hour,
					RefreshTokenExpiration:     1 * time.Hour,
				},
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "no changes, invalid argument error 3",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "INSTANCE"),
				oidcConfig: &domain.OIDCSettings{
					AccessTokenLifetime:        1 * time.Hour,
					IdTokenLifetime:            1 * time.Hour,
					RefreshTokenIdleExpiration: 0 * time.Hour,
					RefreshTokenExpiration:     1 * time.Hour,
				},
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "no changes, invalid argument error 4",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "INSTANCE"),
				oidcConfig: &domain.OIDCSettings{
					AccessTokenLifetime:        1 * time.Hour,
					IdTokenLifetime:            1 * time.Hour,
					RefreshTokenIdleExpiration: 1 * time.Hour,
					RefreshTokenExpiration:     0 * time.Hour,
				},
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "no changes, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewOIDCSettingsAddedEvent(
								context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								time.Hour*1,
								time.Hour*1,
								time.Hour*1,
								time.Hour*1,
							),
						),
					),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "INSTANCE"),
				oidcConfig: &domain.OIDCSettings{
					AccessTokenLifetime:        1 * time.Hour,
					IdTokenLifetime:            1 * time.Hour,
					RefreshTokenIdleExpiration: 1 * time.Hour,
					RefreshTokenExpiration:     1 * time.Hour,
				},
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "oidc settings change, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewOIDCSettingsAddedEvent(
								context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								time.Hour*1,
								time.Hour*1,
								time.Hour*1,
								time.Hour*1,
							),
						),
					),
					expectPush(
						newOIDCConfigChangedEvent(
							context.Background(),
							time.Hour*2,
							time.Hour*2,
							time.Hour*2,
							time.Hour*2),
					),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "INSTANCE"),
				oidcConfig: &domain.OIDCSettings{
					AccessTokenLifetime:        2 * time.Hour,
					IdTokenLifetime:            2 * time.Hour,
					RefreshTokenIdleExpiration: 2 * time.Hour,
					RefreshTokenExpiration:     2 * time.Hour,
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "INSTANCE",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.ChangeOIDCSettings(tt.args.ctx, tt.args.oidcConfig)
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

func newOIDCConfigChangedEvent(ctx context.Context, accessTokenLifetime, idTokenLifetime, refreshTokenIdleExpiration, refreshTokenExpiration time.Duration) *instance.OIDCSettingsChangedEvent {
	changes := []instance.OIDCSettingsChanges{
		instance.ChangeOIDCSettingsAccessTokenLifetime(accessTokenLifetime),
		instance.ChangeOIDCSettingsIdTokenLifetime(idTokenLifetime),
		instance.ChangeOIDCSettingsRefreshTokenIdleExpiration(refreshTokenIdleExpiration),
		instance.ChangeOIDCSettingsRefreshTokenExpiration(refreshTokenExpiration),
	}
	event, _ := instance.NewOIDCSettingsChangeEvent(ctx,
		&instance.NewAggregate("INSTANCE").Aggregate,
		changes,
	)
	return event
}
