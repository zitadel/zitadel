package query

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/domain"
)

var (
	prepareFlowStmt = `SELECT projections.actions3.id,` +
		` projections.actions3.creation_date,` +
		` projections.actions3.change_date,` +
		` projections.actions3.resource_owner,` +
		` projections.actions3.action_state,` +
		` projections.actions3.sequence,` +
		` projections.actions3.name,` +
		` projections.actions3.script,` +
		` projections.actions3.allowed_to_fail,` +
		` projections.actions3.timeout,` +
		` projections.flow_triggers2.trigger_type,` +
		` projections.flow_triggers2.trigger_sequence,` +
		` projections.flow_triggers2.flow_type,` +
		` projections.flow_triggers2.change_date,` +
		` projections.flow_triggers2.sequence,` +
		` projections.flow_triggers2.resource_owner` +
		` FROM projections.flow_triggers2` +
		` LEFT JOIN projections.actions3 ON projections.flow_triggers2.action_id = projections.actions3.id AND projections.flow_triggers2.instance_id = projections.actions3.instance_id` +
		` AS OF SYSTEM TIME '-1 ms'`
	prepareFlowCols = []string{
		"id",
		"creation_date",
		"change_date",
		"resource_owner",
		"state",
		"sequence",
		"name",
		"script",
		"allowed_to_fail",
		"timeout",
		// flow
		"trigger_type",
		"trigger_sequence",
		"flow_type",
		"change_date",
		"sequence",
		"resource_owner",
	}

	prepareTriggerActionStmt = `SELECT projections.actions3.id,` +
		` projections.actions3.creation_date,` +
		` projections.actions3.change_date,` +
		` projections.actions3.resource_owner,` +
		` projections.actions3.action_state,` +
		` projections.actions3.sequence,` +
		` projections.actions3.name,` +
		` projections.actions3.script,` +
		` projections.actions3.allowed_to_fail,` +
		` projections.actions3.timeout` +
		` FROM projections.flow_triggers2` +
		` LEFT JOIN projections.actions3 ON projections.flow_triggers2.action_id = projections.actions3.id AND projections.flow_triggers2.instance_id = projections.actions3.instance_id` +
		` AS OF SYSTEM TIME '-1 ms'`

	prepareTriggerActionCols = []string{
		"id",
		"creation_date",
		"change_date",
		"resource_owner",
		"state",
		"sequence",
		"name",
		"script",
		"allowed_to_fail",
		"timeout",
	}

	prepareFlowTypeStmt = `SELECT projections.flow_triggers2.flow_type` +
		` FROM projections.flow_triggers2` +
		` AS OF SYSTEM TIME '-1 ms'`

	prepareFlowTypeCols = []string{
		"flow_type",
	}
)

