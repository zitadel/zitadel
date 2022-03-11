package db

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
)

type Client interface {
	GetConnectionInfo(monitor mntr.Monitor, k8sClient kubernetes.ClientInt) (string, string, error)
	DeleteUser(monitor mntr.Monitor, user string, k8sClient kubernetes.ClientInt) error
	AddUser(monitor mntr.Monitor, user string, k8sClient kubernetes.ClientInt) error
	ListUsers(monitor mntr.Monitor, k8sClient kubernetes.ClientInt) ([]string, error)
}
