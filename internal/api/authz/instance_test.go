package authz

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/execution/target"
	"github.com/zitadel/zitadel/internal/feature"
)

func Test_Instance(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type res struct {
		instanceID string
		projectID  string
		consoleID  string
		features   feature.Features
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			"empty context",
			args{
				context.Background(),
			},
			res{
				instanceID: "",
				projectID:  "",
				consoleID:  "",
			},
		},
		{
			"WithInstanceID",
			args{
				WithInstanceID(context.Background(), "id"),
			},
			res{
				instanceID: "id",
				projectID:  "",
				consoleID:  "",
			},
		},
		{
			"WithInstance",
			args{
				WithInstance(context.Background(), &mockInstance{}),
			},
			res{
				instanceID: "instanceID",
				projectID:  "projectID",
				consoleID:  "managementConsoleID",
			},
		},
		{
			"WithFeatures",
			args{
				WithFeatures(context.Background(), feature.Features{
					LoginDefaultOrg: true,
				}),
			},
			res{
				features: feature.Features{
					LoginDefaultOrg: true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetInstance(tt.args.ctx)
			assert.Equal(t, tt.res.instanceID, got.InstanceID())
			assert.Equal(t, tt.res.projectID, got.ProjectID())
			assert.Equal(t, tt.res.consoleID, got.ManagementConsoleClientID())
			assert.Equal(t, tt.res.features, got.Features())
		})
	}
}

type mockInstance struct{}

func (m *mockInstance) Block() *bool {
	panic("shouldn't be called here")
}

func (m *mockInstance) AuditLogRetention() *time.Duration {
	panic("shouldn't be called here")
}

func (m *mockInstance) InstanceID() string {
	return "instanceID"
}

func (m *mockInstance) ProjectID() string {
	return "projectID"
}

func (m *mockInstance) ManagementConsoleClientID() string {
	return "managementConsoleID"
}

func (m *mockInstance) ManagementConsoleApplicationID() string {
	return "appID"
}

func (m *mockInstance) DefaultLanguage() language.Tag {
	return language.English
}

func (m *mockInstance) AllowedLanguages() []language.Tag {
	return []language.Tag{language.English}
}

func (m *mockInstance) DefaultOrganisationID() string {
	return "orgID"
}

func (m *mockInstance) SecurityPolicyAllowedOrigins() []string {
	return nil
}

func (m *mockInstance) EnableImpersonation() bool {
	return false
}

func (m *mockInstance) Features() feature.Features {
	return feature.Features{}
}

func (m *mockInstance) ExecutionRouter() target.Router {
	return target.NewRouter(nil)
}
