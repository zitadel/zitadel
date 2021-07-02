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
	secretKey := "testKey"
	secretName := "testSecretName"
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
						Name:      internalSecretName,
						MountPath: certPath,
					}, {
						Name:      secretKey,
						SubPath:   secretKey,
						MountPath: secretPath,
					}},
					ImagePullPolicy: corev1.PullAlways,
				}},
				Volumes: []corev1.Volume{{
					Name: internalSecretName,
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName:  rootSecretName,
							DefaultMode: helpers.PointerInt32(defaultMode),
						},
					},
				}, {
					Name: secretKey,
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName: secretName,
						},
					},
				}},
			},
		},
	}

	assert.Equal(t, equals, getJobSpecDef(nodeselector, tolerations, secretName, secretKey, backupName, command, image))
}

func TestBackup_JobSpec2(t *testing.T) {
	nodeselector := map[string]string{"test2": "test2"}
	tolerations := []corev1.Toleration{
		{Key: "testKey2", Operator: "testOp2"}}
	backupName := "testName2"
	command := "test2"
	secretKey := "testKey2"
	secretName := "testSecretName2"
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
						Name:      internalSecretName,
						MountPath: certPath,
					}, {
						Name:      secretKey,
						SubPath:   secretKey,
						MountPath: secretPath,
					}},
					ImagePullPolicy: corev1.PullAlways,
				}},
				Volumes: []corev1.Volume{{
					Name: internalSecretName,
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName:  rootSecretName,
							DefaultMode: helpers.PointerInt32(defaultMode),
						},
					},
				}, {
					Name: secretKey,
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName: secretName,
						},
					},
				}},
			},
		},
	}

	assert.Equal(t, equals, getJobSpecDef(nodeselector, tolerations, secretName, secretKey, backupName, command, image))
}
