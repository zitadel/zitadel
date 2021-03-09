package admin

import (
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/view/model"
	admin_pb "github.com/caos/zitadel/pkg/grpc/admin"
	"github.com/golang/protobuf/ptypes"
)

func ViewsToPb(views []*model.View) []*admin_pb.View {
	v := make([]*admin_pb.View, len(views))
	for i, view := range views {
		v[i] = ViewToPb(view)
	}
	return v
}

func ViewToPb(view *model.View) *admin_pb.View {
	lastSuccessfulSpoolerRun, err := ptypes.TimestampProto(view.LastSuccessfulSpoolerRun)
	logging.Log("ADMIN-4zs01").OnError(err).Debug("unable to parse last successful spooler run")

	eventTs, err := ptypes.TimestampProto(view.EventTimestamp)
	logging.Log("ADMIN-q2Wzj").OnError(err).Debug("unable to parse event timestamp")

	return &admin_pb.View{
		Database:                 view.Database,
		ViewName:                 view.ViewName,
		LastSuccessfulSpoolerRun: lastSuccessfulSpoolerRun,
		ProcessedSequence:        view.CurrentSequence,
		EventTimestamp:           eventTs,
	}
}
