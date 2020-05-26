package eventsourcing

import (
	"context"
	"testing"
	"time"

	"github.com/caos/zitadel/internal/api/auth"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
)

func TestUserByIDQuery(t *testing.T) {
	type args struct {
		id       string
		sequence uint64
	}
	type res struct {
		filterLen int
		wantErr   bool
		errFunc   func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "user by id query ok",
			args: args{
				id:       "ID",
				sequence: 1,
			},
			res: res{
				filterLen: 3,
			},
		},
		{
			name: "no id",
			args: args{
				sequence: 1,
			},
			res: res{
				filterLen: 3,
				wantErr:   true,
				errFunc:   caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query, err := UserByIDQuery(tt.args.id, tt.args.sequence)
			if !tt.res.wantErr && query == nil {
				t.Errorf("query should not be nil")
			}
			if !tt.res.wantErr && len(query.Filters) != tt.res.filterLen {
				t.Errorf("got wrong filter len: expected: %v, actual: %v ", tt.res.filterLen, len(query.Filters))
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestUserQuery(t *testing.T) {
	type args struct {
		sequence uint64
	}
	type res struct {
		filterLen int
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "user query ok",
			args: args{
				sequence: 1,
			},
			res: res{
				filterLen: 2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := UserQuery(tt.args.sequence)
			if query == nil {
				t.Errorf("query should not be nil")
			}
			if len(query.Filters) != tt.res.filterLen {
				t.Errorf("got wrong filter len: expected: %v, actual: %v ", tt.res.filterLen, len(query.Filters))
			}
		})
	}
}

func TestUserCreateAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		new        *model.User
		initCode   *model.InitUserCode
		phoneCode  *model.PhoneCode
		aggCreator *models.AggregateCreator
	}
	type res struct {
		eventLen   int
		eventTypes []models.EventType
		checkData  []bool
		wantErr    bool
		errFunc    func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "user create aggregate ok",
			args: args{
				ctx: auth.NewMockContext("orgID", "userID"),
				new: &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					Profile: &model.Profile{UserName: "UserName"},
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{model.UserAdded},
				checkData:  []bool{true},
			},
		},
		{
			name: "new user nil",
			args: args{
				ctx:        auth.NewMockContext("orgID", "userID"),
				new:        nil,
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "create with init code",
			args: args{
				ctx: auth.NewMockContext("orgID", "userID"),
				new: &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					Profile: &model.Profile{UserName: "UserName"},
				},
				initCode:   &model.InitUserCode{},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   2,
				eventTypes: []models.EventType{model.UserAdded, model.InitializedUserCodeAdded},
				checkData:  []bool{true, true},
			},
		},
		{
			name: "create with phone code",
			args: args{
				ctx: auth.NewMockContext("orgID", "userID"),
				new: &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					Profile: &model.Profile{UserName: "UserName"},
				},
				phoneCode:  &model.PhoneCode{},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   2,
				eventTypes: []models.EventType{model.UserAdded, model.UserPhoneCodeAdded},
				checkData:  []bool{true, true},
			},
		},
		{
			name: "create with email verified",
			args: args{
				ctx: auth.NewMockContext("orgID", "userID"),
				new: &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					Profile: &model.Profile{UserName: "UserName"},
					Email:   &model.Email{EmailAddress: "EmailAddress", IsEmailVerified: true},
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   2,
				eventTypes: []models.EventType{model.UserAdded, model.UserEmailVerified},
				checkData:  []bool{true, false},
			},
		},
		{
			name: "create with phone verified",
			args: args{
				ctx: auth.NewMockContext("orgID", "userID"),
				new: &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					Profile: &model.Profile{UserName: "UserName"},
					Phone:   &model.Phone{PhoneNumber: "PhoneNumber", IsPhoneVerified: true},
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   2,
				eventTypes: []models.EventType{model.UserAdded, model.UserPhoneVerified},
				checkData:  []bool{true, false},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := UserCreateAggregate(tt.args.ctx, tt.args.aggCreator, tt.args.new, tt.args.initCode, tt.args.phoneCode, "")

			if !tt.res.wantErr && len(agg.Events) != tt.res.eventLen {
				t.Errorf("got wrong event len: expected: %v, actual: %v ", tt.res.eventLen, len(agg.Events))
			}
			for i := 0; i < tt.res.eventLen; i++ {
				if !tt.res.wantErr && agg.Events[i].Type != tt.res.eventTypes[i] {
					t.Errorf("got wrong event type: expected: %v, actual: %v ", tt.res.eventTypes[i], agg.Events[i].Type.String())
				}
				if !tt.res.wantErr && tt.res.checkData[i] && agg.Events[i].Data == nil {
					t.Errorf("should have data in event")
				}
				if !tt.res.wantErr && !tt.res.checkData[i] && agg.Events[i].Data != nil {
					t.Errorf("should not have data in event")
				}
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestUserRegisterAggregate(t *testing.T) {
	type args struct {
		ctx           context.Context
		new           *model.User
		emailCode     *model.EmailCode
		resourceOwner string
		aggCreator    *models.AggregateCreator
	}
	type res struct {
		eventLen   int
		eventTypes []models.EventType
		errFunc    func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "user register aggregate ok",
			args: args{
				ctx: auth.NewMockContext("", ""),
				new: &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					Profile: &model.Profile{UserName: "UserName"},
				},
				emailCode:     &model.EmailCode{},
				resourceOwner: "newResourceowner",
				aggCreator:    models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   2,
				eventTypes: []models.EventType{model.UserRegistered, model.UserEmailCodeAdded},
			},
		},
		{
			name: "new user nil",
			args: args{
				ctx:           auth.NewMockContext("orgID", "userID"),
				new:           nil,
				emailCode:     &model.EmailCode{},
				resourceOwner: "newResourceowner",
				aggCreator:    models.NewAggregateCreator("Test"),
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "code nil",
			args: args{
				ctx:           auth.NewMockContext("orgID", "userID"),
				resourceOwner: "newResourceowner",
				new: &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					Profile: &model.Profile{UserName: "UserName"},
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "create with email code",
			args: args{
				ctx: auth.NewMockContext("orgID", "userID"),
				new: &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					Profile: &model.Profile{UserName: "UserName"},
				},
				resourceOwner: "newResourceowner",
				emailCode:     &model.EmailCode{},
				aggCreator:    models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   2,
				eventTypes: []models.EventType{model.UserRegistered, model.UserEmailCodeAdded},
			},
		},
		{
			name: "create no resourceowner",
			args: args{
				ctx: auth.NewMockContext("orgID", "userID"),
				new: &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					Profile: &model.Profile{UserName: "UserName"},
				},
				emailCode:  &model.EmailCode{},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := UserRegisterAggregate(tt.args.ctx, tt.args.aggCreator, tt.args.new, tt.args.resourceOwner, tt.args.emailCode)

			if tt.res.errFunc == nil && len(agg.Events) != tt.res.eventLen {
				t.Errorf("got wrong event len: expected: %v, actual: %v ", tt.res.eventLen, len(agg.Events))
			}
			for i := 0; i < tt.res.eventLen; i++ {
				if tt.res.errFunc == nil && agg.Events[i].Type != tt.res.eventTypes[i] {
					t.Errorf("got wrong event type: expected: %v, actual: %v ", tt.res.eventTypes[i], agg.Events[i].Type.String())
				}
				if tt.res.errFunc == nil && agg.Events[i].Data == nil {
					t.Errorf("should have data in event")
				}
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestUserDeactivateAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		new        *model.User
		aggCreator *models.AggregateCreator
	}
	type res struct {
		eventLen  int
		eventType models.EventType
		errFunc   func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "user deactivate aggregate ok",
			args: args{
				ctx: auth.NewMockContext("orgID", "userID"),
				new: &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					Profile: &model.Profile{UserName: "UserName"},
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.UserDeactivated,
			},
		},
		{
			name: "new user nil",
			args: args{
				ctx:        auth.NewMockContext("orgID", "userID"),
				new:        nil,
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := UserDeactivateAggregate(tt.args.aggCreator, tt.args.new)(tt.args.ctx)

			if tt.res.errFunc == nil && len(agg.Events) != tt.res.eventLen {
				t.Errorf("got wrong event len: expected: %v, actual: %v ", tt.res.eventLen, len(agg.Events))
			}
			if tt.res.errFunc == nil && agg.Events[0].Type != tt.res.eventType {
				t.Errorf("got wrong event type: expected: %v, actual: %v ", tt.res.eventType, agg.Events[0].Type.String())
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestUserReactivateAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		new        *model.User
		aggCreator *models.AggregateCreator
	}
	type res struct {
		eventLen  int
		eventType models.EventType
		errFunc   func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "user reactivate aggregate ok",
			args: args{
				ctx: auth.NewMockContext("orgID", "userID"),
				new: &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					Profile: &model.Profile{UserName: "UserName"},
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.UserReactivated,
			},
		},
		{
			name: "new user nil",
			args: args{
				ctx:        auth.NewMockContext("orgID", "userID"),
				new:        nil,
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := UserReactivateAggregate(tt.args.aggCreator, tt.args.new)(tt.args.ctx)

			if tt.res.errFunc == nil && len(agg.Events) != tt.res.eventLen {
				t.Errorf("got wrong event len: expected: %v, actual: %v ", tt.res.eventLen, len(agg.Events))
			}
			if tt.res.errFunc == nil && agg.Events[0].Type != tt.res.eventType {
				t.Errorf("got wrong event type: expected: %v, actual: %v ", tt.res.eventType, agg.Events[0].Type.String())
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestUserLockedAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		new        *model.User
		aggCreator *models.AggregateCreator
	}
	type res struct {
		eventLen  int
		eventType models.EventType
		errFunc   func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "user locked aggregate ok",
			args: args{
				ctx: auth.NewMockContext("orgID", "userID"),
				new: &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					Profile: &model.Profile{UserName: "UserName"},
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.UserLocked,
			},
		},
		{
			name: "new user nil",
			args: args{
				ctx:        auth.NewMockContext("orgID", "userID"),
				new:        nil,
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := UserLockAggregate(tt.args.aggCreator, tt.args.new)(tt.args.ctx)

			if tt.res.errFunc == nil && len(agg.Events) != tt.res.eventLen {
				t.Errorf("got wrong event len: expected: %v, actual: %v ", tt.res.eventLen, len(agg.Events))
			}
			if tt.res.errFunc == nil && agg.Events[0].Type != tt.res.eventType {
				t.Errorf("got wrong event type: expected: %v, actual: %v ", tt.res.eventType, agg.Events[0].Type.String())
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestUserUnlockedAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		new        *model.User
		aggCreator *models.AggregateCreator
	}
	type res struct {
		eventLen  int
		eventType models.EventType
		errFunc   func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "user unlocked aggregate ok",
			args: args{
				ctx: auth.NewMockContext("orgID", "userID"),
				new: &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					Profile: &model.Profile{UserName: "UserName"},
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.UserUnlocked,
			},
		},
		{
			name: "new user nil",
			args: args{
				ctx:        auth.NewMockContext("orgID", "userID"),
				new:        nil,
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := UserUnlockAggregate(tt.args.aggCreator, tt.args.new)(tt.args.ctx)

			if tt.res.errFunc == nil && len(agg.Events) != tt.res.eventLen {
				t.Errorf("got wrong event len: expected: %v, actual: %v ", tt.res.eventLen, len(agg.Events))
			}
			if tt.res.errFunc == nil && agg.Events[0].Type != tt.res.eventType {
				t.Errorf("got wrong event type: expected: %v, actual: %v ", tt.res.eventType, agg.Events[0].Type.String())
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestUserInitCodeAggregate(t *testing.T) {
	type args struct {
		ctx      context.Context
		existing *model.User
		code     *model.InitUserCode

		aggCreator *models.AggregateCreator
	}
	type res struct {
		eventLen  int
		eventType models.EventType
		errFunc   func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "user unlocked aggregate ok",
			args: args{
				ctx: auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					Profile: &model.Profile{UserName: "UserName"},
				},
				code:       &model.InitUserCode{Expiry: time.Hour * 1},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.InitializedUserCodeAdded,
			},
		},
		{
			name: "code nil",
			args: args{
				ctx: auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					Profile: &model.Profile{UserName: "UserName"},
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := UserInitCodeAggregate(tt.args.aggCreator, tt.args.existing, tt.args.code)(tt.args.ctx)

			if tt.res.errFunc == nil && len(agg.Events) != tt.res.eventLen {
				t.Errorf("got wrong event len: expected: %v, actual: %v ", tt.res.eventLen, len(agg.Events))
			}
			if tt.res.errFunc == nil && agg.Events[0].Type != tt.res.eventType {
				t.Errorf("got wrong event type: expected: %v, actual: %v ", tt.res.eventType, agg.Events[0].Type.String())
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestInitCodeSentAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		existing   *model.User
		aggCreator *models.AggregateCreator
	}
	type res struct {
		eventLen   int
		eventTypes []models.EventType
		errFunc    func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "user init code sent aggregate ok",
			args: args{
				ctx:        auth.NewMockContext("orgID", "userID"),
				existing:   &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"}},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{model.InitializedUserCodeSent},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := UserInitCodeSentAggregate(tt.args.aggCreator, tt.args.existing)(tt.args.ctx)

			if tt.res.errFunc == nil && len(agg.Events) != tt.res.eventLen {
				t.Errorf("got wrong event len: expected: %v, actual: %v ", tt.res.eventLen, len(agg.Events))
			}
			for i := 0; i < tt.res.eventLen; i++ {
				if tt.res.errFunc == nil && agg.Events[i].Type != tt.res.eventTypes[i] {
					t.Errorf("got wrong event type: expected: %v, actual: %v ", tt.res.eventTypes[i], agg.Events[i].Type.String())
				}
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestInitCodeVerifiedAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		existing   *model.User
		password   *model.Password
		aggCreator *models.AggregateCreator
	}
	type res struct {
		eventLen   int
		eventTypes []models.EventType
		errFunc    func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "user init code only email verify",
			args: args{
				ctx: auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					Profile: &model.Profile{UserName: "UserName"},
					Email:   &model.Email{EmailAddress: "EmailAddress"},
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{model.UserEmailVerified},
			},
		},
		{
			name: "user init code only password",
			args: args{
				ctx: auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					Profile: &model.Profile{UserName: "UserName"},
					Email:   &model.Email{EmailAddress: "EmailAddress", IsEmailVerified: true},
				},
				password:   &model.Password{ChangeRequired: false},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{model.UserPasswordChanged},
			},
		},
		{
			name: "user init code email and pw",
			args: args{
				ctx: auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					Profile: &model.Profile{UserName: "UserName"},
					Email:   &model.Email{EmailAddress: "EmailAddress"},
				},
				password:   &model.Password{ChangeRequired: false},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   2,
				eventTypes: []models.EventType{model.UserEmailVerified, model.UserPasswordChanged},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := InitCodeVerifiedAggregate(tt.args.aggCreator, tt.args.existing, tt.args.password)(tt.args.ctx)

			if tt.res.errFunc == nil && len(agg.Events) != tt.res.eventLen {
				t.Errorf("got wrong event len: expected: %v, actual: %v ", tt.res.eventLen, len(agg.Events))
			}
			for i := 0; i < tt.res.eventLen; i++ {
				if tt.res.errFunc == nil && agg.Events[i].Type != tt.res.eventTypes[i] {
					t.Errorf("got wrong event type: expected: %v, actual: %v ", tt.res.eventTypes[i], agg.Events[i].Type.String())
				}
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestSkipMfaAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		existing   *model.User
		aggCreator *models.AggregateCreator
	}
	type res struct {
		eventLen  int
		eventType models.EventType
		errFunc   func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "mfa skipped init ok",
			args: args{
				ctx: auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					Profile: &model.Profile{UserName: "UserName"},
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.MfaInitSkipped,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := SkipMfaAggregate(tt.args.aggCreator, tt.args.existing)(tt.args.ctx)

			if tt.res.errFunc == nil && len(agg.Events) != tt.res.eventLen {
				t.Errorf("got wrong event len: expected: %v, actual: %v ", tt.res.eventLen, len(agg.Events))
			}
			if tt.res.errFunc == nil && agg.Events[0].Type != tt.res.eventType {
				t.Errorf("got wrong event type: expected: %v, actual: %v ", tt.res.eventType, agg.Events[0].Type.String())
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestChangePasswordAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		existing   *model.User
		password   *model.Password
		aggCreator *models.AggregateCreator
	}
	type res struct {
		eventLen  int
		eventType models.EventType
		errFunc   func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "user password aggregate ok",
			args: args{
				ctx: auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					Profile: &model.Profile{UserName: "UserName"},
				},
				password:   &model.Password{ChangeRequired: true},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.UserPasswordChanged,
			},
		},
		{
			name: "password nil",
			args: args{
				ctx: auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					Profile: &model.Profile{UserName: "UserName"},
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := PasswordChangeAggregate(tt.args.aggCreator, tt.args.existing, tt.args.password)(tt.args.ctx)

			if tt.res.errFunc == nil && len(agg.Events) != tt.res.eventLen {
				t.Errorf("got wrong event len: expected: %v, actual: %v ", tt.res.eventLen, len(agg.Events))
			}
			if tt.res.errFunc == nil && agg.Events[0].Type != tt.res.eventType {
				t.Errorf("got wrong event type: expected: %v, actual: %v ", tt.res.eventType, agg.Events[0].Type.String())
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestRequestSetPasswordAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		existing   *model.User
		request    *model.PasswordCode
		aggCreator *models.AggregateCreator
	}
	type res struct {
		eventLen  int
		eventType models.EventType
		errFunc   func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "user password aggregate ok",
			args: args{
				ctx: auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					Profile: &model.Profile{UserName: "UserName"},
				},
				request:    &model.PasswordCode{Expiry: time.Hour * 1},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.UserPasswordCodeAdded,
			},
		},
		{
			name: "request nil",
			args: args{
				ctx: auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					Profile: &model.Profile{UserName: "UserName"},
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := RequestSetPassword(tt.args.aggCreator, tt.args.existing, tt.args.request)(tt.args.ctx)

			if tt.res.errFunc == nil && len(agg.Events) != tt.res.eventLen {
				t.Errorf("got wrong event len: expected: %v, actual: %v ", tt.res.eventLen, len(agg.Events))
			}
			if tt.res.errFunc == nil && agg.Events[0].Type != tt.res.eventType {
				t.Errorf("got wrong event type: expected: %v, actual: %v ", tt.res.eventType, agg.Events[0].Type.String())
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestPasswordCodeSentAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		existing   *model.User
		aggCreator *models.AggregateCreator
	}
	type res struct {
		eventLen   int
		eventTypes []models.EventType
		errFunc    func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "user password code sent aggregate ok",
			args: args{
				ctx:        auth.NewMockContext("orgID", "userID"),
				existing:   &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"}},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{model.UserPasswordCodeSent},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := PasswordCodeSentAggregate(tt.args.aggCreator, tt.args.existing)(tt.args.ctx)

			if tt.res.errFunc == nil && len(agg.Events) != tt.res.eventLen {
				t.Errorf("got wrong event len: expected: %v, actual: %v ", tt.res.eventLen, len(agg.Events))
			}
			for i := 0; i < tt.res.eventLen; i++ {
				if tt.res.errFunc == nil && agg.Events[i].Type != tt.res.eventTypes[i] {
					t.Errorf("got wrong event type: expected: %v, actual: %v ", tt.res.eventTypes[i], agg.Events[i].Type.String())
				}
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestChangeProfileAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		existing   *model.User
		profile    *model.Profile
		aggCreator *models.AggregateCreator
	}
	type res struct {
		eventLen  int
		eventType models.EventType
		errFunc   func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "user profile aggregate ok",
			args: args{
				ctx: auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					Profile: &model.Profile{FirstName: "FirstName"},
				},
				profile:    &model.Profile{FirstName: "FirstNameChanged"},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.UserProfileChanged,
			},
		},
		{
			name: "profile nil",
			args: args{
				ctx: auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					Profile: &model.Profile{UserName: "UserName"},
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := ProfileChangeAggregate(tt.args.aggCreator, tt.args.existing, tt.args.profile)(tt.args.ctx)

			if tt.res.errFunc == nil && len(agg.Events) != tt.res.eventLen {
				t.Errorf("got wrong event len: expected: %v, actual: %v ", tt.res.eventLen, len(agg.Events))
			}
			if tt.res.errFunc == nil && agg.Events[0].Type != tt.res.eventType {
				t.Errorf("got wrong event type: expected: %v, actual: %v ", tt.res.eventType, agg.Events[0].Type.String())
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestChangeEmailAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		existing   *model.User
		email      *model.Email
		code       *model.EmailCode
		aggCreator *models.AggregateCreator
	}
	type res struct {
		eventLen   int
		eventTypes []models.EventType
		errFunc    func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "change email aggregate, verified email",
			args: args{
				ctx: auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					Email: &model.Email{EmailAddress: "EmailAddress"},
				},
				email:      &model.Email{EmailAddress: "Changed", IsEmailVerified: true},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   2,
				eventTypes: []models.EventType{model.UserEmailChanged, model.UserEmailVerified},
			},
		},
		{
			name: "with code",
			args: args{
				ctx: auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					Email: &model.Email{EmailAddress: "EmailAddress"},
				},
				email:      &model.Email{EmailAddress: "Changed"},
				code:       &model.EmailCode{Expiry: time.Hour * 1},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   2,
				eventTypes: []models.EventType{model.UserEmailChanged, model.UserEmailCodeAdded},
			},
		},
		{
			name: "email nil",
			args: args{
				ctx: auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					Email: &model.Email{EmailAddress: "Changed"},
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "email verified and code not nil",
			args: args{
				ctx: auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					Email: &model.Email{EmailAddress: "Changed", IsEmailVerified: true},
				},
				code:       &model.EmailCode{Expiry: time.Hour * 1},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "email not verified and code nil",
			args: args{
				ctx: auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					Email: &model.Email{EmailAddress: "Changed"},
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := EmailChangeAggregate(tt.args.aggCreator, tt.args.existing, tt.args.email, tt.args.code)(tt.args.ctx)

			if tt.res.errFunc == nil && len(agg.Events) != tt.res.eventLen {
				t.Errorf("got wrong event len: expected: %v, actual: %v ", tt.res.eventLen, len(agg.Events))
			}
			for i := 0; i < tt.res.eventLen; i++ {
				if tt.res.errFunc == nil && agg.Events[i].Type != tt.res.eventTypes[i] {
					t.Errorf("got wrong event type: expected: %v, actual: %v ", tt.res.eventTypes[i], agg.Events[i].Type.String())
				}
				if tt.res.errFunc == nil && agg.Events[i].Data == nil {
					t.Errorf("should have data in event")
				}
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestVerifyEmailAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		existing   *model.User
		aggCreator *models.AggregateCreator
	}
	type res struct {
		eventLen   int
		eventTypes []models.EventType
		errFunc    func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "user email verified aggregate ok",
			args: args{
				ctx:        auth.NewMockContext("orgID", "userID"),
				existing:   &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"}},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{model.UserEmailVerified},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := EmailVerifiedAggregate(tt.args.aggCreator, tt.args.existing)(tt.args.ctx)

			if tt.res.errFunc == nil && len(agg.Events) != tt.res.eventLen {
				t.Errorf("got wrong event len: expected: %v, actual: %v ", tt.res.eventLen, len(agg.Events))
			}
			for i := 0; i < tt.res.eventLen; i++ {
				if tt.res.errFunc == nil && agg.Events[i].Type != tt.res.eventTypes[i] {
					t.Errorf("got wrong event type: expected: %v, actual: %v ", tt.res.eventTypes[i], agg.Events[i].Type.String())
				}
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestVerificationFailedEmailAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		existing   *model.User
		aggCreator *models.AggregateCreator
	}
	type res struct {
		eventLen   int
		eventTypes []models.EventType
		errFunc    func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "user email verification failed aggregate ok",
			args: args{
				ctx:        auth.NewMockContext("orgID", "userID"),
				existing:   &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"}},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{model.UserEmailVerificationFailed},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := EmailVerificationFailedAggregate(tt.args.aggCreator, tt.args.existing)(tt.args.ctx)

			if tt.res.errFunc == nil && len(agg.Events) != tt.res.eventLen {
				t.Errorf("got wrong event len: expected: %v, actual: %v ", tt.res.eventLen, len(agg.Events))
			}
			for i := 0; i < tt.res.eventLen; i++ {
				if tt.res.errFunc == nil && agg.Events[i].Type != tt.res.eventTypes[i] {
					t.Errorf("got wrong event type: expected: %v, actual: %v ", tt.res.eventTypes[i], agg.Events[i].Type.String())
				}
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestCreateEmailCodeAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		existing   *model.User
		code       *model.EmailCode
		aggCreator *models.AggregateCreator
	}
	type res struct {
		eventLen   int
		eventTypes []models.EventType
		errFunc    func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "with code",
			args: args{
				ctx: auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					Email: &model.Email{EmailAddress: "EmailAddress"},
				},
				code:       &model.EmailCode{Expiry: time.Hour * 1},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{model.UserEmailCodeAdded},
			},
		},
		{
			name: "code nil",
			args: args{
				ctx: auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					Email: &model.Email{EmailAddress: "Changed"},
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := EmailVerificationCodeAggregate(tt.args.aggCreator, tt.args.existing, tt.args.code)(tt.args.ctx)

			if tt.res.errFunc == nil && len(agg.Events) != tt.res.eventLen {
				t.Errorf("got wrong event len: expected: %v, actual: %v ", tt.res.eventLen, len(agg.Events))
			}
			for i := 0; i < tt.res.eventLen; i++ {
				if tt.res.errFunc == nil && agg.Events[i].Type != tt.res.eventTypes[i] {
					t.Errorf("got wrong event type: expected: %v, actual: %v ", tt.res.eventTypes[i], agg.Events[i].Type.String())
				}
				if tt.res.errFunc == nil && agg.Events[i].Data == nil {
					t.Errorf("should have data in event")
				}
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestEmailCodeSentAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		existing   *model.User
		aggCreator *models.AggregateCreator
	}
	type res struct {
		eventLen   int
		eventTypes []models.EventType
		errFunc    func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "user email code sent aggregate ok",
			args: args{
				ctx:        auth.NewMockContext("orgID", "userID"),
				existing:   &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"}},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{model.UserEmailCodeSent},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := EmailCodeSentAggregate(tt.args.aggCreator, tt.args.existing)(tt.args.ctx)

			if tt.res.errFunc == nil && len(agg.Events) != tt.res.eventLen {
				t.Errorf("got wrong event len: expected: %v, actual: %v ", tt.res.eventLen, len(agg.Events))
			}
			for i := 0; i < tt.res.eventLen; i++ {
				if tt.res.errFunc == nil && agg.Events[i].Type != tt.res.eventTypes[i] {
					t.Errorf("got wrong event type: expected: %v, actual: %v ", tt.res.eventTypes[i], agg.Events[i].Type.String())
				}
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestChangePhoneAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		existing   *model.User
		phone      *model.Phone
		code       *model.PhoneCode
		aggCreator *models.AggregateCreator
	}
	type res struct {
		eventLen   int
		eventTypes []models.EventType
		errFunc    func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "phone change aggregate verified phone",
			args: args{
				ctx: auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					Phone: &model.Phone{PhoneNumber: "+41791234567"},
				},
				phone:      &model.Phone{PhoneNumber: "+41799876543", IsPhoneVerified: true},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   2,
				eventTypes: []models.EventType{model.UserPhoneChanged, model.UserPhoneVerified},
			},
		},
		{
			name: "with code",
			args: args{
				ctx: auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					Phone: &model.Phone{PhoneNumber: "PhoneNumber"},
				},
				phone:      &model.Phone{PhoneNumber: "Changed"},
				code:       &model.PhoneCode{Expiry: time.Hour * 1},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   2,
				eventTypes: []models.EventType{model.UserPhoneChanged, model.UserPhoneCodeAdded},
			},
		},
		{
			name: "phone nil",
			args: args{
				ctx: auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					Phone: &model.Phone{PhoneNumber: "Changed"},
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "phone verified and code not nil",
			args: args{
				ctx: auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					Phone: &model.Phone{PhoneNumber: "Changed"},
				},
				phone:      &model.Phone{PhoneNumber: "Changed", IsPhoneVerified: true},
				code:       &model.PhoneCode{Expiry: time.Hour * 1},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "phone not verified and code nil",
			args: args{
				ctx: auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					Phone: &model.Phone{PhoneNumber: "Changed"},
				},
				phone:      &model.Phone{PhoneNumber: "Changed"},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := PhoneChangeAggregate(tt.args.aggCreator, tt.args.existing, tt.args.phone, tt.args.code)(tt.args.ctx)

			if tt.res.errFunc == nil && len(agg.Events) != tt.res.eventLen {
				t.Errorf("got wrong event len: expected: %v, actual: %v ", tt.res.eventLen, len(agg.Events))
			}
			for i := 0; i < tt.res.eventLen; i++ {
				if tt.res.errFunc == nil && agg.Events[i].Type != tt.res.eventTypes[i] {
					t.Errorf("got wrong event type: expected: %v, actual: %v ", tt.res.eventTypes[i], agg.Events[i].Type.String())
				}
				if tt.res.errFunc == nil && agg.Events[i].Data == nil {
					t.Errorf("should have data in event")
				}
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestVerifyPhoneAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		existing   *model.User
		aggCreator *models.AggregateCreator
	}
	type res struct {
		eventLen   int
		eventTypes []models.EventType
		errFunc    func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "user phone verified aggregate ok",
			args: args{
				ctx: auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					Phone: &model.Phone{PhoneNumber: "PhoneNumber"},
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{model.UserPhoneVerified},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := PhoneVerifiedAggregate(tt.args.aggCreator, tt.args.existing)(tt.args.ctx)

			if tt.res.errFunc == nil && len(agg.Events) != tt.res.eventLen {
				t.Errorf("got wrong event len: expected: %v, actual: %v ", tt.res.eventLen, len(agg.Events))
			}
			for i := 0; i < tt.res.eventLen; i++ {
				if tt.res.errFunc == nil && agg.Events[i].Type != tt.res.eventTypes[i] {
					t.Errorf("got wrong event type: expected: %v, actual: %v ", tt.res.eventTypes[i], agg.Events[i].Type.String())
				}
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestVerificationFailedPhoneAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		existing   *model.User
		aggCreator *models.AggregateCreator
	}
	type res struct {
		eventLen   int
		eventTypes []models.EventType
		errFunc    func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "user phone verification failed aggregate ok",
			args: args{
				ctx: auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					Phone: &model.Phone{PhoneNumber: "PhoneNumber"},
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{model.UserPhoneVerificationFailed},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := PhoneVerificationFailedAggregate(tt.args.aggCreator, tt.args.existing)(tt.args.ctx)

			if tt.res.errFunc == nil && len(agg.Events) != tt.res.eventLen {
				t.Errorf("got wrong event len: expected: %v, actual: %v ", tt.res.eventLen, len(agg.Events))
			}
			for i := 0; i < tt.res.eventLen; i++ {
				if tt.res.errFunc == nil && agg.Events[i].Type != tt.res.eventTypes[i] {
					t.Errorf("got wrong event type: expected: %v, actual: %v ", tt.res.eventTypes[i], agg.Events[i].Type.String())
				}
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestCreatePhoneCodeAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		existing   *model.User
		code       *model.PhoneCode
		aggCreator *models.AggregateCreator
	}
	type res struct {
		eventLen   int
		eventTypes []models.EventType
		errFunc    func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "with code",
			args: args{
				ctx: auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					Email: &model.Email{EmailAddress: "EmailAddress"},
				},
				code:       &model.PhoneCode{Expiry: time.Hour * 1},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{model.UserPhoneCodeAdded},
			},
		},
		{
			name: "code nil",
			args: args{
				ctx: auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					Email: &model.Email{EmailAddress: "Changed"},
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := PhoneVerificationCodeAggregate(tt.args.aggCreator, tt.args.existing, tt.args.code)(tt.args.ctx)

			if tt.res.errFunc == nil && len(agg.Events) != tt.res.eventLen {
				t.Errorf("got wrong event len: expected: %v, actual: %v ", tt.res.eventLen, len(agg.Events))
			}
			for i := 0; i < tt.res.eventLen; i++ {
				if tt.res.errFunc == nil && agg.Events[i].Type != tt.res.eventTypes[i] {
					t.Errorf("got wrong event type: expected: %v, actual: %v ", tt.res.eventTypes[i], agg.Events[i].Type.String())
				}
				if tt.res.errFunc == nil && agg.Events[i].Data == nil {
					t.Errorf("should have data in event")
				}
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestPhoneCodeSentAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		existing   *model.User
		aggCreator *models.AggregateCreator
	}
	type res struct {
		eventLen   int
		eventTypes []models.EventType
		errFunc    func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "user phone code sent aggregate ok",
			args: args{
				ctx:        auth.NewMockContext("orgID", "userID"),
				existing:   &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"}},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{model.UserPhoneCodeSent},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := PhoneCodeSentAggregate(tt.args.aggCreator, tt.args.existing)(tt.args.ctx)

			if tt.res.errFunc == nil && len(agg.Events) != tt.res.eventLen {
				t.Errorf("got wrong event len: expected: %v, actual: %v ", tt.res.eventLen, len(agg.Events))
			}
			for i := 0; i < tt.res.eventLen; i++ {
				if tt.res.errFunc == nil && agg.Events[i].Type != tt.res.eventTypes[i] {
					t.Errorf("got wrong event type: expected: %v, actual: %v ", tt.res.eventTypes[i], agg.Events[i].Type.String())
				}
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestChangeAddressAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		existing   *model.User
		address    *model.Address
		aggCreator *models.AggregateCreator
	}
	type res struct {
		eventLen   int
		eventTypes []models.EventType
		errFunc    func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "user address change aggregate ok",
			args: args{
				ctx: auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					Address: &model.Address{Locality: "Locality"},
				},
				address:    &model.Address{Locality: "Changed"},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{model.UserAddressChanged},
			},
		},
		{
			name: "address nil",
			args: args{
				ctx: auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					Address: &model.Address{Locality: "Changed"},
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := AddressChangeAggregate(tt.args.aggCreator, tt.args.existing, tt.args.address)(tt.args.ctx)

			if tt.res.errFunc == nil && len(agg.Events) != tt.res.eventLen {
				t.Errorf("got wrong event len: expected: %v, actual: %v ", tt.res.eventLen, len(agg.Events))
			}
			for i := 0; i < tt.res.eventLen; i++ {
				if tt.res.errFunc == nil && agg.Events[i].Type != tt.res.eventTypes[i] {
					t.Errorf("got wrong event type: expected: %v, actual: %v ", tt.res.eventTypes[i], agg.Events[i].Type.String())
				}
				if tt.res.errFunc == nil && agg.Events[i].Data == nil {
					t.Errorf("should have data in event")
				}
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestOtpAddAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		existing   *model.User
		otp        *model.OTP
		aggCreator *models.AggregateCreator
	}
	type res struct {
		eventLen   int
		eventTypes []models.EventType
		errFunc    func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "user otp change aggregate ok",
			args: args{
				ctx:        auth.NewMockContext("orgID", "userID"),
				existing:   &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"}},
				otp:        &model.OTP{},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{model.MfaOtpAdded},
			},
		},
		{
			name: "otp nil",
			args: args{
				ctx:        auth.NewMockContext("orgID", "userID"),
				existing:   &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"}},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := MfaOTPAddAggregate(tt.args.aggCreator, tt.args.existing, tt.args.otp)(tt.args.ctx)

			if tt.res.errFunc == nil && len(agg.Events) != tt.res.eventLen {
				t.Errorf("got wrong event len: expected: %v, actual: %v ", tt.res.eventLen, len(agg.Events))
			}
			for i := 0; i < tt.res.eventLen; i++ {
				if tt.res.errFunc == nil && agg.Events[i].Type != tt.res.eventTypes[i] {
					t.Errorf("got wrong event type: expected: %v, actual: %v ", tt.res.eventTypes[i], agg.Events[i].Type.String())
				}
				if tt.res.errFunc == nil && agg.Events[i].Data == nil {
					t.Errorf("should have data in event")
				}
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestOtpVerifyAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		existing   *model.User
		aggCreator *models.AggregateCreator
	}
	type res struct {
		eventLen   int
		eventTypes []models.EventType
		errFunc    func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "user otp change aggregate ok",
			args: args{
				ctx:        auth.NewMockContext("orgID", "userID"),
				existing:   &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"}},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{model.MfaOtpVerified},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := MfaOTPVerifyAggregate(tt.args.aggCreator, tt.args.existing)(tt.args.ctx)

			if tt.res.errFunc == nil && len(agg.Events) != tt.res.eventLen {
				t.Errorf("got wrong event len: expected: %v, actual: %v ", tt.res.eventLen, len(agg.Events))
			}
			for i := 0; i < tt.res.eventLen; i++ {
				if tt.res.errFunc == nil && agg.Events[i].Type != tt.res.eventTypes[i] {
					t.Errorf("got wrong event type: expected: %v, actual: %v ", tt.res.eventTypes[i], agg.Events[i].Type.String())
				}
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestOtpRemoveAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		existing   *model.User
		aggCreator *models.AggregateCreator
	}
	type res struct {
		eventLen   int
		eventTypes []models.EventType
		errFunc    func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "user otp change aggregate ok",
			args: args{
				ctx:        auth.NewMockContext("orgID", "userID"),
				existing:   &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"}},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{model.MfaOtpRemoved},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := MfaOTPRemoveAggregate(tt.args.aggCreator, tt.args.existing)(tt.args.ctx)

			if tt.res.errFunc == nil && len(agg.Events) != tt.res.eventLen {
				t.Errorf("got wrong event len: expected: %v, actual: %v ", tt.res.eventLen, len(agg.Events))
			}
			for i := 0; i < tt.res.eventLen; i++ {
				if tt.res.errFunc == nil && agg.Events[i].Type != tt.res.eventTypes[i] {
					t.Errorf("got wrong event type: expected: %v, actual: %v ", tt.res.eventTypes[i], agg.Events[i].Type.String())
				}
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}
