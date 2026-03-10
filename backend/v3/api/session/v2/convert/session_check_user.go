package convert

import (
	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/internal/zerrors"
	session_grpc "github.com/zitadel/zitadel/pkg/grpc/session/v2"
)

func CheckUserGRPCToQueryOpt(checkUser *session_grpc.CheckUser) (*domain.CheckUserType, error) {
	var toReturn *domain.CheckUserType
	if checkUser == nil {
		return toReturn, nil
	}

	switch searchType := checkUser.GetSearch().(type) {
	case *session_grpc.CheckUser_UserId:
		toReturn = &domain.CheckUserType{UserID: searchType.UserId}
	case *session_grpc.CheckUser_LoginName:
		toReturn = &domain.CheckUserType{LoginName: searchType.LoginName}
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "CONV-7B2m0b", "user search %T not implemented", searchType)
	}

	return toReturn, nil
}
