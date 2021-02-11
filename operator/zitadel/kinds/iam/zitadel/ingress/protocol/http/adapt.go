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
	controllerSpecifics map[string]interface{},
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

	queryAdminR, destroyAdminR, err := apiAdapter(
		monitor,
		namespace,
		labels.MustForName(componentLabels, AdminRName),
		false,
		originCASecretName,
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
		return nil, nil, err
	}

	queryMgmtRest, destroyMgmtRest, err := apiAdapter(
		monitor,
		namespace,
		labels.MustForName(componentLabels, MgmtName),
		false,
		originCASecretName,
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
		return nil, nil, err
	}

	queryOAuthv2, destroyOAuthv2, err := apiAdapter(
		monitor,
		namespace,
		labels.MustForName(componentLabels, OauthName),
		false,
		originCASecretName,
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
		return nil, nil, err
	}

	queryAuthR, destroyAuthR, err := apiAdapter(
		monitor,
		namespace,
		labels.MustForName(componentLabels, AuthRName),
		false,
		originCASecretName,
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
		return nil, nil, err
	}

	queryAuthorize, destroyAuthorize, err := accountsAdapter(
		monitor,
		namespace,
		labels.MustForName(componentLabels, AuthorizeName),
		false,
		originCASecretName,
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
		return nil, nil, err
	}

	queryEndsession, destroyEndsession, err := accountsAdapter(
		monitor,
		namespace,
		labels.MustForName(componentLabels, EndsessionName),
		false,
		originCASecretName,
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
		return nil, nil, err
	}

	queryIssuer, destroyIssuer, err := issuerAdapter(
		monitor,
		namespace,
		labels.MustForName(componentLabels, IssuerName),
		false,
		originCASecretName,
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
