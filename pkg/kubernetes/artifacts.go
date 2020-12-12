package kubernetes

import (
	"fmt"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/labels"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	rbac "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	mach "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func EnsureZitadelOperatorArtifacts(
	monitor mntr.Monitor,
	apiLabels *labels.API,
	client kubernetes.ClientInt,
	version string,
	nodeselector map[string]string,
	tolerations []core.Toleration,
) error {

	monitor.WithFields(map[string]interface{}{
		"zitadel": version,
	}).Debug("Ensuring zitadel artifacts")

	nameLabels := labels.MustForName(labels.MustForComponent(apiLabels, "operator"), "zitadel-operator")
	k8sNameLabels := labels.MustK8sMap(nameLabels)
	k8sPodSelector := labels.MustK8sMap(labels.DeriveNameSelector(nameLabels, false))

	if version == "" {
		return nil
	}

	if err := client.ApplyServiceAccount(&core.ServiceAccount{
		ObjectMeta: mach.ObjectMeta{
			Name:      nameLabels.Name(),
			Namespace: "caos-system",
		},
	}); err != nil {
		return err
	}

	if err := client.ApplyClusterRole(&rbac.ClusterRole{
		ObjectMeta: mach.ObjectMeta{
			Name:   nameLabels.Name(),
			Labels: k8sNameLabels,
		},
		Rules: []rbac.PolicyRule{{
			APIGroups: []string{"*"},
			Resources: []string{"*"},
			Verbs:     []string{"*"},
		}},
	}); err != nil {
		return err
	}

	if err := client.ApplyClusterRoleBinding(&rbac.ClusterRoleBinding{
		ObjectMeta: mach.ObjectMeta{
			Name:   nameLabels.Name(),
			Labels: k8sNameLabels,
		},

		RoleRef: rbac.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     nameLabels.Name(),
		},
		Subjects: []rbac.Subject{{
			Kind:      "ServiceAccount",
			Name:      nameLabels.Name(),
			Namespace: "caos-system",
		}},
	}); err != nil {
		return err
	}

	deployment := &apps.Deployment{
		ObjectMeta: mach.ObjectMeta{
			Name:      nameLabels.Name(),
			Namespace: "caos-system",
			Labels:    k8sNameLabels,
		},
		Spec: apps.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &mach.LabelSelector{
				MatchLabels: k8sPodSelector,
			},
			Template: core.PodTemplateSpec{
				ObjectMeta: mach.ObjectMeta{
					Labels: labels.MustK8sMap(labels.AsSelectable(nameLabels)),
				},
				Spec: core.PodSpec{
					ServiceAccountName: nameLabels.Name(),
					Containers: []core.Container{{
						Name:            "zitadel",
						ImagePullPolicy: core.PullIfNotPresent,
						Image:           fmt.Sprintf("ghcr.io/caos/zitadel-operator:%s", version),
						Command:         []string{"/zitadelctl", "operator", "-f", "/secrets/orbconfig"},
						Args:            []string{},
						Ports: []core.ContainerPort{{
							Name:          "metrics",
							ContainerPort: 2112,
							Protocol:      "TCP",
						}},
						VolumeMounts: []core.VolumeMount{{
							Name:      "orbconfig",
							ReadOnly:  true,
							MountPath: "/secrets",
						}},
						Resources: core.ResourceRequirements{
							Limits: core.ResourceList{
								"cpu":    resource.MustParse("500m"),
								"memory": resource.MustParse("500Mi"),
							},
							Requests: core.ResourceList{
								"cpu":    resource.MustParse("250m"),
								"memory": resource.MustParse("250Mi"),
							},
						},
					}},
					NodeSelector: nodeselector,
					Tolerations:  tolerations,
					Volumes: []core.Volume{{
						Name: "orbconfig",
						VolumeSource: core.VolumeSource{
							Secret: &core.SecretVolumeSource{
								SecretName: "caos",
							},
						},
					}},
					TerminationGracePeriodSeconds: int64Ptr(10),
				},
			},
		},
	}
	if err := client.ApplyDeployment(deployment, true); err != nil {
		return err
	}
	monitor.WithFields(map[string]interface{}{
		"version": version,
	}).Debug("Zitadel Operator deployment ensured")

	return nil
}

func ScaleZitadelOperator(
	monitor mntr.Monitor,
	client *kubernetes.Client,
	replicaCount int,
) error {
	monitor.Debug("Scaling zitadel-operator")
	return client.ScaleDeployment("caos-system", "zitadel-operator", replicaCount)
}

func int32Ptr(i int32) *int32 { return &i }
func int64Ptr(i int64) *int64 { return &i }
