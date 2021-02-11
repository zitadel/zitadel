package grpc

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/labels"
	"github.com/caos/zitadel/operator"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/ingress/protocol/core"
)

const (
	AdminIName = "admin-grpc-v1"
	AuthIName  = "auth-grpc-v1"
	MgmtIName  = "mgmt-grpc-v1"
)

func AdaptFunc(
	monitor mntr.Monitor,
	componentLabels *labels.Component,
	namespace string,
	grpcService string,
	grpcPort uint16,
	controllerSpecifics map[string]interface{},
	originCASecretName string,
	apiAdapter core.PathAdapter,
) (
	operator.QueryFunc,
	operator.DestroyFunc,
	error,
) {
	internalMonitor := monitor.WithField("part", "grpc")

	cors := &core.CORS{
		Origins:        "*",
		Methods:        "POST, GET, OPTIONS, DELETE, PUT",
		Headers:        "*",
		Credentials:    true,
		ExposedHeaders: "*",
		MaxAge:         "86400",
	}

	queryAdminG, destroyAdminG, err := apiAdapter(
		monitor,
		namespace,
		labels.MustForName(componentLabels, AdminIName),
		true,
		originCASecretName,
		"/caos.zitadel.admin.api.v1.AdminService/",
		"/caos.zitadel.admin.api.v1.AdminService/",
		grpcService,
		grpcPort,
		30000,
		30000,
		cors,
		controllerSpecifics,
	)
	if err != nil {
		return nil, nil, err
	}

	queryAuthG, destroyAuthG, err := apiAdapter(
		monitor,
		namespace,
		labels.MustForName(componentLabels, AuthIName),
		true,
		originCASecretName,
		"/caos.zitadel.auth.api.v1.AuthService/",
		"/caos.zitadel.auth.api.v1.AuthService/",
		grpcService,
		grpcPort,
		30000,
		30000,
		cors,
		controllerSpecifics,
	)
	if err != nil {
		return nil, nil, err
	}

	queryMgmtGRPC, destroyMgmtGRPC, err := apiAdapter(
		monitor,
		namespace,
		labels.MustForName(componentLabels, MgmtIName),
		true,
		originCASecretName,
		"/caos.zitadel.management.api.v1.ManagementService/",
		"/caos.zitadel.management.api.v1.ManagementService/",
		grpcService,
		grpcPort,
		30000,
		30000,
		cors,
		controllerSpecifics,
	)
	if err != nil {
		return nil, nil, err
	}

	return func(k8sClient kubernetes.ClientInt, queried map[string]interface{}) (operator.EnsureFunc, error) {
			return operator.QueriersToEnsureFunc(internalMonitor, false, []operator.QueryFunc{
				queryAdminG,
				queryAuthG,
				queryMgmtGRPC,
			}, k8sClient, queried)
		},
		operator.DestroyersToDestroyFunc(internalMonitor, []operator.DestroyFunc{
			destroyAdminG,
			destroyAuthG,
			destroyMgmtGRPC,
		}),
		nil
}
