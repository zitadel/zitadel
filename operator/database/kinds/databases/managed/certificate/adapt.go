package certificate

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/labels"
	"github.com/caos/zitadel/operator"
	"github.com/caos/zitadel/operator/database/kinds/databases/managed/certificate/client"
	"github.com/caos/zitadel/operator/database/kinds/databases/managed/certificate/node"
)

var (
	nodeSecret = "cockroachdb.node"
)

func AdaptFunc(
	monitor mntr.Monitor,
	namespace string,
	componentLabels *labels.Component,
	clusterDns string,
	generateNodeIfNotExists bool,
) (
	operator.QueryFunc,
	operator.DestroyFunc,
	func(user string) (operator.QueryFunc, error),
	func(user string) (operator.DestroyFunc, error),
	func(k8sClient kubernetes.ClientInt) ([]string, error),
	error,
) {
	cMonitor := monitor.WithField("type", "certificates")

	queryNode, destroyNode, err := node.AdaptFunc(
		cMonitor,
		namespace,
		labels.MustForName(componentLabels, nodeSecret),
		clusterDns,
		generateNodeIfNotExists,
	)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	queriers := []operator.QueryFunc{
		queryNode,
	}

	destroyers := []operator.DestroyFunc{
		destroyNode,
	}

	return func(k8sClient kubernetes.ClientInt, queried map[string]interface{}) (operator.EnsureFunc, error) {
			return operator.QueriersToEnsureFunc(cMonitor, false, queriers, k8sClient, queried)
		},
		operator.DestroyersToDestroyFunc(cMonitor, destroyers),
		func(user string) (operator.QueryFunc, error) {
			query, _, err := client.AdaptFunc(
				cMonitor,
				namespace,
				componentLabels,
			)
			if err != nil {
				return nil, err
			}
			queryClient := query(user)

			return func(k8sClient kubernetes.ClientInt, queried map[string]interface{}) (operator.EnsureFunc, error) {
				_, err := queryNode(k8sClient, queried)
				if err != nil {
					return nil, err
				}

				return queryClient(k8sClient, queried)
			}, nil
		},
		func(user string) (operator.DestroyFunc, error) {
			_, destroy, err := client.AdaptFunc(
				cMonitor,
				namespace,
				componentLabels,
			)
			if err != nil {
				return nil, err
			}

			return destroy(user), nil
		},
		func(k8sClient kubernetes.ClientInt) ([]string, error) {
			return client.QueryCertificates(namespace, labels.DeriveComponentSelector(componentLabels, false), k8sClient)
		},
		nil
}
