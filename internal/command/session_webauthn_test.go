package command

import (
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestSessionCommands_getHumanWebAuthNTokens(t *testing.T) {
	userAggr := &user.NewAggregate("user1", "org1").Aggregate

	type fields struct {
		eventstore        *eventstore.Eventstore
		sessionWriteModel *SessionWriteModel
	}
	type args struct {
		userVerification domain.UserVerificationRequirement
	}
	type res struct {
		want *humanWebAuthNTokens
		err  error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "missing UID",
			fields: fields{
				eventstore:        &eventstore.Eventstore{},
				sessionWriteModel: &SessionWriteModel{},
			},
			args: args{
				domain.UserVerificationRequirementDiscouraged,
			},
			res: res{
				want: nil,
				err:  zerrors.ThrowPreconditionFailed(nil, "COMMAND-eeR2e", "Errors.User.UserIDMissing"),
			},
		},
		{
			name: "passwordless filter error",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								userAggr,
								"", "", "", "", "", language.Georgian,
								domain.GenderDiverse, "", true,
							),
						),
					),
					expectFilterError(io.ErrClosedPipe),
				),
				sessionWriteModel: &SessionWriteModel{
					UserID: "user1",
				},
			},
			args: args{
				domain.UserVerificationRequirementDiscouraged,
			},
			res: res{
				want: nil,
				err:  io.ErrClosedPipe,
			},
		},
		{
			name: "ok, discouraged, u2f",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								userAggr,
								"", "", "", "", "", language.Georgian,
								domain.GenderDiverse, "", true,
							),
						),
					),
					expectFilter(eventFromEventPusher(
						user.NewHumanWebAuthNAddedEvent(eventstore.NewBaseEventForPush(
							context.Background(), &org.NewAggregate("org1").Aggregate, user.HumanU2FTokenAddedType,
						), "111", "challenge", "rpID"),
					)),
				),
				sessionWriteModel: &SessionWriteModel{
					UserID: "user1",
				},
			},
			args: args{
				domain.UserVerificationRequirementDiscouraged,
			},
			res: res{
				want: &humanWebAuthNTokens{
					human: &domain.Human{
						ObjectRoot: models.ObjectRoot{
							AggregateID:   "user1",
							ResourceOwner: "org1",
						},
						State: domain.UserStateActive,
						Profile: &domain.Profile{
							PreferredLanguage: language.Georgian,
							Gender:            domain.GenderDiverse,
						},
						Email: &domain.Email{},
					},
					tokens: []*domain.WebAuthNToken{{
						ObjectRoot: models.ObjectRoot{
							AggregateID: "org1",
						},
						WebAuthNTokenID: "111",
						State:           domain.MFAStateNotReady,
						Challenge:       "challenge",
						RPID:            "rpID",
					}},
				},
				err: nil,
			},
		},
		{
			name: "ok, preferred, u2f",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								userAggr,
								"", "", "", "", "", language.Georgian,
								domain.GenderDiverse, "", true,
							),
						),
					),
					expectFilter(eventFromEventPusher(
						user.NewHumanWebAuthNAddedEvent(eventstore.NewBaseEventForPush(
							context.Background(), &org.NewAggregate("org1").Aggregate, user.HumanU2FTokenAddedType,
						), "111", "challenge", "rpID"),
					)),
				),
				sessionWriteModel: &SessionWriteModel{
					UserID: "user1",
				},
			},
			args: args{
				domain.UserVerificationRequirementPreferred,
			},
			res: res{
				want: &humanWebAuthNTokens{
					human: &domain.Human{
						ObjectRoot: models.ObjectRoot{
							AggregateID:   "user1",
							ResourceOwner: "org1",
						},
						State: domain.UserStateActive,
						Profile: &domain.Profile{
							PreferredLanguage: language.Georgian,
							Gender:            domain.GenderDiverse,
						},
						Email: &domain.Email{},
					},
					tokens: []*domain.WebAuthNToken{{
						ObjectRoot: models.ObjectRoot{
							AggregateID: "org1",
						},
						WebAuthNTokenID: "111",
						State:           domain.MFAStateNotReady,
						Challenge:       "challenge",
						RPID:            "rpID",
					}},
				},
				err: nil,
			},
		},
		{
			name: "ok, required, u2f",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								userAggr,
								"", "", "", "", "", language.Georgian,
								domain.GenderDiverse, "", true,
							),
						),
					),
					expectFilter(eventFromEventPusher(
						user.NewHumanWebAuthNAddedEvent(eventstore.NewBaseEventForPush(
							context.Background(), &org.NewAggregate("org1").Aggregate, user.HumanPasswordlessTokenAddedType,
						), "111", "challenge", "rpID"),
					)),
				),
				sessionWriteModel: &SessionWriteModel{
					UserID: "user1",
				},
			},
			args: args{
				domain.UserVerificationRequirementRequired,
			},
			res: res{
				want: &humanWebAuthNTokens{
					human: &domain.Human{
						ObjectRoot: models.ObjectRoot{
							AggregateID:   "user1",
							ResourceOwner: "org1",
						},
						State: domain.UserStateActive,
						Profile: &domain.Profile{
							PreferredLanguage: language.Georgian,
							Gender:            domain.GenderDiverse,
						},
						Email: &domain.Email{},
					},
					tokens: []*domain.WebAuthNToken{{
						ObjectRoot: models.ObjectRoot{
							AggregateID: "org1",
						},
						WebAuthNTokenID: "111",
						State:           domain.MFAStateNotReady,
						Challenge:       "challenge",
						RPID:            "rpID",
					}},
				},
				err: nil,
			},
		},
	}
	for _, tt := range tests {
		s := &SessionCommands{
			eventstore:        tt.fields.eventstore,
			sessionWriteModel: tt.fields.sessionWriteModel,
		}
		got, err := s.getHumanWebAuthNTokens(context.Background(), tt.args.userVerification)
		require.ErrorIs(t, err, tt.res.err)
		assert.Equal(t, tt.res.want, got)
	}
}
