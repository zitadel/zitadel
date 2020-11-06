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

	rootSecret := "client-root"
	secretMode := int32(0777)
	replicas := int32(replicaCount)
	runAsUser := int64(1000)
	runAsNonRoot := true
	certMountPath := "/dbsecrets"
	containerName := "zitadel"

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
	}}
	volMounts := []corev1.VolumeMount{
		{Name: secretName, MountPath: secretPath},
		{Name: consoleCMName, MountPath: "/console/environment.json", SubPath: "environment.json"},
		{Name: rootSecret, MountPath: certMountPath + "/ca.crt", SubPath: "ca.crt"},
	}

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
		volMounts = append(volMounts, corev1.VolumeMount{
			Name: internalName,
			//ReadOnly:  true,
			MountPath: certMountPath + "/client." + user + ".crt",
			SubPath:   "client." + user + ".crt",
		})
		volMounts = append(volMounts, corev1.VolumeMount{
			Name: internalName,
			//ReadOnly:  true,
			MountPath: certMountPath + "/client." + user + ".key",
			SubPath:   "client." + user + ".key",
		})
	}

	envVars := []corev1.EnvVar{
		{Name: "POD_IP",
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					FieldPath: "status.podIP",
				},
			}},
		{Name: "CHAT_URL",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{Name: secretVarsName},
					Key:                  "ZITADEL_GOOGLE_CHAT_URL",
				},
			}},
		{Name: "TWILIO_TOKEN",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{Name: secretVarsName},
					Key:                  "ZITADEL_TWILIO_AUTH_TOKEN",
				},
			}},
		{Name: "TWILIO_SERVICE_SID",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{Name: secretVarsName},
					Key:                  "ZITADEL_TWILIO_SID",
				},
			}},
		{Name: "SMTP_PASSWORD",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{Name: secretVarsName},
					Key:                  "ZITADEL_EMAILAPPKEY",
				},
			}},
	}

	for _, user := range users {
		envVars = append(envVars, corev1.EnvVar{
			Name: "CR_" + strings.ToUpper(user) + "_PASSWORD",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{Name: secretPasswordsName},
					Key:                  user,
				},
			},
		})
	}

	deployName := "zitadel"
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
					SecurityContext: &corev1.PodSecurityContext{
						RunAsUser:    &runAsUser,
						RunAsNonRoot: &runAsNonRoot,
					},
					Containers: []corev1.Container{{
						Resources: corev1.ResourceRequirements(*resources),
						Lifecycle: &corev1.Lifecycle{
							PostStart: &corev1.Handler{
								Exec: &corev1.ExecAction{
									// TODO: until proper fix of https://github.com/kubernetes/kubernetes/issues/2630
									Command: []string{"sh", "-c",
										"mkdir -p " + certPath + "/ && cp " + certMountPath + "/* " + certPath + "/ && chmod 400 " + certPath + "/*"},
								},
							},
						},
						Args: []string{"start"},
						SecurityContext: &corev1.SecurityContext{
							RunAsUser:    &runAsUser,
							RunAsNonRoot: &runAsNonRoot,
						},
						Name:            containerName,
						Image:           zitadelImage,
						ImagePullPolicy: "IfNotPresent",
						Ports: []corev1.ContainerPort{
							{Name: "grpc", ContainerPort: 50001},
							{Name: "http", ContainerPort: 50002},
							{Name: "ui", ContainerPort: 50003},
						},
						Env: envVars,
						EnvFrom: []corev1.EnvFromSource{
							{ConfigMapRef: &corev1.ConfigMapEnvSource{
								LocalObjectReference: corev1.LocalObjectReference{Name: cmName},
							}}},
						VolumeMounts: volMounts,
						/*LivenessProbe: &corev1.Probe{
							Handler: corev1.Handler{
								HTTPGet: &corev1.HTTPGetAction{
									Path:   "/healthz",
									Port:   intstr.Parse("http"),
									Scheme: "HTTP",
								},
							},
							PeriodSeconds:       5,
							FailureThreshold:    2,
							InitialDelaySeconds: 60,
						},*/
						ReadinessProbe: &corev1.Probe{
							Handler: corev1.Handler{
								HTTPGet: &corev1.HTTPGetAction{
									Path:   "/ready",
									Port:   intstr.Parse("http"),
									Scheme: "HTTP",
								},
							},
							PeriodSeconds:    5,
							FailureThreshold: 2,
						},
					}},
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
		if err := k8sClient.WaitUntilStatefulsetIsReady(namespace, deployName, true, true, 1); err != nil {
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
