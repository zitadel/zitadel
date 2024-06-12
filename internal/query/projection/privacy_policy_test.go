package projection

import (
	"testing"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestPrivacyPolicyProjection_reduces(t *testing.T) {
	type args struct {
		event func(t *testing.T) eventstore.Event
	}
	tests := []struct {
		name   string
		args   args
		reduce func(event eventstore.Event) (*handler.Statement, error)
		want   wantReduce
	}{
		{
			name: "org reduceAdded",
			args: args{
				event: getEvent(
					testEvent(
						org.PrivacyPolicyAddedEventType,
						org.AggregateType,
						[]byte(`{
						"tosLink": "http://tos.link",
						"privacyLink": "http://privacy.link",
						"helpLink": "http://help.link",
						"docsLink": "http://docs.link",
						"customLink": "http://custom.link",
						"customLinkText": "Custom Link",
						"supportEmail": "support@example.com"}`),
					), org.PrivacyPolicyAddedEventMapper),
			},
			reduce: (&privacyPolicyProjection{}).reduceAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.privacy_policies4 (creation_date, change_date, sequence, id, state, privacy_link, tos_link, help_link, support_email, docs_link, custom_link, custom_link_text, is_default, resource_owner, instance_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)",
							expectedArgs: []interface{}{
								anyArg{},
								anyArg{},
								uint64(15),
								"agg-id",
								domain.PolicyStateActive,
								"http://privacy.link",
								"http://tos.link",
								"http://help.link",
								domain.EmailAddress("support@example.com"),
								"http://docs.link",
								"http://custom.link",
								"Custom Link",
								false,
								"ro-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name:   "org reduceChanged",
			reduce: (&privacyPolicyProjection{}).reduceChanged,
			args: args{
				event: getEvent(
					testEvent(
						org.PrivacyPolicyChangedEventType,
						org.AggregateType,
						[]byte(`{
						"tosLink": "http://tos.link",
						"privacyLink": "http://privacy.link",
						"helpLink": "http://help.link",
						"docsLink": "http://docs.link",
						"customLink": "http://custom.link",
						"customLinkText": "Custom Link",
						"supportEmail": "support@example.com"}`),
					), org.PrivacyPolicyChangedEventMapper),
			},
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.privacy_policies4 SET (change_date, sequence, privacy_link, tos_link, help_link, support_email, docs_link, custom_link, custom_link_text) = ($1, $2, $3, $4, $5, $6, $7, $8, $9) WHERE (id = $10) AND (instance_id = $11)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"http://privacy.link",
								"http://tos.link",
								"http://help.link",
								domain.EmailAddress("support@example.com"),
								"http://docs.link",
								"http://custom.link",
								"Custom Link",
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name:   "org reduceRemoved",
			reduce: (&privacyPolicyProjection{}).reduceRemoved,
			args: args{
				event: getEvent(
					testEvent(
						org.PrivacyPolicyRemovedEventType,
						org.AggregateType,
						nil,
					), org.PrivacyPolicyRemovedEventMapper),
			},
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.privacy_policies4 WHERE (id = $1) AND (instance_id = $2)",
							expectedArgs: []interface{}{
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		}, {
			name: "instance reduceInstanceRemoved",
			args: args{
				event: getEvent(
					testEvent(
						instance.InstanceRemovedEventType,
						instance.AggregateType,
						nil,
					), instance.InstanceRemovedEventMapper),
			},
			reduce: reduceInstanceRemovedHelper(PrivacyPolicyInstanceIDCol),
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.privacy_policies4 WHERE (instance_id = $1)",
							expectedArgs: []interface{}{
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name:   "instance reduceAdded",
			reduce: (&privacyPolicyProjection{}).reduceAdded,
			args: args{
				event: getEvent(
					testEvent(
						instance.PrivacyPolicyAddedEventType,
						instance.AggregateType,
						[]byte(`{
						"tosLink": "http://tos.link",
						"privacyLink": "http://privacy.link",
						"helpLink": "http://help.link",
						"docsLink": "http://docs.link",
						"customLink": "http://custom.link",
						"customLinkText": "Custom Link",
						"supportEmail": "support@example.com"}`),
					), instance.PrivacyPolicyAddedEventMapper),
			},
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.privacy_policies4 (creation_date, change_date, sequence, id, state, privacy_link, tos_link, help_link, support_email, docs_link, custom_link, custom_link_text, is_default, resource_owner, instance_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)",
							expectedArgs: []interface{}{
								anyArg{},
								anyArg{},
								uint64(15),
								"agg-id",
								domain.PolicyStateActive,
								"http://privacy.link",
								"http://tos.link",
								"http://help.link",
								domain.EmailAddress("support@example.com"),
								"http://docs.link",
								"http://custom.link",
								"Custom Link",
								true,
								"ro-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name:   "instance reduceChanged",
			reduce: (&privacyPolicyProjection{}).reduceChanged,
			args: args{
				event: getEvent(
					testEvent(
						instance.PrivacyPolicyChangedEventType,
						instance.AggregateType,
						[]byte(`{
						"tosLink": "http://tos.link",
						"privacyLink": "http://privacy.link",
						"helpLink": "http://help.link",
						"docsLink": "http://docs.link",
						"customLink": "http://custom.link",
						"customLinkText": "Custom Link",
						"supportEmail": "support@example.com"}`),
					), instance.PrivacyPolicyChangedEventMapper),
			},
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.privacy_policies4 SET (change_date, sequence, privacy_link, tos_link, help_link, support_email, docs_link, custom_link, custom_link_text) = ($1, $2, $3, $4, $5, $6, $7, $8, $9) WHERE (id = $10) AND (instance_id = $11)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"http://privacy.link",
								"http://tos.link",
								"http://help.link",
								domain.EmailAddress("support@example.com"),
								"http://docs.link",
								"http://custom.link",
								"Custom Link",
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name:   "org.reduceOwnerRemoved",
			reduce: (&privacyPolicyProjection{}).reduceOwnerRemoved,
			args: args{
				event: getEvent(
					testEvent(
						org.OrgRemovedEventType,
						org.AggregateType,
						nil,
					), org.OrgRemovedEventMapper),
			},
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.privacy_policies4 WHERE (instance_id = $1) AND (resource_owner = $2)",
							expectedArgs: []interface{}{
								"instance-id",
								"agg-id",
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := baseEvent(t)
			got, err := tt.reduce(event)
			if ok := zerrors.IsErrorInvalidArgument(err); !ok {
				t.Errorf("no wrong event mapping: %v, got: %v", err, got)
			}

			event = tt.args.event(t)
			got, err = tt.reduce(event)
			assertReduce(t, got, err, PrivacyPolicyTable, tt.want)
		})
	}
}
