package cockroachdb

import (
	"fmt"

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
	certKey   = "client.root.crt"
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
		desiredKind, err := parseDesiredV0(desired)
		if err != nil {
			return nil, nil, nil, nil, nil, false, fmt.Errorf("parsing desired state failed: %w", err)
		}
		desired.Parsed = desiredKind

		if err := desiredKind.validate(); err != nil {
			return nil, nil, nil, nil, nil, false, err
		}

		allSecrets, allExisting := getSecretsMap(desiredKind)

		currentDB := &Current{
			Common: tree.NewCommon("zitadel.caos.ch/CockroachDB", "v0", false),
		}
		current.Parsed = currentDB

		certLabels := labels.MustForName(labels.MustForComponent(apiLabels, component), "cockroachdb.client.root")

		return func(k8sClient kubernetes.ClientInt, queried map[string]interface{}) (operator.EnsureFunc, error) {

				if err := desiredKind.validateSecrets(); err != nil {
					return nil, err
				}

				certificate, err := read.GetSecretValue(k8sClient, desiredKind.Spec.Certificate, desiredKind.Spec.ExistingCertificate)
				if err != nil {
					return nil, err
				}

				certQuerier, err := k8sSecret.AdaptFuncToEnsure(namespace, labels.AsSelectable(certLabels), map[string]string{
					certKey: certificate,
				})
				if err != nil {
					return nil, err
				}
				queriers := []operator.QueryFunc{operator.ResourceQueryToZitadelQuery(certQuerier)}

				currentDB.Current.URL = desiredKind.Spec.URL
				currentDB.Current.Port = desiredKind.Spec.Port
				db.SetQueriedForDatabase(queried, current)
				return operator.QueriersToEnsureFunc(monitor, true, queriers, k8sClient, queried)
			}, func(k8sClient kubernetes.ClientInt) error { return nil },
			func(kubernetes.ClientInt, map[string]interface{}, bool) error { return nil },
			allSecrets,
			allExisting,
			false,
			nil
	}
}
