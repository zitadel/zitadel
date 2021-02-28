package change

import (
	"github.com/caos/zitadel/internal/user/model"
	change_pb "github.com/caos/zitadel/pkg/grpc/change"
	"github.com/caos/zitadel/pkg/grpc/message"
)

func UserChangesToPb(changes []*model.UserChange) []*change_pb.Change {
	c := make([]*change_pb.Change, len(changes))
	for i, change := range changes {
		c[i] = UserChangeToPb(change)
	}
	return c
}

func UserChangeToPb(change *model.UserChange) *change_pb.Change {
	return &change_pb.Change{
		ChangeDate:        change.ChangeDate,
		EventType:         message.NewLocalizedEventType(change.EventType),
		Sequence:          change.Sequence,
		EditorId:          change.ModifierID,
		EditorDisplayName: change.ModifierName,
		// ResourceOwnerId: change.,TODO: resource owner not returned
	}
}
