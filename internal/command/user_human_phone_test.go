package command

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommandSide_ChangeHumanPhone(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx             context.Context
		email           *domain.Phone
		resourceOwner   string
		secretGenerator crypto.Generator
	}
	type res struct {
		want *domain.Phone
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "invalid phone, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: context.Background(),
				email: &domain.Phone{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "user1",
					},
				},
				resourceOwner: "org1",
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
				ctx: context.Background(),
				email: &domain.Phone{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "user1",
					},
					PhoneNumber: "+41711234567",
				},
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "phone not changed, precondition error",
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
						eventFromEventPusher(
							user.NewHumanPhoneChangedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"+41711234567",
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				email: &domain.Phone{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "user1",
					},
					PhoneNumber: "+41711234567",
				},
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "verified phone changed, ok",
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
						eventFromEventPusher(
							user.NewHumanPhoneChangedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"+41711234567",
							),
						),
					),
					expectPush(
						user.NewHumanPhoneChangedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							"+41719876543",
						),
						user.NewHumanPhoneVerifiedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				email: &domain.Phone{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "user1",
					},
					PhoneNumber:     "+41719876543",
					IsPhoneVerified: true,
				},
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.Phone{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "user1",
						ResourceOwner: "org1",
					},
					PhoneNumber:     "+41719876543",
					IsPhoneVerified: true,
				},
			},
		},
		{
			name: "phone changed to verified, ok",
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
						eventFromEventPusher(
							user.NewHumanPhoneChangedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"+41711234567",
							),
						),
					),
					expectPush(
						user.NewHumanPhoneVerifiedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				email: &domain.Phone{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "user1",
					},
					PhoneNumber:     "+41711234567",
					IsPhoneVerified: true,
				},
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.Phone{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "user1",
						ResourceOwner: "org1",
					},
					PhoneNumber:     "+41711234567",
					IsPhoneVerified: true,
				},
			},
		},
		{
			name: "phone changed to verified, ok",
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
						eventFromEventPusher(
							user.NewHumanPhoneChangedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"+41711234567",
							),
						),
					),
					expectPush(
						user.NewHumanPhoneVerifiedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				email: &domain.Phone{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "user1",
					},
					PhoneNumber:     "+41711234567",
					IsPhoneVerified: true,
				},
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.Phone{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "user1",
						ResourceOwner: "org1",
					},
					PhoneNumber:     "+41711234567",
					IsPhoneVerified: true,
				},
			},
		},
		{
			name: "phone changed with code, ok",
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
						user.NewHumanPhoneChangedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							"+41711234567",
						),
						user.NewHumanPhoneCodeAddedEvent(context.Background(),
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
				),
			},
			args: args{
				ctx: context.Background(),
				email: &domain.Phone{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "user1",
					},
					PhoneNumber: "+41711234567",
				},
				resourceOwner:   "org1",
				secretGenerator: GetMockSecretGenerator(t),
			},
			res: res{
				want: &domain.Phone{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "user1",
						ResourceOwner: "org1",
					},
					PhoneNumber: "+41711234567",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.ChangeHumanPhone(tt.args.ctx, tt.args.email, tt.args.resourceOwner, tt.args.secretGenerator)
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

func TestCommandSide_VerifyHumanPhone(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx             context.Context
		userID          string
		code            string
		resourceOwner   string
		secretGenerator crypto.Generator
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
			name: "userid missing, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:           context.Background(),
				code:          "aa",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "code missing, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				resourceOwner: "org1",
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
				code:          "aa",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "code not existing, not found error",
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
				ctx:           context.Background(),
				userID:        "user1",
				code:          "aa",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "invalid code, invalid argument error",
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
						eventFromEventPusher(
							user.NewHumanPhoneChangedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"+411234567",
							),
						),
						eventFromEventPusher(
							user.NewHumanPhoneCodeAddedEvent(context.Background(),
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
					),
					expectPush(
						user.NewHumanPhoneVerificationFailedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
						),
					),
				),
			},
			args: args{
				ctx:             context.Background(),
				userID:          "user1",
				code:            "test",
				resourceOwner:   "org1",
				secretGenerator: GetMockSecretGenerator(t),
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "valid code, ok",
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
						eventFromEventPusher(
							user.NewHumanPhoneChangedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"+411234567",
							),
						),
						eventFromEventPusherWithCreationDateNow(
							user.NewHumanPhoneCodeAddedEvent(context.Background(),
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
					),
					expectPush(
						user.NewHumanPhoneVerifiedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
						),
					),
				),
			},
			args: args{
				ctx:             context.Background(),
				userID:          "user1",
				code:            "a",
				resourceOwner:   "org1",
				secretGenerator: GetMockSecretGenerator(t),
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
			got, err := r.VerifyHumanPhone(tt.args.ctx, tt.args.userID, tt.args.code, tt.args.resourceOwner, tt.args.secretGenerator)
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

func TestCommandSide_CreateVerificationCodeHumanPhone(t *testing.T) {
	type fields struct {
		eventstore     *eventstore.Eventstore
		userEncryption crypto.EncryptionAlgorithm
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
			name: "userid missing, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
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
			name: "phone already verified, precondition error",
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
						eventFromEventPusher(
							user.NewHumanPhoneChangedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"+411234567",
							),
						),
						eventFromEventPusher(
							user.NewHumanPhoneVerifiedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
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
			name: "new code, ok",
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
						eventFromEventPusher(
							user.NewHumanPhoneChangedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"+411234567",
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							instance.NewSecretGeneratorAddedEvent(context.Background(),
								&instance.NewAggregate("instanceID").Aggregate,
								domain.SecretGeneratorTypeVerifyPhoneCode,
								8,
								time.Hour,
								true,
								true,
								true,
								true,
							)),
					),
					expectPush(
						user.NewHumanPhoneCodeAddedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("12345678"),
							},
							time.Hour*1,
						),
					),
				),
				userEncryption: crypto.CreateMockEncryptionAlgWithCode(gomock.NewController(t), "12345678"),
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
				eventstore:     tt.fields.eventstore,
				userEncryption: tt.fields.userEncryption,
			}
			got, err := r.CreateHumanPhoneVerificationCode(tt.args.ctx, tt.args.userID, tt.args.resourceOwner)
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

func TestCommandSide_PhoneVerificationCodeSent(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		userID        string
		resourceOwner string
	}
	type res struct {
		err func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "userid missing, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
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
			name: "code sent, ok",
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
						eventFromEventPusher(
							user.NewHumanPhoneChangedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"+411234567",
							),
						),
					),
					expectPush(
						user.NewHumanPhoneCodeSentEvent(context.Background(),
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
			res: res{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			err := r.HumanPhoneVerificationCodeSent(tt.args.ctx, tt.args.resourceOwner, tt.args.userID)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestCommandSide_RemoveHumanPhone(t *testing.T) {
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
			name: "userid missing, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
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
			name: "phone not existing, precondition error",
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
				ctx:           context.Background(),
				userID:        "user1",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "remove phone, ok",
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
						eventFromEventPusherWithCreationDateNow(
							user.NewHumanPhoneChangedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"+411234567",
							),
						),
					),
					expectPush(
						user.NewHumanPhoneRemovedEvent(context.Background(),
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
			got, err := r.RemoveHumanPhone(tt.args.ctx, tt.args.userID, tt.args.resourceOwner)
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
