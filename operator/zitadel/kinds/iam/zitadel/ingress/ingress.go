package ingress

import (
	"fmt"

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
	dns *configuration.DNS,
	spec *Spec,
) (
	operator.QueryFunc,
	operator.DestroyFunc,
	error,
) {

	const component = "apiGateway"
	componentLabels := labels.MustForComponent(apiLabels, component)
	monitor = monitor.WithField("component", component)

	var (
		queriers            []operator.QueryFunc
		controller          = "None"
		controllerSpecifics map[string]interface{}
	)

	if spec != nil {
		controller = spec.Controller
		controllerSpecifics = spec.ControllerSpecifics
	}

	controllerIngressFunc := func(
		suffix string,
		queryIngressFunc core.IngressDefinitionQueryFunc,
		destroyDestroyFunc core.IngressDefinitionDestroyFunc) (
		operator.QueryFunc,
		operator.DestroyFunc,
		error,
	) {
		return protocol.AdaptFunc(
			monitor,
			componentLabels,
			namespaceStr,
			suffix,
			grpcServiceName,
			grpcPort,
			httpServiceName,
			httpPort,
			uiServiceName,
			uiPort,
			dns,
			controllerSpecifics,
			queryIngressFunc,
			destroyDestroyFunc,
		)
	}

	ambassadorIngQ, ambassadorIngD, err := controllerIngressFunc("", ambassador.QueryMapping, ambassador.DestroyMapping)
	if err != nil {
		return nil, nil, err
	}

	nginxIngQ, nginxIngD, err := controllerIngressFunc("-nginx", nginx.QueryIngress, nginx.DestroyIngress)
	if err != nil {
		return nil, nil, err
	}

	ambassadorQ, ambassadorD, err := ambassador.AdaptFunc(
		monitor,
		componentLabels,
		namespaceStr,
		dns,
	)
	if err != nil {
		return nil, nil, err
	}

	destroyers := []operator.DestroyFunc{
		ambassadorIngD,
		ambassadorD,
		nginxIngD,
	}

	switch controller {
	case "Ambassador":
		queriers = append(queriers, ambassadorQ, ambassadorIngQ, operator.DestroyerToQueryFunc(nginxIngD))
	case "NGINX":
		queriers = append(queriers, nginxIngQ, operator.DestroyerToQueryFunc(ambassadorD), operator.DestroyerToQueryFunc(ambassadorIngD))
	case "None":
		queriers = operator.DestroyersToQueryFuncs(destroyers)
	default:
		return nil, nil, fmt.Errorf("unknown contoller type %s. possible values: Ambassador, NGINX, None", controller)
	}

	return func(k8sClient kubernetes.ClientInt, queried map[string]interface{}) (operator.EnsureFunc, error) {
			return operator.QueriersToEnsureFunc(monitor, true, queriers, k8sClient, queried)
		},
		operator.DestroyersToDestroyFunc(monitor, destroyers),
		nil
}
