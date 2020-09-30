package http

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
	httpUrl string,
	dns *configuration.DNS,
) (
	operator.QueryFunc,
	operator.DestroyFunc,
	error,
) {
	internalMonitor := monitor.WithField("part", "http")

	adminRName := "admin-rest-v1"
	mgmtName := "mgmt-v1"
	oauthName := "oauth-v1"
	authRName := "auth-rest-v1"
	authorizeName := "authorize-v1"
	endsessionName := "endsession-v1"
	issuerName := "issuer-v1"

	destroyAdminR, err := mapping.AdaptFuncToDestroy(namespace, adminRName)
	if err != nil {
		return nil, nil, err
	}

	destroyMgmtRest, err := mapping.AdaptFuncToDestroy(namespace, mgmtName)
	if err != nil {
		return nil, nil, err
	}

	destroyOAuthv2, err := mapping.AdaptFuncToDestroy(namespace, oauthName)
	if err != nil {
		return nil, nil, err
	}

	destroyAuthR, err := mapping.AdaptFuncToDestroy(namespace, authRName)
	if err != nil {
		return nil, nil, err
	}

	destroyAuthorize, err := mapping.AdaptFuncToDestroy(namespace, authorizeName)
	if err != nil {
		return nil, nil, err
	}

	destroyEndsession, err := mapping.AdaptFuncToDestroy(namespace, endsessionName)
	if err != nil {
		return nil, nil, err
	}

	destroyIssuer, err := mapping.AdaptFuncToDestroy(namespace, issuerName)
	if err != nil {
		return nil, nil, err
	}

	destroyers := []operator.DestroyFunc{
		operator.ResourceDestroyToZitadelDestroy(destroyAdminR),
		operator.ResourceDestroyToZitadelDestroy(destroyMgmtRest),
		operator.ResourceDestroyToZitadelDestroy(destroyOAuthv2),
		operator.ResourceDestroyToZitadelDestroy(destroyAuthR),
		operator.ResourceDestroyToZitadelDestroy(destroyAuthorize),
		operator.ResourceDestroyToZitadelDestroy(destroyEndsession),
		operator.ResourceDestroyToZitadelDestroy(destroyIssuer),
	}

	return func(k8sClient *kubernetes.Client, queried map[string]interface{}) (operator.EnsureFunc, error) {
			crd, err := k8sClient.CheckCRD("mappings.getambassador.io")
			if crd == nil || err != nil {
				return func(k8sClient *kubernetes.Client) error { return nil }, nil
			}

			accountsDomain := dns.Subdomains.Accounts + "." + dns.Domain
			apiDomain := dns.Subdomains.API + "." + dns.Domain
			issuerDomain := dns.Subdomains.Issuer + "." + dns.Domain

			cors := &mapping.CORS{
				Origins:        "*",
				Methods:        "POST, GET, OPTIONS, DELETE, PUT",
				Headers:        "*",
				Credentials:    true,
				ExposedHeaders: "*",
				MaxAge:         "86400",
			}

			queryAdminR, err := mapping.AdaptFuncToEnsure(
				namespace,
				adminRName,
				labels,
				false,
				apiDomain,
				"/admin/v1",
				"",
				httpUrl,
				"30000",
				"30000",
				cors,
			)
			if err != nil {
				return nil, err
			}

			queryMgmtRest, err := mapping.AdaptFuncToEnsure(
				namespace,
				mgmtName,
				labels,
				false,
				apiDomain,
				"/management/v1/",
				"",
				httpUrl,
				"30000",
				"30000",
				cors,
			)
			if err != nil {
				return nil, err
			}

			queryOAuthv2, err := mapping.AdaptFuncToEnsure(
				namespace,
				oauthName,
				labels,
				false,
				apiDomain,
				"/oauth/v2/",
				"",
				httpUrl,
				"30000",
				"30000",
				cors,
			)
			if err != nil {
				return nil, err
			}

			queryAuthR, err := mapping.AdaptFuncToEnsure(
				namespace,
				authRName,
				labels,
				false,
				apiDomain,
				"/auth/v1/",
				"",
				httpUrl,
				"30000",
				"30000",
				cors,
			)
			if err != nil {
				return nil, err
			}

			queryAuthorize, err := mapping.AdaptFuncToEnsure(
				namespace,
				authorizeName,
				labels,
				false,
				accountsDomain,
				"/oauth/v2/authorize",
				"",
				httpUrl,
				"30000",
				"30000",
				cors,
			)
			if err != nil {
				return nil, err
			}

			queryEndsession, err := mapping.AdaptFuncToEnsure(
				namespace,
				endsessionName,
				labels,
				false,
				accountsDomain,
				"/oauth/v2/endsession",
				"",
				httpUrl,
				"30000",
				"30000",
				cors,
			)
			if err != nil {
				return nil, err
			}

			queryIssuer, err := mapping.AdaptFuncToEnsure(
				namespace,
				issuerName,
				labels,
				false,
				issuerDomain,
				"/.well-known/openid-configuration",
				"/oauth/v2/.well-known/openid-configuration",
				httpUrl,
				"30000",
				"30000",
				cors,
			)
			if err != nil {
				return nil, err
			}

			queriers := []operator.QueryFunc{
				operator.ResourceQueryToZitadelQuery(queryAdminR),
				operator.ResourceQueryToZitadelQuery(queryMgmtRest),
				operator.ResourceQueryToZitadelQuery(queryOAuthv2),
				operator.ResourceQueryToZitadelQuery(queryAuthR),
				operator.ResourceQueryToZitadelQuery(queryAuthorize),
				operator.ResourceQueryToZitadelQuery(queryEndsession),
				operator.ResourceQueryToZitadelQuery(queryIssuer),
			}

			return operator.QueriersToEnsureFunc(internalMonitor, false, queriers, k8sClient, queried)
		},
		operator.DestroyersToDestroyFunc(internalMonitor, destroyers),
		nil
}
