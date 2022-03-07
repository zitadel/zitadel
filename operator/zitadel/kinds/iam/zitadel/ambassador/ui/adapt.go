package ui

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/kubernetes/resources/ambassador/mapping"
	"github.com/caos/orbos/pkg/labels"
	"github.com/caos/zitadel/operator"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/configuration"
)

const (
	ConsoleName  = "console-v1"
	AccountsName = "accounts-v1"
)

func AdaptFunc(
	monitor mntr.Monitor,
	componentLabels *labels.Component,
	namespace string,
	uiURL string,
	dns *configuration.DNS,
) (
	operator.QueryFunc,
	operator.DestroyFunc,
	error,
) {
	internalMonitor := monitor.WithField("part", "ui")

	destroyAcc, err := mapping.AdaptFuncToDestroy(namespace, AccountsName)
	if err != nil {
		return nil, nil, err
	}

	destroyConsole, err := mapping.AdaptFuncToDestroy(namespace, ConsoleName)
	if err != nil {
		return nil, nil, err
	}

	destroyers := []operator.DestroyFunc{
		operator.ResourceDestroyToZitadelDestroy(destroyAcc),
		operator.ResourceDestroyToZitadelDestroy(destroyConsole),
	}

	return func(k8sClient kubernetes.ClientInt, queried map[string]interface{}) (operator.EnsureFunc, error) {
			_, found, err := k8sClient.CheckCRD("mappings.getambassador.io")
			if err != nil || !found {
				return func(k8sClient kubernetes.ClientInt) error { return nil }, err
			}

			accountsDomain := dns.Subdomains.Accounts + "." + dns.Domain
			consoleDomain := dns.Subdomains.Console + "." + dns.Domain

			queryConsole, err := mapping.AdaptFuncToEnsure(&mapping.Arguments{
				Monitor:          internalMonitor,
				Namespace:        namespace,
				ID:               labels.MustForName(componentLabels, ConsoleName),
				GRPC:             false,
				Host:             consoleDomain,
				Prefix:           "/",
				Rewrite:          "/console/",
				Service:          uiURL,
				TimeoutMS:        0,
				ConnectTimeoutMS: 0,
				CORS:             nil,
			})
			if err != nil {
				return nil, err
			}

			queryAcc, err := mapping.AdaptFuncToEnsure(&mapping.Arguments{
				Monitor:          internalMonitor,
				Namespace:        namespace,
				ID:               labels.MustForName(componentLabels, AccountsName),
				GRPC:             false,
				Host:             accountsDomain,
				Prefix:           "/",
				Rewrite:          "/login/",
				Service:          uiURL,
				TimeoutMS:        30000,
				ConnectTimeoutMS: 30000,
				CORS:             nil,
			})
			if err != nil {
				return nil, err
			}

			queriers := []operator.QueryFunc{
				operator.ResourceQueryToZitadelQuery(queryConsole),
				operator.ResourceQueryToZitadelQuery(queryAcc),
			}

			return operator.QueriersToEnsureFunc(internalMonitor, false, queriers, k8sClient, queried)
		},
		operator.DestroyersToDestroyFunc(internalMonitor, destroyers),
		nil
}
