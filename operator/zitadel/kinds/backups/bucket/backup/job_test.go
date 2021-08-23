package backup

import (
	"testing"

	"github.com/caos/zitadel/operator/helpers"
	"github.com/stretchr/testify/assert"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
)

func TestBackup_JobSpec1(t *testing.T) {
	nodeselector := map[string]string{"test": "test"}
	tolerations := []corev1.Toleration{
		{Key: "testKey", Operator: "testOp"}}
	backupName := "testName"
	command := "test"
	saSecretName := "testSaSecretName"
	saSecretKey := "testSaSecretKey"
	configSecretName := "testConfigSecretName"
	configSecretKey := "testConfigSecretKey"
	image := "testImage"

	equals := batchv1.JobSpec{
		Template: corev1.PodTemplateSpec{
			Spec: corev1.PodSpec{
				RestartPolicy: corev1.RestartPolicyNever,
				NodeSelector:  nodeselector,
				Tolerations:   tolerations,
				Containers: []corev1.Container{{
					Name:  backupName,
					Image: image,
					Command: []string{
						"/bin/bash",
						"-c",
						command,
					},
					VolumeMounts: []corev1.VolumeMount{{
						Name:      saInternalSecretName,
						SubPath:   saSecretKey,
						MountPath: saSecretPath,
					}, {
						Name:      configInternalSecretName,
						SubPath:   configSecretKey,
						MountPath: configSecretPath,
					}},
					ImagePullPolicy: corev1.PullAlways,
				}},
				Volumes: []corev1.Volume{{
					Name: saInternalSecretName,
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName:  saSecretName,
							DefaultMode: helpers.PointerInt32(defaultMode),
						},
					},
				}, {
					Name: configInternalSecretName,
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName:  configSecretName,
							DefaultMode: helpers.PointerInt32(defaultMode),
						},
					},
				}},
			},
		},
	}

	assert.Equal(t, equals, getJobSpecDef(nodeselector, tolerations, saSecretName, saSecretKey, configSecretName, configSecretKey, backupName, command, image))
}

func TestBackup_JobSpec2(t *testing.T) {
	nodeselector := map[string]string{"test2": "test2"}
	tolerations := []corev1.Toleration{
		{Key: "testKey2", Operator: "testOp2"}}
	backupName := "testName2"
	command := "test2"
	saSecretName := "testSaSecretName2"
	saSecretKey := "testSaSecretKey2"
	configSecretName := "testConfigSecretName2"
	configSecretKey := "testConfigSecretKey2"
	image := "testImage2"

	equals := batchv1.JobSpec{
		Template: corev1.PodTemplateSpec{
			Spec: corev1.PodSpec{
				RestartPolicy: corev1.RestartPolicyNever,
				NodeSelector:  nodeselector,
				Tolerations:   tolerations,
				Containers: []corev1.Container{{
					Name:  backupName,
					Image: image,
					Command: []string{
						"/bin/bash",
						"-c",
						command,
					},
					VolumeMounts: []corev1.VolumeMount{{
						Name:      saInternalSecretName,
						SubPath:   saSecretKey,
						MountPath: saSecretPath,
					}, {
						Name:      configInternalSecretName,
						SubPath:   configSecretKey,
						MountPath: configSecretPath,
					}},
					ImagePullPolicy: corev1.PullAlways,
				}},
				Volumes: []corev1.Volume{{
					Name: saInternalSecretName,
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName:  saSecretName,
							DefaultMode: helpers.PointerInt32(defaultMode),
						},
					},
				}, {
					Name: configInternalSecretName,
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName:  configSecretName,
							DefaultMode: helpers.PointerInt32(defaultMode),
						},
					},
				}},
			},
		},
	}

	assert.Equal(t, equals, getJobSpecDef(nodeselector, tolerations, saSecretName, saSecretKey, configSecretName, configSecretKey, backupName, command, image))
}
