package cockroachdb

import (
	"fmt"
	"strconv"

	"github.com/caos/orbos/pkg/secret/read"

	k8sSecret "github.com/caos/orbos/pkg/kubernetes/resources/secret"

	"github.com/caos/orbos/pkg/labels"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/secret"
	"github.com/caos/orbos/pkg/tree"

	"github.com/caos/zitadel/operator"
	"github.com/caos/zitadel/pkg/databases/db"
)

const (
	namespace = "caos-zitadel"
	component = "dbconnection"
)

func Adapter(apiLabels *labels.API) operator.AdaptFunc {
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

		internalMonitor := monitor.WithField("kind", "cockroachdb")

		desiredKind, err := parseDesiredV0(desired)
		if err != nil {
			return nil, nil, nil, nil, nil, false, fmt.Errorf("parsing desired state failed: %w", err)
		}
		desired.Parsed = desiredKind

		if desiredKind.Spec.Verbose {
			internalMonitor = internalMonitor.Verbose()
		}

		if err := desiredKind.validate(); err != nil {
			return nil, nil, nil, nil, nil, false, err
		}

		allSecrets, allExisting := getSecretsMap(desiredKind)

		currentDB := &Current{
			Common: tree.NewCommon("zitadel.caos.ch/CockroachDB", "v0", false),
		}
		current.Parsed = currentDB

		componentLabels := labels.MustForComponent(apiLabels, component)
		certLabels := labels.MustForName(componentLabels, db.CertsSecret)
		pwLabels := labels.AsSelectable(labels.MustForName(componentLabels, "db-connection-password"))

		return func(k8sClient kubernetes.ClientInt, queried map[string]interface{}) (operator.EnsureFunc, error) {

				currentDB.Current.Host = desiredKind.Spec.Host
				currentDB.Current.Cluster = desiredKind.Spec.Cluster
				currentDB.Current.Port = strconv.Itoa(int(desiredKind.Spec.Port))
				if currentDB.Current.Port == "" {
					currentDB.Current.Port = "26257"
				}

				currentDB.Current.User = desiredKind.Spec.User
				if currentDB.Current.User == "" {
					currentDB.Current.User = "root"
				}

				certificate, err := read.GetSecretValue(k8sClient, desiredKind.Spec.Certificate, desiredKind.Spec.ExistingCertificate)
				if err != nil {
					return nil, err
				}

				var queriers []operator.QueryFunc
				if certificate != "" {
					certQuerier, err := k8sSecret.AdaptFuncToEnsure(namespace, labels.AsSelectable(certLabels), map[string]string{
						db.CACert: certificate,
					})
					if err != nil {
						return nil, err
					}
					queriers = append(queriers, operator.ResourceQueryToZitadelQuery(certQuerier))
				}
				currentDB.Current.Secure = certificate != ""

				password, err := read.GetSecretValue(k8sClient, desiredKind.Spec.Password, desiredKind.Spec.ExistingPassword)
				if err != nil {
					return nil, err
				}

				if password != "" {
					currentDB.Current.PasswordSecret = pwLabels
					currentDB.Current.PasswordSecretKey = currentDB.Current.User
					pwQuerier, err := k8sSecret.AdaptFuncToEnsure(namespace, pwLabels, map[string]string{
						currentDB.Current.PasswordSecretKey: password,
					})
					if err != nil {
						return nil, err
					}
					queriers = append(queriers, operator.ResourceQueryToZitadelQuery(pwQuerier))
				}

				/*

					if password != "" {
						certQuerier, err := k8sSecret.AdaptFuncToEnsure(namespace, labels.AsSelectable(certLabels), map[string]string{
							: certificate,
						})
						if err != nil {
							return nil, err
						}
						queriers = append(queriers, operator.ResourceQueryToZitadelQuery(certQuerier))

					}

				*/

				db.SetQueriedForDatabase(queried, current)
				return operator.QueriersToEnsureFunc(internalMonitor, false, queriers, k8sClient, queried)
			}, func(k8sClient kubernetes.ClientInt) error { return nil },
			func(kubernetes.ClientInt, map[string]interface{}, bool) error { return nil },
			allSecrets,
			allExisting,
			false,
			nil
	}
}
