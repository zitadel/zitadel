package ui

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/labels"
	"github.com/caos/zitadel/operator"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/configuration"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/ingress/protocol/core"
)

const (
	ConsoleName  = "console-v1"
	AccountsName = "accounts-v1"
)

func AdaptFunc(
	monitor mntr.Monitor,
	componentLabels *labels.Component,
	namespace string,
	ingressDefinitionSuffix string,
	uiService string,
	uiPort uint16,
	dns *configuration.DNS,
	controllerSpecifics map[string]interface{},
	queryIngress core.IngressDefinitionQueryFunc,
	destroyIngress core.IngressDefinitionDestroyFunc,
) (
	operator.QueryFunc,
	operator.DestroyFunc,
	error,
) {
	internalMonitor := monitor.WithField("part", "ui")

	fullConsoleName := ConsoleName + ingressDefinitionSuffix
	fullAccountsName := AccountsName + ingressDefinitionSuffix

	destroyAcc, err := destroyIngress(namespace, fullAccountsName)
	if err != nil {
		return nil, nil, err
	}

	destroyConsole, err := destroyIngress(namespace, fullConsoleName)
	if err != nil {
		return nil, nil, err
	}

	destroyers := []operator.DestroyFunc{
		operator.ResourceDestroyToZitadelDestroy(destroyAcc),
		operator.ResourceDestroyToZitadelDestroy(destroyConsole),
	}

	return func(k8sClient kubernetes.ClientInt, queried map[string]interface{}) (operator.EnsureFunc, error) {
			crd, err := k8sClient.CheckCRD("mappings.getambassador.io")
			if crd == nil || err != nil {
				return func(k8sClient kubernetes.ClientInt) error { return nil }, nil
			}

			accountsDomain := dns.Subdomains.Accounts + "." + dns.Domain
			consoleDomain := dns.Subdomains.Console + "." + dns.Domain

			queryConsole, err := queryIngress(
				namespace,
				labels.MustForName(componentLabels, fullConsoleName),
				false,
				consoleDomain,
				"/",
				"/console/",
				uiService,
				uiPort,
				0,
				0,
				nil,
				controllerSpecifics,
			)
			if err != nil {
				return nil, err
			}

			queryAcc, err := queryIngress(
				namespace,
				labels.MustForName(componentLabels, fullAccountsName),
				false,
				accountsDomain,
				"/",
				"/login/",
				uiService,
				uiPort,
				30000,
				30000,
				nil,
				controllerSpecifics,
			)
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
