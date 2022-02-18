package certificate

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/labels"
	"github.com/caos/zitadel/operator"
	"github.com/caos/zitadel/operator/database/kinds/databases/managed/certificate/client"
	"github.com/caos/zitadel/operator/database/kinds/databases/managed/certificate/node"
	"github.com/caos/zitadel/operator/database/kinds/databases/managed/user"
	"github.com/caos/zitadel/pkg/databases/db"
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
	userName string,
	pwSecretLabels *labels.Selectable,
	pwSecretKey string,
) (
	operator.QueryFunc,
	operator.DestroyFunc,
	operator.QueryFunc,
	operator.DestroyFunc,
	// func(user, secretName, userCrtFilename, userKeyFilename string) (operator.QueryFunc, error),
	// func(secretName string) (operator.DestroyFunc, error),
	//	func(k8sClient kubernetes.ClientInt) ([]string, error),
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
		return nil, nil, nil, nil, err
	}

	/*TODO: dynamic variables*/
	queryDBUser, destroyDBUser, err := user.AdaptFunc(
		monitor,
		namespace,
		"cockroachdb-0",
		"cockroachdb",
		"/cockroach/cockroach-client-certs/",
		userName,
		"verysecret",
		db.CertsSecret,
		db.UserCert,
		db.UserKey,
		pwSecretLabels,
		pwSecretKey,
		componentLabels,
	)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	queryCert, destroyCert, err := client.AdaptFunc(
		cMonitor,
		namespace,
		componentLabels,
	)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	/*TODO: dynamic variables*/
	beforeCRqueriers := []operator.QueryFunc{
		queryNode,
		queryCert("root", "rootcerts", db.RootUserCert, db.RootUserKey),
		queryCert(userName, db.CertsSecret, db.UserCert, db.UserKey),
	}

	beforeCRdestroyers := []operator.DestroyFunc{
		destroyCert(db.CertsSecret),
		destroyCert("rootcerts"),
		destroyNode,
	}

	afterCRqueriers := []operator.QueryFunc{
		queryDBUser,
	}

	afterCRdestroyers := []operator.DestroyFunc{
		destroyDBUser,
	}

	return func(k8sClient kubernetes.ClientInt, queried map[string]interface{}) (operator.EnsureFunc, error) {
			return operator.QueriersToEnsureFunc(cMonitor, false, beforeCRqueriers, k8sClient, queried)
		},
		operator.DestroyersToDestroyFunc(cMonitor, beforeCRdestroyers),
		func(k8sClient kubernetes.ClientInt, queried map[string]interface{}) (operator.EnsureFunc, error) {
			return operator.QueriersToEnsureFunc(cMonitor, false, afterCRqueriers, k8sClient, queried)
		},
		operator.DestroyersToDestroyFunc(cMonitor, afterCRdestroyers),
		/*func(user, secretName, userCrtFilename, userKeyFilename string) (operator.QueryFunc, error) {
			query, _, err := client.AdaptFunc(
				cMonitor,
				namespace,
				componentLabels,
			)
			if err != nil {
				return nil, err
			}
			queryClient := query(user, secretName, userCrtFilename, userKeyFilename)

			return func(k8sClient kubernetes.ClientInt, queried map[string]interface{}) (operator.EnsureFunc, error) {
				_, err := queryNode(k8sClient, queried)
				if err != nil {
					return nil, err
				}

				return queryClient(k8sClient, queried)
			}, nil
		},
		func(secretName string) (operator.DestroyFunc, error) {
			_, destroy, err := client.AdaptFunc(
				cMonitor,
				namespace,
				componentLabels,
			)
			if err != nil {
				return nil, err
			}

			return destroy(secretName), nil
		},
				func(k8sClient kubernetes.ClientInt) ([]string, error) {
				return client.QueryCertificates(namespace, labels.DeriveComponentSelector(componentLabels, false), k8sClient)
			},*/
		nil
}
