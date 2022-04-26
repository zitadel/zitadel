package db

import (
	"crypto/rsa"
	"fmt"
	"sort"
	"strings"

	"github.com/caos/orbos/pkg/labels"

	"github.com/zitadel/zitadel/operator/common"

	"github.com/zitadel/zitadel/operator/helpers"
	corev1 "k8s.io/api/core/v1"
)

const (
	CACert = "ca.crt"
	CAKey  = "ca.key"
)

type Connection interface {
	Host() string
	Port() string
	User() string
	PasswordSecret() (*labels.Selectable, string)
	SSL() *SSL
	Options() string
	CACert() []byte
	CAKey() *rsa.PrivateKey
}

type SSL struct {
	RootCert       bool
	UserCertAndKey bool
}

func CertsSecret(user string) string {
	return fmt.Sprintf("cockroachdb.client.%s", user)
}

func UserCert(user string) string {
	return fmt.Sprintf("client.%s.crt", user)
}

func UserKey(user string) string {
	return fmt.Sprintf("client.%s.key", user)
}

func InitChownCerts(
	customImageRegistry string,
	permissions string,
	users []string,
	to corev1.VolumeMount,
) (
	volumes []corev1.Volume,
	init corev1.Container,
) {

	sort.Strings(users)
	volumeMounts := make([]corev1.VolumeMount, len(users)+1)
	volumeMounts[0] = to
	volumes = make([]corev1.Volume, len(users)+1)
	volumes[0] = corev1.Volume{
		Name: to.Name,
		VolumeSource: corev1.VolumeSource{
			EmptyDir: &corev1.EmptyDirVolumeSource{},
		},
	}
	copyCmd := make([]string, len(users))
	for i := range users {
		user := users[i]
		volumeName := user + "-certs"
		mountPath := "/originalCerts/" + volumeName
		copyCmd[i] = fmt.Sprintf("cp %s/* %s/", mountPath, to.MountPath)
		volumeMounts[i+1] = corev1.VolumeMount{
			Name:      volumeName,
			MountPath: mountPath,
		}
		volumes[i+1] = corev1.Volume{
			Name: volumeName,
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName:  "cockroachdb.client." + user,
					DefaultMode: helpers.PointerInt32(0400),
				},
			},
		}
	}

	return volumes, corev1.Container{
		Name:         "chown",
		Image:        common.AlpineImage.Reference(customImageRegistry),
		Command:      []string{"sh", "-c"},
		Args:         []string{fmt.Sprintf("%s && chown -R %s %s/* && chmod 600 %s/*", strings.Join(copyCmd, " && "), permissions, to.MountPath, to.MountPath)},
		VolumeMounts: volumeMounts,
	}
}
