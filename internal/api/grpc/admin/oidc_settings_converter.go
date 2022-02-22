package admin

import (
	obj_grpc "github.com/caos/zitadel/internal/api/grpc/object"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/query"
	admin_pb "github.com/caos/zitadel/pkg/grpc/admin"
	settings_pb "github.com/caos/zitadel/pkg/grpc/settings"
	"google.golang.org/protobuf/types/known/durationpb"
)

func OIDCSettingsToPb(config *query.OIDCSettings) *settings_pb.OIDCSettings {
	return &settings_pb.OIDCSettings{
		Details:                    obj_grpc.ToViewDetailsPb(config.Sequence, config.CreationDate, config.ChangeDate, config.AggregateID),
		AccessTokenLifetime:        durationpb.New(config.AccessTokenLifetime),
		IdTokenLifetime:            durationpb.New(config.IdTokenLifetime),
		RefreshTokenIdleExpiration: durationpb.New(config.RefreshTokenIdleExpiration),
		RefreshTokenExpiration:     durationpb.New(config.RefreshTokenExpiration),
	}
}

func UpdateOIDCConfigToConfig(req *admin_pb.UpdateOIDCSettingsRequest) *domain.OIDCSettings {
	return &domain.OIDCSettings{
		AccessTokenLifetime:        req.AccessTokenLifetime.AsDuration(),
		IdTokenLifetime:            req.IdTokenLifetime.AsDuration(),
		RefreshTokenIdleExpiration: req.RefreshTokenIdleExpiration.AsDuration(),
		RefreshTokenExpiration:     req.RefreshTokenExpiration.AsDuration(),
	}
}
