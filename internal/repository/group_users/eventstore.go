package groupusers

import (
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/group"
)

func init() {
	eventstore.RegisterFilterEventMapper(group.AggregateType, AddedEventType, GroupUserAddedEventMapper)
	eventstore.RegisterFilterEventMapper(group.AggregateType, ChangedEventType, GroupUserChangedEventMapper)
	eventstore.RegisterFilterEventMapper(group.AggregateType, RemovedEventType, GroupUserRemovedEventMapper)
	eventstore.RegisterFilterEventMapper(group.AggregateType, CascadeRemovedEventType, GroupUserCascadeRemovedEventMapper)
}
