package object

import (
	"time"

	"github.com/caos/logging"
	object_pb "github.com/caos/zitadel/pkg/grpc/object"
	"github.com/golang/protobuf/ptypes"
)

func ToDetailsPb(
	sequence uint64,
	creationDate,
	changeDate time.Time,
	resourceOwner string,
) *object_pb.ObjectDetails {
	creationDatePb, err := ptypes.TimestampProto(creationDate)
	logging.Log("ADMIN-yzma4").OnError(err).Debug("unable to parse creation date")
	changeDatePb, err := ptypes.TimestampProto(changeDate)
	logging.Log("ADMIN-NTgjY").OnError(err).Debug("unable to parse change date")

	return &object_pb.ObjectDetails{
		Sequence:      sequence,
		CreationDate:  creationDatePb,
		ChangeDate:    changeDatePb,
		ResourceOwner: resourceOwner,
	}
}
