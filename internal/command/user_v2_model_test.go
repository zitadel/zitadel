package command

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user"
)

func TestCommandSide_userExistsWriteModel(t *testing.T) {
	type fields struct {
		eventstore func(t *testing.T) *eventstore.Eventstore
	}
	type args struct {
		ctx    context.Context
		userID string
	}
	type res struct {
		want *UserV2WriteModel
		err  func(error) bool
	}

	userAgg := user.NewAggregate("user1", "org1")

	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "user added",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				userID: "user1",
			},
			res: res{
				want: &UserV2WriteModel{
					HumanWriteModel:   true,
					MachineWriteModel: true,
					WriteModel: eventstore.WriteModel{
						AggregateID:       "user1",
						Events:            []eventstore.Event{},
						ProcessedSequence: 0,
						ResourceOwner:     "org1",
					},
					UserName:               "username",
					FirstName:              "firstname",
					LastName:               "lastname",
					DisplayName:            "firstname lastname",
					PreferredLanguage:      language.English,
					PasswordEncodedHash:    "$plain$x$password",
					PasswordChangeRequired: true,
					Email:                  "email@test.ch",
					IsEmailVerified:        false,
					InitCode:               nil,
					InitCodeCreationDate:   time.Time{},
					InitCodeExpiry:         0,
					UserState:              domain.UserStateActive,
				},
			},
		},
		{
			name: "user registered",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newRegisterHumanEvent("username", "$plain$x$password", true, true, "", language.English),
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				userID: "user1",
			},
			res: res{
				want: &UserV2WriteModel{
					HumanWriteModel:   true,
					MachineWriteModel: true,
					WriteModel: eventstore.WriteModel{
						AggregateID:   "user1",
						Events:        []eventstore.Event{},
						ResourceOwner: "org1",
					},
					UserName:               "username",
					FirstName:              "firstname",
					LastName:               "lastname",
					DisplayName:            "firstname lastname",
					PasswordEncodedHash:    "$plain$x$password",
					PasswordChangeRequired: true,
					Email:                  "email@test.ch",
					IsEmailVerified:        false,
					InitCode:               nil,
					InitCodeCreationDate:   time.Time{},
					InitCodeExpiry:         0,
					PreferredLanguage:      language.English,
					UserState:              domain.UserStateActive,
				},
			},
		},
		{
			name: "user machine added",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddMachineEvent(true, domain.OIDCTokenTypeBearer),
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				userID: "user1",
			},
			res: res{
				want: &UserV2WriteModel{
					HumanWriteModel:   true,
					MachineWriteModel: true,
					WriteModel: eventstore.WriteModel{
						AggregateID:       "user1",
						Events:            []eventstore.Event{},
						ProcessedSequence: 0,
						ResourceOwner:     "org1",
					},
					UserName:        "username",
					Name:            "name",
					Description:     "description",
					AccessTokenType: domain.OIDCTokenTypeBearer,
					UserState:       domain.UserStateActive,
				},
			},
		},
		{
			name: "user added with init code",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newRegisterHumanEvent("username", "$plain$x$password", true, true, "", language.English),
						),
						eventFromEventPusher(
							user.NewHumanInitialCodeAddedEvent(context.Background(),
								&userAgg.Aggregate,
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
				),
			},
			args: args{
				ctx:    context.Background(),
				userID: "user1",
			},
			res: res{
				want: &UserV2WriteModel{
					HumanWriteModel:   true,
					MachineWriteModel: true,
					WriteModel: eventstore.WriteModel{
						AggregateID:       "user1",
						Events:            []eventstore.Event{},
						ProcessedSequence: 0,
						ResourceOwner:     "org1",
					},
					UserName:               "username",
					FirstName:              "firstname",
					LastName:               "lastname",
					DisplayName:            "firstname lastname",
					PreferredLanguage:      language.English,
					PasswordEncodedHash:    "$plain$x$password",
					PasswordChangeRequired: true,
					Email:                  "email@test.ch",
					IsEmailVerified:        false,
					InitCode: &crypto.CryptoValue{
						CryptoType: crypto.TypeEncryption,
						Algorithm:  "enc",
						KeyID:      "id",
						Crypted:    []byte("a"),
					},
					InitCodeCreationDate: time.Time{},
					InitCodeExpiry:       time.Hour * 1,
					UserState:            domain.UserStateInitial,
				},
			},
		},
		{
			name: "user added with initialized",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newRegisterHumanEvent("username", "$plain$x$password", true, true, "", language.English),
						),
						eventFromEventPusher(
							user.NewHumanInitialCodeAddedEvent(context.Background(),
								&userAgg.Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("a"),
								},
								time.Hour*1,
							),
						),
						eventFromEventPusher(
							user.NewHumanInitializedCheckSucceededEvent(context.Background(),
								&userAgg.Aggregate,
							),
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				userID: "user1",
			},
			res: res{
				want: &UserV2WriteModel{
					HumanWriteModel:   true,
					MachineWriteModel: true,
					WriteModel: eventstore.WriteModel{
						AggregateID:       "user1",
						Events:            []eventstore.Event{},
						ProcessedSequence: 0,
						ResourceOwner:     "org1",
					},
					UserName:               "username",
					FirstName:              "firstname",
					LastName:               "lastname",
					DisplayName:            "firstname lastname",
					PreferredLanguage:      language.English,
					PasswordEncodedHash:    "$plain$x$password",
					PasswordChangeRequired: true,
					Email:                  "email@test.ch",
					IsEmailVerified:        false,
					UserState:              domain.UserStateActive,
				},
			},
		},
		{
			name: "user added with initialized failed",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newRegisterHumanEvent("username", "$plain$x$password", true, true, "", language.English),
						),
						eventFromEventPusher(
							user.NewHumanInitialCodeAddedEvent(context.Background(),
								&userAgg.Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("a"),
								},
								time.Hour*1,
							),
						),
						eventFromEventPusher(
							user.NewHumanInitializedCheckFailedEvent(context.Background(),
								&userAgg.Aggregate,
							),
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				userID: "user1",
			},
			res: res{
				want: &UserV2WriteModel{
					HumanWriteModel:   true,
					MachineWriteModel: true,
					WriteModel: eventstore.WriteModel{
						AggregateID:       "user1",
						Events:            []eventstore.Event{},
						ProcessedSequence: 0,
						ResourceOwner:     "org1",
					},
					UserName:               "username",
					FirstName:              "firstname",
					LastName:               "lastname",
					DisplayName:            "firstname lastname",
					PreferredLanguage:      language.English,
					PasswordEncodedHash:    "$plain$x$password",
					PasswordChangeRequired: true,
					Email:                  "email@test.ch",
					IsEmailVerified:        false,
					InitCode: &crypto.CryptoValue{
						CryptoType: crypto.TypeEncryption,
						Algorithm:  "enc",
						KeyID:      "id",
						Crypted:    []byte("a"),
					},
					InitCodeCreationDate: time.Time{},
					InitCodeExpiry:       time.Hour * 1,
					InitCheckFailedCount: 1,
					UserState:            domain.UserStateInitial,
				},
			},
		},
		{
			name: "user added with username changed",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newRegisterHumanEvent("username", "$plain$x$password", true, true, "", language.English),
						),
						eventFromEventPusher(
							user.NewUsernameChangedEvent(context.Background(),
								&userAgg.Aggregate,
								"username",
								"changed",
								true,
							),
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				userID: "user1",
			},
			res: res{
				want: &UserV2WriteModel{
					HumanWriteModel:   true,
					MachineWriteModel: true,
					WriteModel: eventstore.WriteModel{
						AggregateID:       "user1",
						Events:            []eventstore.Event{},
						ProcessedSequence: 0,
						ResourceOwner:     "org1",
					},
					UserName:               "changed",
					FirstName:              "firstname",
					LastName:               "lastname",
					DisplayName:            "firstname lastname",
					PreferredLanguage:      language.English,
					PasswordEncodedHash:    "$plain$x$password",
					PasswordChangeRequired: true,
					Email:                  "email@test.ch",
					IsEmailVerified:        false,
					UserState:              domain.UserStateActive,
				},
			},
		},
		{
			name: "user removed",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						),
						eventFromEventPusher(
							user.NewUserRemovedEvent(context.Background(),
								&userAgg.Aggregate,
								"username",
								[]*domain.UserIDPLink{},
								true,
							),
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				userID: "user1",
			},
			res: res{
				want: &UserV2WriteModel{
					HumanWriteModel:   true,
					MachineWriteModel: true,
					WriteModel: eventstore.WriteModel{
						AggregateID:       "user1",
						Events:            []eventstore.Event{},
						ProcessedSequence: 0,
						ResourceOwner:     "org1",
					},
					UserName:               "username",
					FirstName:              "firstname",
					LastName:               "lastname",
					DisplayName:            "firstname lastname",
					PreferredLanguage:      language.English,
					PasswordEncodedHash:    "$plain$x$password",
					PasswordChangeRequired: true,
					Email:                  "email@test.ch",
					IsEmailVerified:        false,
					UserState:              domain.UserStateDeleted,
				},
			},
		},
		{
			name: "user machine removed",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddMachineEvent(true, domain.OIDCTokenTypeBearer),
						),
						eventFromEventPusher(
							user.NewUserRemovedEvent(context.Background(),
								&userAgg.Aggregate,
								"username",
								[]*domain.UserIDPLink{},
								true,
							),
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				userID: "user1",
			},
			res: res{
				want: &UserV2WriteModel{
					HumanWriteModel:   true,
					MachineWriteModel: true,
					WriteModel: eventstore.WriteModel{
						AggregateID:       "user1",
						Events:            []eventstore.Event{},
						ProcessedSequence: 0,
						ResourceOwner:     "org1",
					},
					UserName:        "username",
					Name:            "name",
					Description:     "description",
					AccessTokenType: domain.OIDCTokenTypeBearer,
					UserState:       domain.UserStateDeleted,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore(t),
			}
			wm, err := r.userExistsWriteModel(tt.args.ctx, tt.args.userID)
			if tt.res.err == nil {
				if !assert.NoError(t, err) {
					t.FailNow()
				}
			} else if !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
				return
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, wm)
			}
		})
	}
}

