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
	accessKeyIDName := "testAKIN"
	accessKeyIDKey := "testAKIK"
	secretAccessKeyName := "testSAKN"
	secretAccessKeyKey := "testSAKK"
	sessionTokenName := "testSTN"
	sessionTokenKey := "testSTK"
	image := common.ZITADELCockroachImage.Reference("", version)
	runAsUser := int64(100)

	equals := batchv1.JobSpec{
		Template: corev1.PodTemplateSpec{
			Spec: corev1.PodSpec{
				RestartPolicy: corev1.RestartPolicyNever,
				NodeSelector:  nodeselector,
				Tolerations:   tolerations,
				InitContainers: []corev1.Container{
					common.GetInitContainer(
						"backup",
						internalSecretName,
						dbSecrets,
						[]string{"root"},
						runAsUser,
						image,
					),
				},
				Containers: []corev1.Container{{
					Name:  backupName,
					Image: image,
					Command: []string{
						"/bin/bash",
						"-c",
						command,
					},
					VolumeMounts: []corev1.VolumeMount{{
						Name:      dbSecrets,
						MountPath: certPath,
					}, {
						Name:      accessKeyIDKey,
						SubPath:   accessKeyIDKey,
						MountPath: accessKeyIDPath,
					}, {
						Name:      secretAccessKeyKey,
						SubPath:   secretAccessKeyKey,
						MountPath: secretAccessKeyPath,
					}, {
						Name:      sessionTokenKey,
						SubPath:   sessionTokenKey,
						MountPath: sessionTokenPath,
					}},
					ImagePullPolicy: corev1.PullIfNotPresent,
				}},
				Volumes: []corev1.Volume{{
					Name: internalSecretName,
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName:  rootSecretName,
							DefaultMode: helpers.PointerInt32(0444),
						},
					},
				}, {
					Name: accessKeyIDKey,
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName:  accessKeyIDName,
							DefaultMode: helpers.PointerInt32(0444),
						},
					},
				}, {
					Name: secretAccessKeyKey,
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName:  secretAccessKeyName,
							DefaultMode: helpers.PointerInt32(0444),
						},
					},
				}, {
					Name: sessionTokenKey,
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName:  sessionTokenName,
							DefaultMode: helpers.PointerInt32(0444),
						},
					},
				}, {
					Name: dbSecrets,
					VolumeSource: corev1.VolumeSource{
						EmptyDir: &corev1.EmptyDirVolumeSource{},
					},
				},
				},
			},
		},
	}

	assert.Equal(t, equals, getJobSpecDef(
		nodeselector,
		tolerations,
		accessKeyIDName,
		accessKeyIDKey,
		secretAccessKeyName,
		secretAccessKeyKey,
		sessionTokenName,
		sessionTokenKey,
		backupName,
		command,
		image,
		runAsUser,
	))
}

func TestBackup_JobSpec2(t *testing.T) {
	nodeselector := map[string]string{"test2": "test2"}
	tolerations := []corev1.Toleration{
		{Key: "testKey2", Operator: "testOp2"}}
	backupName := "testName2"
	version := "testVersion2"
	command := "test2"
	accessKeyIDName := "testAKIN2"
	accessKeyIDKey := "testAKIK2"
	secretAccessKeyName := "testSAKN2"
	secretAccessKeyKey := "testSAKK2"
	sessionTokenName := "testSTN2"
	sessionTokenKey := "testSTK2"
	image := common.ZITADELCockroachImage.Reference("", version)
	runAsUser := int64(100)

	equals := batchv1.JobSpec{
		Template: corev1.PodTemplateSpec{
			Spec: corev1.PodSpec{
				RestartPolicy: corev1.RestartPolicyNever,
				NodeSelector:  nodeselector,
				Tolerations:   tolerations,
				InitContainers: []corev1.Container{
					common.GetInitContainer(
						"backup",
						internalSecretName,
						dbSecrets,
						[]string{"root"},
						runAsUser,
						image,
					),
				},
				Containers: []corev1.Container{{
					Name:  backupName,
					Image: image,
					Command: []string{
						"/bin/bash",
						"-c",
						command,
					},
					VolumeMounts: []corev1.VolumeMount{{
						Name:      dbSecrets,
						MountPath: certPath,
					}, {
						Name:      accessKeyIDKey,
						SubPath:   accessKeyIDKey,
						MountPath: accessKeyIDPath,
					}, {
						Name:      secretAccessKeyKey,
						SubPath:   secretAccessKeyKey,
						MountPath: secretAccessKeyPath,
					}, {
						Name:      sessionTokenKey,
						SubPath:   sessionTokenKey,
						MountPath: sessionTokenPath,
					}},
					ImagePullPolicy: corev1.PullIfNotPresent,
				}},
				Volumes: []corev1.Volume{{
					Name: internalSecretName,
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName:  rootSecretName,
							DefaultMode: helpers.PointerInt32(0444),
						},
					},
				}, {
					Name: accessKeyIDKey,
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName:  accessKeyIDName,
							DefaultMode: helpers.PointerInt32(0444),
						},
					},
				}, {
					Name: secretAccessKeyKey,
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName:  secretAccessKeyName,
							DefaultMode: helpers.PointerInt32(0444),
						},
					},
				}, {
					Name: sessionTokenKey,
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName:  sessionTokenName,
							DefaultMode: helpers.PointerInt32(0444),
						},
					},
				}, {
					Name: dbSecrets,
					VolumeSource: corev1.VolumeSource{
						EmptyDir: &corev1.EmptyDirVolumeSource{},
					},
				},
				},
			},
		},
	}

	assert.Equal(t, equals, getJobSpecDef(
		nodeselector,
		tolerations,
		accessKeyIDName,
		accessKeyIDKey,
		secretAccessKeyName,
		secretAccessKeyKey,
		sessionTokenName,
		sessionTokenKey,
		backupName,
		command,
		image,
		runAsUser,
	))
}
