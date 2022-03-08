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
	Upload         = "assets-v1"
	MgmtName       = "mgmt-v1"
	OauthName      = "oauth-v1"
	AuthRName      = "auth-rest-v1"
	AuthorizeName  = "authorize-v1"
	EndsessionName = "endsession-v1"
	IssuerName     = "issuer-v1"
	OpenAPIName    = "openapi"
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

	destroyUpload, err := mapping.AdaptFuncToDestroy(namespace, Upload)
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

	destroySwagger, err := mapping.AdaptFuncToDestroy(namespace, OpenAPIName)
	if err != nil {
		return nil, nil, err
	}

	destroyers := []operator.DestroyFunc{
		operator.ResourceDestroyToZitadelDestroy(destroyAdminR),
		operator.ResourceDestroyToZitadelDestroy(destroyUpload),
		operator.ResourceDestroyToZitadelDestroy(destroyMgmtRest),
		operator.ResourceDestroyToZitadelDestroy(destroyOAuthv2),
		operator.ResourceDestroyToZitadelDestroy(destroyAuthR),
		operator.ResourceDestroyToZitadelDestroy(destroyAuthorize),
		operator.ResourceDestroyToZitadelDestroy(destroyEndsession),
		operator.ResourceDestroyToZitadelDestroy(destroyIssuer),
		operator.ResourceDestroyToZitadelDestroy(destroySwagger),
	}

	return func(k8sClient kubernetes.ClientInt, queried map[string]interface{}) (operator.EnsureFunc, error) {

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

			queryAdminR, err := mapping.AdaptFuncToEnsure(&mapping.Arguments{
				Monitor:          internalMonitor,
				Namespace:        namespace,
				ID:               labels.MustForName(componentLabels, AdminRName),
				GRPC:             false,
				Host:             apiDomain,
				Prefix:           "/admin/v1",
				Rewrite:          "",
				Service:          httpUrl,
				TimeoutMS:        30000,
				ConnectTimeoutMS: 30000,
				CORS:             cors,
			})
			if err != nil {
				return nil, err
			}

			queryUpload, err := mapping.AdaptFuncToEnsure(&mapping.Arguments{
				Monitor:          internalMonitor,
				Namespace:        namespace,
				ID:               labels.MustForName(componentLabels, Upload),
				GRPC:             false,
				Host:             apiDomain,
				Prefix:           "/assets/v1",
				Rewrite:          "",
				Service:          httpUrl,
				TimeoutMS:        30000,
				ConnectTimeoutMS: 30000,
				CORS:             cors,
			})
			if err != nil {
				return nil, err
			}

			queryMgmtRest, err := mapping.AdaptFuncToEnsure(&mapping.Arguments{
				Monitor:          internalMonitor,
				Namespace:        namespace,
				ID:               labels.MustForName(componentLabels, MgmtName),
				GRPC:             false,
				Host:             apiDomain,
				Prefix:           "/management/v1/",
				Rewrite:          "",
				Service:          httpUrl,
				TimeoutMS:        30000,
				ConnectTimeoutMS: 30000,
				CORS:             cors,
			})

			if err != nil {
				return nil, err
			}

			queryOAuthv2, err := mapping.AdaptFuncToEnsure(&mapping.Arguments{
				Monitor:          internalMonitor,
				Namespace:        namespace,
				ID:               labels.MustForName(componentLabels, OauthName),
				GRPC:             false,
				Host:             apiDomain,
				Prefix:           "/oauth/v2/",
				Rewrite:          "",
				Service:          httpUrl,
				TimeoutMS:        30000,
				ConnectTimeoutMS: 30000,
				CORS:             cors,
			})
			if err != nil {
				return nil, err
			}

			queryAuthR, err := mapping.AdaptFuncToEnsure(&mapping.Arguments{
				Monitor:          internalMonitor,
				Namespace:        namespace,
				ID:               labels.MustForName(componentLabels, AuthRName),
				GRPC:             false,
				Host:             apiDomain,
				Prefix:           "/auth/v1/",
				Rewrite:          "",
				Service:          httpUrl,
				TimeoutMS:        30000,
				ConnectTimeoutMS: 30000,
				CORS:             cors,
			})
			if err != nil {
				return nil, err
			}

			queryAuthorize, err := mapping.AdaptFuncToEnsure(&mapping.Arguments{
				Monitor:          internalMonitor,
				Namespace:        namespace,
				ID:               labels.MustForName(componentLabels, AuthorizeName),
				GRPC:             false,
				Host:             accountsDomain,
				Prefix:           "/oauth/v2/authorize",
				Rewrite:          "",
				Service:          httpUrl,
				TimeoutMS:        30000,
				ConnectTimeoutMS: 30000,
				CORS:             cors,
			})
			if err != nil {
				return nil, err
			}

			queryEndsession, err := mapping.AdaptFuncToEnsure(&mapping.Arguments{
				Monitor:          internalMonitor,
				Namespace:        namespace,
				ID:               labels.MustForName(componentLabels, EndsessionName),
				GRPC:             false,
				Host:             accountsDomain,
				Prefix:           "/oauth/v2/endsession",
				Rewrite:          "",
				Service:          httpUrl,
				TimeoutMS:        30000,
				ConnectTimeoutMS: 30000,
				CORS:             cors,
			})
			if err != nil {
				return nil, err
			}

			queryIssuer, err := mapping.AdaptFuncToEnsure(&mapping.Arguments{
				Monitor:          internalMonitor,
				Namespace:        namespace,
				ID:               labels.MustForName(componentLabels, IssuerName),
				GRPC:             false,
				Host:             issuerDomain,
				Prefix:           "/.well-known/openid-configuration",
				Rewrite:          "/oauth/v2/.well-known/openid-configuration",
				Service:          httpUrl,
				TimeoutMS:        30000,
				ConnectTimeoutMS: 30000,
				CORS:             cors,
			})
			if err != nil {
				return nil, err
			}

			queryOpenAPI, err := mapping.AdaptFuncToEnsure(&mapping.Arguments{
				Monitor:          internalMonitor,
				Namespace:        namespace,
				ID:               labels.MustForName(componentLabels, OpenAPIName),
				GRPC:             false,
				Host:             apiDomain,
				Prefix:           "/openapi/v2/swagger",
				Rewrite:          "",
				Service:          httpUrl,
				TimeoutMS:        30000,
				ConnectTimeoutMS: 30000,
				CORS:             nil,
			})
			if err != nil {
				return nil, err
			}

			queriers := []operator.QueryFunc{
				operator.ResourceQueryToZitadelQuery(queryAdminR),
				operator.ResourceQueryToZitadelQuery(queryUpload),
				operator.ResourceQueryToZitadelQuery(queryMgmtRest),
				operator.ResourceQueryToZitadelQuery(queryOAuthv2),
				operator.ResourceQueryToZitadelQuery(queryAuthR),
				operator.ResourceQueryToZitadelQuery(queryAuthorize),
				operator.ResourceQueryToZitadelQuery(queryEndsession),
				operator.ResourceQueryToZitadelQuery(queryIssuer),
				operator.ResourceQueryToZitadelQuery(queryOpenAPI),
			}

			return operator.QueriersToEnsureFunc(internalMonitor, false, queriers, k8sClient, queried)
		},
		operator.DestroyersToDestroyFunc(internalMonitor, destroyers),
		nil
}