func TestCommandSide_userHumanWriteModel_profile(t *testing.T) {
	type fields struct {
		eventstore func(t *testing.T) *eventstore.Eventstore
	}
	type args struct {
		ctx    context.Context
		userID string
	}
	type res struct {
		want *UserV2WriteModel
		err  func(error) bool
	}

	userAgg := user.NewAggregate("user1", "org1")
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "user added with profile changed",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newRegisterHumanEvent("username", "$plain$x$password", true, true, "", language.English),
						),
						eventFromEventPusher(
							func() eventstore.Command {
								cmd, _ := user.NewHumanProfileChangedEvent(context.Background(),
									&userAgg.Aggregate,
									[]user.ProfileChanges{
										user.ChangeFirstName("changedfn"),
										user.ChangeLastName("changedln"),
										user.ChangeNickName("changednn"),
										user.ChangeDisplayName("changeddn"),
										user.ChangePreferredLanguage(language.Afrikaans),
										user.ChangeGender(domain.GenderDiverse),
									},
								)
								return cmd
							}(),
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				userID: "user1",
			},
			res: res{
				want: &UserV2WriteModel{
					HumanWriteModel:   true,
					ProfileWriteModel: true,
					StateWriteModel:   true,
					WriteModel: eventstore.WriteModel{
						AggregateID:       "user1",
						Events:            []eventstore.Event{},
						ProcessedSequence: 0,
						ResourceOwner:     "org1",
					},
					UserName:               "username",
					FirstName:              "changedfn",
					LastName:               "changedln",
					DisplayName:            "changeddn",
					NickName:               "changednn",
					PreferredLanguage:      language.Afrikaans,
					Gender:                 domain.GenderDiverse,
					PasswordEncodedHash:    "$plain$x$password",
					PasswordChangeRequired: true,
					Email:                  "email@test.ch",
					IsEmailVerified:        false,
					UserState:              domain.UserStateActive,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore(t),
			}
			wm, err := r.userHumanWriteModel(tt.args.ctx, tt.args.userID, true, false, false, false, false, false)
			if tt.res.err == nil {
				if !assert.NoError(t, err) {
					t.FailNow()
				}
			} else if !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
				return
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, wm)
			}
		})
	}
}

