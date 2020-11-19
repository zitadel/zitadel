package configuration

import (
	"errors"
	"github.com/caos/orbos/mntr"
	kubernetesmock "github.com/caos/orbos/pkg/kubernetes/mock"
	"github.com/caos/zitadel/operator/kinds/iam/zitadel/database"
	databasemock "github.com/caos/zitadel/operator/kinds/iam/zitadel/database/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func SetConfigMap(
	k8sClient *kubernetesmock.MockClientInt,
	namespace string,
	cmName string,
	labels map[string]string,
	queried map[string]interface{},
	desired *Configuration,
	users map[string]string,
	certPath, secretPath, googleServiceAccountJSONPath, zitadelKeysPath string,
) {

	k8sClient.EXPECT().ApplyConfigmap(&corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      cmName,
			Labels:    labels,
		},
		Data: literalsConfigMap(desired, users, certPath, secretPath, googleServiceAccountJSONPath, zitadelKeysPath, queried),
	})
}

func SetSecretVars(
	k8sClient *kubernetesmock.MockClientInt,
	namespace string,
	secretVarsName string,
	labels map[string]string,
	desired *Configuration,
) {

	k8sClient.EXPECT().ApplySecret(&corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      secretVarsName,
			Labels:    labels,
		},
		Type:       "Opaque",
		StringData: literalsSecretVars(desired),
	}).Times(1)
}
func SetConsoleCM(
	k8sClient *kubernetesmock.MockClientInt,
	namespace string,
	consoleCMName string,
	labels map[string]string,
	getClientID func() string,
	desired *Configuration,
) {

	k8sClient.EXPECT().GetConfigMap(namespace, consoleCMName).Times(2).Return(nil, errors.New("Not Found"))
	consoleCM := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      consoleCMName,
			Labels:    labels,
		},
		Data: literalsConsoleCM(getClientID(), desired.DNS, k8sClient, namespace, consoleCMName),
	}
	k8sClient.EXPECT().ApplyConfigmap(consoleCM).Times(1)
}
func SetSecrets(
	k8sClient *kubernetesmock.MockClientInt,
	namespace string,
	secretName string,
	labels map[string]string,
	desired *Configuration,
) {
	k8sClient.EXPECT().ApplySecret(&corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      secretName,
			Labels:    labels,
		},
		Type:       "Opaque",
		StringData: literalsSecret(desired, googleServiceAccountJSONPath, zitadelKeysPath),
	}).Times(1)
}

func SetSecretPasswords(
	k8sClient *kubernetesmock.MockClientInt,
	namespace string,
	secretPasswordName string,
	labels map[string]string,
	users map[string]string,
) {

	k8sClient.EXPECT().ApplySecret(&corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      secretPasswordName,
			Labels:    labels,
		},
		Type:       "Opaque",
		StringData: users,
	}).Times(1)
}

func TestConfiguration_Adapt(t *testing.T) {
	k8sClient := kubernetesmock.NewMockClientInt(gomock.NewController(t))
	dbClient := databasemock.NewMockClientInt(gomock.NewController(t))

	monitor := mntr.Monitor{Fields: map[string]interface{}{"component": "configuration"}}
	namespace := "test"
	labels := map[string]string{"test": "test"}
	cmName := "cm"
	secretName := "secret"
	consoleCMName := "console"
	secretVarsName := "vars"
	secretPasswordName := "passwords"
	getClientID := func() string { return "test" }
	certPath := "test"
	secretPath := "test"
	users := map[string]string{
		"migration":    "migration",
		"management":   "management",
		"auth":         "auth",
		"authz":        "authz",
		"adminapi":     "adminapi",
		"notification": "notification",
		"eventstore":   "eventstore",
	}
	queried := map[string]interface{}{}
	database.SetDatabaseInQueried(queried, &database.Current{
		Host:  "host",
		Port:  "port",
		Users: []string{},
	})
	for user := range users {
		dbClient.EXPECT().AddUser(monitor, user, k8sClient).Times(1)
	}

	SetConfigMap(
		k8sClient,
		namespace,
		cmName,
		labels,
		queried,
		desiredEmpty,
		users,
		certPath,
		secretPath,
		googleServiceAccountJSONPath,
		zitadelKeysPath)

	SetSecretVars(
		k8sClient,
		namespace,
		secretVarsName,
		labels,
		desiredEmpty,
	)

	SetConsoleCM(
		k8sClient,
		namespace,
		consoleCMName,
		labels,
		getClientID,
		desiredEmpty,
	)

	SetSecrets(
		k8sClient,
		namespace,
		secretName,
		labels,
		desiredEmpty,
	)

	SetSecretPasswords(
		k8sClient,
		namespace,
		secretPasswordName,
		labels,
		users,
	)

	query, _, _, err := AdaptFunc(
		monitor,
		namespace,
		labels,
		desiredEmpty,
		cmName,
		certPath,
		secretName,
		secretPath,
		consoleCMName,
		secretVarsName,
		secretPasswordName,
		users,
		getClientID,
		dbClient,
	)

	assert.NoError(t, err)
	ensure, err := query(k8sClient, queried)
	assert.NoError(t, err)
	assert.NoError(t, ensure(k8sClient))

}

func TestConfiguration_AdaptFull(t *testing.T) {
	k8sClient := kubernetesmock.NewMockClientInt(gomock.NewController(t))
	dbClient := databasemock.NewMockClientInt(gomock.NewController(t))

	monitor := mntr.Monitor{Fields: map[string]interface{}{"component": "configuration"}}
	namespace := "test2"
	labels := map[string]string{"test2": "test2"}
	cmName := "cm2"
	secretName := "secret2"
	consoleCMName := "console2"
	secretVarsName := "vars2"
	secretPasswordName := "passwords2"
	getClientID := func() string { return "test2" }
	certPath := "test2"
	secretPath := "test2"
	users := map[string]string{
		"migration":    "migration",
		"management":   "management",
		"auth":         "auth",
		"authz":        "authz",
		"adminapi":     "adminapi",
		"notification": "notification",
		"eventstore":   "eventstore",
	}

	queried := map[string]interface{}{}
	database.SetDatabaseInQueried(queried, &database.Current{
		Host:  "host2",
		Port:  "port2",
		Users: []string{},
	})
	for user := range users {
		dbClient.EXPECT().AddUser(monitor, user, k8sClient).Times(1)
	}

	SetConfigMap(
		k8sClient,
		namespace,
		cmName,
		labels,
		queried,
		desiredFull,
		users,
		certPath,
		secretPath,
		googleServiceAccountJSONPath,
		zitadelKeysPath)

	SetSecretVars(
		k8sClient,
		namespace,
		secretVarsName,
		labels,
		desiredFull,
	)

	SetConsoleCM(
		k8sClient,
		namespace,
		consoleCMName,
		labels,
		getClientID,
		desiredFull,
	)

	SetSecrets(
		k8sClient,
		namespace,
		secretName,
		labels,
		desiredFull,
	)

	SetSecretPasswords(
		k8sClient,
		namespace,
		secretPasswordName,
		labels,
		users,
	)

	query, _, _, err := AdaptFunc(
		monitor,
		namespace,
		labels,
		desiredFull,
		cmName,
		certPath,
		secretName,
		secretPath,
		consoleCMName,
		secretVarsName,
		secretPasswordName,
		users,
		getClientID,
		dbClient,
	)

	assert.NoError(t, err)
	ensure, err := query(k8sClient, queried)
	assert.NoError(t, err)
	assert.NoError(t, ensure(k8sClient))

}
