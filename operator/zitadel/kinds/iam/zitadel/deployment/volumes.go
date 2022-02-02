package deployment

import (
	"github.com/caos/zitadel/operator/helpers"
	corev1 "k8s.io/api/core/v1"
)

func GetVolumes(
	secretName string,
	secretPasswordsName string,
	consoleCMName string,
	users []string,
) []corev1.Volume {

	return []corev1.Volume{{
		Name: secretName,
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName:  secretName,
				DefaultMode: helpers.PointerInt32(420),
			},
		},
	}, {
		Name: secretPasswordsName,
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName:  secretPasswordsName,
				DefaultMode: helpers.PointerInt32(384),
			},
		},
	}, {
		Name: consoleCMName,
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{Name: consoleCMName},
				DefaultMode:          helpers.PointerInt32(420),
			},
		},
	}, {
		Name: dbSecrets,
		VolumeSource: corev1.VolumeSource{
			EmptyDir: &corev1.EmptyDirVolumeSource{},
		},
	}}
}

/*
func userVolumes(
	users []string,
) []corev1.Volume {
	volumes := make([]corev1.Volume, 0)

	sort.Strings(users)
	for _, user := range users {
		userReplaced := strings.ReplaceAll(user, "_", "-")
		internalName := "client-" + userReplaced
		volumes = append(volumes, corev1.Volume{
			Name: internalName,
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName:  "cockroachdb.client." + userReplaced,
					DefaultMode: helpers.PointerInt32(384),
				},
			},
		})
	}
	return volumes
}
*/