func TestCommandSide_userHumanWriteModel_email(t *testing.T) {
	type fields struct {
		eventstore func(t *testing.T) *eventstore.Eventstore
	}
	type args struct {
		ctx    context.Context
		userID string
	}
	type res struct {
		want *UserV2WriteModel
		err  func(error) bool
	}

	userAgg := user.NewAggregate("user1", "org1")

	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "user added email changed",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						),
						eventFromEventPusher(
							user.NewHumanEmailChangedEvent(context.Background(),
								&userAgg.Aggregate,
								"changed@test.com",
							),
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				userID: "user1",
			},
			res: res{
				want: &UserV2WriteModel{
					HumanWriteModel: true,
					EmailWriteModel: true,
					StateWriteModel: true,
					WriteModel: eventstore.WriteModel{
						AggregateID:       "user1",
						Events:            []eventstore.Event{},
						ProcessedSequence: 0,
						ResourceOwner:     "org1",
					},
					UserName:               "username",
					FirstName:              "firstname",
					LastName:               "lastname",
					DisplayName:            "firstname lastname",
					PreferredLanguage:      language.English,
					PasswordEncodedHash:    "$plain$x$password",
					PasswordChangeRequired: true,
					Email:                  "changed@test.com",
					IsEmailVerified:        false,
					InitCode:               nil,
					InitCodeCreationDate:   time.Time{},
					InitCodeExpiry:         0,
					UserState:              domain.UserStateActive,
				},
			},
		},
		{
			name: "user added with email code",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newRegisterHumanEvent("username", "$plain$x$password", true, true, "", language.English),
						),
						eventFromEventPusher(
							user.NewHumanEmailCodeAddedEventV2(context.Background(),
								&userAgg.Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("a"),
								},
								time.Hour*1,
								"",
								false,
							),
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				userID: "user1",
			},
			res: res{
				want: &UserV2WriteModel{
					HumanWriteModel: true,
					EmailWriteModel: true,
					StateWriteModel: true,
					WriteModel: eventstore.WriteModel{
						AggregateID:       "user1",
						Events:            []eventstore.Event{},
						ProcessedSequence: 0,
						ResourceOwner:     "org1",
					},
					UserName:               "username",
					FirstName:              "firstname",
					LastName:               "lastname",
					DisplayName:            "firstname lastname",
					PreferredLanguage:      language.English,
					PasswordEncodedHash:    "$plain$x$password",
					PasswordChangeRequired: true,
					Email:                  "email@test.ch",
					IsEmailVerified:        false,
					EmailCode: &crypto.CryptoValue{
						CryptoType: crypto.TypeEncryption,
						Algorithm:  "enc",
						KeyID:      "id",
						Crypted:    []byte("a"),
					},
					EmailCodeCreationDate: time.Time{},
					EmailCodeExpiry:       time.Hour * 1,
					UserState:             domain.UserStateActive,
				},
			},
		},
		{
			name: "user added with email code verified",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newRegisterHumanEvent("username", "$plain$x$password", true, true, "", language.English),
						),
						eventFromEventPusher(
							user.NewHumanEmailCodeAddedEventV2(context.Background(),
								&userAgg.Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("a"),
								},
								time.Hour*1,
								"",
								false,
							),
						),
						eventFromEventPusher(
							user.NewHumanEmailVerifiedEvent(context.Background(),
								&userAgg.Aggregate,
							),
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				userID: "user1",
			},
			res: res{
				want: &UserV2WriteModel{
					HumanWriteModel: true,
					EmailWriteModel: true,
					StateWriteModel: true,
					WriteModel: eventstore.WriteModel{
						AggregateID:       "user1",
						Events:            []eventstore.Event{},
						ProcessedSequence: 0,
						ResourceOwner:     "org1",
					},
					UserName:               "username",
					FirstName:              "firstname",
					LastName:               "lastname",
					DisplayName:            "firstname lastname",
					PreferredLanguage:      language.English,
					PasswordEncodedHash:    "$plain$x$password",
					PasswordChangeRequired: true,
					Email:                  "email@test.ch",
					IsEmailVerified:        true,
					UserState:              domain.UserStateActive,
				},
			},
		},
		{
			name: "user added with email code verified failed",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newRegisterHumanEvent("username", "$plain$x$password", true, true, "", language.English),
						),
						eventFromEventPusher(
							user.NewHumanEmailCodeAddedEventV2(context.Background(),
								&userAgg.Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("a"),
								},
								time.Hour*1,
								"",
								false,
							),
						),
						eventFromEventPusher(
							user.NewHumanEmailVerificationFailedEvent(context.Background(),
								&userAgg.Aggregate,
							),
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				userID: "user1",
			},
			res: res{
				want: &UserV2WriteModel{
					HumanWriteModel: true,
					EmailWriteModel: true,
					StateWriteModel: true,
					WriteModel: eventstore.WriteModel{
						AggregateID:       "user1",
						Events:            []eventstore.Event{},
						ProcessedSequence: 0,
						ResourceOwner:     "org1",
					},
					UserName:               "username",
					FirstName:              "firstname",
					LastName:               "lastname",
					DisplayName:            "firstname lastname",
					PreferredLanguage:      language.English,
					PasswordEncodedHash:    "$plain$x$password",
					PasswordChangeRequired: true,
					Email:                  "email@test.ch",
					IsEmailVerified:        false,
					EmailCode: &crypto.CryptoValue{
						CryptoType: crypto.TypeEncryption,
						Algorithm:  "enc",
						KeyID:      "id",
						Crypted:    []byte("a"),
					},
					EmailCodeCreationDate: time.Time{},
					EmailCodeExpiry:       time.Hour * 1,
					EmailCheckFailedCount: 1,
					UserState:             domain.UserStateActive,
				},
			},
		},
		{
			name: "user added with email code verified, then changed",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newRegisterHumanEvent("username", "$plain$x$password", true, true, "", language.English),
						),
						eventFromEventPusher(
							user.NewHumanEmailCodeAddedEventV2(context.Background(),
								&userAgg.Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("a"),
								},
								time.Hour*1,
								"",
								false,
							),
						),
						eventFromEventPusher(
							user.NewHumanEmailVerifiedEvent(context.Background(),
								&userAgg.Aggregate,
							),
						),
						eventFromEventPusher(
							user.NewHumanEmailChangedEvent(context.Background(),
								&userAgg.Aggregate,
								"changed@test.com",
							),
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				userID: "user1",
			},
			res: res{
				want: &UserV2WriteModel{
					HumanWriteModel: true,
					EmailWriteModel: true,
					StateWriteModel: true,
					WriteModel: eventstore.WriteModel{
						AggregateID:       "user1",
						Events:            []eventstore.Event{},
						ProcessedSequence: 0,
						ResourceOwner:     "org1",
					},
					UserName:               "username",
					FirstName:              "firstname",
					LastName:               "lastname",
					DisplayName:            "firstname lastname",
					PreferredLanguage:      language.English,
					PasswordEncodedHash:    "$plain$x$password",
					PasswordChangeRequired: true,
					Email:                  "changed@test.com",
					IsEmailVerified:        false,
					UserState:              domain.UserStateActive,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore(t),
			}
			wm, err := r.userHumanWriteModel(tt.args.ctx, tt.args.userID, false, true, false, false, false, false)
			if tt.res.err == nil {
				if !assert.NoError(t, err) {
					t.FailNow()
				}
			} else if !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
				return
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, wm)
			}
		})
	}
}

