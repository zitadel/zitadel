package change

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/config/systemdefaults"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	change_pb "github.com/zitadel/zitadel/pkg/grpc/change"
	"github.com/zitadel/zitadel/pkg/grpc/message"
)

// ChangeQueryToModel returns the limit, the sequence and the sorting order of the change query.
// The limit is defaulted and checked by [systemdefaults.SystemDefaults.V1QueryLimit].
func ChangeQueryToModel(defaults systemdefaults.SystemDefaults, query *change_pb.ChangeQuery) (limit, sequence uint64, asc bool, err error) {
	limit, err = defaults.V1QueryLimit(uint64(query.GetLimit()))
	if err != nil {
		return 0, 0, false, err
	}
	return limit, query.GetSequence(), query.GetAsc(), nil
}

func EventsToChangesPb(changes []*query.Event, assetAPIPrefix string) []*change_pb.Change {
	c := make([]*change_pb.Change, len(changes))
	for i, change := range changes {
		c[i] = EventToChangePb(change, assetAPIPrefix)
	}
	return c
}

func EventToChangePb(change *query.Event, assetAPIPrefix string) *change_pb.Change {
	return &change_pb.Change{
		ChangeDate:               timestamppb.New(change.CreationDate),
		EventType:                message.NewLocalizedEventType(change.Type),
		Sequence:                 change.Sequence,
		EditorId:                 change.Editor.ID,
		EditorDisplayName:        change.Editor.DisplayName,
		EditorPreferredLoginName: change.Editor.PreferedLoginName,
		EditorAvatarUrl:          domain.AvatarURL(assetAPIPrefix, change.Aggregate.ResourceOwner, change.Editor.AvatarKey),
		ResourceOwnerId:          change.Aggregate.ResourceOwner,
	}
}
