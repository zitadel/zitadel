package change

import (
	org_model "github.com/caos/zitadel/internal/org/model"
	proj_model "github.com/caos/zitadel/internal/project/model"
	user_model "github.com/caos/zitadel/internal/user/model"
	change_pb "github.com/caos/zitadel/pkg/grpc/change"
	"github.com/caos/zitadel/pkg/grpc/message"
)

func UserChangesToPb(changes []*user_model.UserChange) []*change_pb.Change {
	c := make([]*change_pb.Change, len(changes))
	for i, change := range changes {
		c[i] = UserChangeToPb(change)
	}
	return c
}

func UserChangeToPb(change *user_model.UserChange) *change_pb.Change {
	return &change_pb.Change{
		ChangeDate:        change.ChangeDate,
		EventType:         message.NewLocalizedEventType(change.EventType),
		Sequence:          change.Sequence,
		EditorId:          change.ModifierID,
		EditorDisplayName: change.ModifierName,
		// ResourceOwnerId: change.,TODO: resource owner not returned
	}
}

func OrgChangesToPb(changes []*org_model.OrgChange) []*change_pb.Change {
	c := make([]*change_pb.Change, len(changes))
	for i, change := range changes {
		c[i] = OrgChangeToPb(change)
	}
	return c
}

func OrgChangeToPb(change *org_model.OrgChange) *change_pb.Change {
	return &change_pb.Change{
		ChangeDate:        change.ChangeDate,
		EventType:         message.NewLocalizedEventType(change.EventType),
		Sequence:          change.Sequence,
		EditorId:          change.ModifierId,
		EditorDisplayName: change.ModifierName,
		// ResourceOwnerId: change.,TODO: resource owner not returned
	}
}

func ProjectChangesToPb(changes []*proj_model.ProjectChange) []*change_pb.Change {
	c := make([]*change_pb.Change, len(changes))
	for i, change := range changes {
		c[i] = ProjectChangeToPb(change)
	}
	return c
}

func ProjectChangeToPb(change *proj_model.ProjectChange) *change_pb.Change {
	return &change_pb.Change{
		ChangeDate:        change.ChangeDate,
		EventType:         message.NewLocalizedEventType(change.EventType),
		Sequence:          change.Sequence,
		EditorId:          change.ModifierId,
		EditorDisplayName: change.ModifierName,
		// ResourceOwnerId: change.,TODO: resource owner not returned
	}
}

func AppChangesToPb(changes []*proj_model.ApplicationChange) []*change_pb.Change {
	c := make([]*change_pb.Change, len(changes))
	for i, change := range changes {
		c[i] = AppChangeToPb(change)
	}
	return c
}

func AppChangeToPb(change *proj_model.ApplicationChange) *change_pb.Change {
	return &change_pb.Change{
		ChangeDate:        change.ChangeDate,
		EventType:         message.NewLocalizedEventType(change.EventType),
		Sequence:          change.Sequence,
		EditorId:          change.ModifierId,
		EditorDisplayName: change.ModifierName,
		// ResourceOwnerId: change.,TODO: resource owner not returned
	}
}