func TestCommandSide_userHumanWriteModel_phone(t *testing.T) {
	type fields struct {
		eventstore func(t *testing.T) *eventstore.Eventstore
	}
	type args struct {
		ctx    context.Context
		userID string
	}
	type res struct {
		want *UserV2WriteModel
		err  func(error) bool
	}

	userAgg := user.NewAggregate("user1", "org1")

	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "user added phone changed",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						),
						eventFromEventPusher(
							user.NewHumanPhoneChangedEvent(context.Background(),
								&userAgg.Aggregate,
								"+41791234567",
							),
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				userID: "user1",
			},
			res: res{
				want: &UserV2WriteModel{
					HumanWriteModel: true,
					PhoneWriteModel: true,
					StateWriteModel: true,
					WriteModel: eventstore.WriteModel{
						AggregateID:       "user1",
						Events:            []eventstore.Event{},
						ProcessedSequence: 0,
						ResourceOwner:     "org1",
					},
					UserName:               "username",
					FirstName:              "firstname",
					LastName:               "lastname",
					DisplayName:            "firstname lastname",
					PreferredLanguage:      language.English,
					PasswordEncodedHash:    "$plain$x$password",
					PasswordChangeRequired: true,
					Email:                  "email@test.ch",
					IsEmailVerified:        false,
					Phone:                  "+41791234567",
					IsPhoneVerified:        false,
					UserState:              domain.UserStateActive,
				},
			},
		},
		{
			name: "user added with phone code",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newRegisterHumanEvent("username", "$plain$x$password", true, true, "", language.English),
						),
						eventFromEventPusher(
							user.NewHumanPhoneChangedEvent(context.Background(),
								&userAgg.Aggregate,
								"+41791234567",
							),
						),
						eventFromEventPusher(
							user.NewHumanPhoneCodeAddedEventV2(context.Background(),
								&userAgg.Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("a"),
								},
								time.Hour*1,
								false,
							),
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				userID: "user1",
			},
			res: res{
				want: &UserV2WriteModel{
					HumanWriteModel: true,
					PhoneWriteModel: true,
					StateWriteModel: true,
					WriteModel: eventstore.WriteModel{
						AggregateID:       "user1",
						Events:            []eventstore.Event{},
						ProcessedSequence: 0,
						ResourceOwner:     "org1",
					},
					UserName:               "username",
					FirstName:              "firstname",
					LastName:               "lastname",
					DisplayName:            "firstname lastname",
					PreferredLanguage:      language.English,
					PasswordEncodedHash:    "$plain$x$password",
					PasswordChangeRequired: true,
					Email:                  "email@test.ch",
					IsEmailVerified:        false,
					Phone:                  "+41791234567",
					IsPhoneVerified:        false,
					PhoneCode: &crypto.CryptoValue{
						CryptoType: crypto.TypeEncryption,
						Algorithm:  "enc",
						KeyID:      "id",
						Crypted:    []byte("a"),
					},
					PhoneCodeCreationDate: time.Time{},
					PhoneCodeExpiry:       time.Hour * 1,
					UserState:             domain.UserStateActive,
				},
			},
		},
		{
			name: "user added with phone code verified",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newRegisterHumanEvent("username", "$plain$x$password", true, true, "", language.English),
						),
						eventFromEventPusher(
							user.NewHumanPhoneChangedEvent(context.Background(),
								&userAgg.Aggregate,
								"+41791234567",
							),
						),
						eventFromEventPusher(
							user.NewHumanPhoneCodeAddedEventV2(context.Background(),
								&userAgg.Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("a"),
								},
								time.Hour*1,
								false,
							),
						),
						eventFromEventPusher(
							user.NewHumanPhoneVerifiedEvent(context.Background(),
								&userAgg.Aggregate,
							),
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				userID: "user1",
			},
			res: res{
				want: &UserV2WriteModel{
					HumanWriteModel: true,
					PhoneWriteModel: true,
					StateWriteModel: true,
					WriteModel: eventstore.WriteModel{
						AggregateID:       "user1",
						Events:            []eventstore.Event{},
						ProcessedSequence: 0,
						ResourceOwner:     "org1",
					},
					UserName:               "username",
					FirstName:              "firstname",
					LastName:               "lastname",
					DisplayName:            "firstname lastname",
					PreferredLanguage:      language.English,
					PasswordEncodedHash:    "$plain$x$password",
					PasswordChangeRequired: true,
					Email:                  "email@test.ch",
					IsEmailVerified:        false,
					Phone:                  "+41791234567",
					IsPhoneVerified:        true,
					UserState:              domain.UserStateActive,
				},
			},
		},
		{
			name: "user added with phone code verified failed",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newRegisterHumanEvent("username", "$plain$x$password", true, true, "", language.English),
						),

						eventFromEventPusher(
							user.NewHumanPhoneChangedEvent(context.Background(),
								&userAgg.Aggregate,
								"+41791234567",
							),
						),
						eventFromEventPusher(
							user.NewHumanPhoneCodeAddedEventV2(context.Background(),
								&userAgg.Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("a"),
								},
								time.Hour*1,
								false,
							),
						),
						eventFromEventPusher(
							user.NewHumanPhoneVerificationFailedEvent(context.Background(),
								&userAgg.Aggregate,
							),
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				userID: "user1",
			},
			res: res{
				want: &UserV2WriteModel{
					HumanWriteModel: true,
					PhoneWriteModel: true,
					StateWriteModel: true,
					WriteModel: eventstore.WriteModel{
						AggregateID:       "user1",
						Events:            []eventstore.Event{},
						ProcessedSequence: 0,
						ResourceOwner:     "org1",
					},
					UserName:               "username",
					FirstName:              "firstname",
					LastName:               "lastname",
					DisplayName:            "firstname lastname",
					PreferredLanguage:      language.English,
					PasswordEncodedHash:    "$plain$x$password",
					PasswordChangeRequired: true,
					Email:                  "email@test.ch",
					IsEmailVerified:        false,
					Phone:                  "+41791234567",
					IsPhoneVerified:        false,
					PhoneCode: &crypto.CryptoValue{
						CryptoType: crypto.TypeEncryption,
						Algorithm:  "enc",
						KeyID:      "id",
						Crypted:    []byte("a"),
					},
					PhoneCodeCreationDate: time.Time{},
					PhoneCodeExpiry:       time.Hour * 1,
					PhoneCheckFailedCount: 1,
					UserState:             domain.UserStateActive,
				},
			},
		},
		{
			name: "user added with email code verified, then changed",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newRegisterHumanEvent("username", "$plain$x$password", true, true, "", language.English),
						),
						eventFromEventPusher(
							user.NewHumanPhoneChangedEvent(context.Background(),
								&userAgg.Aggregate,
								"+41791234567",
							),
						),
						eventFromEventPusher(
							user.NewHumanPhoneCodeAddedEventV2(context.Background(),
								&userAgg.Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("a"),
								},
								time.Hour*1,
								false,
							),
						),
						eventFromEventPusher(
							user.NewHumanPhoneVerifiedEvent(context.Background(),
								&userAgg.Aggregate,
							),
						),
						eventFromEventPusher(
							user.NewHumanPhoneChangedEvent(context.Background(),
								&userAgg.Aggregate,
								"+41797654321",
							),
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				userID: "user1",
			},
			res: res{
				want: &UserV2WriteModel{
					HumanWriteModel: true,
					PhoneWriteModel: true,
					StateWriteModel: true,
					WriteModel: eventstore.WriteModel{
						AggregateID:       "user1",
						Events:            []eventstore.Event{},
						ProcessedSequence: 0,
						ResourceOwner:     "org1",
					},
					UserName:               "username",
					FirstName:              "firstname",
					LastName:               "lastname",
					DisplayName:            "firstname lastname",
					PreferredLanguage:      language.English,
					PasswordEncodedHash:    "$plain$x$password",
					PasswordChangeRequired: true,
					Email:                  "email@test.ch",
					IsEmailVerified:        false,
					Phone:                  "+41797654321",
					IsPhoneVerified:        false,
					UserState:              domain.UserStateActive,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore(t),
			}
			wm, err := r.userHumanWriteModel(tt.args.ctx, tt.args.userID, false, false, true, false, false, false)
			if tt.res.err == nil {
				if !assert.NoError(t, err) {
					t.FailNow()
				}
			} else if !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
				return
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, wm)
			}
		})
	}
}

