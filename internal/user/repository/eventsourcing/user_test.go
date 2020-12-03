package eventsourcing

import (
	"context"
	"testing"
	"time"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/crypto"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
	usr_model "github.com/caos/zitadel/internal/user/model"
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

func TestHumanCreateAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		user       *model.User
		initCode   *model.InitUserCode
		phoneCode  *model.PhoneCode
		aggCreator *models.AggregateCreator
	}
	type res struct {
		eventLen      int
		eventTypes    []models.EventType
		aggregatesLen int
		checkData     []bool
		wantErr       bool
		errFunc       func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "user create aggregate ok",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				user: &model.User{
					ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					UserName:   "UserName",
					Human: &model.Human{
						Profile: &model.Profile{DisplayName: "DisplayName"},
						Email:   &model.Email{EmailAddress: "EmailAddress"},
					},
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:      1,
				eventTypes:    []models.EventType{model.HumanAdded},
				checkData:     []bool{true},
				aggregatesLen: 2,
			},
		},
		{
			name: "new user nil",
			args: args{
				ctx:        authz.NewMockContext("orgID", "userID"),
				user:       nil,
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
				ctx: authz.NewMockContext("orgID", "userID"),
				user: &model.User{
					ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					UserName:   "UserName",
					Human: &model.Human{
						Profile: &model.Profile{DisplayName: "DisplayName"},
						Email:   &model.Email{EmailAddress: "EmailAddress"},
					},
				},
				initCode:   &model.InitUserCode{},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:      2,
				eventTypes:    []models.EventType{model.HumanAdded, model.InitializedUserCodeAdded},
				checkData:     []bool{true, true},
				aggregatesLen: 2,
			},
		},
		{
			name: "create with phone code",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				user: &model.User{
					ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					UserName:   "UserName",
					Human: &model.Human{
						Profile: &model.Profile{DisplayName: "DisplayName"},
						Email:   &model.Email{EmailAddress: "EmailAddress"},
					},
				},
				phoneCode:  &model.PhoneCode{},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:      2,
				eventTypes:    []models.EventType{model.HumanAdded, model.UserPhoneCodeAdded},
				checkData:     []bool{true, true},
				aggregatesLen: 2,
			},
		},
		{
			name: "create with email verified",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				user: &model.User{
					ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					UserName:   "UserName",
					Human: &model.Human{
						Profile: &model.Profile{DisplayName: "DisplayName"},
						Email:   &model.Email{EmailAddress: "EmailAddress", IsEmailVerified: true},
					},
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:      2,
				eventTypes:    []models.EventType{model.HumanAdded, model.UserEmailVerified},
				checkData:     []bool{true, false},
				aggregatesLen: 2,
			},
		},
		{
			name: "create with phone verified",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				user: &model.User{
					ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					UserName:   "UserName",
					Human: &model.Human{
						Profile: &model.Profile{DisplayName: "DisplayName"},
						Email:   &model.Email{EmailAddress: "EmailAddress"},
						Phone:   &model.Phone{PhoneNumber: "PhoneNumber", IsPhoneVerified: true},
					},
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:      2,
				eventTypes:    []models.EventType{model.HumanAdded, model.UserPhoneVerified},
				checkData:     []bool{true, false},
				aggregatesLen: 2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aggregates, err := HumanCreateAggregate(tt.args.ctx, tt.args.aggCreator, tt.args.user, tt.args.initCode, tt.args.phoneCode, "", true)

			if !tt.res.wantErr && len(aggregates) != tt.res.aggregatesLen {
				t.Errorf("got wrong event len: expected: %v, actual: %v ", tt.res.aggregatesLen, len(aggregates))
			}

			if !tt.res.wantErr && len(aggregates[1].Events) != tt.res.eventLen {
				t.Errorf("got wrong event len: expected: %v, actual: %v ", tt.res.eventLen, len(aggregates[1].Events))
			}
			for i := 0; i < tt.res.eventLen; i++ {
				if !tt.res.wantErr && aggregates[1].Events[i].Type != tt.res.eventTypes[i] {
					t.Errorf("got wrong event type: expected: %v, actual: %v ", tt.res.eventTypes[i], aggregates[1].Events[i].Type.String())
				}
				if !tt.res.wantErr && tt.res.checkData[i] && aggregates[1].Events[i].Data == nil {
					t.Errorf("should have data in event")
				}
				if !tt.res.wantErr && !tt.res.checkData[i] && aggregates[1].Events[i].Data != nil {
					t.Errorf("should not have data in event")
				}
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestMachineCreateAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		user       *model.User
		aggCreator *models.AggregateCreator
	}
	type res struct {
		eventLen      int
		eventTypes    []models.EventType
		aggregatesLen int
		checkData     []bool
		wantErr       bool
		errFunc       func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "user create aggregate ok",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				user: &model.User{
					ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					UserName:   "UserName",
					Machine: &model.Machine{
						Description: "Description",
						Name:        "Name",
					},
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:      1,
				eventTypes:    []models.EventType{model.MachineAdded},
				checkData:     []bool{true},
				aggregatesLen: 2,
			},
		},
		{
			name: "new user nil",
			args: args{
				ctx:        authz.NewMockContext("orgID", "userID"),
				user:       nil,
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aggregates, err := MachineCreateAggregate(tt.args.ctx, tt.args.aggCreator, tt.args.user, "", true)

			if !tt.res.wantErr && len(aggregates) != tt.res.aggregatesLen {
				t.Errorf("got wrong event len: expected: %v, actual: %v ", tt.res.aggregatesLen, len(aggregates))
			}

			if !tt.res.wantErr && len(aggregates[0].Events) != tt.res.eventLen {
				t.Errorf("got wrong event len: expected: %v, actual: %v ", tt.res.eventLen, len(aggregates[1].Events))
			}
			for i := 0; i < tt.res.eventLen; i++ {
				if !tt.res.wantErr && aggregates[0].Events[i].Type != tt.res.eventTypes[i] {
					t.Errorf("got wrong event type: expected: %v, actual: %v ", tt.res.eventTypes[i], aggregates[0].Events[i].Type.String())
				}
				if !tt.res.wantErr && tt.res.checkData[i] && aggregates[0].Events[i].Data == nil {
					t.Errorf("should have data in event")
				}
				if !tt.res.wantErr && !tt.res.checkData[i] && aggregates[0].Events[i].Data != nil {
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
		user          *model.User
		externalIDP   *model.ExternalIDP
		initCode      *model.InitUserCode
		resourceOwner string
		aggCreator    *models.AggregateCreator
	}
	type res struct {
		eventLen   int
		eventTypes []models.EventType
		aggLen     int
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
				ctx: authz.NewMockContext("orgID", "userID"),
				user: &model.User{
					ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					UserName:   "UserName",
					Human: &model.Human{
						Profile: &model.Profile{DisplayName: "DisplayName"},
						Email:   &model.Email{EmailAddress: "EmailAddress"},
					},
				},
				initCode:      &model.InitUserCode{},
				resourceOwner: "newResourceowner",
				aggCreator:    models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   2,
				eventTypes: []models.EventType{model.HumanRegistered, model.InitializedHumanCodeAdded},
				aggLen:     2,
			},
		},
		{
			name: "user register with erxternalIDP aggregate ok",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				user: &model.User{
					ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					UserName:   "UserName",
					Human: &model.Human{
						Profile: &model.Profile{DisplayName: "DisplayName"},
						Email:   &model.Email{EmailAddress: "EmailAddress"},
					},
				},
				externalIDP:   &model.ExternalIDP{IDPConfigID: "IDPConfigID"},
				resourceOwner: "newResourceowner",
				aggCreator:    models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   2,
				eventTypes: []models.EventType{model.HumanRegistered, model.HumanExternalIDPAdded},
				aggLen:     3,
			},
		},
		{
			name: "new user nil",
			args: args{
				ctx:           authz.NewMockContext("orgID", "userID"),
				user:          nil,
				initCode:      &model.InitUserCode{},
				resourceOwner: "newResourceowner",
				aggCreator:    models.NewAggregateCreator("Test"),
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "create with init code",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				user: &model.User{
					ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					UserName:   "UserName",
					Human: &model.Human{
						Profile: &model.Profile{DisplayName: "DisplayName"},
						Email:   &model.Email{EmailAddress: "EmailAddress"},
					},
				},
				resourceOwner: "newResourceowner",
				initCode:      &model.InitUserCode{},
				aggCreator:    models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   2,
				eventTypes: []models.EventType{model.HumanRegistered, model.InitializedHumanCodeAdded},
				aggLen:     2,
			},
		},
		{
			name: "create no resourceowner",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				user: &model.User{
					ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					UserName:   "UserName",
					Human: &model.Human{
						Profile: &model.Profile{DisplayName: "DisplayName"},
						Email:   &model.Email{EmailAddress: "EmailAddress"},
					},
				},
				initCode:   &model.InitUserCode{},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aggregates, err := UserRegisterAggregate(tt.args.ctx, tt.args.aggCreator, tt.args.user, tt.args.externalIDP, tt.args.resourceOwner, tt.args.initCode, false)

			if tt.res.errFunc == nil && len(aggregates) != tt.res.aggLen {
				t.Errorf("got wrong aggregates len: expected: %v, actual: %v ", tt.res.aggLen, len(aggregates))
			}

			if tt.res.errFunc == nil && len(aggregates[tt.res.aggLen-1].Events) != tt.res.eventLen {
				t.Errorf("got wrong event len: expected: %v, actual: %v ", tt.res.eventLen, len(aggregates[tt.res.aggLen-1].Events))
			}
			for i := 0; i < tt.res.eventLen; i++ {
				if tt.res.errFunc == nil && aggregates[tt.res.aggLen-1].Events[i].Type != tt.res.eventTypes[i] {
					t.Errorf("got wrong event type: expected: %v, actual: %v ", tt.res.eventTypes[i], aggregates[tt.res.aggLen-1].Events[i].Type.String())
				}
				if tt.res.errFunc == nil && aggregates[tt.res.aggLen-1].Events[i].Data == nil {
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
		user       *model.User
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
				ctx: authz.NewMockContext("orgID", "userID"),
				user: &model.User{
					ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					UserName:   "UserName",
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
				ctx:        authz.NewMockContext("orgID", "userID"),
				user:       nil,
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := UserDeactivateAggregate(tt.args.aggCreator, tt.args.user)(tt.args.ctx)

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
		user       *model.User
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
				ctx: authz.NewMockContext("orgID", "userID"),
				user: &model.User{
					ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					UserName:   "UserName",
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
				ctx:        authz.NewMockContext("orgID", "userID"),
				user:       nil,
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := UserReactivateAggregate(tt.args.aggCreator, tt.args.user)(tt.args.ctx)

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
		user       *model.User
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
				ctx: authz.NewMockContext("orgID", "userID"),
				user: &model.User{
					ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					UserName:   "UserName",
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
				ctx:        authz.NewMockContext("orgID", "userID"),
				user:       nil,
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := UserLockAggregate(tt.args.aggCreator, tt.args.user)(tt.args.ctx)

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
		user       *model.User
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
				ctx: authz.NewMockContext("orgID", "userID"),
				user: &model.User{
					ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					UserName:   "UserName",
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
				ctx:        authz.NewMockContext("orgID", "userID"),
				user:       nil,
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := UserUnlockAggregate(tt.args.aggCreator, tt.args.user)(tt.args.ctx)

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
		ctx  context.Context
		user *model.User
		code *model.InitUserCode

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
				ctx: authz.NewMockContext("orgID", "userID"),
				user: &model.User{
					ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					UserName:   "UserName",
					Human: &model.Human{
						Profile: &model.Profile{DisplayName: "DisplayName"},
					},
				},
				code:       &model.InitUserCode{Expiry: time.Hour * 1},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.InitializedHumanCodeAdded,
			},
		},
		{
			name: "code nil",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				user: &model.User{
					ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					UserName:   "UserName",
					Human: &model.Human{
						Profile: &model.Profile{DisplayName: "DisplayName"},
					},
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
			agg, err := UserInitCodeAggregate(tt.args.aggCreator, tt.args.user, tt.args.code)(tt.args.ctx)

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
		user       *model.User
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
				ctx:        authz.NewMockContext("orgID", "userID"),
				user:       &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"}},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{model.InitializedHumanCodeSent},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := UserInitCodeSentAggregate(tt.args.aggCreator, tt.args.user)(tt.args.ctx)

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
		user       *model.User
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
				ctx: authz.NewMockContext("orgID", "userID"),
				user: &model.User{
					ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					UserName:   "UserName",
					Human: &model.Human{
						Profile: &model.Profile{DisplayName: "DisplayName"},
						Email:   &model.Email{EmailAddress: "EmailAddress"},
					},
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   2,
				eventTypes: []models.EventType{model.HumanEmailVerified, model.InitializedHumanCheckSucceeded},
			},
		},
		{
			name: "user init code only password",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				user: &model.User{
					ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					UserName:   "UserName",
					Human: &model.Human{
						Profile: &model.Profile{DisplayName: "DisplayName"},
						Email:   &model.Email{EmailAddress: "EmailAddress", IsEmailVerified: true},
					},
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{model.InitializedHumanCheckSucceeded},
			},
		},
		{
			name: "user init code email and pw",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				user: &model.User{
					ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					UserName:   "UserName",
					Human: &model.Human{
						Profile: &model.Profile{DisplayName: "DisplayName"},
						Email:   &model.Email{EmailAddress: "EmailAddress"},
					},
				},
				password:   &model.Password{Secret: &crypto.CryptoValue{}, ChangeRequired: false},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   3,
				eventTypes: []models.EventType{model.HumanEmailVerified, model.HumanPasswordChanged, model.InitializedHumanCheckSucceeded},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := InitCodeVerifiedAggregate(tt.args.aggCreator, tt.args.user, tt.args.password)(tt.args.ctx)

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

func TestInitCodeCheckFailedAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		user       *model.User
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
				ctx: authz.NewMockContext("orgID", "userID"),
				user: &model.User{
					ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					UserName:   "UserName",
					Human: &model.Human{
						Profile: &model.Profile{DisplayName: "DisplayName"},
					},
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.HumanMFAInitSkipped,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := SkipMFAAggregate(tt.args.aggCreator, tt.args.user)(tt.args.ctx)

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

func TestSkipMFAAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		user       *model.User
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
			name: "init code check failed",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				user: &model.User{
					ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					UserName:   "UserName",
					Human: &model.Human{
						Profile: &model.Profile{DisplayName: "DisplayName"},
					},
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.InitializedHumanCheckFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := InitCodeCheckFailedAggregate(tt.args.aggCreator, tt.args.user)(tt.args.ctx)

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
		user       *model.User
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
				ctx: authz.NewMockContext("orgID", "userID"),
				user: &model.User{
					ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					UserName:   "UserName",
					Human: &model.Human{
						Profile: &model.Profile{DisplayName: "DisplayName"},
					},
				},
				password:   &model.Password{ChangeRequired: true},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.HumanPasswordChanged,
			},
		},
		{
			name: "password nil",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					UserName: "UserName",
					Human: &model.Human{
						Profile: &model.Profile{DisplayName: "DisplayName"},
					},
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
			agg, err := PasswordChangeAggregate(tt.args.aggCreator, tt.args.user, tt.args.password)(tt.args.ctx)

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
		user       *model.User
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
				ctx: authz.NewMockContext("orgID", "userID"),
				user: &model.User{
					ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					UserName:   "UserName",
					Human: &model.Human{
						Profile: &model.Profile{DisplayName: "DisplayName"},
					},
				},
				request:    &model.PasswordCode{Expiry: time.Hour * 1},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.HumanPasswordCodeAdded,
			},
		},
		{
			name: "request nil",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				user: &model.User{
					ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					UserName:   "UserName",
					Human: &model.Human{
						Profile: &model.Profile{DisplayName: "DisplayName"},
					},
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
			agg, err := RequestSetPassword(tt.args.aggCreator, tt.args.user, tt.args.request)(tt.args.ctx)

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

func TestResendInitialPasswordAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		user       *model.User
		aggCreator *models.AggregateCreator
		initcode   *usr_model.InitUserCode
		email      string
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
			name: "resend initial password aggregate ok",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				user: &model.User{
					ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					UserName:   "UserName",
					Human: &model.Human{
						Profile: &model.Profile{DisplayName: "DisplayName"},
						Email:   &model.Email{},
					},
				},
				aggCreator: models.NewAggregateCreator("Test"),
				initcode:   &usr_model.InitUserCode{Expiry: time.Hour * 1},
			},
			res: res{
				eventLen: 1,
			},
		},
		{
			name: "resend initial password with same email ok",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				user: &model.User{
					ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					UserName:   "UserName",
					Human: &model.Human{
						Profile: &model.Profile{DisplayName: "DisplayName"},
						Email:   &model.Email{EmailAddress: "email"},
					},
				},
				aggCreator: models.NewAggregateCreator("Test"),
				initcode:   &usr_model.InitUserCode{Expiry: time.Hour * 1},
				email:      "email",
			},
			res: res{
				eventLen: 1,
			},
		},
		{
			name: "resend initial password with new email ok",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				user: &model.User{
					ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					UserName:   "UserName",
					Human: &model.Human{
						Profile: &model.Profile{DisplayName: "DisplayName"},
						Email:   &model.Email{EmailAddress: "old"},
					},
				},
				aggCreator: models.NewAggregateCreator("Test"),
				initcode:   &usr_model.InitUserCode{Expiry: time.Hour * 1},
				email:      "new",
			},
			res: res{
				eventLen: 2,
			},
		},
		{
			name: "request nil",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				user: &model.User{
					ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					UserName:   "UserName",
					Human: &model.Human{
						Profile: &model.Profile{DisplayName: "DisplayName"},
					},
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
			agg, err := ResendInitialPasswordAggregate(tt.args.aggCreator, tt.args.user, tt.args.initcode, tt.args.email)(tt.args.ctx)
			if (tt.res.errFunc == nil && err != nil) || (tt.res.errFunc != nil && !tt.res.errFunc(err)) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.errFunc == nil && len(agg.Events) != tt.res.eventLen {
				t.Errorf("got wrong event len: expected: %v, actual: %v ", tt.res.eventLen, len(agg.Events))
			}
		})
	}
}

func TestPasswordCodeSentAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		user       *model.User
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
				ctx:        authz.NewMockContext("orgID", "userID"),
				user:       &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"}},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{model.HumanPasswordCodeSent},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := PasswordCodeSentAggregate(tt.args.aggCreator, tt.args.user)(tt.args.ctx)

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
		user       *model.User
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
				ctx: authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					Human: &model.Human{
						Profile: &model.Profile{FirstName: "FirstName"},
					},
				},
				profile:    &model.Profile{FirstName: "FirstNameChanged"},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.HumanProfileChanged,
			},
		},
		{
			name: "profile nil",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				user: &model.User{
					ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					UserName:   "UserName",
					Human: &model.Human{
						Profile: &model.Profile{DisplayName: "DisplayName"},
					},
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
			agg, err := ProfileChangeAggregate(tt.args.aggCreator, tt.args.user, tt.args.profile)(tt.args.ctx)

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
		user       *model.User
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
				ctx: authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					Human: &model.Human{
						Email: &model.Email{EmailAddress: "EmailAddress"},
					},
				},
				email:      &model.Email{EmailAddress: "Changed", IsEmailVerified: true},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   2,
				eventTypes: []models.EventType{model.HumanEmailChanged, model.HumanEmailVerified},
			},
		},
		{
			name: "with code",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					Human: &model.Human{
						Email: &model.Email{EmailAddress: "EmailAddress"},
					},
				},
				email:      &model.Email{EmailAddress: "Changed"},
				code:       &model.EmailCode{Expiry: time.Hour * 1},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   2,
				eventTypes: []models.EventType{model.HumanEmailChanged, model.HumanEmailCodeAdded},
			},
		},
		{
			name: "email nil",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					Human: &model.Human{
						Email: &model.Email{EmailAddress: "Changed"},
					},
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
				ctx: authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					Human: &model.Human{
						Email: &model.Email{EmailAddress: "Changed", IsEmailVerified: true},
					},
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
				ctx: authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					Human: &model.Human{
						Email: &model.Email{EmailAddress: "Changed"},
					},
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
			aggregate, err := EmailChangeAggregate(tt.args.ctx, tt.args.aggCreator, tt.args.user, tt.args.email, tt.args.code)

			if tt.res.errFunc == nil && len(aggregate.Events) != tt.res.eventLen {
				t.Errorf("got wrong event len: expected: %v, actual: %v ", tt.res.eventLen, len(aggregate.Events))
			}
			for i := 0; i < tt.res.eventLen; i++ {
				if tt.res.errFunc == nil && aggregate.Events[i].Type != tt.res.eventTypes[i] {
					t.Errorf("got wrong event type: expected: %v, actual: %v ", tt.res.eventTypes[i], aggregate.Events[i].Type.String())
				}
				if tt.res.errFunc == nil && aggregate.Events[i].Data == nil {
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
		user       *model.User
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
				ctx:        authz.NewMockContext("orgID", "userID"),
				user:       &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"}},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{model.HumanEmailVerified},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := EmailVerifiedAggregate(tt.args.aggCreator, tt.args.user)(tt.args.ctx)

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
		user       *model.User
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
				ctx:        authz.NewMockContext("orgID", "userID"),
				user:       &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"}},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{model.HumanEmailVerificationFailed},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := EmailVerificationFailedAggregate(tt.args.aggCreator, tt.args.user)(tt.args.ctx)

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
		user       *model.User
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
				ctx: authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					Human: &model.Human{
						Email: &model.Email{EmailAddress: "EmailAddress"},
					},
				},
				code:       &model.EmailCode{Expiry: time.Hour * 1},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{model.HumanEmailCodeAdded},
			},
		},
		{
			name: "code nil",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					Human: &model.Human{
						Email: &model.Email{EmailAddress: "Changed"},
					},
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
			agg, err := EmailVerificationCodeAggregate(tt.args.aggCreator, tt.args.user, tt.args.code)(tt.args.ctx)

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
		user       *model.User
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
				ctx:        authz.NewMockContext("orgID", "userID"),
				user:       &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"}},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{model.HumanEmailCodeSent},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := EmailCodeSentAggregate(tt.args.aggCreator, tt.args.user)(tt.args.ctx)

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
		user       *model.User
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
				ctx: authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					Human: &model.Human{
						Phone: &model.Phone{PhoneNumber: "+41791234567"},
					},
				},
				phone:      &model.Phone{PhoneNumber: "+41799876543", IsPhoneVerified: true},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   2,
				eventTypes: []models.EventType{model.HumanPhoneChanged, model.HumanPhoneVerified},
			},
		},
		{
			name: "with code",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					Human: &model.Human{
						Phone: &model.Phone{PhoneNumber: "PhoneNumber"},
					},
				},
				phone:      &model.Phone{PhoneNumber: "Changed"},
				code:       &model.PhoneCode{Expiry: time.Hour * 1},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   2,
				eventTypes: []models.EventType{model.HumanPhoneChanged, model.HumanPhoneCodeAdded},
			},
		},
		{
			name: "phone nil",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					Human: &model.Human{
						Phone: &model.Phone{PhoneNumber: "Changed"},
					},
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
				ctx: authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					Human: &model.Human{
						Phone: &model.Phone{PhoneNumber: "Changed"},
					},
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
				ctx: authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					Human: &model.Human{
						Phone: &model.Phone{PhoneNumber: "Changed"},
					},
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
			agg, err := PhoneChangeAggregate(tt.args.aggCreator, tt.args.user, tt.args.phone, tt.args.code)(tt.args.ctx)

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
		user       *model.User
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
				ctx: authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					Human: &model.Human{
						Phone: &model.Phone{PhoneNumber: "PhoneNumber"},
					},
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{model.HumanPhoneVerified},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := PhoneVerifiedAggregate(tt.args.aggCreator, tt.args.user)(tt.args.ctx)

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

func TestRemovePhoneAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		user       *model.User
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
			name: "user phone removed aggregate ok",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					Human: &model.Human{
						Phone: &model.Phone{PhoneNumber: "PhoneNumber"},
					},
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{model.HumanPhoneRemoved},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := PhoneRemovedAggregate(tt.args.aggCreator, tt.args.user)(tt.args.ctx)

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
		user       *model.User
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
				ctx: authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					Human: &model.Human{
						Phone: &model.Phone{PhoneNumber: "PhoneNumber"},
					},
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{model.HumanPhoneVerificationFailed},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := PhoneVerificationFailedAggregate(tt.args.aggCreator, tt.args.user)(tt.args.ctx)

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
		user       *model.User
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
				ctx: authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					Human: &model.Human{
						Email: &model.Email{EmailAddress: "EmailAddress"},
					},
				},
				code:       &model.PhoneCode{Expiry: time.Hour * 1},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{model.HumanPhoneCodeAdded},
			},
		},
		{
			name: "code nil",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					Human: &model.Human{
						Email: &model.Email{EmailAddress: "Changed"},
					},
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
			agg, err := PhoneVerificationCodeAggregate(tt.args.aggCreator, tt.args.user, tt.args.code)(tt.args.ctx)

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
		user       *model.User
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
				ctx:        authz.NewMockContext("orgID", "userID"),
				user:       &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"}},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{model.HumanPhoneCodeSent},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := PhoneCodeSentAggregate(tt.args.aggCreator, tt.args.user)(tt.args.ctx)

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
		user       *model.User
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
				ctx: authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					Human: &model.Human{
						Address: &model.Address{Locality: "Locality"},
					},
				},
				address:    &model.Address{Locality: "Changed"},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{model.HumanAddressChanged},
			},
		},
		{
			name: "address nil",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"},
					Human: &model.Human{
						Address: &model.Address{Locality: "Changed"},
					},
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
			agg, err := AddressChangeAggregate(tt.args.aggCreator, tt.args.user, tt.args.address)(tt.args.ctx)

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

func TestOTPAddAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		user       *model.User
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
				ctx:        authz.NewMockContext("orgID", "userID"),
				user:       &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"}},
				otp:        &model.OTP{},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{model.HumanMFAOTPAdded},
			},
		},
		{
			name: "otp nil",
			args: args{
				ctx:        authz.NewMockContext("orgID", "userID"),
				user:       &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"}},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := MFAOTPAddAggregate(tt.args.aggCreator, tt.args.user, tt.args.otp)(tt.args.ctx)

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

func TestOTPVerifyAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		user       *model.User
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
				ctx:        authz.NewMockContext("orgID", "userID"),
				user:       &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"}},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{model.HumanMFAOTPVerified},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := MFAOTPVerifyAggregate(tt.args.aggCreator, tt.args.user)(tt.args.ctx)

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

func TestOTPRemoveAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		user       *model.User
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
				ctx:        authz.NewMockContext("orgID", "userID"),
				user:       &model.User{ObjectRoot: models.ObjectRoot{AggregateID: "ID"}},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{model.HumanMFAOTPRemoved},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := MFAOTPRemoveAggregate(tt.args.aggCreator, tt.args.user)(tt.args.ctx)

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

func TestExternalIDPAddedAggregates(t *testing.T) {
	type res struct {
		aggregateCount int
		isErr          func(error) bool
	}
	type args struct {
		ctx         context.Context
		aggCreator  *models.AggregateCreator
		user        *model.User
		externalIDP *model.ExternalIDP
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "no user error",
			args: args{
				ctx:        authz.NewMockContext("org", "user"),
				aggCreator: models.NewAggregateCreator("test"),
				user:       nil,
			},
			res: res{
				aggregateCount: 0,
				isErr:          caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "user add external idp successful",
			args: args{
				ctx:        authz.NewMockContext("org", "user"),
				aggCreator: models.NewAggregateCreator("test"),
				user: &model.User{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "AggregateID",
						Sequence:    5,
					},
				},
				externalIDP: &model.ExternalIDP{
					IDPConfigID: "IDPConfigID",
					UserID:      "UserID",
					DisplayName: "DisplayName",
				},
			},
			res: res{
				aggregateCount: 2,
				isErr:          nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExternalIDPAddedAggregate(tt.args.ctx, tt.args.aggCreator, tt.args.user, tt.args.externalIDP)
			if tt.res.isErr == nil && err != nil {
				t.Errorf("no error expected got %T: %v", err, err)
			}
			if tt.res.isErr != nil && !tt.res.isErr(err) {
				t.Errorf("wrong error got %T: %v", err, err)
			}
			if tt.res.isErr == nil && len(got) != tt.res.aggregateCount {
				t.Errorf("ExternalIDPAddedAggregate() aggregate count = %d, wanted count %d", len(got), tt.res.aggregateCount)
			}
		})
	}
}

func TestExternalIDPRemovedAggregates(t *testing.T) {
	type res struct {
		aggregateCount int
		isErr          func(error) bool
	}
	type args struct {
		ctx         context.Context
		aggCreator  *models.AggregateCreator
		user        *model.User
		externalIDP *model.ExternalIDP
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "no user error",
			args: args{
				ctx:        authz.NewMockContext("org", "user"),
				aggCreator: models.NewAggregateCreator("test"),
				user:       nil,
			},
			res: res{
				aggregateCount: 0,
				isErr:          caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "user removed external idp successful",
			args: args{
				ctx:        authz.NewMockContext("org", "user"),
				aggCreator: models.NewAggregateCreator("test"),
				user: &model.User{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "AggregateID",
						Sequence:      5,
						ResourceOwner: "ResourceOwner",
					},
				},
				externalIDP: &model.ExternalIDP{
					IDPConfigID: "IDPConfigID",
					UserID:      "UserID",
				},
			},
			res: res{
				aggregateCount: 2,
				isErr:          nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExternalIDPRemovedAggregate(tt.args.ctx, tt.args.aggCreator, tt.args.user, tt.args.externalIDP, false)
			if tt.res.isErr == nil && err != nil {
				t.Errorf("no error expected got %T: %v", err, err)
			}
			if tt.res.isErr != nil && !tt.res.isErr(err) {
				t.Errorf("wrong error got %T: %v", err, err)
			}
			if tt.res.isErr == nil && len(got) != tt.res.aggregateCount {
				t.Errorf("ExternalIDPRemovedAggregate() aggregate count = %d, wanted count %d", len(got), tt.res.aggregateCount)
			}
		})
	}
}
