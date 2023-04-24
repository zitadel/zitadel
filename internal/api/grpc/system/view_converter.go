package system

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/view/model"
	system_pb "github.com/zitadel/zitadel/pkg/grpc/system"
)

func ViewsToPb(views []*model.View) []*system_pb.View {
	v := make([]*system_pb.View, len(views))
	for i, view := range views {
		v[i] = ViewToPb(view)
	}
	return v
}

func ViewToPb(view *model.View) *system_pb.View {
	return &system_pb.View{
		Database:                 view.Database,
		ViewName:                 view.ViewName,
		LastSuccessfulSpoolerRun: timestamppb.New(view.LastSuccessfulSpoolerRun),
		ProcessedSequence:        view.CurrentSequence,
		EventTimestamp:           timestamppb.New(view.EventTimestamp),
	}
}

func CurrentSequencesToPb(database string, currentSequences *query.CurrentStates) []*system_pb.View {
	v := make([]*system_pb.View, len(currentSequences.CurrentStates))
	for i, currentSequence := range currentSequences.CurrentStates {
		v[i] = CurrentSequenceToPb(database, currentSequence)
	}
	return v
}

func CurrentSequenceToPb(database string, currentSequence *query.CurrentState) *system_pb.View {
	return &system_pb.View{
		Database:                 database,
		ViewName:                 currentSequence.ProjectionName,
		ProcessedSequence:        currentSequence.CurrentSequence,
		LastSuccessfulSpoolerRun: timestamppb.New(currentSequence.EventTimestamp),
	}
}
