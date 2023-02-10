package command

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/id"
	id_mock "github.com/zitadel/zitadel/internal/id/mock"
	"github.com/zitadel/zitadel/internal/repository/user"
)

func TestCommands_AddMachineKey(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		idGenerator  id.Generator
		keyAlgorithm crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx context.Context
		key *MachineKey
	}
	type res struct {
		want *domain.ObjectDetails
		key  bool
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"user does not exist, error",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "token1"),
			},
			args{
				ctx: context.Background(),
				key: &MachineKey{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "user1",
						ResourceOwner: "org1",
					},
					Type:           domain.AuthNKeyTypeJSON,
					ExpirationDate: time.Time{},
				},
			},
			res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			"invalid expiration date, error",
			fields{
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "token1"),
			},
			args{
				ctx: context.Background(),
				key: &MachineKey{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "user1",
						ResourceOwner: "org1",
					},
					Type:           domain.AuthNKeyTypeJSON,
					ExpirationDate: time.Now().Add(-24 * time.Hour),
				},
			},
			res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			"no userID, error",
			fields{
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "token1"),
			},
			args{
				ctx: context.Background(),
				key: &MachineKey{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "",
						ResourceOwner: "org1",
					},
					Type:           domain.AuthNKeyTypeJSON,
					ExpirationDate: time.Time{},
				},
			},
			res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			"no resourceowner, error",
			fields{
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "token1"),
			},
			args{
				ctx: context.Background(),
				key: &MachineKey{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "user1",
						ResourceOwner: "",
					},
					Type:           domain.AuthNKeyTypeJSON,
					ExpirationDate: time.Time{},
				},
			},
			res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			"key added with public key",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							user.NewMachineAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"machine",
								"Machine",
								"",
								true,
								domain.OIDCTokenTypeBearer,
							),
						),
					),
					expectFilter(),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								user.NewMachineKeyAddedEvent(context.Background(),
									&user.NewAggregate("user1", "org1").Aggregate,
									"key1",
									domain.AuthNKeyTypeJSON,
									time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC),
									[]byte("public"),
								),
							),
						},
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "key1"),
				keyAlgorithm: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args{
				ctx: context.Background(),
				key: &MachineKey{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "user1",
						ResourceOwner: "org1",
					},
					Type:           domain.AuthNKeyTypeJSON,
					ExpirationDate: time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC),
					PublicKey:      []byte("public"),
				},
			},
			res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
				key: true,
			},
		},
		{
			"key added with ID and public key",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							user.NewMachineAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"machine",
								"Machine",
								"",
								true,
								domain.OIDCTokenTypeBearer,
							),
						),
					),
					expectFilter(),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								user.NewMachineKeyAddedEvent(context.Background(),
									&user.NewAggregate("user1", "org1").Aggregate,
									"key1",
									domain.AuthNKeyTypeJSON,
									time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC),
									[]byte("public"),
								),
							),
						},
					),
				),
				keyAlgorithm: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args{
				ctx: context.Background(),
				key: &MachineKey{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "user1",
						ResourceOwner: "org1",
					},
					KeyID:          "key1",
					Type:           domain.AuthNKeyTypeJSON,
					ExpirationDate: time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC),
					PublicKey:      []byte("public"),
				},
			},
			res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
				key: true,
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
			got, err := c.AddUserMachineKey(tt.args.ctx, tt.args.key)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, got)
				if tt.res.key {
					assert.NotEqual(t, "", tt.args.key.PrivateKey)
				}
			}
		})
	}
}
