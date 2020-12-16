package eventsourcing

import (
	"context"
	"encoding/json"
	"net"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/caos/zitadel/internal/api/authz"
	req_model "github.com/caos/zitadel/internal/auth_request/model"
	"github.com/caos/zitadel/internal/crypto"
	caos_errs "github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/user/model"
	repo_model "github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
)

func TestUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es   *UserEventstore
		user *model.User
	}
	type res struct {
		user    *model.User
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
				es:   GetMockUserByIDOK(ctrl, repo_model.User{Human: &repo_model.Human{Profile: &repo_model.Profile{DisplayName: "DisplayName"}}, UserName: "UserName"}),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				user: &model.User{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					UserName:   "UserName",
					Human: &model.Human{
						Profile: &model.Profile{DisplayName: "DisplayName"},
					},
				},
			},
		},
		{
			name: "no events found",
			args: args{
				es:   GetMockUserByIDNoEvents(ctrl),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
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
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.UserByID(nil, tt.args.user.AggregateID)

			if tt.res.errFunc == nil && result.AggregateID != tt.res.user.AggregateID {
				t.Errorf("got wrong result aggregateID: expected: %v, actual: %v ", tt.res.user.AggregateID, result.AggregateID)
			}
			if tt.res.errFunc == nil && result.UserName != tt.res.user.UserName {
				t.Errorf("got wrong result userName: expected: %v, actual: %v ", tt.res.user.UserName, result.UserName)
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es        *UserEventstore
		ctx       context.Context
		user      *model.User
		policy    *iam_model.PasswordComplexityPolicyView
		orgPolicy *iam_model.OrgIAMPolicyView
	}
	type res struct {
		user    *model.User
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "init mail because no pw",
			args: args{
				es: GetMockManipulateUserWithInitCodeGen(ctrl, repo_model.User{
					UserName: "",
					Human: &repo_model.Human{
						Profile: &repo_model.Profile{DisplayName: "DisplayName", FirstName: "FirstName", LastName: "LastName"},
						Email:   &repo_model.Email{EmailAddress: "EmailAddress", IsEmailVerified: true},
					},
				}),
				ctx: authz.NewMockContext("orgID", "userID"),
				user: &model.User{
					ObjectRoot: es_models.ObjectRoot{Sequence: 1},
					UserName:   "UserName",
					Human: &model.Human{
						Profile: &model.Profile{
							FirstName:   "FirstName",
							LastName:    "LastName",
							DisplayName: "DisplayName",
						},
						Email: &model.Email{
							EmailAddress:    "EmailAddress",
							IsEmailVerified: true,
						},
					},
				},
				policy:    &iam_model.PasswordComplexityPolicyView{},
				orgPolicy: &iam_model.OrgIAMPolicyView{},
			},
			res: res{
				user: &model.User{
					ObjectRoot: es_models.ObjectRoot{Sequence: 1},
					UserName:   "UserName",
					Human: &model.Human{
						Profile: &model.Profile{
							FirstName:   "FirstName",
							LastName:    "LastName",
							DisplayName: "DisplayName",
						},
						Email: &model.Email{
							EmailAddress:    "EmailAddress",
							IsEmailVerified: true,
						},
					},
				},
			},
		},
		{
			name: "email as username",
			args: args{
				es: GetMockManipulateUserWithInitCodeGen(ctrl, repo_model.User{
					Human: &repo_model.Human{
						Profile: &repo_model.Profile{DisplayName: "DisplayName", FirstName: "FirstName", LastName: "LastName"},
						Email:   &repo_model.Email{EmailAddress: "EmailAddress", IsEmailVerified: true},
					},
				}),
				ctx: authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{Sequence: 1},
					Human: &model.Human{
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
				policy:    &iam_model.PasswordComplexityPolicyView{},
				orgPolicy: &iam_model.OrgIAMPolicyView{UserLoginMustBeDomain: false},
			},
			res: res{
				user: &model.User{ObjectRoot: es_models.ObjectRoot{Sequence: 1},
					UserName: "EmailAddress",
					Human: &model.Human{
						Profile: &model.Profile{
							FirstName:   "FirstName",
							LastName:    "LastName",
							DisplayName: "DisplayName",
						},
						Email: &model.Email{
							EmailAddress:    "EmailAddress",
							IsEmailVerified: true,
						},
					},
				},
			},
		},
		{
			name: "with verified phone number",
			args: args{
				es: GetMockManipulateUserWithInitCodeGen(ctrl, repo_model.User{
					UserName: "EmailAddress",
					Human: &repo_model.Human{
						Profile: &repo_model.Profile{DisplayName: "DisplayName", FirstName: "FirstName", LastName: "LastName"},
						Email:   &repo_model.Email{EmailAddress: "EmailAddress", IsEmailVerified: true},
						Phone:   &repo_model.Phone{PhoneNumber: "+41711234567", IsPhoneVerified: true},
					},
				}),
				ctx: authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{Sequence: 1},
					UserName: "UserName",
					Human: &model.Human{
						Profile: &model.Profile{
							FirstName:   "FirstName",
							LastName:    "LastName",
							DisplayName: "DisplayName",
						},
						Email: &model.Email{
							EmailAddress:    "UserName",
							IsEmailVerified: true,
						},
						Phone: &model.Phone{
							PhoneNumber:     "+41711234567",
							IsPhoneVerified: true,
						},
					},
				},
				policy:    &iam_model.PasswordComplexityPolicyView{},
				orgPolicy: &iam_model.OrgIAMPolicyView{},
			},
			res: res{
				user: &model.User{ObjectRoot: es_models.ObjectRoot{Sequence: 1},
					UserName: "UserName",
					Human: &model.Human{
						Profile: &model.Profile{
							FirstName:   "FirstName",
							LastName:    "LastName",
							DisplayName: "DisplayName",
						},
						Email: &model.Email{
							EmailAddress:    "EmailAddress",
							IsEmailVerified: true,
						},
						Phone: &model.Phone{
							PhoneNumber:     "+41711234567",
							IsPhoneVerified: true,
						},
					},
				},
			},
		},
		{
			name: "with password",
			args: args{
				es: GetMockManipulateUserWithPasswordAndEmailCodeGen(ctrl, repo_model.User{
					UserName: "UserName",
					Human: &repo_model.Human{
						Profile: &repo_model.Profile{DisplayName: "DisplayName", FirstName: "FirstName", LastName: "LastName"},
						Email:   &repo_model.Email{EmailAddress: "EmailAddress", IsEmailVerified: true}},
				}),
				ctx: authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{Sequence: 1},
					UserName: "UserName",
					Human: &model.Human{
						Profile: &model.Profile{
							FirstName:   "FirstName",
							LastName:    "LastName",
							DisplayName: "DisplayName",
						},
						Password: &model.Password{SecretString: "Password"},
						Email: &model.Email{
							EmailAddress:    "UserName",
							IsEmailVerified: true,
						},
					},
				},
				policy:    &iam_model.PasswordComplexityPolicyView{},
				orgPolicy: &iam_model.OrgIAMPolicyView{},
			},
			res: res{
				user: &model.User{ObjectRoot: es_models.ObjectRoot{Sequence: 1},
					UserName: "UserName",
					Human: &model.Human{
						Profile: &model.Profile{
							FirstName:   "FirstName",
							LastName:    "LastName",
							DisplayName: "DisplayName",
						},
						Email: &model.Email{
							EmailAddress:    "EmailAddress",
							IsEmailVerified: true,
						},
					},
				},
			},
		},
		{
			name: "create user invalid",
			args: args{
				es:        GetMockManipulateUser(ctrl),
				ctx:       authz.NewMockContext("orgID", "userID"),
				user:      &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, Human: &model.Human{}},
				policy:    &iam_model.PasswordComplexityPolicyView{},
				orgPolicy: &iam_model.OrgIAMPolicyView{},
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "create user pw policy nil",
			args: args{
				es:        GetMockManipulateUser(ctrl),
				ctx:       authz.NewMockContext("orgID", "userID"),
				user:      &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, Human: &model.Human{}},
				orgPolicy: &iam_model.OrgIAMPolicyView{},
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "create user org policy nil",
			args: args{
				es:     GetMockManipulateUser(ctrl),
				ctx:    authz.NewMockContext("orgID", "userID"),
				user:   &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, Human: &model.Human{}},
				policy: &iam_model.PasswordComplexityPolicyView{},
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.CreateUser(tt.args.ctx, tt.args.user, tt.args.policy, tt.args.orgPolicy)

			if tt.res.errFunc == nil && result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if tt.res.errFunc == nil && result.UserName != tt.res.user.UserName {
				t.Errorf("got wrong result username: expected: %v, actual: %v ", tt.res.user.UserName, result.UserName)
			}
			if tt.res.errFunc == nil && tt.res.user.Email != nil {
				if result.IsEmailVerified != tt.res.user.IsEmailVerified {
					t.Errorf("got wrong result IsEmailVerified: expected: %v, actual: %v ", tt.res.user.IsEmailVerified, result.IsEmailVerified)
				}
			}
			if tt.res.errFunc == nil && tt.res.user.Phone != nil {
				if result.IsPhoneVerified != tt.res.user.IsPhoneVerified {
					t.Errorf("got wrong result IsPhoneVerified: expected: %v, actual: %v ", tt.res.user.IsPhoneVerified, result.IsPhoneVerified)
				}
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
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
		policy        *iam_model.PasswordComplexityPolicyView
		orgPolicy     *iam_model.OrgIAMPolicyView
	}
	type res struct {
		user    *model.User
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
				es: GetMockManipulateUserWithPasswordInitCodeGen(ctrl, repo_model.User{
					UserName: "UserName",
					Human: &repo_model.Human{
						Profile: &repo_model.Profile{DisplayName: "DisplayName", FirstName: "FirstName", LastName: "LastName"},
						Email:   &repo_model.Email{EmailAddress: "EmailAddress"}}},
				),
				ctx: authz.NewMockContext("orgID", "userID"),
				user: &model.User{
					ObjectRoot: es_models.ObjectRoot{Sequence: 1},
					UserName:   "UserName",
					Human: &model.Human{
						Profile: &model.Profile{
							FirstName:   "FirstName",
							LastName:    "LastName",
							DisplayName: "DisplayName",
						},
						Email: &model.Email{
							EmailAddress: "EmailAddress",
						},
						Password: &model.Password{
							SecretString: "Password",
						},
					},
				},
				policy:        &iam_model.PasswordComplexityPolicyView{},
				orgPolicy:     &iam_model.OrgIAMPolicyView{UserLoginMustBeDomain: true},
				resourceOwner: "ResourceOwner",
			},
			res: res{
				user: &model.User{ObjectRoot: es_models.ObjectRoot{Sequence: 1},
					UserName: "UserName",
					Human: &model.Human{
						Profile: &model.Profile{
							FirstName:   "FirstName",
							LastName:    "LastName",
							DisplayName: "DisplayName",
						},
						Email: &model.Email{
							EmailAddress: "EmailAddress",
						},
					},
				},
			},
		},
		{
			name: "email as username",
			args: args{
				es: GetMockManipulateUserWithPasswordInitCodeGen(ctrl, repo_model.User{
					UserName: "UserName",
					Human: &repo_model.Human{
						Profile: &repo_model.Profile{DisplayName: "DisplayName", FirstName: "FirstName", LastName: "LastName"},
						Email:   &repo_model.Email{EmailAddress: "EmailAddress"}}},
				),
				ctx: authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{Sequence: 1},
					Human: &model.Human{
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
				},
				policy:        &iam_model.PasswordComplexityPolicyView{},
				orgPolicy:     &iam_model.OrgIAMPolicyView{UserLoginMustBeDomain: false},
				resourceOwner: "ResourceOwner",
			},
			res: res{
				user: &model.User{ObjectRoot: es_models.ObjectRoot{Sequence: 1},
					UserName: "EmailAddress",
					Human: &model.Human{
						Profile: &model.Profile{
							FirstName:   "FirstName",
							LastName:    "LastName",
							DisplayName: "DisplayName",
						},
						Email: &model.Email{
							EmailAddress: "EmailAddress",
						},
					},
				},
			},
		},
		{
			name: "invalid user",
			args: args{
				es:            GetMockManipulateUser(ctrl),
				ctx:           authz.NewMockContext("orgID", "userID"),
				user:          &model.User{ObjectRoot: es_models.ObjectRoot{Sequence: 1}, Human: &model.Human{}},
				policy:        &iam_model.PasswordComplexityPolicyView{},
				orgPolicy:     &iam_model.OrgIAMPolicyView{},
				resourceOwner: "ResourceOwner",
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "register user no password",
			args: args{
				es:  GetMockManipulateUser(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{Sequence: 1},
					UserName: "EmailAddress",
					Human: &model.Human{
						Profile: &model.Profile{
							FirstName:   "FirstName",
							LastName:    "LastName",
							DisplayName: "DisplayName",
						},
						Email: &model.Email{
							EmailAddress: "EmailAddress",
						},
					},
				},
				policy:        &iam_model.PasswordComplexityPolicyView{},
				orgPolicy:     &iam_model.OrgIAMPolicyView{},
				resourceOwner: "ResourceOwner",
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "no resourceowner",
			args: args{
				es:  GetMockManipulateUser(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{Sequence: 1},
					UserName: "EmailAddress",
					Human: &model.Human{
						Profile: &model.Profile{
							FirstName:   "FirstName",
							LastName:    "LastName",
							DisplayName: "DisplayName",
						},
						Email: &model.Email{
							EmailAddress: "EmailAddress",
						},
					},
				},
				policy:    &iam_model.PasswordComplexityPolicyView{},
				orgPolicy: &iam_model.OrgIAMPolicyView{},
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "no pw policy",
			args: args{
				es:  GetMockManipulateUser(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{Sequence: 1},
					UserName: "EmailAddress",
					Human: &model.Human{
						Profile: &model.Profile{
							FirstName:   "FirstName",
							LastName:    "LastName",
							DisplayName: "DisplayName",
						},
						Email: &model.Email{
							EmailAddress: "EmailAddress",
						},
					},
				},
				orgPolicy: &iam_model.OrgIAMPolicyView{},
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "no org policy",
			args: args{
				es:  GetMockManipulateUser(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{Sequence: 1},
					UserName: "EmailAddress",
					Human: &model.Human{
						Profile: &model.Profile{
							FirstName:   "FirstName",
							LastName:    "LastName",
							DisplayName: "DisplayName",
						},
						Email: &model.Email{
							EmailAddress: "EmailAddress",
						},
					},
				},
				policy: &iam_model.PasswordComplexityPolicyView{},
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.RegisterUser(tt.args.ctx, tt.args.user, tt.args.policy, tt.args.orgPolicy, tt.args.resourceOwner)

			if tt.res.errFunc == nil && result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if tt.res.errFunc == nil && result.UserName != tt.res.user.UserName {
				t.Errorf("got wrong result name: expected: %v, actual: %v ", tt.res.user.UserName, result.UserName)
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestDeactivateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es   *UserEventstore
		ctx  context.Context
		user *model.User
	}
	type res struct {
		user    *model.User
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
				es:   GetMockManipulateUser(ctrl),
				ctx:  authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, State: model.UserStateInactive},
			},
		},
		{
			name: "deactivate user with inactive state",
			args: args{
				es:   GetMockManipulateInactiveUser(ctrl),
				ctx:  authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing not found",
			args: args{
				es:   GetMockManipulateUserNoEvents(ctrl),
				ctx:  authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.DeactivateUser(tt.args.ctx, tt.args.user.AggregateID)

			if tt.res.errFunc == nil && result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if tt.res.errFunc == nil && result.State != tt.res.user.State {
				t.Errorf("got wrong result name: expected: %v, actual: %v ", tt.res.user.State, result.State)
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestReactivateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es   *UserEventstore
		ctx  context.Context
		user *model.User
	}
	type res struct {
		user    *model.User
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
				es:   GetMockManipulateInactiveUser(ctrl),
				ctx:  authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, State: model.UserStateActive},
			},
		},
		{
			name: "reactivate user with inital state",
			args: args{
				es:   GetMockManipulateUser(ctrl),
				ctx:  authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing user not found",
			args: args{
				es:   GetMockManipulateUserNoEvents(ctrl),
				ctx:  authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.ReactivateUser(tt.args.ctx, tt.args.user.AggregateID)

			if tt.res.errFunc == nil && result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if tt.res.errFunc == nil && result.State != tt.res.user.State {
				t.Errorf("got wrong result name: expected: %v, actual: %v ", tt.res.user.State, result.State)
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestLockUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es   *UserEventstore
		ctx  context.Context
		user *model.User
	}
	type res struct {
		user    *model.User
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
				es:   GetMockManipulateUser(ctrl),
				ctx:  authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, State: model.UserStateLocked},
			},
		},
		{
			name: "lock user with locked state",
			args: args{
				es:   GetMockManipulateLockedUser(ctrl),
				ctx:  authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing user not found",
			args: args{
				es:   GetMockManipulateUserNoEvents(ctrl),
				ctx:  authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.LockUser(tt.args.ctx, tt.args.user.AggregateID)

			if tt.res.errFunc == nil && result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if tt.res.errFunc == nil && result.State != tt.res.user.State {
				t.Errorf("got wrong result name: expected: %v, actual: %v ", tt.res.user.State, result.State)
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestUnlockUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es   *UserEventstore
		ctx  context.Context
		user *model.User
	}
	type res struct {
		user    *model.User
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "unlock user, ok",
			args: args{
				es:   GetMockManipulateLockedUser(ctrl),
				ctx:  authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, State: model.UserStateActive},
			},
		},
		{
			name: "unlock user not locked state",
			args: args{
				es:   GetMockManipulateUser(ctrl),
				ctx:  authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing user not found",
			args: args{
				es:   GetMockManipulateUserNoEvents(ctrl),
				ctx:  authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.UnlockUser(tt.args.ctx, tt.args.user.AggregateID)

			if tt.res.errFunc == nil && result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if tt.res.errFunc == nil && result.State != tt.res.user.State {
				t.Errorf("got wrong result name: expected: %v, actual: %v ", tt.res.user.State, result.State)
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestGetInitCodeByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es   *UserEventstore
		ctx  context.Context
		user *model.User
	}
	type res struct {
		code    *model.InitUserCode
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "get by id, ok",
			args: args{
				es: GetMockManipulateUserWithInitCode(ctrl,
					repo_model.User{
						ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID"},
						UserName:   "UserName",
						Human: &repo_model.Human{
							Profile: &repo_model.Profile{
								DisplayName: "DisplayName",
							},
						},
					}),
				ctx:  authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				code: &model.InitUserCode{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, Expiry: time.Hour * 30},
			},
		},
		{
			name: "empty userid",
			args: args{
				es:   GetMockManipulateUser(ctrl),
				ctx:  authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "", Sequence: 1}},
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing user not found",
			args: args{
				es:   GetMockManipulateUserNoEvents(ctrl),
				ctx:  authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.InitializeUserCodeByID(tt.args.ctx, tt.args.user.AggregateID)

			if tt.res.errFunc == nil && result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if tt.res.errFunc == nil && result == nil {
				t.Error("got wrong result code should not be nil", result)
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestCreateInitCode(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es   *UserEventstore
		ctx  context.Context
		user *model.User
	}
	type res struct {
		code    *model.InitUserCode
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "create init code",
			args: args{
				es:   GetMockManipulateUserWithPasswordInitCodeGen(ctrl, repo_model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}}),
				ctx:  authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, Human: &model.Human{}},
			},
			res: res{
				code: &model.InitUserCode{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, Expiry: time.Hour * 1},
			},
		},
		{
			name: "empty userid",
			args: args{
				es:   GetMockManipulateUser(ctrl),
				ctx:  authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "", Sequence: 1}, Human: &model.Human{}},
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing user not found",
			args: args{
				es:   GetMockManipulateUserNoEvents(ctrl),
				ctx:  authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, Human: &model.Human{}},
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.CreateInitializeUserCodeByID(tt.args.ctx, tt.args.user.AggregateID)

			if tt.res.errFunc == nil && result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if tt.res.errFunc == nil && result == nil {
				t.Errorf("got wrong result code is nil")
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestInitCodeSent(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es   *UserEventstore
		ctx  context.Context
		user *model.User
	}
	type res struct {
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "sent init",
			args: args{
				es:   GetMockManipulateUser(ctrl),
				ctx:  authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{},
		},
		{
			name: "empty userid",
			args: args{
				es:   GetMockManipulateUser(ctrl),
				ctx:  authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "", Sequence: 1}},
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing user not found",
			args: args{
				es:   GetMockManipulateUserNoEvents(ctrl),
				ctx:  authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.es.InitCodeSent(tt.args.ctx, tt.args.user.AggregateID)

			if tt.res.errFunc == nil && err != nil {
				t.Errorf("rshould not get err")
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestInitCodeVerify(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es         *UserEventstore
		ctx        context.Context
		policy     *iam_model.PasswordComplexityPolicyView
		userID     string
		verifyCode string
		password   string
	}
	type res struct {
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "verify init code, no pw",
			args: args{
				es: GetMockManipulateUserWithInitCode(ctrl,
					repo_model.User{
						ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID"},
						Human: &repo_model.Human{
							Email: &repo_model.Email{
								EmailAddress: "EmailAddress",
							},
						},
					},
				),
				ctx:        authz.NewMockContext("orgID", "userID"),
				policy:     &iam_model.PasswordComplexityPolicyView{},
				verifyCode: "code",
				userID:     "userID",
			},
		},
		{
			name: "verify init code, pw",
			args: args{
				es: GetMockManipulateUserWithInitCode(ctrl,
					repo_model.User{
						ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID"},
						Human: &repo_model.Human{
							Email: &repo_model.Email{
								EmailAddress:    "EmailAddress",
								IsEmailVerified: true,
							},
						},
					},
				),
				ctx:        authz.NewMockContext("orgID", "userID"),
				policy:     &iam_model.PasswordComplexityPolicyView{},
				userID:     "userID",
				verifyCode: "code",
				password:   "password",
			},
		},
		{
			name: "verify init code, email and pw",
			args: args{
				es: GetMockManipulateUserWithInitCode(ctrl,
					repo_model.User{
						ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID"},
						Human: &repo_model.Human{
							Email: &repo_model.Email{
								EmailAddress: "EmailAddress",
							},
						},
					},
				),
				ctx:        authz.NewMockContext("orgID", "userID"),
				policy:     &iam_model.PasswordComplexityPolicyView{},
				userID:     "userID",
				verifyCode: "code",
				password:   "password",
			},
		},
		{
			name: "empty userid",
			args: args{
				es:         GetMockManipulateUser(ctrl),
				ctx:        authz.NewMockContext("orgID", "userID"),
				policy:     &iam_model.PasswordComplexityPolicyView{},
				userID:     "",
				verifyCode: "code",
				password:   "password",
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "password policy not matched",
			args: args{
				es:         GetMockManipulateUser(ctrl),
				ctx:        authz.NewMockContext("orgID", "userID"),
				policy:     &iam_model.PasswordComplexityPolicyView{HasNumber: true},
				userID:     "userID",
				verifyCode: "code",
				password:   "password",
			},
			res: res{
				errFunc: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "existing user not found",
			args: args{
				es:         GetMockManipulateUserNoEventsWithPw(ctrl),
				ctx:        authz.NewMockContext("orgID", "userID"),
				policy:     &iam_model.PasswordComplexityPolicyView{},
				userID:     "userID",
				password:   "password",
				verifyCode: "code",
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.es.VerifyInitCode(tt.args.ctx, tt.args.policy, tt.args.userID, tt.args.verifyCode, tt.args.password)

			if tt.res.errFunc == nil && err != nil {
				t.Errorf("should not have err: %v", err)
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v", err)
			}
		})
	}
}

func TestSkipMFAInit(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es   *UserEventstore
		ctx  context.Context
		user *model.User
	}
	type res struct {
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "skip mfa init",
			args: args{
				es:   GetMockManipulateUser(ctrl),
				ctx:  authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{},
		},
		{
			name: "empty userid",
			args: args{
				es:   GetMockManipulateUser(ctrl),
				ctx:  authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "", Sequence: 1}},
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing user not found",
			args: args{
				es:   GetMockManipulateUserNoEvents(ctrl),
				ctx:  authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.es.SkipMFAInit(tt.args.ctx, tt.args.user.AggregateID)

			if tt.res.errFunc == nil && err != nil {
				t.Errorf("rshould not get err")
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestPasswordID(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es   *UserEventstore
		ctx  context.Context
		user *model.User
	}
	type res struct {
		password *model.Password
		errFunc  func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "get by id, ok",
			args: args{
				es:   GetMockManipulateUserFull(ctrl, false),
				ctx:  authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				password: &model.Password{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, ChangeRequired: true},
			},
		},
		{
			name: "empty userid",
			args: args{
				es:   GetMockManipulateUser(ctrl),
				ctx:  authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "", Sequence: 1}},
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing user not found",
			args: args{
				es:   GetMockManipulateUserNoEvents(ctrl),
				ctx:  authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
		{
			name: "existing pw not found",
			args: args{
				es:   GetMockManipulateUser(ctrl),
				ctx:  authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.UserPasswordByID(tt.args.ctx, tt.args.user.AggregateID)

			if tt.res.errFunc == nil && result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if tt.res.errFunc == nil && result.ChangeRequired != tt.res.password.ChangeRequired {
				t.Errorf("got wrong result change required: expected: %v, actual: %v ", tt.res.password.ChangeRequired, result.ChangeRequired)
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestSetOneTimePassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es       *UserEventstore
		ctx      context.Context
		policy   *iam_model.PasswordComplexityPolicyView
		password *model.Password
	}
	type res struct {
		password *model.Password
		errFunc  func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "create one time pw",
			args: args{
				es:       GetMockManipulateUserWithPasswordCodeGen(ctrl, repo_model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID"}, Human: &repo_model.Human{Email: &repo_model.Email{EmailAddress: "email"}}}, true),
				ctx:      authz.NewMockContext("orgID", "userID"),
				policy:   &iam_model.PasswordComplexityPolicyView{},
				password: &model.Password{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID"}, SecretString: "Password"},
			},
			res: res{
				password: &model.Password{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, ChangeRequired: true},
			},
		},
		{
			name: "empty userid",
			args: args{
				es:       GetMockManipulateUser(ctrl),
				ctx:      authz.NewMockContext("orgID", "userID"),
				policy:   &iam_model.PasswordComplexityPolicyView{},
				password: &model.Password{ObjectRoot: es_models.ObjectRoot{AggregateID: ""}, SecretString: "Password"},
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing user not found",
			args: args{
				es:       GetMockManipulateUserNoEvents(ctrl),
				ctx:      authz.NewMockContext("orgID", "userID"),
				policy:   &iam_model.PasswordComplexityPolicyView{},
				password: &model.Password{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID"}, SecretString: "Password"},
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.SetOneTimePassword(tt.args.ctx, tt.args.policy, tt.args.password)

			if tt.res.errFunc == nil && result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if tt.res.errFunc == nil && result.ChangeRequired != true {
				t.Errorf("should be one time")
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestCheckPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es          *UserEventstore
		ctx         context.Context
		userID      string
		password    string
		authRequest *req_model.AuthRequest
	}
	type res struct {
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "check pw ok",
			args: args{
				es: GetMockManipulateUserWithPasswordAndEmailCodeGen(ctrl,
					repo_model.User{
						ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID"},
						Human: &repo_model.Human{
							Password: &repo_model.Password{Secret: &crypto.CryptoValue{
								CryptoType: crypto.TypeHash,
								Algorithm:  "hash",
								Crypted:    []byte("password"),
							}},
						},
					},
				),
				ctx:      authz.NewMockContext("orgID", "userID"),
				userID:   "userID",
				password: "password",
				authRequest: &req_model.AuthRequest{
					ID:      "id",
					AgentID: "agentID",
					BrowserInfo: &req_model.BrowserInfo{
						UserAgent:      "user agent",
						AcceptLanguage: "accept langugage",
						RemoteIP:       net.IPv4(29, 4, 20, 19),
					},
				},
			},
			res: res{},
		},
		{
			name: "empty userid",
			args: args{
				es:       GetMockManipulateUser(ctrl),
				ctx:      authz.NewMockContext("orgID", "userID"),
				userID:   "",
				password: "password",
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing user not found",
			args: args{
				es:       GetMockManipulateUserNoEvents(ctrl),
				ctx:      authz.NewMockContext("orgID", "userID"),
				userID:   "userID",
				password: "password",
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
		{
			name: "no password",
			args: args{
				es: GetMockUserByIDOK(ctrl,
					repo_model.User{
						ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID"},
					},
				),
				ctx:      authz.NewMockContext("orgID", "userID"),
				userID:   "userID",
				password: "password",
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "invalid password",
			args: args{
				es: GetMockManipulateUserWithPasswordCodeGen(ctrl,
					repo_model.User{
						ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID"},
						Human: &repo_model.Human{
							Password: &repo_model.Password{Secret: &crypto.CryptoValue{
								CryptoType: crypto.TypeHash,
								Algorithm:  "hash",
								Crypted:    []byte("password"),
							}},
							Email: &repo_model.Email{EmailAddress: "email"},
						},
					},
					true,
				),
				ctx:      authz.NewMockContext("orgID", "userID"),
				userID:   "userID",
				password: "wrong password",
				authRequest: &req_model.AuthRequest{
					ID:      "id",
					AgentID: "agentID",
					BrowserInfo: &req_model.BrowserInfo{
						UserAgent:      "user agent",
						AcceptLanguage: "accept langugage",
						RemoteIP:       net.IPv4(29, 4, 20, 19),
					},
				},
			},
			res: res{
				errFunc: caos_errs.IsErrorInvalidArgument,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.es.CheckPassword(tt.args.ctx, tt.args.userID, tt.args.password, tt.args.authRequest)

			if tt.res.errFunc == nil && err != nil {
				t.Errorf("result has error: %v", err)
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v", err)
			}
		})
	}
}

func TestSetPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es       *UserEventstore
		ctx      context.Context
		policy   *iam_model.PasswordComplexityPolicyView
		userID   string
		code     string
		password string
	}
	type res struct {
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "create pw",
			args: args{
				es: GetMockManipulateUserWithPasswordCodeGen(ctrl,
					repo_model.User{
						ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID"},
						Human: &repo_model.Human{
							PasswordCode: &repo_model.PasswordCode{Code: &crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("code"),
							}},
							Email: &repo_model.Email{EmailAddress: "email"},
						},
					},
					true,
				),
				ctx:      authz.NewMockContext("orgID", "userID"),
				policy:   &iam_model.PasswordComplexityPolicyView{},
				userID:   "userID",
				code:     "code",
				password: "password",
			},
			res: res{},
		},
		{
			name: "empty userid",
			args: args{
				es:       GetMockManipulateUser(ctrl),
				ctx:      authz.NewMockContext("orgID", "userID"),
				policy:   &iam_model.PasswordComplexityPolicyView{},
				userID:   "",
				code:     "code",
				password: "password",
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing user not found",
			args: args{
				es:       GetMockManipulateUserNoEvents(ctrl),
				ctx:      authz.NewMockContext("orgID", "userID"),
				policy:   &iam_model.PasswordComplexityPolicyView{},
				userID:   "userID",
				code:     "code",
				password: "password",
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
		{
			name: "no passcode",
			args: args{
				es: GetMockManipulateUserWithPasswordCodeGen(ctrl,
					repo_model.User{
						ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID"},
						Human: &repo_model.Human{
							Email: &repo_model.Email{EmailAddress: "email"},
						},
					},
					true,
				),
				ctx:      authz.NewMockContext("orgID", "userID"),
				policy:   &iam_model.PasswordComplexityPolicyView{},
				userID:   "userID",
				code:     "code",
				password: "password",
			},
			res: res{
				errFunc: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "invalid passcode",
			args: args{
				es: GetMockManipulateUserWithPasswordCodeGen(ctrl,
					repo_model.User{
						ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID"},
						Human: &repo_model.Human{
							PasswordCode: &repo_model.PasswordCode{Code: &crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc2",
								KeyID:      "id",
								Crypted:    []byte("code2"),
							}},
							Email: &repo_model.Email{EmailAddress: "email"},
						},
					},
					true,
				),
				ctx:      authz.NewMockContext("orgID", "userID"),
				policy:   &iam_model.PasswordComplexityPolicyView{},
				userID:   "userID",
				code:     "code",
				password: "password",
			},
			res: res{
				errFunc: caos_errs.IsErrorInvalidArgument,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.es.SetPassword(tt.args.ctx, tt.args.policy, tt.args.userID, tt.args.code, tt.args.password, "")

			if tt.res.errFunc == nil && err != nil {
				t.Errorf("result has error: %v", err)
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v", err)
			}
		})
	}
}

func TestChangePassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es     *UserEventstore
		ctx    context.Context
		policy *iam_model.PasswordComplexityPolicyView
		userID string
		old    string
		new    string
	}
	type res struct {
		password string
		errFunc  func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "change pw",
			args: args{
				es: GetMockManipulateUserWithPasswordAndEmailCodeGen(ctrl,
					repo_model.User{
						ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID"},
						Human: &repo_model.Human{
							Password: &repo_model.Password{Secret: &crypto.CryptoValue{
								CryptoType: crypto.TypeHash,
								Algorithm:  "hash",
								Crypted:    []byte("old"),
							}},
						},
					},
				),
				ctx:    authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.PasswordComplexityPolicyView{},
				userID: "userID",
				old:    "old",
				new:    "new",
			},
			res: res{
				password: "new",
			},
		},
		{
			name: "empty userid",
			args: args{
				es:     GetMockManipulateUser(ctrl),
				ctx:    authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.PasswordComplexityPolicyView{},
				userID: "",
				old:    "old",
				new:    "new",
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing user not found",
			args: args{
				es:     GetMockManipulateUserNoEvents(ctrl),
				ctx:    authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.PasswordComplexityPolicyView{},
				userID: "userID",
				old:    "old",
				new:    "new",
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
		{
			name: "no password",
			args: args{
				es: GetMockManipulateUserWithPasswordCodeGen(ctrl,
					repo_model.User{
						ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID"},
						Human: &repo_model.Human{
							Email: &repo_model.Email{EmailAddress: "email"},
						},
					},
					true,
				),
				ctx:    authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.PasswordComplexityPolicyView{},
				userID: "userID",
				old:    "old",
				new:    "new",
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "invalid password",
			args: args{
				es: GetMockManipulateUserWithPasswordCodeGen(ctrl,
					repo_model.User{
						ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID"},
						Human: &repo_model.Human{
							Password: &repo_model.Password{Secret: &crypto.CryptoValue{
								CryptoType: crypto.TypeHash,
								Algorithm:  "hash",
								Crypted:    []byte("older"),
							}},
							Email: &repo_model.Email{EmailAddress: "email"},
						},
					},
					true,
				),
				ctx:    authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.PasswordComplexityPolicyView{},
				userID: "userID",
				old:    "old",
				new:    "new",
			},
			res: res{
				errFunc: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "no policy",
			args: args{
				es: GetMockManipulateUserWithPasswordAndEmailCodeGen(ctrl,
					repo_model.User{
						ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID"},
						Human: &repo_model.Human{
							Password: &repo_model.Password{Secret: &crypto.CryptoValue{
								CryptoType: crypto.TypeHash,
								Algorithm:  "hash",
								Crypted:    []byte("old"),
							}},
						},
					},
				),
				ctx:    authz.NewMockContext("orgID", "userID"),
				userID: "userID",
				old:    "old",
				new:    "new",
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.ChangePassword(tt.args.ctx, tt.args.policy, tt.args.userID, tt.args.old, tt.args.new, "")

			if tt.res.errFunc == nil && result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if tt.res.errFunc == nil && string(result.SecretCrypto.Crypted) != tt.res.password {
				t.Errorf("got wrong result crypted: expected: %v, actual: %v ", tt.res.password, result.SecretString)
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v", err)
			}
		})
	}
}

func TestRequestSetPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es         *UserEventstore
		ctx        context.Context
		userID     string
		notifyType model.NotificationType
	}
	type res struct {
		password *model.Password
		errFunc  func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "create pw",
			args: args{
				es: GetMockManipulateUserWithPasswordCodeGen(ctrl, repo_model.User{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID"},
					Human: &repo_model.Human{
						Email: &repo_model.Email{EmailAddress: "email"},
					}},
					true),
				ctx:        authz.NewMockContext("orgID", "userID"),
				userID:     "AggregateID",
				notifyType: model.NotificationTypeEmail,
			},
			res: res{
				password: &model.Password{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, ChangeRequired: false},
			},
		},
		{
			name: "initial state",
			args: args{
				es:         GetMockManipulateUserWithPasswordCodeGen(ctrl, repo_model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID"}, Human: &repo_model.Human{}}, false),
				ctx:        authz.NewMockContext("orgID", "userID"),
				userID:     "AggregateID",
				notifyType: model.NotificationTypeEmail,
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "empty userid",
			args: args{
				es:         GetMockManipulateUser(ctrl),
				notifyType: model.NotificationTypeEmail,
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing user not found",
			args: args{
				es:         GetMockManipulateUserNoEvents(ctrl),
				ctx:        authz.NewMockContext("orgID", "userID"),
				userID:     "AggregateID",
				notifyType: model.NotificationTypeEmail,
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.es.RequestSetPassword(tt.args.ctx, tt.args.userID, tt.args.notifyType)

			if tt.res.errFunc == nil && err != nil {
				t.Errorf("should not get err")
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestResendInitialMail(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es     *UserEventstore
		ctx    context.Context
		userID string
		mail   string
	}
	type res struct {
		password *model.Password
		errFunc  func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "resend ok",
			args: args{
				es:     GetMockManipulateUserWithInitCode(ctrl, repo_model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID"}, Human: &repo_model.Human{}}),
				ctx:    authz.NewMockContext("orgID", "userID"),
				userID: "AggregateID",
				mail:   "",
			},
			res: res{
				password: &model.Password{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, ChangeRequired: false},
			},
		},
		{
			name: "resend with email ok",
			args: args{
				es:     GetMockManipulateUserWithInitCode(ctrl, repo_model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID"}, Human: &repo_model.Human{}}),
				ctx:    authz.NewMockContext("orgID", "userID"),
				userID: "AggregateID",
				mail:   "email",
			},
			res: res{
				password: &model.Password{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, ChangeRequired: false},
			},
		},
		{
			name: "empty userid",
			args: args{
				es:   GetMockManipulateUser(ctrl),
				mail: "",
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing user not found",
			args: args{
				es:     GetMockManipulateUserNoEvents(ctrl),
				ctx:    authz.NewMockContext("orgID", "userID"),
				userID: "AggregateID",
				mail:   "",
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.es.ResendInitialMail(tt.args.ctx, tt.args.userID, tt.args.mail)

			if tt.res.errFunc == nil && err != nil {
				t.Errorf("should not get err")
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestPasswordCodeSent(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es   *UserEventstore
		ctx  context.Context
		user *model.User
	}
	type res struct {
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "sent password code",
			args: args{
				es:   GetMockManipulateUser(ctrl),
				ctx:  authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, Human: &model.Human{}},
			},
			res: res{},
		},
		{
			name: "empty userid",
			args: args{
				es:   GetMockManipulateUser(ctrl),
				ctx:  authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "", Sequence: 1}},
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing user not found",
			args: args{
				es:   GetMockManipulateUserNoEvents(ctrl),
				ctx:  authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.es.PasswordCodeSent(tt.args.ctx, tt.args.user.AggregateID)

			if tt.res.errFunc == nil && err != nil {
				t.Errorf("rshould not get err")
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestAddExternalIDP(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es          *UserEventstore
		ctx         context.Context
		externalIDP *model.ExternalIDP
	}
	type res struct {
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "add ok",
			args: args{
				es:  GetMockManipulateUser(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				externalIDP: &model.ExternalIDP{
					ObjectRoot: es_models.ObjectRoot{
						AggregateID: "AggregateID",
					},
					IDPConfigID: "IDPConfigID",
					UserID:      "UserID",
					DisplayName: "DisplayName",
				},
			},
		},
		{
			name: "invalid idp",
			args: args{
				es:  GetMockManipulateUser(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				externalIDP: &model.ExternalIDP{
					ObjectRoot: es_models.ObjectRoot{
						AggregateID: "AggregateID",
					},
					UserID:      "UserID",
					DisplayName: "DisplayName",
				},
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing user not found",
			args: args{
				es:  GetMockManipulateUserNoEvents(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				externalIDP: &model.ExternalIDP{
					ObjectRoot: es_models.ObjectRoot{
						AggregateID: "AggregateID",
					},
					IDPConfigID: "IDPConfigID",
					UserID:      "UserID",
					DisplayName: "DisplayName",
				},
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.AddExternalIDP(tt.args.ctx, tt.args.externalIDP)

			if tt.res.errFunc == nil && result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if tt.res.errFunc == nil && result.IDPConfigID == "" {
				t.Errorf("result has no idpconfig")
			}
			if tt.res.errFunc == nil && result.UserID == "" {
				t.Errorf("result has no UserID")
			}
			if tt.res.errFunc == nil && result == nil {
				t.Errorf("got wrong result change required: actual: %v ", result)
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestRemoveExternalIDP(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es          *UserEventstore
		ctx         context.Context
		externalIDP *model.ExternalIDP
	}
	type res struct {
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "remove ok",
			args: args{
				es:  GetMockManipulateUserWithExternalIDP(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				externalIDP: &model.ExternalIDP{
					ObjectRoot: es_models.ObjectRoot{
						AggregateID: "AggregateID",
					},
					IDPConfigID: "IDPConfigID",
					UserID:      "UserID",
				},
			},
		},
		{
			name: "invalid idp",
			args: args{
				es:  GetMockManipulateUser(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				externalIDP: &model.ExternalIDP{
					ObjectRoot: es_models.ObjectRoot{
						AggregateID: "AggregateID",
					},
					UserID:      "UserID",
					DisplayName: "DisplayName",
				},
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "remove external idp not existing",
			args: args{
				es:  GetMockManipulateUser(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				externalIDP: &model.ExternalIDP{
					ObjectRoot: es_models.ObjectRoot{
						AggregateID: "AggregateID",
					},
					IDPConfigID: "IDPConfigID",
					UserID:      "UserID",
				},
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing user not found",
			args: args{
				es:  GetMockManipulateUserNoEvents(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				externalIDP: &model.ExternalIDP{
					ObjectRoot: es_models.ObjectRoot{
						AggregateID: "AggregateID",
					},
					IDPConfigID: "IDPConfigID",
					UserID:      "UserID",
					DisplayName: "DisplayName",
				},
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.es.RemoveExternalIDP(tt.args.ctx, tt.args.externalIDP)

			if tt.res.errFunc == nil && err != nil {
				t.Errorf("should not get err, %v", err)
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestProfileByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es   *UserEventstore
		ctx  context.Context
		user *model.User
	}
	type res struct {
		profile *model.Profile
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "get by id, ok",
			args: args{
				es:   GetMockManipulateUserFull(ctrl, false),
				ctx:  authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, Human: &model.Human{}},
			},
			res: res{
				profile: &model.Profile{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, DisplayName: "DisplayName"},
			},
		},
		{
			name: "empty userid",
			args: args{
				es:   GetMockManipulateUser(ctrl),
				ctx:  authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "", Sequence: 1}},
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing user not found",
			args: args{
				es:   GetMockManipulateUserNoEvents(ctrl),
				ctx:  authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profile, err := tt.args.es.ProfileByID(tt.args.ctx, tt.args.user.AggregateID)

			if tt.res.errFunc == nil && profile.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if tt.res.errFunc == nil && profile.DisplayName != tt.res.profile.DisplayName {
				t.Errorf("got wrong result change required: expected: %v, actual: %v ", tt.res.profile.DisplayName, profile.DisplayName)
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestChangeProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es      *UserEventstore
		ctx     context.Context
		profile *model.Profile
	}
	type res struct {
		profile *model.Profile
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "get by id, ok",
			args: args{
				es:      GetMockManipulateUserFull(ctrl, false),
				ctx:     authz.NewMockContext("orgID", "userID"),
				profile: &model.Profile{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, FirstName: "FirstName Changed", LastName: "LastName Changed"},
			},
			res: res{
				profile: &model.Profile{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, FirstName: "FirstName Changed", LastName: "LastName Changed"},
			},
		},
		{
			name: "invalid profile",
			args: args{
				es:      GetMockManipulateUser(ctrl),
				ctx:     authz.NewMockContext("orgID", "userID"),
				profile: &model.Profile{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing user not found",
			args: args{
				es:      GetMockManipulateUserNoEvents(ctrl),
				ctx:     authz.NewMockContext("orgID", "userID"),
				profile: &model.Profile{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, FirstName: "FirstName Changed", LastName: "LastName Changed"},
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.ChangeProfile(tt.args.ctx, tt.args.profile)

			if tt.res.errFunc == nil && result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if tt.res.errFunc == nil && result.FirstName != tt.res.profile.FirstName {
				t.Errorf("got wrong result change required: expected: %v, actual: %v ", tt.res.profile.FirstName, result.FirstName)
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestEmailByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es   *UserEventstore
		ctx  context.Context
		user *model.User
	}
	type res struct {
		email   *model.Email
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "get by id, ok",
			args: args{
				es:   GetMockManipulateUserFull(ctrl, false),
				ctx:  authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				email: &model.Email{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, EmailAddress: "EmailAddress"},
			},
		},
		{
			name: "empty userid",
			args: args{
				es:   GetMockManipulateUser(ctrl),
				ctx:  authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "", Sequence: 1}},
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing user not found",
			args: args{
				es:   GetMockManipulateUserNoEvents(ctrl),
				ctx:  authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.EmailByID(tt.args.ctx, tt.args.user.AggregateID)

			if tt.res.errFunc == nil && result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if tt.res.errFunc == nil && result.EmailAddress != tt.res.email.EmailAddress {
				t.Errorf("got wrong result change required: expected: %v, actual: %v ", tt.res.email.EmailAddress, result.EmailAddress)
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestChangeEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es    *UserEventstore
		ctx   context.Context
		email *model.Email
	}
	type res struct {
		email   *model.Email
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "change email address, verified",
			args: args{
				es:    GetMockManipulateUserFull(ctrl, true),
				ctx:   authz.NewMockContext("orgID", "userID"),
				email: &model.Email{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, EmailAddress: "EmailAddressChanged", IsEmailVerified: true},
			},
			res: res{
				email: &model.Email{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, EmailAddress: "EmailAddressChanged", IsEmailVerified: true},
			},
		},
		{
			name: "change email not verified, getting code",
			args: args{
				es: GetMockManipulateUserWithEmailCodeGen(ctrl, repo_model.User{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					UserName:   "UserName",
					Human: &repo_model.Human{
						Profile: &repo_model.Profile{DisplayName: "DisplayName"},
						Email:   &repo_model.Email{EmailAddress: "EmailAddress"},
					},
				}, true),
				ctx:   authz.NewMockContext("orgID", "userID"),
				email: &model.Email{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, EmailAddress: "EmailAddressChanged", IsEmailVerified: false},
			},
			res: res{
				email: &model.Email{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, EmailAddress: "EmailAddressChanged", IsEmailVerified: false},
			},
		},
		{
			name: "user state initial",
			args: args{
				es: GetMockManipulateUserWithEmailCodeGen(ctrl, repo_model.User{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					UserName:   "UserName",
					Human: &repo_model.Human{
						Profile: &repo_model.Profile{DisplayName: "DisplayName"},
						Email:   &repo_model.Email{EmailAddress: "EmailAddress"},
					},
				}, false),
				ctx:   authz.NewMockContext("orgID", "userID"),
				email: &model.Email{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, EmailAddress: "EmailAddressChanged", IsEmailVerified: false},
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "empty userid",
			args: args{
				es:    GetMockManipulateUser(ctrl),
				ctx:   authz.NewMockContext("orgID", "userID"),
				email: &model.Email{ObjectRoot: es_models.ObjectRoot{AggregateID: "", Sequence: 1}, EmailAddress: "EmailAddressChanged", IsEmailVerified: true},
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing user not found",
			args: args{
				es:    GetMockManipulateUserNoEvents(ctrl),
				ctx:   authz.NewMockContext("orgID", "userID"),
				email: &model.Email{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, EmailAddress: "EmailAddressChanged", IsEmailVerified: true},
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.ChangeEmail(tt.args.ctx, tt.args.email)
			if (tt.res.errFunc != nil && !tt.res.errFunc(err)) || (tt.res.errFunc == nil && err != nil) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.errFunc == nil && result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if tt.res.errFunc == nil && result.EmailAddress != tt.res.email.EmailAddress {
				t.Errorf("got wrong result change required: expected: %v, actual: %v ", tt.res.email.EmailAddress, result.EmailAddress)
			}
			if tt.res.errFunc == nil && result.IsEmailVerified != tt.res.email.IsEmailVerified {
				t.Errorf("got wrong result change required: expected: %v, actual: %v ", tt.res.email.IsEmailVerified, result.IsEmailVerified)
			}
		})
	}
}

func TestVerifyEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es     *UserEventstore
		ctx    context.Context
		userID string
		code   string
	}
	type res struct {
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "verify email code ok",
			args: args{
				es:     GetMockManipulateUserWithEmailCode(ctrl),
				ctx:    authz.NewMockContext("orgID", "userID"),
				userID: "AggregateID",
				code:   "code",
			},
			res: res{},
		},
		{
			name: "verify email code wrong",
			args: args{
				es:     GetMockManipulateUserWithEmailCode(ctrl),
				ctx:    authz.NewMockContext("orgID", "userID"),
				userID: "AggregateID",
				code:   "wrong",
			},
			res: res{
				errFunc: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "empty userid",
			args: args{
				es:   GetMockManipulateUser(ctrl),
				ctx:  authz.NewMockContext("orgID", "userID"),
				code: "code",
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "empty code",
			args: args{
				es:     GetMockManipulateUser(ctrl),
				ctx:    authz.NewMockContext("orgID", "userID"),
				userID: "AggregateID",
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing user not found",
			args: args{
				es:     GetMockManipulateUserNoEvents(ctrl),
				ctx:    authz.NewMockContext("orgID", "userID"),
				userID: "AggregateID",
				code:   "code",
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.es.VerifyEmail(tt.args.ctx, tt.args.userID, tt.args.code)

			if tt.res.errFunc == nil && err != nil {
				t.Errorf("should not get err %v", err)
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestCreateEmailVerificationCode(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es     *UserEventstore
		ctx    context.Context
		userID string
	}
	type res struct {
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "create email verification code ok",
			args: args{
				es: GetMockManipulateUserWithEmailCodeGen(ctrl, repo_model.User{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					UserName:   "UserName",
					Human: &repo_model.Human{
						Profile: &repo_model.Profile{DisplayName: "DisplayName"},
						Email:   &repo_model.Email{EmailAddress: "EmailAddress"},
					}}, true),
				ctx:    authz.NewMockContext("orgID", "userID"),
				userID: "userID",
			},
			res: res{},
		},
		{
			name: "initial state",
			args: args{
				es: GetMockManipulateUserWithEmailCodeGen(ctrl, repo_model.User{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					UserName:   "UserName",
					Human: &repo_model.Human{
						Profile: &repo_model.Profile{DisplayName: "DisplayName"},
						Email:   &repo_model.Email{EmailAddress: "EmailAddress"},
					},
				}, false),
				ctx:    authz.NewMockContext("orgID", "userID"),
				userID: "userID",
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "empty userid",
			args: args{
				es:  GetMockManipulateUser(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing user not found",
			args: args{
				es:     GetMockManipulateUserNoEvents(ctrl),
				ctx:    authz.NewMockContext("orgID", "userID"),
				userID: "userID",
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
		{
			name: "no email found",
			args: args{
				es:     GetMockManipulateUser(ctrl),
				ctx:    authz.NewMockContext("orgID", "userID"),
				userID: "userID",
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "already verified",
			args: args{
				es:     GetMockManipulateUserVerifiedEmail(ctrl),
				ctx:    authz.NewMockContext("orgID", "userID"),
				userID: "userID",
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.es.CreateEmailVerificationCode(tt.args.ctx, tt.args.userID)

			if tt.res.errFunc == nil && err != nil {
				t.Errorf("should not get err, got: %v", err)
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestEmailVerificationCodeSent(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es   *UserEventstore
		ctx  context.Context
		user *model.User
	}
	type res struct {
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "sent email verify code",
			args: args{
				es:   GetMockManipulateUser(ctrl),
				ctx:  authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{},
		},
		{
			name: "empty userid",
			args: args{
				es:   GetMockManipulateUser(ctrl),
				ctx:  authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "", Sequence: 1}},
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing user not found",
			args: args{
				es:   GetMockManipulateUserNoEvents(ctrl),
				ctx:  authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.es.EmailVerificationCodeSent(tt.args.ctx, tt.args.user.AggregateID)

			if tt.res.errFunc == nil && err != nil {
				t.Errorf("rshould not get err")
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestPhoneByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es   *UserEventstore
		ctx  context.Context
		user *model.User
	}
	type res struct {
		phone   *model.Phone
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "get by id, ok",
			args: args{
				es:   GetMockManipulateUserFull(ctrl, false),
				ctx:  authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				phone: &model.Phone{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, PhoneNumber: "PhoneNumber"},
			},
		},
		{
			name: "empty userid",
			args: args{
				es:   GetMockManipulateUser(ctrl),
				ctx:  authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "", Sequence: 1}},
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing user not found",
			args: args{
				es:   GetMockManipulateUserNoEvents(ctrl),
				ctx:  authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.PhoneByID(tt.args.ctx, tt.args.user.AggregateID)

			if tt.res.errFunc == nil && result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if tt.res.errFunc == nil && result.PhoneNumber != tt.res.phone.PhoneNumber {
				t.Errorf("got wrong result change required: expected: %v, actual: %v ", tt.res.phone.PhoneNumber, result.PhoneNumber)
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestChangePhone(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es    *UserEventstore
		ctx   context.Context
		phone *model.Phone
	}
	type res struct {
		phone   *model.Phone
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "change phone, verified",
			args: args{
				es:    GetMockManipulateUserFull(ctrl, false),
				ctx:   authz.NewMockContext("orgID", "userID"),
				phone: &model.Phone{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, PhoneNumber: "0711234567", IsPhoneVerified: true},
			},
			res: res{
				phone: &model.Phone{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, PhoneNumber: "+41711234567", IsPhoneVerified: true},
			},
		},
		{
			name: "change phone not verified, getting code",
			args: args{
				es: GetMockManipulateUserWithPhoneCodeGen(ctrl, repo_model.User{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					UserName:   "UserName",
					Human: &repo_model.Human{
						Profile: &repo_model.Profile{DisplayName: "DisplayName"},
						Phone:   &repo_model.Phone{PhoneNumber: "PhoneNumber"}},
				}),
				ctx:   authz.NewMockContext("orgID", "userID"),
				phone: &model.Phone{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, PhoneNumber: "+41711234567", IsPhoneVerified: false},
			},
			res: res{
				phone: &model.Phone{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, PhoneNumber: "+41711234567", IsPhoneVerified: false},
			},
		},
		{
			name: "empty userid",
			args: args{
				es:    GetMockManipulateUser(ctrl),
				ctx:   authz.NewMockContext("orgID", "userID"),
				phone: &model.Phone{ObjectRoot: es_models.ObjectRoot{AggregateID: "", Sequence: 1}, PhoneNumber: "PhoneNumberChanged", IsPhoneVerified: true},
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing user not found",
			args: args{
				es:    GetMockManipulateUserNoEvents(ctrl),
				ctx:   authz.NewMockContext("orgID", "userID"),
				phone: &model.Phone{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, PhoneNumber: "+41711234567", IsPhoneVerified: true},
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.ChangePhone(tt.args.ctx, tt.args.phone)

			if tt.res.errFunc == nil && result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if tt.res.errFunc == nil && result.PhoneNumber != tt.res.phone.PhoneNumber {
				t.Errorf("got wrong result change required: expected: %v, actual: %v ", tt.res.phone.PhoneNumber, result.PhoneNumber)
			}
			if tt.res.errFunc == nil && result.IsPhoneVerified != tt.res.phone.IsPhoneVerified {
				t.Errorf("got wrong result change required: expected: %v, actual: %v ", tt.res.phone.IsPhoneVerified, result.IsPhoneVerified)
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestVerifyPhone(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es     *UserEventstore
		ctx    context.Context
		userID string
		code   string
	}
	type res struct {
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "verify code ok",
			args: args{
				es:     GetMockManipulateUserWithPhoneCode(ctrl),
				ctx:    authz.NewMockContext("orgID", "userID"),
				userID: "AggregateID",
				code:   "code",
			},
			res: res{},
		},
		{
			name: "verify code wrong",
			args: args{
				es:     GetMockManipulateUserWithPhoneCode(ctrl),
				ctx:    authz.NewMockContext("orgID", "userID"),
				userID: "AggregateID",
				code:   "wrong",
			},
			res: res{
				errFunc: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "empty userid",
			args: args{
				es:   GetMockManipulateUser(ctrl),
				ctx:  authz.NewMockContext("orgID", "userID"),
				code: "Code",
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "empty code",
			args: args{
				es:     GetMockManipulateUser(ctrl),
				ctx:    authz.NewMockContext("orgID", "userID"),
				userID: "AggregateID",
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing user not found",
			args: args{
				es:     GetMockManipulateUserNoEvents(ctrl),
				ctx:    authz.NewMockContext("orgID", "userID"),
				userID: "AggregateID",
				code:   "Code",
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.es.VerifyPhone(tt.args.ctx, tt.args.userID, tt.args.code)

			if tt.res.errFunc == nil && err != nil {
				t.Errorf("should not get err %v", err)
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestRemovePhone(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es     *UserEventstore
		ctx    context.Context
		userID string
		code   string
	}
	type res struct {
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "remove phone ok",
			args: args{
				es:     GetMockManipulateUserVerifiedPhone(ctrl),
				ctx:    authz.NewMockContext("orgID", "userID"),
				userID: "AggregateID",
				code:   "code",
			},
			res: res{},
		},
		{
			name: "empty userid",
			args: args{
				es:   GetMockManipulateUser(ctrl),
				ctx:  authz.NewMockContext("orgID", "userID"),
				code: "Code",
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing user not found",
			args: args{
				es:     GetMockManipulateUserNoEvents(ctrl),
				ctx:    authz.NewMockContext("orgID", "userID"),
				userID: "AggregateID",
				code:   "Code",
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.es.RemovePhone(tt.args.ctx, tt.args.userID)

			if tt.res.errFunc == nil && err != nil {
				t.Errorf("should not get err %v", err)
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestCreatePhoneVerificationCode(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es     *UserEventstore
		ctx    context.Context
		userID string
	}
	type res struct {
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "create phone verification code okk",
			args: args{
				es: GetMockManipulateUserWithPhoneCodeGen(ctrl, repo_model.User{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					UserName:   "UserName",
					Human: &repo_model.Human{
						Profile: &repo_model.Profile{DisplayName: "DisplayName"},
						Phone:   &repo_model.Phone{PhoneNumber: "PhoneNumber"},
					}}),
				ctx:    authz.NewMockContext("orgID", "userID"),
				userID: "userID",
			},
			res: res{},
		},
		{
			name: "no userid",
			args: args{
				es:  GetMockManipulateUser(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing user not found",
			args: args{
				es:     GetMockManipulateUserNoEvents(ctrl),
				ctx:    authz.NewMockContext("orgID", "userID"),
				userID: "userID",
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
		{
			name: "no phone found",
			args: args{
				es:     GetMockManipulateUser(ctrl),
				ctx:    authz.NewMockContext("orgID", "userID"),
				userID: "userID",
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "already verified",
			args: args{
				es:     GetMockManipulateUserVerifiedPhone(ctrl),
				ctx:    authz.NewMockContext("orgID", "userID"),
				userID: "userID",
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.es.CreatePhoneVerificationCode(tt.args.ctx, tt.args.userID)

			if tt.res.errFunc == nil && err != nil {
				t.Errorf("should not ger err")
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestPhoneVerificationCodeSent(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es   *UserEventstore
		ctx  context.Context
		user *model.User
	}
	type res struct {
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "sent phone verification code",
			args: args{
				es:   GetMockManipulateUser(ctrl),
				ctx:  authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{},
		},
		{
			name: "empty userid",
			args: args{
				es:   GetMockManipulateUser(ctrl),
				ctx:  authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "", Sequence: 1}},
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing user not found",
			args: args{
				es:   GetMockManipulateUserNoEvents(ctrl),
				ctx:  authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.es.PhoneVerificationCodeSent(tt.args.ctx, tt.args.user.AggregateID)

			if tt.res.errFunc == nil && err != nil {
				t.Errorf("rshould not get err")
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestAddressByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es   *UserEventstore
		ctx  context.Context
		user *model.User
	}
	type res struct {
		address *model.Address
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "get by id, ok",
			args: args{
				es:   GetMockManipulateUserFull(ctrl, false),
				ctx:  authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				address: &model.Address{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, Country: "Country"},
			},
		},
		{
			name: "empty userid",
			args: args{
				es:   GetMockManipulateUser(ctrl),
				ctx:  authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "", Sequence: 1}},
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing user not found",
			args: args{
				es:   GetMockManipulateUserNoEvents(ctrl),
				ctx:  authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.AddressByID(tt.args.ctx, tt.args.user.AggregateID)

			if tt.res.errFunc == nil && result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if tt.res.errFunc == nil && result.Country != tt.res.address.Country {
				t.Errorf("got wrong result change required: expected: %v, actual: %v ", tt.res.address.Country, result.Country)
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestChangeAddress(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es      *UserEventstore
		ctx     context.Context
		address *model.Address
	}
	type res struct {
		address *model.Address
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "change address ok",
			args: args{
				es:      GetMockManipulateUserFull(ctrl, false),
				ctx:     authz.NewMockContext("orgID", "userID"),
				address: &model.Address{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, Country: "CountryChanged"},
			},
			res: res{
				address: &model.Address{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, Country: "CountryChanged"},
			},
		},
		{
			name: "empty userid",
			args: args{
				es:      GetMockManipulateUser(ctrl),
				ctx:     authz.NewMockContext("orgID", "userID"),
				address: &model.Address{ObjectRoot: es_models.ObjectRoot{AggregateID: "", Sequence: 1}, Country: "CountryChanged"},
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing user not found",
			args: args{
				es:      GetMockManipulateUserNoEvents(ctrl),
				ctx:     authz.NewMockContext("orgID", "userID"),
				address: &model.Address{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, Country: "CountryCountry"},
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.ChangeAddress(tt.args.ctx, tt.args.address)

			if tt.res.errFunc == nil && result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if tt.res.errFunc == nil && result.Country != tt.res.address.Country {
				t.Errorf("got wrong result change required: expected: %v, actual: %v ", tt.res.address.Country, result.Country)
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestAddOTP(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es          *UserEventstore
		ctx         context.Context
		userID      string
		accountName string
	}
	type res struct {
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "add ok",
			args: args{
				es:          GetMockManipulateUserWithOTPGen(ctrl),
				ctx:         authz.NewMockContext("orgID", "userID"),
				userID:      "AggregateID",
				accountName: "AccountName",
			},
		},
		{
			name: "empty userid",
			args: args{
				es:     GetMockManipulateUser(ctrl),
				ctx:    authz.NewMockContext("orgID", "userID"),
				userID: "",
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing user not found",
			args: args{
				es:     GetMockManipulateUserNoEvents(ctrl),
				ctx:    authz.NewMockContext("orgID", "userID"),
				userID: "AggregateID",
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.AddOTP(tt.args.ctx, tt.args.userID, tt.args.accountName)

			if tt.res.errFunc == nil && result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if tt.res.errFunc == nil && result.Url == "" {
				t.Errorf("result has no url")
			}
			if tt.res.errFunc == nil && result.SecretString == "" {
				t.Errorf("result has no url")
			}
			if tt.res.errFunc == nil && result == nil {
				t.Errorf("got wrong result change required: actual: %v ", result)
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestCheckMFAOTPSetup(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es     *UserEventstore
		ctx    context.Context
		userID string
		code   string
	}
	type res struct {
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "setup ok",
			args: args{
				es: func() *UserEventstore {
					es := GetMockManipulateUserWithOTP(ctrl, true, false)
					es.validateTOTP = func(string, string) bool {
						return true
					}
					return es
				}(),
				ctx:    authz.NewMockContext("orgID", "userID"),
				userID: "id",
				code:   "code",
			},
			res: res{},
		},
		{
			name: "wrong code",
			args: args{
				es: func() *UserEventstore {
					es := GetMockManipulateUserWithOTP(ctrl, true, false)
					es.validateTOTP = func(string, string) bool {
						return false
					}
					return es
				}(),
				ctx:    authz.NewMockContext("orgID", "userID"),
				userID: "id",
				code:   "code",
			},
			res: res{
				errFunc: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "empty userid",
			args: args{
				es:   GetMockManipulateUser(ctrl),
				ctx:  authz.NewMockContext("orgID", "userID"),
				code: "code",
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "empty code",
			args: args{
				es:     GetMockManipulateUser(ctrl),
				ctx:    authz.NewMockContext("orgID", "userID"),
				userID: "userID",
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing user not found",
			args: args{
				es:     GetMockManipulateUserNoEvents(ctrl),
				ctx:    authz.NewMockContext("orgID", "userID"),
				userID: "userID",
				code:   "code",
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
		{
			name: "user has no otp",
			args: args{
				es:     GetMockManipulateUser(ctrl),
				ctx:    authz.NewMockContext("orgID", "userID"),
				userID: "userID",
				code:   "code",
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.es.CheckMFAOTPSetup(tt.args.ctx, tt.args.userID, tt.args.code, "")

			if tt.res.errFunc == nil && err != nil {
				t.Errorf("result should not get err")
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestCheckMFAOTP(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es          *UserEventstore
		ctx         context.Context
		userID      string
		code        string
		authRequest *req_model.AuthRequest
	}
	type res struct {
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "check ok",
			args: args{
				es: func() *UserEventstore {
					es := GetMockManipulateUserWithOTP(ctrl, true, true)
					es.validateTOTP = func(string, string) bool {
						return true
					}
					return es
				}(),
				ctx:    authz.NewMockContext("orgID", "userID"),
				userID: "id",
				code:   "code",
				authRequest: &req_model.AuthRequest{
					ID:      "id",
					AgentID: "agentID",
					BrowserInfo: &req_model.BrowserInfo{
						UserAgent:      "user agent",
						AcceptLanguage: "accept langugage",
						RemoteIP:       net.IPv4(29, 4, 20, 19),
					},
				},
			},
			res: res{},
		},
		{
			name: "wrong code",
			args: args{
				es: func() *UserEventstore {
					es := GetMockManipulateUserWithOTP(ctrl, true, true)
					es.validateTOTP = func(string, string) bool {
						return false
					}
					return es
				}(),
				ctx:    authz.NewMockContext("orgID", "userID"),
				userID: "id",
				code:   "code",
				authRequest: &req_model.AuthRequest{
					ID:      "id",
					AgentID: "agentID",
					BrowserInfo: &req_model.BrowserInfo{
						UserAgent:      "user agent",
						AcceptLanguage: "accept langugage",
						RemoteIP:       net.IPv4(29, 4, 20, 19),
					},
				},
			},
			res: res{
				errFunc: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "empty userid",
			args: args{
				es:   GetMockManipulateUser(ctrl),
				ctx:  authz.NewMockContext("orgID", "userID"),
				code: "code",
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "empty code",
			args: args{
				es:     GetMockManipulateUser(ctrl),
				ctx:    authz.NewMockContext("orgID", "userID"),
				userID: "userID",
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing user not found",
			args: args{
				es:     GetMockManipulateUserNoEvents(ctrl),
				ctx:    authz.NewMockContext("orgID", "userID"),
				userID: "userID",
				code:   "code",
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
		{
			name: "user has no otp",
			args: args{
				es:     GetMockManipulateUser(ctrl),
				ctx:    authz.NewMockContext("orgID", "userID"),
				userID: "userID",
				code:   "code",
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.es.CheckMFAOTP(tt.args.ctx, tt.args.userID, tt.args.code, tt.args.authRequest)

			if tt.res.errFunc == nil && err != nil {
				t.Errorf("result should not get err, got : %v", err)
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestRemoveOTP(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es   *UserEventstore
		ctx  context.Context
		user *model.User
	}
	type res struct {
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "remove ok",
			args: args{
				es:   GetMockManipulateUserWithOTP(ctrl, false, true),
				ctx:  authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
		},
		{
			name: "empty userid",
			args: args{
				es:   GetMockManipulateUser(ctrl),
				ctx:  authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "", Sequence: 1}},
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing user not found",
			args: args{
				es:   GetMockManipulateUserNoEvents(ctrl),
				ctx:  authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
		{
			name: "user has no otp",
			args: args{
				es:   GetMockManipulateUser(ctrl),
				ctx:  authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.es.RemoveOTP(tt.args.ctx, tt.args.user.AggregateID)

			if tt.res.errFunc == nil && err != nil {
				t.Errorf("result should not get err")
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestChangesUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es           *UserEventstore
		id           string
		secID        string
		lastSequence uint64
		limit        uint64
	}
	type res struct {
		changes *model.UserChanges
		profile *model.Profile
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "changes from events, ok",
			args: args{
				es:           GetMockChangesUserOK(ctrl),
				id:           "1",
				lastSequence: 0,
				limit:        0,
			},
			res: res{
				changes: &model.UserChanges{Changes: []*model.UserChange{{EventType: "", Sequence: 1, ModifierID: ""}}, LastSequence: 1},
				profile: &model.Profile{FirstName: "Hans", LastName: "Muster"},
			},
		},
		{
			name: "changes from events, no events",
			args: args{
				es:           GetMockChangesUserNoEvents(ctrl),
				id:           "2",
				lastSequence: 0,
				limit:        0,
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.UserChanges(nil, tt.args.id, tt.args.lastSequence, tt.args.limit, false)

			profile := &model.Profile{}
			if result != nil && len(result.Changes) > 0 {
				b, err := json.Marshal(result.Changes[0].Data)
				json.Unmarshal(b, profile)
				if err != nil {
				}
			}
			if !tt.res.wantErr && result.LastSequence != tt.res.changes.LastSequence && profile.DisplayName != tt.res.profile.DisplayName {
				t.Errorf("got wrong result name: expected: %v, actual: %v ", tt.res.changes.LastSequence, result.LastSequence)
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}
