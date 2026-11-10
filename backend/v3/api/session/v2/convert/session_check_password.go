package convert

import (
	"github.com/zitadel/zitadel/backend/v3/domain"
	session_grpc "github.com/zitadel/zitadel/pkg/grpc/session/v2"
)

func CheckPasswordGRPCToDomain(checkPsw *session_grpc.CheckPassword) *domain.CheckPasswordType {
	if checkPsw == nil {
		return nil
	}

	return &domain.CheckPasswordType{
		Password: checkPsw.GetPassword(),
	}
}
