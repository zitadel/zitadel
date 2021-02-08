package eventsourcing

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/caos/zitadel/internal/api/authz"
	caos_errs "github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
)

func TestIamByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es  *IAMEventstore
		iam *model.IAM
	}
	type res struct {
		iam     *model.IAM
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "iam from events, ok",
			args: args{
				es:  GetMockIAMByIDOK(ctrl),
				iam: &model.IAM{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				iam: &model.IAM{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
		},
		{
			name: "iam from events, no events",
			args: args{
				es:  GetMockIamByIDNoEvents(ctrl),
				iam: &model.IAM{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
		{
			name: "iam from events, no id",
			args: args{
				es:  GetMockIamByIDNoEvents(ctrl),
				iam: &model.IAM{ObjectRoot: es_models.ObjectRoot{AggregateID: "", Sequence: 1}},
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.IAMByID(nil, tt.args.iam.AggregateID)
			if (tt.res.errFunc != nil && !tt.res.errFunc(err)) || (err != nil && tt.res.errFunc == nil) {
				t.Errorf("got wrong err: %v ", err)
				return
			}
			if tt.res.errFunc != nil && tt.res.errFunc(err) {
				return
			}
			if result.AggregateID != tt.res.iam.AggregateID {
				t.Errorf("got wrong result name: expected: %v, actual: %v ", tt.res.iam.AggregateID, result.AggregateID)
			}
		})
	}
}

//
//func TestSetUpStarted(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	type args struct {
//		es    *IAMEventstore
//		ctx   context.Context
//		iamID string
//		step  iam_model.Step
//	}
//	type res struct {
//		iam     *iam_model.IAM
//		errFunc func(err error) bool
//	}
//	tests := []struct {
//		name string
//		args args
//		res  res
//	}{
//		{
//			name: "setup started iam, ok",
//			args: args{
//				es:    GetMockManipulateIAMNotExisting(ctrl),
//				ctx:   authz.NewMockContext("orgID", "userID"),
//				iamID: "iamID",
//				step:  iam_model.Step1,
//			},
//			res: res{
//				iam: &iam_model.IAM{ObjectRoot: es_models.ObjectRoot{AggregateID: "iamID", Sequence: 1}, SetUpStarted: iam_model.Step1},
//			},
//		},
//		{
//			name: "setup already started",
//			args: args{
//				es:    GetMockManipulateIAM(ctrl),
//				ctx:   authz.NewMockContext("orgID", "userID"),
//				iamID: "iamID",
//				step:  iam_model.Step1,
//			},
//			res: res{
//				errFunc: caos_errs.IsPreconditionFailed,
//			},
//		},
//		{
//			name: "setup iam no id",
//			args: args{
//				es:   GetMockManipulateIAM(ctrl),
//				ctx:  authz.NewMockContext("orgID", "userID"),
//				step: iam_model.Step1,
//			},
//			res: res{
//				errFunc: caos_errs.IsPreconditionFailed,
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			result, err := tt.args.es.StartSetup(tt.args.ctx, tt.args.iamID, tt.args.step)
//			if (tt.res.errFunc != nil && !tt.res.errFunc(err)) || (err != nil && tt.res.errFunc == nil) {
//				t.Errorf("got wrong err: %v ", err)
//				return
//			}
//			if tt.res.errFunc != nil && tt.res.errFunc(err) {
//				return
//			}
//			if result.AggregateID == "" {
//				t.Errorf("result has no id")
//			}
//			if result.SetUpStarted != tt.res.iam.SetUpStarted {
//				t.Errorf("got wrong result setupStarted: expected: %v, actual: %v ", tt.res.iam.SetUpStarted, result.SetUpStarted)
//			}
//		})
//	}
//}

func TestSetUpDone(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es    *IAMEventstore
		ctx   context.Context
		iamID string
		step  iam_model.Step
	}
	type res struct {
		iam     *iam_model.IAM
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "setup done iam, ok",
			args: args{
				es:    GetMockManipulateIAM(ctrl),
				ctx:   authz.NewMockContext("orgID", "userID"),
				iamID: "iamID",
				step:  iam_model.Step1,
			},
			res: res{
				iam: &iam_model.IAM{ObjectRoot: es_models.ObjectRoot{AggregateID: "iamID", Sequence: 1}, SetUpStarted: iam_model.Step1, SetUpDone: iam_model.Step1},
			},
		},
		{
			name: "setup iam no id",
			args: args{
				es:   GetMockManipulateIAM(ctrl),
				ctx:  authz.NewMockContext("orgID", "userID"),
				step: iam_model.Step1,
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "iam not found",
			args: args{
				es:    GetMockManipulateIAMNotExisting(ctrl),
				ctx:   authz.NewMockContext("orgID", "userID"),
				iamID: "iamID",
				step:  iam_model.Step1,
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.SetupDone(tt.args.ctx, tt.args.iamID, tt.args.step)
			if (tt.res.errFunc != nil && !tt.res.errFunc(err)) || (err != nil && tt.res.errFunc == nil) {
				t.Errorf("got wrong err: %v ", err)
				return
			}
			if tt.res.errFunc != nil && tt.res.errFunc(err) {
				return
			}
			if result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if result.SetUpDone != tt.res.iam.SetUpDone {
				t.Errorf("got wrong result SetUpDone: expected: %v, actual: %v ", tt.res.iam.SetUpDone, result.SetUpDone)
			}
		})
	}
}

func TestSetGlobalOrg(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es        *IAMEventstore
		ctx       context.Context
		iamID     string
		globalOrg string
	}
	type res struct {
		iam     *model.IAM
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "global org set, ok",
			args: args{
				es:        GetMockManipulateIAM(ctrl),
				ctx:       authz.NewMockContext("orgID", "userID"),
				iamID:     "iamID",
				globalOrg: "globalOrg",
			},
			res: res{
				iam: &model.IAM{ObjectRoot: es_models.ObjectRoot{AggregateID: "iamID", Sequence: 1}, SetUpStarted: model.Step1, GlobalOrgID: "globalOrg"},
			},
		},
		{
			name: "no iam id",
			args: args{
				es:        GetMockManipulateIAM(ctrl),
				ctx:       authz.NewMockContext("orgID", "userID"),
				globalOrg: "",
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "no global org",
			args: args{
				es:    GetMockManipulateIAM(ctrl),
				ctx:   authz.NewMockContext("orgID", "userID"),
				iamID: "iamID",
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "iam not found",
			args: args{
				es:        GetMockManipulateIAMNotExisting(ctrl),
				ctx:       authz.NewMockContext("orgID", "userID"),
				iamID:     "iamID",
				globalOrg: "globalOrg",
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.SetGlobalOrg(tt.args.ctx, tt.args.iamID, tt.args.globalOrg)
			if (tt.res.errFunc != nil && !tt.res.errFunc(err)) || (err != nil && tt.res.errFunc == nil) {
				t.Errorf("got wrong err: %v ", err)
				return
			}
			if tt.res.errFunc != nil && tt.res.errFunc(err) {
				return
			}
			if result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if result.GlobalOrgID != tt.res.iam.GlobalOrgID {
				t.Errorf("got wrong result GlobalOrgID: expected: %v, actual: %v ", tt.res.iam.GlobalOrgID, result.GlobalOrgID)
			}
		})
	}
}

func TestSetIamProjectID(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es           *IAMEventstore
		ctx          context.Context
		iamID        string
		iamProjectID string
	}
	type res struct {
		iam     *model.IAM
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "iam project set, ok",
			args: args{
				es:           GetMockManipulateIAM(ctrl),
				ctx:          authz.NewMockContext("orgID", "userID"),
				iamID:        "iamID",
				iamProjectID: "iamProjectID",
			},
			res: res{
				iam: &model.IAM{ObjectRoot: es_models.ObjectRoot{AggregateID: "iamID", Sequence: 1}, SetUpStarted: model.Step1, IAMProjectID: "iamProjectID"},
			},
		},
		{
			name: "no iam id",
			args: args{
				es:           GetMockManipulateIAM(ctrl),
				ctx:          authz.NewMockContext("orgID", "userID"),
				iamProjectID: "",
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "no global org",
			args: args{
				es:    GetMockManipulateIAM(ctrl),
				ctx:   authz.NewMockContext("orgID", "userID"),
				iamID: "iamID",
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "iam not found",
			args: args{
				es:           GetMockManipulateIAMNotExisting(ctrl),
				ctx:          authz.NewMockContext("orgID", "userID"),
				iamID:        "iamID",
				iamProjectID: "iamProjectID",
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.SetIAMProject(tt.args.ctx, tt.args.iamID, tt.args.iamProjectID)
			if (tt.res.errFunc != nil && !tt.res.errFunc(err)) || (err != nil && tt.res.errFunc == nil) {
				t.Errorf("got wrong err: %v ", err)
				return
			}
			if tt.res.errFunc != nil && tt.res.errFunc(err) {
				return
			}
			if result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if result.IAMProjectID != tt.res.iam.IAMProjectID {
				t.Errorf("got wrong result IAMProjectID: expected: %v, actual: %v ", tt.res.iam.IAMProjectID, result.IAMProjectID)
			}
		})
	}
}

func TestAddIamMember(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es     *IAMEventstore
		ctx    context.Context
		member *iam_model.IAMMember
	}
	type res struct {
		result  *iam_model.IAMMember
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "add iam member, ok",
			args: args{
				es:     GetMockManipulateIAM(ctrl),
				ctx:    authz.NewMockContext("orgID", "userID"),
				member: &iam_model.IAMMember{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, UserID: "UserID", Roles: []string{"Roles"}},
			},
			res: res{
				result: &iam_model.IAMMember{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, UserID: "UserID", Roles: []string{"Roles"}},
			},
		},
		{
			name: "no userid",
			args: args{
				es:     GetMockManipulateIAM(ctrl),
				ctx:    authz.NewMockContext("orgID", "userID"),
				member: &iam_model.IAMMember{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, Roles: []string{"Roles"}},
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "no roles",
			args: args{
				es:     GetMockManipulateIAM(ctrl),
				ctx:    authz.NewMockContext("orgID", "userID"),
				member: &iam_model.IAMMember{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, UserID: "UserID"},
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "member already existing",
			args: args{
				es:     GetMockManipulateIAMWithMember(ctrl),
				ctx:    authz.NewMockContext("orgID", "userID"),
				member: &iam_model.IAMMember{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, UserID: "UserID", Roles: []string{"Roles"}},
			},
			res: res{
				errFunc: caos_errs.IsErrorAlreadyExists,
			},
		},
		{
			name: "existing iam not found",
			args: args{
				es:     GetMockManipulateIAMNotExisting(ctrl),
				ctx:    authz.NewMockContext("orgID", "userID"),
				member: &iam_model.IAMMember{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, UserID: "UserID", Roles: []string{"Roles"}},
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.AddIAMMember(tt.args.ctx, tt.args.member)
			if (tt.res.errFunc != nil && !tt.res.errFunc(err)) || (err != nil && tt.res.errFunc == nil) {
				t.Errorf("got wrong err: %v ", err)
				return
			}
			if tt.res.errFunc != nil && tt.res.errFunc(err) {
				return
			}
			if result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if result.UserID != tt.res.result.UserID {
				t.Errorf("got wrong result userid: expected: %v, actual: %v ", tt.res.result.UserID, result.UserID)
			}
			if len(result.Roles) != len(tt.res.result.Roles) {
				t.Errorf("got wrong result roles: expected: %v, actual: %v ", tt.res.result.Roles, result.Roles)
			}
		})
	}
}

func TestChangeIamMember(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es     *IAMEventstore
		ctx    context.Context
		member *iam_model.IAMMember
	}
	type res struct {
		result  *iam_model.IAMMember
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "add iam member, ok",
			args: args{
				es:     GetMockManipulateIAMWithMember(ctrl),
				ctx:    authz.NewMockContext("orgID", "userID"),
				member: &iam_model.IAMMember{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, UserID: "UserID", Roles: []string{"ChangeRoles"}},
			},
			res: res{
				result: &iam_model.IAMMember{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, UserID: "UserID", Roles: []string{"Roles"}},
			},
		},
		{
			name: "no userid",
			args: args{
				es:     GetMockManipulateIAM(ctrl),
				ctx:    authz.NewMockContext("orgID", "userID"),
				member: &iam_model.IAMMember{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, Roles: []string{"ChangeRoles"}},
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "no roles",
			args: args{
				es:     GetMockManipulateIAM(ctrl),
				ctx:    authz.NewMockContext("orgID", "userID"),
				member: &iam_model.IAMMember{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, UserID: "UserID"},
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "member not existing",
			args: args{
				es:     GetMockManipulateIAM(ctrl),
				ctx:    authz.NewMockContext("orgID", "userID"),
				member: &iam_model.IAMMember{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, UserID: "UserID", Roles: []string{"Roles"}},
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing not found",
			args: args{
				es:     GetMockManipulateIAMNotExisting(ctrl),
				ctx:    authz.NewMockContext("orgID", "userID"),
				member: &iam_model.IAMMember{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, UserID: "UserID", Roles: []string{"ChangeRoles"}},
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.ChangeIAMMember(tt.args.ctx, tt.args.member)
			if (tt.res.errFunc != nil && !tt.res.errFunc(err)) || (err != nil && tt.res.errFunc == nil) {
				t.Errorf("got wrong err: %v ", err)
				return
			}
			if tt.res.errFunc != nil && tt.res.errFunc(err) {
				return
			}
			if result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if result.UserID != tt.res.result.UserID {
				t.Errorf("got wrong result userid: expected: %v, actual: %v ", tt.res.result.UserID, result.UserID)
			}
			if len(result.Roles) != len(tt.res.result.Roles) {
				t.Errorf("got wrong result roles: expected: %v, actual: %v ", tt.res.result.Roles, result.Roles)
			}
		})
	}
}

func TestRemoveIamMember(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es          *IAMEventstore
		ctx         context.Context
		existingIAM *model.IAM
		member      *iam_model.IAMMember
	}
	type res struct {
		result  *iam_model.IAMMember
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "remove iam member, ok",
			args: args{
				es:  GetMockManipulateIAMWithMember(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				existingIAM: &model.IAM{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					Members:    []*model.IAMMember{{UserID: "UserID", Roles: []string{"Roles"}}},
				},
				member: &iam_model.IAMMember{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, UserID: "UserID"},
			},
			res: res{
				result: &iam_model.IAMMember{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, UserID: "UserID", Roles: []string{"Roles"}},
			},
		},
		{
			name: "no userid",
			args: args{
				es:  GetMockManipulateIAM(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				existingIAM: &model.IAM{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					Members:    []*model.IAMMember{{UserID: "UserID", Roles: []string{"Roles"}}},
				},
				member: &iam_model.IAMMember{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, Roles: []string{"ChangeRoles"}},
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "member not existing",
			args: args{
				es:  GetMockManipulateIAM(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				existingIAM: &model.IAM{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
				},
				member: &iam_model.IAMMember{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, UserID: "UserID", Roles: []string{"Roles"}},
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing not found",
			args: args{
				es:     GetMockManipulateIAMNotExisting(ctrl),
				ctx:    authz.NewMockContext("orgID", "userID"),
				member: &iam_model.IAMMember{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, UserID: "UserID", Roles: []string{"ChangeRoles"}},
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.es.RemoveIAMMember(tt.args.ctx, tt.args.member)

			if tt.res.errFunc == nil && err != nil {
				t.Errorf("should not get err")
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestAddIdpConfiguration(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es  *IAMEventstore
		ctx context.Context
		idp *iam_model.IDPConfig
	}
	type res struct {
		result  *iam_model.IDPConfig
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "add idp, ok",
			args: args{
				es:  GetMockManipulateIAMWithCrypto(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				idp: &iam_model.IDPConfig{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					IDPConfigID: "IDPConfigID",
					Name:        "Name",
					Type:        iam_model.IDPConfigTypeOIDC,
					OIDCConfig: &iam_model.OIDCIDPConfig{
						ClientID:           "ClientID",
						ClientSecretString: "ClientSecret",
						Issuer:             "Issuer",
						Scopes:             []string{"scope"},
					},
				},
			},
			res: res{
				result: &iam_model.IDPConfig{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					Name: "Name",
					Type: iam_model.IDPConfigTypeOIDC,
					OIDCConfig: &iam_model.OIDCIDPConfig{
						ClientID: "ClientID",
						Issuer:   "Issuer",
						Scopes:   []string{"scope"},
					},
				},
			},
		},
		{
			name: "invalid idp config",
			args: args{
				es:  GetMockManipulateIAM(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				idp: &iam_model.IDPConfig{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing iam not found",
			args: args{
				es:  GetMockManipulateIAMNotExisting(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				idp: &iam_model.IDPConfig{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					IDPConfigID: "IDPConfigID",
					Name:        "Name",
					OIDCConfig: &iam_model.OIDCIDPConfig{
						ClientID:           "ClientID",
						ClientSecretString: "ClientSecret",
						Issuer:             "Issuer",
						Scopes:             []string{"scope"},
					},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.AddIDPConfig(tt.args.ctx, tt.args.idp)
			if (tt.res.wantErr && !tt.res.errFunc(err)) || (err != nil && !tt.res.wantErr) {
				t.Errorf("got wrong err: %v ", err)
				return
			}
			if tt.res.wantErr && tt.res.errFunc(err) {
				return
			}
			if result.IDPConfigID == "" {
				t.Errorf("result has no id")
			}
			if result.OIDCConfig.IDPConfigID == "" {
				t.Errorf("result has no id")
			}
			if result.OIDCConfig == nil && result.OIDCConfig.ClientSecret == nil {
				t.Errorf("result has no client secret")
			}
			if result.Name != tt.res.result.Name {
				t.Errorf("got wrong result key: expected: %v, actual: %v ", tt.res.result.Name, result.Name)
			}
			if result.OIDCConfig.ClientID != tt.res.result.OIDCConfig.ClientID {
				t.Errorf("got wrong result key: expected: %v, actual: %v ", tt.res.result.OIDCConfig.ClientID, result.OIDCConfig.ClientID)
			}
		})
	}
}

func TestChangeIdpConfiguration(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es  *IAMEventstore
		ctx context.Context
		idp *iam_model.IDPConfig
	}
	type res struct {
		result  *iam_model.IDPConfig
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "change idp, ok",
			args: args{
				es:  GetMockManipulateIAMWithOIDCIdp(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				idp: &iam_model.IDPConfig{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					IDPConfigID: "IDPConfigID",
					Name:        "NameChanged",
				},
			},
			res: res{
				result: &iam_model.IDPConfig{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					IDPConfigID: "IDPConfigID",
					Name:        "NameChanged",
					OIDCConfig: &iam_model.OIDCIDPConfig{
						ClientID: "ClientID",
					},
				},
			},
		},
		{
			name: "invalid idp",
			args: args{
				es:  GetMockManipulateIAM(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				idp: &iam_model.IDPConfig{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					IDPConfigID: "IDPConfigID",
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "idp not existing",
			args: args{
				es:  GetMockManipulateIAM(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				idp: &iam_model.IDPConfig{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					IDPConfigID: "IDPConfigID",
					Name:        "Name",
					OIDCConfig: &iam_model.OIDCIDPConfig{
						ClientID: "ClientID",
					},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing project not found",
			args: args{
				es:  GetMockManipulateIAMNotExisting(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				idp: &iam_model.IDPConfig{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					IDPConfigID: "IDPConfigID",
					Name:        "Name",
					OIDCConfig: &iam_model.OIDCIDPConfig{
						ClientID: "ClientID",
					},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.ChangeIDPConfig(tt.args.ctx, tt.args.idp)
			if (tt.res.wantErr && !tt.res.errFunc(err)) || (err != nil && !tt.res.wantErr) {
				t.Errorf("got wrong err: %v ", err)
				return
			}
			if tt.res.wantErr && tt.res.errFunc(err) {
				return
			}
			if !tt.res.wantErr && result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if result.IDPConfigID != tt.res.result.IDPConfigID {
				t.Errorf("got wrong result AppID: expected: %v, actual: %v ", tt.res.result.IDPConfigID, result.IDPConfigID)
			}
			if result.Name != tt.res.result.Name {
				t.Errorf("got wrong result name: expected: %v, actual: %v ", tt.res.result.Name, result.Name)
			}
		})
	}
}

func TestRemoveIdpConfiguration(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es  *IAMEventstore
		ctx context.Context
		idp *iam_model.IDPConfig
	}
	type res struct {
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "remove idp, ok",
			args: args{
				es:  GetMockManipulateIAMWithOIDCIdp(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				idp: &iam_model.IDPConfig{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					IDPConfigID: "IDPConfigID",
				},
			},
		},
		{
			name: "no IDPConfigID",
			args: args{
				es:  GetMockManipulateIAM(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				idp: &iam_model.IDPConfig{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "idp not existing",
			args: args{
				es:  GetMockManipulateIAM(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				idp: &iam_model.IDPConfig{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					IDPConfigID: "IDPConfigID",
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing idp not found",
			args: args{
				es:  GetMockManipulateIAMNotExisting(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				idp: &iam_model.IDPConfig{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					IDPConfigID: "IDPConfigID",
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.es.RemoveIDPConfig(tt.args.ctx, tt.args.idp)

			if !tt.res.wantErr && err != nil {
				t.Errorf("should not get err")
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}
func TestDeactivateIdpConfiguration(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es  *IAMEventstore
		ctx context.Context
		idp *iam_model.IDPConfig
	}
	type res struct {
		result  *iam_model.IDPConfig
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "deactivate, ok",
			args: args{
				es:  GetMockManipulateIAMWithOIDCIdp(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				idp: &iam_model.IDPConfig{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					IDPConfigID: "IDPConfigID",
					Name:        "Name",
				},
			},
			res: res{
				result: &iam_model.IDPConfig{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					IDPConfigID: "IDPConfigID",
					Name:        "Name",
					State:       iam_model.IDPConfigStateInactive,
					OIDCConfig: &iam_model.OIDCIDPConfig{
						ClientID: "ClientID",
					},
				},
			},
		},
		{
			name: "no idp id",
			args: args{
				es:  GetMockManipulateIAM(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				idp: &iam_model.IDPConfig{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "idp not existing",
			args: args{
				es:  GetMockManipulateIAM(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				idp: &iam_model.IDPConfig{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					IDPConfigID: "IDPConfigID",
					Name:        "Name",
					OIDCConfig: &iam_model.OIDCIDPConfig{
						ClientID: "ClientID",
					},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing iam not found",
			args: args{
				es:  GetMockManipulateIAMNotExisting(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				idp: &iam_model.IDPConfig{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					IDPConfigID: "IDPConfigID",
					Name:        "Name",
					OIDCConfig: &iam_model.OIDCIDPConfig{
						ClientID: "ClientID",
					},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.DeactivateIDPConfig(tt.args.ctx, tt.args.idp.AggregateID, tt.args.idp.IDPConfigID)
			if (tt.res.wantErr && !tt.res.errFunc(err)) || (err != nil && !tt.res.wantErr) {
				t.Errorf("got wrong err: %v ", err)
				return
			}
			if tt.res.wantErr && tt.res.errFunc(err) {
				return
			}
			if !tt.res.wantErr && result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if result.IDPConfigID != tt.res.result.IDPConfigID {
				t.Errorf("got wrong result IDPConfigID: expected: %v, actual: %v ", tt.res.result.IDPConfigID, result.IDPConfigID)
			}
			if result.State != tt.res.result.State {
				t.Errorf("got wrong result state: expected: %v, actual: %v ", tt.res.result.State, result.State)
			}
		})
	}
}

func TestReactivateIdpConfiguration(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es  *IAMEventstore
		ctx context.Context
		idp *iam_model.IDPConfig
	}
	type res struct {
		result  *iam_model.IDPConfig
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "reactivate, ok",
			args: args{
				es:  GetMockManipulateIAMWithOIDCIdp(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				idp: &iam_model.IDPConfig{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					IDPConfigID: "IDPConfigID",
					Name:        "Name",
				},
			},
			res: res{
				result: &iam_model.IDPConfig{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					IDPConfigID: "IDPConfigID",
					Name:        "Name",
					State:       iam_model.IDPConfigStateActive,
					OIDCConfig: &iam_model.OIDCIDPConfig{
						ClientID: "ClientID",
					},
				},
			},
		},
		{
			name: "no idp id",
			args: args{
				es:  GetMockManipulateIAM(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				idp: &iam_model.IDPConfig{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "idp not existing",
			args: args{
				es:  GetMockManipulateIAM(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				idp: &iam_model.IDPConfig{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					IDPConfigID: "IDPConfigID",
					Name:        "Name",
					OIDCConfig: &iam_model.OIDCIDPConfig{
						ClientID: "ClientID",
					},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing iam not found",
			args: args{
				es:  GetMockManipulateIAMNotExisting(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				idp: &iam_model.IDPConfig{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					IDPConfigID: "IDPConfigID",
					Name:        "Name",
					OIDCConfig: &iam_model.OIDCIDPConfig{
						ClientID: "ClientID",
					},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.ReactivateIDPConfig(tt.args.ctx, tt.args.idp.AggregateID, tt.args.idp.IDPConfigID)
			if (tt.res.wantErr && !tt.res.errFunc(err)) || (err != nil && !tt.res.wantErr) {
				t.Errorf("got wrong err: %v ", err)
				return
			}
			if tt.res.wantErr && tt.res.errFunc(err) {
				return
			}
			if !tt.res.wantErr && result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if result.IDPConfigID != tt.res.result.IDPConfigID {
				t.Errorf("got wrong result IDPConfigID: expected: %v, actual: %v ", tt.res.result.IDPConfigID, result.IDPConfigID)
			}
			if result.State != tt.res.result.State {
				t.Errorf("got wrong result state: expected: %v, actual: %v ", tt.res.result.State, result.State)
			}
		})
	}
}

func TestChangeOIDCIDPConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es     *IAMEventstore
		ctx    context.Context
		config *iam_model.OIDCIDPConfig
	}
	type res struct {
		result  *iam_model.OIDCIDPConfig
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "change oidc config, ok",
			args: args{
				es:  GetMockManipulateIAMWithOIDCIdp(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				config: &iam_model.OIDCIDPConfig{
					ObjectRoot:  es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					IDPConfigID: "IDPConfigID",
					ClientID:    "ClientIDChange",
					Issuer:      "Issuer",
					Scopes:      []string{"scope"},
				},
			},
			res: res{
				result: &iam_model.OIDCIDPConfig{
					ObjectRoot:  es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					IDPConfigID: "IDPConfigID",
					ClientID:    "ClientIDChange",
				},
			},
		},
		{
			name: "invalid config",
			args: args{
				es:  GetMockManipulateIAM(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				config: &iam_model.OIDCIDPConfig{
					ObjectRoot:  es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					IDPConfigID: "IDPConfigID",
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "idp not existing",
			args: args{
				es:  GetMockManipulateIAM(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				config: &iam_model.OIDCIDPConfig{
					ObjectRoot:  es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					IDPConfigID: "IDPConfigID",
					ClientID:    "ClientID",
					Issuer:      "Issuer",
					Scopes:      []string{"scope"},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing iam not found",
			args: args{
				es:  GetMockManipulateIAMNotExisting(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				config: &iam_model.OIDCIDPConfig{
					ObjectRoot:  es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					IDPConfigID: "IDPConfigID",
					ClientID:    "ClientID",
					Issuer:      "Issuer",
					Scopes:      []string{"scope"},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.ChangeIDPOIDCConfig(tt.args.ctx, tt.args.config)
			if (tt.res.wantErr && !tt.res.errFunc(err)) || (err != nil && !tt.res.wantErr) {
				t.Errorf("got wrong err: %v ", err)
				return
			}
			if tt.res.wantErr && tt.res.errFunc(err) {
				return
			}
			if result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if result.IDPConfigID != tt.res.result.IDPConfigID {
				t.Errorf("got wrong result AppID: expected: %v, actual: %v ", tt.res.result.IDPConfigID, result.IDPConfigID)
			}
			if result.ClientID != tt.res.result.ClientID {
				t.Errorf("got wrong result responsetype: expected: %v, actual: %v ", tt.res.result.ClientID, result.ClientID)
			}
		})
	}
}

func TestAddLoginPolicy(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es     *IAMEventstore
		ctx    context.Context
		policy *iam_model.LoginPolicy
	}
	type res struct {
		result  *iam_model.LoginPolicy
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "add login policy, ok",
			args: args{
				es:  GetMockManipulateIAM(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.LoginPolicy{
					ObjectRoot:    es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					AllowRegister: true,
				},
			},
			res: res{
				result: &iam_model.LoginPolicy{
					ObjectRoot:    es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					AllowRegister: true,
				},
			},
		},
		{
			name: "invalid policy",
			args: args{
				es:  GetMockManipulateIAM(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.LoginPolicy{
					ObjectRoot: es_models.ObjectRoot{Sequence: 0},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing iam not found",
			args: args{
				es:  GetMockManipulateIAMNotExisting(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.LoginPolicy{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.AddLoginPolicy(tt.args.ctx, tt.args.policy)
			if (tt.res.wantErr && !tt.res.errFunc(err)) || (err != nil && !tt.res.wantErr) {
				t.Errorf("got wrong err: %v ", err)
				return
			}
			if tt.res.wantErr && tt.res.errFunc(err) {
				return
			}
			if result.AllowRegister != tt.res.result.AllowRegister {
				t.Errorf("got wrong result AllowRegister: expected: %v, actual: %v ", tt.res.result.AllowRegister, result.AllowRegister)
			}
		})
	}
}

func TestChangeLoginPolicy(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es     *IAMEventstore
		ctx    context.Context
		policy *iam_model.LoginPolicy
	}
	type res struct {
		result  *iam_model.LoginPolicy
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "add login policy, ok",
			args: args{
				es:  GetMockManipulateIAMWithLoginPolicy(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.LoginPolicy{
					ObjectRoot:            es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					AllowRegister:         true,
					AllowExternalIdp:      false,
					AllowUsernamePassword: false,
				},
			},
			res: res{
				result: &iam_model.LoginPolicy{
					ObjectRoot:            es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					AllowRegister:         true,
					AllowExternalIdp:      false,
					AllowUsernamePassword: false,
				},
			},
		},
		{
			name: "invalid policy",
			args: args{
				es:  GetMockManipulateIAM(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.LoginPolicy{
					ObjectRoot: es_models.ObjectRoot{Sequence: 0},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing iam not found",
			args: args{
				es:  GetMockManipulateIAMNotExisting(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.LoginPolicy{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.ChangeLoginPolicy(tt.args.ctx, tt.args.policy)
			if (tt.res.wantErr && !tt.res.errFunc(err)) || (err != nil && !tt.res.wantErr) {
				t.Errorf("got wrong err: %v ", err)
				return
			}
			if tt.res.wantErr && tt.res.errFunc(err) {
				return
			}
			if result.AllowRegister != tt.res.result.AllowRegister {
				t.Errorf("got wrong result AllowRegister: expected: %v, actual: %v ", tt.res.result.AllowRegister, result.AllowRegister)
			}
			if result.AllowUsernamePassword != tt.res.result.AllowUsernamePassword {
				t.Errorf("got wrong result AllowUsernamePassword: expected: %v, actual: %v ", tt.res.result.AllowUsernamePassword, result.AllowUsernamePassword)
			}
			if result.AllowExternalIdp != tt.res.result.AllowExternalIdp {
				t.Errorf("got wrong result AllowExternalIDP: expected: %v, actual: %v ", tt.res.result.AllowExternalIdp, result.AllowExternalIdp)
			}
		})
	}
}

func TestAddIdpProviderToLoginPolicy(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es       *IAMEventstore
		ctx      context.Context
		provider *iam_model.IDPProvider
	}
	type res struct {
		result  *iam_model.IDPProvider
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "add idp to login policy, ok",
			args: args{
				es:  GetMockManipulateIAMWithLoginPolicy(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				provider: &iam_model.IDPProvider{
					ObjectRoot:  es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					IDPConfigID: "IdpConfigID2",
					Type:        iam_model.IDPProviderTypeSystem,
				},
			},
			res: res{
				result: &iam_model.IDPProvider{IDPConfigID: "IdpConfigID2"},
			},
		},
		{
			name: "add idp to login policy, already existing",
			args: args{
				es:  GetMockManipulateIAMWithLoginPolicy(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				provider: &iam_model.IDPProvider{
					ObjectRoot:  es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					IDPConfigID: "IDPConfigID",
					Type:        iam_model.IDPProviderTypeSystem,
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsErrorAlreadyExists,
			},
		},
		{
			name: "invalid provider",
			args: args{
				es:  GetMockManipulateIAM(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				provider: &iam_model.IDPProvider{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing iam not found",
			args: args{
				es:  GetMockManipulateIAMNotExisting(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				provider: &iam_model.IDPProvider{
					ObjectRoot:  es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					IDPConfigID: "IdpConfigID2",
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.AddIDPProviderToLoginPolicy(tt.args.ctx, tt.args.provider)
			if (tt.res.wantErr && !tt.res.errFunc(err)) || (err != nil && !tt.res.wantErr) {
				t.Errorf("got wrong err: %v ", err)
				return
			}
			if tt.res.wantErr && tt.res.errFunc(err) {
				return
			}
			if result.IDPConfigID != tt.res.result.IDPConfigID {
				t.Errorf("got wrong result IDPConfigID: expected: %v, actual: %v ", tt.res.result.IDPConfigID, result.IDPConfigID)
			}
			if result.Type != tt.res.result.Type {
				t.Errorf("got wrong result KeyType: expected: %v, actual: %v ", tt.res.result.Type, result.Type)
			}
		})
	}
}

func TestRemoveIdpProviderFromLoginPolicy(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es       *IAMEventstore
		ctx      context.Context
		provider *iam_model.IDPProvider
	}
	type res struct {
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "remove idp to login policy, ok",
			args: args{
				es:  GetMockManipulateIAMWithLoginPolicy(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				provider: &iam_model.IDPProvider{
					ObjectRoot:  es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					IDPConfigID: "IDPConfigID",
					Type:        iam_model.IDPProviderTypeSystem,
				},
			},
			res: res{},
		},
		{
			name: "remove idp to login policy, not existing",
			args: args{
				es:  GetMockManipulateIAMWithLoginPolicy(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				provider: &iam_model.IDPProvider{
					ObjectRoot:  es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					IDPConfigID: "IdpConfigID2",
					Type:        iam_model.IDPProviderTypeSystem,
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "invalid provider",
			args: args{
				es:  GetMockManipulateIAM(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				provider: &iam_model.IDPProvider{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing iam not found",
			args: args{
				es:  GetMockManipulateIAMNotExisting(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				provider: &iam_model.IDPProvider{
					ObjectRoot:  es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					IDPConfigID: "IdpConfigID2",
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.es.RemoveIDPProviderFromLoginPolicy(tt.args.ctx, tt.args.provider)

			if !tt.res.wantErr && err != nil {
				t.Errorf("should not get err: %v ", err)
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestAddSecondFactorToLoginPolicy(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es         *IAMEventstore
		ctx        context.Context
		aggreageID string
		mfa        iam_model.SecondFactorType
	}
	type res struct {
		result  iam_model.SecondFactorType
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "add second factor to login policy, ok",
			args: args{
				es:         GetMockManipulateIAMWithLoginPolicy(ctrl),
				ctx:        authz.NewMockContext("orgID", "userID"),
				aggreageID: "AggregateID",
				mfa:        iam_model.SecondFactorTypeOTP,
			},
			res: res{
				result: iam_model.SecondFactorTypeOTP,
			},
		},
		{
			name: "add second factor to login policy, already existing",
			args: args{
				es:         GetMockManipulateIAMWithLoginPolicyWithMFAs(ctrl),
				ctx:        authz.NewMockContext("orgID", "userID"),
				aggreageID: "AggregateID",
				mfa:        iam_model.SecondFactorTypeOTP,
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsErrorAlreadyExists,
			},
		},
		{
			name: "invalid mfa",
			args: args{
				es:  GetMockManipulateIAM(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				mfa: iam_model.SecondFactorTypeUnspecified,
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing iam not found",
			args: args{
				es:         GetMockManipulateIAMNotExisting(ctrl),
				ctx:        authz.NewMockContext("orgID", "userID"),
				aggreageID: "Test",
				mfa:        iam_model.SecondFactorTypeOTP,
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.AddSecondFactorToLoginPolicy(tt.args.ctx, tt.args.aggreageID, tt.args.mfa)
			if (tt.res.wantErr && !tt.res.errFunc(err)) || (err != nil && !tt.res.wantErr) {
				t.Errorf("got wrong err: %v ", err)
				return
			}
			if tt.res.wantErr && tt.res.errFunc(err) {
				return
			}
			if result != tt.res.result {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.res.result, result)
			}
		})
	}
}

func TestRemoveSecondFactorFromLoginPolicy(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es          *IAMEventstore
		ctx         context.Context
		aggregateID string
		mfa         iam_model.SecondFactorType
	}
	type res struct {
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "remove second factor from login policy, ok",
			args: args{
				es:          GetMockManipulateIAMWithLoginPolicyWithMFAs(ctrl),
				ctx:         authz.NewMockContext("orgID", "userID"),
				aggregateID: "AggregateID",
				mfa:         iam_model.SecondFactorTypeOTP,
			},
			res: res{},
		},
		{
			name: "remove second factor to login policy, not existing",
			args: args{
				es:          GetMockManipulateIAMWithLoginPolicy(ctrl),
				ctx:         authz.NewMockContext("orgID", "userID"),
				aggregateID: "AggregateID",
				mfa:         iam_model.SecondFactorTypeOTP,
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "invalid provider",
			args: args{
				es:          GetMockManipulateIAM(ctrl),
				ctx:         authz.NewMockContext("orgID", "userID"),
				aggregateID: "AggregateID",
				mfa:         iam_model.SecondFactorTypeUnspecified,
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing iam not found",
			args: args{
				es:          GetMockManipulateIAMNotExisting(ctrl),
				ctx:         authz.NewMockContext("orgID", "userID"),
				aggregateID: "Test",
				mfa:         iam_model.SecondFactorTypeOTP,
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.es.RemoveSecondFactorFromLoginPolicy(tt.args.ctx, tt.args.aggregateID, tt.args.mfa)

			if !tt.res.wantErr && err != nil {
				t.Errorf("should not get err: %v ", err)
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestAddMultiFactorToLoginPolicy(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es         *IAMEventstore
		ctx        context.Context
		aggreageID string
		mfa        iam_model.MultiFactorType
	}
	type res struct {
		result  iam_model.MultiFactorType
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "add second factor to login policy, ok",
			args: args{
				es:         GetMockManipulateIAMWithLoginPolicy(ctrl),
				ctx:        authz.NewMockContext("orgID", "userID"),
				aggreageID: "AggregateID",
				mfa:        iam_model.MultiFactorTypeU2FWithPIN,
			},
			res: res{
				result: iam_model.MultiFactorTypeU2FWithPIN,
			},
		},
		{
			name: "add second factor to login policy, already existing",
			args: args{
				es:         GetMockManipulateIAMWithLoginPolicyWithMFAs(ctrl),
				ctx:        authz.NewMockContext("orgID", "userID"),
				aggreageID: "AggregateID",
				mfa:        iam_model.MultiFactorTypeU2FWithPIN,
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsErrorAlreadyExists,
			},
		},
		{
			name: "invalid mfa",
			args: args{
				es:  GetMockManipulateIAM(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				mfa: iam_model.MultiFactorTypeUnspecified,
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing iam not found",
			args: args{
				es:         GetMockManipulateIAMNotExisting(ctrl),
				ctx:        authz.NewMockContext("orgID", "userID"),
				aggreageID: "Test",
				mfa:        iam_model.MultiFactorTypeU2FWithPIN,
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.AddMultiFactorToLoginPolicy(tt.args.ctx, tt.args.aggreageID, tt.args.mfa)
			if (tt.res.wantErr && !tt.res.errFunc(err)) || (err != nil && !tt.res.wantErr) {
				t.Errorf("got wrong err: %v ", err)
				return
			}
			if tt.res.wantErr && tt.res.errFunc(err) {
				return
			}
			if result != tt.res.result {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.res.result, result)
			}
		})
	}
}

func TestRemoveMultiFactorFromLoginPolicy(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es          *IAMEventstore
		ctx         context.Context
		aggregateID string
		mfa         iam_model.MultiFactorType
	}
	type res struct {
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "remove second factor from login policy, ok",
			args: args{
				es:          GetMockManipulateIAMWithLoginPolicyWithMFAs(ctrl),
				ctx:         authz.NewMockContext("orgID", "userID"),
				aggregateID: "AggregateID",
				mfa:         iam_model.MultiFactorTypeU2FWithPIN,
			},
			res: res{},
		},
		{
			name: "remove second factor to login policy, not existing",
			args: args{
				es:          GetMockManipulateIAMWithLoginPolicy(ctrl),
				ctx:         authz.NewMockContext("orgID", "userID"),
				aggregateID: "AggregateID",
				mfa:         iam_model.MultiFactorTypeU2FWithPIN,
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "invalid provider",
			args: args{
				es:          GetMockManipulateIAM(ctrl),
				ctx:         authz.NewMockContext("orgID", "userID"),
				aggregateID: "AggregateID",
				mfa:         iam_model.MultiFactorTypeUnspecified,
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing iam not found",
			args: args{
				es:          GetMockManipulateIAMNotExisting(ctrl),
				ctx:         authz.NewMockContext("orgID", "userID"),
				aggregateID: "Test",
				mfa:         iam_model.MultiFactorTypeU2FWithPIN,
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.es.RemoveMultiFactorFromLoginPolicy(tt.args.ctx, tt.args.aggregateID, tt.args.mfa)

			if !tt.res.wantErr && err != nil {
				t.Errorf("should not get err: %v ", err)
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestAddLabelPolicy(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es     *IAMEventstore
		ctx    context.Context
		policy *iam_model.LabelPolicy
	}
	type res struct {
		result  *iam_model.LabelPolicy
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "add label policy, ok",
			args: args{
				es:  GetMockManipulateIAM(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.LabelPolicy{
					ObjectRoot:   es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					PrimaryColor: "000000",
				},
			},
			res: res{
				result: &iam_model.LabelPolicy{
					ObjectRoot:   es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					PrimaryColor: "000000",
				},
			},
		},
		{
			name: "invalid policy",
			args: args{
				es:  GetMockManipulateIAM(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.LabelPolicy{
					ObjectRoot: es_models.ObjectRoot{Sequence: 0},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing iam not found",
			args: args{
				es:  GetMockManipulateIAMNotExisting(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.LabelPolicy{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.AddLabelPolicy(tt.args.ctx, tt.args.policy)
			if (tt.res.wantErr && !tt.res.errFunc(err)) || (err != nil && !tt.res.wantErr) {
				t.Errorf("got wrong err: %v ", err)
				return
			}
			if tt.res.wantErr && tt.res.errFunc(err) {
				return
			}
			if result.PrimaryColor != tt.res.result.PrimaryColor {
				t.Errorf("got wrong result PrimaryColor: expected: %v, actual: %v ", tt.res.result.PrimaryColor, result.PrimaryColor)
			}
		})
	}
}

func TestChangeLabelPolicy(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es     *IAMEventstore
		ctx    context.Context
		policy *iam_model.LabelPolicy
	}
	type res struct {
		result  *iam_model.LabelPolicy
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "change label policy, ok",
			args: args{
				es:  GetMockManipulateIAMWithLabelPolicy(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.LabelPolicy{
					ObjectRoot:     es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					PrimaryColor:   "000000",
					SecondaryColor: "FFFFFF",
				},
			},
			res: res{
				result: &iam_model.LabelPolicy{
					ObjectRoot:     es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					PrimaryColor:   "000000",
					SecondaryColor: "FFFFFF",
				},
			},
		},
		{
			name: "invalid policy",
			args: args{
				es:  GetMockManipulateIAM(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.LabelPolicy{
					ObjectRoot: es_models.ObjectRoot{Sequence: 0},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing iam not found",
			args: args{
				es:  GetMockManipulateIAMNotExisting(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.LabelPolicy{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.ChangeLabelPolicy(tt.args.ctx, tt.args.policy)
			if (tt.res.wantErr && !tt.res.errFunc(err)) || (err != nil && !tt.res.wantErr) {
				t.Errorf("got wrong err: %v ", err)
				return
			}
			if tt.res.wantErr && tt.res.errFunc(err) {
				return
			}
			if result.PrimaryColor != tt.res.result.PrimaryColor {
				t.Errorf("got wrong result PrimaryColor: expected: %v, actual: %v ", tt.res.result.PrimaryColor, result.PrimaryColor)
			}
			if result.SecondaryColor != tt.res.result.SecondaryColor {
				t.Errorf("got wrong result SecondaryColor: expected: %v, actual: %v ", tt.res.result.SecondaryColor, result.SecondaryColor)
			}
		})
	}
}
func TestAddPasswordComplexityPolicy(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es     *IAMEventstore
		ctx    context.Context
		policy *iam_model.PasswordComplexityPolicy
	}
	type res struct {
		result  *iam_model.PasswordComplexityPolicy
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "add password complexity policy, ok",
			args: args{
				es:  GetMockManipulateIAM(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.PasswordComplexityPolicy{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					MinLength:  10,
				},
			},
			res: res{
				result: &iam_model.PasswordComplexityPolicy{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					MinLength:  10,
				},
			},
		},
		{
			name: "invalid policy",
			args: args{
				es:  GetMockManipulateIAM(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.PasswordComplexityPolicy{
					ObjectRoot: es_models.ObjectRoot{Sequence: 0},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "existing iam not found",
			args: args{
				es:  GetMockManipulateIAMNotExisting(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.PasswordComplexityPolicy{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					MinLength:  10,
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.AddPasswordComplexityPolicy(tt.args.ctx, tt.args.policy)
			if (tt.res.wantErr && !tt.res.errFunc(err)) || (err != nil && !tt.res.wantErr) {
				t.Errorf("got wrong err: %v ", err)
				return
			}
			if tt.res.wantErr && tt.res.errFunc(err) {
				return
			}
			if result.MinLength != tt.res.result.MinLength {
				t.Errorf("got wrong result MinLength: expected: %v, actual: %v ", tt.res.result.MinLength, result.MinLength)
			}
		})
	}
}

func TestChangePasswordComplexityPolicy(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es     *IAMEventstore
		ctx    context.Context
		policy *iam_model.PasswordComplexityPolicy
	}
	type res struct {
		result  *iam_model.PasswordComplexityPolicy
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "change password complexity policy, ok",
			args: args{
				es:  GetMockManipulateIAMWithPasswodComplexityPolicy(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.PasswordComplexityPolicy{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					MinLength:  5,
				},
			},
			res: res{
				result: &iam_model.PasswordComplexityPolicy{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					MinLength:  5,
				},
			},
		},
		{
			name: "invalid policy",
			args: args{
				es:  GetMockManipulateIAM(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.PasswordComplexityPolicy{
					ObjectRoot: es_models.ObjectRoot{Sequence: 0},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "existing iam not found",
			args: args{
				es:  GetMockManipulateIAMNotExisting(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.PasswordComplexityPolicy{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					MinLength:  10,
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.ChangePasswordComplexityPolicy(tt.args.ctx, tt.args.policy)
			if (tt.res.wantErr && !tt.res.errFunc(err)) || (err != nil && !tt.res.wantErr) {
				t.Errorf("got wrong err: %v ", err)
				return
			}
			if tt.res.wantErr && tt.res.errFunc(err) {
				return
			}
			if result.MinLength != tt.res.result.MinLength {
				t.Errorf("got wrong result MinLength: expected: %v, actual: %v ", tt.res.result.MinLength, result.MinLength)
			}
		})
	}
}

func TestAddPasswordAgePolicy(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es     *IAMEventstore
		ctx    context.Context
		policy *iam_model.PasswordAgePolicy
	}
	type res struct {
		result  *iam_model.PasswordAgePolicy
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "add password age policy, ok",
			args: args{
				es:  GetMockManipulateIAM(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.PasswordAgePolicy{
					ObjectRoot:     es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					MaxAgeDays:     10,
					ExpireWarnDays: 10,
				},
			},
			res: res{
				result: &iam_model.PasswordAgePolicy{
					ObjectRoot:     es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					MaxAgeDays:     10,
					ExpireWarnDays: 10,
				},
			},
		},
		{
			name: "empty policy",
			args: args{
				es:     GetMockManipulateIAM(ctrl),
				ctx:    authz.NewMockContext("orgID", "userID"),
				policy: nil,
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing iam not found",
			args: args{
				es:  GetMockManipulateIAMNotExisting(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.PasswordAgePolicy{
					ObjectRoot:     es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					MaxAgeDays:     10,
					ExpireWarnDays: 10,
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.AddPasswordAgePolicy(tt.args.ctx, tt.args.policy)
			if (tt.res.wantErr && !tt.res.errFunc(err)) || (err != nil && !tt.res.wantErr) {
				t.Errorf("got wrong err: %v ", err)
				return
			}
			if tt.res.wantErr && tt.res.errFunc(err) {
				return
			}
			if result.MaxAgeDays != tt.res.result.MaxAgeDays {
				t.Errorf("got wrong result MaxAgeDays: expected: %v, actual: %v ", tt.res.result.MaxAgeDays, result.MaxAgeDays)
			}

			if result.ExpireWarnDays != tt.res.result.ExpireWarnDays {
				t.Errorf("got wrong result.ExpireWarnDays: expected: %v, actual: %v ", tt.res.result.ExpireWarnDays, result.ExpireWarnDays)
			}
		})
	}
}

func TestChangePasswordAgePolicy(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es     *IAMEventstore
		ctx    context.Context
		policy *iam_model.PasswordAgePolicy
	}
	type res struct {
		result  *iam_model.PasswordAgePolicy
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "change password age policy, ok",
			args: args{
				es:  GetMockManipulateIAMWithPasswordAgePolicy(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.PasswordAgePolicy{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					MaxAgeDays: 5,
				},
			},
			res: res{
				result: &iam_model.PasswordAgePolicy{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					MaxAgeDays: 5,
				},
			},
		},
		{
			name: "invalid policy",
			args: args{
				es:     GetMockManipulateIAM(ctrl),
				ctx:    authz.NewMockContext("orgID", "userID"),
				policy: nil,
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing iam not found",
			args: args{
				es:  GetMockManipulateIAMNotExisting(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.PasswordAgePolicy{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					MaxAgeDays: 10,
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.ChangePasswordAgePolicy(tt.args.ctx, tt.args.policy)
			if (tt.res.wantErr && !tt.res.errFunc(err)) || (err != nil && !tt.res.wantErr) {
				t.Errorf("got wrong err: %v ", err)
				return
			}
			if tt.res.wantErr && tt.res.errFunc(err) {
				return
			}
			if result.MaxAgeDays != tt.res.result.MaxAgeDays {
				t.Errorf("got wrong result MaxAgeDays: expected: %v, actual: %v ", tt.res.result.MaxAgeDays, result.MaxAgeDays)
			}

			if result.ExpireWarnDays != tt.res.result.ExpireWarnDays {
				t.Errorf("got wrong result.ExpireWarnDays: expected: %v, actual: %v ", tt.res.result.ExpireWarnDays, result.ExpireWarnDays)
			}
		})
	}
}

func TestAddPasswordLockoutPolicy(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es     *IAMEventstore
		ctx    context.Context
		policy *iam_model.PasswordLockoutPolicy
	}
	type res struct {
		result  *iam_model.PasswordLockoutPolicy
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "add password lockout policy, ok",
			args: args{
				es:  GetMockManipulateIAM(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.PasswordLockoutPolicy{
					ObjectRoot:          es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					MaxAttempts:         10,
					ShowLockOutFailures: true,
				},
			},
			res: res{
				result: &iam_model.PasswordLockoutPolicy{
					ObjectRoot:          es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					MaxAttempts:         10,
					ShowLockOutFailures: true,
				},
			},
		},
		{
			name: "empty policy",
			args: args{
				es:     GetMockManipulateIAM(ctrl),
				ctx:    authz.NewMockContext("orgID", "userID"),
				policy: nil,
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing iam not found",
			args: args{
				es:  GetMockManipulateIAMNotExisting(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.PasswordLockoutPolicy{
					ObjectRoot:          es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					MaxAttempts:         10,
					ShowLockOutFailures: true,
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.AddPasswordLockoutPolicy(tt.args.ctx, tt.args.policy)
			if (tt.res.wantErr && !tt.res.errFunc(err)) || (err != nil && !tt.res.wantErr) {
				t.Errorf("got wrong err: %v ", err)
				return
			}
			if tt.res.wantErr && tt.res.errFunc(err) {
				return
			}

			if result.MaxAttempts != tt.res.result.MaxAttempts {
				t.Errorf("got wrong result MaxAttempts: expected: %v, actual: %v ", tt.res.result.MaxAttempts, result.MaxAttempts)
			}

			if result.ShowLockOutFailures != tt.res.result.ShowLockOutFailures {
				t.Errorf("got wrong result.ShowLockOutFailures: expected: %v, actual: %v ", tt.res.result.ShowLockOutFailures, result.ShowLockOutFailures)
			}
		})
	}
}

func TestChangePasswordLockoutPolicy(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es     *IAMEventstore
		ctx    context.Context
		policy *iam_model.PasswordLockoutPolicy
	}
	type res struct {
		result  *iam_model.PasswordLockoutPolicy
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "change password lockout policy, ok",
			args: args{
				es:  GetMockManipulateIAMWithPasswordLockoutPolicy(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.PasswordLockoutPolicy{
					ObjectRoot:  es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					MaxAttempts: 5,
				},
			},
			res: res{
				result: &iam_model.PasswordLockoutPolicy{
					ObjectRoot:  es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					MaxAttempts: 5,
				},
			},
		},
		{
			name: "invalid policy",
			args: args{
				es:     GetMockManipulateIAM(ctrl),
				ctx:    authz.NewMockContext("orgID", "userID"),
				policy: nil,
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing iam not found",
			args: args{
				es:  GetMockManipulateIAMNotExisting(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.PasswordLockoutPolicy{
					ObjectRoot:  es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					MaxAttempts: 10,
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.ChangePasswordLockoutPolicy(tt.args.ctx, tt.args.policy)

			if (tt.res.wantErr && !tt.res.errFunc(err)) || (err != nil && !tt.res.wantErr) {
				t.Errorf("got wrong err: %v ", err)
				return
			}
			if tt.res.wantErr && tt.res.errFunc(err) {
				return
			}
			if result.MaxAttempts != tt.res.result.MaxAttempts {
				t.Errorf("got wrong result MaxAttempts: expected: %v, actual: %v ", tt.res.result.MaxAttempts, result.MaxAttempts)
			}

			if result.ShowLockOutFailures != tt.res.result.ShowLockOutFailures {
				t.Errorf("got wrong result.ShowLockOutFailures: expected: %v, actual: %v ", tt.res.result.ShowLockOutFailures, result.ShowLockOutFailures)
			}
		})
	}
}

func TestAddOrgIAMPolicy(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es     *IAMEventstore
		ctx    context.Context
		policy *iam_model.OrgIAMPolicy
	}
	type res struct {
		result  *iam_model.OrgIAMPolicy
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "add org iam policy, ok",
			args: args{
				es:  GetMockManipulateIAM(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.OrgIAMPolicy{
					ObjectRoot:            es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					UserLoginMustBeDomain: true,
				},
			},
			res: res{
				result: &iam_model.OrgIAMPolicy{
					ObjectRoot:            es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					UserLoginMustBeDomain: true,
				},
			},
		},
		{
			name: "empty policy",
			args: args{
				es:     GetMockManipulateIAM(ctrl),
				ctx:    authz.NewMockContext("orgID", "userID"),
				policy: nil,
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing iam not found",
			args: args{
				es:  GetMockManipulateIAMNotExisting(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.OrgIAMPolicy{
					ObjectRoot:            es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					UserLoginMustBeDomain: true,
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.AddOrgIAMPolicy(tt.args.ctx, tt.args.policy)

			if (tt.res.wantErr && !tt.res.errFunc(err)) || (err != nil && !tt.res.wantErr) {
				t.Errorf("got wrong err: %v ", err)
				return
			}
			if tt.res.wantErr && tt.res.errFunc(err) {
				return
			}
			if result.UserLoginMustBeDomain != tt.res.result.UserLoginMustBeDomain {
				t.Errorf("got wrong result userLoginMustBeDomain: expected: %v, actual: %v ", tt.res.result.UserLoginMustBeDomain, result.UserLoginMustBeDomain)
			}
		})
	}
}

func TestChangeOrgIAMPolicy(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es     *IAMEventstore
		ctx    context.Context
		policy *iam_model.OrgIAMPolicy
	}
	type res struct {
		result  *iam_model.OrgIAMPolicy
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "change org iam policy, ok",
			args: args{
				es:  GetMockManipulateIAMWithOrgIAMPolicy(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.OrgIAMPolicy{
					ObjectRoot:            es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					UserLoginMustBeDomain: false,
				},
			},
			res: res{
				result: &iam_model.OrgIAMPolicy{
					ObjectRoot:            es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					UserLoginMustBeDomain: false,
				},
			},
		},
		{
			name: "invalid policy",
			args: args{
				es:     GetMockManipulateIAM(ctrl),
				ctx:    authz.NewMockContext("orgID", "userID"),
				policy: nil,
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing iam not found",
			args: args{
				es:  GetMockManipulateIAMNotExisting(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.OrgIAMPolicy{
					ObjectRoot:            es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					UserLoginMustBeDomain: true,
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.ChangeOrgIAMPolicy(tt.args.ctx, tt.args.policy)
			if (tt.res.wantErr && !tt.res.errFunc(err)) || (err != nil && !tt.res.wantErr) {
				t.Errorf("got wrong err: %v ", err)
				return
			}
			if tt.res.wantErr && tt.res.errFunc(err) {
				return
			}
			if result.UserLoginMustBeDomain != tt.res.result.UserLoginMustBeDomain {
				t.Errorf("got wrong result userLoginMustBeDomain: expected: %v, actual: %v ", tt.res.result.UserLoginMustBeDomain, result.UserLoginMustBeDomain)
			}
		})
	}
}
func TestAddMailTemplate(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es     *IAMEventstore
		ctx    context.Context
		policy *iam_model.MailTemplate
	}
	type res struct {
		result  *iam_model.MailTemplate
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "add mailtemplate, ok",
			args: args{
				es:  GetMockManipulateIAM(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.MailTemplate{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					Template:   []byte("<!doctype html>"),
				},
			},
			res: res{
				result: &iam_model.MailTemplate{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					Template:   []byte("<!doctype html>"),
				},
			},
		},
		{
			name: "invalid policy",
			args: args{
				es:  GetMockManipulateIAM(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.MailTemplate{
					ObjectRoot: es_models.ObjectRoot{Sequence: 0},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing iam not found",
			args: args{
				es:  GetMockManipulateIAMNotExisting(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.MailTemplate{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.AddMailTemplate(tt.args.ctx, tt.args.policy)
			if (tt.res.wantErr && !tt.res.errFunc(err)) || (err != nil && !tt.res.wantErr) {
				t.Errorf("got wrong err: %v ", err)
				return
			}
			if tt.res.wantErr && tt.res.errFunc(err) {
				return
			}
			if string(result.Template) != string(tt.res.result.Template) {
				t.Errorf("got wrong result Template: expected: %v, actual: %v ", tt.res.result.Template, result.Template)
			}
		})
	}
}

func TestChangeMailTemplate(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es       *IAMEventstore
		ctx      context.Context
		template *iam_model.MailTemplate
	}
	type res struct {
		result  *iam_model.MailTemplate
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "add mail template, ok",
			args: args{
				es:  GetMockManipulateIAMWithMailTemplate(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				template: &iam_model.MailTemplate{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					Template:   []byte("<!doctype html>"),
				},
			},
			res: res{
				result: &iam_model.MailTemplate{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					Template:   []byte("<!doctype html>"),
				},
			},
		},
		{
			name: "invalid mail template",
			args: args{
				es:  GetMockManipulateIAM(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				template: &iam_model.MailTemplate{
					ObjectRoot: es_models.ObjectRoot{Sequence: 0},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing iam not found",
			args: args{
				es:  GetMockManipulateIAMNotExisting(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				template: &iam_model.MailTemplate{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.ChangeMailTemplate(tt.args.ctx, tt.args.template)
			if (tt.res.wantErr && !tt.res.errFunc(err)) || (err != nil && !tt.res.wantErr) {
				t.Errorf("got wrong err: %v ", err)
				return
			}
			if tt.res.wantErr && tt.res.errFunc(err) {
				return
			}
			if string(result.Template) != string(tt.res.result.Template) {
				t.Errorf("got wrong result Template: expected: %v, actual: %v ", tt.res.result.Template, result.Template)
			}
		})
	}
}
func TestAddMailText(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es     *IAMEventstore
		ctx    context.Context
		policy *iam_model.MailText
	}
	type res struct {
		result  *iam_model.MailText
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "add mailtemplate, ok",
			args: args{
				es:  GetMockManipulateIAM(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.MailText{
					ObjectRoot:   es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					MailTextType: "Type", Language: "DE",
				},
			},
			res: res{
				result: &iam_model.MailText{
					ObjectRoot:   es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					MailTextType: "Type", Language: "DE",
				},
			},
		},
		{
			name: "invalid policy",
			args: args{
				es:  GetMockManipulateIAM(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.MailText{
					ObjectRoot: es_models.ObjectRoot{Sequence: 0},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing iam not found",
			args: args{
				es:  GetMockManipulateIAMNotExisting(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.MailText{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.AddMailText(tt.args.ctx, tt.args.policy)
			if (tt.res.wantErr && !tt.res.errFunc(err)) || (err != nil && !tt.res.wantErr) {
				t.Errorf("got wrong err: %v ", err)
				return
			}
			if tt.res.wantErr && tt.res.errFunc(err) {
				return
			}
			if string(result.MailTextType) != string(tt.res.result.MailTextType) {
				t.Errorf("got wrong result MailTextType: expected: %v, actual: %v ", tt.res.result.MailTextType, result.MailTextType)
			}
		})
	}
}

func TestChangeMailText(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es     *IAMEventstore
		ctx    context.Context
		policy *iam_model.MailText
	}
	type res struct {
		result  *iam_model.MailText
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "change mailtemplate, ok",
			args: args{
				es:  GetMockManipulateIAMWithMailText(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.MailText{
					ObjectRoot:   es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					MailTextType: "Type", Language: "DE",
				},
			},
			res: res{
				result: &iam_model.MailText{
					ObjectRoot:   es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					MailTextType: "Type", Language: "DE",
				},
			},
		},
		{
			name: "invalid policy",
			args: args{
				es:  GetMockManipulateIAM(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.MailText{
					ObjectRoot: es_models.ObjectRoot{Sequence: 0},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing iam not found",
			args: args{
				es:  GetMockManipulateIAMNotExisting(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.MailText{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.ChangeMailText(tt.args.ctx, tt.args.policy)
			if (tt.res.wantErr && !tt.res.errFunc(err)) || (err != nil && !tt.res.wantErr) {
				t.Errorf("got wrong err: %v ", err)
				return
			}
			if tt.res.wantErr && tt.res.errFunc(err) {
				return
			}
			if string(result.MailTextType) != string(tt.res.result.MailTextType) {
				t.Errorf("got wrong result MailTextType: expected: %v, actual: %v ", tt.res.result.MailTextType, result.MailTextType)
			}
		})
	}
}
