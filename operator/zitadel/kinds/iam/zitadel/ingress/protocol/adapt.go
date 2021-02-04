package protocol

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/labels"
	"github.com/caos/zitadel/operator"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/configuration"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/ingress/protocol/core"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/ingress/protocol/grpc"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/ingress/protocol/http"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/ingress/protocol/ui"
)

func AdaptFunc(
	monitor mntr.Monitor,
	componentLabels *labels.Component,
	namespaceStr string,
	ingressDefinitionSuffix string,
	grpcService string,
	grpcPort uint16,
	httpService string,
	httpPort uint16,
	uiService string,
	uiPort uint16,
	dns *configuration.DNS,
	controllerSpecifics map[string]interface{},
	queryIngress core.IngressDefinitionQueryFunc,
	destroyIngress core.IngressDefinitionDestroyFunc,

) (
	operator.QueryFunc,
	operator.DestroyFunc,
	error,
) {
	queryGRPC, destroyGRPC, err := grpc.AdaptFunc(
		monitor,
		componentLabels,
		namespaceStr,
		ingressDefinitionSuffix,
		grpcService,
		grpcPort,
		dns,
		controllerSpecifics,
		queryIngress,
		destroyIngress,
	)
	if err != nil {
		return nil, nil, err
	}

	queryUI, destroyHTTP, err := ui.AdaptFunc(
		monitor,
		componentLabels,
		namespaceStr,
		ingressDefinitionSuffix,
		uiService,
		uiPort,
		dns,
		controllerSpecifics,
		queryIngress,
		destroyIngress,
	)
	if err != nil {
		return nil, nil, err
	}

	queryHTTP, destroyUI, err := http.AdaptFunc(
		monitor,
		componentLabels,
		namespaceStr,
		ingressDefinitionSuffix,
		httpService,
		httpPort,
		dns,
		controllerSpecifics,
		queryIngress,
		destroyIngress,
	)
	if err != nil {
		return nil, nil, err
	}

	destroyers := []operator.DestroyFunc{
		destroyGRPC,
		destroyUI,
		destroyHTTP,
	}

	queriers := []operator.QueryFunc{
		queryHTTP,
		queryUI,
		queryGRPC,
	}

	return func(k8sClient kubernetes.ClientInt, queried map[string]interface{}) (operator.EnsureFunc, error) {
			return operator.QueriersToEnsureFunc(monitor, true, queriers, k8sClient, queried)
		},
		operator.DestroyersToDestroyFunc(monitor, destroyers),
		nil
}
