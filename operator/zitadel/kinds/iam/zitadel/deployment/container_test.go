package deployment

import (
	"github.com/caos/orbos/pkg/kubernetes/k8s"
	"github.com/caos/zitadel/operator/helpers"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/intstr"
	"testing"
)

func TestDeployment_GetContainer(t *testing.T) {
	secretVarsName := "testVars"
	version := "test"
	secretPasswordsName := "testPasswords"
	secretPath := "testSecretPath"
	certPath := "testCert"
	secretName := "testSecret"
	consoleCMName := "testConsoleCM"
	cmName := "testCM"
	users := []string{"test"}

	resources := &k8s.Resources{
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("2"),
			corev1.ResourceMemory: resource.MustParse("2Gi"),
		},
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("500m"),
			corev1.ResourceMemory: resource.MustParse("512Mi"),
		},
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
		{Name: "ASSET_STORAGE_ACCESS_KEY_ID",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{Name: secretVarsName},
					Key:                  "ZITADEL_ASSET_STORAGE_ACCESS_KEY_ID",
				},
			}},
		{Name: "ASSET_STORAGE_SECRET_ACCESS_KEY",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{Name: secretVarsName},
					Key:                  "ZITADEL_ASSET_STORAGE_SECRET_ACCESS_KEY",
				},
			}},
		{
			Name: "CR_TEST_PASSWORD",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{Name: secretPasswordsName},
					Key:                  "test",
				},
			},
		},
	}

	volMounts := []corev1.VolumeMount{
		{Name: secretName, MountPath: secretPath},
		{Name: consoleCMName, MountPath: "/console/environment.json", SubPath: "environment.json"},
		{Name: dbSecrets, MountPath: certPath},
	}

	equals := corev1.Container{
		Resources: corev1.ResourceRequirements(*resources),
		//Command:   []string{"/bin/sh", "-c"},
		//Args:      []string{"tail -f /dev/null;"},
		Args: []string{"start"},
		SecurityContext: &corev1.SecurityContext{
			RunAsUser:    helpers.PointerInt64(RunAsUser),
			RunAsNonRoot: helpers.PointerBool(true),
		},
		Name:            containerName,
		Image:           zitadelImage + ":" + version,
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
		TerminationMessagePolicy: "File",
		TerminationMessagePath:   "/dev/termination-log",
	}

	container := GetContainer(
		containerName,
		version,
		RunAsUser,
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
		"start",
	)

	assert.Equal(t, equals, container)
}
