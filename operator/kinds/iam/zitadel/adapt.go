package zitadel

import (
	"strconv"

	"github.com/caos/orbos/pkg/labels"
	"github.com/caos/orbos/pkg/orb"
	"github.com/caos/orbos/pkg/secret"
	"github.com/caos/zitadel/operator/kinds/iam/zitadel/database"
	"github.com/caos/zitadel/operator/kinds/iam/zitadel/setup"

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
	apiLabels *labels.API,
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

		//		labels := getLabels()
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

		zitadelComponent := labels.MustForComponent(apiLabels, "ZITADEL")
		zitadelDeploymentName := labels.MustForName(zitadelComponent, "zitadel")
		zitadelPodSelector := labels.DeriveNameSelector(zitadelDeploymentName, false)
		queryS, destroyS, err := services.AdaptFunc(
			internalMonitor,
			zitadelComponent,
			zitadelPodSelector,
			namespaceStr,
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
			zitadelComponent,
			namespaceStr,
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

		queryM, destroyM, err := migration.AdaptFunc(
			internalMonitor,
			labels.MustForComponent(apiLabels, "database"),
			namespaceStr,
			action,
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

		querySetup, destroySetup, err := setup.AdaptFunc(
			internalMonitor,
			zitadelComponent,
			namespaceStr,
			action,
			desiredKind.Spec.NodeSelector,
			desiredKind.Spec.Tolerations,
			desiredKind.Spec.Resources,
			version,
			cmName,
			certPath,
			secretName,
			secretPath,
			consoleCMName,
			secretVarsName,
			secretPasswordName,
			allZitadelUsers,
			migration.GetDoneFunc(monitor, namespaceStr, action),
			configuration.GetReadyFunc(monitor, namespaceStr, secretName, secretVarsName, secretPasswordName, cmName, consoleCMName),
			getConfigurationHashes,
		)

		queryD, destroyD, err := deployment.AdaptFunc(
			internalMonitor,
			zitadelDeploymentName,
			zitadelPodSelector,
			desiredKind.Spec.Force,
			version,
			namespaceStr,
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
			migration.GetDoneFunc(monitor, namespaceStr, action),
			configuration.GetReadyFunc(monitor, namespaceStr, secretName, secretVarsName, secretPasswordName, cmName, consoleCMName),
			setup.GetDoneFunc(monitor, namespaceStr, action),
			getConfigurationHashes,
		)
		if err != nil {
			return nil, nil, allSecrets, err
		}

		queryAmbassador, destroyAmbassador, err := ambassador.AdaptFunc(
			internalMonitor,
			labels.MustForComponent(apiLabels, "apiGateway"),
			namespaceStr,
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
					querySetup,
					queryD,
					operator.EnsureFuncToQueryFunc(deployment.GetReadyFunc(monitor, namespaceStr)),
					queryAmbassador,
				)
				destroyers = append(destroyers,
					destroyAmbassador,
					destroyS,
					destroyM,
					destroyD,
					destroySetup,
					destroyC,
					operator.ResourceDestroyToZitadelDestroy(destroyNS),
				)
			case "scaledown":
				queriers = append(queriers,
					operator.EnsureFuncToQueryFunc(deployment.GetScaleFunc(monitor, namespaceStr)(0)),
				)
			case "scaleup":
				queriers = append(queriers,
					operator.EnsureFuncToQueryFunc(deployment.GetScaleFunc(monitor, namespaceStr)(desiredKind.Spec.ReplicaCount)),
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
