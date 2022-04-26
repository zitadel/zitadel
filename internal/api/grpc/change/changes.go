package change

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/domain"

	"github.com/zitadel/zitadel/internal/query"
	change_pb "github.com/zitadel/zitadel/pkg/grpc/change"
	"github.com/zitadel/zitadel/pkg/grpc/message"
)

func ChangeQueryToQuery(query *change_pb.ChangeQuery) (sequence uint64, limit uint64, asc bool) {
	if query == nil {
		return 0, 0, false
	}
	return query.Sequence, uint64(query.Limit), query.Asc
}

func ChangesToPb(changes []*query.Change, assetAPIPrefix string) []*change_pb.Change {
	c := make([]*change_pb.Change, len(changes))
	for i, change := range changes {
		c[i] = ChangeToPb(change, assetAPIPrefix)
	}
	return c
}

func ChangeToPb(change *query.Change, assetAPIPrefix string) *change_pb.Change {
	return &change_pb.Change{
		ChangeDate:               timestamppb.New(change.ChangeDate),
		EventType:                message.NewLocalizedEventType(change.EventType),
		Sequence:                 change.Sequence,
		EditorId:                 change.ModifierId,
		EditorDisplayName:        change.ModifierName,
		EditorPreferredLoginName: change.ModifierLoginName,
		EditorAvatarUrl:          domain.AvatarURL(assetAPIPrefix, change.ModifierResourceOwner, change.ModifierAvatarKey),
		ResourceOwnerId:          change.ResourceOwner,
	}
}
