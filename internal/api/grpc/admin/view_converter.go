package admin

import (
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/view/model"
	admin_pb "github.com/zitadel/zitadel/pkg/grpc/admin"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ViewsToPb(views []*model.View) []*admin_pb.View {
	v := make([]*admin_pb.View, len(views))
	for i, view := range views {
		v[i] = ViewToPb(view)
	}
	return v
}

func ViewToPb(view *model.View) *admin_pb.View {
	return &admin_pb.View{
		Database:                 view.Database,
		ViewName:                 view.ViewName,
		LastSuccessfulSpoolerRun: timestamppb.New(view.LastSuccessfulSpoolerRun),
		ProcessedSequence:        view.CurrentSequence,
		EventTimestamp:           timestamppb.New(view.EventTimestamp),
	}
}

func CurrentSequencesToPb(currentSequences *query.CurrentSequences) []*admin_pb.View {
	v := make([]*admin_pb.View, len(currentSequences.CurrentSequences))
	for i, currentSequence := range currentSequences.CurrentSequences {
		v[i] = CurrentSequenceToPb(currentSequence)
	}
	return v
}

func CurrentSequenceToPb(currentSequence *query.CurrentSequence) *admin_pb.View {
	return &admin_pb.View{
		Database:          "zitadel",
		ViewName:          currentSequence.ProjectionName,
		ProcessedSequence: currentSequence.CurrentSequence,
		EventTimestamp:    timestamppb.New(currentSequence.Timestamp),
	}
}
