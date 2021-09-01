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
	backupSecretName := "testSecret"
	saSecretKey := "testSaKey"
	assetAKIDKey := "testAkidKey"
	assetSAKKey := "testSakKey"
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
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      certsInternalSecretName,
							MountPath: certPath,
						}, {
							Name:      saInternalSecretName,
							SubPath:   saSecretKey,
							MountPath: saSecretPath,
						}, {
							Name:      akidInternalSecretName,
							SubPath:   assetAKIDKey,
							MountPath: akidSecretPath,
						}, {
							Name:      sakInternalSecretName,
							SubPath:   assetSAKKey,
							MountPath: sakSecretPath,
						},
					},
					ImagePullPolicy: corev1.PullAlways,
				}},
				Volumes: []corev1.Volume{
					{
						Name: certsInternalSecretName,
						VolumeSource: corev1.VolumeSource{
							Secret: &corev1.SecretVolumeSource{
								SecretName:  rootSecretName,
								DefaultMode: helpers.PointerInt32(defaultMode),
							},
						},
					}, {
						Name: saInternalSecretName,
						VolumeSource: corev1.VolumeSource{
							Secret: &corev1.SecretVolumeSource{
								SecretName:  backupSecretName,
								DefaultMode: helpers.PointerInt32(defaultMode),
							},
						},
					}, {
						Name: akidInternalSecretName,
						VolumeSource: corev1.VolumeSource{
							Secret: &corev1.SecretVolumeSource{
								SecretName:  backupSecretName,
								DefaultMode: helpers.PointerInt32(defaultMode),
							},
						},
					}, {
						Name: sakInternalSecretName,
						VolumeSource: corev1.VolumeSource{
							Secret: &corev1.SecretVolumeSource{
								SecretName:  backupSecretName,
								DefaultMode: helpers.PointerInt32(defaultMode),
							},
						},
					},
				},
			},
		},
	}

	assert.Equal(t, equals, getJobSpecDef(
		nodeselector,
		tolerations,
		backupSecretName,
		saSecretKey,
		assetAKIDKey,
		assetSAKKey,
		backupName,
		command,
		image,
	))
}

func TestBackup_JobSpec2(t *testing.T) {
	nodeselector := map[string]string{"test2": "test2"}
	tolerations := []corev1.Toleration{
		{Key: "testKey2", Operator: "testOp2"}}
	backupName := "testName2"
	command := "test2"
	backupSecretName := "testSecret2"
	saSecretKey := "testSaKey2"
	assetAKIDKey := "testAkidKey2"
	assetSAKKey := "testSakKey2"
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
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      certsInternalSecretName,
							MountPath: certPath,
						}, {
							Name:      saInternalSecretName,
							SubPath:   saSecretKey,
							MountPath: saSecretPath,
						}, {
							Name:      akidInternalSecretName,
							SubPath:   assetAKIDKey,
							MountPath: akidSecretPath,
						}, {
							Name:      sakInternalSecretName,
							SubPath:   assetSAKKey,
							MountPath: sakSecretPath,
						},
					},
					ImagePullPolicy: corev1.PullAlways,
				}},
				Volumes: []corev1.Volume{
					{
						Name: certsInternalSecretName,
						VolumeSource: corev1.VolumeSource{
							Secret: &corev1.SecretVolumeSource{
								SecretName:  rootSecretName,
								DefaultMode: helpers.PointerInt32(defaultMode),
							},
						},
					}, {
						Name: saInternalSecretName,
						VolumeSource: corev1.VolumeSource{
							Secret: &corev1.SecretVolumeSource{
								SecretName:  backupSecretName,
								DefaultMode: helpers.PointerInt32(defaultMode),
							},
						},
					}, {
						Name: akidInternalSecretName,
						VolumeSource: corev1.VolumeSource{
							Secret: &corev1.SecretVolumeSource{
								SecretName:  backupSecretName,
								DefaultMode: helpers.PointerInt32(defaultMode),
							},
						},
					}, {
						Name: sakInternalSecretName,
						VolumeSource: corev1.VolumeSource{
							Secret: &corev1.SecretVolumeSource{
								SecretName:  backupSecretName,
								DefaultMode: helpers.PointerInt32(defaultMode),
							},
						},
					},
				},
			},
		},
	}

	assert.Equal(t, equals, getJobSpecDef(
		nodeselector,
		tolerations,
		backupSecretName,
		saSecretKey,
		assetAKIDKey,
		assetSAKKey,
		backupName,
		command,
		image,
	))
}
