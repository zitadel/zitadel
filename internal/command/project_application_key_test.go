package command

import (
	"context"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/id"
	id_mock "github.com/caos/zitadel/internal/id/mock"
	"github.com/caos/zitadel/internal/repository/project"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCommandSide_AddAPIApplicationKey(t *testing.T) {
	type fields struct {
		eventstore      *eventstore.Eventstore
		idGenerator     id.Generator
		secretGenerator crypto.Generator
		keySize         int
	}
	type args struct {
		ctx           context.Context
		key           *domain.ApplicationKey
		resourceOwner string
	}
	type res struct {
		want *domain.APIApp
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "no aggregateid, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: context.Background(),
				key: &domain.ApplicationKey{
					ApplicationID: "app1",
				},
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "no appid, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: context.Background(),
				key: &domain.ApplicationKey{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
				},
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "app not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx: context.Background(),
				key: &domain.ApplicationKey{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					ApplicationID: "app1",
				},
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "create key not allowed, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							project.NewApplicationAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"app1",
								"app",
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							project.NewAPIConfigAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"app1",
								"client1@project",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("a"),
								},
								domain.APIAuthMethodTypeBasic),
						),
					),
				),
				idGenerator:     id_mock.NewIDGeneratorExpectIDs(t, "key1"),
				secretGenerator: GetMockSecretGenerator(t),
			},
			args: args{
				ctx: context.Background(),
				key: &domain.ApplicationKey{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					ApplicationID: "app1",
				},
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "create key not allowed, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							project.NewApplicationAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"app1",
								"app",
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							project.NewAPIConfigAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"app1",
								"client1@project",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("a"),
								},
								domain.APIAuthMethodTypeBasic),
						),
					),
				),
				idGenerator:     id_mock.NewIDGeneratorExpectIDs(t, "key1"),
				secretGenerator: GetMockSecretGenerator(t),
				keySize:         10,
			},
			args: args{
				ctx: context.Background(),
				key: &domain.ApplicationKey{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					ApplicationID: "app1",
				},
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore:                 tt.fields.eventstore,
				idGenerator:                tt.fields.idGenerator,
				applicationSecretGenerator: tt.fields.secretGenerator,
				applicationKeySize:         tt.fields.keySize,
			}
			got, err := r.AddApplicationKey(tt.args.ctx, tt.args.key, tt.args.resourceOwner)
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
