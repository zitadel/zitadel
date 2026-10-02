package projection

import (
	"testing"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestSecurityPolicyProjection_reducesClientIDMetadataDocumentAllowedURLs(t *testing.T) {
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
			// The policy row is upserted first because it only exists once the policy was set,
			// and the url is removed before it is appended so it is never in the array twice.
			name:   "reduceClientIDMetadataDocumentAllowedURLAdded",
			reduce: (&securityPolicyProjection{}).reduceClientIDMetadataDocumentAllowedURLAdded,
			args: args{
				event: getEvent(
					testEvent(
						instance.SecurityPolicyClientIDMetadataDocumentAllowedURLAddedEventType,
						instance.AggregateType,
						[]byte(`{"url": "https://clients.example.com/"}`),
					), eventstore.GenericEventMapper[instance.SecurityPolicyClientIDMetadataDocumentAllowedURLAddedEvent]),
			},
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.security_policies3 (creation_date, change_date, instance_id, sequence) VALUES ($1, $2, $3, $4) ON CONFLICT (instance_id) DO UPDATE SET (creation_date, change_date, sequence) = (projections.security_policies3.creation_date, EXCLUDED.change_date, EXCLUDED.sequence)",
							expectedArgs: []interface{}{
								anyArg{},
								anyArg{},
								"instance-id",
								uint64(15),
							},
						},
						{
							expectedStmt: "UPDATE projections.security_policies3 SET client_id_metadata_document_allowed_urls = array_append(array_remove(client_id_metadata_document_allowed_urls, $1), $1) WHERE (instance_id = $2)",
							expectedArgs: []interface{}{
								"https://clients.example.com/",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name:   "reduceClientIDMetadataDocumentAllowedURLRemoved",
			reduce: (&securityPolicyProjection{}).reduceClientIDMetadataDocumentAllowedURLRemoved,
			args: args{
				event: getEvent(
					testEvent(
						instance.SecurityPolicyClientIDMetadataDocumentAllowedURLRemovedEventType,
						instance.AggregateType,
						[]byte(`{"url": "https://clients.example.com/"}`),
					), eventstore.GenericEventMapper[instance.SecurityPolicyClientIDMetadataDocumentAllowedURLRemovedEvent]),
			},
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.security_policies3 (creation_date, change_date, instance_id, sequence) VALUES ($1, $2, $3, $4) ON CONFLICT (instance_id) DO UPDATE SET (creation_date, change_date, sequence) = (projections.security_policies3.creation_date, EXCLUDED.change_date, EXCLUDED.sequence)",
							expectedArgs: []interface{}{
								anyArg{},
								anyArg{},
								"instance-id",
								uint64(15),
							},
						},
						{
							expectedStmt: "UPDATE projections.security_policies3 SET client_id_metadata_document_allowed_urls = array_remove(client_id_metadata_document_allowed_urls, $1) WHERE (instance_id = $2)",
							expectedArgs: []interface{}{
								"https://clients.example.com/",
								"instance-id",
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
			assertReduce(t, got, err, SecurityPolicyProjectionTable, tt.want)
		})
	}
}
