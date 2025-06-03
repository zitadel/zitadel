package migration

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/database"
)

const (
	expCountTriggerQuery = `CREATE OR REPLACE TRIGGER count_resource
    AFTER INSERT OR DELETE
    ON table
    FOR EACH ROW
    EXECUTE FUNCTION projections.count_resource(
        'instance', 
        'instance_id', 
        'parent_id',
        'resource'
    );

CREATE OR REPLACE TRIGGER truncate_resource_counts
    AFTER TRUNCATE
    ON table
    FOR EACH STATEMENT
    EXECUTE FUNCTION projections.delete_table_counts();

-- Prevent inserts and deletes while we populate the counts.
LOCK TABLE table IN SHARE MODE;

-- Populate the resource counts for the existing data in the table.
INSERT INTO projections.resource_counts(
	instance_id,
    table_name,
    parent_type,
    parent_id,
    resource_name,
    amount
)
SELECT
    instance_id,
    'table',
    'instance',
    parent_id,
    'resource',
    COUNT(*) AS amount
FROM table
GROUP BY (instance_id, parent_id)
ON CONFLICT (instance_id, table_name, parent_type, parent_id) DO
UPDATE SET updated_at = now(), amount = EXCLUDED.amount;`

	expDeleteParentCountsQuery = `CREATE OR REPLACE TRIGGER delete_parent_counts_trigger
    AFTER DELETE
    ON table
    FOR EACH ROW
    EXECUTE FUNCTION projections.delete_parent_counts(
        'instance', 
        'instance_id', 
        'parent_id'
    );`
)

func Test_triggerMigration_Execute(t *testing.T) {
	type fields struct {
		triggerConfig triggerConfig
		templateName  string
	}
	tests := []struct {
		name    string
		fields  fields
		expects func(sqlmock.Sqlmock)
		wantErr bool
	}{
		{
			name: "template error",
			fields: fields{
				triggerConfig: triggerConfig{
					Table:            "table",
					ParentType:       "instance",
					InstanceIDColumn: "instance_id",
					ParentIDColumn:   "parent_id",
					Resource:         "resource",
				},
				templateName: "foo",
			},
			expects: func(_ sqlmock.Sqlmock) {},
			wantErr: true,
		},
		{
			name: "db error",
			fields: fields{
				triggerConfig: triggerConfig{
					Table:            "table",
					ParentType:       "instance",
					InstanceIDColumn: "instance_id",
					ParentIDColumn:   "parent_id",
					Resource:         "resource",
				},
				templateName: countTriggerTmpl,
			},
			expects: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(expCountTriggerQuery)).
					WillReturnError(assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "count trigger",
			fields: fields{
				triggerConfig: triggerConfig{
					Table:            "table",
					ParentType:       "instance",
					InstanceIDColumn: "instance_id",
					ParentIDColumn:   "parent_id",
					Resource:         "resource",
				},
				templateName: countTriggerTmpl,
			},
			expects: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(expCountTriggerQuery)).
					WithoutArgs().
					WillReturnResult(
						sqlmock.NewResult(1, 1),
					)
			},
		},
		{
			name: "count trigger",
			fields: fields{
				triggerConfig: triggerConfig{
					Table:            "table",
					ParentType:       "instance",
					InstanceIDColumn: "instance_id",
					ParentIDColumn:   "parent_id",
					Resource:         "resource",
				},
				templateName: deleteParentCountsTmpl,
			},
			expects: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(expDeleteParentCountsQuery)).
					WithoutArgs().
					WillReturnResult(
						sqlmock.NewResult(1, 1),
					)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer func() {
				err := mock.ExpectationsWereMet()
				require.NoError(t, err)
			}()
			defer db.Close()
			tt.expects(mock)
			mock.ExpectClose()

			m := &triggerMigration{
				db: &database.DB{
					DB: db,
				},
				triggerConfig: tt.fields.triggerConfig,
				templateName:  tt.fields.templateName,
			}
			err = m.Execute(context.Background(), nil)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func Test_triggerConfig_Check(t *testing.T) {
	type fields struct {
		Table            string
		ParentType       string
		InstanceIDColumn string
		ParentIDColumn   string
		Resource         string
	}
	type args struct {
		lastRun map[string]any
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "should",
			fields: fields{
				Table:            "users2",
				ParentType:       "instance",
				InstanceIDColumn: "instance_id",
				ParentIDColumn:   "parent_id",
				Resource:         "user",
			},
			args: args{
				lastRun: map[string]any{
					"table":              "users1",
					"parent_type":        "instance",
					"instance_id_column": "instance_id",
					"parent_id_column":   "parent_id",
					"resource":           "user",
				},
			},
			want: true,
		},
		{
			name: "should not",
			fields: fields{
				Table:            "users1",
				ParentType:       "instance",
				InstanceIDColumn: "instance_id",
				ParentIDColumn:   "parent_id",
				Resource:         "user",
			},
			args: args{
				lastRun: map[string]any{
					"table":              "users1",
					"parent_type":        "instance",
					"instance_id_column": "instance_id",
					"parent_id_column":   "parent_id",
					"resource":           "user",
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &triggerConfig{
				Table:            tt.fields.Table,
				ParentType:       tt.fields.ParentType,
				InstanceIDColumn: tt.fields.InstanceIDColumn,
				ParentIDColumn:   tt.fields.ParentIDColumn,
				Resource:         tt.fields.Resource,
			}
			got := c.Check(tt.args.lastRun)
			assert.Equal(t, tt.want, got)
		})
	}
}
