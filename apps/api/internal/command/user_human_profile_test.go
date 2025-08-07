package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommandSide_ChangeHumanProfile(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		address       *domain.Profile
		resourceOwner string
	}
	type res struct {
		want *domain.Profile
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "user not existing, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
				),
			},
			args: args{
				ctx: context.Background(),
				address: &domain.Profile{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "user1",
					},
					FirstName:         "firstname",
					LastName:          "lastname",
					NickName:          "nickname",
					DisplayName:       "displayname",
					PreferredLanguage: AllowedLanguage,
					Gender:            domain.GenderFemale,
				},
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "profile not changed, precondition error",
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
								AllowedLanguage,
								domain.GenderFemale,
								"email",
								true,
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				address: &domain.Profile{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "user1",
					},
					FirstName:         "firstname",
					LastName:          "lastname",
					NickName:          "nickname",
					DisplayName:       "displayname",
					PreferredLanguage: AllowedLanguage,
					Gender:            domain.GenderFemale,
				},
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "profile changed, ok",
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
								DisallowedLanguage,
								domain.GenderUnspecified,
								"email",
								true,
							),
						),
					),
					expectPush(
						newProfileChangedEvent(context.Background(),
							"user1", "org1",
							"firstname2",
							"lastname2",
							"nickname2",
							"displayname2",
							AllowedLanguage,
							domain.GenderMale,
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				address: &domain.Profile{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "user1",
					},
					FirstName:         "firstname2",
					LastName:          "lastname2",
					NickName:          "nickname2",
					DisplayName:       "displayname2",
					PreferredLanguage: AllowedLanguage,
					Gender:            domain.GenderMale,
				},
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.Profile{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "user1",
						ResourceOwner: "org1",
					},
					FirstName:         "firstname2",
					LastName:          "lastname2",
					NickName:          "nickname2",
					DisplayName:       "displayname2",
					PreferredLanguage: AllowedLanguage,
					Gender:            domain.GenderMale,
				},
			},
		},
		{
			name: "undefined preferred language, ok",
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
								DisallowedLanguage,
								domain.GenderUnspecified,
								"email",
								true,
							),
						),
					),
					expectPush(
						newProfileChangedEvent(context.Background(),
							"user1", "org1",
							"firstname2",
							"lastname2",
							"nickname2",
							"displayname2",
							language.Und,
							domain.GenderMale,
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				address: &domain.Profile{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "user1",
					},
					FirstName:   "firstname2",
					LastName:    "lastname2",
					NickName:    "nickname2",
					DisplayName: "displayname2",
					Gender:      domain.GenderMale,
				},
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.Profile{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "user1",
						ResourceOwner: "org1",
					},
					FirstName:         "firstname2",
					LastName:          "lastname2",
					NickName:          "nickname2",
					DisplayName:       "displayname2",
					PreferredLanguage: language.Und,
					Gender:            domain.GenderMale,
				},
			},
		}, {
			name: "unsupported preferred language, ok",
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
								DisallowedLanguage,
								domain.GenderUnspecified,
								"email",
								true,
							),
						),
					),
					expectPush(
						newProfileChangedEvent(context.Background(),
							"user1", "org1",
							"firstname2",
							"lastname2",
							"nickname2",
							"displayname2",
							UnsupportedLanguage,
							domain.GenderMale,
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				address: &domain.Profile{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "user1",
					},
					FirstName:         "firstname2",
					LastName:          "lastname2",
					NickName:          "nickname2",
					DisplayName:       "displayname2",
					PreferredLanguage: UnsupportedLanguage,
					Gender:            domain.GenderMale,
				},
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.Profile{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "user1",
						ResourceOwner: "org1",
					},
					FirstName:         "firstname2",
					LastName:          "lastname2",
					NickName:          "nickname2",
					DisplayName:       "displayname2",
					PreferredLanguage: UnsupportedLanguage,
					Gender:            domain.GenderMale,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.ChangeHumanProfile(tt.args.ctx, tt.args.address)
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

func newProfileChangedEvent(ctx context.Context, userID, resourceOwner, fistName, lastName, nickName, displayName string, lang language.Tag, gender domain.Gender) *user.HumanProfileChangedEvent {
	event, _ := user.NewHumanProfileChangedEvent(ctx,
		&user.NewAggregate(userID, resourceOwner).Aggregate,
		[]user.ProfileChanges{
			user.ChangeFirstName(fistName),
			user.ChangeLastName(lastName),
			user.ChangeNickName(nickName),
			user.ChangeDisplayName(displayName),
			user.ChangePreferredLanguage(lang),
			user.ChangeGender(gender),
		},
	)
	return event
}
