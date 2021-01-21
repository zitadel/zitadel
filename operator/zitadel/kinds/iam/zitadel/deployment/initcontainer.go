package deployment

import (
	corev1 "k8s.io/api/core/v1"
	"strconv"
	"strings"
)

const (
	certMountPath     = "/dbsecrets"
	certTempMountPath = "/tmp/dbsecrets"
)

func GetInitContainer(
	rootSecret string,
	dbSecrets string,
	users []string,
	runAsUser int64,
) corev1.Container {

	initVolumeMounts := []corev1.VolumeMount{
		{Name: rootSecret, MountPath: certMountPath + "/client_root"},
		{Name: dbSecrets, MountPath: certTempMountPath},
	}

	copySecrets := append([]string{}, "cp "+certMountPath+"/client_root/ca.crt "+certTempMountPath+"/ca.crt")
	for _, user := range users {
		userReplaced := strings.ReplaceAll(user, "_", "-")
		internalName := "client-" + userReplaced
		initVolumeMounts = append(initVolumeMounts, corev1.VolumeMount{
			Name: internalName,
			//ReadOnly:  true,
			MountPath: certMountPath + "/client_" + user,
		})
		copySecrets = append(copySecrets, "cp "+certMountPath+"/client_"+user+"/client."+user+".crt "+certTempMountPath+"/client."+user+".crt")
		copySecrets = append(copySecrets, "cp "+certMountPath+"/client_"+user+"/client."+user+".key "+certTempMountPath+"/client."+user+".key")
	}

	initCommands := append(
		copySecrets,
		"chown -R "+strconv.FormatInt(runAsUser, 10)+":"+strconv.FormatInt(runAsUser, 10)+" "+certTempMountPath+"",
		"chmod 0600 "+certTempMountPath+"/*",
	)

	return corev1.Container{
		Name:         "fix-permissions",
		Image:        "alpine:3.11",
		Command:      []string{"/bin/sh", "-c"},
		Args:         []string{strings.Join(initCommands, " && ")},
		VolumeMounts: initVolumeMounts,
	}
}
