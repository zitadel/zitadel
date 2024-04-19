package eventstore_test

import (
	"testing"

	. "github.com/zitadel/zitadel/internal/v2/eventstore"
)

func TestNewFilter(t *testing.T) {
	type args struct {
		filters []*Filter
	}
	tests := []struct {
		name string
		args args
		want *Filter
	}{
		{
			name: "milestones filter",
			args: args{
				filters: []*Filter{
					NewFilter(
						AppendAggregateFilter(
							"instance",
							AppendEvent(
								EventType("instance.added"),
							),
						),
						FilterPagination(
							Limit(1),
						),
					),
					NewFilter(
						AppendAggregateFilter(
							"instance",
							AppendEvent(
								EventType("instance.removed"),
							),
						),
						FilterPagination(
							Limit(1),
						),
					),
					NewFilter(
						AppendAggregateFilter(
							"instance",
							AppendEvent(
								EventType("instance.domain.primary.set"),
								EventCreatorsNotContains("", "SYSTEM"),
							),
						),
						FilterPagination(
							Limit(1),
						),
					),
					NewFilter(
						AppendAggregateFilter(
							"project",
							AppendEvent(
								EventType("project.added"),
								EventCreatorsNotContains("", "SYSTEM"),
							),
						),
					),
					NewFilter(
						AppendAggregateFilter(
							"project",
							AppendEvent(
								EventCreatorsNotContains("", "SYSTEM"),
								EventType("project.application.added"),
							),
						),
					),
					NewFilter(
						AppendAggregateFilter(
							"user",
							AppendEvent(
								EventType("user.token.added"),
							),
						),
						FilterPagination(
							// used because we need to check for first login and an app which is not console
							PositionGreater(12, 4),
						),
					),
					NewFilter(
						AppendAggregateFilter(
							"instance",
							AppendEvent(
								AppendEventTypes(
									"instance.idp.config.added",
									"instance.idp.oauth.added",
									"instance.idp.oidc.added",
									"instance.idp.jwt.added",
									"instance.idp.azure.added",
									"instance.idp.github.added",
									"instance.idp.github.enterprise.added",
									"instance.idp.gitlab.added",
									"instance.idp.gitlab.selfhosted.added",
									"instance.idp.google.added",
									"instance.idp.ldap.added",
									"instance.idp.config.apple.added",
									"instance.idp.saml.added",
								),
							),
						),
						AppendAggregateFilter(
							"org",
							AppendEvent(
								AppendEventTypes(
									"org.idp.config.added",
									"org.idp.oauth.added",
									"org.idp.oidc.added",
									"org.idp.jwt.added",
									"org.idp.azure.added",
									"org.idp.github.added",
									"org.idp.github.enterprise.added",
									"org.idp.gitlab.added",
									"org.idp.gitlab.selfhosted.added",
									"org.idp.google.added",
									"org.idp.ldap.added",
									"org.idp.config.apple.added",
									"org.idp.saml.added",
								),
							),
						),
						FilterPagination(
							Limit(1),
						),
					),
					NewFilter(
						AppendAggregateFilter(
							"instance",
							AppendEvent(
								EventType("instance.login.policy.idp.added"),
							),
						),
						AppendAggregateFilter(
							"org",
							AppendEvent(
								EventType("org.login.policy.idp.added"),
							),
						),
						FilterPagination(
							Limit(1),
						),
					),
					NewFilter(
						AppendAggregateFilter(
							"instance",
							AppendEvent(
								EventType("instance.smtp.config.added"),
								EventCreatorsNotContains("", "SYSTEM", "<SYSTEM-USER>"),
							),
						),
						FilterPagination(
							Limit(1),
						),
					),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// got := NewFilter(tt.args.instance, tt.args.filters...)
			// log.Println(got)
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("NewFilter() = %v, want %v", got, tt.want)
			// }
		})
	}
}
