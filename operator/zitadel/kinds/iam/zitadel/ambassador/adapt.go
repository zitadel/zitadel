package ambassador

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/labels"
	"github.com/caos/zitadel/operator"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/ambassador/grpc"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/ambassador/hosts"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/ambassador/http"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/ambassador/ui"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/configuration"
)

func AdaptFunc(
	monitor mntr.Monitor,
	componentLabels *labels.Component,
	namespace string,
	grpcURL string,
	httpURL string,
	uiURL string,
	dns *configuration.DNS,
) (
	operator.QueryFunc,
	operator.DestroyFunc,
	error,
) {
	internalMonitor := monitor.WithField("type", "ambassador")

	queryGRPC, destroyGRPC, err := grpc.AdaptFunc(internalMonitor, componentLabels, namespace, grpcURL, dns)
	if err != nil {
		return nil, nil, err
	}

	queryUI, destroyHTTP, err := ui.AdaptFunc(internalMonitor, componentLabels, namespace, uiURL, dns)
	if err != nil {
		return nil, nil, err
	}

	queryHTTP, destroyUI, err := http.AdaptFunc(internalMonitor, componentLabels, namespace, httpURL, dns)
	if err != nil {
		return nil, nil, err
	}

	queryHosts, destroyHosts, err := hosts.AdaptFunc(internalMonitor, componentLabels, namespace, dns)
	if err != nil {
		return nil, nil, err
	}

	destroyers := []operator.DestroyFunc{
		destroyGRPC,
		destroyHTTP,
		destroyUI,
		destroyHosts,
	}

	queriers := []operator.QueryFunc{
		queryHosts,
		queryGRPC,
		queryUI,
		queryHTTP,
	}

	return func(k8sClient kubernetes.ClientInt, queried map[string]interface{}) (operator.EnsureFunc, error) {
			return operator.QueriersToEnsureFunc(internalMonitor, true, queriers, k8sClient, queried)
		},
		operator.DestroyersToDestroyFunc(internalMonitor, destroyers),
		nil
}