func TestCommandSide_userHumanWriteModel_password(t *testing.T) {
	type fields struct {
		eventstore func(t *testing.T) *eventstore.Eventstore
	}
	type args struct {
		ctx    context.Context
		userID string
	}
	type res struct {
		want *UserV2WriteModel
		err  func(error) bool
	}

	userAgg := user.NewAggregate("user1", "org1")

	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "user added password hashchanged",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						),
						eventFromEventPusher(
							user.NewHumanPasswordHashUpdatedEvent(context.Background(),
								&userAgg.Aggregate,
								"hash",
							),
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				userID: "user1",
			},
			res: res{
				want: &UserV2WriteModel{
					HumanWriteModel:    true,
					PasswordWriteModel: true,
					StateWriteModel:    true,
					WriteModel: eventstore.WriteModel{
						AggregateID:       "user1",
						Events:            []eventstore.Event{},
						ProcessedSequence: 0,
						ResourceOwner:     "org1",
					},
					UserName:               "username",
					FirstName:              "firstname",
					LastName:               "lastname",
					DisplayName:            "firstname lastname",
					PreferredLanguage:      language.English,
					PasswordEncodedHash:    "hash",
					PasswordChangeRequired: true,
					Email:                  "email@test.ch",
					IsEmailVerified:        false,
					UserState:              domain.UserStateActive,
				},
			},
		},
		{
			name: "user added password changed",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						),
						eventFromEventPusher(
							user.NewHumanPasswordChangedEvent(context.Background(),
								&userAgg.Aggregate,
								"hash",
								false,
								"",
							),
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				userID: "user1",
			},
			res: res{
				want: &UserV2WriteModel{
					HumanWriteModel:    true,
					PasswordWriteModel: true,
					StateWriteModel:    true,
					WriteModel: eventstore.WriteModel{
						AggregateID:       "user1",
						Events:            []eventstore.Event{},
						ProcessedSequence: 0,
						ResourceOwner:     "org1",
					},
					UserName:               "username",
					FirstName:              "firstname",
					LastName:               "lastname",
					DisplayName:            "firstname lastname",
					PreferredLanguage:      language.English,
					PasswordEncodedHash:    "hash",
					PasswordChangeRequired: false,
					Email:                  "email@test.ch",
					IsEmailVerified:        false,
					UserState:              domain.UserStateActive,
				},
			},
		},
		{
			name: "user added with password code",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newRegisterHumanEvent("username", "$plain$x$password", true, true, "", language.English),
						),
						eventFromEventPusher(
							user.NewHumanPasswordCodeAddedEventV2(context.Background(),
								&userAgg.Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("a"),
								},
								time.Hour*1,
								domain.NotificationTypeEmail,
								"",
								false,
							),
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				userID: "user1",
			},
			res: res{
				want: &UserV2WriteModel{
					HumanWriteModel:    true,
					PasswordWriteModel: true,
					StateWriteModel:    true,
					WriteModel: eventstore.WriteModel{
						AggregateID:       "user1",
						Events:            []eventstore.Event{},
						ProcessedSequence: 0,
						ResourceOwner:     "org1",
					},
					UserName:               "username",
					FirstName:              "firstname",
					LastName:               "lastname",
					DisplayName:            "firstname lastname",
					PreferredLanguage:      language.English,
					PasswordEncodedHash:    "$plain$x$password",
					PasswordChangeRequired: true,
					Email:                  "email@test.ch",
					IsEmailVerified:        false,
					PasswordCode: &crypto.CryptoValue{
						CryptoType: crypto.TypeEncryption,
						Algorithm:  "enc",
						KeyID:      "id",
						Crypted:    []byte("a"),
					},
					PasswordCodeCreationDate: time.Time{},
					PasswordCodeExpiry:       time.Hour * 1,
					UserState:                domain.UserStateActive,
				},
			},
		},
		{
			name: "user added with password code and then change",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newRegisterHumanEvent("username", "$plain$x$password", true, true, "", language.English),
						),
						eventFromEventPusher(
							user.NewHumanPasswordCodeAddedEventV2(context.Background(),
								&userAgg.Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("a"),
								},
								time.Hour*1,
								domain.NotificationTypeEmail,
								"",
								false,
							),
						),
						eventFromEventPusher(
							user.NewHumanPasswordChangedEvent(context.Background(),
								&userAgg.Aggregate,
								"hash",
								true,
								"",
							),
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				userID: "user1",
			},
			res: res{
				want: &UserV2WriteModel{
					HumanWriteModel:    true,
					PasswordWriteModel: true,
					StateWriteModel:    true,
					WriteModel: eventstore.WriteModel{
						AggregateID:       "user1",
						Events:            []eventstore.Event{},
						ProcessedSequence: 0,
						ResourceOwner:     "org1",
					},
					UserName:               "username",
					FirstName:              "firstname",
					LastName:               "lastname",
					DisplayName:            "firstname lastname",
					PreferredLanguage:      language.English,
					PasswordEncodedHash:    "hash",
					PasswordChangeRequired: true,
					Email:                  "email@test.ch",
					IsEmailVerified:        false,
					UserState:              domain.UserStateActive,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore(t),
			}
			wm, err := r.userHumanWriteModel(tt.args.ctx, tt.args.userID, false, false, false, true, false, false)
			if tt.res.err == nil {
				if !assert.NoError(t, err) {
					t.FailNow()
				}
			} else if !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
				return
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, wm)
			}
		})
	}
}

