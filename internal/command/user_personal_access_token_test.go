package command

import (
	"context"
	"encoding/base64"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/eventstore/repository"

	"github.com/zitadel/zitadel/internal/eventstore/v1/models"

	id_mock "github.com/zitadel/zitadel/internal/id/mock"

	"github.com/zitadel/zitadel/internal/repository/user"

	caos_errs "github.com/zitadel/zitadel/internal/errors"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id"
)

func TestCommands_AddPersonalAccessToken(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		idGenerator  id.Generator
		keyAlgorithm crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx context.Context
		pat *PersonalAccessToken
	}
	type res struct {
		want  *domain.ObjectDetails
		token string
		err   func(error) bool
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
				pat: &PersonalAccessToken{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "user1",
						ResourceOwner: "org1",
					},
					Scopes:         []string{"openid"},
					ExpirationDate: time.Time{},
				},
			},
			res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			"user type not allowed, error",
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
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "token1"),
			},
			args{
				ctx: context.Background(),
				pat: &PersonalAccessToken{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "user1",
						ResourceOwner: "org1",
					},
					Scopes:          []string{"openid"},
					ExpirationDate:  time.Time{},
					AllowedUserType: domain.UserTypeHuman,
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
				pat: &PersonalAccessToken{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "user1",
						ResourceOwner: "org1",
					},
					Scopes:         []string{"openid"},
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
				pat: &PersonalAccessToken{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "",
						ResourceOwner: "org1",
					},
					Scopes:         []string{"openid"},
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
				pat: &PersonalAccessToken{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "user1",
						ResourceOwner: "",
					},
					Scopes:         []string{"openid"},
					ExpirationDate: time.Time{},
				},
			},
			res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			"token added",
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
								user.NewPersonalAccessTokenAddedEvent(context.Background(),
									&user.NewAggregate("user1", "org1").Aggregate,
									"token1",
									time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC),
									[]string{"openid"},
								),
							),
						},
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "token1"),
				keyAlgorithm: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args{
				ctx: context.Background(),
				pat: &PersonalAccessToken{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "user1",
						ResourceOwner: "org1",
					},
					Scopes:          []string{"openid"},
					ExpirationDate:  time.Time{},
					AllowedUserType: domain.UserTypeMachine,
				},
			},
			res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
				token: base64.RawURLEncoding.EncodeToString([]byte("token1:user1")),
			},
		},
		{
			"token added with ID",
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
								user.NewPersonalAccessTokenAddedEvent(context.Background(),
									&user.NewAggregate("user1", "org1").Aggregate,
									"token1",
									time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC),
									[]string{"openid"},
								),
							),
						},
					),
				),
				keyAlgorithm: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args{
				ctx: context.Background(),
				pat: &PersonalAccessToken{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "user1",
						ResourceOwner: "org1",
					},
					TokenID:         "token1",
					Scopes:          []string{"openid"},
					ExpirationDate:  time.Time{},
					AllowedUserType: domain.UserTypeMachine,
				},
			},
			res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
				token: base64.RawURLEncoding.EncodeToString([]byte("token1:user1")),
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
			got, err := c.AddPersonalAccessToken(tt.args.ctx, tt.args.pat)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, got)
				assert.Equal(t, tt.res.token, tt.args.pat.Token)
			}
		})
	}
}

func TestCommands_RemovePersonalAccessToken(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx context.Context
		pat *PersonalAccessToken
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
			"token does not exist, error",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
				),
			},
			args{
				ctx: context.Background(),
				pat: &PersonalAccessToken{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "user1",
						ResourceOwner: "org1",
					},
					TokenID: "token1",
				},
			},
			res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			"remove token, ok",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							user.NewPersonalAccessTokenAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"token1",
								time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC),
								[]string{"openid"},
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								user.NewPersonalAccessTokenRemovedEvent(context.Background(),
									&user.NewAggregate("user1", "org1").Aggregate,
									"token1",
								),
							),
						},
					),
				),
			},
			args{
				ctx: context.Background(),
				pat: &PersonalAccessToken{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "user1",
						ResourceOwner: "org1",
					},
					TokenID: "token1",
				},
			},
			res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := c.RemovePersonalAccessToken(tt.args.ctx, tt.args.pat)
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
