package convert

import (
	"github.com/zitadel/zitadel/backend/v3/domain"
	session_grpc "github.com/zitadel/zitadel/pkg/grpc/session/v2"
)

func CheckIDPIntentGRPCToDomain(checkIDPIntent *session_grpc.CheckIDPIntent) *domain.CheckIDPIntentType {
	if checkIDPIntent == nil {
		return nil
	}

	return &domain.CheckIDPIntentType{
		ID:    checkIDPIntent.GetIdpIntentId(),
		Token: checkIDPIntent.GetIdpIntentToken(),
	}
}
