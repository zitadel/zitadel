package ambassador

import (
	"fmt"
	"strings"

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

	return func(args core.PathArguments) (operator.QueryFunc, operator.DestroyFunc, error) {

		destroyMapping, err := mapping.AdaptFuncToDestroy(args.Namespace, args.ID.Name())
		if err != nil {
			return nil, nil, err
		}

		queryMapping, err := mapping.AdaptFuncToEnsure(
			args.Monitor,
			args.Namespace,
			args.ID,
			args.GRPC,
			virtualHost,
			args.Prefix,
			args.Rewrite,
			fmt.Sprintf("%s:%d", args.Service, args.ServicePort),
			args.TimeoutMS,
			args.ConnectTimeoutMS,
			args.CORS.ToAmassadorCORS(),
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

			destroyHost, err := host.AdaptFuncToDestroy(args.Namespace, hostName)
			if err != nil {
				return nil, nil, err
			}

			queryHost, err := host.AdaptFuncToEnsure(
				args.Monitor,
				args.Namespace,
				hostName,
				labels.MustK8sMap(args.ID),
				virtualHost,
				"none",
				"",
				map[string]string{
					"hostname": virtualHost,
				},
				args.OriginCASecretName,
			)
			if err != nil {
				return nil, nil, err
			}
			queriers = append(queriers, operator.ResourceQueryToZitadelQuery(queryHost))
			destroyers = append(destroyers, operator.ResourceDestroyToZitadelDestroy(destroyHost))
		}

		return func(k8sClient kubernetes.ClientInt, queried map[string]interface{}) (operator.EnsureFunc, error) {
				return operator.QueriersToEnsureFunc(args.Monitor, false, queriers, k8sClient, queried)
			},
			operator.DestroyersToDestroyFunc(args.Monitor, destroyers),
			nil
	}
}
