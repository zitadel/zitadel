package user

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/labels"
	"github.com/caos/orbos/pkg/secret"
	"github.com/caos/orbos/pkg/secret/read"
	"github.com/caos/zitadel/operator"
	"github.com/caos/zitadel/operator/database/kinds/databases/managed/user/dbuser"
	"github.com/caos/zitadel/operator/database/kinds/databases/managed/user/node"
	"github.com/caos/zitadel/pkg/databases/certs/client"
	"github.com/caos/zitadel/pkg/databases/db"
)

const (
	execDBPod       = "cockroachdb-0"
	execDBContainer = "cockroachdb"
	rootUserName    = "root"
)

func AdaptFunc(
	monitor mntr.Monitor,
	namespace string,
	componentLabels *labels.Component,
	clusterDns string,
	generateNodeIfNotExists bool,
	userName string,
	dbPasswd *secret.Secret,
	dbPasswdExisting *secret.Existing,
	pwSecretLabels *labels.Selectable,
	pwSecretKey string,
	rootCertsSecret string,
	containerCertsDir string,
	nodeSecret string,
	dbConn db.Connection,
) (
	operator.QueryFunc,
	operator.DestroyFunc,
	operator.QueryFunc,
	operator.DestroyFunc,
	error,
) {
	cMonitor := monitor.WithField("type", "users")

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

	queryDBUser, destroyDBUser, err := dbuser.AdaptFunc(
		monitor,
		namespace,
		execDBPod,
		execDBContainer,
		containerCertsDir,
		userName,
		pwSecretLabels,
		pwSecretKey,
		func(k8sClient kubernetes.ClientInt) (string, error) {
			pwValue, err := read.GetSecretValue(k8sClient, dbPasswd, dbPasswdExisting)
			if err != nil {
				return "", err
			}
			if pwValue == "" {
				pwValue = userName
			}
			return pwValue, nil
		},
	)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	queryCert, destroyCert, err := client.AdaptFunc(
		cMonitor,
		namespace,
		componentLabels,
		dbConn,
	)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	beforeCRqueriers := []operator.QueryFunc{
		queryNode,
		queryCert(rootUserName, rootCertsSecret, db.RootUserCert, db.RootUserKey),
		queryCert(userName, db.CertsSecret(userName), db.UserCert, db.UserKey),
	}

	beforeCRdestroyers := []operator.DestroyFunc{
		destroyCert(db.CertsSecret(userName)),
		destroyCert(rootCertsSecret),
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
		nil
}
