package ingress

import (
	"fmt"

	anyingress "github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/ingress/controllers/ingress"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/ingress/controllers/nginx"

	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/ingress/protocol"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/ingress/protocol/core"

	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/configuration"

	"github.com/caos/orbos/pkg/kubernetes"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/labels"
	"github.com/caos/zitadel/operator"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/ingress/controllers/ambassador"
)

func AdaptFunc(
	monitor mntr.Monitor,
	apiLabels *labels.API,
	namespaceStr string,
	grpcServiceName string,
	grpcPort uint16,
	httpServiceName string,
	httpPort uint16,
	uiServiceName string,
	uiPort uint16,
	ingress *configuration.Ingress,
) (
	operator.QueryFunc,
	operator.DestroyFunc,
	error,
) {

	const component = "apiGateway"
	componentLabels := labels.MustForComponent(apiLabels, component)
	monitor = monitor.WithField("component", component)

	controllerIngressFunc := func(
		adapter core.HostAdapter) (
		operator.QueryFunc,
		operator.DestroyFunc,
		error,
	) {
		return protocol.AdaptFunc(
			monitor,
			componentLabels,
			namespaceStr,
			grpcServiceName,
			grpcPort,
			httpServiceName,
			httpPort,
			uiServiceName,
			uiPort,
			ingress,
			ingress.ControllerSpecifics,
			adapter,
		)
	}

	ambassadorIngQ, ambassadorIngD, err := controllerIngressFunc(ambassador.Adapt)
	if err != nil {
		return nil, nil, err
	}

	nginxIngQ, nginxIngD, err := controllerIngressFunc(nginx.Adapt)
	if err != nil {
		return nil, nil, err
	}

	anyIngQ, anyIngD, err := controllerIngressFunc(anyingress.Adapt)
	if err != nil {
		return nil, nil, err
	}

	destroyers := []operator.DestroyFunc{
		ambassadorIngD,
		nginxIngD,
		anyIngD,
	}

	var queriers []operator.QueryFunc
	switch ingress.Controller {
	case "Ambassador":
		queriers = append(queriers, ambassadorIngQ, operator.DestroyerToQueryFunc(nginxIngD), operator.DestroyerToQueryFunc(anyIngD))
	case "NGINX":
		queriers = append(queriers, nginxIngQ, operator.DestroyerToQueryFunc(ambassadorIngD))
	case "Any":
		queriers = append(queriers, anyIngQ, operator.DestroyerToQueryFunc(ambassadorIngD))
	case "None":
		queriers = operator.DestroyersToQueryFuncs(destroyers)
	default:
		return nil, nil, fmt.Errorf("unknown contoller type %s. possible values: Ambassador, NGINX, None", ingress.Controller)
	}

	return func(k8sClient kubernetes.ClientInt, queried map[string]interface{}) (operator.EnsureFunc, error) {
			return operator.QueriersToEnsureFunc(monitor, true, queriers, k8sClient, queried)
		},
		operator.DestroyersToDestroyFunc(monitor, destroyers),
		nil
}
