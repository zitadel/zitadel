package operator

import (
	"fmt"

	"gopkg.in/yaml.v3"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/git"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/kubernetes/resources"
	"github.com/caos/orbos/pkg/secret"
	"github.com/caos/orbos/pkg/tree"
)

type AdaptFunc func(
	monitor mntr.Monitor,
	desired *tree.Tree,
	current *tree.Tree,
) (
	QueryFunc,
	DestroyFunc,
	ConfigureFunc,
	map[string]*secret.Secret,
	map[string]*secret.Existing,
	bool,
	error,
)

type EnsureFunc func(k8sClient kubernetes.ClientInt) error

type DestroyFunc func(k8sClient kubernetes.ClientInt) error

type ConfigureFunc func(k8sClient kubernetes.ClientInt, queried map[string]interface{}, gitops bool) error

type QueryFunc func(k8sClient kubernetes.ClientInt, queried map[string]interface{}) (EnsureFunc, error)

func Parse(gitClient *git.Client, file string) (*tree.Tree, error) {
	if err := gitClient.Clone(); err != nil {
		return nil, err
	}

	tree := &tree.Tree{}
	if err := yaml.Unmarshal(gitClient.Read(file), tree); err != nil {
		return nil, err
	}

	return tree, nil
}

func ResourceDestroyToZitadelDestroy(destroyFunc resources.DestroyFunc) DestroyFunc {
	return func(k8sClient kubernetes.ClientInt) error {
		return destroyFunc(k8sClient)
	}
}

func ResourceQueryToZitadelQuery(queryFunc resources.QueryFunc) QueryFunc {
	return func(k8sClient kubernetes.ClientInt, _ map[string]interface{}) (EnsureFunc, error) {
		ensure, err := queryFunc(k8sClient)
		ensureInternal := ResourceEnsureToZitadelEnsure(ensure)

		return func(k8sClient kubernetes.ClientInt) error {
			return ensureInternal(k8sClient)
		}, err
	}
}

func ResourceEnsureToZitadelEnsure(ensureFunc resources.EnsureFunc) EnsureFunc {
	return func(k8sClient kubernetes.ClientInt) error {
		return ensureFunc(k8sClient)
	}
}
func EnsureFuncToQueryFunc(ensure EnsureFunc) QueryFunc {
	return func(k8sClient kubernetes.ClientInt, queried map[string]interface{}) (ensureFunc EnsureFunc, err error) {
		return ensure, err
	}
}

func QueriersToEnsureFunc(monitor mntr.Monitor, infoLogs bool, queriers []QueryFunc, k8sClient kubernetes.ClientInt, queried map[string]interface{}) (EnsureFunc, error) {
	if infoLogs {
		monitor.Info("querying...")
	} else {
		monitor.Debug("querying...")
	}
	ensurers := make([]EnsureFunc, 0)
	for _, querier := range queriers {
		ensurer, err := querier(k8sClient, queried)
		if err != nil {
			return nil, fmt.Errorf("error while querying: %w", err)
		}
		ensurers = append(ensurers, ensurer)
	}
	if infoLogs {
		monitor.Info("queried")
	} else {
		monitor.Debug("queried")
	}
	return func(k8sClient kubernetes.ClientInt) error {
		if infoLogs {
			monitor.Info("ensuring...")
		} else {
			monitor.Debug("ensuring...")
		}
		for _, ensurer := range ensurers {
			if err := ensurer(k8sClient); err != nil {
				return fmt.Errorf("error while ensuring: %w", err)
			}
		}
		if infoLogs {
			monitor.Info("ensured")
		} else {
			monitor.Debug("ensured")
		}
		return nil
	}, nil
}

func DestroyersToDestroyFunc(monitor mntr.Monitor, destroyers []DestroyFunc) DestroyFunc {
	return func(k8sClient kubernetes.ClientInt) error {
		monitor.Info("destroying...")
		for _, destroyer := range destroyers {
			if err := destroyer(k8sClient); err != nil {
				return fmt.Errorf("error while destroying: %w", err)
			}
		}
		monitor.Info("destroyed")
		return nil
	}
}
