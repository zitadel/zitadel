package deployment

import (
	"github.com/caos/zitadel/operator/helpers"
	"sort"
	"strings"

	"github.com/caos/zitadel/operator/common"

	corev1 "k8s.io/api/core/v1"
)

const (
	certMountPath     = "/dbsecrets"
	certTempMountPath = "/tmp/dbsecrets"
)

func GetInitContainer(
	reason string,
	rootSecret string,
	dbSecrets string,
	users []string,
	runAsUser int64,
	customImageRegistry string,
	version string,
) corev1.Container {
	initVolumeMounts := []corev1.VolumeMount{
		{Name: rootSecret, MountPath: certMountPath + "/client_root"},
		{Name: dbSecrets, MountPath: certTempMountPath},
	}
	containsRoot := false
	for _, user := range users {
		if user == "root" {
			containsRoot = true
		}
	}
	copySecrets := append([]string{}, "cp "+certMountPath+"/client_root/ca.crt "+certTempMountPath+"/ca.crt")
	if containsRoot {
		copySecrets = append(copySecrets, "cp "+certMountPath+"/client_root/client.root.crt "+certTempMountPath+"/client.root.crt")
		copySecrets = append(copySecrets, "cp "+certMountPath+"/client_root/client.root.key "+certTempMountPath+"/client.root.key")
	}

	sort.Strings(users)
	for _, user := range users {
		if user == "root" {
			continue
		}
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
		//"chown -R "+strconv.FormatInt(runAsUser, 10)+":"+strconv.FormatInt(runAsUser, 10)+" "+certTempMountPath+"",
		"chmod 0600 "+certTempMountPath+"/*",
	)

	name := "fix-permissions"
	if reason != "" {
		name = name + "-" + reason
	}

	return corev1.Container{
		Name:                     name,
		Image:                    common.ZITADELCockroachImage.Reference(customImageRegistry, version),
		Command:                  []string{"/bin/sh", "-c"},
		Args:                     []string{strings.Join(initCommands, " && ")},
		VolumeMounts:             initVolumeMounts,
		TerminationMessagePolicy: corev1.TerminationMessageReadFile,
		ImagePullPolicy:          corev1.PullIfNotPresent,
		TerminationMessagePath:   "/dev/termination-log",
		SecurityContext: &corev1.SecurityContext{
			RunAsUser:  helpers.PointerInt64(runAsUser),
			RunAsGroup: helpers.PointerInt64(InitRunAsUser),
		},
	}
}
