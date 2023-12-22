package authz

import (
	"testing"

	"github.com/zitadel/zitadel/internal/zerrors"
)

type TestRequest struct {
	Test string
}

func Test_CheckUserPermissions(t *testing.T) {
	type args struct {
		req     *TestRequest
		perms   []string
		authOpt Option
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "no permissions",
			args: args{
				req:   &TestRequest{},
				perms: []string{},
			},
			wantErr: true,
		},
		{
			name: "has permission and no context requested",
			args: args{
				req:     &TestRequest{},
				perms:   []string{"project.read"},
				authOpt: Option{CheckParam: ""},
			},
			wantErr: false,
		},
		{
			name: "context requested and has global permission",
			args: args{
				req:     &TestRequest{Test: "Test"},
				perms:   []string{"project.read", "project.read:1"},
				authOpt: Option{CheckParam: "Test"},
			},
			wantErr: false,
		},
		{
			name: "context requested and has specific permission",
			args: args{
				req:     &TestRequest{Test: "Test"},
				perms:   []string{"project.read:Test"},
				authOpt: Option{CheckParam: "Test"},
			},
			wantErr: false,
		},
		{
			name: "context requested and has no permission",
			args: args{
				req:     &TestRequest{Test: "Hodor"},
				perms:   []string{"project.read:Test"},
				authOpt: Option{CheckParam: "Test"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkUserPermissions(tt.args.req, tt.args.perms, tt.args.authOpt)
			if tt.wantErr && err == nil {
				t.Errorf("got wrong result, should get err: actual: %v ", err)
			}

			if !tt.wantErr && err != nil {
				t.Errorf("shouldn't get err: %v ", err)
			}

			if tt.wantErr && !zerrors.IsPermissionDenied(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func Test_SplitPermission(t *testing.T) {
	type args struct {
		perm string
	}
	tests := []struct {
		name      string
		args      args
		permName  string
		permCtxID string
	}{
		{
			name: "permission with context id",
			args: args{
				perm: "project.read:ctxID",
			},
			permName:  "project.read",
			permCtxID: "ctxID",
		},
		{
			name: "permission without context id",
			args: args{
				perm: "project.read",
			},
			permName:  "project.read",
			permCtxID: "",
		},
		{
			name: "permission to many parts",
			args: args{
				perm: "project.read:1:0",
			},
			permName:  "project.read",
			permCtxID: "1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name, id := SplitPermission(tt.args.perm)
			if name != tt.permName {
				t.Errorf("got wrong result on name, expecting: %v, actual: %v ", tt.permName, name)
			}
			if id != tt.permCtxID {
				t.Errorf("got wrong result on id, expecting: %v, actual: %v ", tt.permCtxID, id)
			}
		})
	}
}

func Test_HasContextPermission(t *testing.T) {
	type args struct {
		req       *TestRequest
		fieldname string
		perms     []string
	}
	tests := []struct {
		name   string
		args   args
		result bool
	}{
		{
			name: "existing context permission",
			args: args{
				req:       &TestRequest{Test: "right"},
				fieldname: "Test",
				perms:     []string{"test:wrong", "test:right"},
			},
			result: true,
		},
		{
			name: "not existing context permission",
			args: args{
				req:       &TestRequest{Test: "test"},
				fieldname: "Test",
				perms:     []string{"test:wrong", "test:wrong2"},
			},
			result: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasContextPermission(tt.args.req, tt.args.fieldname, tt.args.perms)
			if result != tt.result {
				t.Errorf("got wrong result, expecting: %v, actual: %v ", tt.result, result)
			}
		})
	}
}

func Test_GetFieldFromReq(t *testing.T) {
	type args struct {
		req       *TestRequest
		fieldname string
	}
	tests := []struct {
		name   string
		args   args
		result string
	}{
		{
			name: "existing field",
			args: args{
				req:       &TestRequest{Test: "TestValue"},
				fieldname: "Test",
			},
			result: "TestValue",
		},
		{
			name: "not existing field",
			args: args{
				req:       &TestRequest{Test: "TestValue"},
				fieldname: "Test2",
			},
			result: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getFieldFromReq(tt.args.req, tt.args.fieldname)
			if result != tt.result {
				t.Errorf("got wrong result, expecting: %v, actual: %v ", tt.result, result)
			}
		})
	}
}

func Test_HasGlobalPermission(t *testing.T) {
	type args struct {
		perms []string
	}
	tests := []struct {
		name   string
		args   args
		result bool
	}{
		{
			name: "global perm existing",
			args: args{
				perms: []string{"perm:1", "perm:2", "perm"},
			},
			result: true,
		},
		{
			name: "global perm not existing",
			args: args{
				perms: []string{"perm:1", "perm:2", "perm:3"},
			},
			result: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HasGlobalPermission(tt.args.perms)
			if result != tt.result {
				t.Errorf("got wrong result, expecting: %v, actual: %v ", tt.result, result)
			}
		})
	}
}

func Test_GetPermissionCtxIDs(t *testing.T) {
	type args struct {
		perms []string
	}
	tests := []struct {
		name   string
		args   args
		result []string
	}{
		{
			name: "no specific permission",
			args: args{
				perms: []string{"perm"},
			},
			result: []string{},
		},
		{
			name: "ctx id",
			args: args{
				perms: []string{"perm:1", "perm", "perm:3"},
			},
			result: []string{"1", "3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetAllPermissionCtxIDs(tt.args.perms)
			if !equalStringArray(result, tt.result) {
				t.Errorf("got wrong result, expecting: %v, actual: %v ", tt.result, result)
			}
		})
	}
}
