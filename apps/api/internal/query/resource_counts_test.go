package query

import (
	"context"
	_ "embed"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
)

func TestQueries_ListResourceCounts(t *testing.T) {
	columns := []string{"id", "instance_id", "table_name", "parent_type", "parent_id", "resource_name", "updated_at", "amount"}
	type args struct {
		lastID int
		limit  int
	}
	tests := []struct {
		name       string
		args       args
		expects    func(sqlmock.Sqlmock)
		wantResult []ResourceCount
		wantErr    bool
	}{
		{
			name: "query error",
			args: args{
				lastID: 0,
				limit:  10,
			},
			expects: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(resourceCountsListQuery)).
					WithArgs(0, 10).
					WillReturnError(assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "success",
			args: args{
				lastID: 0,
				limit:  10,
			},
			expects: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(resourceCountsListQuery)).
					WithArgs(0, 10).
					WillReturnRows(
						sqlmock.NewRows(columns).
							AddRow(1, "instance_1", "table", "instance", "parent_1", "resource_name", time.Unix(1, 2), 5).
							AddRow(2, "instance_2", "table", "instance", "parent_2", "resource_name", time.Unix(1, 2), 6),
					)
			},
			wantResult: []ResourceCount{
				{
					ID:         1,
					InstanceID: "instance_1",
					TableName:  "table",
					ParentType: domain.CountParentTypeInstance,
					ParentID:   "parent_1",
					Resource:   "resource_name",
					UpdatedAt:  time.Unix(1, 2),
					Amount:     5,
				},
				{
					ID:         2,
					InstanceID: "instance_2",
					TableName:  "table",
					ParentType: domain.CountParentTypeInstance,
					ParentID:   "parent_2",
					Resource:   "resource_name",
					UpdatedAt:  time.Unix(1, 2),
					Amount:     6,
				},
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
			q := &Queries{
				client: &database.DB{
					DB: db,
				},
			}

			gotResult, err := q.ListResourceCounts(context.Background(), tt.args.lastID, tt.args.limit)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantResult, gotResult, "ListResourceCounts() result mismatch")
		})
	}
}
