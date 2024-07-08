package command

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/id_generator"
	id_mock "github.com/zitadel/zitadel/internal/id_generator/mock"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const fakePubkey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAp4qNBuUu/HekF2E5bOtA
oEL76zS0NQdZL3ByEJ3hZplJhE30ITPIOLW3+uaMMM+obl/LLapwG2vdhvutQtx/
FOLJmXysbG3RL9zjXDBT5IE+nGFC7ctsi5FGbHQbAm45E3HHCSk7gfmTy9hxyk1K
GsyU8BDeOWasJO6aeXqpOnRM8vw/fY+6mHVC9CxcIroSfrIabFGe/mP6qpBGeFSn
APymBc/8lca4JaPv2/u/rBhnaAHZiUuCS1+MonWelOb+MSfq48VgtpiaYIVY9szI
esorA6EJ9pO17ROEUpX5wP5Oir+yGJU27jSvLCjvK6fOFX+OwUM9L8047JKoo+Nf
PwIDAQAB
-----END PUBLIC KEY-----`

func TestCommands_AddMachineKey(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		idGenerator  id_generator.Generator
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
				err: zerrors.IsPreconditionFailed,
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
				err: zerrors.IsErrorInvalidArgument,
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
				err: zerrors.IsErrorInvalidArgument,
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
				err: zerrors.IsErrorInvalidArgument,
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
						user.NewMachineKeyAddedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							"key1",
							domain.AuthNKeyTypeJSON,
							time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC),
							[]byte(fakePubkey),
						),
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
					PublicKey:      []byte(fakePubkey),
				},
			},
			res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
				key: false,
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
						user.NewMachineKeyAddedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							"key1",
							domain.AuthNKeyTypeJSON,
							time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC),
							[]byte(fakePubkey),
						),
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
					PublicKey:      []byte(fakePubkey),
				},
			},
			res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
				key: false,
			},
		},
		{
			"key added with invalid public key",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx: context.Background(),
				key: &MachineKey{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "user1",
						ResourceOwner: "org1",
					},
					KeyID:     "key1",
					Type:      domain.AuthNKeyTypeJSON,
					PublicKey: []byte("incorrect"),
				},
			},
			res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:   tt.fields.eventstore,
				keyAlgorithm: tt.fields.keyAlgorithm,
			}
			id_generator.SetGenerator(tt.fields.idGenerator)
			got, err := c.AddUserMachineKey(tt.args.ctx, tt.args.key)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, got)
				receivedKey := len(tt.args.key.PrivateKey) > 0
				assert.Equal(t, tt.res.key, receivedKey)
			}
		})
	}
}
