package provided

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/secret"
	"github.com/caos/orbos/pkg/tree"
	"github.com/caos/zitadel/operator"
	"github.com/pkg/errors"
)

func AdaptFunc() func(
	monitor mntr.Monitor,
	desired *tree.Tree,
	current *tree.Tree,
) (
	operator.QueryFunc,
	operator.DestroyFunc,
	map[string]*secret.Secret,
	map[string]*secret.Existing,
	error,
) {
	return func(
		monitor mntr.Monitor,
		desired *tree.Tree,
		current *tree.Tree,
	) (
		operator.QueryFunc,
		operator.DestroyFunc,
		map[string]*secret.Secret,
		map[string]*secret.Existing,
		error,
	) {
		desiredKind, err := parseDesiredV0(desired)
		if err != nil {
			return nil, nil, nil, nil, errors.Wrap(err, "parsing desired state failed")
		}
		desired.Parsed = desiredKind

		currentDB := &Current{
			Common: &tree.Common{
				Kind:    "databases.caos.ch/ProvidedDatabase",
				Version: "v0",
			},
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
			make(map[string]*secret.Secret),
			make(map[string]*secret.Existing),
			nil
	}
}
