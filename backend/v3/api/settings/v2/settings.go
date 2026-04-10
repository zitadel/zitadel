package v2

import (
	"context"

	"connectrpc.com/connect"

	"github.com/zitadel/zitadel/backend/v3/api/settings/v2/convert"
	"github.com/zitadel/zitadel/backend/v3/domain"
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
	linkSettings, err := convert.DomainLinkSettingsModelToGRPCResponse(res)
	if err != nil {
		return nil, err
	}

	return &connect.Response[settings.GetLinkSettingsResponse]{
		Msg: &settings.GetLinkSettingsResponse{
			Settings: linkSettings,
			Source:   0,
		},
	}, nil
}

func SetLinkSettings(ctx context.Context, request *connect.Request[settings.SetLinkSettingsRequest]) (*connect.Response[settings.GetLinkSettingsResponse], error) {
	ls, err := convert.GrpcLinksToDomain(request.Msg.GetLinks())
	if err != nil {
		return nil, err
	}

	cmd := domain.NewSetLinkSettingsCommand(
		request.Msg.Ctx.GetInstance(),
		request.Msg.Ctx.GetOrgId(),
		ls,
	)

	err = domain.Invoke(ctx, cmd)
	if err != nil {
		return nil, err
	}

	// TODO(wim): should a query be invoked to retrieve the return object?
}

func ResetLinkSettings(ctx context.Context, request *connect.Request[settings.ResetLinkSettingsRequest]) (*connect.Response[settings.ResetLinkSettingsResponse], error) {
	cmd := domain.NewResetLinkSettingsCommand(
		request.Msg.Ctx.GetInstance(),
		request.Msg.Ctx.GetOrgId(),
	)

	err := domain.Invoke(ctx, cmd)
	if err != nil {
		return nil, err
	}

	// TODO(wim): should a query be invoked to retrieve the return object?
}
