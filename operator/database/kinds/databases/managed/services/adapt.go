package services

import (
	"github.com/caos/zitadel/operator"
	"strconv"

	"github.com/caos/orbos/pkg/labels"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/kubernetes/resources/service"
)

func AdaptFunc(
	monitor mntr.Monitor,
	namespace string,
	publicServiceNameLabels *labels.Name,
	privateServiceNameLabels *labels.Name,
	cockroachSelector *labels.Selector,
	cockroachPort int32,
	cockroachHTTPPort int32,
) (
	operator.QueryFunc,
	operator.DestroyFunc,
	error,
) {
	internalMonitor := monitor.WithField("type", "services")

	publicServiceSelectable := labels.AsSelectable(publicServiceNameLabels)

	destroySPD, err := service.AdaptFuncToDestroy("default", publicServiceSelectable.Name())
	if err != nil {
		return nil, nil, err
	}
	destroySP, err := service.AdaptFuncToDestroy(namespace, publicServiceSelectable.Name())
	if err != nil {
		return nil, nil, err
	}
	destroyS, err := service.AdaptFuncToDestroy(namespace, privateServiceNameLabels.Name())
	if err != nil {
		return nil, nil, err
	}
	destroyers := []operator.DestroyFunc{
		operator.ResourceDestroyToZitadelDestroy(destroySPD),
		operator.ResourceDestroyToZitadelDestroy(destroySP),
		operator.ResourceDestroyToZitadelDestroy(destroyS),
	}

	ports := []service.Port{
		{Port: 26257, TargetPort: strconv.Itoa(int(cockroachPort)), Name: "grpc"},
		{Port: 8080, TargetPort: strconv.Itoa(int(cockroachHTTPPort)), Name: "http"},
	}
	querySPD, err := service.AdaptFuncToEnsure("default", publicServiceSelectable, ports, "", cockroachSelector, false, "", "")
	if err != nil {
		return nil, nil, err
	}
	querySP, err := service.AdaptFuncToEnsure(namespace, publicServiceSelectable, ports, "", cockroachSelector, false, "", "")
	if err != nil {
		return nil, nil, err
	}
	queryS, err := service.AdaptFuncToEnsure(namespace, privateServiceNameLabels, ports, "", cockroachSelector, true, "None", "")
	if err != nil {
		return nil, nil, err
	}

	queriers := []operator.QueryFunc{
		operator.ResourceQueryToZitadelQuery(querySPD),
		operator.ResourceQueryToZitadelQuery(querySP),
		operator.ResourceQueryToZitadelQuery(queryS),
	}

	return func(k8sClient kubernetes.ClientInt, queried map[string]interface{}) (operator.EnsureFunc, error) {

			return operator.QueriersToEnsureFunc(internalMonitor, false, queriers, k8sClient, queried)
		},
		operator.DestroyersToDestroyFunc(internalMonitor, destroyers),
		nil

}
