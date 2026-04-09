package convert

import (
	"github.com/zitadel/zitadel/backend/v3/domain"
	session_grpc "github.com/zitadel/zitadel/pkg/grpc/session/v2"
)

func CheckRecoveryCodeGRPCToDomain(checkRecoveryCode *session_grpc.CheckRecoveryCode) *domain.CheckTypeRecoveryCode {
	if checkRecoveryCode == nil {
		return nil
	}
	return &domain.CheckTypeRecoveryCode{
		RecoveryCode: checkRecoveryCode.GetCode(),
	}
}
