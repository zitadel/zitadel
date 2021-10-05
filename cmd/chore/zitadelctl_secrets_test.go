package chore_test

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

func getSecretKeysWithName(kubectl kubectlCmd, namespace, name string) []string {
	keys := make([]string, 0)
	secret := getSecretWithName(kubectl, namespace, name)
	for k := range secret.Data {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	return keys
}

func getSecretWithName(kubectl kubectlCmd, namespace, name string) secret {
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
