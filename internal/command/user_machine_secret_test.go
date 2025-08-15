package command

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommandSide_GenerateMachineSecret(t *testing.T) {
	type fields struct {
		eventstore func(*testing.T) *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		userID        string
		resourceOwner string
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
				eventstore: expectEventstore(),
			},
			args: args{
				ctx:           context.Background(),
				userID:        "",
				resourceOwner: "org1",
				set:           new(GenerateMachineSecret),
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "user invalid, invalid argument error resourceowner",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				resourceOwner: "",
				set:           new(GenerateMachineSecret),
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "user not existing, precondition error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				resourceOwner: "org1",
				set:           new(GenerateMachineSecret),
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "add machine secret, ok",
			fields: fields{
				eventstore: expectEventstore(
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
							"secret",
						),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				resourceOwner: "org1",
				set:           &GenerateMachineSecret{},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
				secret: &GenerateMachineSecret{
					ClientSecret: "secret",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore:      tt.fields.eventstore(t),
				newHashedSecret: mockHashedSecret("secret"),
				defaultSecretGenerators: &SecretGenerators{
					ClientSecret: emptyConfig,
				},
			}
			got, err := r.GenerateMachineSecret(tt.args.ctx, tt.args.userID, tt.args.resourceOwner, tt.args.set)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assertObjectDetails(t, tt.res.want, got)
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
								"secret",
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
			got, err := r.RemoveMachineSecret(tt.args.ctx, tt.args.userID, tt.args.resourceOwner, nil)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assertObjectDetails(t, tt.res.want, got)
			}
		})
	}
}

func TestCommands_MachineSecretCheckSucceeded(t *testing.T) {

	tt := []struct {
		testName               string
		eventStoreExpectations func(ctx context.Context, updatedSecret string) []expect

		inputUpdatedMachineSecret string
	}{
		{
			testName:                  "when machine secret is updated should meet expectations and return no error",
			inputUpdatedMachineSecret: "upd4t3dS3cr3t",
			eventStoreExpectations: func(ctx context.Context, updatedSecret string) []expect {
				return []expect{expectPushSlow(time.Second/100, user.NewMachineSecretHashUpdatedEvent(
					ctx,
					&user.NewAggregate("userID", "orgID").Aggregate,
					updatedSecret,
				))}
			},
		},
		{
			testName: "when machine secret is not update should have no expectations and return no error",
			eventStoreExpectations: func(_ context.Context, _ string) []expect {
				return []expect{}
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			// Given
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			c := &Commands{eventstore: expectEventstore(tc.eventStoreExpectations(ctx, tc.inputUpdatedMachineSecret)...)(t)}

			// Test
			c.MachineSecretCheckSucceeded(ctx, "userID", "orgID", tc.inputUpdatedMachineSecret)

			// Verify
			require.NoError(t, c.Close(ctx))
		})
	}
}
