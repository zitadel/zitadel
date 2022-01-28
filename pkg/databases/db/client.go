package db

/*
type Client interface {
	GetConnectionInfo(monitor mntr.Monitor, k8sClient kubernetes.ClientInt) (Connection, error)
	DeleteUser(monitor mntr.Monitor, user string, k8sClient kubernetes.ClientInt) error
	AddUser(monitor mntr.Monitor, user string, k8sClient kubernetes.ClientInt) error
	ListUsers(monitor mntr.Monitor, k8sClient kubernetes.ClientInt) ([]string, error)
}

*/

type Connection interface {
	Host() string
	Port() string
	User() string
	PasswordSecret() (string, string)
	ConnectionParams(certsDir string) string
}
