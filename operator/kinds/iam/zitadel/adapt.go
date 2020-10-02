package zitadel

import (
	"github.com/caos/orbos/pkg/orb"
	"sort"
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
	features []string,
) operator.AdaptFunc {
	return func(
		monitor mntr.Monitor,
		desired *tree.Tree,
		current *tree.Tree,
	) (
		operator.QueryFunc,
		operator.DestroyFunc,
		error,
	) {

		internalMonitor := monitor.WithField("kind", "iam")

		desiredKind, err := parseDesiredV0(desired)
		if err != nil {
			return nil, nil, errors.Wrap(err, "parsing desired state failed")
		}
		desired.Parsed = desiredKind

		if !monitor.IsVerbose() && desiredKind.Spec.Verbose {
			internalMonitor.Verbose()
		}

		namespaceStr := "caos-zitadel"
		labels := map[string]string{
			"app.kubernetes.io/managed-by": "zitadel.caos.ch",
			"app.kubernetes.io/part-of":    "zitadel",
		}
		internalLabels := map[string]string{}
		for k, v := range labels {
			internalLabels[k] = v
		}
		internalLabels["app.kubernetes.io/component"] = "iam"

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

		users, migrationUser := getUsers(desiredKind)

		allZitadelUsers := make([]string, 0)
		for k := range users {
			if k != migrationUser {
				allZitadelUsers = append(allZitadelUsers, k)
			}
		}
		sort.Slice(allZitadelUsers, func(i, j int) bool {
			return allZitadelUsers[i] < allZitadelUsers[j]
		})

		allUsers := make([]string, 0)
		for k := range users {
			allUsers = append(allUsers, k)
		}
		sort.Slice(allUsers, func(i, j int) bool {
			return allUsers[i] < allUsers[j]
		})

		queryNS, err := namespace.AdaptFuncToEnsure(namespaceStr)
		if err != nil {
			return nil, nil, err
		}
		destroyNS, err := namespace.AdaptFuncToDestroy(namespaceStr)
		if err != nil {
			return nil, nil, err
		}

		queryS, destroyS, getClientID, err := services.AdaptFunc(internalMonitor, namespaceStr, internalLabels, grpcServiceName, grpcPort, httpServiceName, httpPort, uiServiceName, uiPort)
		if err != nil {
			return nil, nil, err
		}

		queryC, destroyC, configurationDone, getConfigurationHashes, err := configuration.AdaptFunc(
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
			getClientID,
			orbconfig.URL,
			orbconfig.Repokey,
		)
		if err != nil {
			return nil, nil, err
		}

		queryM, destroyM, migrationDone, _, err := migration.AdaptFunc(
			internalMonitor,
			namespaceStr,
			action,
			internalLabels,
			secretPasswordName,
			migrationUser,
			allZitadelUsers,
			nodeselector,
			tolerations,
			orbconfig.URL,
			orbconfig.Repokey,
		)
		if err != nil {
			return nil, nil, err
		}

		queryD, destroyD, deploymentReady, scaleDeployment, ensureInit, err := deployment.AdaptFunc(
			internalMonitor,
			namespaceStr,
			internalLabels,
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
			configurationDone,
			getConfigurationHashes,
		)
		if err != nil {
			return nil, nil, err
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
			return nil, nil, err
		}

		destroyers := make([]operator.DestroyFunc, 0)
		queriers := make([]operator.QueryFunc, 0)
		for _, feature := range features {
			switch feature {
			case "migration":
				queriers = append(queriers,
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

		return func(k8sClient *kubernetes.Client, queried map[string]interface{}) (operator.EnsureFunc, error) {
				return operator.QueriersToEnsureFunc(internalMonitor, true, queriers, k8sClient, queried)
			},
			operator.DestroyersToDestroyFunc(monitor, destroyers),
			nil
	}
}

func getUsers(desired *DesiredV0) (map[string]string, string) {
	passwords := &configuration.Passwords{}
	if desired.Spec != nil && desired.Spec.Configuration != nil && desired.Spec.Configuration.Passwords != nil {
		passwords = desired.Spec.Configuration.Passwords
	}
	users := make(map[string]string, 0)

	migrationUser := "flyway"
	migrationPassword := migrationUser
	if passwords.Migration != nil {
		migrationPassword = passwords.Migration.Value
	}
	users[migrationUser] = migrationPassword

	mgmtUser := "management"
	mgmtPassword := mgmtUser
	if passwords != nil && passwords.Management != nil {
		mgmtPassword = passwords.Management.Value
	}
	users[mgmtUser] = mgmtPassword

	adminUser := "adminapi"
	adminPassword := adminUser
	if passwords != nil && passwords.Adminapi != nil {
		adminPassword = passwords.Adminapi.Value
	}
	users[adminUser] = adminPassword

	authUser := "auth"
	authPassword := authUser
	if passwords != nil && passwords.Auth != nil {
		authPassword = passwords.Auth.Value
	}
	users[authUser] = authPassword

	authzUser := "authz"
	authzPassword := authzUser
	if passwords != nil && passwords.Authz != nil {
		authzPassword = passwords.Authz.Value
	}
	users[authzUser] = authzPassword

	notUser := "notification"
	notPassword := notUser
	if passwords != nil && passwords.Notification != nil {
		notPassword = passwords.Notification.Value
	}
	users[notUser] = notPassword

	esUser := "eventstore"
	esPassword := esUser
	if passwords != nil && passwords.Eventstore != nil {
		esPassword = passwords.Notification.Value
	}
	users[esUser] = esPassword

	return users, migrationUser
}
