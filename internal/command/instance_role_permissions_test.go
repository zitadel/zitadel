package command

import (
	"context"
	"database/sql"
	"database/sql/driver"
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/database/mock"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/permission"
)

func Test_rolePermissionMappingsToDatabaseMap(t *testing.T) {
	type args struct {
		mappings []authz.RoleMapping
		system   bool
	}
	tests := []struct {
		name string
		args args
		want database.Map[[]string]
	}{
		{
			name: "instance",
			args: args{
				mappings: []authz.RoleMapping{
					{Role: "role1", Permissions: []string{"permission1", "permission2"}},
					{Role: "role2", Permissions: []string{"permission3", "permission4"}},
					{Role: "SYSTEM_ROLE", Permissions: []string{"permission5", "permission6"}},
				},
				system: false,
			},
			want: database.Map[[]string]{
				"role1": []string{"permission1", "permission2"},
				"role2": []string{"permission3", "permission4"},
			},
		},
		{
			name: "system",
			args: args{
				mappings: []authz.RoleMapping{
					{Role: "role1", Permissions: []string{"permission1", "permission2"}},
					{Role: "role2", Permissions: []string{"permission3", "permission4"}},
					{Role: "SYSTEM_ROLE", Permissions: []string{"permission5", "permission6"}},
				},
				system: true,
			},
			want: database.Map[[]string]{
				"SYSTEM_ROLE": []string{"permission5", "permission6"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := rolePermissionMappingsToDatabaseMap(tt.args.mappings, tt.args.system)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_synchronizeRolePermissionCommands(t *testing.T) {
	const aggregateID = "aggregateID"
	aggregate := permission.NewAggregate(aggregateID)
	target := database.Map[[]string]{
		"role1": []string{"permission1", "permission2"},
		"role2": []string{"permission3", "permission4"},
	}
	tests := []struct {
		name     string
		mock     func(*testing.T) *mock.SQLMock
		wantCmds []eventstore.Command
		wantErr  error
	}{
		{
			name: "query error",
			mock: func(t *testing.T) *mock.SQLMock {
				return mock.NewSQLMock(t,
					mock.ExpectQuery(instanceRolePermissionsSyncQuery,
						mock.WithQueryArgs(aggregateID, target),
						mock.WithQueryErr(sql.ErrConnDone),
					),
				)
			},
			wantErr: sql.ErrConnDone,
		},
		{
			name: "no rows",
			mock: func(t *testing.T) *mock.SQLMock {
				return mock.NewSQLMock(t,
					mock.ExpectQuery(instanceRolePermissionsSyncQuery,
						mock.WithQueryArgs(aggregateID, target),
						mock.WithQueryResult([]string{"operation", "role", "permission"}, [][]driver.Value{}),
					),
				)
			},
		},
		{
			name: "add and remove operations",
			mock: func(t *testing.T) *mock.SQLMock {
				return mock.NewSQLMock(t,
					mock.ExpectQuery(instanceRolePermissionsSyncQuery,
						mock.WithQueryArgs(aggregateID, target),
						mock.WithQueryResult([]string{"operation", "role", "permission"}, [][]driver.Value{
							{"add", "role1", "permission1"},
							{"add", "role1", "permission2"},
							{"remove", "role3", "permission5"},
							{"remove", "role3", "permission6"},
						}),
					),
				)
			},
			wantCmds: []eventstore.Command{
				permission.NewAddedEvent(context.Background(), aggregate, "role1", "permission1"),
				permission.NewAddedEvent(context.Background(), aggregate, "role1", "permission2"),
				permission.NewRemovedEvent(context.Background(), aggregate, "role3", "permission5"),
				permission.NewRemovedEvent(context.Background(), aggregate, "role3", "permission6"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := tt.mock(t)
			defer mock.Assert(t)
			db := &database.DB{
				DB: mock.DB,
			}
			gotCmds, err := synchronizeRolePermissionCommands(context.Background(), db, aggregateID, target)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.wantCmds, gotCmds)
		})
	}
}
