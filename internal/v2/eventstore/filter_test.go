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
								WithEventType("instance.added"),
							),
						),
						WithLimit(1),
					),
					NewFilter(
						AppendAggregateFilter(
							"instance",
							AppendEvent(
								WithEventType("instance.removed"),
							),
						),
						WithLimit(1),
					),
					NewFilter(
						AppendAggregateFilter(
							"instance",
							AppendEvent(
								WithEventType("instance.domain.primary.set"),
								WithCreatorList(database.NewListNotContains("", "SYSTEM")),
							),
						),
						WithLimit(1),
					),
					NewFilter(
						AppendAggregateFilter(
							"project",
							AppendEvent(
								WithEventType("project.added"),
								WithCreatorList(database.NewListNotContains("", "SYSTEM")),
							),
						),
					),
					NewFilter(
						AppendAggregateFilter(
							"project",
							AppendEvent(
								WithCreatorList(database.NewListNotContains("", "SYSTEM")),
								WithEventType("project.application.added"),
							),
						),
					),
					NewFilter(
						AppendAggregateFilter(
							"user",
							AppendEvent(
								WithEventType("user.token.added"),
							),
						),
						// used because we need to check for first login and an app which is not console
						WithPosition(database.NewNumberAtLeast(12), database.NewNumberGreater(4)),
					),
					NewFilter(
						AppendAggregateFilter(
							"instance",
							AppendEvent(
								WithEventType("instance.idp.config.added"),
							),
							AppendEvent(
								WithEventType("instance.idp.oauth.added"),
							),
							AppendEvent(
								WithEventType("instance.idp.oidc.added"),
							),
							AppendEvent(
								WithEventType("instance.idp.jwt.added"),
							),
							AppendEvent(
								WithEventType("instance.idp.azure.added"),
							),
							AppendEvent(
								WithEventType("instance.idp.github.added"),
							),
							AppendEvent(
								WithEventType("instance.idp.github.enterprise.added"),
							),
							AppendEvent(
								WithEventType("instance.idp.gitlab.added"),
							),
							AppendEvent(
								WithEventType("instance.idp.gitlab.selfhosted.added"),
							),
							AppendEvent(
								WithEventType("instance.idp.google.added"),
							),
							AppendEvent(
								WithEventType("instance.idp.ldap.added"),
							),
							AppendEvent(
								WithEventType("instance.idp.config.apple.added"),
							),
							AppendEvent(
								WithEventType("instance.idp.saml.added"),
							),
						),
						AppendAggregateFilter(
							"org",
							AppendEvent(
								WithEventType("org.idp.config.added"),
							),
							AppendEvent(
								WithEventType("org.idp.oauth.added"),
							),
							AppendEvent(
								WithEventType("org.idp.oidc.added"),
							),
							AppendEvent(
								WithEventType("org.idp.jwt.added"),
							),
							AppendEvent(
								WithEventType("org.idp.azure.added"),
							),
							AppendEvent(
								WithEventType("org.idp.github.added"),
							),
							AppendEvent(
								WithEventType("org.idp.github.enterprise.added"),
							),
							AppendEvent(
								WithEventType("org.idp.gitlab.added"),
							),
							AppendEvent(
								WithEventType("org.idp.gitlab.selfhosted.added"),
							),
							AppendEvent(
								WithEventType("org.idp.google.added"),
							),
							AppendEvent(
								WithEventType("org.idp.ldap.added"),
							),
							AppendEvent(
								WithEventType("org.idp.config.apple.added"),
							),
							AppendEvent(
								WithEventType("org.idp.saml.added"),
							),
						),
						WithLimit(1),
					),
					NewFilter(
						AppendAggregateFilter(
							"instance",
							AppendEvent(
								WithEventType("instance.login.policy.idp.added"),
							),
						),
						AppendAggregateFilter(
							"org",
							AppendEvent(
								WithEventType("org.login.policy.idp.added"),
							),
						),
						WithLimit(1),
					),
					NewFilter(
						AppendAggregateFilter(
							"instance",
							AppendEvent(
								WithEventType("instance.smtp.config.added"),
								WithCreatorList(database.NewListNotContains("", "SYSTEM", "<SYSTEM-USER>")),
							),
						),
						WithLimit(1),
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
