package command

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommandSide_GenerateMachineSecret(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		userID        string
		resourceOwner string
		generator     crypto.Generator
		set           *GenerateMachineSecret
	}
	type res struct {
		want   *domain.ObjectDetails
		secret *GenerateMachineSecret
		err    func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "user invalid, invalid argument error userID",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:           context.Background(),
				userID:        "",
				resourceOwner: "org1",
				generator:     GetMockSecretGenerator(t),
				set:           nil,
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "user invalid, invalid argument error resourceowner",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				resourceOwner: "",
				generator:     GetMockSecretGenerator(t),
				set:           nil,
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "user not existing, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				resourceOwner: "org1",
				generator:     GetMockSecretGenerator(t),
				set:           nil,
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "add machine secret, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							user.NewMachineAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"user1",
								"username",
								"user",
								false,
								domain.OIDCTokenTypeBearer,
							),
						),
					),
					expectPush(
						user.NewMachineSecretSetEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("a"),
							},
						),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				resourceOwner: "org1",
				generator:     GetMockSecretGenerator(t),
				set:           &GenerateMachineSecret{},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
				secret: &GenerateMachineSecret{
					ClientSecret: "a",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.GenerateMachineSecret(tt.args.ctx, tt.args.userID, tt.args.resourceOwner, tt.args.generator, tt.args.set)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, got)
				assert.Equal(t, tt.args.set.ClientSecret, tt.res.secret.ClientSecret)
			}
		})
	}
}

func TestCommandSide_RemoveMachineSecret(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		userID        string
		resourceOwner string
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
			name: "user invalid, invalid argument error userID",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:           context.Background(),
				userID:        "",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "user invalid, invalid argument error resourceowner",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				resourceOwner: "",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "user not existing, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "user existing without secret, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							user.NewMachineAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"user1",
								"username",
								"user",
								false,
								domain.OIDCTokenTypeBearer,
							),
						),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "remove machine secret, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							user.NewMachineAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"user1",
								"username",
								"user",
								false,
								domain.OIDCTokenTypeBearer,
							),
						),
						eventFromEventPusher(
							user.NewMachineSecretSetEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("a"),
								},
							),
						),
					),
					expectPush(
						user.NewMachineSecretRemovedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
						),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.RemoveMachineSecret(tt.args.ctx, tt.args.userID, tt.args.resourceOwner)
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

func TestCommands_MachineSecretCheckSucceeded(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	agg := user.NewAggregate("userID", "orgID")
	cmd := user.NewMachineSecretCheckSucceededEvent(ctx, &agg.Aggregate)

	c := &Commands{
		eventstore: eventstoreExpect(t,
			expectPushSlow(time.Second/100, cmd),
		),
	}
	c.MachineSecretCheckSucceeded(ctx, "userID", "orgID")
	require.NoError(t, c.Close(ctx))
}

func TestCommands_MachineSecretCheckFailed(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	agg := user.NewAggregate("userID", "orgID")
	cmd := user.NewMachineSecretCheckFailedEvent(ctx, &agg.Aggregate)

	c := &Commands{
		eventstore: eventstoreExpect(t,
			expectPushSlow(time.Second/100, cmd),
		),
	}
	c.MachineSecretCheckFailed(ctx, "userID", "orgID")
	require.NoError(t, c.Close(ctx))
}
