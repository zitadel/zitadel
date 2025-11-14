package projection

import (
	"testing"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestSMTPConfigProjection_reduces(t *testing.T) {
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
			name: "reduceSMTPConfigChanged (no id)",
			args: args{
				event: getEvent(
					testEvent(
						instance.SMTPConfigChangedEventType,
						instance.AggregateType,
						[]byte(`{
						"instance_id": "instance-id",	
						"resource_owner": "ro-id",	
						"aggregate_id": "agg-id",	
						"description": "test",
						"tls": true,
						"senderAddress": "sender",
						"senderName": "name",
						"replyToAddress": "reply-to",
						"host": "host",
						"user": "user"		
					}`,
						),
					), eventstore.GenericEventMapper[instance.SMTPConfigChangedEvent]),
			},
			reduce: (&smtpConfigProjection{}).reduceSMTPConfigChanged,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.smtp_configs5 SET (change_date, sequence, description) = ($1, $2, $3) WHERE (id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"test",
								"ro-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.smtp_configs5_smtp SET (tls, sender_address, sender_name, reply_to_address, host, username) = ($1, $2, $3, $4, $5, $6) WHERE (id = $7) AND (instance_id = $8)",
							expectedArgs: []interface{}{
								true,
								"sender",
								"name",
								"reply-to",
								"host",
								"user",
								"ro-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceSMTPConfigChanged",
			args: args{
				event: getEvent(
					testEvent(
						instance.SMTPConfigChangedEventType,
						instance.AggregateType,
						[]byte(`{
						"instance_id": "instance-id",	
						"resource_owner": "ro-id",	
						"aggregate_id": "agg-id",
						"id": "config-id",		
						"description": "test",
						"tls": true,
						"senderAddress": "sender",
						"senderName": "name",
						"replyToAddress": "reply-to",
						"host": "host",
						"user": "user"		
					}`,
						),
					), eventstore.GenericEventMapper[instance.SMTPConfigChangedEvent]),
			},
			reduce: (&smtpConfigProjection{}).reduceSMTPConfigChanged,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.smtp_configs5 SET (change_date, sequence, description) = ($1, $2, $3) WHERE (id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"test",
								"config-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.smtp_configs5_smtp SET (tls, sender_address, sender_name, reply_to_address, host, username) = ($1, $2, $3, $4, $5, $6) WHERE (id = $7) AND (instance_id = $8)",
							expectedArgs: []interface{}{
								true,
								"sender",
								"name",
								"reply-to",
								"host",
								"user",
								"config-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceSMTPConfigChanged, description",
			args: args{
				event: getEvent(
					testEvent(
						instance.SMTPConfigChangedEventType,
						instance.AggregateType,
						[]byte(`{
						"instance_id": "instance-id",	
						"resource_owner": "ro-id",	
						"aggregate_id": "agg-id",
						"id": "config-id",	
						"description": "test"					
					}`,
						),
					), eventstore.GenericEventMapper[instance.SMTPConfigChangedEvent]),
			},
			reduce: (&smtpConfigProjection{}).reduceSMTPConfigChanged,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.smtp_configs5 SET (change_date, sequence, description) = ($1, $2, $3) WHERE (id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"test",
								"config-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceSMTPConfigChanged, senderAddress",
			args: args{
				event: getEvent(
					testEvent(
						instance.SMTPConfigChangedEventType,
						instance.AggregateType,
						[]byte(`{
						"instance_id": "instance-id",	
						"resource_owner": "ro-id",	
						"aggregate_id": "agg-id",
						"id": "config-id",	
						"senderAddress": "sender"
					}`,
						),
					), eventstore.GenericEventMapper[instance.SMTPConfigChangedEvent]),
			},
			reduce: (&smtpConfigProjection{}).reduceSMTPConfigChanged,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.smtp_configs5 SET (change_date, sequence) = ($1, $2) WHERE (id = $3) AND (instance_id = $4)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"config-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.smtp_configs5_smtp SET sender_address = $1 WHERE (id = $2) AND (instance_id = $3)",
							expectedArgs: []interface{}{
								"sender",
								"config-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceSMTPConfigHTTPChanged",
			args: args{
				event: getEvent(
					testEvent(
						instance.SMTPConfigHTTPChangedEventType,
						instance.AggregateType,
						[]byte(`{
						"instance_id": "instance-id",	
						"resource_owner": "ro-id",	
						"aggregate_id": "agg-id",
						"id": "config-id",		
						"description": "test",
						"endpoint": "endpoint",
						"signingKey": { "cryptoType": 0, "algorithm": "RSA-265", "keyId": "key-id" }
					}`,
						),
					), eventstore.GenericEventMapper[instance.SMTPConfigHTTPChangedEvent]),
			},
			reduce: (&smtpConfigProjection{}).reduceSMTPConfigHTTPChanged,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.smtp_configs5 SET (change_date, sequence, description) = ($1, $2, $3) WHERE (id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"test",
								"config-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.smtp_configs5_http SET (endpoint, signing_key) = ($1, $2) WHERE (id = $3) AND (instance_id = $4)",
							expectedArgs: []interface{}{
								"endpoint",
								anyArg{},
								"config-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceSMTPConfigHTTPChanged, description",
			args: args{
				event: getEvent(
					testEvent(
						instance.SMTPConfigHTTPChangedEventType,
						instance.AggregateType,
						[]byte(`{
						"instance_id": "instance-id",	
						"resource_owner": "ro-id",	
						"aggregate_id": "agg-id",
						"id": "config-id",	
						"description": "test"					
					}`,
						),
					), eventstore.GenericEventMapper[instance.SMTPConfigHTTPChangedEvent]),
			},
			reduce: (&smtpConfigProjection{}).reduceSMTPConfigHTTPChanged,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.smtp_configs5 SET (change_date, sequence, description) = ($1, $2, $3) WHERE (id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"test",
								"config-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceSMTPConfigHTTPChanged, endpoint",
			args: args{
				event: getEvent(
					testEvent(
						instance.SMTPConfigHTTPChangedEventType,
						instance.AggregateType,
						[]byte(`{
						"instance_id": "instance-id",	
						"resource_owner": "ro-id",	
						"aggregate_id": "agg-id",
						"id": "config-id",	
						"endpoint": "endpoint"
					}`,
						),
					), eventstore.GenericEventMapper[instance.SMTPConfigHTTPChangedEvent]),
			},
			reduce: (&smtpConfigProjection{}).reduceSMTPConfigHTTPChanged,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.smtp_configs5 SET (change_date, sequence) = ($1, $2) WHERE (id = $3) AND (instance_id = $4)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"config-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.smtp_configs5_http SET endpoint = $1 WHERE (id = $2) AND (instance_id = $3)",
							expectedArgs: []interface{}{
								"endpoint",
								"config-id",
								"instance-id",
							},
						},
					},
				},
			},
		}, {
			name: "reduceSMTPConfigHTTPChanged, signing key",
			args: args{
				event: getEvent(
					testEvent(
						instance.SMTPConfigHTTPChangedEventType,
						instance.AggregateType,
						[]byte(`{
						"instance_id": "instance-id",	
						"resource_owner": "ro-id",	
						"aggregate_id": "agg-id",
						"id": "config-id",	
						"signingKey": { "cryptoType": 0, "algorithm": "RSA-265", "keyId": "key-id" }
					}`,
						),
					), eventstore.GenericEventMapper[instance.SMTPConfigHTTPChangedEvent]),
			},
			reduce: (&smtpConfigProjection{}).reduceSMTPConfigHTTPChanged,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.smtp_configs5 SET (change_date, sequence) = ($1, $2) WHERE (id = $3) AND (instance_id = $4)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"config-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.smtp_configs5_http SET signing_key = $1 WHERE (id = $2) AND (instance_id = $3)",
							expectedArgs: []interface{}{
								anyArg{},
								"config-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceSMTPConfigAdded (no id)",
			args: args{
				event: getEvent(
					testEvent(
						instance.SMTPConfigAddedEventType,
						instance.AggregateType,
						[]byte(`{
						"instance_id": "instance-id",	
						"resource_owner": "ro-id",	
						"aggregate_id": "agg-id",
						"tls": true,
						"senderAddress": "sender",
						"senderName": "name",
						"replyToAddress": "reply-to",
						"host": "host",
						"user": "user",
						"password": {
							"cryptoType": 0,
							"algorithm": "RSA-265",
							"keyId": "key-id"
						}
					}`),
					), eventstore.GenericEventMapper[instance.SMTPConfigAddedEvent]),
			},
			reduce: (&smtpConfigProjection{}).reduceSMTPConfigAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.smtp_configs5 (creation_date, change_date, instance_id, resource_owner, aggregate_id, id, sequence, state, description) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
							expectedArgs: []interface{}{
								anyArg{},
								anyArg{},
								"instance-id",
								"ro-id",
								"agg-id",
								"ro-id",
								uint64(15),
								domain.SMTPConfigStateActive,
								"generic",
							},
						},
						{
							expectedStmt: "INSERT INTO projections.smtp_configs5_smtp (instance_id, id, tls, sender_address, sender_name, reply_to_address, host, username, password) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
							expectedArgs: []interface{}{
								"instance-id",
								"ro-id",
								true,
								"sender",
								"name",
								"reply-to",
								"host",
								"user",
								anyArg{},
							},
						},
					},
				},
			},
		},
		{
			name: "reduceSMTPConfigAdded",
			args: args{
				event: getEvent(
					testEvent(
						instance.SMTPConfigAddedEventType,
						instance.AggregateType,
						[]byte(`{
						"instance_id": "instance-id",	
						"resource_owner": "ro-id",	
						"aggregate_id": "agg-id",
						"id": "config-id",		
						"description": "test",
						"tls": true,
						"senderAddress": "sender",
						"senderName": "name",
						"replyToAddress": "reply-to",
						"host": "host",
						"user": "user",
						"password": {
							"cryptoType": 0,
							"algorithm": "RSA-265",
							"keyId": "key-id"
						}
					}`),
					), eventstore.GenericEventMapper[instance.SMTPConfigAddedEvent]),
			},
			reduce: (&smtpConfigProjection{}).reduceSMTPConfigAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.smtp_configs5 (creation_date, change_date, instance_id, resource_owner, aggregate_id, id, sequence, state, description) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
							expectedArgs: []interface{}{
								anyArg{},
								anyArg{},
								"instance-id",
								"ro-id",
								"agg-id",
								"config-id",
								uint64(15),
								domain.SMTPConfigStateInactive,
								"test",
							},
						},
						{
							expectedStmt: "INSERT INTO projections.smtp_configs5_smtp (instance_id, id, tls, sender_address, sender_name, reply_to_address, host, username, password) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
							expectedArgs: []interface{}{
								"instance-id",
								"config-id",
								true,
								"sender",
								"name",
								"reply-to",
								"host",
								"user",
								anyArg{},
							},
						},
					},
				},
			},
		},
		{
			name: "reduceSMTPConfigHTTPAdded",
			args: args{
				event: getEvent(
					testEvent(
						instance.SMTPConfigHTTPAddedEventType,
						instance.AggregateType,
						[]byte(`{
						"instance_id": "instance-id",	
						"resource_owner": "ro-id",	
						"aggregate_id": "agg-id",
						"id": "config-id",		
						"description": "test",
						"senderAddress": "sender",
						"endpoint": "endpoint",
						"signingKey": { "cryptoType": 0, "algorithm": "RSA-265", "keyId": "key-id" }
					}`),
					), eventstore.GenericEventMapper[instance.SMTPConfigHTTPAddedEvent]),
			},
			reduce: (&smtpConfigProjection{}).reduceSMTPConfigHTTPAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.smtp_configs5 (creation_date, change_date, instance_id, resource_owner, aggregate_id, id, sequence, state, description) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
							expectedArgs: []interface{}{
								anyArg{},
								anyArg{},
								"instance-id",
								"ro-id",
								"agg-id",
								"config-id",
								uint64(15),
								domain.SMTPConfigStateInactive,
								"test",
							},
						},
						{
							expectedStmt: "INSERT INTO projections.smtp_configs5_http (instance_id, id, endpoint, signing_key) VALUES ($1, $2, $3, $4)",
							expectedArgs: []interface{}{
								"instance-id",
								"config-id",
								"endpoint",
								anyArg{},
							},
						},
					},
				},
			},
		},
		{
			name: "reduceSMTPConfigActivated",
			args: args{
				event: getEvent(testEvent(
					instance.SMTPConfigActivatedEventType,
					instance.AggregateType,
					[]byte(`{ 
						"instance_id": "instance-id",	
						"resource_owner": "ro-id",	
						"aggregate_id": "agg-id",
						"id": "config-id"		
					}`),
				), eventstore.GenericEventMapper[instance.SMTPConfigActivatedEvent]),
			},
			reduce: (&smtpConfigProjection{}).reduceSMTPConfigActivated,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.smtp_configs5 SET (change_date, sequence, state) = ($1, $2, $3) WHERE (NOT (id = $4)) AND (state = $5) AND (instance_id = $6)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								domain.SMTPConfigStateInactive,
								"config-id",
								domain.SMTPConfigStateActive,
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.smtp_configs5 SET (change_date, sequence, state) = ($1, $2, $3) WHERE (id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								domain.SMTPConfigStateActive,
								"config-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceSMTPConfigActivated (no id)",
			args: args{
				event: getEvent(testEvent(
					instance.SMTPConfigActivatedEventType,
					instance.AggregateType,
					[]byte(`{ 
						"instance_id": "instance-id",	
						"resource_owner": "ro-id",	
						"aggregate_id": "agg-id"
					}`),
				), eventstore.GenericEventMapper[instance.SMTPConfigActivatedEvent]),
			},
			reduce: (&smtpConfigProjection{}).reduceSMTPConfigActivated,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.smtp_configs5 SET (change_date, sequence, state) = ($1, $2, $3) WHERE (NOT (id = $4)) AND (state = $5) AND (instance_id = $6)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								domain.SMTPConfigStateInactive,
								"ro-id",
								domain.SMTPConfigStateActive,
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.smtp_configs5 SET (change_date, sequence, state) = ($1, $2, $3) WHERE (id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								domain.SMTPConfigStateActive,
								"ro-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceSMTPConfigDeactivated",
			args: args{
				event: getEvent(testEvent(
					instance.SMTPConfigDeactivatedEventType,
					instance.AggregateType,
					[]byte(`{ 
						"instance_id": "instance-id",	
						"resource_owner": "ro-id",	
						"aggregate_id": "agg-id",
						"id": "config-id"		
					}`),
				), eventstore.GenericEventMapper[instance.SMTPConfigDeactivatedEvent]),
			},
			reduce: (&smtpConfigProjection{}).reduceSMTPConfigDeactivated,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.smtp_configs5 SET (change_date, sequence, state) = ($1, $2, $3) WHERE (id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								domain.SMTPConfigStateInactive,
								"config-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceSMTPConfigDeactivated (no id)",
			args: args{
				event: getEvent(testEvent(
					instance.SMTPConfigDeactivatedEventType,
					instance.AggregateType,
					[]byte(`{ 
						"instance_id": "instance-id",	
						"resource_owner": "ro-id",	
						"aggregate_id": "agg-id"	
					}`),
				), eventstore.GenericEventMapper[instance.SMTPConfigDeactivatedEvent]),
			},
			reduce: (&smtpConfigProjection{}).reduceSMTPConfigDeactivated,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.smtp_configs5 SET (change_date, sequence, state) = ($1, $2, $3) WHERE (id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								domain.SMTPConfigStateInactive,
								"ro-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceSMTPConfigPasswordChanged",
			args: args{
				event: getEvent(
					testEvent(
						instance.SMTPConfigPasswordChangedEventType,
						instance.AggregateType,
						[]byte(`{
						"instance_id": "instance-id",	
						"resource_owner": "ro-id",	
						"aggregate_id": "agg-id",
						"id": "config-id",		
						"password": {
							"cryptoType": 0,
							"algorithm": "RSA-265",
							"keyId": "key-id"
						}
					}`),
					), eventstore.GenericEventMapper[instance.SMTPConfigPasswordChangedEvent]),
			},
			reduce: (&smtpConfigProjection{}).reduceSMTPConfigPasswordChanged,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.smtp_configs5_smtp SET password = $1 WHERE (id = $2) AND (instance_id = $3)",
							expectedArgs: []interface{}{
								anyArg{},
								"config-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.smtp_configs5 SET (change_date, sequence) = ($1, $2) WHERE (id = $3) AND (instance_id = $4)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"config-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceSMTPConfigRemoved",
			args: args{
				event: getEvent(testEvent(
					instance.SMTPConfigRemovedEventType,
					instance.AggregateType,
					[]byte(`{ 
						"instance_id": "instance-id",	
						"resource_owner": "ro-id",	
						"aggregate_id": "agg-id",
						"id": "config-id"
}`),
				), eventstore.GenericEventMapper[instance.SMTPConfigRemovedEvent]),
			},
			reduce: (&smtpConfigProjection{}).reduceSMTPConfigRemoved,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.smtp_configs5 WHERE (id = $1) AND (instance_id = $2)",
							expectedArgs: []interface{}{
								"config-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceSMTPConfigRemoved (no id)",
			args: args{
				event: getEvent(testEvent(
					instance.SMTPConfigRemovedEventType,
					instance.AggregateType,
					[]byte(`{ 
						"instance_id": "instance-id",	
						"resource_owner": "ro-id",	
						"aggregate_id": "agg-id"
}`),
				), eventstore.GenericEventMapper[instance.SMTPConfigRemovedEvent]),
			},
			reduce: (&smtpConfigProjection{}).reduceSMTPConfigRemoved,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.smtp_configs5 WHERE (id = $1) AND (instance_id = $2)",
							expectedArgs: []interface{}{
								"ro-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "instance reduceInstanceRemoved",
			args: args{
				event: getEvent(
					testEvent(
						instance.InstanceRemovedEventType,
						instance.AggregateType,
						nil,
					), instance.InstanceRemovedEventMapper),
			},
			reduce: reduceInstanceRemovedHelper(SMTPConfigColumnInstanceID),
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.smtp_configs5 WHERE (instance_id = $1)",
							expectedArgs: []interface{}{
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
			assertReduce(t, got, err, SMTPConfigProjectionTable, tt.want)
		})
	}
}
