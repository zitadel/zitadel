package ambassador

import (
	"fmt"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/kubernetes/resources"
	"github.com/caos/orbos/pkg/kubernetes/resources/ambassador/mapping"
	"github.com/caos/orbos/pkg/labels"
	"github.com/caos/zitadel/operator"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/configuration"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/ingress/controllers/ambassador/hosts"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/ingress/protocol/core"
)

var (
	DestroyMapping core.IngressDefinitionDestroyFunc = mapping.AdaptFuncToDestroy
	_              core.IngressDefinitionQueryFunc   = QueryMapping
)

func QueryMapping(
	namespace string,
	id labels.IDLabels,
	grpc bool,
	host,
	prefix,
	rewrite,
	service string,
	servicePort uint16,
	timeoutMS,
	connectTimeoutMS int,
	cors *core.CORS,
	_ map[string]interface{},
) (resources.QueryFunc, error) {

	return mapping.AdaptFuncToEnsure(
		namespace,
		id,
		grpc,
		host,
		prefix,
		rewrite,
		fmt.Sprintf("%s:%d", service, servicePort),
		timeoutMS,
		connectTimeoutMS,
		cors.ToAmassadorCORS(),
	)
}

func AdaptFunc(
	monitor mntr.Monitor,
	componentLabels *labels.Component,
	namespace string,
	dns *configuration.DNS,
) (
	operator.QueryFunc,
	operator.DestroyFunc,
	error,
) {
	internalMonitor := monitor.WithField("type", "ambassador")

	queryHosts, destroyHosts, err := hosts.AdaptFunc(internalMonitor, componentLabels, namespace, dns)
	if err != nil {
		return nil, nil, err
	}

	destroyers := []operator.DestroyFunc{
		destroyHosts,
	}

	queriers := []operator.QueryFunc{
		queryHosts,
	}

	return func(k8sClient kubernetes.ClientInt, queried map[string]interface{}) (operator.EnsureFunc, error) {
			return operator.QueriersToEnsureFunc(internalMonitor, true, queriers, k8sClient, queried)
		},
		operator.DestroyersToDestroyFunc(internalMonitor, destroyers),
		nil
}
