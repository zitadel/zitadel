package command

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/repository/idpconfig"
	"github.com/caos/zitadel/internal/repository/org"
)

func TestCommandSide_ChangeIDPJWTConfig(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		secretCrypto crypto.EncryptionAlgorithm
	}
	type (
		args struct {
			ctx           context.Context
			config        *domain.JWTIDPConfig
			resourceOwner string
		}
	)
	type res struct {
		want *domain.JWTIDPConfig
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "resourceowner missing, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: context.Background(),
				config: &domain.JWTIDPConfig{
					IDPConfigID: "config1",
				},
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "invalid config, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:           context.Background(),
				config:        &domain.JWTIDPConfig{},
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "idp config not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx: context.Background(),
				config: &domain.JWTIDPConfig{
					IDPConfigID: "config1",
				},
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "idp config removed, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewIDPConfigAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"config1",
								"name1",
								domain.IDPConfigTypeJWT,
								domain.IDPConfigStylingTypeGoogle,
							),
						),
						eventFromEventPusher(
							org.NewIDPJWTConfigAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"config1",
								"issuer",
								"keys-endpoint",
							),
						),
						eventFromEventPusher(
							org.NewIDPConfigRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"config1",
								"name",
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				config: &domain.JWTIDPConfig{
					IDPConfigID: "config1",
				},
				resourceOwner: "org1",
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
							org.NewIDPConfigAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"config1",
								"name1",
								domain.IDPConfigTypeJWT,
								domain.IDPConfigStylingTypeGoogle,
							),
						),
						eventFromEventPusher(
							org.NewIDPJWTConfigAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"config1",
								"issuer",
								"keys-endpoint",
							),
						),
					),
				),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: context.Background(),
				config: &domain.JWTIDPConfig{
					IDPConfigID:  "config1",
					Issuer:       "issuer",
					KeysEndpoint: "keys-endpoint",
				},
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "idp config jwt add, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewIDPConfigAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"config1",
								"name1",
								domain.IDPConfigTypeJWT,
								domain.IDPConfigStylingTypeGoogle,
							),
						),
						eventFromEventPusher(
							org.NewIDPJWTConfigAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"config1",
								"issuer",
								"keys-endpoint",
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								newIDPJWTConfigChangedEvent(context.Background(),
									"org1",
									"config1",
									"issuer-changed",
									"keys-endpoint-changed",
								),
							),
						},
					),
				),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: context.Background(),
				config: &domain.JWTIDPConfig{
					IDPConfigID:  "config1",
					Issuer:       "issuer-changed",
					KeysEndpoint: "keys-endpoint-changed",
				},
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.JWTIDPConfig{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "org1",
						ResourceOwner: "org1",
					},
					IDPConfigID:  "config1",
					Issuer:       "issuer-changed",
					KeysEndpoint: "keys-endpoint-changed",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore:            tt.fields.eventstore,
				idpConfigSecretCrypto: tt.fields.secretCrypto,
			}
			got, err := r.ChangeIDPJWTConfig(tt.args.ctx, tt.args.config, tt.args.resourceOwner)
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

func newIDPJWTConfigChangedEvent(ctx context.Context, orgID, configID, issuer, keysEndpoint string) *org.IDPJWTConfigChangedEvent {
	event, _ := org.NewIDPJWTConfigChangedEvent(ctx,
		&org.NewAggregate(orgID, orgID).Aggregate,
		configID,
		[]idpconfig.JWTConfigChanges{
			idpconfig.ChangeJWTIssuer(issuer),
			idpconfig.ChangeKeysEndpoint(keysEndpoint),
		},
	)
	return event
}
