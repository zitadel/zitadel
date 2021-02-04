package grpc

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/labels"
	"github.com/caos/zitadel/operator"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/configuration"
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
	ingressDefinitionSuffix string,
	grpcService string,
	grpcPort uint16,
	dns *configuration.DNS,
	controllerSpecifics map[string]interface{},
	queryIngress core.IngressDefinitionQueryFunc,
	destroyIngress core.IngressDefinitionDestroyFunc,
) (
	operator.QueryFunc,
	operator.DestroyFunc,
	error,
) {
	internalMonitor := monitor.WithField("part", "grpc")

	fullAdminIName := AdminIName + ingressDefinitionSuffix
	fullAuthIName := AuthIName + ingressDefinitionSuffix
	fullMgmtIName := MgmtIName + ingressDefinitionSuffix

	destroyAdminG, err := destroyIngress(namespace, fullAdminIName)
	if err != nil {
		return nil, nil, err
	}
	destroyAuthG, err := destroyIngress(namespace, fullAuthIName)
	if err != nil {
		return nil, nil, err
	}
	destroyMgmtGRPC, err := destroyIngress(namespace, fullMgmtIName)
	if err != nil {
		return nil, nil, err
	}

	destroyers := []operator.DestroyFunc{
		operator.ResourceDestroyToZitadelDestroy(destroyAdminG),
		operator.ResourceDestroyToZitadelDestroy(destroyAuthG),
		operator.ResourceDestroyToZitadelDestroy(destroyMgmtGRPC),
	}

	return func(k8sClient kubernetes.ClientInt, queried map[string]interface{}) (operator.EnsureFunc, error) {
			crd, err := k8sClient.CheckCRD("mappings.getambassador.io")
			if crd == nil || err != nil {
				return func(k8sClient kubernetes.ClientInt) error { return nil }, nil
			}

			apiDomain := dns.Subdomains.API + "." + dns.Domain
			consoleDomain := dns.Subdomains.Console + "." + dns.Domain
			_ = consoleDomain

			cors := &core.CORS{
				Origins:        "*",
				Methods:        "POST, GET, OPTIONS, DELETE, PUT",
				Headers:        "*",
				Credentials:    true,
				ExposedHeaders: "*",
				MaxAge:         "86400",
			}

			queryAdminG, err := queryIngress(
				namespace,
				labels.MustForName(componentLabels, fullAdminIName),
				true,
				apiDomain,
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
				return nil, err
			}

			queryAuthG, err := queryIngress(
				namespace,
				labels.MustForName(componentLabels, fullAuthIName),
				true,
				apiDomain,
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
				return nil, err
			}

			queryMgmtGRPC, err := queryIngress(
				namespace,
				labels.MustForName(componentLabels, fullMgmtIName),
				true,
				apiDomain,
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
