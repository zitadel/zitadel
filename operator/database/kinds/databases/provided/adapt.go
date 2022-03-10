package provided

import (
	"fmt"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/secret"
	"github.com/caos/orbos/pkg/tree"

	"github.com/caos/zitadel/operator"
)

func Adapter() operator.AdaptFunc {
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

		currentDB := &Current{
			Common: tree.NewCommon("databases.caos.ch/ProvidedDatabase", "v0", false),
		}
		current.Parsed = currentDB

		return func(k8sClient kubernetes.ClientInt, _ map[string]interface{}) (operator.EnsureFunc, error) {
				currentDB.Current.URL = desiredKind.Spec.URL
				currentDB.Current.Port = desiredKind.Spec.Port

				return func(k8sClient kubernetes.ClientInt) error {
					return nil
				}, nil
			}, func(k8sClient kubernetes.ClientInt) error {
				return nil
			},
			func(kubernetes.ClientInt, map[string]interface{}, bool) error { return nil },
			make(map[string]*secret.Secret),
			make(map[string]*secret.Existing),
			false,
			nil
	}
}
