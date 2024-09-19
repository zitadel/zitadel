package management

import (
	"context"

	"github.com/zitadel/zitadel/v2/internal/domain"
	"github.com/zitadel/zitadel/v2/internal/i18n"
	mgmt_pb "github.com/zitadel/zitadel/v2/pkg/grpc/management"
)

func (s *Server) GetSupportedLanguages(context.Context, *mgmt_pb.GetSupportedLanguagesRequest) (*mgmt_pb.GetSupportedLanguagesResponse, error) {
	return &mgmt_pb.GetSupportedLanguagesResponse{Languages: domain.LanguagesToStrings(i18n.SupportedLanguages())}, nil
}
