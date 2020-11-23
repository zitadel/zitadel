package zitadel

import (
	"github.com/caos/orbos/pkg/orb"
	"github.com/caos/orbos/pkg/secret"
	"github.com/caos/zitadel/operator/kinds/iam/zitadel/database"
	"strconv"

	core "k8s.io/api/core/v1"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/kubernetes/resources/namespace"
	"github.com/caos/orbos/pkg/tree"
	"github.com/caos/zitadel/operator"
	"github.com/caos/zitadel/operator/kinds/iam/zitadel/ambassador"
	"github.com/caos/zitadel/operator/kinds/iam/zitadel/configuration"
	"github.com/caos/zitadel/operator/kinds/iam/zitadel/deployment"
	"github.com/caos/zitadel/operator/kinds/iam/zitadel/migration"
	"github.com/caos/zitadel/operator/kinds/iam/zitadel/services"
	"github.com/pkg/errors"
)

func AdaptFunc(
	nodeselector map[string]string,
	tolerations []core.Toleration,
	orbconfig *orb.Orb,
	action string,
	migrationsPath string,
	version string,
	features []string,
) operator.AdaptFunc {
	return func(
		monitor mntr.Monitor,
		desired *tree.Tree,
		current *tree.Tree,
	) (
		operator.QueryFunc,
		operator.DestroyFunc,
		map[string]*secret.Secret,
		error,
	) {

		allSecrets := make(map[string]*secret.Secret)
		internalMonitor := monitor.WithField("kind", "iam")

		desiredKind, err := parseDesiredV0(desired)
		if err != nil {
			return nil, nil, allSecrets, errors.Wrap(err, "parsing desired state failed")
		}
		desired.Parsed = desiredKind
		secret.AppendSecrets("", allSecrets, getSecretsMap(desiredKind))

		if !monitor.IsVerbose() && desiredKind.Spec.Verbose {
			internalMonitor.Verbose()
		}

		namespaceStr := "caos-zitadel"
		// shared elements
		cmName := "zitadel-vars"
		secretName := "zitadel-secret"
		consoleCMName := "console-config"
		secretVarsName := "zitadel-secrets-vars"
		secretPasswordName := "zitadel-passwords"
		//paths which are used in the configuration and also are used for mounting the used files
		certPath := "/home/zitadel/dbsecrets-zitadel"
		secretPath := "/secret"
		//services which are kubernetes resources and are used in the ambassador elements
		grpcServiceName := "grpc-v1"
		grpcPort := 80
		httpServiceName := "http-v1"
		httpPort := 80
		uiServiceName := "ui-v1"
		uiPort := 80

		labels := getLabels()
		users := getAllUsers(desiredKind)
		allZitadelUsers := getZitadelUserList()
		dbClient, err := database.NewClient(monitor, orbconfig.URL, orbconfig.Repokey)
		if err != nil {
			return nil, nil, allSecrets, err
		}

		queryNS, err := namespace.AdaptFuncToEnsure(namespaceStr)
		if err != nil {
			return nil, nil, allSecrets, err
		}
		destroyNS, err := namespace.AdaptFuncToDestroy(namespaceStr)
		if err != nil {
			return nil, nil, allSecrets, err
		}

		queryS, destroyS, err := services.AdaptFunc(
			internalMonitor,
			namespaceStr,
			labels,
			grpcServiceName,
			grpcPort,
			httpServiceName,
			httpPort,
			uiServiceName,
			uiPort)
		if err != nil {
			return nil, nil, allSecrets, err
		}

		queryC, destroyC, getConfigurationHashes, err := configuration.AdaptFunc(
			internalMonitor,
			namespaceStr,
			labels,
			desiredKind.Spec.Configuration,
			cmName,
			certPath,
			secretName,
			secretPath,
			consoleCMName,
			secretVarsName,
			secretPasswordName,
			users,
			services.GetClientIDFunc(namespaceStr, httpServiceName, httpPort),
			dbClient,
		)
		if err != nil {
			return nil, nil, allSecrets, err
		}

		queryDB, err := database.AdaptFunc(
			monitor,
			dbClient,
		)
		if err != nil {
			return nil, nil, allSecrets, err
		}

		queryM, destroyM, migrationDone, _, err := migration.AdaptFunc(
			internalMonitor,
			namespaceStr,
			action,
			labels,
			secretPasswordName,
			migrationUser,
			allZitadelUsers,
			nodeselector,
			tolerations,
			migrationsPath,
		)
		if err != nil {
			return nil, nil, allSecrets, err
		}

		queryD, destroyD, deploymentReady, scaleDeployment, ensureInit, err := deployment.AdaptFunc(
			internalMonitor,
			version,
			namespaceStr,
			labels,
			desiredKind.Spec.ReplicaCount,
			desiredKind.Spec.Affinity,
			cmName,
			certPath,
			secretName,
			secretPath,
			consoleCMName,
			secretVarsName,
			secretPasswordName,
			allZitadelUsers,
			desiredKind.Spec.NodeSelector,
			desiredKind.Spec.Tolerations,
			desiredKind.Spec.Resources,
			migrationDone,
			configuration.GetReadyFunc(monitor, namespaceStr, secretName, secretVarsName, secretPasswordName, cmName, consoleCMName),
			getConfigurationHashes,
		)
		if err != nil {
			return nil, nil, allSecrets, err
		}

		queryAmbassador, destroyAmbassador, err := ambassador.AdaptFunc(
			internalMonitor,
			namespaceStr,
			labels,
			grpcServiceName+"."+namespaceStr+":"+strconv.Itoa(grpcPort),
			"http://"+httpServiceName+"."+namespaceStr+":"+strconv.Itoa(httpPort),
			"http://"+uiServiceName+"."+namespaceStr,
			desiredKind.Spec.Configuration.DNS,
		)
		if err != nil {
			return nil, nil, allSecrets, err
		}

		destroyers := make([]operator.DestroyFunc, 0)
		queriers := make([]operator.QueryFunc, 0)
		for _, feature := range features {
			switch feature {
			case "migration":
				queriers = append(queriers,
					queryDB,
					//configuration
					queryC,
					//migration
					queryM,
				)
				destroyers = append(destroyers,
					destroyM,
				)
			case "iam":
				queriers = append(queriers,
					operator.ResourceQueryToZitadelQuery(queryNS),
					queryDB,
					//configuration
					queryC,
					//migration
					queryM,
					//services
					queryS,
					queryD,
					operator.EnsureFuncToQueryFunc(ensureInit),
					operator.EnsureFuncToQueryFunc(deploymentReady),
					queryAmbassador,
				)
				destroyers = append(destroyers,
					destroyAmbassador,
					destroyS,
					destroyM,
					destroyD,
					destroyC,
					operator.ResourceDestroyToZitadelDestroy(destroyNS),
				)
			case "scaledown":
				queriers = append(queriers,
					operator.EnsureFuncToQueryFunc(scaleDeployment(0)),
				)
			case "scaleup":
				queriers = append(queriers,
					operator.EnsureFuncToQueryFunc(scaleDeployment(desiredKind.Spec.ReplicaCount)),
				)
			}
		}

		return func(k8sClient kubernetes.ClientInt, queried map[string]interface{}) (operator.EnsureFunc, error) {
				return operator.QueriersToEnsureFunc(internalMonitor, true, queriers, k8sClient, queried)
			},
			operator.DestroyersToDestroyFunc(monitor, destroyers),
			allSecrets,
			nil
	}
}
