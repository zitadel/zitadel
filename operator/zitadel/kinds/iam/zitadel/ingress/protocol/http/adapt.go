package http

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/labels"
	"github.com/caos/zitadel/operator"
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
	httpService string,
	httpPort uint16,
	controllerSpecifics map[string]string,
	originCASecretName string,
	apiAdapter core.PathAdapter,
	accountsAdapter core.PathAdapter,
	issuerAdapter core.PathAdapter,
) (
	operator.QueryFunc,
	operator.DestroyFunc,
	error,
) {
	internalMonitor := monitor.WithField("part", "http")

	cors := &core.CORS{
		Origins:        "*",
		Methods:        "POST, GET, OPTIONS, DELETE, PUT",
		Headers:        "*",
		Credentials:    true,
		ExposedHeaders: "*",
		MaxAge:         "86400",
	}

	queryAdminR, destroyAdminR, err := apiAdapter(core.PathArguments{
		Monitor:             monitor,
		Namespace:           namespace,
		ID:                  labels.MustForName(componentLabels, AdminRName),
		GRPC:                false,
		OriginCASecretName:  originCASecretName,
		Prefix:              "/admin/v1",
		Rewrite:             "/admin/v1",
		Service:             httpService,
		ServicePort:         httpPort,
		TimeoutMS:           30000,
		ConnectTimeoutMS:    30000,
		CORS:                cors,
		ControllerSpecifics: controllerSpecifics,
	})
	if err != nil {
		return nil, nil, err
	}

	queryMgmtRest, destroyMgmtRest, err := apiAdapter(core.PathArguments{
		Monitor:             monitor,
		Namespace:           namespace,
		ID:                  labels.MustForName(componentLabels, MgmtName),
		GRPC:                false,
		OriginCASecretName:  originCASecretName,
		Prefix:              "/management/v1/",
		Rewrite:             "/management/v1/",
		Service:             httpService,
		ServicePort:         httpPort,
		TimeoutMS:           30000,
		ConnectTimeoutMS:    30000,
		CORS:                cors,
		ControllerSpecifics: controllerSpecifics,
	})
	if err != nil {
		return nil, nil, err
	}

	queryOAuthv2, destroyOAuthv2, err := apiAdapter(core.PathArguments{
		Monitor:             monitor,
		Namespace:           namespace,
		ID:                  labels.MustForName(componentLabels, OauthName),
		GRPC:                false,
		OriginCASecretName:  originCASecretName,
		Prefix:              "/oauth/v2/",
		Rewrite:             "/oauth/v2/",
		Service:             httpService,
		ServicePort:         httpPort,
		TimeoutMS:           30000,
		ConnectTimeoutMS:    30000,
		CORS:                cors,
		ControllerSpecifics: controllerSpecifics,
	})
	if err != nil {
		return nil, nil, err
	}

	queryAuthR, destroyAuthR, err := apiAdapter(core.PathArguments{
		Monitor:             monitor,
		Namespace:           namespace,
		ID:                  labels.MustForName(componentLabels, AuthRName),
		GRPC:                false,
		OriginCASecretName:  originCASecretName,
		Prefix:              "/auth/v1/",
		Rewrite:             "/auth/v1/",
		Service:             httpService,
		ServicePort:         httpPort,
		TimeoutMS:           30000,
		ConnectTimeoutMS:    30000,
		CORS:                cors,
		ControllerSpecifics: controllerSpecifics,
	})
	if err != nil {
		return nil, nil, err
	}

	queryAuthorize, destroyAuthorize, err := accountsAdapter(core.PathArguments{
		Monitor:             monitor,
		Namespace:           namespace,
		ID:                  labels.MustForName(componentLabels, AuthorizeName),
		GRPC:                false,
		OriginCASecretName:  originCASecretName,
		Prefix:              "/oauth/v2/authorize",
		Rewrite:             "/oauth/v2/authorize",
		Service:             httpService,
		ServicePort:         httpPort,
		TimeoutMS:           30000,
		ConnectTimeoutMS:    30000,
		CORS:                cors,
		ControllerSpecifics: controllerSpecifics,
	})
	if err != nil {
		return nil, nil, err
	}

	queryEndsession, destroyEndsession, err := accountsAdapter(core.PathArguments{
		Monitor:             monitor,
		Namespace:           namespace,
		ID:                  labels.MustForName(componentLabels, EndsessionName),
		GRPC:                false,
		OriginCASecretName:  originCASecretName,
		Prefix:              "/oauth/v2/endsession",
		Rewrite:             "/oauth/v2/endsession",
		Service:             httpService,
		ServicePort:         httpPort,
		TimeoutMS:           30000,
		ConnectTimeoutMS:    30000,
		CORS:                cors,
		ControllerSpecifics: controllerSpecifics,
	})
	if err != nil {
		return nil, nil, err
	}

	queryIssuer, destroyIssuer, err := issuerAdapter(core.PathArguments{
		Monitor:             monitor,
		Namespace:           namespace,
		ID:                  labels.MustForName(componentLabels, IssuerName),
		GRPC:                false,
		OriginCASecretName:  originCASecretName,
		Prefix:              "/.well-known/openid-configuration",
		Rewrite:             "/oauth/v2/.well-known/openid-configuration",
		Service:             httpService,
		ServicePort:         httpPort,
		TimeoutMS:           30000,
		ConnectTimeoutMS:    30000,
		CORS:                cors,
		ControllerSpecifics: controllerSpecifics,
	})
	if err != nil {
		return nil, nil, err
	}

	return func(k8sClient kubernetes.ClientInt, queried map[string]interface{}) (operator.EnsureFunc, error) {
			return operator.QueriersToEnsureFunc(internalMonitor, false, []operator.QueryFunc{
				queryMgmtRest,
				queryOAuthv2,
				queryAuthR,
				queryAdminR,
				queryAuthorize,
				queryEndsession,
				queryIssuer,
			}, k8sClient, queried)
		},
		operator.DestroyersToDestroyFunc(internalMonitor, []operator.DestroyFunc{
			destroyAdminR,
			destroyMgmtRest,
			destroyOAuthv2,
			destroyAuthR,
			destroyAuthorize,
			destroyEndsession,
			destroyIssuer,
		}),
		nil
}
