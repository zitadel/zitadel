package deployment

import (
	"github.com/caos/zitadel/operator/helpers"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"testing"
)

func TestDeployment_Volumes(t *testing.T) {
	secretName := "testSecret"
	secretPasswordsName := "testPasswords"
	consoleCMName := "testCM"
	users := []string{"test"}

	equals := []corev1.Volume{{
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
				DefaultMode: helpers.PointerInt32(384),
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
	}, {Name: "client-test",
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName:  "cockroachdb.client.test",
				DefaultMode: helpers.PointerInt32(384),
			},
		},
	}}

	assert.ElementsMatch(t, equals, GetVolumes(secretName, secretPasswordsName, consoleCMName, users))

}

func TestDeployment_UserVolumes(t *testing.T) {
	users := []string{"test"}
	equals := []corev1.Volume{
		{Name: "client-test",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName:  "cockroachdb.client.test",
					DefaultMode: helpers.PointerInt32(384),
				},
			},
		}}

	assert.ElementsMatch(t, equals, userVolumes(users))

	users = []string{"te_st"}
	equals = []corev1.Volume{
		{Name: "client-te-st",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName:  "cockroachdb.client.te-st",
					DefaultMode: helpers.PointerInt32(384),
				},
			},
		}}

	assert.ElementsMatch(t, equals, userVolumes(users))

	users = []string{"test", "te-st"}
	equals = []corev1.Volume{
		{Name: "client-test",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName:  "cockroachdb.client.test",
					DefaultMode: helpers.PointerInt32(384),
				},
			},
		}, {Name: "client-te-st",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName:  "cockroachdb.client.te-st",
					DefaultMode: helpers.PointerInt32(384),
				},
			},
		}}

	assert.ElementsMatch(t, equals, userVolumes(users))

}
