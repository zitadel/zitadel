package management

import (
	svcacc_model "github.com/caos/zitadel/internal/service_account/model"
	"github.com/caos/zitadel/pkg/grpc/management"
)

func createServiceAccountToModel(account *management.CreateServiceAccountRequest) *svcacc_model.ServiceAccount {
	return nil
}

func updateServiceAccountToModel(account *management.UpdateServiceAccountRequest) *svcacc_model.ServiceAccount {
	return nil
}

func serviceAccountFromModel(account *svcacc_model.ServiceAccount) *management.ServiceAccountResponse {
	return nil
}
