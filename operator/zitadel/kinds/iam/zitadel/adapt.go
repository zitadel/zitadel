package zitadel

import (
	"strconv"

	"github.com/caos/orbos/pkg/labels"
	"github.com/caos/orbos/pkg/secret"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/database"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/setup"

	core "k8s.io/api/core/v1"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/tree"
	"github.com/caos/zitadel/operator"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/ambassador"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/configuration"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/deployment"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/migration"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/services"
	"github.com/pkg/errors"
)

func AdaptFunc(
	apiLabels *labels.API,
	nodeselector map[string]string,
	tolerations []core.Toleration,
	dbClient database.Client,
	namespace string,
	action string,
	version *string,
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
		usersWithoutPWs := getUserListWithoutPasswords(desiredKind)

		zitadelComponent := labels.MustForComponent(apiLabels, "ZITADEL")
		zitadelDeploymentName := labels.MustForName(zitadelComponent, "zitadel")
		zitadelPodSelector := labels.DeriveNameSelector(zitadelDeploymentName, false)
		queryS, destroyS, err := services.AdaptFunc(
			internalMonitor,
			zitadelComponent,
			zitadelPodSelector,
			namespace,
			grpcServiceName,
			grpcPort,
			httpServiceName,
			httpPort,
			uiServiceName,
			uiPort)
		if err != nil {
			return nil, nil, allSecrets, err
		}

		getQueryC, destroyC, getConfigurationHashes, err := configuration.AdaptFunc(
			internalMonitor,
			zitadelComponent,
			namespace,
			desiredKind.Spec.Configuration,
			cmName,
			certPath,
			secretName,
			secretPath,
			consoleCMName,
			secretVarsName,
			secretPasswordName,
			dbClient,
			services.GetClientIDFunc(namespace, httpServiceName, httpPort),
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
			namespace,
			action,
			secretPasswordName,
			migrationUser,
			usersWithoutPWs,
			nodeselector,
			tolerations,
		)
		if err != nil {
			return nil, nil, allSecrets, err
		}

		getQuerySetup, destroySetup, err := setup.AdaptFunc(
			internalMonitor,
			zitadelComponent,
			namespace,
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
		)
		if err != nil {
			return nil, nil, allSecrets, err
		}

		queryD, destroyD, err := deployment.AdaptFunc(
			internalMonitor,
			zitadelDeploymentName,
			zitadelPodSelector,
			desiredKind.Spec.Force,
			version,
			namespace,
			desiredKind.Spec.ReplicaCount,
			desiredKind.Spec.Affinity,
			cmName,
			certPath,
			secretName,
			secretPath,
			consoleCMName,
			secretVarsName,
			secretPasswordName,
			desiredKind.Spec.NodeSelector,
			desiredKind.Spec.Tolerations,
			desiredKind.Spec.Resources,
			migration.GetDoneFunc(monitor, namespace, action),
			configuration.GetReadyFunc(monitor, namespace, secretName, secretVarsName, secretPasswordName, cmName, consoleCMName),
			setup.GetDoneFunc(monitor, namespace, action),
		)
		if err != nil {
			return nil, nil, allSecrets, err
		}

		queryAmbassador, destroyAmbassador, err := ambassador.AdaptFunc(
			internalMonitor,
			labels.MustForComponent(apiLabels, "apiGateway"),
			namespace,
			grpcServiceName+"."+namespace+":"+strconv.Itoa(grpcPort),
			"http://"+httpServiceName+"."+namespace+":"+strconv.Itoa(httpPort),
			"http://"+uiServiceName+"."+namespace,
			desiredKind.Spec.Configuration.DNS,
		)
		if err != nil {
			return nil, nil, allSecrets, err
		}

		destroyers := make([]operator.DestroyFunc, 0)
		for _, feature := range features {
			switch feature {
			case "migration":
				destroyers = append(destroyers,
					destroyM,
				)
			case "iam":
				destroyers = append(destroyers,
					destroyAmbassador,
					destroyS,
					destroyM,
					destroyD,
					destroySetup,
					destroyC,
				)
			}
		}

		return func(k8sClient kubernetes.ClientInt, queried map[string]interface{}) (operator.EnsureFunc, error) {
				users, err := getAllUsers(k8sClient, desiredKind)
				if err != nil {
					return nil, err
				}
				allZitadelUsers, err := getZitadelUserList(k8sClient, desiredKind)
				if err != nil {
					return nil, err
				}

				queryReadyM := operator.EnsureFuncToQueryFunc(migration.GetDoneFunc(monitor, namespace, action))
				queryC := getQueryC(users)
				queryReadyC := operator.EnsureFuncToQueryFunc(configuration.GetReadyFunc(monitor, namespace, secretName, secretVarsName, secretPasswordName, cmName, consoleCMName))
				querySetup := getQuerySetup(allZitadelUsers, getConfigurationHashes)
				queryReadySetup := operator.EnsureFuncToQueryFunc(setup.GetDoneFunc(monitor, namespace, action))
				queryD := queryD(allZitadelUsers, getConfigurationHashes)
				queryReadyD := operator.EnsureFuncToQueryFunc(deployment.GetReadyFunc(monitor, namespace, zitadelDeploymentName))

				queriers := make([]operator.QueryFunc, 0)
				for _, feature := range features {
					switch feature {
					case "migration":
						queriers = append(queriers,
							queryDB,
							//configuration
							queryC,
							queryReadyC,
							//migration
							queryM,
							queryReadyM,
						)
					case "iam":
						queriers = append(queriers,
							queryDB,
							//configuration
							queryC,
							queryReadyC,
							//migration
							queryM,
							queryReadyM,
							//services
							queryS,
							//setup
							querySetup,
							queryReadySetup,
							//deployment
							queryD,
							queryReadyD,
							//handle change if necessary for clientID
							queryC,
							queryReadyC,
							//again apply deployment if config changed
							queryD,
							queryReadyD,
							//apply ambassador crds after zitadel is ready
							queryAmbassador,
						)
					case "scaledown":
						queriers = append(queriers,
							operator.EnsureFuncToQueryFunc(deployment.GetScaleFunc(monitor, namespace, zitadelDeploymentName)(0)),
						)
					case "scaleup":
						queriers = append(queriers,
							operator.EnsureFuncToQueryFunc(deployment.GetScaleFunc(monitor, namespace, zitadelDeploymentName)(desiredKind.Spec.ReplicaCount)),
						)
					}
				}

				return operator.QueriersToEnsureFunc(internalMonitor, true, queriers, k8sClient, queried)
			},
			operator.DestroyersToDestroyFunc(monitor, destroyers),
			allSecrets,
			nil
	}
}
