package zitadel

import (
	"strconv"

	"github.com/caos/orbos/pkg/helper"

	"gopkg.in/yaml.v3"

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
	customImageRegistry string,
) operator.AdaptFunc {
	return func(
		monitor mntr.Monitor,
		desired *tree.Tree,
		current *tree.Tree,
	) (
		operator.QueryFunc,
		operator.DestroyFunc,
		operator.ConfigureFunc,
		map[string]*secret.Secret,
		map[string]*secret.Existing,
		bool,
		error,
	) {

		internalMonitor := monitor.WithField("kind", "iam")

		desiredKind, err := parseDesiredV0(desired)
		if err != nil {
			return nil, nil, nil, nil, nil, false, errors.Wrap(err, "parsing desired state failed")
		}
		desired.Parsed = desiredKind

		if err := desiredKind.Spec.validate(); err != nil {
			return nil, nil, nil, nil, nil, false, err
		}

		allSecrets, allExisting := getSecretsMap(desiredKind)

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
			uint16(grpcPort),
			httpServiceName,
			uint16(httpPort),
			uiServiceName,
			uint16(uiPort))
		if err != nil {
			return nil, nil, nil, nil, nil, false, err
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
			return nil, nil, nil, nil, nil, false, err
		}

		queryDB, err := database.AdaptFunc(
			monitor,
			dbClient,
		)
		if err != nil {
			return nil, nil, nil, nil, nil, false, err
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
			customImageRegistry,
		)
		if err != nil {
			return nil, nil, nil, nil, nil, false, err
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
			customImageRegistry,
		)
		if err != nil {
			return nil, nil, nil, nil, nil, false, err
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
			customImageRegistry,
		)
		if err != nil {
			return nil, nil, nil, nil, nil, false, err
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
			return nil, nil, nil, nil, nil, false, err
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

		concatQueriers := func(queriers ...operator.QueryFunc) operator.QueryFunc {
			return func(k8sClient kubernetes.ClientInt, queried map[string]interface{}) (ensureFunc operator.EnsureFunc, err error) {
				return operator.QueriersToEnsureFunc(
					monitor,
					true,
					queriers,
					k8sClient,
					queried,
				)
			}
		}

		queryCfg := func(k8sClient kubernetes.ClientInt, queried map[string]interface{}) (ensureFunc operator.EnsureFunc, err error) {
			users, err := getAllUsers(k8sClient, desiredKind)
			if err != nil {
				return nil, err
			}
			return concatQueriers(
				queryDB,
				getQueryC(users),
				operator.EnsureFuncToQueryFunc(configuration.GetReadyFunc(
					monitor,
					namespace,
					secretName,
					secretVarsName,
					secretPasswordName,
					cmName,
					consoleCMName,
				)),
			)(k8sClient, queried)
		}

		queryReadyD := operator.EnsureFuncToQueryFunc(deployment.GetReadyFunc(monitor, namespace, zitadelDeploymentName))

		return func(k8sClient kubernetes.ClientInt, queried map[string]interface{}) (operator.EnsureFunc, error) {
				allZitadelUsers, err := getZitadelUserList(k8sClient, desiredKind)
				if err != nil {
					return nil, err
				}

				queryReadyM := operator.EnsureFuncToQueryFunc(migration.GetDoneFunc(monitor, namespace, action))
				querySetup := getQuerySetup(allZitadelUsers, getConfigurationHashes)
				queryReadySetup := operator.EnsureFuncToQueryFunc(setup.GetDoneFunc(monitor, namespace, action))
				queryD := queryD(allZitadelUsers, getConfigurationHashes)

				queriers := make([]operator.QueryFunc, 0)
				for _, feature := range features {
					switch feature {
					case "migration":
						queriers = append(queriers,
							//configuration
							queryCfg,
							//migration
							queryM,
							queryReadyM,
							operator.EnsureFuncToQueryFunc(migration.GetCleanupFunc(monitor, namespace, action)),
						)
					case "iam":
						queriers = append(queriers,
							//configuration
							queryCfg,
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
							queryCfg,
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
			func(k8sClient kubernetes.ClientInt, queried map[string]interface{}, gitops bool) error {

				if desiredKind.Spec == nil {
					desiredKind.Spec = &Spec{}
				}
				if desiredKind.Spec.Configuration == nil {
					desiredKind.Spec.Configuration = &configuration.Configuration{}
				}
				if desiredKind.Spec.Configuration.Secrets == nil {
					desiredKind.Spec.Configuration.Secrets = &configuration.Secrets{}
				}
				if desiredKind.Spec.Configuration.Secrets.CookieID == "" {
					desiredKind.Spec.Configuration.Secrets.CookieID = "cookiekey_1"
				}
				if desiredKind.Spec.Configuration.Secrets.OTPVerificationID == "" {
					desiredKind.Spec.Configuration.Secrets.OTPVerificationID = "otpverificationkey_1"
				}
				if desiredKind.Spec.Configuration.Secrets.DomainVerificationID == "" {
					desiredKind.Spec.Configuration.Secrets.DomainVerificationID = "domainverificationkey_1"
				}
				if desiredKind.Spec.Configuration.Secrets.IDPConfigVerificationID == "" {
					desiredKind.Spec.Configuration.Secrets.IDPConfigVerificationID = "idpconfigverificationkey_1"
				}
				if desiredKind.Spec.Configuration.Secrets.OIDCKeysID == "" {
					desiredKind.Spec.Configuration.Secrets.OIDCKeysID = "oidckey_1"
				}
				if desiredKind.Spec.Configuration.Secrets.UserVerificationID == "" {
					desiredKind.Spec.Configuration.Secrets.UserVerificationID = "userverificationkey_1"
				}
				if gitops && desiredKind.Spec.Configuration.Secrets.Keys == nil {
					desiredKind.Spec.Configuration.Secrets.Keys = &secret.Secret{}
				}
				if !gitops && desiredKind.Spec.Configuration.Secrets.ExistingKeys == nil {
					desiredKind.Spec.Configuration.Secrets.ExistingKeys = &secret.Existing{}
				}

				keys := make(map[string]string)
				if gitops {
					if err := yaml.Unmarshal([]byte(desiredKind.Spec.Configuration.Secrets.Keys.Value), keys); err != nil {
						return err
					}
				} else {
					return errors.New("configure is not yet implemented for CRD mode")
				}

				if _, ok := keys[desiredKind.Spec.Configuration.Secrets.CookieID]; !ok {
					keys[desiredKind.Spec.Configuration.Secrets.CookieID] = helper.RandStringBytes(32)
				}
				if _, ok := keys[desiredKind.Spec.Configuration.Secrets.OTPVerificationID]; !ok {
					keys[desiredKind.Spec.Configuration.Secrets.OTPVerificationID] = helper.RandStringBytes(32)
				}
				if _, ok := keys[desiredKind.Spec.Configuration.Secrets.DomainVerificationID]; !ok {
					keys[desiredKind.Spec.Configuration.Secrets.DomainVerificationID] = helper.RandStringBytes(32)
				}
				if _, ok := keys[desiredKind.Spec.Configuration.Secrets.IDPConfigVerificationID]; !ok {
					keys[desiredKind.Spec.Configuration.Secrets.IDPConfigVerificationID] = helper.RandStringBytes(32)
				}
				if _, ok := keys[desiredKind.Spec.Configuration.Secrets.OIDCKeysID]; !ok {
					keys[desiredKind.Spec.Configuration.Secrets.OIDCKeysID] = helper.RandStringBytes(32)
				}
				if _, ok := keys[desiredKind.Spec.Configuration.Secrets.UserVerificationID]; !ok {
					keys[desiredKind.Spec.Configuration.Secrets.UserVerificationID] = helper.RandStringBytes(32)
				}

				newKeys, err := yaml.Marshal(keys)
				if err != nil {
					return err
				}

				desiredKind.Spec.Configuration.Secrets.Keys.Value = string(newKeys)
				return nil
			},
			allSecrets,
			allExisting,
			false,
			nil
	}
}
