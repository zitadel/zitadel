package backup

import (
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

	equals := batchv1.JobSpec{
		Template: corev1.PodTemplateSpec{
			Spec: corev1.PodSpec{
				RestartPolicy: corev1.RestartPolicyNever,
				NodeSelector:  nodeselector,
				Tolerations:   tolerations,
				Containers: []corev1.Container{{
					Name:  backupName,
					Image: image + ":" + version,
					Command: []string{
						"/bin/bash",
						"-c",
						command,
					},
					VolumeMounts: []corev1.VolumeMount{{
						Name:      internalSecretName,
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
					Name: accessKeyIDKey,
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName: accessKeyIDName,
						},
					},
				}, {
					Name: secretAccessKeyKey,
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName: secretAccessKeyName,
						},
					},
				}, {
					Name: sessionTokenKey,
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName: sessionTokenName,
						},
					},
				}},
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
		version,
		command))
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

	equals := batchv1.JobSpec{
		Template: corev1.PodTemplateSpec{
			Spec: corev1.PodSpec{
				RestartPolicy: corev1.RestartPolicyNever,
				NodeSelector:  nodeselector,
				Tolerations:   tolerations,
				Containers: []corev1.Container{{
					Name:  backupName,
					Image: image + ":" + version,
					Command: []string{
						"/bin/bash",
						"-c",
						command,
					},
					VolumeMounts: []corev1.VolumeMount{{
						Name:      internalSecretName,
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
					Name: accessKeyIDKey,
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName: accessKeyIDName,
						},
					},
				}, {
					Name: secretAccessKeyKey,
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName: secretAccessKeyName,
						},
					},
				}, {
					Name: sessionTokenKey,
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName: sessionTokenName,
						},
					},
				}},
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
		version,
		command,
	))
}
