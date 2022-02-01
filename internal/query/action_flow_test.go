package query

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/caos/zitadel/internal/domain"
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
			name:    "prepareFlowQuery no result",
			prepare: prepareFlowQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(`SELECT zitadel.projections.actions.id,`+
						` zitadel.projections.actions.creation_date,`+
						` zitadel.projections.actions.change_date,`+
						` zitadel.projections.actions.resource_owner,`+
						` zitadel.projections.actions.action_state,`+
						` zitadel.projections.actions.sequence,`+
						` zitadel.projections.actions.name,`+
						` zitadel.projections.actions.script,`+
						` zitadel.projections.flows_triggers.trigger_type,`+
						` zitadel.projections.flows_triggers.trigger_sequence,`+
						` zitadel.projections.flows_triggers.flow_type`+
						` FROM zitadel.projections.flows_triggers`+
						` LEFT JOIN zitadel.projections.actions ON zitadel.projections.flows_triggers.action_id = zitadel.projections.actions.id`),
					nil,
					nil,
				),
			},
			object: &Flow{TriggerActions: map[domain.TriggerType][]*Action{}},
		},
		{
			name:    "prepareFlowQuery one action",
			prepare: prepareFlowQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(`SELECT zitadel.projections.actions.id,`+
						` zitadel.projections.actions.creation_date,`+
						` zitadel.projections.actions.change_date,`+
						` zitadel.projections.actions.resource_owner,`+
						` zitadel.projections.actions.action_state,`+
						` zitadel.projections.actions.sequence,`+
						` zitadel.projections.actions.name,`+
						` zitadel.projections.actions.script,`+
						` zitadel.projections.flows_triggers.trigger_type,`+
						` zitadel.projections.flows_triggers.trigger_sequence,`+
						` zitadel.projections.flows_triggers.flow_type`+
						` FROM zitadel.projections.flows_triggers`+
						` LEFT JOIN zitadel.projections.actions ON zitadel.projections.flows_triggers.action_id = zitadel.projections.actions.id`),
					[]string{
						"id",
						"creation_date",
						"change_date",
						"resource_owner",
						"state",
						"sequence",
						"name",
						"script",
						"trigger_type",
						"trigger_sequence",
						"flow_type",
					},
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
							domain.TriggerTypePreCreation,
							uint64(20211109),
							domain.FlowTypeExternalAuthentication,
						},
					},
				),
			},
			object: &Flow{
				Type: domain.FlowTypeExternalAuthentication,
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
						},
					},
				},
			},
		},
		{
			name:    "prepareFlowQuery multiple actions",
			prepare: prepareFlowQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(`SELECT zitadel.projections.actions.id,`+
						` zitadel.projections.actions.creation_date,`+
						` zitadel.projections.actions.change_date,`+
						` zitadel.projections.actions.resource_owner,`+
						` zitadel.projections.actions.action_state,`+
						` zitadel.projections.actions.sequence,`+
						` zitadel.projections.actions.name,`+
						` zitadel.projections.actions.script,`+
						` zitadel.projections.flows_triggers.trigger_type,`+
						` zitadel.projections.flows_triggers.trigger_sequence,`+
						` zitadel.projections.flows_triggers.flow_type`+
						` FROM zitadel.projections.flows_triggers`+
						` LEFT JOIN zitadel.projections.actions ON zitadel.projections.flows_triggers.action_id = zitadel.projections.actions.id`),
					[]string{
						"id",
						"creation_date",
						"change_date",
						"resource_owner",
						"state",
						"sequence",
						"name",
						"script",
						"trigger_type",
						"trigger_sequence",
						"flow_type",
					},
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
							domain.TriggerTypePreCreation,
							uint64(20211109),
							domain.FlowTypeExternalAuthentication,
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
							domain.TriggerTypePostCreation,
							uint64(20211109),
							domain.FlowTypeExternalAuthentication,
						},
					},
				),
			},
			object: &Flow{
				Type: domain.FlowTypeExternalAuthentication,
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
						},
					},
				},
			},
		},
		{
			name:    "prepareFlowQuery no action",
			prepare: prepareFlowQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(`SELECT zitadel.projections.actions.id,`+
						` zitadel.projections.actions.creation_date,`+
						` zitadel.projections.actions.change_date,`+
						` zitadel.projections.actions.resource_owner,`+
						` zitadel.projections.actions.action_state,`+
						` zitadel.projections.actions.sequence,`+
						` zitadel.projections.actions.name,`+
						` zitadel.projections.actions.script,`+
						` zitadel.projections.flows_triggers.trigger_type,`+
						` zitadel.projections.flows_triggers.trigger_sequence,`+
						` zitadel.projections.flows_triggers.flow_type`+
						` FROM zitadel.projections.flows_triggers`+
						` LEFT JOIN zitadel.projections.actions ON zitadel.projections.flows_triggers.action_id = zitadel.projections.actions.id`),
					[]string{
						"id",
						"creation_date",
						"change_date",
						"resource_owner",
						"state",
						"sequence",
						"name",
						"script",
						"trigger_type",
						"trigger_sequence",
						"flow_type",
					},
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
							domain.TriggerTypePostCreation,
							uint64(20211109),
							domain.FlowTypeExternalAuthentication,
						},
					},
				),
			},
			object: &Flow{
				Type:           domain.FlowTypeExternalAuthentication,
				TriggerActions: map[domain.TriggerType][]*Action{},
			},
		},
		{
			name:    "prepareFlowQuery sql err",
			prepare: prepareFlowQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(`SELECT zitadel.projections.actions.id,`+
						` zitadel.projections.actions.creation_date,`+
						` zitadel.projections.actions.change_date,`+
						` zitadel.projections.actions.resource_owner,`+
						` zitadel.projections.actions.action_state,`+
						` zitadel.projections.actions.sequence,`+
						` zitadel.projections.actions.name,`+
						` zitadel.projections.actions.script,`+
						` zitadel.projections.flows_triggers.trigger_type,`+
						` zitadel.projections.flows_triggers.trigger_sequence,`+
						` zitadel.projections.flows_triggers.flow_type`+
						` FROM zitadel.projections.flows_triggers`+
						` LEFT JOIN zitadel.projections.actions ON zitadel.projections.flows_triggers.action_id = zitadel.projections.actions.id`),
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
					regexp.QuoteMeta(`SELECT zitadel.projections.actions.id,`+
						` zitadel.projections.actions.creation_date,`+
						` zitadel.projections.actions.change_date,`+
						` zitadel.projections.actions.resource_owner,`+
						` zitadel.projections.actions.action_state,`+
						` zitadel.projections.actions.sequence,`+
						` zitadel.projections.actions.name,`+
						` zitadel.projections.actions.script`+
						` FROM zitadel.projections.flows_triggers`+
						` LEFT JOIN zitadel.projections.actions ON zitadel.projections.flows_triggers.action_id = zitadel.projections.actions.id`),
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
					regexp.QuoteMeta(`SELECT zitadel.projections.actions.id,`+
						` zitadel.projections.actions.creation_date,`+
						` zitadel.projections.actions.change_date,`+
						` zitadel.projections.actions.resource_owner,`+
						` zitadel.projections.actions.action_state,`+
						` zitadel.projections.actions.sequence,`+
						` zitadel.projections.actions.name,`+
						` zitadel.projections.actions.script`+
						` FROM zitadel.projections.flows_triggers`+
						` LEFT JOIN zitadel.projections.actions ON zitadel.projections.flows_triggers.action_id = zitadel.projections.actions.id`),
					[]string{
						"id",
						"creation_date",
						"change_date",
						"resource_owner",
						"state",
						"sequence",
						"name",
						"script",
					},
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
				},
			},
		},
		{
			name:    "prepareTriggerActionsQuery multiple results",
			prepare: prepareTriggerActionsQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(`SELECT zitadel.projections.actions.id,`+
						` zitadel.projections.actions.creation_date,`+
						` zitadel.projections.actions.change_date,`+
						` zitadel.projections.actions.resource_owner,`+
						` zitadel.projections.actions.action_state,`+
						` zitadel.projections.actions.sequence,`+
						` zitadel.projections.actions.name,`+
						` zitadel.projections.actions.script`+
						` FROM zitadel.projections.flows_triggers`+
						` LEFT JOIN zitadel.projections.actions ON zitadel.projections.flows_triggers.action_id = zitadel.projections.actions.id`),
					[]string{
						"id",
						"creation_date",
						"change_date",
						"resource_owner",
						"state",
						"sequence",
						"name",
						"script",
					},
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
				},
			},
		},
		{
			name:    "prepareTriggerActionsQuery sql err",
			prepare: prepareTriggerActionsQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(`SELECT zitadel.projections.actions.id,`+
						` zitadel.projections.actions.creation_date,`+
						` zitadel.projections.actions.change_date,`+
						` zitadel.projections.actions.resource_owner,`+
						` zitadel.projections.actions.action_state,`+
						` zitadel.projections.actions.sequence,`+
						` zitadel.projections.actions.name,`+
						` zitadel.projections.actions.script`+
						` FROM zitadel.projections.flows_triggers`+
						` LEFT JOIN zitadel.projections.actions ON zitadel.projections.flows_triggers.action_id = zitadel.projections.actions.id`),
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
					regexp.QuoteMeta(`SELECT zitadel.projections.flows_triggers.flow_type`+
						` FROM zitadel.projections.flows_triggers`),
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
					regexp.QuoteMeta(`SELECT zitadel.projections.flows_triggers.flow_type`+
						` FROM zitadel.projections.flows_triggers`),
					[]string{
						"flow_type",
					},
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
					regexp.QuoteMeta(`SELECT zitadel.projections.flows_triggers.flow_type`+
						` FROM zitadel.projections.flows_triggers`),
					[]string{
						"flow_type",
					},
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
					regexp.QuoteMeta(`SELECT zitadel.projections.flows_triggers.flow_type`+
						` FROM zitadel.projections.flows_triggers`),
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
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err)
		})
	}
}
