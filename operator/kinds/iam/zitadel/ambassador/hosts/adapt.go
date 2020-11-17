package hosts

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/kubernetes/resources/ambassador/host"
	"github.com/caos/zitadel/operator"
	"github.com/caos/zitadel/operator/kinds/iam/zitadel/configuration"
)

func AdaptFunc(
	monitor mntr.Monitor,
	namespace string,
	labels map[string]string,
	dns *configuration.DNS,
) (
	operator.QueryFunc,
	operator.DestroyFunc,
	error,
) {
	internalMonitor := monitor.WithField("part", "hosts")

	accountsHostName := "accounts"
	apiHostName := "api"
	consoleHostName := "console"
	issuerHostName := "issuer"

	destroyAccounts, err := host.AdaptFuncToDestroy(namespace, accountsHostName)
	if err != nil {
		return nil, nil, err
	}

	destroyAPI, err := host.AdaptFuncToDestroy(namespace, apiHostName)
	if err != nil {
		return nil, nil, err
	}

	destroyConsole, err := host.AdaptFuncToDestroy(namespace, consoleHostName)
	if err != nil {
		return nil, nil, err
	}

	destroyIssuer, err := host.AdaptFuncToDestroy(namespace, issuerHostName)
	if err != nil {
		return nil, nil, err
	}

	destroyers := []operator.DestroyFunc{
		operator.ResourceDestroyToZitadelDestroy(destroyAccounts),
		operator.ResourceDestroyToZitadelDestroy(destroyAPI),
		operator.ResourceDestroyToZitadelDestroy(destroyConsole),
		operator.ResourceDestroyToZitadelDestroy(destroyIssuer),
	}

	return func(k8sClient kubernetes.ClientInt, queried map[string]interface{}) (operator.EnsureFunc, error) {
			crd, err := k8sClient.CheckCRD("hosts.getambassador.io")
			if crd == nil || err != nil {
				return func(k8sClient kubernetes.ClientInt) error { return nil }, nil
			}

			accountsDomain := dns.Subdomains.Accounts + "." + dns.Domain
			apiDomain := dns.Subdomains.API + "." + dns.Domain
			consoleDomain := dns.Subdomains.Console + "." + dns.Domain
			issuerDomain := dns.Subdomains.Issuer + "." + dns.Domain
			originCASecretName := dns.TlsSecret

			accountsSelector := map[string]string{
				"hostname": accountsDomain,
			}
			queryAccounts, err := host.AdaptFuncToEnsure(namespace, accountsHostName, labels, accountsDomain, "none", "", accountsSelector, originCASecretName)
			if err != nil {
				return nil, err
			}

			apiSelector := map[string]string{
				"hostname": apiDomain,
			}
			queryAPI, err := host.AdaptFuncToEnsure(namespace, apiHostName, labels, apiDomain, "none", "", apiSelector, originCASecretName)
			if err != nil {
				return nil, err
			}

			consoleSelector := map[string]string{
				"hostname": consoleDomain,
			}
			queryConsole, err := host.AdaptFuncToEnsure(namespace, consoleHostName, labels, consoleDomain, "none", "", consoleSelector, originCASecretName)
			if err != nil {
				return nil, err
			}

			issuerSelector := map[string]string{
				"hostname": issuerDomain,
			}
			queryIssuer, err := host.AdaptFuncToEnsure(namespace, issuerHostName, labels, issuerDomain, "none", "", issuerSelector, originCASecretName)
			if err != nil {
				return nil, err
			}

			queriers := []operator.QueryFunc{
				operator.ResourceQueryToZitadelQuery(queryAccounts),
				operator.ResourceQueryToZitadelQuery(queryAPI),
				operator.ResourceQueryToZitadelQuery(queryConsole),
				operator.ResourceQueryToZitadelQuery(queryIssuer),
			}

			return operator.QueriersToEnsureFunc(internalMonitor, false, queriers, k8sClient, queried)
		},
		operator.DestroyersToDestroyFunc(internalMonitor, destroyers),
		nil
}
