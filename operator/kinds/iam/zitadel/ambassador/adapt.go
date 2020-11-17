package ambassador

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/zitadel/operator"
	"github.com/caos/zitadel/operator/kinds/iam/zitadel/ambassador/grpc"
	"github.com/caos/zitadel/operator/kinds/iam/zitadel/ambassador/hosts"
	"github.com/caos/zitadel/operator/kinds/iam/zitadel/ambassador/http"
	"github.com/caos/zitadel/operator/kinds/iam/zitadel/ambassador/ui"
	"github.com/caos/zitadel/operator/kinds/iam/zitadel/configuration"
)

func AdaptFunc(
	monitor mntr.Monitor,
	namespace string,
	labels map[string]string,
	grpcURL string,
	httpURL string,
	uiURL string,
	dns *configuration.DNS,
) (
	operator.QueryFunc,
	operator.DestroyFunc,
	error,
) {
	internalMonitor := monitor.WithField("component", "ambassador")

	internalLabels := make(map[string]string, 0)
	for k, v := range labels {
		internalLabels[k] = v
	}
	internalLabels["app.kubernetes.io/component"] = "ambassador"

	queryGRPC, destroyGRPC, err := grpc.AdaptFunc(internalMonitor, namespace, labels, grpcURL, dns)
	if err != nil {
		return nil, nil, err
	}

	queryUI, destroyHTTP, err := ui.AdaptFunc(internalMonitor, namespace, labels, uiURL, dns)
	if err != nil {
		return nil, nil, err
	}

	queryHTTP, destroyUI, err := http.AdaptFunc(internalMonitor, namespace, labels, httpURL, dns)
	if err != nil {
		return nil, nil, err
	}

	queryHosts, destroyHosts, err := hosts.AdaptFunc(internalMonitor, namespace, labels, dns)
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
