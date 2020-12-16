package eventsourcing

import (
	"context"
	"testing"

	"github.com/caos/zitadel/internal/api/authz"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
)

func TestSetUpStartedAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		iam        *model.IAM
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
			name: "setupstarted aggregate ok",
			args: args{
				ctx:        authz.NewMockContext("orgID", "userID"),
				iam:        &model.IAM{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.IAMSetupStarted,
			},
		},
		{
			name: "iam nil",
			args: args{
				ctx:        authz.NewMockContext("orgID", "userID"),
				iam:        nil,
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.IAMSetupStarted,
				errFunc:   caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := IAMSetupStartedAggregate(tt.args.aggCreator, tt.args.iam)(tt.args.ctx)

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

func TestSetUpDoneAggregate(t *testing.T) {
	type args struct {
		ctx         context.Context
		existingIAM *model.IAM
		aggCreator  *models.AggregateCreator
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
			name: "setup done aggregate ok",
			args: args{
				ctx:         authz.NewMockContext("orgID", "userID"),
				existingIAM: &model.IAM{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}},
				aggCreator:  models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.IAMSetupDone,
			},
		},
		{
			name: "existing iam nil",
			args: args{
				ctx:         authz.NewMockContext("orgID", "userID"),
				existingIAM: nil,
				aggCreator:  models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.IAMSetupDone,
				errFunc:   caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := IAMSetupDoneAggregate(tt.args.aggCreator, tt.args.existingIAM)(tt.args.ctx)

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

func TestGlobalOrgAggregate(t *testing.T) {
	type args struct {
		ctx         context.Context
		existingIAM *model.IAM
		orgID       string
		aggCreator  *models.AggregateCreator
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
			name: "global org set aggregate ok",
			args: args{
				ctx:         authz.NewMockContext("orgID", "userID"),
				existingIAM: &model.IAM{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}},
				orgID:       "orgID",
				aggCreator:  models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.GlobalOrgSet,
			},
		},
		{
			name: "existing iam nil",
			args: args{
				ctx:         authz.NewMockContext("orgID", "userID"),
				existingIAM: nil,
				orgID:       "orgID",
				aggCreator:  models.NewAggregateCreator("Test"),
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "global org empty",
			args: args{
				ctx:         authz.NewMockContext("orgID", "userID"),
				existingIAM: &model.IAM{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}},
				aggCreator:  models.NewAggregateCreator("Test"),
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := IAMSetGlobalOrgAggregate(tt.args.aggCreator, tt.args.existingIAM, tt.args.orgID)(tt.args.ctx)

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

func TestIamProjectAggregate(t *testing.T) {
	type args struct {
		ctx         context.Context
		existingIAM *model.IAM
		projectID   string
		aggCreator  *models.AggregateCreator
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
			name: "iam project id set aggregate ok",
			args: args{
				ctx:         authz.NewMockContext("orgID", "userID"),
				existingIAM: &model.IAM{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}},
				projectID:   "projectID",
				aggCreator:  models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.IAMProjectSet,
			},
		},
		{
			name: "existing iam nil",
			args: args{
				ctx:         authz.NewMockContext("orgID", "userID"),
				existingIAM: nil,
				projectID:   "projectID",
				aggCreator:  models.NewAggregateCreator("Test"),
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "project id empty",
			args: args{
				ctx:         authz.NewMockContext("orgID", "userID"),
				existingIAM: &model.IAM{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}},
				aggCreator:  models.NewAggregateCreator("Test"),
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := IAMSetIamProjectAggregate(tt.args.aggCreator, tt.args.existingIAM, tt.args.projectID)(tt.args.ctx)

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

func TestIamMemberAddedAggregate(t *testing.T) {
	type args struct {
		ctx         context.Context
		existingIAM *model.IAM
		newMember   *model.IAMMember
		aggCreator  *models.AggregateCreator
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
			name: "iammember added ok",
			args: args{
				ctx:         authz.NewMockContext("orgID", "userID"),
				existingIAM: &model.IAM{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}},
				newMember:   &model.IAMMember{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, UserID: "UserID", Roles: []string{"Roles"}},
				aggCreator:  models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.IAMMemberAdded,
			},
		},
		{
			name: "existing iam nil",
			args: args{
				ctx:         authz.NewMockContext("orgID", "userID"),
				existingIAM: nil,
				aggCreator:  models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.IAMMemberAdded,
				errFunc:   caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "member nil",
			args: args{
				ctx:         authz.NewMockContext("orgID", "userID"),
				existingIAM: &model.IAM{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}},
				newMember:   nil,
				aggCreator:  models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.IAMMemberAdded,
				errFunc:   caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := IAMMemberAddedAggregate(tt.args.aggCreator, tt.args.existingIAM, tt.args.newMember)(tt.args.ctx)

			if tt.res.errFunc == nil && len(agg.Events) != tt.res.eventLen {
				t.Errorf("got wrong event len: expected: %v, actual: %v ", tt.res.eventLen, len(agg.Events))
			}
			if tt.res.errFunc == nil && agg.Events[0].Type != tt.res.eventType {
				t.Errorf("got wrong event type: expected: %v, actual: %v ", tt.res.eventType, agg.Events[0].Type.String())
			}
			if tt.res.errFunc == nil && agg.Events[0].Data == nil {
				t.Errorf("should have data in event")
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestIamMemberChangedAggregate(t *testing.T) {
	type args struct {
		ctx         context.Context
		existingIAM *model.IAM
		newMember   *model.IAMMember
		aggCreator  *models.AggregateCreator
	}
	type res struct {
		eventLen  int
		eventType models.EventType
		wantErr   bool
		errFunc   func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "iammember changed ok",
			args: args{
				ctx:         authz.NewMockContext("orgID", "userID"),
				existingIAM: &model.IAM{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}},
				newMember:   &model.IAMMember{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, UserID: "UserID", Roles: []string{"Roles"}},
				aggCreator:  models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.IAMMemberChanged,
			},
		},
		{
			name: "existing iam nil",
			args: args{
				ctx:         authz.NewMockContext("orgID", "userID"),
				existingIAM: nil,
				aggCreator:  models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.IAMMemberChanged,
				wantErr:   true,
				errFunc:   caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "member nil",
			args: args{
				ctx:         authz.NewMockContext("orgID", "userID"),
				existingIAM: &model.IAM{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}},
				newMember:   nil,
				aggCreator:  models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.IAMMemberChanged,
				wantErr:   true,
				errFunc:   caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := IAMMemberChangedAggregate(tt.args.aggCreator, tt.args.existingIAM, tt.args.newMember)(tt.args.ctx)

			if !tt.res.wantErr && len(agg.Events) != tt.res.eventLen {
				t.Errorf("got wrong event len: expected: %v, actual: %v ", tt.res.eventLen, len(agg.Events))
			}
			if !tt.res.wantErr && agg.Events[0].Type != tt.res.eventType {
				t.Errorf("got wrong event type: expected: %v, actual: %v ", tt.res.eventType, agg.Events[0].Type.String())
			}
			if !tt.res.wantErr && agg.Events[0].Data == nil {
				t.Errorf("should have data in event")
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestIamMemberRemovedAggregate(t *testing.T) {
	type args struct {
		ctx         context.Context
		existingIAM *model.IAM
		newMember   *model.IAMMember
		aggCreator  *models.AggregateCreator
	}
	type res struct {
		eventLen  int
		eventType models.EventType
		wantErr   bool
		errFunc   func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "iammember removed ok",
			args: args{
				ctx:         authz.NewMockContext("orgID", "userID"),
				existingIAM: &model.IAM{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}},
				newMember:   &model.IAMMember{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, UserID: "UserID", Roles: []string{"Roles"}},
				aggCreator:  models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.IAMMemberRemoved,
			},
		},
		{
			name: "existing iam nil",
			args: args{
				ctx:         authz.NewMockContext("orgID", "userID"),
				existingIAM: nil,
				aggCreator:  models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.IAMMemberRemoved,
				wantErr:   true,
				errFunc:   caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "member nil",
			args: args{
				ctx:         authz.NewMockContext("orgID", "userID"),
				existingIAM: &model.IAM{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}},
				newMember:   nil,
				aggCreator:  models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.IAMMemberRemoved,
				wantErr:   true,
				errFunc:   caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := IAMMemberRemovedAggregate(tt.args.aggCreator, tt.args.existingIAM, tt.args.newMember)(tt.args.ctx)

			if !tt.res.wantErr && len(agg.Events) != tt.res.eventLen {
				t.Errorf("got wrong event len: expected: %v, actual: %v ", tt.res.eventLen, len(agg.Events))
			}
			if !tt.res.wantErr && agg.Events[0].Type != tt.res.eventType {
				t.Errorf("got wrong event type: expected: %v, actual: %v ", tt.res.eventType, agg.Events[0].Type.String())
			}
			if !tt.res.wantErr && agg.Events[0].Data == nil {
				t.Errorf("should have data in event")
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestIdpConfigAddedAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		existing   *model.IAM
		newConfig  *model.IDPConfig
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
			name: "add oidc idp configuration",
			args: args{
				ctx:      authz.NewMockContext("orgID", "userID"),
				existing: &model.IAM{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, IAMProjectID: "IAMProjectID"},
				newConfig: &model.IDPConfig{
					ObjectRoot:    models.ObjectRoot{AggregateID: "AggregateID"},
					IDPConfigID:   "IDPConfigID",
					Name:          "Name",
					OIDCIDPConfig: &model.OIDCIDPConfig{IDPConfigID: "IDPConfigID", ClientID: "ClientID"},
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   2,
				eventTypes: []models.EventType{model.IDPConfigAdded, model.OIDCIDPConfigAdded},
			},
		},
		{
			name: "existing iam nil",
			args: args{
				ctx:        authz.NewMockContext("orgID", "userID"),
				existing:   nil,
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "idp config nil",
			args: args{
				ctx:        authz.NewMockContext("orgID", "userID"),
				existing:   &model.IAM{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, IAMProjectID: "IAMProjectID"},
				newConfig:  nil,
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
			agg, err := IDPConfigAddedAggregate(tt.args.aggCreator, tt.args.existing, tt.args.newConfig)(tt.args.ctx)

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

func TestIdpConfigurationChangedAggregate(t *testing.T) {
	type args struct {
		ctx         context.Context
		existingIAM *model.IAM
		newConfig   *model.IDPConfig
		aggCreator  *models.AggregateCreator
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
			name: "change idp configuration",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				existingIAM: &model.IAM{
					ObjectRoot:   models.ObjectRoot{AggregateID: "AggregateID"},
					IAMProjectID: "IAMProjectID",
					IDPs: []*model.IDPConfig{
						{IDPConfigID: "IDPConfigID", Name: "IDPName"},
					}},
				newConfig: &model.IDPConfig{
					ObjectRoot:  models.ObjectRoot{AggregateID: "AggregateID"},
					IDPConfigID: "IDPConfigID",
					Name:        "NameChanged",
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{model.IDPConfigChanged},
			},
		},
		{
			name: "existing iam nil",
			args: args{
				ctx:         authz.NewMockContext("orgID", "userID"),
				existingIAM: nil,
				aggCreator:  models.NewAggregateCreator("Test"),
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "idp config nil",
			args: args{
				ctx:         authz.NewMockContext("orgID", "userID"),
				existingIAM: &model.IAM{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, IAMProjectID: "IAMProjectID"},
				newConfig:   nil,
				aggCreator:  models.NewAggregateCreator("Test"),
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := IDPConfigChangedAggregate(tt.args.aggCreator, tt.args.existingIAM, tt.args.newConfig)(tt.args.ctx)

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

func TestIdpConfigurationRemovedAggregate(t *testing.T) {
	type args struct {
		ctx         context.Context
		existingIAM *model.IAM
		newConfig   *model.IDPConfig
		provider    *model.IDPProvider
		aggCreator  *models.AggregateCreator
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
			name: "remove idp config",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				existingIAM: &model.IAM{
					ObjectRoot:   models.ObjectRoot{AggregateID: "AggregateID"},
					IAMProjectID: "IAMProjectID",
					IDPs: []*model.IDPConfig{
						{IDPConfigID: "IDPConfigID", Name: "Name"},
					}},
				newConfig: &model.IDPConfig{
					ObjectRoot:  models.ObjectRoot{AggregateID: "AggregateID"},
					IDPConfigID: "IDPConfigID",
					Name:        "Name",
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{model.IDPConfigRemoved},
			},
		},
		{
			name: "remove idp config with provider",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				existingIAM: &model.IAM{
					ObjectRoot:   models.ObjectRoot{AggregateID: "AggregateID"},
					IAMProjectID: "IAMProjectID",
					IDPs: []*model.IDPConfig{
						{IDPConfigID: "IDPConfigID", Name: "Name"},
					}},
				newConfig: &model.IDPConfig{
					ObjectRoot:  models.ObjectRoot{AggregateID: "AggregateID"},
					IDPConfigID: "IDPConfigID",
					Name:        "Name",
				},
				provider: &model.IDPProvider{
					IDPConfigID: "IDPConfigID",
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   2,
				eventTypes: []models.EventType{model.IDPConfigRemoved, model.LoginPolicyIDPProviderCascadeRemoved},
			},
		},
		{
			name: "existing iam nil",
			args: args{
				ctx:         authz.NewMockContext("orgID", "userID"),
				existingIAM: nil,
				aggCreator:  models.NewAggregateCreator("Test"),
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "idp config nil",
			args: args{
				ctx:         authz.NewMockContext("orgID", "userID"),
				existingIAM: &model.IAM{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, IAMProjectID: "IAMProjectID"},
				newConfig:   nil,
				aggCreator:  models.NewAggregateCreator("Test"),
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := IDPConfigRemovedAggregate(tt.args.ctx, tt.args.aggCreator, tt.args.existingIAM, tt.args.newConfig, tt.args.provider)

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

func TestIdpConfigurationDeactivatedAggregate(t *testing.T) {
	type args struct {
		ctx         context.Context
		existingIAM *model.IAM
		newConfig   *model.IDPConfig
		aggCreator  *models.AggregateCreator
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
			name: "deactivate idp config",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				existingIAM: &model.IAM{
					ObjectRoot:   models.ObjectRoot{AggregateID: "AggregateID"},
					IAMProjectID: "IAMProjectID",
					IDPs: []*model.IDPConfig{
						{IDPConfigID: "IDPConfigID", Name: "Name"},
					}},
				newConfig: &model.IDPConfig{
					ObjectRoot:  models.ObjectRoot{AggregateID: "AggregateID"},
					IDPConfigID: "IDPConfigID",
					Name:        "Name",
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{model.IDPConfigDeactivated},
			},
		},
		{
			name: "existing iam nil",
			args: args{
				ctx:         authz.NewMockContext("orgID", "userID"),
				existingIAM: nil,
				aggCreator:  models.NewAggregateCreator("Test"),
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "idp config nil",
			args: args{
				ctx:         authz.NewMockContext("orgID", "userID"),
				existingIAM: &model.IAM{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, IAMProjectID: "IAMProjectID"},
				newConfig:   nil,
				aggCreator:  models.NewAggregateCreator("Test"),
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := IDPConfigDeactivatedAggregate(tt.args.aggCreator, tt.args.existingIAM, tt.args.newConfig)(tt.args.ctx)

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

func TestIdpConfigurationReactivatedAggregate(t *testing.T) {
	type args struct {
		ctx         context.Context
		existingIAM *model.IAM
		newConfig   *model.IDPConfig
		aggCreator  *models.AggregateCreator
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
			name: "deactivate app",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				existingIAM: &model.IAM{
					ObjectRoot:   models.ObjectRoot{AggregateID: "AggregateID"},
					IAMProjectID: "IAMProjectID",
					IDPs: []*model.IDPConfig{
						{IDPConfigID: "IDPConfigID", Name: "Name"},
					}},
				newConfig: &model.IDPConfig{
					ObjectRoot:  models.ObjectRoot{AggregateID: "AggregateID"},
					IDPConfigID: "IDPConfigID",
					Name:        "Name",
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{model.IDPConfigReactivated},
			},
		},
		{
			name: "existing iam nil",
			args: args{
				ctx:         authz.NewMockContext("orgID", "userID"),
				existingIAM: nil,
				aggCreator:  models.NewAggregateCreator("Test"),
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "idp config nil",
			args: args{
				ctx:         authz.NewMockContext("orgID", "userID"),
				existingIAM: &model.IAM{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, IAMProjectID: "IAMProjectID"},
				newConfig:   nil,
				aggCreator:  models.NewAggregateCreator("Test"),
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := IDPConfigReactivatedAggregate(tt.args.aggCreator, tt.args.existingIAM, tt.args.newConfig)(tt.args.ctx)

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

func TestOIDCConfigChangedAggregate(t *testing.T) {
	type args struct {
		ctx         context.Context
		existingIAM *model.IAM
		newConfig   *model.OIDCIDPConfig
		aggCreator  *models.AggregateCreator
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
			name: "change oidc config",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				existingIAM: &model.IAM{
					ObjectRoot:   models.ObjectRoot{AggregateID: "AggregateID"},
					IAMProjectID: "IAMProjectID",
					IDPs: []*model.IDPConfig{
						{IDPConfigID: "IDPConfigID", Name: "Name", OIDCIDPConfig: &model.OIDCIDPConfig{IDPConfigID: "IDPConfigID", ClientID: "ClientID"}},
					}},
				newConfig: &model.OIDCIDPConfig{
					ObjectRoot:  models.ObjectRoot{AggregateID: "AggregateID"},
					IDPConfigID: "IDPConfigID",
					ClientID:    "ClientIDChanged",
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{model.OIDCIDPConfigChanged},
			},
		},
		{
			name: "no changes",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				existingIAM: &model.IAM{
					ObjectRoot:   models.ObjectRoot{AggregateID: "AggregateID"},
					IAMProjectID: "IAMProjectID",
					IDPs: []*model.IDPConfig{
						{IDPConfigID: "IDPConfigID", Name: "Name", OIDCIDPConfig: &model.OIDCIDPConfig{IDPConfigID: "IDPConfigID", ClientID: "ClientID"}},
					}},
				newConfig: &model.OIDCIDPConfig{
					ObjectRoot:  models.ObjectRoot{AggregateID: "AggregateID"},
					IDPConfigID: "IDPConfigID",
					ClientID:    "ClientID",
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing iam nil",
			args: args{
				ctx:         authz.NewMockContext("orgID", "userID"),
				existingIAM: nil,
				aggCreator:  models.NewAggregateCreator("Test"),
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "oidc config nil",
			args: args{
				ctx:         authz.NewMockContext("orgID", "userID"),
				existingIAM: &model.IAM{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, IAMProjectID: "IAMProjectID"},
				newConfig:   nil,
				aggCreator:  models.NewAggregateCreator("Test"),
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := OIDCIDPConfigChangedAggregate(tt.args.aggCreator, tt.args.existingIAM, tt.args.newConfig)(tt.args.ctx)

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

func TestLoginPolicyAddedAggregate(t *testing.T) {
	type args struct {
		ctx         context.Context
		existingIAM *model.IAM
		newPolicy   *model.LoginPolicy
		aggCreator  *models.AggregateCreator
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
			name: "add login polciy",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				existingIAM: &model.IAM{
					ObjectRoot:   models.ObjectRoot{AggregateID: "AggregateID"},
					IAMProjectID: "IAMProjectID",
					IDPs: []*model.IDPConfig{
						{IDPConfigID: "IDPConfigID", Name: "Name", OIDCIDPConfig: &model.OIDCIDPConfig{IDPConfigID: "IDPConfigID", ClientID: "ClientID"}},
					}},
				newPolicy: &model.LoginPolicy{
					ObjectRoot:            models.ObjectRoot{AggregateID: "AggregateID"},
					AllowUsernamePassword: true,
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{model.LoginPolicyAdded},
			},
		},
		{
			name: "existing iam nil",
			args: args{
				ctx:         authz.NewMockContext("orgID", "userID"),
				existingIAM: nil,
				aggCreator:  models.NewAggregateCreator("Test"),
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "login policy config nil",
			args: args{
				ctx:         authz.NewMockContext("orgID", "userID"),
				existingIAM: &model.IAM{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, IAMProjectID: "IAMProjectID"},
				newPolicy:   nil,
				aggCreator:  models.NewAggregateCreator("Test"),
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := LoginPolicyAddedAggregate(tt.args.ctx, tt.args.aggCreator, tt.args.existingIAM, tt.args.newPolicy)

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

func TestLoginPolicyChangedAggregate(t *testing.T) {
	type args struct {
		ctx         context.Context
		existingIAM *model.IAM
		newPolicy   *model.LoginPolicy
		aggCreator  *models.AggregateCreator
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
			name: "change login policy",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				existingIAM: &model.IAM{
					ObjectRoot:   models.ObjectRoot{AggregateID: "AggregateID"},
					IAMProjectID: "IAMProjectID",
					DefaultLoginPolicy: &model.LoginPolicy{
						AllowUsernamePassword: true,
					}},
				newPolicy: &model.LoginPolicy{
					ObjectRoot:            models.ObjectRoot{AggregateID: "AggregateID"},
					AllowUsernamePassword: true,
					AllowRegister:         true,
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{model.LoginPolicyChanged},
			},
		},
		{
			name: "no changes",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				existingIAM: &model.IAM{
					ObjectRoot:   models.ObjectRoot{AggregateID: "AggregateID"},
					IAMProjectID: "IAMProjectID",
					DefaultLoginPolicy: &model.LoginPolicy{
						AllowUsernamePassword: true,
					}},
				newPolicy: &model.LoginPolicy{
					ObjectRoot:            models.ObjectRoot{AggregateID: "AggregateID"},
					AllowUsernamePassword: true,
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing iam nil",
			args: args{
				ctx:         authz.NewMockContext("orgID", "userID"),
				existingIAM: nil,
				aggCreator:  models.NewAggregateCreator("Test"),
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "login policy config nil",
			args: args{
				ctx:         authz.NewMockContext("orgID", "userID"),
				existingIAM: &model.IAM{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, IAMProjectID: "IAMProjectID"},
				newPolicy:   nil,
				aggCreator:  models.NewAggregateCreator("Test"),
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := LoginPolicyChangedAggregate(tt.args.aggCreator, tt.args.existingIAM, tt.args.newPolicy)(tt.args.ctx)

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

func TestLoginPolicyIdpProviderAddedAggregate(t *testing.T) {
	type args struct {
		ctx         context.Context
		existingIAM *model.IAM
		newProvider *model.IDPProvider
		aggCreator  *models.AggregateCreator
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
			name: "add idp provider to login policy",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				existingIAM: &model.IAM{
					ObjectRoot:   models.ObjectRoot{AggregateID: "AggregateID"},
					IAMProjectID: "IAMProjectID",
					DefaultLoginPolicy: &model.LoginPolicy{
						AllowUsernamePassword: true,
					}},
				newProvider: &model.IDPProvider{
					ObjectRoot:  models.ObjectRoot{AggregateID: "AggregateID"},
					Type:        int32(iam_model.IDPProviderTypeSystem),
					IDPConfigID: "IDPConfigID",
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{model.LoginPolicyIDPProviderAdded},
			},
		},
		{
			name: "existing iam nil",
			args: args{
				ctx:         authz.NewMockContext("orgID", "userID"),
				existingIAM: nil,
				aggCreator:  models.NewAggregateCreator("Test"),
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "idp config config nil",
			args: args{
				ctx:         authz.NewMockContext("orgID", "userID"),
				existingIAM: &model.IAM{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, IAMProjectID: "IAMProjectID"},
				newProvider: nil,
				aggCreator:  models.NewAggregateCreator("Test"),
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := LoginPolicyIDPProviderAddedAggregate(tt.args.aggCreator, tt.args.existingIAM, tt.args.newProvider)(tt.args.ctx)

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

func TestLoginPolicyIdpProviderRemovedAggregate(t *testing.T) {
	type args struct {
		ctx           context.Context
		existingIAM   *model.IAM
		newProviderID *model.IDPProviderID
		aggCreator    *models.AggregateCreator
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
			name: "remove idp provider to login policy",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				existingIAM: &model.IAM{
					ObjectRoot:   models.ObjectRoot{AggregateID: "AggregateID"},
					IAMProjectID: "IAMProjectID",
					DefaultLoginPolicy: &model.LoginPolicy{
						AllowUsernamePassword: true,
						IDPProviders: []*model.IDPProvider{
							{IDPConfigID: "IDPConfigID", Type: int32(iam_model.IDPProviderTypeSystem)},
						},
					}},
				newProviderID: &model.IDPProviderID{
					IDPConfigID: "IDPConfigID",
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{model.LoginPolicyIDPProviderRemoved},
			},
		},
		{
			name: "existing iam nil",
			args: args{
				ctx:         authz.NewMockContext("orgID", "userID"),
				existingIAM: nil,
				aggCreator:  models.NewAggregateCreator("Test"),
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "idp config config nil",
			args: args{
				ctx:           authz.NewMockContext("orgID", "userID"),
				existingIAM:   &model.IAM{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, IAMProjectID: "IAMProjectID"},
				newProviderID: nil,
				aggCreator:    models.NewAggregateCreator("Test"),
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := LoginPolicyIDPProviderRemovedAggregate(tt.args.ctx, tt.args.aggCreator, tt.args.existingIAM, tt.args.newProviderID)

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

func TestLoginPolicySecondFactorAddedAggregate(t *testing.T) {
	type args struct {
		ctx         context.Context
		existingIAM *model.IAM
		newMFA      *model.MFA
		aggCreator  *models.AggregateCreator
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
			name: "add second factor to login policy",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				existingIAM: &model.IAM{
					ObjectRoot:   models.ObjectRoot{AggregateID: "AggregateID"},
					IAMProjectID: "IAMProjectID",
					DefaultLoginPolicy: &model.LoginPolicy{
						AllowUsernamePassword: true,
					}},
				newMFA: &model.MFA{
					MFAType: int32(iam_model.SecondFactorTypeOTP),
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{model.LoginPolicySecondFactorAdded},
			},
		},
		{
			name: "existing iam nil",
			args: args{
				ctx:         authz.NewMockContext("orgID", "userID"),
				existingIAM: nil,
				aggCreator:  models.NewAggregateCreator("Test"),
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "mfa config nil",
			args: args{
				ctx:         authz.NewMockContext("orgID", "userID"),
				existingIAM: &model.IAM{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, IAMProjectID: "IAMProjectID"},
				newMFA:      nil,
				aggCreator:  models.NewAggregateCreator("Test"),
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := LoginPolicySecondFactorAddedAggregate(tt.args.aggCreator, tt.args.existingIAM, tt.args.newMFA)(tt.args.ctx)

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

func TestLoginPolicySecondFactorRemovedAggregate(t *testing.T) {
	type args struct {
		ctx         context.Context
		existingIAM *model.IAM
		mfa         *model.MFA
		aggCreator  *models.AggregateCreator
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
			name: "remove second factor to login policy",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				existingIAM: &model.IAM{
					ObjectRoot:   models.ObjectRoot{AggregateID: "AggregateID"},
					IAMProjectID: "IAMProjectID",
					DefaultLoginPolicy: &model.LoginPolicy{
						AllowUsernamePassword: true,
						SecondFactors: []int32{
							int32(iam_model.SecondFactorTypeOTP),
						},
					}},
				mfa: &model.MFA{
					MFAType: int32(iam_model.SecondFactorTypeOTP),
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{model.LoginPolicySecondFactorRemoved},
			},
		},
		{
			name: "existing iam nil",
			args: args{
				ctx:         authz.NewMockContext("orgID", "userID"),
				existingIAM: nil,
				aggCreator:  models.NewAggregateCreator("Test"),
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "idp config config nil",
			args: args{
				ctx:         authz.NewMockContext("orgID", "userID"),
				existingIAM: &model.IAM{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, IAMProjectID: "IAMProjectID"},
				mfa:         nil,
				aggCreator:  models.NewAggregateCreator("Test"),
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := LoginPolicySecondFactorRemovedAggregate(tt.args.aggCreator, tt.args.existingIAM, tt.args.mfa)(tt.args.ctx)

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

func TestLoginPolicyMultiFactorAddedAggregate(t *testing.T) {
	type args struct {
		ctx         context.Context
		existingIAM *model.IAM
		newMFA      *model.MFA
		aggCreator  *models.AggregateCreator
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
			name: "add mfa to login policy",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				existingIAM: &model.IAM{
					ObjectRoot:   models.ObjectRoot{AggregateID: "AggregateID"},
					IAMProjectID: "IAMProjectID",
					DefaultLoginPolicy: &model.LoginPolicy{
						AllowUsernamePassword: true,
					}},
				newMFA: &model.MFA{
					MFAType: int32(iam_model.MultiFactorTypeU2FWithPIN),
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{model.LoginPolicyMultiFactorAdded},
			},
		},
		{
			name: "existing iam nil",
			args: args{
				ctx:         authz.NewMockContext("orgID", "userID"),
				existingIAM: nil,
				aggCreator:  models.NewAggregateCreator("Test"),
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "mfa config nil",
			args: args{
				ctx:         authz.NewMockContext("orgID", "userID"),
				existingIAM: &model.IAM{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, IAMProjectID: "IAMProjectID"},
				newMFA:      nil,
				aggCreator:  models.NewAggregateCreator("Test"),
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := LoginPolicyMultiFactorAddedAggregate(tt.args.aggCreator, tt.args.existingIAM, tt.args.newMFA)(tt.args.ctx)

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

func TestLoginPolicyMultiFactorRemovedAggregate(t *testing.T) {
	type args struct {
		ctx         context.Context
		existingIAM *model.IAM
		mfa         *model.MFA
		aggCreator  *models.AggregateCreator
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
			name: "remove mfa to login policy",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				existingIAM: &model.IAM{
					ObjectRoot:   models.ObjectRoot{AggregateID: "AggregateID"},
					IAMProjectID: "IAMProjectID",
					DefaultLoginPolicy: &model.LoginPolicy{
						AllowUsernamePassword: true,
						SecondFactors: []int32{
							int32(iam_model.SecondFactorTypeOTP),
						},
					}},
				mfa: &model.MFA{
					MFAType: int32(iam_model.SecondFactorTypeOTP),
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{model.LoginPolicyMultiFactorRemoved},
			},
		},
		{
			name: "existing iam nil",
			args: args{
				ctx:         authz.NewMockContext("orgID", "userID"),
				existingIAM: nil,
				aggCreator:  models.NewAggregateCreator("Test"),
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "idp config config nil",
			args: args{
				ctx:         authz.NewMockContext("orgID", "userID"),
				existingIAM: &model.IAM{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, IAMProjectID: "IAMProjectID"},
				mfa:         nil,
				aggCreator:  models.NewAggregateCreator("Test"),
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := LoginPolicyMultiFactorRemovedAggregate(tt.args.aggCreator, tt.args.existingIAM, tt.args.mfa)(tt.args.ctx)

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

func TestPasswordComplexityPolicyAddedAggregate(t *testing.T) {
	type args struct {
		ctx         context.Context
		existingIAM *model.IAM
		newPolicy   *model.PasswordComplexityPolicy
		aggCreator  *models.AggregateCreator
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
			name: "add password complexity policy",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				existingIAM: &model.IAM{
					ObjectRoot:   models.ObjectRoot{AggregateID: "AggregateID"},
					IAMProjectID: "IAMProjectID"},
				newPolicy: &model.PasswordComplexityPolicy{
					ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"},
					MinLength:  10,
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{model.PasswordComplexityPolicyAdded},
			},
		},
		{
			name: "existing iam nil",
			args: args{
				ctx:         authz.NewMockContext("orgID", "userID"),
				existingIAM: nil,
				aggCreator:  models.NewAggregateCreator("Test"),
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "complexity policy config nil",
			args: args{
				ctx:         authz.NewMockContext("orgID", "userID"),
				existingIAM: &model.IAM{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, IAMProjectID: "IAMProjectID"},
				newPolicy:   nil,
				aggCreator:  models.NewAggregateCreator("Test"),
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := PasswordComplexityPolicyAddedAggregate(tt.args.ctx, tt.args.aggCreator, tt.args.existingIAM, tt.args.newPolicy)

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

func TestPasswordComplexityPolicyChangedAggregate(t *testing.T) {
	type args struct {
		ctx         context.Context
		existingIAM *model.IAM
		newPolicy   *model.PasswordComplexityPolicy
		aggCreator  *models.AggregateCreator
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
			name: "change password complexity policy",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				existingIAM: &model.IAM{
					ObjectRoot:   models.ObjectRoot{AggregateID: "AggregateID"},
					IAMProjectID: "IAMProjectID",
					DefaultPasswordComplexityPolicy: &model.PasswordComplexityPolicy{
						MinLength: 10,
					}},
				newPolicy: &model.PasswordComplexityPolicy{
					ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"},
					MinLength:  5,
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{model.PasswordComplexityPolicyChanged},
			},
		},
		{
			name: "no changes",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				existingIAM: &model.IAM{
					ObjectRoot:   models.ObjectRoot{AggregateID: "AggregateID"},
					IAMProjectID: "IAMProjectID",
					DefaultPasswordComplexityPolicy: &model.PasswordComplexityPolicy{
						MinLength: 10,
					}},
				newPolicy: &model.PasswordComplexityPolicy{
					ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"},
					MinLength:  10,
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing iam nil",
			args: args{
				ctx:         authz.NewMockContext("orgID", "userID"),
				existingIAM: nil,
				aggCreator:  models.NewAggregateCreator("Test"),
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "complexity policy config nil",
			args: args{
				ctx:         authz.NewMockContext("orgID", "userID"),
				existingIAM: &model.IAM{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, IAMProjectID: "IAMProjectID"},
				newPolicy:   nil,
				aggCreator:  models.NewAggregateCreator("Test"),
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := PasswordComplexityPolicyChangedAggregate(tt.args.aggCreator, tt.args.existingIAM, tt.args.newPolicy)(tt.args.ctx)

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

func TestPasswordAgePolicyAddedAggregate(t *testing.T) {
	type args struct {
		ctx         context.Context
		existingIAM *model.IAM
		newPolicy   *model.PasswordAgePolicy
		aggCreator  *models.AggregateCreator
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
			name: "add password age policy",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				existingIAM: &model.IAM{
					ObjectRoot:   models.ObjectRoot{AggregateID: "AggregateID"},
					IAMProjectID: "IAMProjectID"},
				newPolicy: &model.PasswordAgePolicy{
					ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"},
					MaxAgeDays: 10,
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{model.PasswordAgePolicyAdded},
			},
		},
		{
			name: "existing iam nil",
			args: args{
				ctx:         authz.NewMockContext("orgID", "userID"),
				existingIAM: nil,
				aggCreator:  models.NewAggregateCreator("Test"),
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "age policy config nil",
			args: args{
				ctx:         authz.NewMockContext("orgID", "userID"),
				existingIAM: &model.IAM{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, IAMProjectID: "IAMProjectID"},
				newPolicy:   nil,
				aggCreator:  models.NewAggregateCreator("Test"),
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := PasswordAgePolicyAddedAggregate(tt.args.ctx, tt.args.aggCreator, tt.args.existingIAM, tt.args.newPolicy)

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

func TestPasswordAgePolicyChangedAggregate(t *testing.T) {
	type args struct {
		ctx         context.Context
		existingIAM *model.IAM
		newPolicy   *model.PasswordAgePolicy
		aggCreator  *models.AggregateCreator
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
			name: "change password age policy",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				existingIAM: &model.IAM{
					ObjectRoot:   models.ObjectRoot{AggregateID: "AggregateID"},
					IAMProjectID: "IAMProjectID",
					DefaultPasswordAgePolicy: &model.PasswordAgePolicy{
						MaxAgeDays: 10,
					}},
				newPolicy: &model.PasswordAgePolicy{
					ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"},
					MaxAgeDays: 5,
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{model.PasswordAgePolicyChanged},
			},
		},
		{
			name: "no changes",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				existingIAM: &model.IAM{
					ObjectRoot:   models.ObjectRoot{AggregateID: "AggregateID"},
					IAMProjectID: "IAMProjectID",
					DefaultPasswordAgePolicy: &model.PasswordAgePolicy{
						MaxAgeDays: 10,
					}},
				newPolicy: &model.PasswordAgePolicy{
					ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"},
					MaxAgeDays: 10,
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing iam nil",
			args: args{
				ctx:         authz.NewMockContext("orgID", "userID"),
				existingIAM: nil,
				aggCreator:  models.NewAggregateCreator("Test"),
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "age policy config nil",
			args: args{
				ctx:         authz.NewMockContext("orgID", "userID"),
				existingIAM: &model.IAM{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, IAMProjectID: "IAMProjectID"},
				newPolicy:   nil,
				aggCreator:  models.NewAggregateCreator("Test"),
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := PasswordAgePolicyChangedAggregate(tt.args.aggCreator, tt.args.existingIAM, tt.args.newPolicy)(tt.args.ctx)

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

func TestPasswordLockoutPolicyAddedAggregate(t *testing.T) {
	type args struct {
		ctx         context.Context
		existingIAM *model.IAM
		newPolicy   *model.PasswordLockoutPolicy
		aggCreator  *models.AggregateCreator
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
			name: "add password lockout policy",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				existingIAM: &model.IAM{
					ObjectRoot:   models.ObjectRoot{AggregateID: "AggregateID"},
					IAMProjectID: "IAMProjectID"},
				newPolicy: &model.PasswordLockoutPolicy{
					ObjectRoot:  models.ObjectRoot{AggregateID: "AggregateID"},
					MaxAttempts: 10,
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{model.PasswordLockoutPolicyAdded},
			},
		},
		{
			name: "existing iam nil",
			args: args{
				ctx:         authz.NewMockContext("orgID", "userID"),
				existingIAM: nil,
				aggCreator:  models.NewAggregateCreator("Test"),
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "lockout policy config nil",
			args: args{
				ctx:         authz.NewMockContext("orgID", "userID"),
				existingIAM: &model.IAM{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, IAMProjectID: "IAMProjectID"},
				newPolicy:   nil,
				aggCreator:  models.NewAggregateCreator("Test"),
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := PasswordLockoutPolicyAddedAggregate(tt.args.ctx, tt.args.aggCreator, tt.args.existingIAM, tt.args.newPolicy)

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

func TestPasswordLockoutPolicyChangedAggregate(t *testing.T) {
	type args struct {
		ctx         context.Context
		existingIAM *model.IAM
		newPolicy   *model.PasswordLockoutPolicy
		aggCreator  *models.AggregateCreator
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
			name: "change password lockout policy",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				existingIAM: &model.IAM{
					ObjectRoot:   models.ObjectRoot{AggregateID: "AggregateID"},
					IAMProjectID: "IAMProjectID",
					DefaultPasswordLockoutPolicy: &model.PasswordLockoutPolicy{
						MaxAttempts: 10,
					}},
				newPolicy: &model.PasswordLockoutPolicy{
					ObjectRoot:  models.ObjectRoot{AggregateID: "AggregateID"},
					MaxAttempts: 5,
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{model.PasswordLockoutPolicyChanged},
			},
		},
		{
			name: "no changes",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				existingIAM: &model.IAM{
					ObjectRoot:   models.ObjectRoot{AggregateID: "AggregateID"},
					IAMProjectID: "IAMProjectID",
					DefaultPasswordLockoutPolicy: &model.PasswordLockoutPolicy{
						MaxAttempts: 10,
					}},
				newPolicy: &model.PasswordLockoutPolicy{
					ObjectRoot:  models.ObjectRoot{AggregateID: "AggregateID"},
					MaxAttempts: 10,
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing iam nil",
			args: args{
				ctx:         authz.NewMockContext("orgID", "userID"),
				existingIAM: nil,
				aggCreator:  models.NewAggregateCreator("Test"),
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "lockout policy config nil",
			args: args{
				ctx:         authz.NewMockContext("orgID", "userID"),
				existingIAM: &model.IAM{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, IAMProjectID: "IAMProjectID"},
				newPolicy:   nil,
				aggCreator:  models.NewAggregateCreator("Test"),
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := PasswordLockoutPolicyChangedAggregate(tt.args.aggCreator, tt.args.existingIAM, tt.args.newPolicy)(tt.args.ctx)

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

func TestOrgIAMPolicyAddedAggregate(t *testing.T) {
	type args struct {
		ctx         context.Context
		existingIAM *model.IAM
		newPolicy   *model.OrgIAMPolicy
		aggCreator  *models.AggregateCreator
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
			name: "add org iam policy",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				existingIAM: &model.IAM{
					ObjectRoot:   models.ObjectRoot{AggregateID: "AggregateID"},
					IAMProjectID: "IAMProjectID"},
				newPolicy: &model.OrgIAMPolicy{
					ObjectRoot:            models.ObjectRoot{AggregateID: "AggregateID"},
					UserLoginMustBeDomain: true,
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{model.OrgIAMPolicyAdded},
			},
		},
		{
			name: "existing iam nil",
			args: args{
				ctx:         authz.NewMockContext("orgID", "userID"),
				existingIAM: nil,
				aggCreator:  models.NewAggregateCreator("Test"),
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "lockout policy config nil",
			args: args{
				ctx:         authz.NewMockContext("orgID", "userID"),
				existingIAM: &model.IAM{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, IAMProjectID: "IAMProjectID"},
				newPolicy:   nil,
				aggCreator:  models.NewAggregateCreator("Test"),
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := OrgIAMPolicyAddedAggregate(tt.args.ctx, tt.args.aggCreator, tt.args.existingIAM, tt.args.newPolicy)

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

func TestOrgIAMPolicyChangedAggregate(t *testing.T) {
	type args struct {
		ctx         context.Context
		existingIAM *model.IAM
		newPolicy   *model.OrgIAMPolicy
		aggCreator  *models.AggregateCreator
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
			name: "change org iam policy",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				existingIAM: &model.IAM{
					ObjectRoot:   models.ObjectRoot{AggregateID: "AggregateID"},
					IAMProjectID: "IAMProjectID",
					DefaultOrgIAMPolicy: &model.OrgIAMPolicy{
						UserLoginMustBeDomain: true,
					}},
				newPolicy: &model.OrgIAMPolicy{
					ObjectRoot:            models.ObjectRoot{AggregateID: "AggregateID"},
					UserLoginMustBeDomain: false,
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:   1,
				eventTypes: []models.EventType{model.OrgIAMPolicyChanged},
			},
		},
		{
			name: "no changes",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				existingIAM: &model.IAM{
					ObjectRoot:   models.ObjectRoot{AggregateID: "AggregateID"},
					IAMProjectID: "IAMProjectID",
					DefaultOrgIAMPolicy: &model.OrgIAMPolicy{
						UserLoginMustBeDomain: true,
					}},
				newPolicy: &model.OrgIAMPolicy{
					ObjectRoot:            models.ObjectRoot{AggregateID: "AggregateID"},
					UserLoginMustBeDomain: true,
				},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing iam nil",
			args: args{
				ctx:         authz.NewMockContext("orgID", "userID"),
				existingIAM: nil,
				aggCreator:  models.NewAggregateCreator("Test"),
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "org iam policy config nil",
			args: args{
				ctx:         authz.NewMockContext("orgID", "userID"),
				existingIAM: &model.IAM{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, IAMProjectID: "IAMProjectID"},
				newPolicy:   nil,
				aggCreator:  models.NewAggregateCreator("Test"),
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := OrgIAMPolicyChangedAggregate(tt.args.aggCreator, tt.args.existingIAM, tt.args.newPolicy)(tt.args.ctx)

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
