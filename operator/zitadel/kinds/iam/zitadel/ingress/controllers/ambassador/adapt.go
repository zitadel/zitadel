package ambassador

import (
	"fmt"
	"strings"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/kubernetes/resources/ambassador/host"
	"github.com/caos/orbos/pkg/kubernetes/resources/ambassador/mapping"
	"github.com/caos/orbos/pkg/labels"
	"github.com/caos/zitadel/operator"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/ingress/protocol/core"
)

var _ core.HostAdapter = Adapt

func Adapt(virtualHost string) core.PathAdapter {

	seenHosts := make(map[string]struct{})

	return func(
		monitor mntr.Monitor,
		namespace string,
		id labels.IDLabels,
		grpc bool,
		originCASecretName,
		prefix,
		rewrite,
		service string,
		servicePort uint16,
		timeoutMS,
		connectTimeoutMS int,
		cors *core.CORS,
		controllerSpecifics map[string]interface{},
	) (operator.QueryFunc, operator.DestroyFunc, error) {

		destroyMapping, err := mapping.AdaptFuncToDestroy(namespace, id.Name())
		if err != nil {
			return nil, nil, err
		}

		queryMapping, err := mapping.AdaptFuncToEnsure(
			monitor,
			namespace,
			id,
			grpc,
			virtualHost,
			prefix,
			rewrite,
			fmt.Sprintf("%s:%d", service, servicePort),
			timeoutMS,
			connectTimeoutMS,
			cors.ToAmassadorCORS(),
		)
		if err != nil {
			return nil, nil, err
		}

		queriers := []operator.QueryFunc{
			operator.ResourceQueryToZitadelQuery(queryMapping),
		}

		destroyers := []operator.DestroyFunc{
			operator.ResourceDestroyToZitadelDestroy(destroyMapping),
		}

		if _, ok := seenHosts[virtualHost]; !ok {
			seenHosts[virtualHost] = struct{}{}

			hostName := strings.ReplaceAll(virtualHost, ".", "-")

			destroyHost, err := host.AdaptFuncToDestroy(namespace, hostName)
			if err != nil {
				return nil, nil, err
			}

			queryHost, err := host.AdaptFuncToEnsure(
				monitor,
				namespace,
				hostName,
				labels.MustK8sMap(id),
				virtualHost,
				"none",
				"",
				map[string]string{
					"hostname": virtualHost,
				},
				originCASecretName,
			)
			if err != nil {
				return nil, nil, err
			}
			queriers = append(queriers, operator.ResourceQueryToZitadelQuery(queryHost))
			destroyers = append(destroyers, operator.ResourceDestroyToZitadelDestroy(destroyHost))
		}

		return func(k8sClient kubernetes.ClientInt, queried map[string]interface{}) (operator.EnsureFunc, error) {
				return operator.QueriersToEnsureFunc(monitor, false, queriers, k8sClient, queried)
			},
			operator.DestroyersToDestroyFunc(monitor, destroyers),
			nil
	}
}