func TestCommandSide_userStateWriteModel(t *testing.T) {
	type fields struct {
		eventstore func(t *testing.T) *eventstore.Eventstore
	}
	type args struct {
		ctx    context.Context
		userID string
	}
	type res struct {
		want *UserV2WriteModel
		err  func(error) bool
	}

	userAgg := user.NewAggregate("user1", "org1")

	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "user added",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				userID: "user1",
			},
			res: res{
				want: &UserV2WriteModel{
					HumanWriteModel:   true,
					MachineWriteModel: true,
					StateWriteModel:   true,
					WriteModel: eventstore.WriteModel{
						AggregateID:       "user1",
						Events:            []eventstore.Event{},
						ProcessedSequence: 0,
						ResourceOwner:     "org1",
					},
					UserName:                 "username",
					FirstName:                "firstname",
					LastName:                 "lastname",
					DisplayName:              "firstname lastname",
					PreferredLanguage:        language.English,
					PasswordEncodedHash:      "$plain$x$password",
					PasswordChangeRequired:   true,
					PasswordCheckFailedCount: 0,
					Email:                    "email@test.ch",
					IsEmailVerified:          false,
					UserState:                domain.UserStateActive,
				},
			},
		},
		{
			name: "user added initialized",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						),
						eventFromEventPusher(
							user.NewHumanInitializedCheckSucceededEvent(context.Background(),
								&userAgg.Aggregate,
							),
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				userID: "user1",
			},
			res: res{
				want: &UserV2WriteModel{
					HumanWriteModel:   true,
					MachineWriteModel: true,
					StateWriteModel:   true,
					WriteModel: eventstore.WriteModel{
						AggregateID:       "user1",
						Events:            []eventstore.Event{},
						ProcessedSequence: 0,
						ResourceOwner:     "org1",
					},
					UserName:                 "username",
					FirstName:                "firstname",
					LastName:                 "lastname",
					DisplayName:              "firstname lastname",
					PreferredLanguage:        language.English,
					PasswordEncodedHash:      "$plain$x$password",
					PasswordChangeRequired:   true,
					PasswordCheckFailedCount: 0,
					Email:                    "email@test.ch",
					IsEmailVerified:          false,
					UserState:                domain.UserStateActive,
				},
			},
		},
		{
			name: "user machine added",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddMachineEvent(true, domain.OIDCTokenTypeBearer),
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				userID: "user1",
			},
			res: res{
				want: &UserV2WriteModel{
					HumanWriteModel:   true,
					MachineWriteModel: true,
					StateWriteModel:   true,
					WriteModel: eventstore.WriteModel{
						AggregateID:       "user1",
						Events:            []eventstore.Event{},
						ProcessedSequence: 0,
						ResourceOwner:     "org1",
					},
					UserName:        "username",
					Name:            "name",
					Description:     "description",
					AccessTokenType: domain.OIDCTokenTypeBearer,
					UserState:       domain.UserStateActive,
				},
			},
		},
		{
			name: "user added locked",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						),
						eventFromEventPusher(
							user.NewHumanPasswordCheckFailedEvent(context.Background(),
								&userAgg.Aggregate,
								nil,
							),
						),
						eventFromEventPusher(
							user.NewHumanPasswordCheckFailedEvent(context.Background(),
								&userAgg.Aggregate,
								nil,
							),
						),
						eventFromEventPusher(
							user.NewHumanPasswordCheckFailedEvent(context.Background(),
								&userAgg.Aggregate,
								nil,
							),
						),
						eventFromEventPusher(
							user.NewUserLockedEvent(context.Background(),
								&userAgg.Aggregate,
							),
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				userID: "user1",
			},
			res: res{
				want: &UserV2WriteModel{
					HumanWriteModel:   true,
					MachineWriteModel: true,
					StateWriteModel:   true,
					WriteModel: eventstore.WriteModel{
						AggregateID:       "user1",
						Events:            []eventstore.Event{},
						ProcessedSequence: 0,
						ResourceOwner:     "org1",
					},
					UserName:                 "username",
					FirstName:                "firstname",
					LastName:                 "lastname",
					DisplayName:              "firstname lastname",
					PreferredLanguage:        language.English,
					PasswordEncodedHash:      "$plain$x$password",
					PasswordChangeRequired:   true,
					PasswordCheckFailedCount: 3,
					Email:                    "email@test.ch",
					IsEmailVerified:          false,
					UserState:                domain.UserStateLocked,
				},
			},
		},
		{
			name: "user added locked and unlocked",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newRegisterHumanEvent("username", "$plain$x$password", true, true, "", language.English),
						),
						eventFromEventPusher(
							user.NewHumanPasswordCheckFailedEvent(context.Background(),
								&userAgg.Aggregate,
								nil,
							),
						),
						eventFromEventPusher(
							user.NewHumanPasswordCheckFailedEvent(context.Background(),
								&userAgg.Aggregate,
								nil,
							),
						),
						eventFromEventPusher(
							user.NewHumanPasswordCheckFailedEvent(context.Background(),
								&userAgg.Aggregate,
								nil,
							),
						),
						eventFromEventPusher(
							user.NewUserLockedEvent(context.Background(),
								&userAgg.Aggregate,
							),
						),
						eventFromEventPusher(
							user.NewUserUnlockedEvent(context.Background(),
								&userAgg.Aggregate,
							),
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				userID: "user1",
			},
			res: res{
				want: &UserV2WriteModel{
					HumanWriteModel:   true,
					MachineWriteModel: true,
					StateWriteModel:   true,
					WriteModel: eventstore.WriteModel{
						AggregateID:       "user1",
						Events:            []eventstore.Event{},
						ProcessedSequence: 0,
						ResourceOwner:     "org1",
					},
					UserName:                 "username",
					FirstName:                "firstname",
					LastName:                 "lastname",
					DisplayName:              "firstname lastname",
					PreferredLanguage:        language.English,
					PasswordEncodedHash:      "$plain$x$password",
					PasswordChangeRequired:   true,
					PasswordCheckFailedCount: 0,
					Email:                    "email@test.ch",
					IsEmailVerified:          false,
					UserState:                domain.UserStateActive,
				},
			},
		},
		{
			name: "user added deactivated",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newRegisterHumanEvent("username", "$plain$x$password", true, true, "", language.English),
						),
						eventFromEventPusher(
							user.NewUserDeactivatedEvent(context.Background(),
								&userAgg.Aggregate,
							),
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				userID: "user1",
			},
			res: res{
				want: &UserV2WriteModel{
					HumanWriteModel:   true,
					MachineWriteModel: true,
					StateWriteModel:   true,
					WriteModel: eventstore.WriteModel{
						AggregateID:       "user1",
						Events:            []eventstore.Event{},
						ProcessedSequence: 0,
						ResourceOwner:     "org1",
					},
					UserName:               "username",
					FirstName:              "firstname",
					LastName:               "lastname",
					DisplayName:            "firstname lastname",
					PreferredLanguage:      language.English,
					PasswordEncodedHash:    "$plain$x$password",
					PasswordChangeRequired: true,
					Email:                  "email@test.ch",
					IsEmailVerified:        false,
					UserState:              domain.UserStateInactive,
				},
			},
		},
		{
			name: "user added deactivated and reactived",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newRegisterHumanEvent("username", "$plain$x$password", true, true, "", language.English),
						),
						eventFromEventPusher(
							user.NewUserDeactivatedEvent(context.Background(),
								&userAgg.Aggregate,
							),
						),
						eventFromEventPusher(
							user.NewUserReactivatedEvent(context.Background(),
								&userAgg.Aggregate,
							),
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				userID: "user1",
			},
			res: res{
				want: &UserV2WriteModel{
					HumanWriteModel:   true,
					MachineWriteModel: true,
					StateWriteModel:   true,
					WriteModel: eventstore.WriteModel{
						AggregateID:       "user1",
						Events:            []eventstore.Event{},
						ProcessedSequence: 0,
						ResourceOwner:     "org1",
					},
					UserName:               "username",
					FirstName:              "firstname",
					LastName:               "lastname",
					DisplayName:            "firstname lastname",
					PreferredLanguage:      language.English,
					PasswordEncodedHash:    "$plain$x$password",
					PasswordChangeRequired: true,
					Email:                  "email@test.ch",
					IsEmailVerified:        false,
					UserState:              domain.UserStateActive,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore(t),
			}
			wm, err := r.userStateWriteModel(tt.args.ctx, tt.args.userID)
			if tt.res.err == nil {
				if !assert.NoError(t, err) {
					t.FailNow()
				}
			} else if !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
				return
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, wm)
			}
		})
	}
}

