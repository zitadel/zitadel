package backup

import (
	"github.com/caos/zitadel/operator/common"
	"github.com/caos/zitadel/operator/helpers"
	"github.com/stretchr/testify/assert"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"testing"
)

func TestBackup_JobSpec1(t *testing.T) {
	nodeselector := map[string]string{"test": "test"}
	tolerations := []corev1.Toleration{
		{Key: "testKey", Operator: "testOp"}}
	backupName := "testName"
	version := "testVersion"
	command := "test"
	secretName := "testSecret"

	equals := batchv1.JobSpec{
		Template: corev1.PodTemplateSpec{
			Spec: corev1.PodSpec{
				RestartPolicy: corev1.RestartPolicyNever,
				NodeSelector:  nodeselector,
				Tolerations:   tolerations,
				Containers: []corev1.Container{{
					Name:  backupName,
					Image: common.BackupImage.Reference("", version),
					Command: []string{
						"/bin/bash",
						"-c",
						command,
					},
					VolumeMounts: []corev1.VolumeMount{{
						Name:      internalCertsSecretName,
						MountPath: certPath,
					}, {
						Name:      internalSecretName,
						MountPath: secretsPath,
					}},
					ImagePullPolicy: corev1.PullIfNotPresent,
				}},
				Volumes: []corev1.Volume{{
					Name: internalCertsSecretName,
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName:  rootSecretName,
							DefaultMode: helpers.PointerInt32(defaultMode),
						},
					},
				}, {
					Name: internalSecretName,
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName: secretName,
						},
					},
				}},
			},
		},
	}

	assert.Equal(t, equals, getJobSpecDef(
		nodeselector,
		tolerations,
		secretName,
		backupName,
		common.BackupImage.Reference("", version),
		command))
}

func TestBackup_JobSpec2(t *testing.T) {
	nodeselector := map[string]string{"test2": "test2"}
	tolerations := []corev1.Toleration{
		{Key: "testKey2", Operator: "testOp2"}}
	backupName := "testName2"
	version := "testVersion2"
	command := "test2"
	secretName := "testSecret"

	equals := batchv1.JobSpec{
		Template: corev1.PodTemplateSpec{
			Spec: corev1.PodSpec{
				RestartPolicy: corev1.RestartPolicyNever,
				NodeSelector:  nodeselector,
				Tolerations:   tolerations,
				Containers: []corev1.Container{{
					Name:  backupName,
					Image: common.BackupImage.Reference("", version),
					Command: []string{
						"/bin/bash",
						"-c",
						command,
					},
					VolumeMounts: []corev1.VolumeMount{{
						Name:      internalCertsSecretName,
						MountPath: certPath,
					}, {
						Name:      internalSecretName,
						MountPath: secretsPath,
					}},
					ImagePullPolicy: corev1.PullIfNotPresent,
				}},
				Volumes: []corev1.Volume{{
					Name: internalCertsSecretName,
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName:  rootSecretName,
							DefaultMode: helpers.PointerInt32(defaultMode),
						},
					},
				}, {
					Name: internalSecretName,
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName: secretName,
						},
					},
				}},
			},
		},
	}

	assert.Equal(t, equals, getJobSpecDef(
		nodeselector,
		tolerations,
		secretName,
		backupName,
		common.BackupImage.Reference("", version),
		command,
	))
}
