package http

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/labels"
	"github.com/caos/zitadel/operator"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/configuration"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/ingress/protocol/core"
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
	ingressDefinitionSuffix string,
	httpService string,
	httpPort uint16,
	dns *configuration.DNS,
	controllerSpecifics map[string]interface{},
	queryIngress core.IngressDefinitionQueryFunc,
	destroyIngress core.IngressDefinitionDestroyFunc,
) (
	operator.QueryFunc,
	operator.DestroyFunc,
	error,
) {
	internalMonitor := monitor.WithField("part", "http")

	fulladminRName := AdminRName + ingressDefinitionSuffix
	fullmgmtName := MgmtName + ingressDefinitionSuffix
	fulloauthName := OauthName + ingressDefinitionSuffix
	fullauthRName := AuthRName + ingressDefinitionSuffix
	fullauthorizeName := AuthorizeName + ingressDefinitionSuffix
	fullendsessionName := EndsessionName + ingressDefinitionSuffix
	fullissuerName := IssuerName + ingressDefinitionSuffix

	destroyAdminR, err := destroyIngress(namespace, fulladminRName)
	if err != nil {
		return nil, nil, err
	}

	destroyMgmtRest, err := destroyIngress(namespace, fullmgmtName)
	if err != nil {
		return nil, nil, err
	}

	destroyOAuthv2, err := destroyIngress(namespace, fulloauthName)
	if err != nil {
		return nil, nil, err
	}

	destroyAuthR, err := destroyIngress(namespace, fullauthRName)
	if err != nil {
		return nil, nil, err
	}

	destroyAuthorize, err := destroyIngress(namespace, fullauthorizeName)
	if err != nil {
		return nil, nil, err
	}

	destroyEndsession, err := destroyIngress(namespace, fullendsessionName)
	if err != nil {
		return nil, nil, err
	}

	destroyIssuer, err := destroyIngress(namespace, fullissuerName)
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

			cors := &core.CORS{
				Origins:        "*",
				Methods:        "POST, GET, OPTIONS, DELETE, PUT",
				Headers:        "*",
				Credentials:    true,
				ExposedHeaders: "*",
				MaxAge:         "86400",
			}

			queryAdminR, err := queryIngress(
				namespace,
				labels.MustForName(componentLabels, fulladminRName),
				false,
				apiDomain,
				"/admin/v1",
				"/admin/v1",
				httpService,
				httpPort,
				30000,
				30000,
				cors,
				controllerSpecifics,
			)
			if err != nil {
				return nil, err
			}

			queryMgmtRest, err := queryIngress(
				namespace,
				labels.MustForName(componentLabels, fullmgmtName),
				false,
				apiDomain,
				"/management/v1/",
				"/management/v1/",
				httpService,
				httpPort,
				30000,
				30000,
				cors,
				controllerSpecifics,
			)
			if err != nil {
				return nil, err
			}

			queryOAuthv2, err := queryIngress(
				namespace,
				labels.MustForName(componentLabels, fulloauthName),
				false,
				apiDomain,
				"/oauth/v2/",
				"/oauth/v2/",
				httpService,
				httpPort,
				30000,
				30000,
				cors,
				controllerSpecifics,
			)
			if err != nil {
				return nil, err
			}

			queryAuthR, err := queryIngress(
				namespace,
				labels.MustForName(componentLabels, fullauthRName),
				false,
				apiDomain,
				"/auth/v1/",
				"/auth/v1/",
				httpService,
				httpPort,
				30000,
				30000,
				cors,
				controllerSpecifics,
			)
			if err != nil {
				return nil, err
			}

			queryAuthorize, err := queryIngress(
				namespace,
				labels.MustForName(componentLabels, fullauthorizeName),
				false,
				accountsDomain,
				"/oauth/v2/authorize",
				"/oauth/v2/authorize",
				httpService,
				httpPort,
				30000,
				30000,
				cors,
				controllerSpecifics,
			)
			if err != nil {
				return nil, err
			}

			queryEndsession, err := queryIngress(
				namespace,
				labels.MustForName(componentLabels, fullendsessionName),
				false,
				accountsDomain,
				"/oauth/v2/endsession",
				"/oauth/v2/endsession",
				httpService,
				httpPort,
				30000,
				30000,
				cors,
				controllerSpecifics,
			)
			if err != nil {
				return nil, err
			}

			queryIssuer, err := queryIngress(
				namespace,
				labels.MustForName(componentLabels, fullissuerName),
				false,
				issuerDomain,
				"/.well-known/openid-configuration",
				"/oauth/v2/.well-known/openid-configuration",
				httpService,
				httpPort,
				30000,
				30000,
				cors,
				controllerSpecifics,
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
