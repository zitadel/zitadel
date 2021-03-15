package command

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/repository/user"
)

func TestCommandSide_ChangeHumanEmail(t *testing.T) {
	type fields struct {
		eventstore      *eventstore.Eventstore
		secretGenerator crypto.Generator
	}
	type args struct {
		ctx           context.Context
		email         *domain.Email
		resourceOwner string
	}
	type res struct {
		want *domain.Email
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "invalid email, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: context.Background(),
				email: &domain.Email{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "user1",
					},
				},
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
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
				ctx: context.Background(),
				email: &domain.Email{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "user1",
					},
					EmailAddress: "email@test.com",
				},
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "email not changed, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.German,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				email: &domain.Email{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "user1",
					},
					EmailAddress: "email@test.ch",
				},
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "verified email changed, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.German,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								user.NewHumanEmailChangedEvent(context.Background(),
									&user.NewAggregate("user1", "org1").Aggregate,
									"email-changed@test.ch",
								),
							),
							eventFromEventPusher(
								user.NewHumanEmailVerifiedEvent(context.Background(),
									&user.NewAggregate("user1", "org1").Aggregate,
								),
							),
						},
					),
				),
			},
			args: args{
				ctx: context.Background(),
				email: &domain.Email{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "user1",
					},
					EmailAddress:    "email-changed@test.ch",
					IsEmailVerified: true,
				},
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.Email{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "user1",
						ResourceOwner: "org1",
					},
					EmailAddress:    "email-changed@test.ch",
					IsEmailVerified: true,
				},
			},
		},
		{
			name: "email changed with code, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.German,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								user.NewHumanEmailChangedEvent(context.Background(),
									&user.NewAggregate("user1", "org1").Aggregate,
									"email-changed@test.ch",
								),
							),
							eventFromEventPusher(
								user.NewHumanEmailCodeAddedEvent(context.Background(),
									&user.NewAggregate("user1", "org1").Aggregate,
									&crypto.CryptoValue{
										CryptoType: crypto.TypeEncryption,
										Algorithm:  "enc",
										KeyID:      "id",
										Crypted:    []byte("a"),
									},
									time.Hour*1,
								),
							),
						},
					),
				),
				secretGenerator: GetMockSecretGenerator(t),
			},
			args: args{
				ctx: context.Background(),
				email: &domain.Email{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "user1",
					},
					EmailAddress: "email-changed@test.ch",
				},
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.Email{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "user1",
						ResourceOwner: "org1",
					},
					EmailAddress: "email-changed@test.ch",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore:            tt.fields.eventstore,
				emailVerificationCode: tt.fields.secretGenerator,
			}
			got, err := r.ChangeHumanEmail(tt.args.ctx, tt.args.email)
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
