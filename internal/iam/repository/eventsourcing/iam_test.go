package eventsourcing

import (
	"context"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"testing"

	"github.com/caos/zitadel/internal/api/authz"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
)

func TestSetUpStartedAggregate(t *testing.T) {
	type args struct {
		ctx        context.Context
		iam        *model.Iam
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
				iam:        &model.Iam{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.IamSetupStarted,
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
				eventType: model.IamSetupStarted,
				errFunc:   caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := IamSetupStartedAggregate(tt.args.aggCreator, tt.args.iam)(tt.args.ctx)

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
		ctx        context.Context
		existing   *model.Iam
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
			name: "setup done aggregate ok",
			args: args{
				ctx:        authz.NewMockContext("orgID", "userID"),
				existing:   &model.Iam{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.IamSetupDone,
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
				eventLen:  1,
				eventType: model.IamSetupDone,
				errFunc:   caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := IamSetupDoneAggregate(tt.args.aggCreator, tt.args.existing)(tt.args.ctx)

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
		ctx        context.Context
		existing   *model.Iam
		orgID      string
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
			name: "global org set aggregate ok",
			args: args{
				ctx:        authz.NewMockContext("orgID", "userID"),
				existing:   &model.Iam{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}},
				orgID:      "orgID",
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.GlobalOrgSet,
			},
		},
		{
			name: "existing iam nil",
			args: args{
				ctx:        authz.NewMockContext("orgID", "userID"),
				existing:   nil,
				orgID:      "orgID",
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "global org empty",
			args: args{
				ctx:        authz.NewMockContext("orgID", "userID"),
				existing:   &model.Iam{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := IamSetGlobalOrgAggregate(tt.args.aggCreator, tt.args.existing, tt.args.orgID)(tt.args.ctx)

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
		ctx        context.Context
		existing   *model.Iam
		projectID  string
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
			name: "iam project id set aggregate ok",
			args: args{
				ctx:        authz.NewMockContext("orgID", "userID"),
				existing:   &model.Iam{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}},
				projectID:  "projectID",
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.IamProjectSet,
			},
		},
		{
			name: "existing iam nil",
			args: args{
				ctx:        authz.NewMockContext("orgID", "userID"),
				existing:   nil,
				projectID:  "projectID",
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "project id empty",
			args: args{
				ctx:        authz.NewMockContext("orgID", "userID"),
				existing:   &model.Iam{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := IamSetIamProjectAggregate(tt.args.aggCreator, tt.args.existing, tt.args.projectID)(tt.args.ctx)

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
		ctx        context.Context
		existing   *model.Iam
		new        *model.IamMember
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
			name: "iammember added ok",
			args: args{
				ctx:        authz.NewMockContext("orgID", "userID"),
				existing:   &model.Iam{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}},
				new:        &model.IamMember{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, UserID: "UserID", Roles: []string{"Roles"}},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.IamMemberAdded,
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
				eventLen:  1,
				eventType: model.IamMemberAdded,
				errFunc:   caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "member nil",
			args: args{
				ctx:        authz.NewMockContext("orgID", "userID"),
				existing:   &model.Iam{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}},
				new:        nil,
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.IamMemberAdded,
				errFunc:   caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := IamMemberAddedAggregate(tt.args.aggCreator, tt.args.existing, tt.args.new)(tt.args.ctx)

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
		ctx        context.Context
		existing   *model.Iam
		new        *model.IamMember
		aggCreator *models.AggregateCreator
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
				ctx:        authz.NewMockContext("orgID", "userID"),
				existing:   &model.Iam{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}},
				new:        &model.IamMember{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, UserID: "UserID", Roles: []string{"Roles"}},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.IamMemberChanged,
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
				eventLen:  1,
				eventType: model.IamMemberChanged,
				wantErr:   true,
				errFunc:   caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "member nil",
			args: args{
				ctx:        authz.NewMockContext("orgID", "userID"),
				existing:   &model.Iam{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}},
				new:        nil,
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.IamMemberChanged,
				wantErr:   true,
				errFunc:   caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := IamMemberChangedAggregate(tt.args.aggCreator, tt.args.existing, tt.args.new)(tt.args.ctx)

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
		ctx        context.Context
		existing   *model.Iam
		new        *model.IamMember
		aggCreator *models.AggregateCreator
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
				ctx:        authz.NewMockContext("orgID", "userID"),
				existing:   &model.Iam{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}},
				new:        &model.IamMember{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, UserID: "UserID", Roles: []string{"Roles"}},
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.IamMemberRemoved,
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
				eventLen:  1,
				eventType: model.IamMemberRemoved,
				wantErr:   true,
				errFunc:   caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "member nil",
			args: args{
				ctx:        authz.NewMockContext("orgID", "userID"),
				existing:   &model.Iam{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}},
				new:        nil,
				aggCreator: models.NewAggregateCreator("Test"),
			},
			res: res{
				eventLen:  1,
				eventType: model.IamMemberRemoved,
				wantErr:   true,
				errFunc:   caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg, err := IamMemberRemovedAggregate(tt.args.aggCreator, tt.args.existing, tt.args.new)(tt.args.ctx)

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
		existing   *model.Iam
		new        *model.IDPConfig
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
				existing: &model.Iam{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, IamProjectID: "IamProjectID"},
				new: &model.IDPConfig{
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
				existing:   &model.Iam{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, IamProjectID: "IamProjectID"},
				new:        nil,
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
			agg, err := IDPConfigAddedAggregate(tt.args.aggCreator, tt.args.existing, tt.args.new)(tt.args.ctx)

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
		ctx        context.Context
		existing   *model.Iam
		new        *model.IDPConfig
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
			name: "change idp configuration",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				existing: &model.Iam{
					ObjectRoot:   models.ObjectRoot{AggregateID: "AggregateID"},
					IamProjectID: "IamProjectID",
					IDPs: []*model.IDPConfig{
						{IDPConfigID: "IDPConfigID", Name: "IDPName"},
					}},
				new: &model.IDPConfig{
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
				existing:   &model.Iam{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, IamProjectID: "IamProjectID"},
				new:        nil,
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
			agg, err := IDPConfigChangedAggregate(tt.args.aggCreator, tt.args.existing, tt.args.new)(tt.args.ctx)

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
		ctx        context.Context
		existing   *model.Iam
		new        *model.IDPConfig
		provider   *model.IDPProvider
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
			name: "remove idp config",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				existing: &model.Iam{
					ObjectRoot:   models.ObjectRoot{AggregateID: "AggregateID"},
					IamProjectID: "IamProjectID",
					IDPs: []*model.IDPConfig{
						{IDPConfigID: "IDPConfigID", Name: "Name"},
					}},
				new: &model.IDPConfig{
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
				existing: &model.Iam{
					ObjectRoot:   models.ObjectRoot{AggregateID: "AggregateID"},
					IamProjectID: "IamProjectID",
					IDPs: []*model.IDPConfig{
						{IDPConfigID: "IDPConfigID", Name: "Name"},
					}},
				new: &model.IDPConfig{
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
				existing:   &model.Iam{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, IamProjectID: "IamProjectID"},
				new:        nil,
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
			agg, err := IDPConfigRemovedAggregate(tt.args.ctx, tt.args.aggCreator, tt.args.existing, tt.args.new, tt.args.provider)

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
		ctx        context.Context
		existing   *model.Iam
		new        *model.IDPConfig
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
			name: "deactivate idp config",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				existing: &model.Iam{
					ObjectRoot:   models.ObjectRoot{AggregateID: "AggregateID"},
					IamProjectID: "IamProjectID",
					IDPs: []*model.IDPConfig{
						{IDPConfigID: "IDPConfigID", Name: "Name"},
					}},
				new: &model.IDPConfig{
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
				existing:   &model.Iam{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, IamProjectID: "IamProjectID"},
				new:        nil,
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
			agg, err := IDPConfigDeactivatedAggregate(tt.args.aggCreator, tt.args.existing, tt.args.new)(tt.args.ctx)

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
		ctx        context.Context
		existing   *model.Iam
		new        *model.IDPConfig
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
			name: "deactivate app",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				existing: &model.Iam{
					ObjectRoot:   models.ObjectRoot{AggregateID: "AggregateID"},
					IamProjectID: "IamProjectID",
					IDPs: []*model.IDPConfig{
						{IDPConfigID: "IDPConfigID", Name: "Name"},
					}},
				new: &model.IDPConfig{
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
				existing:   &model.Iam{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, IamProjectID: "IamProjectID"},
				new:        nil,
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
			agg, err := IDPConfigReactivatedAggregate(tt.args.aggCreator, tt.args.existing, tt.args.new)(tt.args.ctx)

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
		ctx        context.Context
		existing   *model.Iam
		new        *model.OIDCIDPConfig
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
			name: "change oidc config",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				existing: &model.Iam{
					ObjectRoot:   models.ObjectRoot{AggregateID: "AggregateID"},
					IamProjectID: "IamProjectID",
					IDPs: []*model.IDPConfig{
						{IDPConfigID: "IDPConfigID", Name: "Name", OIDCIDPConfig: &model.OIDCIDPConfig{IDPConfigID: "IDPConfigID", ClientID: "ClientID"}},
					}},
				new: &model.OIDCIDPConfig{
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
			name: "oidc config nil",
			args: args{
				ctx:        authz.NewMockContext("orgID", "userID"),
				existing:   &model.Iam{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, IamProjectID: "IamProjectID"},
				new:        nil,
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
			agg, err := OIDCIDPConfigChangedAggregate(tt.args.aggCreator, tt.args.existing, tt.args.new)(tt.args.ctx)

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
		ctx        context.Context
		existing   *model.Iam
		new        *model.LoginPolicy
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
			name: "add login polciy",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				existing: &model.Iam{
					ObjectRoot:   models.ObjectRoot{AggregateID: "AggregateID"},
					IamProjectID: "IamProjectID",
					IDPs: []*model.IDPConfig{
						{IDPConfigID: "IDPConfigID", Name: "Name", OIDCIDPConfig: &model.OIDCIDPConfig{IDPConfigID: "IDPConfigID", ClientID: "ClientID"}},
					}},
				new: &model.LoginPolicy{
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
			name: "login policy config nil",
			args: args{
				ctx:        authz.NewMockContext("orgID", "userID"),
				existing:   &model.Iam{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, IamProjectID: "IamProjectID"},
				new:        nil,
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
			agg, err := LoginPolicyAddedAggregate(tt.args.aggCreator, tt.args.existing, tt.args.new)(tt.args.ctx)

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
		ctx        context.Context
		existing   *model.Iam
		new        *model.LoginPolicy
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
			name: "change login policy",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				existing: &model.Iam{
					ObjectRoot:   models.ObjectRoot{AggregateID: "AggregateID"},
					IamProjectID: "IamProjectID",
					DefaultLoginPolicy: &model.LoginPolicy{
						AllowUsernamePassword: true,
					}},
				new: &model.LoginPolicy{
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
				existing: &model.Iam{
					ObjectRoot:   models.ObjectRoot{AggregateID: "AggregateID"},
					IamProjectID: "IamProjectID",
					DefaultLoginPolicy: &model.LoginPolicy{
						AllowUsernamePassword: true,
					}},
				new: &model.LoginPolicy{
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
			name: "login policy config nil",
			args: args{
				ctx:        authz.NewMockContext("orgID", "userID"),
				existing:   &model.Iam{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, IamProjectID: "IamProjectID"},
				new:        nil,
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
			agg, err := LoginPolicyChangedAggregate(tt.args.aggCreator, tt.args.existing, tt.args.new)(tt.args.ctx)

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
		ctx        context.Context
		existing   *model.Iam
		new        *model.IDPProvider
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
			name: "add idp provider to login policy",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				existing: &model.Iam{
					ObjectRoot:   models.ObjectRoot{AggregateID: "AggregateID"},
					IamProjectID: "IamProjectID",
					DefaultLoginPolicy: &model.LoginPolicy{
						AllowUsernamePassword: true,
					}},
				new: &model.IDPProvider{
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
			name: "idp config config nil",
			args: args{
				ctx:        authz.NewMockContext("orgID", "userID"),
				existing:   &model.Iam{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, IamProjectID: "IamProjectID"},
				new:        nil,
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
			agg, err := LoginPolicyIDPProviderAddedAggregate(tt.args.aggCreator, tt.args.existing, tt.args.new)(tt.args.ctx)

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
		ctx        context.Context
		existing   *model.Iam
		new        *model.IDPProviderID
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
			name: "remove idp provider to login policy",
			args: args{
				ctx: authz.NewMockContext("orgID", "userID"),
				existing: &model.Iam{
					ObjectRoot:   models.ObjectRoot{AggregateID: "AggregateID"},
					IamProjectID: "IamProjectID",
					DefaultLoginPolicy: &model.LoginPolicy{
						AllowUsernamePassword: true,
						IDPProviders: []*model.IDPProvider{
							{IDPConfigID: "IDPConfigID", Type: int32(iam_model.IDPProviderTypeSystem)},
						},
					}},
				new: &model.IDPProviderID{
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
			name: "idp config config nil",
			args: args{
				ctx:        authz.NewMockContext("orgID", "userID"),
				existing:   &model.Iam{ObjectRoot: models.ObjectRoot{AggregateID: "AggregateID"}, IamProjectID: "IamProjectID"},
				new:        nil,
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
			agg, err := LoginPolicyIDPProviderRemovedAggregate(tt.args.aggCreator, tt.args.existing, tt.args.new)(tt.args.ctx)

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
