package eventsourcing

import (
	"context"
	"github.com/caos/zitadel/internal/api/auth"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
	"testing"
	"time"
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
				eventTypes: []models.EventType{usr_model.UserAdded},
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
				eventLen:   1,
				eventTypes: []models.EventType{usr_model.UserAdded},
				wantErr:    true,
				errFunc:    caos_errs.IsPreconditionFailed,
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
				eventTypes: []models.EventType{usr_model.UserAdded, usr_model.InitializedUserCodeCreated},
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
				eventTypes: []models.EventType{usr_model.UserAdded, usr_model.UserPhoneCodeAdded},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := UserCreateAggregate(tt.args.aggCreator, tt.args.new, tt.args.initCode, tt.args.phoneCode)(tt.args.ctx)

			if !tt.res.wantErr && len(agg.Events) != tt.res.eventLen {
				t.Errorf("got wrong event len: expected: %v, actual: %v ", tt.res.eventLen, len(agg.Events))
			}
			for i := 0; i < tt.res.eventLen; i++ {
				if !tt.res.wantErr && agg.Events[i].Type != tt.res.eventTypes[i] {
					t.Errorf("got wrong event type: expected: %v, actual: %v ", tt.res.eventTypes[i], agg.Events[i].Type.String())
				}
				if !tt.res.wantErr && agg.Events[i].Data == nil {
					t.Errorf("should have data in event")
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
				resourceOwner: "newResourceowner",
				aggCreator:    models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{usr_model.UserRegistered},
			},
		},
		{
			name: "new user nil",
			args: args{
				ctx:           auth.NewMockContext("orgID", "userID"),
				new:           nil,
				resourceOwner: "newResourceowner",
				aggCreator:    models.NewAggregateCreator("Test"),
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
				eventTypes: []models.EventType{usr_model.UserRegistered, usr_model.UserEmailCodeAdded},
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
			agg, err := UserRegisterAggregate(tt.args.aggCreator, tt.args.new, tt.args.resourceOwner, tt.args.emailCode)(tt.args.ctx)

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
				eventType: usr_model.UserDeactivated,
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
				eventType: usr_model.UserReactivated,
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
				eventType: usr_model.UserLocked,
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
				eventType: usr_model.UserUnlocked,
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
				eventType: usr_model.InitializedUserCodeCreated,
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
			name: "user unlocked aggregate ok",
			args: args{
				ctx: auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					Profile: &model.Profile{UserName: "UserName"},
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: usr_model.MfaInitSkipped,
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
				eventType: usr_model.UserPasswordChanged,
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
		request    *model.RequestPasswordSet
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
				request:    &model.RequestPasswordSet{Expiry: time.Hour * 1},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: usr_model.UserPasswordSetRequested,
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
