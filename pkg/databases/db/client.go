package db

import (
	"fmt"
	"github.com/caos/orbos/pkg/labels"

	"github.com/caos/zitadel/operator/common"

	"github.com/caos/zitadel/operator/helpers"
	corev1 "k8s.io/api/core/v1"
)

const (
	CertsSecret  = "user-certs" // TODO: make dynamic
	CACert       = "ca.crt"
	RootUserCert = "client.root.crt"
	RootUserKey  = "client.root.key"
	UserCert     = "client.zitadel.crt"
	UserKey      = "client.zitadel.key"
)

type Connection interface {
	Host() string
	Port() string
	User() string
	PasswordSecret() (*labels.Selectable, string)
	SSL() *SSL
	Options() string
}

type SSL struct {
	RootCert       bool
	UserCertAndKey bool
}

func InitChownCerts(customImageRegistry string, permissions string, from, to corev1.VolumeMount) (source, chowned corev1.Volume, init corev1.Container) {

	return corev1.Volume{
			Name: from.Name,
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName:  CertsSecret,
					DefaultMode: helpers.PointerInt32(0400),
				},
			},
		}, corev1.Volume{
			Name: to.Name,
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{},
			},
		}, corev1.Container{
			Name:         "chown",
			Image:        common.AlpineImage.Reference(customImageRegistry),
			Command:      []string{"sh", "-c"},
			Args:         []string{fmt.Sprintf("cp %s/* %s/ && chown -R %s %s/* && chmod 600 %s/*", from.MountPath, to.MountPath, permissions, to.MountPath, to.MountPath)},
			VolumeMounts: []corev1.VolumeMount{from, to},
		}
}