func Test_FlowPrepares(t *testing.T) {
	type want struct {
		sqlExpectations sqlExpectation
		err             checkErr
	}
	tests := []struct {
		name    string
		prepare interface{}
		want    want
		object  interface{}
	}{
		{
			name: "prepareFlowQuery no result",
			prepare: func(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Rows) (*Flow, error)) {
				return prepareFlowQuery(ctx, db, domain.FlowTypeExternalAuthentication)
			},
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(prepareFlowStmt),
					nil,
					nil,
				),
			},
			object: &Flow{
				TriggerActions: map[domain.TriggerType][]*Action{},
				Type:           domain.FlowTypeExternalAuthentication,
			},
		},
		{
			name: "prepareFlowQuery one action",
			prepare: func(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Rows) (*Flow, error)) {
				return prepareFlowQuery(ctx, db, domain.FlowTypeExternalAuthentication)
			},
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(prepareFlowStmt),
					prepareFlowCols,
					[][]driver.Value{
						{
							"action-id",
							testNow,
							testNow,
							"ro",
							domain.ActionStateActive,
							uint64(20211115),
							"action-name",
							"script",
							true,
							10000000000,
							domain.TriggerTypePreCreation,
							uint64(20211109),
							domain.FlowTypeExternalAuthentication,
							testNow,
							uint64(20211115),
							"owner",
						},
					},
				),
			},
			object: &Flow{
				ChangeDate:    testNow,
				ResourceOwner: "owner",
				Sequence:      20211115,
				Type:          domain.FlowTypeExternalAuthentication,
				TriggerActions: map[domain.TriggerType][]*Action{
					domain.TriggerTypePreCreation: {
						{
							ID:            "action-id",
							CreationDate:  testNow,
							ChangeDate:    testNow,
							ResourceOwner: "ro",
							State:         domain.ActionStateActive,
							Sequence:      20211115,
							Name:          "action-name",
							Script:        "script",
							AllowedToFail: true,
							timeout:       10 * time.Second,
						},
					},
				},
			},
		},
		{
			name: "prepareFlowQuery multiple actions",
			prepare: func(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Rows) (*Flow, error)) {
				return prepareFlowQuery(ctx, db, domain.FlowTypeExternalAuthentication)
			},
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(prepareFlowStmt),
					prepareFlowCols,
					[][]driver.Value{
						{
							"action-id-pre",
							testNow,
							testNow,
							"ro",
							domain.ActionStateActive,
							uint64(20211115),
							"action-name-pre",
							"script",
							true,
							10000000000,
							domain.TriggerTypePreCreation,
							uint64(20211109),
							domain.FlowTypeExternalAuthentication,
							testNow,
							uint64(20211115),
							"owner",
						},
						{
							"action-id-post",
							testNow,
							testNow,
							"ro",
							domain.ActionStateActive,
							uint64(20211115),
							"action-name-post",
							"script",
							false,
							5000000000,
							domain.TriggerTypePostCreation,
							uint64(20211109),
							domain.FlowTypeExternalAuthentication,
							testNow,
							uint64(20211115),
							"owner",
						},
					},
				),
			},
			object: &Flow{
				ChangeDate:    testNow,
				ResourceOwner: "owner",
				Sequence:      20211115,
				Type:          domain.FlowTypeExternalAuthentication,
				TriggerActions: map[domain.TriggerType][]*Action{
					domain.TriggerTypePreCreation: {
						{
							ID:            "action-id-pre",
							CreationDate:  testNow,
							ChangeDate:    testNow,
							ResourceOwner: "ro",
							State:         domain.ActionStateActive,
							Sequence:      20211115,
							Name:          "action-name-pre",
							Script:        "script",
							AllowedToFail: true,
							timeout:       10 * time.Second,
						},
					},
					domain.TriggerTypePostCreation: {
						{
							ID:            "action-id-post",
							CreationDate:  testNow,
							ChangeDate:    testNow,
							ResourceOwner: "ro",
							State:         domain.ActionStateActive,
							Sequence:      20211115,
							Name:          "action-name-post",
							Script:        "script",
							AllowedToFail: false,
							timeout:       5 * time.Second,
						},
					},
				},
			},
		},
		{
			name: "prepareFlowQuery no action",
			prepare: func(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Rows) (*Flow, error)) {
				return prepareFlowQuery(ctx, db, domain.FlowTypeExternalAuthentication)
			},
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(prepareFlowStmt),
					prepareFlowCols,
					[][]driver.Value{
						{
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							domain.TriggerTypePostCreation,
							uint64(20211109),
							domain.FlowTypeExternalAuthentication,
							testNow,
							uint64(20211115),
							"owner",
						},
					},
				),
			},
			object: &Flow{
				ChangeDate:     testNow,
				ResourceOwner:  "owner",
				Sequence:       20211115,
				Type:           domain.FlowTypeExternalAuthentication,
				TriggerActions: map[domain.TriggerType][]*Action{},
			},
		},
		{
			name: "prepareFlowQuery sql err",
			prepare: func(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Rows) (*Flow, error)) {
				return prepareFlowQuery(ctx, db, domain.FlowTypeExternalAuthentication)
			},
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(prepareFlowStmt),
					sql.ErrConnDone,
				),
				err: func(err error) (error, bool) {
					if !errors.Is(err, sql.ErrConnDone) {
						return fmt.Errorf("err should be sql.ErrConnDone got: %w", err), false
					}
					return nil, true
				},
			},
			object: nil,
		},
		{
			name:    "prepareTriggerActionsQuery no result",
			prepare: prepareTriggerActionsQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(prepareTriggerActionStmt),
					nil,
					nil,
				),
			},
			object: []*Action{},
		},
		{
			name:    "prepareTriggerActionsQuery one result",
			prepare: prepareTriggerActionsQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(prepareTriggerActionStmt),
					prepareTriggerActionCols,
					[][]driver.Value{
						{
							"action-id",
							testNow,
							testNow,
							"ro",
							domain.AddressStateActive,
							uint64(20211115),
							"action-name",
							"script",
							true,
							10000000000,
						},
					},
				),
			},
			object: []*Action{
				{
					ID:            "action-id",
					CreationDate:  testNow,
					ChangeDate:    testNow,
					ResourceOwner: "ro",
					State:         domain.ActionStateActive,
					Sequence:      20211115,
					Name:          "action-name",
					Script:        "script",
					AllowedToFail: true,
					timeout:       10 * time.Second,
				},
			},
		},
		{
			name:    "prepareTriggerActionsQuery multiple results",
			prepare: prepareTriggerActionsQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(prepareTriggerActionStmt),
					prepareTriggerActionCols,
					[][]driver.Value{
						{
							"action-id-1",
							testNow,
							testNow,
							"ro",
							domain.AddressStateActive,
							uint64(20211115),
							"action-name-1",
							"script",
							true,
							10000000000,
						},
						{
							"action-id-2",
							testNow,
							testNow,
							"ro",
							domain.ActionStateActive,
							uint64(20211115),
							"action-name-2",
							"script",
							false,
							5000000000,
						},
					},
				),
			},
			object: []*Action{
				{
					ID:            "action-id-1",
					CreationDate:  testNow,
					ChangeDate:    testNow,
					ResourceOwner: "ro",
					State:         domain.ActionStateActive,
					Sequence:      20211115,
					Name:          "action-name-1",
					Script:        "script",
					AllowedToFail: true,
					timeout:       10 * time.Second,
				},
				{
					ID:            "action-id-2",
					CreationDate:  testNow,
					ChangeDate:    testNow,
					ResourceOwner: "ro",
					State:         domain.ActionStateActive,
					Sequence:      20211115,
					Name:          "action-name-2",
					Script:        "script",
					AllowedToFail: false,
					timeout:       5 * time.Second,
				},
			},
		},
		{
			name:    "prepareTriggerActionsQuery sql err",
			prepare: prepareTriggerActionsQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(prepareTriggerActionStmt),
					sql.ErrConnDone,
				),
				err: func(err error) (error, bool) {
					if !errors.Is(err, sql.ErrConnDone) {
						return fmt.Errorf("err should be sql.ErrConnDone got: %w", err), false
					}
					return nil, true
				},
			},
			object: nil,
		},
		{
			name:    "prepareFlowTypesQuery no result",
			prepare: prepareFlowTypesQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(prepareFlowTypeStmt),
					nil,
					nil,
				),
			},
			object: []domain.FlowType{},
		},
		{
			name:    "prepareFlowTypesQuery one result",
			prepare: prepareFlowTypesQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(prepareFlowTypeStmt),
					prepareFlowTypeCols,
					[][]driver.Value{
						{
							domain.FlowTypeExternalAuthentication,
						},
					},
				),
			},
			object: []domain.FlowType{
				domain.FlowTypeExternalAuthentication,
			},
		},
		{
			name:    "prepareFlowTypesQuery multiple results",
			prepare: prepareFlowTypesQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(prepareFlowTypeStmt),
					prepareFlowTypeCols,
					[][]driver.Value{
						{
							domain.FlowTypeExternalAuthentication,
						},
						{
							domain.FlowTypeUnspecified,
						},
					},
				),
			},
			object: []domain.FlowType{
				domain.FlowTypeExternalAuthentication,
				domain.FlowTypeUnspecified,
			},
		},
		{
			name:    "prepareFlowTypesQuery sql err",
			prepare: prepareFlowTypesQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(prepareFlowTypeStmt),
					sql.ErrConnDone,
				),
				err: func(err error) (error, bool) {
					if !errors.Is(err, sql.ErrConnDone) {
						return fmt.Errorf("err should be sql.ErrConnDone got: %w", err), false
					}
					return nil, true
				},
			},
			object: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err, defaultPrepareArgs...)
		})
	}
}
