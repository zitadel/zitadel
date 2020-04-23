package eventsourcing

import (
	"context"
	"github.com/caos/zitadel/internal/api/auth"
	caos_errs "github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/user/model"
	"github.com/golang/mock/gomock"
	"testing"
)

func TestUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es   *UserEventstore
		user *model.User
	}
	type res struct {
		user    *model.User
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "user from events, ok",
			args: args{
				es:   GetMockUserByIDOK(ctrl),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
		},
		{
			name: "no events found",
			args: args{
				es:   GetMockUserByIDNoEvents(ctrl),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
		{
			name: "no id",
			args: args{
				es:   GetMockUserByIDNoEvents(ctrl),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "", Sequence: 1}},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.UserByID(nil, tt.args.user.AggregateID)

			if !tt.res.wantErr && result.AggregateID != tt.res.user.AggregateID {
				t.Errorf("got wrong result name: expected: %v, actual: %v ", tt.res.user.AggregateID, result.AggregateID)
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es   *UserEventstore
		ctx  context.Context
		user *model.User
	}
	type res struct {
		user    *model.User
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "create user, ok",
			args: args{
				es:  GetMockManipulateUserWithPW(ctrl, true, false, false, false),
				ctx: auth.NewMockContext("orgID", "userID"),
				user: &model.User{
					ObjectRoot: es_models.ObjectRoot{Sequence: 1},
					Profile: &model.Profile{
						UserName:  "UserName",
						FirstName: "FirstName",
						LastName:  "LastName",
					},
					Email: &model.Email{
						EmailAddress:    "EmailAddress",
						IsEmailVerified: true,
					},
				},
			},
			res: res{
				user: &model.User{ObjectRoot: es_models.ObjectRoot{Sequence: 1},
					Profile: &model.Profile{
						UserName:  "UserName",
						FirstName: "FirstName",
						LastName:  "LastName",
					},
					Email: &model.Email{
						EmailAddress: "EmailAddress",
					},
				},
			},
		},
		{
			name: "no username, should use email",
			args: args{
				es:  GetMockManipulateUserWithPW(ctrl, true, false, false, false),
				ctx: auth.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{Sequence: 1},
					Profile: &model.Profile{
						FirstName: "FirstName",
						LastName:  "LastName",
					},
					Email: &model.Email{
						EmailAddress:    "EmailAddress",
						IsEmailVerified: true,
					},
				},
			},
			res: res{
				user: &model.User{ObjectRoot: es_models.ObjectRoot{Sequence: 1},
					Profile: &model.Profile{
						UserName:  "EmailAddress",
						FirstName: "FirstName",
						LastName:  "LastName",
					},
					Email: &model.Email{
						EmailAddress: "EmailAddress",
					},
				},
			},
		},
		{
			name: "with phone code",
			args: args{
				es:  GetMockManipulateUserWithPW(ctrl, true, false, false, false),
				ctx: auth.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{Sequence: 1},
					Profile: &model.Profile{
						FirstName: "FirstName",
						LastName:  "LastName",
						UserName:  "UserName",
					},
					Email: &model.Email{
						EmailAddress:    "UserName",
						IsEmailVerified: true,
					},
					Phone: &model.Phone{
						PhoneNumber:     "UserName",
						IsPhoneVerified: true,
					},
				},
			},
			res: res{
				user: &model.User{ObjectRoot: es_models.ObjectRoot{Sequence: 1},
					Profile: &model.Profile{
						UserName:  "UserName",
						FirstName: "FirstName",
						LastName:  "LastName",
					},
					Email: &model.Email{
						EmailAddress: "EmailAddress",
					},
				},
			},
		},
		{
			name: "with password",
			args: args{
				es:  GetMockManipulateUserWithPW(ctrl, false, false, false, true),
				ctx: auth.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{Sequence: 1},
					Profile: &model.Profile{
						FirstName: "FirstName",
						LastName:  "LastName",
						UserName:  "UserName",
					},
					Password: &model.Password{SecretString: "Password"},
					Email: &model.Email{
						EmailAddress:    "UserName",
						IsEmailVerified: true,
					},
					Phone: &model.Phone{
						PhoneNumber:     "UserName",
						IsPhoneVerified: true,
					},
				},
			},
			res: res{
				user: &model.User{ObjectRoot: es_models.ObjectRoot{Sequence: 1},
					Profile: &model.Profile{
						UserName:  "UserName",
						FirstName: "FirstName",
						LastName:  "LastName",
					},
					Email: &model.Email{
						EmailAddress: "EmailAddress",
					},
				},
			},
		},
		{
			name: "create user invalid",
			args: args{
				es:   GetMockManipulateUser(ctrl),
				ctx:  auth.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.CreateUser(tt.args.ctx, tt.args.user)

			if !tt.res.wantErr && result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if !tt.res.wantErr && result.UserName != tt.res.user.UserName {
				t.Errorf("got wrong result name: expected: %v, actual: %v ", tt.res.user.UserName, result.UserName)
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestRegisterUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es            *UserEventstore
		ctx           context.Context
		user          *model.User
		resourceOwner string
	}
	type res struct {
		user    *model.User
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "register user, ok",
			args: args{
				es:  GetMockManipulateUserWithPW(ctrl, false, true, false, true),
				ctx: auth.NewMockContext("orgID", "userID"),
				user: &model.User{
					ObjectRoot: es_models.ObjectRoot{Sequence: 1},
					Profile: &model.Profile{
						UserName:  "UserName",
						FirstName: "FirstName",
						LastName:  "LastName",
					},
					Email: &model.Email{
						EmailAddress: "EmailAddress",
					},
					Password: &model.Password{
						SecretString: "Password",
					},
				},
				resourceOwner: "ResourceOwner",
			},
			res: res{
				user: &model.User{ObjectRoot: es_models.ObjectRoot{Sequence: 1},
					Profile: &model.Profile{
						UserName:  "UserName",
						FirstName: "FirstName",
						LastName:  "LastName",
					},
					Email: &model.Email{
						EmailAddress: "EmailAddress",
					},
				},
			},
		},
		{
			name: "no username, should use email",
			args: args{
				es:  GetMockManipulateUserWithPW(ctrl, false, true, false, true),
				ctx: auth.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{Sequence: 1},
					Profile: &model.Profile{
						FirstName: "FirstName",
						LastName:  "LastName",
					},
					Email: &model.Email{
						EmailAddress: "EmailAddress",
					},
					Password: &model.Password{
						SecretString: "Password",
					},
				},
				resourceOwner: "ResourceOwner",
			},
			res: res{
				user: &model.User{ObjectRoot: es_models.ObjectRoot{Sequence: 1},
					Profile: &model.Profile{
						UserName:  "EmailAddress",
						FirstName: "FirstName",
						LastName:  "LastName",
					},
					Email: &model.Email{
						EmailAddress: "EmailAddress",
					},
				},
			},
		},
		{
			name: "register user invalid",
			args: args{
				es:            GetMockManipulateUser(ctrl),
				ctx:           auth.NewMockContext("orgID", "userID"),
				user:          &model.User{ObjectRoot: es_models.ObjectRoot{Sequence: 1}},
				resourceOwner: "ResourceOwner",
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "register user no password",
			args: args{
				es:  GetMockManipulateUser(ctrl),
				ctx: auth.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{Sequence: 1},
					Profile: &model.Profile{
						UserName:  "EmailAddress",
						FirstName: "FirstName",
						LastName:  "LastName",
					},
					Email: &model.Email{
						EmailAddress: "EmailAddress",
					},
				},
				resourceOwner: "ResourceOwner",
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "no resourceowner",
			args: args{
				es:  GetMockManipulateUser(ctrl),
				ctx: auth.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{Sequence: 1},
					Profile: &model.Profile{
						UserName:  "EmailAddress",
						FirstName: "FirstName",
						LastName:  "LastName",
					},
					Email: &model.Email{
						EmailAddress: "EmailAddress",
					},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.RegisterUser(tt.args.ctx, tt.args.user, tt.args.resourceOwner)

			if !tt.res.wantErr && result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if !tt.res.wantErr && result.UserName != tt.res.user.UserName {
				t.Errorf("got wrong result name: expected: %v, actual: %v ", tt.res.user.UserName, result.UserName)
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestDeactivateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es       *UserEventstore
		ctx      context.Context
		existing *model.User
	}
	type res struct {
		user    *model.User
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "deactivate user, ok",
			args: args{
				es:       GetMockManipulateUser(ctrl),
				ctx:      auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, Profile: &model.Profile{UserName: "UserName"}, State: model.USERSTATE_INACTIVE},
			},
		},
		{
			name: "deactivate user with inactive state",
			args: args{
				es:       GetMockManipulateInactiveUser(ctrl),
				ctx:      auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing not found",
			args: args{
				es:       GetMockManipulateUserNoEvents(ctrl),
				ctx:      auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.DeactivateUser(tt.args.ctx, tt.args.existing.AggregateID)

			if !tt.res.wantErr && result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if !tt.res.wantErr && result.State != tt.res.user.State {
				t.Errorf("got wrong result name: expected: %v, actual: %v ", tt.res.user.State, result.State)
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestReactivateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es       *UserEventstore
		ctx      context.Context
		existing *model.User
	}
	type res struct {
		user    *model.User
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "reactivate user, ok",
			args: args{
				es:       GetMockManipulateInactiveUser(ctrl),
				ctx:      auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, Profile: &model.Profile{UserName: "UserName"}, State: model.USERSTATE_ACTIVE},
			},
		},
		{
			name: "reactivate user with inital state",
			args: args{
				es:       GetMockManipulateUser(ctrl),
				ctx:      auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing not found",
			args: args{
				es:       GetMockManipulateUserNoEvents(ctrl),
				ctx:      auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.ReactivateUser(tt.args.ctx, tt.args.existing.AggregateID)

			if !tt.res.wantErr && result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if !tt.res.wantErr && result.State != tt.res.user.State {
				t.Errorf("got wrong result name: expected: %v, actual: %v ", tt.res.user.State, result.State)
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestLockUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es       *UserEventstore
		ctx      context.Context
		existing *model.User
	}
	type res struct {
		user    *model.User
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "lock user, ok",
			args: args{
				es:       GetMockManipulateUser(ctrl),
				ctx:      auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, Profile: &model.Profile{UserName: "UserName"}, State: model.USERSTATE_LOCKED},
			},
		},
		{
			name: "lock user with locked state",
			args: args{
				es:       GetMockManipulateLockedUser(ctrl),
				ctx:      auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing not found",
			args: args{
				es:       GetMockManipulateUserNoEvents(ctrl),
				ctx:      auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.LockUser(tt.args.ctx, tt.args.existing.AggregateID)

			if !tt.res.wantErr && result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if !tt.res.wantErr && result.State != tt.res.user.State {
				t.Errorf("got wrong result name: expected: %v, actual: %v ", tt.res.user.State, result.State)
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestUnlockUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es       *UserEventstore
		ctx      context.Context
		existing *model.User
	}
	type res struct {
		user    *model.User
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "lock user, ok",
			args: args{
				es:       GetMockManipulateLockedUser(ctrl),
				ctx:      auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, Profile: &model.Profile{UserName: "UserName"}, State: model.USERSTATE_ACTIVE},
			},
		},
		//{
		//	name: "lock user not locked state",
		//	args: args{
		//		es:       GetMockManipulateUser(ctrl),
		//		ctx:      auth.NewMockContext("orgID", "userID"),
		//		existing: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
		//	},
		//	res: res{
		//		wantErr: true,
		//		errFunc: caos_errs.IsPreconditionFailed,
		//	},
		//},
		//{
		//	name: "existing not found",
		//	args: args{
		//		es:       GetMockManipulateUserNoEvents(ctrl),
		//		ctx:      auth.NewMockContext("orgID", "userID"),
		//		existing: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
		//	},
		//	res: res{
		//		wantErr: true,
		//		errFunc: caos_errs.IsNotFound,
		//	},
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.UnlockUser(tt.args.ctx, tt.args.existing.AggregateID)

			if !tt.res.wantErr && result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if !tt.res.wantErr && result.State != tt.res.user.State {
				t.Errorf("got wrong result name: expected: %v, actual: %v ", tt.res.user.State, result.State)
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}
