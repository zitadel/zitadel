package convert

import (
	"encoding/json"

	session_grpc "github.com/zitadel/zitadel/pkg/grpc/session/v2"
)

func CheckPasskeyGRPCToDomain(checkPsw *session_grpc.CheckWebAuthN) ([]byte, error) {
	if checkPsw == nil || checkPsw.GetCredentialAssertionData() == nil {
		return nil, nil
	}

	return json.Marshal(checkPsw.GetCredentialAssertionData())
}
