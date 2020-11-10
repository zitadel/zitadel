package deployment

import (
	"strings"

	"k8s.io/apimachinery/pkg/api/resource"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/kubernetes/k8s"
	"github.com/caos/orbos/pkg/kubernetes/resources/deployment"
	"github.com/caos/zitadel/operator"
	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	//zitadelImage can be found in github.com/caos/zitadel repo
	zitadelImage = "ghcr.io/caos/zitadel:0.101.0"
)

func AdaptFunc(
	monitor mntr.Monitor,
	namespace string,
	labels map[string]string,
	replicaCount int,
	affinity *k8s.Affinity,
	cmName string,
	certPath string,
	secretName string,
	secretPath string,
	consoleCMName string,
	secretVarsName string,
	secretPasswordsName string,
	users []string,
	nodeSelector map[string]string,
	tolerations []corev1.Toleration,
	resources *k8s.Resources,
	migrationDone operator.EnsureFunc,
	configurationDone operator.EnsureFunc,
	getConfigurationHashes func(k8sClient *kubernetes.Client) map[string]string,
) (
	operator.QueryFunc,
	operator.DestroyFunc,
	operator.EnsureFunc,
	func(replicaCount int) operator.EnsureFunc,
	operator.EnsureFunc,
	error,
) {
	internalMonitor := monitor.WithField("component", "deployment")

	secretMode := int32(384)
	replicas := int32(replicaCount)
	runAsUser := int64(1000)

	rootSecret := "client-root"
	dbSecrets := "db-secrets"

	deployName := "zitadel"
	containerName := "zitadel"
	maxUnavailable := intstr.FromInt(1)
	maxSurge := intstr.FromInt(1)

	if resources == nil {
		resources = &k8s.Resources{
			Limits: corev1.ResourceList{
				corev1.ResourceCPU:    resource.MustParse("2"),
				corev1.ResourceMemory: resource.MustParse("2Gi"),
			},
			Requests: corev1.ResourceList{
				corev1.ResourceCPU:    resource.MustParse("500m"),
				corev1.ResourceMemory: resource.MustParse("512Mi"),
			},
		}
	}

	volumes := []corev1.Volume{{
		Name: secretName,
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName: secretName,
			},
		},
	}, {
		Name: rootSecret,
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName:  "cockroachdb.client.root",
				DefaultMode: &secretMode,
			},
		},
	}, {
		Name: secretPasswordsName,
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName: secretPasswordsName,
			},
		},
	}, {
		Name: consoleCMName,
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{Name: consoleCMName},
			},
		},
	}, {
		Name: dbSecrets,
		VolumeSource: corev1.VolumeSource{
			EmptyDir: &corev1.EmptyDirVolumeSource{},
		},
	}}

	for _, user := range users {
		userReplaced := strings.ReplaceAll(user, "_", "-")
		internalName := "client-" + userReplaced
		volumes = append(volumes, corev1.Volume{
			Name: internalName,
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName:  "cockroachdb.client." + userReplaced,
					DefaultMode: &secretMode,
				},
			},
		})
	}

	deploymentDef := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        deployName,
			Namespace:   namespace,
			Labels:      labels,
			Annotations: map[string]string{},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.RollingUpdateDeploymentStrategyType,
				RollingUpdate: &appsv1.RollingUpdateDeployment{
					MaxUnavailable: &maxUnavailable,
					MaxSurge:       &maxSurge,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      labels,
					Annotations: map[string]string{},
				},
				Spec: corev1.PodSpec{
					NodeSelector: nodeSelector,
					Tolerations:  tolerations,
					Affinity:     affinity.K8s(),
					InitContainers: []corev1.Container{
						getInitContainer(
							rootSecret,
							dbSecrets,
							users,
							runAsUser,
						),
					},
					Containers: []corev1.Container{
						getContainer(
							containerName,
							runAsUser,
							true,
							resources,
							cmName,
							certPath,
							secretName,
							secretPath,
							consoleCMName,
							secretVarsName,
							secretPasswordsName,
							users,
							dbSecrets,
						),
					},
					Volumes: volumes,
				},
			},
		},
	}

	destroy, err := deployment.AdaptFuncToDestroy(namespace, deployName)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	destroyers := []operator.DestroyFunc{
		operator.ResourceDestroyToZitadelDestroy(destroy),
	}

	checkDeployRunning := func(k8sClient *kubernetes.Client) error {
		internalMonitor.Info("waiting for deployment to be running")
		if err := k8sClient.WaitUntilDeploymentReady(namespace, deployName, true, false, 60); err != nil {
			internalMonitor.Error(errors.Wrap(err, "error while waiting for deployment to be running"))
			return err
		}
		internalMonitor.Info("deployment is running")
		return nil
	}

	checkDeployNotReady := func(k8sClient *kubernetes.Client) error {
		internalMonitor.Info("checking for deployment to not be ready")
		if err := k8sClient.WaitUntilDeploymentReady(namespace, deployName, true, true, 1); err != nil {
			internalMonitor.Info("deployment is not ready")
			return nil
		}
		internalMonitor.Info("deployment is ready")
		return errors.New("deployment is ready")
	}

	return func(k8sClient *kubernetes.Client, queried map[string]interface{}) (operator.EnsureFunc, error) {
			hashes := getConfigurationHashes(k8sClient)
			if hashes != nil && len(hashes) != 0 {
				for k, v := range hashes {
					deploymentDef.Annotations[k] = v
					deploymentDef.Spec.Template.Annotations[k] = v
				}
			}

			query, err := deployment.AdaptFuncToEnsure(deploymentDef)
			if err != nil {
				return nil, err
			}

			queriers := []operator.QueryFunc{
				operator.EnsureFuncToQueryFunc(migrationDone),
				operator.EnsureFuncToQueryFunc(configurationDone),
				operator.ResourceQueryToZitadelQuery(query),
			}

			return operator.QueriersToEnsureFunc(internalMonitor, false, queriers, k8sClient, queried)
		},
		operator.DestroyersToDestroyFunc(internalMonitor, destroyers),
		func(k8sClient *kubernetes.Client) error {
			internalMonitor.Info("waiting for deployment to be ready")
			if err := k8sClient.WaitUntilDeploymentReady(namespace, deployName, true, true, 300); err != nil {
				internalMonitor.Error(errors.Wrap(err, "error while waiting for deployment to be ready"))
				return err
			}
			internalMonitor.Info("deployment is ready")
			return nil
		},
		func(replicaCount int) operator.EnsureFunc {
			return func(k8sClient *kubernetes.Client) error {
				internalMonitor.Info("Scaling deployment")
				return k8sClient.ScaleDeployment(namespace, deployName, replicaCount)
			}
		},
		func(k8sClient *kubernetes.Client) error {
			if err := checkDeployRunning(k8sClient); err != nil {
				return err
			}

			if err := checkDeployNotReady(k8sClient); err != nil {
				return nil
			}

			command := "/zitadel setup"

			if err := k8sClient.ExecInPodOfDeployment(namespace, deployName, containerName, command); err != nil {
				return err
			}
			return nil
		},
		nil
}
