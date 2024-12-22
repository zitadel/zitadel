package instance

import (
	"context"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/pkg/grpc/object"
	system_pb "github.com/zitadel/zitadel/pkg/grpc/system"
)

type AddInstanceRequest struct {
	system_pb.AddInstanceRequest

	ID        string
	CreatedAt time.Time
}

type AddInstanceResponse struct {
	system_pb.AddInstanceResponse
}

func (bl *BusinessLogic) AddInstance(ctx context.Context, request *AddInstanceRequest) (_ *AddInstanceResponse, err error) {
	tx, err := bl.client.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = tx.End(ctx, err)
	}()

	request.ID = time.Since(time.Time{}).String()

	err = bl.storage.WriteInstanceAdded(ctx, tx, request)
	if err != nil {
		return nil, err
	}

	return &AddInstanceResponse{
		AddInstanceResponse: system_pb.AddInstanceResponse{
			InstanceId: request.ID,
			Details: &object.ObjectDetails{
				Sequence:      1,
				CreationDate:  timestamppb.New(request.CreatedAt),
				ChangeDate:    timestamppb.New(request.CreatedAt),
				ResourceOwner: request.ID,
			},
		},
	}, nil
}
