package system

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/v2/internal/query"
	system_pb "github.com/zitadel/zitadel/v2/pkg/grpc/system"
)

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
		ProcessedSequence:        currentSequence.Sequence,
		LastSuccessfulSpoolerRun: timestamppb.New(currentSequence.LastRun),
	}
}
