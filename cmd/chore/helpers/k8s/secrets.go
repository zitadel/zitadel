package k8s_test

import (
	"gopkg.in/yaml.v3"
	"sort"
)

type secret struct {
	Metadata struct {
		Name string
	}
	Data map[string]string
}

func GetSecretKeysWithName(kubectl KubectlCmd, namespace, name string) func() []string {
	return func() []string {
		keys := make([]string, 0)
		secret := GetSecretWithName(kubectl, namespace, name)()
		for k := range secret.Data {
			keys = append(keys, k)
		}

		sort.Strings(keys)
		return keys
	}
}

func GetSecretWithName(kubectl KubectlCmd, namespace, name string) func() secret {
	return func() secret {
		secret := secret{}
		args := []string{
			"get", "secret", name,
			"--namespace", namespace,
			"--output", "yaml",
		}

		cmd := kubectl(args...)
		out, err := cmd.CombinedOutput()
		if err != nil {
			return secret
		}

		if err := yaml.Unmarshal(out, &secret); err != nil {
			return secret
		}
		return secret
	}
}
