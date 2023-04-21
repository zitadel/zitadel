package admin

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/view/model"
	admin_pb "github.com/zitadel/zitadel/pkg/grpc/admin"
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

func CurrentSequencesToPb(database string, currentSequences *query.CurrentStates) []*admin_pb.View {
	v := make([]*admin_pb.View, len(currentSequences.CurrentStates))
	for i, currentSequence := range currentSequences.CurrentStates {
		v[i] = CurrentSequenceToPb(database, currentSequence)
	}
	return v
}

func CurrentSequenceToPb(database string, currentSequence *query.CurrentSequence) *admin_pb.View {
	return &admin_pb.View{
		Database:                 database,
		ViewName:                 currentSequence.ProjectionName,
		ProcessedSequence:        currentSequence.CurrentSequence,
		LastSuccessfulSpoolerRun: timestamppb.New(currentSequence.Timestamp),
	}
}
