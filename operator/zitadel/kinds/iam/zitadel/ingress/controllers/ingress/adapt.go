package ingress

import (
	"github.com/caos/orbos/pkg/kubernetes"

	"github.com/caos/zitadel/operator"

	"github.com/caos/orbos/pkg/kubernetes/resources/ingress"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/ingress/protocol/core"
)

var _ core.HostAdapter = Adapt

func Adapt(virtualHost string) core.PathAdapter {
	return func(args core.PathArguments) (operator.QueryFunc, operator.DestroyFunc, error) {

		query, err := ingress.AdaptFuncToEnsure(
			args.Namespace,
			args.ID,
			virtualHost,
			args.Prefix,
			args.Service,
			args.ServicePort,
			args.ControllerSpecifics,
		)
		if err != nil {
			return nil, nil, err
		}

		destroy, err := ingress.AdaptFuncToDestroy(args.Namespace, args.ID.Name())
		if err != nil {
			return nil, nil, err
		}

		return func(k8sClient kubernetes.ClientInt, queried map[string]interface{}) (operator.EnsureFunc, error) {
				return operator.QueriersToEnsureFunc(args.Monitor, false, []operator.QueryFunc{
					operator.ResourceQueryToZitadelQuery(query),
				}, k8sClient, queried)
			},
			operator.DestroyersToDestroyFunc(args.Monitor, []operator.DestroyFunc{
				operator.ResourceDestroyToZitadelDestroy(destroy)}),
			nil
	}
}
