package core

import (
	"errors"

	"github.com/caos/orbos/pkg/tree"
	"gopkg.in/yaml.v3"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func UnmarshalUnstructuredSpec(unstruct *unstructured.Unstructured) (*tree.Tree, error) {
	spec, found := unstruct.Object["spec"]
	if !found {
		return nil, errors.New("no spec in crd")
	}
	specMap, ok := spec.(map[string]interface{})
	if !ok {
		return nil, errors.New("no spec in crd")
	}

	data, err := yaml.Marshal(specMap)
	if err != nil {
		return nil, err
	}

	desired := &tree.Tree{}
	if err := yaml.Unmarshal(data, &desired); err != nil {
		return nil, err
	}

	return desired, nil
}

func MarshalToUnstructuredSpec(t *tree.Tree) (*unstructured.Unstructured, error) {
	data, err := yaml.Marshal(t)
	if err != nil {
		return nil, err
	}

	unstruct := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"spec": make(map[string]interface{}),
		},
	}
	return unstruct, yaml.Unmarshal(data, unstruct.Object["spec"])
}
