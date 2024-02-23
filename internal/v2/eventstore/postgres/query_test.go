package postgres

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/zitadel/zitadel/internal/v2/database"
	"github.com/zitadel/zitadel/internal/v2/eventstore"
)

func Test_filterQuery(t *testing.T) {
	type args struct {
		filters  []*eventstore.Filter
		instance string
	}
	type want struct {
		query string
		args  []any
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "aggregate type",
			args: args{
				instance: "instance",
				filters: []*eventstore.Filter{
					eventstore.NewFilter(
						eventstore.AppendAggregateFilter("aggregate"),
					),
				},
			},
			want: want{
				query: `SELECT created_at, event_type, "sequence", "position", payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision FROM eventstore.events2 WHERE instance_id = $1 AND aggregate_type = $2`, //TODO: ORDER BY position, in_tx_order`,
				args: []any{
					"instance",
					"aggregate",
				},
			},
		},
		// {
		// 	name: "descending",
		// 	args: args{
		// 		instance: "instance",
		// 		filters: []*eventstore.Filter{
		// 			eventstore.NewFilter(
		// 				eventstore.AppendAggregateFilter("aggregate"),
		// 			),
		// 		},
		// 	},
		// 	want: want{
		// 		query: `SELECT created_at, event_type, "sequence", "position", payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision FROM eventstore.events2 WHERE instance_id = $1 AND aggregate_type = $2 ORDER BY position DESC, in_tx_order DESC`,
		// 		args: []any{
		// 			"instance",
		// 			"aggregate",
		// 		},
		// 	},
		// },
		{
			name: "multiple aggregates",
			args: args{
				instance: "instance",
				filters: []*eventstore.Filter{
					eventstore.NewFilter(
						eventstore.AppendAggregateFilter("agg1"),
					),
					eventstore.NewFilter(
						eventstore.AppendAggregateFilter("agg2"),
						eventstore.AppendAggregateFilter("agg3"),
					),
				},
			},
			want: want{
				query: `SELECT created_at, event_type, "sequence", "position", payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision FROM eventstore.events2 WHERE instance_id = $1 AND (aggregate_type = $2 OR aggregate_type = ANY($3)) ORDER BY position, in_tx_order`,
				args: []any{
					"instance",
					"agg1",
					"instance",
					"agg2",
					"agg3",
				},
			},
		},
		{
			name: "multiple aggregates with ids",
			args: args{
				instance: "instance",
				filters: []*eventstore.Filter{
					eventstore.NewFilter(
						eventstore.AppendAggregateFilter("agg1", eventstore.WithAggregateID("id")),
					),
					eventstore.NewFilter(
						eventstore.AppendAggregateFilter("agg2", eventstore.WithAggregateID("id2")),
						eventstore.AppendAggregateFilter("agg3"),
					),
				},
			},
			want: want{
				query: `SELECT created_at, event_type, "sequence", "position", payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision FROM eventstore.events2 WHERE instance_id = $1 AND (aggregate_type = $2 OR aggregate_type = ANY($3)) ORDER BY position, in_tx_order`,
				args: []any{
					"instance",
					"agg1",
					"id",
					"instance",
					"agg2",
					"id2",
					"agg3",
				},
			},
		},
		// {
		// 	name: "multiple event queries and multiple filter in queries",
		// 	args: args{
		// 		filters: []*eventstore.Filter{
		// 			eventstore.NewFilter(
		// 				context.Background(),
		// 				eventstore.FilterEventQuery(
		// 					eventstore.FilterAggregateTypes("agg1"),
		// 					eventstore.FilterAggregateIDs("1", "2"),
		// 				),
		// 				eventstore.FilterEventQuery(
		// 					eventstore.FilterAggregateTypes("agg2", "agg3"),
		// 					eventstore.FilterAggregateIDs("3"),
		// 				),
		// 			),
		// 		},
		// 	},
		// 	want: want{
		// 		query: `SELECT created_at, event_type, "sequence", "position", payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision FROM eventstore.events2 WHERE instance_id = $1 AND ((aggregate_type = $2 AND aggregate_id = ANY($3)) OR (aggregate_type = ANY($4) AND aggregate_id = $5)) ORDER BY position, in_tx_order`,
		// 		args: []any{
		// 			"",
		// 			"agg1",
		// 			[]string{"1", "2"},
		// 			[]string{"agg2", "agg3"},
		// 			"3",
		// 		},
		// 	},
		// },
		{
			name: "milestones",
			args: args{
				filters: []*eventstore.Filter{
					eventstore.NewFilter(
						eventstore.AppendAggregateFilter(
							"instance",
							eventstore.AppendEvent(
								eventstore.WithEventType("instance.added"),
							),
						),
						eventstore.WithLimit(1),
					),
					eventstore.NewFilter(
						eventstore.AppendAggregateFilter(
							"instance",
							eventstore.AppendEvent(
								eventstore.WithEventType("instance.removed"),
							),
						),
						eventstore.WithLimit(1),
					),
					eventstore.NewFilter(
						eventstore.AppendAggregateFilter(
							"instance",
							eventstore.AppendEvent(
								eventstore.WithEventType("instance.domain.primary.set"),
								eventstore.WithCreatorList(database.NewListNotContains("", "SYSTEM")),
							),
						),
						eventstore.WithLimit(1),
					),
					eventstore.NewFilter(
						eventstore.AppendAggregateFilter(
							"project",
							eventstore.AppendEvent(
								eventstore.WithEventType("project.added"),
								eventstore.WithCreatorList(database.NewListNotContains("", "SYSTEM")),
							),
						),
						eventstore.WithLimit(1),
					),
					eventstore.NewFilter(
						eventstore.AppendAggregateFilter(
							"project",
							eventstore.AppendEvent(
								eventstore.WithCreatorList(database.NewListNotContains("", "SYSTEM")),
								eventstore.WithEventType("project.application.added"),
							),
						),
						eventstore.WithLimit(1),
					),
					eventstore.NewFilter(
						eventstore.AppendAggregateFilter(
							"user",
							eventstore.AppendEvent(
								eventstore.WithEventType("user.token.added"),
							),
						),
						// used because we need to check for first login and an app which is not console
						eventstore.WithPosition(database.NewNumberAtLeast(12), database.NewNumberGreater(4)),
					),
					eventstore.NewFilter(
						eventstore.AppendAggregateFilter(
							"instance",
							eventstore.AppendEvent(
								eventstore.WithEventType("instance.idp.config.added"),
							),
							eventstore.AppendEvent(
								eventstore.WithEventType("instance.idp.oauth.added"),
							),
							eventstore.AppendEvent(
								eventstore.WithEventType("instance.idp.oidc.added"),
							),
							eventstore.AppendEvent(
								eventstore.WithEventType("instance.idp.jwt.added"),
							),
							eventstore.AppendEvent(
								eventstore.WithEventType("instance.idp.azure.added"),
							),
							eventstore.AppendEvent(
								eventstore.WithEventType("instance.idp.github.added"),
							),
							eventstore.AppendEvent(
								eventstore.WithEventType("instance.idp.github.enterprise.added"),
							),
							eventstore.AppendEvent(
								eventstore.WithEventType("instance.idp.gitlab.added"),
							),
							eventstore.AppendEvent(
								eventstore.WithEventType("instance.idp.gitlab.selfhosted.added"),
							),
							eventstore.AppendEvent(
								eventstore.WithEventType("instance.idp.google.added"),
							),
							eventstore.AppendEvent(
								eventstore.WithEventType("instance.idp.ldap.added"),
							),
							eventstore.AppendEvent(
								eventstore.WithEventType("instance.idp.config.apple.added"),
							),
							eventstore.AppendEvent(
								eventstore.WithEventType("instance.idp.saml.added"),
							),
						),
						eventstore.AppendAggregateFilter(
							"org",
							eventstore.AppendEvent(
								eventstore.WithEventType("org.idp.config.added"),
							),
							eventstore.AppendEvent(
								eventstore.WithEventType("org.idp.oauth.added"),
							),
							eventstore.AppendEvent(
								eventstore.WithEventType("org.idp.oidc.added"),
							),
							eventstore.AppendEvent(
								eventstore.WithEventType("org.idp.jwt.added"),
							),
							eventstore.AppendEvent(
								eventstore.WithEventType("org.idp.azure.added"),
							),
							eventstore.AppendEvent(
								eventstore.WithEventType("org.idp.github.added"),
							),
							eventstore.AppendEvent(
								eventstore.WithEventType("org.idp.github.enterprise.added"),
							),
							eventstore.AppendEvent(
								eventstore.WithEventType("org.idp.gitlab.added"),
							),
							eventstore.AppendEvent(
								eventstore.WithEventType("org.idp.gitlab.selfhosted.added"),
							),
							eventstore.AppendEvent(
								eventstore.WithEventType("org.idp.google.added"),
							),
							eventstore.AppendEvent(
								eventstore.WithEventType("org.idp.ldap.added"),
							),
							eventstore.AppendEvent(
								eventstore.WithEventType("org.idp.config.apple.added"),
							),
							eventstore.AppendEvent(
								eventstore.WithEventType("org.idp.saml.added"),
							),
						),
						eventstore.WithLimit(1),
					),
					eventstore.NewFilter(
						eventstore.AppendAggregateFilter(
							"instance",
							eventstore.AppendEvent(
								eventstore.WithEventType("instance.login.policy.idp.added"),
							),
						),
						eventstore.AppendAggregateFilter(
							"org",
							eventstore.AppendEvent(
								eventstore.WithEventType("org.login.policy.idp.added"),
							),
						),
						eventstore.WithLimit(1),
					),
					eventstore.NewFilter(
						eventstore.AppendAggregateFilter(
							"instance",
							eventstore.AppendEvent(
								eventstore.WithEventType("instance.smtp.config.added"),
								eventstore.WithCreatorList(database.NewListNotContains("", "SYSTEM", "<SYSTEM-USER>")),
							),
						),
						eventstore.WithLimit(1),
					),
				},
			},
			want: want{
				query: ``,
				args:  []any{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stmt database.Statement
			writeFilters(&stmt, tt.args.instance, tt.args.filters)

			if got := stmt.String(); got != tt.want.query {
				t.Errorf("unexpected query:\nwant: %q\n got: %q", tt.want.query, got)
			}
			fmt.Println(stmt.Debug())
			if len(stmt.Args()) != len(tt.want.args) {
				t.Errorf("unexpected length of args, want: %d got %d", len(tt.want.args), len(stmt.Args()))
				return
			}
			for i, got := range stmt.Args() {
				if !reflect.DeepEqual(got, tt.want.args[i]) {
					t.Errorf("unexpected arg at position %d: want %v got %v", i, tt.want.args[i], got)
				}
			}
		})
	}
}
