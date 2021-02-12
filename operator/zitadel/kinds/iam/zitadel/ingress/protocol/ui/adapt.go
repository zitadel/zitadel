package ui

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/labels"
	"github.com/caos/zitadel/operator"
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
	uiService string,
	uiPort uint16,
	controllerSpecifics map[string]interface{},
	originCASecretName string,
	consoleAdapter core.PathAdapter,
	accountsAdapter core.PathAdapter,
) (
	operator.QueryFunc,
	operator.DestroyFunc,
	error,
) {
	internalMonitor := monitor.WithField("part", "ui")

	queryConsole, destroyConsole, err := consoleAdapter(core.PathArguments{
		Monitor:             monitor,
		Namespace:           namespace,
		ID:                  labels.MustForName(componentLabels, ConsoleName),
		GRPC:                false,
		OriginCASecretName:  originCASecretName,
		Prefix:              "/",
		Rewrite:             "/console/",
		Service:             uiService,
		ServicePort:         uiPort,
		TimeoutMS:           0,
		ConnectTimeoutMS:    0,
		CORS:                nil,
		ControllerSpecifics: controllerSpecifics,
	})
	if err != nil {
		return nil, nil, err
	}

	queryAcc, destroyAcc, err := accountsAdapter(core.PathArguments{
		Monitor:             monitor,
		Namespace:           namespace,
		ID:                  labels.MustForName(componentLabels, AccountsName),
		GRPC:                false,
		OriginCASecretName:  originCASecretName,
		Prefix:              "/",
		Rewrite:             "/login/",
		Service:             uiService,
		ServicePort:         uiPort,
		TimeoutMS:           30000,
		ConnectTimeoutMS:    30000,
		CORS:                nil,
		ControllerSpecifics: controllerSpecifics,
	})
	if err != nil {
		return nil, nil, err
	}

	return func(k8sClient kubernetes.ClientInt, queried map[string]interface{}) (operator.EnsureFunc, error) {
			return operator.QueriersToEnsureFunc(internalMonitor, false, []operator.QueryFunc{
				queryConsole,
				queryAcc,
			}, k8sClient, queried)
		},
		operator.DestroyersToDestroyFunc(internalMonitor, []operator.DestroyFunc{
			destroyAcc,
			destroyConsole,
		}),
		nil
}
