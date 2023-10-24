package admin

import (
	"github.com/zitadel/zitadel/internal/api/grpc/object"
	"github.com/zitadel/zitadel/internal/query"
	admin_pb "github.com/zitadel/zitadel/pkg/grpc/admin"
	settings_pb "github.com/zitadel/zitadel/pkg/grpc/settings"
)

func listSMTPConfigsToModel(req *admin_pb.ListSMTPConfigsRequest) (*query.SMTPConfigsSearchQueries, error) {
	offset, limit, asc := object.ListQueryToModel(req.Query)
	return &query.SMTPConfigsSearchQueries{
		SearchRequest: query.SearchRequest{
			Offset: offset,
			Limit:  limit,
			Asc:    asc,
		},
	}, nil
}

func SMTPConfigToProviderPb(config *query.SMTPConfig) *settings_pb.SMTPConfig {
	return &settings_pb.SMTPConfig{
		Details:       object.ToViewDetailsPb(config.Sequence, config.CreationDate, config.ChangeDate, config.ResourceOwner),
		Id:            config.ID,
		Tls:           config.TLS,
		Host:          config.Host,
		User:          config.User,
		ProviderType:  settings_pb.SMTPProviderType(config.ProviderType),
		State:         settings_pb.SMTPConfigState(config.State),
		SenderAddress: config.SenderAddress,
		SenderName:    config.SenderName,
	}
}

func SMTPConfigsToPb(configs []*query.SMTPConfig) []*settings_pb.SMTPConfig {
	c := make([]*settings_pb.SMTPConfig, len(configs))
	for i, config := range configs {
		c[i] = SMTPConfigToProviderPb(config)
	}
	return c
}
