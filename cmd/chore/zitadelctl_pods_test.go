package chore_test

import (
	"gopkg.in/yaml.v3"
)

type pods struct {
	Items []struct {
		Metadata struct {
			Name string
		}
		Status struct {
			Conditions []struct {
				Type   string
				Status string
				Reason string
			}
		}
	}
}

func countCompletedPods(kubectl kubectlCmd, namespace, selector string) func() (readyPodsCount int8) {
	return func() int8 {
		pods, err := getPodsWithSelector(kubectl, namespace, selector)
		if err != nil {
			return -1
		}
		return countCompleted(pods)
	}
}

func countReadyPods(kubectl kubectlCmd, namespace, selector string) func() int8 {
	return func() int8 {
		pods, err := getPodsWithSelector(kubectl, namespace, selector)
		if err != nil {
			return -1
		}
		return countReady(pods)
	}
}

func getPodsWithSelector(kubectl kubectlCmd, namespace, selector string) (pods, error) {
	pods := pods{}
	args := []string{
		"get", "pods",
		"--namespace", namespace,
		"--output", "yaml",
	}

	if selector != "" {
		args = append(args, "--selector", selector)
	}

	cmd := kubectl(args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return pods, err
	}

	if err := yaml.Unmarshal(out, &pods); err != nil {
		return pods, err
	}
	return pods, nil
}

func countReady(pods pods) int8 {
	readyPodsCount := int8(0)
	for i := range pods.Items {
		pod := pods.Items[i]
		for j := range pod.Status.Conditions {
			condition := pod.Status.Conditions[j]
			if condition.Type != "Ready" {
				continue
			}
			if condition.Status == "True" {
				readyPodsCount++
				break
			}
		}
	}

	return readyPodsCount
}

func countCompleted(pods pods) int8 {
	completedPodsCount := int8(0)
	for i := range pods.Items {
		pod := pods.Items[i]
		for j := range pod.Status.Conditions {
			condition := pod.Status.Conditions[j]
			if condition.Type != "Initialized" {
				continue
			}
			if condition.Status == "True" && condition.Reason == "PodCompleted" {
				completedPodsCount++
				break
			}
		}
	}
	return completedPodsCount
}