func TestCommandSide_userHumanWriteModel_avatar(t *testing.T) {
	type fields struct {
		eventstore func(t *testing.T) *eventstore.Eventstore
	}
	type args struct {
		ctx    context.Context
		userID string
	}
	type res struct {
		want *UserV2WriteModel
		err  func(error) bool
	}

	userAgg := user.NewAggregate("user1", "org1")

	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "user added with avatar",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						),
						eventFromEventPusher(
							user.NewHumanAvatarAddedEvent(context.Background(),
								&userAgg.Aggregate,
								"key",
							),
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				userID: "user1",
			},
			res: res{
				want: &UserV2WriteModel{
					HumanWriteModel:  true,
					AvatarWriteModel: true,
					StateWriteModel:  true,
					WriteModel: eventstore.WriteModel{
						AggregateID:       "user1",
						Events:            []eventstore.Event{},
						ProcessedSequence: 0,
						ResourceOwner:     "org1",
					},
					UserName:               "username",
					FirstName:              "firstname",
					LastName:               "lastname",
					DisplayName:            "firstname lastname",
					PreferredLanguage:      language.English,
					PasswordEncodedHash:    "$plain$x$password",
					PasswordChangeRequired: true,
					Email:                  "email@test.ch",
					IsEmailVerified:        false,
					UserState:              domain.UserStateActive,
					Avatar:                 "key",
				},
			},
		},

		{
			name: "user added with avatar and then removed",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						),
						eventFromEventPusher(
							user.NewHumanAvatarAddedEvent(context.Background(),
								&userAgg.Aggregate,
								"key",
							),
						),
						eventFromEventPusher(
							user.NewHumanAvatarRemovedEvent(context.Background(),
								&userAgg.Aggregate,
								"key",
							),
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				userID: "user1",
			},
			res: res{
				want: &UserV2WriteModel{
					HumanWriteModel:  true,
					AvatarWriteModel: true,
					StateWriteModel:  true,
					WriteModel: eventstore.WriteModel{
						AggregateID:       "user1",
						Events:            []eventstore.Event{},
						ProcessedSequence: 0,
						ResourceOwner:     "org1",
					},
					UserName:               "username",
					FirstName:              "firstname",
					LastName:               "lastname",
					DisplayName:            "firstname lastname",
					PreferredLanguage:      language.English,
					PasswordEncodedHash:    "$plain$x$password",
					PasswordChangeRequired: true,
					Email:                  "email@test.ch",
					IsEmailVerified:        false,
					UserState:              domain.UserStateActive,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore(t),
			}
			wm, err := r.userHumanWriteModel(tt.args.ctx, tt.args.userID, false, false, false, false, true, false)
			if tt.res.err == nil {
				if !assert.NoError(t, err) {
					t.FailNow()
				}
			} else if !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
				return
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, wm)
			}
		})
	}
}

