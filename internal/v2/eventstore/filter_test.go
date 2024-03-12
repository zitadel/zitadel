package eventstore_test

import (
	"testing"

	"github.com/zitadel/zitadel/internal/v2/database"
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
								EventCreatorList(database.NewListNotContains("", "SYSTEM")),
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
								EventCreatorList(database.NewListNotContains("", "SYSTEM")),
							),
						),
					),
					NewFilter(
						AppendAggregateFilter(
							"project",
							AppendEvent(
								EventCreatorList(database.NewListNotContains("", "SYSTEM")),
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
								EventType("instance.idp.config.added"),
							),
							AppendEvent(
								EventType("instance.idp.oauth.added"),
							),
							AppendEvent(
								EventType("instance.idp.oidc.added"),
							),
							AppendEvent(
								EventType("instance.idp.jwt.added"),
							),
							AppendEvent(
								EventType("instance.idp.azure.added"),
							),
							AppendEvent(
								EventType("instance.idp.github.added"),
							),
							AppendEvent(
								EventType("instance.idp.github.enterprise.added"),
							),
							AppendEvent(
								EventType("instance.idp.gitlab.added"),
							),
							AppendEvent(
								EventType("instance.idp.gitlab.selfhosted.added"),
							),
							AppendEvent(
								EventType("instance.idp.google.added"),
							),
							AppendEvent(
								EventType("instance.idp.ldap.added"),
							),
							AppendEvent(
								EventType("instance.idp.config.apple.added"),
							),
							AppendEvent(
								EventType("instance.idp.saml.added"),
							),
						),
						AppendAggregateFilter(
							"org",
							AppendEvent(
								EventType("org.idp.config.added"),
							),
							AppendEvent(
								EventType("org.idp.oauth.added"),
							),
							AppendEvent(
								EventType("org.idp.oidc.added"),
							),
							AppendEvent(
								EventType("org.idp.jwt.added"),
							),
							AppendEvent(
								EventType("org.idp.azure.added"),
							),
							AppendEvent(
								EventType("org.idp.github.added"),
							),
							AppendEvent(
								EventType("org.idp.github.enterprise.added"),
							),
							AppendEvent(
								EventType("org.idp.gitlab.added"),
							),
							AppendEvent(
								EventType("org.idp.gitlab.selfhosted.added"),
							),
							AppendEvent(
								EventType("org.idp.google.added"),
							),
							AppendEvent(
								EventType("org.idp.ldap.added"),
							),
							AppendEvent(
								EventType("org.idp.config.apple.added"),
							),
							AppendEvent(
								EventType("org.idp.saml.added"),
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
								EventCreatorList(database.NewListNotContains("", "SYSTEM", "<SYSTEM-USER>")),
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
