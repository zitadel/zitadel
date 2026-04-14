package v2

import (
	"context"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/backend/v3/api/settings/v2/convert"
	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/pkg/grpc/settings/v2"
)

func GetLinkSettings(ctx context.Context, request *connect.Request[settings.GetLinkSettingsRequest]) (*connect.Response[settings.GetLinkSettingsResponse], error) {
	q := domain.NewGetLinkSettingsQuery(
		request.Msg.Ctx.GetInstance(),
		request.Msg.Ctx.GetOrgId(),
	)

	err := domain.Invoke(ctx, q)
	if err != nil {
		return nil, err
	}

	res := q.Result()

	ls, err := convert.DomainLinksModelToGRPCResponse(res.Links)
	if err != nil {
		return nil, err
	}

	source, err := convert.DomainSettingsSourceToSourceToGrpc(res.Source)
	if err != nil {
		return nil, err
	}

	return &connect.Response[settings.GetLinkSettingsResponse]{
		Msg: &settings.GetLinkSettingsResponse{
			Settings: &settings.LinkSettings{Links: ls},
			Source:   source,
		},
	}, nil
}

func SetLinkSettings(ctx context.Context, request *connect.Request[settings.SetLinkSettingsRequest]) (*connect.Response[settings.SetLinkSettingsResponse], error) {
	ls, err := convert.GrpcLinksToDomain(request.Msg.GetLinks())
	if err != nil {
		return nil, err
	}

	cmd := domain.NewSetLinkSettingsCommand(
		authz.GetInstance(ctx).InstanceID(),
		request.Msg.Ctx.GetOrgId(),
		ls,
	)

	err = domain.Invoke(ctx, cmd)
	if err != nil {
		return nil, err
	}

	changeDate := cmd.Result()

	return &connect.Response[settings.SetLinkSettingsResponse]{
		Msg: &settings.SetLinkSettingsResponse{
			ChangeDate: timestamppb.New(changeDate),
			Settings: &settings.LinkSettings{
				Links: request.Msg.GetLinks(),
			},
		},
	}, nil
}

func ResetLinkSettings(ctx context.Context, request *connect.Request[settings.ResetLinkSettingsRequest]) (*connect.Response[settings.ResetLinkSettingsResponse], error) {
	cmd := domain.NewResetLinkSettingsCommand(
		authz.GetInstance(ctx).InstanceID(),
		request.Msg.Ctx.GetOrgId(),
	)

	err := domain.Invoke(ctx, cmd)
	if err != nil {
		return nil, err
	}

	res := cmd.Result()

	links, err := convert.DomainLinksModelToGRPCResponse(res.Links)
	if err != nil {
		return nil, err
	}

	return &connect.Response[settings.ResetLinkSettingsResponse]{
		Msg: &settings.ResetLinkSettingsResponse{
			ChangeDate: timestamppb.New(res.ChangeTime),
			Settings:   &settings.LinkSettings{Links: links},
		},
	}, nil
}
