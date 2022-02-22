package command

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/repository/iam"
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
			name: "oidc config, error already exists",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							iam.NewOIDCSettingsAddedEvent(context.Background(),
								&iam.NewAggregate().Aggregate,
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
				ctx: context.Background(),
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
			name: "add secret generator, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(iam.NewOIDCSettingsAddedEvent(
								context.Background(),
								&iam.NewAggregate().Aggregate,
								time.Hour*1,
								time.Hour*1,
								time.Hour*1,
								time.Hour*1,
							),
							),
						},
					),
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
				want: &domain.ObjectDetails{
					ResourceOwner: "IAM",
				},
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
			name: "oidc config not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx: context.Background(),
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "no changes, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							iam.NewOIDCSettingsAddedEvent(
								context.Background(),
								&iam.NewAggregate().Aggregate,
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
				ctx: context.Background(),
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
			name: "secret generator change, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							iam.NewOIDCSettingsAddedEvent(
								context.Background(),
								&iam.NewAggregate().Aggregate,
								time.Hour*1,
								time.Hour*1,
								time.Hour*1,
								time.Hour*1,
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								newOIDCConfigChangedEvent(context.Background(),
									time.Hour*2,
									time.Hour*2,
									time.Hour*2,
									time.Hour*2),
							),
						},
					),
				),
			},
			args: args{
				ctx: context.Background(),
				oidcConfig: &domain.OIDCSettings{
					AccessTokenLifetime:        2 * time.Hour,
					IdTokenLifetime:            2 * time.Hour,
					RefreshTokenIdleExpiration: 2 * time.Hour,
					RefreshTokenExpiration:     2 * time.Hour,
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "IAM",
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

func newOIDCConfigChangedEvent(ctx context.Context, accessTokenLifetime, idTokenLifetime, refreshTokenIdleExpiration, refreshTokenExpiration time.Duration) *iam.OIDCSettingsChangedEvent {
	changes := []iam.OIDCSettingsChanges{
		iam.ChangeOIDCSettingsAccessTokenLifetime(accessTokenLifetime),
		iam.ChangeOIDCSettingsIdTokenLifetime(idTokenLifetime),
		iam.ChangeOIDCSettingsRefreshTokenIdleExpiration(refreshTokenIdleExpiration),
		iam.ChangeOIDCSettingsRefreshTokenExpiration(refreshTokenExpiration),
	}
	event, _ := iam.NewOIDCSettingsChangeEvent(ctx,
		&iam.NewAggregate().Aggregate,
		changes,
	)
	return event
}