func TestCommandSide_userHumanWriteModel_idpLinks(t *testing.T) {
	type fields struct {
		eventstore func(t *testing.T) *eventstore.Eventstore
	}
	type args struct {
		ctx    context.Context
		userID string
	}
	type res struct {
		want *UserV2WriteModel
		err  func(error) bool
	}

	userAgg := user.NewAggregate("user1", "org1")

	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "user added with idp link",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						),
						eventFromEventPusher(
							user.NewUserIDPLinkAddedEvent(context.Background(),
								&userAgg.Aggregate,
								"idp",
								"name",
								"externalID",
							),
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				userID: "user1",
			},
			res: res{
				want: &UserV2WriteModel{
					HumanWriteModel:   true,
					MachineWriteModel: true,
					StateWriteModel:   true,
					IDPLinkWriteModel: true,
					WriteModel: eventstore.WriteModel{
						AggregateID:       "user1",
						Events:            []eventstore.Event{},
						ProcessedSequence: 0,
						ResourceOwner:     "org1",
					},
					UserName:               "username",
					FirstName:              "firstname",
					LastName:               "lastname",
					DisplayName:            "firstname lastname",
					PreferredLanguage:      language.English,
					PasswordEncodedHash:    "$plain$x$password",
					PasswordChangeRequired: true,
					Email:                  "email@test.ch",
					IsEmailVerified:        false,
					UserState:              domain.UserStateActive,
					IDPLinks: []*domain.UserIDPLink{
						{IDPConfigID: "idp", DisplayName: "name", ExternalUserID: "externalID"},
					},
				},
			},
		},
		{
			name: "user added with idp links",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						),
						eventFromEventPusher(
							user.NewUserIDPLinkAddedEvent(context.Background(),
								&userAgg.Aggregate,
								"idp1",
								"name1",
								"externalID1",
							),
						),
						eventFromEventPusher(
							user.NewUserIDPLinkAddedEvent(context.Background(),
								&userAgg.Aggregate,
								"idp2",
								"name2",
								"externalID2",
							),
						),
						eventFromEventPusher(
							user.NewUserIDPLinkAddedEvent(context.Background(),
								&userAgg.Aggregate,
								"idp3",
								"name3",
								"externalID3",
							),
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				userID: "user1",
			},
			res: res{
				want: &UserV2WriteModel{
					HumanWriteModel:   true,
					MachineWriteModel: true,
					StateWriteModel:   true,
					IDPLinkWriteModel: true,
					WriteModel: eventstore.WriteModel{
						AggregateID:       "user1",
						Events:            []eventstore.Event{},
						ProcessedSequence: 0,
						ResourceOwner:     "org1",
					},
					UserName:               "username",
					FirstName:              "firstname",
					LastName:               "lastname",
					DisplayName:            "firstname lastname",
					PreferredLanguage:      language.English,
					PasswordEncodedHash:    "$plain$x$password",
					PasswordChangeRequired: true,
					Email:                  "email@test.ch",
					IsEmailVerified:        false,
					UserState:              domain.UserStateActive,
					IDPLinks: []*domain.UserIDPLink{
						{IDPConfigID: "idp1", DisplayName: "name1", ExternalUserID: "externalID1"},
						{IDPConfigID: "idp2", DisplayName: "name2", ExternalUserID: "externalID2"},
						{IDPConfigID: "idp3", DisplayName: "name3", ExternalUserID: "externalID3"},
					},
				},
			},
		},
		{
			name: "user added with idp links and removed",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						),
						eventFromEventPusher(
							user.NewUserIDPLinkAddedEvent(context.Background(),
								&userAgg.Aggregate,
								"idp1",
								"name1",
								"externalID1",
							),
						),
						eventFromEventPusher(
							user.NewUserIDPLinkAddedEvent(context.Background(),
								&userAgg.Aggregate,
								"idp2",
								"name2",
								"externalID2",
							),
						),
						eventFromEventPusher(
							user.NewUserIDPLinkAddedEvent(context.Background(),
								&userAgg.Aggregate,
								"idp3",
								"name3",
								"externalID3",
							),
						),
						eventFromEventPusher(
							user.NewUserIDPLinkCascadeRemovedEvent(context.Background(),
								&userAgg.Aggregate,
								"idp2",
								"externalID2",
							),
						),
						eventFromEventPusher(
							user.NewUserIDPLinkRemovedEvent(context.Background(),
								&userAgg.Aggregate,
								"idp3",
								"externalID3",
							),
						),
						eventFromEventPusher(
							user.NewUserIDPLinkAddedEvent(context.Background(),
								&userAgg.Aggregate,
								"idp4",
								"name4",
								"externalID4",
							),
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				userID: "user1",
			},
			res: res{
				want: &UserV2WriteModel{
					HumanWriteModel:   true,
					MachineWriteModel: true,
					StateWriteModel:   true,
					IDPLinkWriteModel: true,
					WriteModel: eventstore.WriteModel{
						AggregateID:       "user1",
						Events:            []eventstore.Event{},
						ProcessedSequence: 0,
						ResourceOwner:     "org1",
					},
					UserName:               "username",
					FirstName:              "firstname",
					LastName:               "lastname",
					DisplayName:            "firstname lastname",
					PreferredLanguage:      language.English,
					PasswordEncodedHash:    "$plain$x$password",
					PasswordChangeRequired: true,
					Email:                  "email@test.ch",
					IsEmailVerified:        false,
					UserState:              domain.UserStateActive,
					IDPLinks: []*domain.UserIDPLink{
						{IDPConfigID: "idp1", DisplayName: "name1", ExternalUserID: "externalID1"},
						{IDPConfigID: "idp4", DisplayName: "name4", ExternalUserID: "externalID4"},
					},
				},
			},
		},
		{
			name: "user added with idp link and removed",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						),
						eventFromEventPusher(
							user.NewUserIDPLinkAddedEvent(context.Background(),
								&userAgg.Aggregate,
								"idp",
								"name",
								"externalID",
							),
						),
						eventFromEventPusher(
							user.NewUserIDPLinkRemovedEvent(context.Background(),
								&userAgg.Aggregate,
								"idp",
								"externalID",
							),
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				userID: "user1",
			},
			res: res{
				want: &UserV2WriteModel{
					HumanWriteModel:   true,
					MachineWriteModel: true,
					StateWriteModel:   true,
					IDPLinkWriteModel: true,
					WriteModel: eventstore.WriteModel{
						AggregateID:       "user1",
						Events:            []eventstore.Event{},
						ProcessedSequence: 0,
						ResourceOwner:     "org1",
					},
					UserName:               "username",
					FirstName:              "firstname",
					LastName:               "lastname",
					DisplayName:            "firstname lastname",
					PreferredLanguage:      language.English,
					PasswordEncodedHash:    "$plain$x$password",
					PasswordChangeRequired: true,
					Email:                  "email@test.ch",
					IsEmailVerified:        false,
					UserState:              domain.UserStateActive,
					IDPLinks:               []*domain.UserIDPLink{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore(t),
			}
			wm, err := r.userRemoveWriteModel(tt.args.ctx, tt.args.userID)
			if tt.res.err == nil {
				if !assert.NoError(t, err) {
					t.FailNow()
				}
			} else if !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
				return
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, wm)
			}
		})
	}
}
