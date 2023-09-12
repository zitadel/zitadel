package admin

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/query"
	admin_pb "github.com/zitadel/zitadel/pkg/grpc/admin"
)

func CurrentSequencesToPb(database string, currentSequences *query.CurrentStates) []*admin_pb.View {
	v := make([]*admin_pb.View, len(currentSequences.CurrentStates))
	for i, currentSequence := range currentSequences.CurrentStates {
		v[i] = CurrentSequenceToPb(database, currentSequence)
	}
	return v
}

func CurrentSequenceToPb(database string, currentSequence *query.CurrentState) *admin_pb.View {
	return &admin_pb.View{
		Database:                 database,
		ViewName:                 currentSequence.ProjectionName,
		ProcessedSequence:        currentSequence.Sequence,
		LastSuccessfulSpoolerRun: timestamppb.New(currentSequence.LastRun),
		EventTimestamp:           timestamppb.New(currentSequence.EventCreatedAt),
	}
}
