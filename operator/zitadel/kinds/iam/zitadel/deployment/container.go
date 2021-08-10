package deployment

import (
	"sort"
	"strings"

	"github.com/caos/zitadel/operator/common"

	"github.com/caos/orbos/pkg/kubernetes/k8s"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func GetContainer(
	containerName string,
	version string,
	runAsUser int64,
	runAsNonRoot bool,
	resources *k8s.Resources,
	cmName string,
	certPath string,
	secretName string,
	secretPath string,
	consoleCMName string,
	secretVarsName string,
	secretPasswordsName string,
	users []string,
	dbSecrets string,
	command string,
	customImageRegistry string,
) corev1.Container {

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
		{Name: "ZITADEL_ASSET_STORAGE_ACCESS_KEY_ID",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{Name: secretVarsName},
					Key:                  "ZITADEL_ASSET_STORAGE_ACCESS_KEY_ID",
				},
			}},
		{Name: "ZITADEL_ASSET_STORAGE_SECRET_ACCESS_KEY",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{Name: secretVarsName},
					Key:                  "ZITADEL_ASSET_STORAGE_SECRET_ACCESS_KEY",
				},
			}},
		{Name: "SENTRY_DSN",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{Name: secretVarsName},
					Key:                  "SENTRY_DSN",
				},
			}},
	}

	sort.Strings(users)
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

	volMounts := []corev1.VolumeMount{
		{Name: secretName, MountPath: secretPath},
		{Name: consoleCMName, MountPath: "/console/environment.json", SubPath: "environment.json"},
		{Name: dbSecrets, MountPath: certPath},
	}

	return corev1.Container{
		Resources: corev1.ResourceRequirements(*resources),
		//Command:   []string{"/bin/sh", "-c"},
		//Args:      []string{"tail -f /dev/null;"},
		Args: []string{command},
		SecurityContext: &corev1.SecurityContext{
			RunAsUser: &runAsUser,
		},
		Name:            containerName,
		Image:           common.ZITADELImage.Reference(customImageRegistry, version),
		ImagePullPolicy: corev1.PullIfNotPresent,
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
		TerminationMessagePolicy: corev1.TerminationMessageReadFile,
		TerminationMessagePath:   "/dev/termination-log",
	}
}
