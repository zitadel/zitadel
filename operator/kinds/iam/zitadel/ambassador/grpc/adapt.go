package grpc

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/kubernetes/resources/ambassador/mapping"
	"github.com/caos/zitadel/operator"
	"github.com/caos/zitadel/operator/kinds/iam/zitadel/configuration"
)

func AdaptFunc(
	monitor mntr.Monitor,
	namespace string,
	labels map[string]string,
	grpcURL string,
	dns *configuration.DNS,
) (
	operator.QueryFunc,
	operator.DestroyFunc,
	error,
) {
	internalMonitor := monitor.WithField("part", "grpc")

	adminMName := "admin-grpc-v1"
	authMName := "auth-grpc-v1"
	mgmtMName := "mgmt-grpc-v1"

	destroyAdminG, err := mapping.AdaptFuncToDestroy(namespace, adminMName)
	if err != nil {
		return nil, nil, err
	}
	destroyAuthG, err := mapping.AdaptFuncToDestroy(namespace, authMName)
	if err != nil {
		return nil, nil, err
	}
	destroyMgmtGRPC, err := mapping.AdaptFuncToDestroy(namespace, mgmtMName)
	if err != nil {
		return nil, nil, err
	}

	destroyers := []operator.DestroyFunc{
		operator.ResourceDestroyToZitadelDestroy(destroyAdminG),
		operator.ResourceDestroyToZitadelDestroy(destroyAuthG),
		operator.ResourceDestroyToZitadelDestroy(destroyMgmtGRPC),
	}

	return func(k8sClient *kubernetes.Client, queried map[string]interface{}) (operator.EnsureFunc, error) {
			crd, err := k8sClient.CheckCRD("mappings.getambassador.io")
			if crd == nil || err != nil {
				return func(k8sClient *kubernetes.Client) error { return nil }, nil
			}

			apiDomain := dns.Subdomains.API + "." + dns.Domain
			consoleDomain := dns.Subdomains.Console + "." + dns.Domain
			_ = consoleDomain

			cors := &mapping.CORS{
				Origins:        "*",
				Methods:        "POST, GET, OPTIONS, DELETE, PUT",
				Headers:        "*",
				Credentials:    true,
				ExposedHeaders: "*",
				MaxAge:         "86400",
			}

			queryAdminG, err := mapping.AdaptFuncToEnsure(
				namespace,
				adminMName,
				labels,
				true,
				apiDomain,
				"/caos.zitadel.admin.api.v1.AdminService/",
				"",
				grpcURL,
				"30000",
				"30000",
				cors,
			)
			if err != nil {
				return nil, err
			}

			queryAuthG, err := mapping.AdaptFuncToEnsure(
				namespace,
				authMName,
				labels,
				true,
				apiDomain,
				"/caos.zitadel.auth.api.v1.AuthService/",
				"",
				grpcURL,
				"30000",
				"30000",
				cors,
			)
			if err != nil {
				return nil, err
			}

			queryMgmtGRPC, err := mapping.AdaptFuncToEnsure(
				namespace,
				mgmtMName,
				labels,
				true,
				apiDomain,
				"/caos.zitadel.management.api.v1.ManagementService/",
				"",
				grpcURL,
				"30000",
				"30000",
				cors,
			)
			if err != nil {
				return nil, err
			}

			queriers := []operator.QueryFunc{
				operator.ResourceQueryToZitadelQuery(queryAdminG),
				operator.ResourceQueryToZitadelQuery(queryAuthG),
				operator.ResourceQueryToZitadelQuery(queryMgmtGRPC),
			}

			return operator.QueriersToEnsureFunc(internalMonitor, false, queriers, k8sClient, queried)
		},
		operator.DestroyersToDestroyFunc(internalMonitor, destroyers),
		nil
}
