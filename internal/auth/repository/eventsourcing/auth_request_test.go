package eventsourcing

import (
	"testing"

	"github.com/stretchr/testify/assert"

	req_model "github.com/caos/zitadel/internal/auth_request/model"
	"github.com/caos/zitadel/internal/user/model"
)

func Test_nextSteps(t *testing.T) {
	type args struct {
		request *req_model.AuthRequest
		user    *model.User
	}
	tests := []struct {
		name    string
		args    args
		want    []req_model.NextStep
		wantErr bool
	}{
		//{
		//	"no auth request, error",
		//	args{
		//		request: nil,
		//		user:    nil,
		//	},
		//	nil,
		//	true,
		//},
		{
			"no user, login step",
			args{
				request: &req_model.AuthRequest{},
				user:    nil,
			},
			[]req_model.NextStep{&req_model.LoginStep{NotFound: true}},
			false,
		},
		{
			"no password, init password step",
			args{
				request: &req_model.AuthRequest{},
				user:    &model.User{},
			},
			[]req_model.NextStep{&req_model.InitPasswordStep{}},
			false,
		},
		//{
		//	"password unverified, password step",
		//	args{
		//		request: &req_model.AuthRequest{},
		//		user:    &model.User{Password: &model.Password{}},
		//	},
		//	[]req_model.NextStep{&req_model.PasswordStep{}},
		//	false,
		//},
		//{
		//	"password wrong, password step",
		//	args{
		//		request: &req_model.AuthRequest{},
		//		user:    &model.User{Password: &model.Password{}}, //TODO: mock verified?
		//	},
		//	[]req_model.NextStep{&req_model.PasswordStep{FailureCount:1}},
		//	false,
		//},
		{
			"password change required, password change step",
			args{
				request: &req_model.AuthRequest{},
				user: &model.User{
					Password: &model.Password{
						ChangeRequired: true,
					},
				},
			},
			[]req_model.NextStep{&req_model.ChangePasswordStep{}},
			false,
		},
		{
			"no email, verify email step", //TODO: ??
			args{
				request: &req_model.AuthRequest{},
				user: &model.User{
					Password: &model.Password{},
				},
			},
			[]req_model.NextStep{&req_model.VerifyEMailStep{}},
			false,
		},
		{
			"email unverified, verify email step",
			args{
				request: &req_model.AuthRequest{},
				user: &model.User{
					Password: &model.Password{},
				},
			},
			[]req_model.NextStep{&req_model.VerifyEMailStep{}},
			false,
		},
		{
			"mfa not set up, mfa prompt step",
			args{
				request: &req_model.AuthRequest{},
				user: &model.User{
					Password: &model.Password{},
					Email: &model.Email{
						IsEmailVerified: true,
					},
				},
			},
			[]req_model.NextStep{&req_model.MfaVerificationStep{
				MfaProviders: []model.MfaType{&model.OTP{}},
			}},
			false,
		},

		{
			"mfa set up and not verified, mfa verification step",
			args{
				request: &req_model.AuthRequest{},
				user: &model.User{
					Password: &model.Password{},
					OTP: &model.OTP{
						State: model.MFASTATE_READY,
					},
				}, //TODO: mock verified?
			},
			[]req_model.NextStep{&req_model.MfaVerificationStep{
				MfaProviders: []model.MfaType{&model.OTP{}},
			}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := nextSteps(tt.args.request, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("nextSteps() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.ElementsMatch(t, got.PossibleSteps, tt.want)
		})
	}
}
