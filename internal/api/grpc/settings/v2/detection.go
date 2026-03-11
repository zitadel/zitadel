package settings

import (
"context"

"connectrpc.com/connect"

"github.com/zitadel/zitadel/internal/api/grpc/object/v2"
"github.com/zitadel/zitadel/pkg/grpc/settings/v2"
)

func (s *Server) GetDetectionSettings(ctx context.Context, _ *connect.Request[settings.GetDetectionSettingsRequest]) (*connect.Response[settings.GetDetectionSettingsResponse], error) {
current, err := s.command.GetEffectiveDetectionSettings(ctx)
if err != nil {
return nil, err
}
return connect.NewResponse(&settings.GetDetectionSettingsResponse{Settings: detectionSettingsToPb(current)}), nil
}

func (s *Server) SetDetectionSettings(ctx context.Context, req *connect.Request[settings.SetDetectionSettingsRequest]) (*connect.Response[settings.SetDetectionSettingsResponse], error) {
details, err := s.command.SetDetectionSettings(ctx, detectionSettingsToCommand(req.Msg.GetSettings()))
if err != nil {
return nil, err
}
return connect.NewResponse(&settings.SetDetectionSettingsResponse{Details: object.DomainToDetailsPb(details)}), nil
}
