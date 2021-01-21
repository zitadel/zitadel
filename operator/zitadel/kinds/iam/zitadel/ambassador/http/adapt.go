package http

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/kubernetes/resources/ambassador/mapping"
	"github.com/caos/orbos/pkg/labels"
	"github.com/caos/zitadel/operator"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/configuration"
)

const (
	AdminRName     = "admin-rest-v1"
	MgmtName       = "mgmt-v1"
	OauthName      = "oauth-v1"
	AuthRName      = "auth-rest-v1"
	AuthorizeName  = "authorize-v1"
	EndsessionName = "endsession-v1"
	IssuerName     = "issuer-v1"
)

func AdaptFunc(
	monitor mntr.Monitor,
	componentLabels *labels.Component,
	namespace string,
	httpUrl string,
	dns *configuration.DNS,
) (
	operator.QueryFunc,
	operator.DestroyFunc,
	error,
) {
	internalMonitor := monitor.WithField("part", "http")

	destroyAdminR, err := mapping.AdaptFuncToDestroy(namespace, AdminRName)
	if err != nil {
		return nil, nil, err
	}

	destroyMgmtRest, err := mapping.AdaptFuncToDestroy(namespace, MgmtName)
	if err != nil {
		return nil, nil, err
	}

	destroyOAuthv2, err := mapping.AdaptFuncToDestroy(namespace, OauthName)
	if err != nil {
		return nil, nil, err
	}

	destroyAuthR, err := mapping.AdaptFuncToDestroy(namespace, AuthRName)
	if err != nil {
		return nil, nil, err
	}

	destroyAuthorize, err := mapping.AdaptFuncToDestroy(namespace, AuthorizeName)
	if err != nil {
		return nil, nil, err
	}

	destroyEndsession, err := mapping.AdaptFuncToDestroy(namespace, EndsessionName)
	if err != nil {
		return nil, nil, err
	}

	destroyIssuer, err := mapping.AdaptFuncToDestroy(namespace, IssuerName)
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

	return func(k8sClient kubernetes.ClientInt, queried map[string]interface{}) (operator.EnsureFunc, error) {
			crd, err := k8sClient.CheckCRD("mappings.getambassador.io")
			if crd == nil || err != nil {
				return func(k8sClient kubernetes.ClientInt) error { return nil }, nil
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
				AdminRName,
				labels.MustForNameK8SMap(componentLabels, AdminRName),
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
				MgmtName,
				labels.MustForNameK8SMap(componentLabels, MgmtName),
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
				OauthName,
				labels.MustForNameK8SMap(componentLabels, OauthName),
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
				AuthRName,
				labels.MustForNameK8SMap(componentLabels, AuthRName),
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
				AuthorizeName,
				labels.MustForNameK8SMap(componentLabels, AuthorizeName),
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
				EndsessionName,
				labels.MustForNameK8SMap(componentLabels, EndsessionName),
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
				IssuerName,
				labels.MustForNameK8SMap(componentLabels, IssuerName),
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
