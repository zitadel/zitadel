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

	queryConsole, destroyConsole, err := consoleAdapter(
		monitor,
		namespace,
		labels.MustForName(componentLabels, ConsoleName),
		false,
		originCASecretName,
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
		return nil, nil, err
	}

	queryAcc, destroyAcc, err := accountsAdapter(
		monitor,
		namespace,
		labels.MustForName(componentLabels, AccountsName),
		false,
		originCASecretName,
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
