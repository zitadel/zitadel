package admin

import (
	"github.com/zitadel/zitadel/internal/query"
	settings_pb "github.com/zitadel/zitadel/pkg/grpc/settings"
)

func SMTPConfigsToPb(configs []*query.SMSConfig) []*settings_pb.SMSProvider {
	c := make([]*settings_pb.SMSProvider, len(configs))
	for i, config := range configs {
		c[i] = SMTPConfigToProviderPb(config)
	}
	return c
}
