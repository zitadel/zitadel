package services

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/kubernetes/resources/service"
	"github.com/caos/orbos/pkg/labels"
	"github.com/caos/zitadel/operator"
)

func AdaptFunc(
	monitor mntr.Monitor,
	componentLabels *labels.Component,
	zitadelPodSelector *labels.Selector,
	namespace string,
	grpcServiceName string,
	grpcPort uint16,
	httpServiceName string,
	httpPort uint16,
	uiServiceName string,
	uiPort uint16,
) (
	operator.QueryFunc,
	operator.DestroyFunc,
	error,
) {
	internalMonitor := monitor.WithField("component", "services")

	destroyGRPC, err := service.AdaptFuncToDestroy(namespace, grpcServiceName)
	if err != nil {
		return nil, nil, err
	}

	destroyHTTP, err := service.AdaptFuncToDestroy(namespace, httpServiceName)
	if err != nil {
		return nil, nil, err
	}

	destroyUI, err := service.AdaptFuncToDestroy(namespace, uiServiceName)
	if err != nil {
		return nil, nil, err
	}

	destroyers := []operator.DestroyFunc{
		operator.ResourceDestroyToZitadelDestroy(destroyGRPC),
		operator.ResourceDestroyToZitadelDestroy(destroyHTTP),
		operator.ResourceDestroyToZitadelDestroy(destroyUI),
	}

	grpcPorts := []service.Port{
		{Name: "grpc", Port: grpcPort, TargetPort: "grpc"},
	}
	queryGRPC, err := service.AdaptFuncToEnsure(namespace, labels.MustForName(componentLabels, grpcServiceName), grpcPorts, "", zitadelPodSelector, false, "", "")
	if err != nil {
		return nil, nil, err
	}

	httpPorts := []service.Port{
		{Name: "http", Port: httpPort, TargetPort: "http"},
	}
	queryHTTP, err := service.AdaptFuncToEnsure(namespace, labels.AsSelectable(labels.MustForName(componentLabels, httpServiceName)), httpPorts, "", zitadelPodSelector, false, "", "")
	if err != nil {
		return nil, nil, err
	}

	uiPorts := []service.Port{
		{Name: "ui", Port: uiPort, TargetPort: "ui"},
	}
	queryUI, err := service.AdaptFuncToEnsure(namespace, labels.MustForName(componentLabels, uiServiceName), uiPorts, "", zitadelPodSelector, false, "", "")
	if err != nil {
		return nil, nil, err
	}

	queriers := []operator.QueryFunc{
		operator.ResourceQueryToZitadelQuery(queryGRPC),
		operator.ResourceQueryToZitadelQuery(queryHTTP),
		operator.ResourceQueryToZitadelQuery(queryUI),
	}

	return func(k8sClient kubernetes.ClientInt, queried map[string]interface{}) (operator.EnsureFunc, error) {
			return operator.QueriersToEnsureFunc(internalMonitor, false, queriers, k8sClient, queried)
		},
		operator.DestroyersToDestroyFunc(internalMonitor, destroyers),

		nil
}
