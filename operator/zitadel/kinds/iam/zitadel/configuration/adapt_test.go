package configuration

import (
	"errors"
	"testing"

	"github.com/caos/orbos/mntr"
	kubernetesmock "github.com/caos/orbos/pkg/kubernetes/mock"
	"github.com/caos/orbos/pkg/labels"
	"github.com/caos/orbos/pkg/labels/mocklabels"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/database"
	databasemock "github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/database/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func SetConfigMap(
	k8sClient *kubernetesmock.MockClientInt,
	namespace string,
	cmName string,
	labels map[string]string,
	queried map[string]interface{},
	desired *Configuration,
	version *string,
	users map[string]string,
	certPath, secretPath, googleServiceAccountJSONPath, zitadelKeysPath string,
) {

	k8sClient.EXPECT().ApplyConfigmap(&corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      cmName,
			Labels:    labels,
		},
		Data: literalsConfigMap(desired, users, certPath, secretPath, googleServiceAccountJSONPath, zitadelKeysPath, version, queried),
	})
}

func SetSecretVars(
	t *testing.T,
	k8sClient *kubernetesmock.MockClientInt,
	namespace string,
	secretVarsName string,
	labels map[string]string,
	desired *Configuration,
) {

	literalsSV, err := literalsSecretVars(k8sClient, desired)
	assert.NoError(t, err)
	k8sClient.EXPECT().ApplySecret(&corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      secretVarsName,
			Labels:    labels,
		},
		Type:       "Opaque",
		StringData: literalsSV,
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
	t *testing.T,
	k8sClient *kubernetesmock.MockClientInt,
	namespace string,
	secretName string,
	labels map[string]string,
	desired *Configuration,
) {
	literalsS, err := literalsSecret(k8sClient, desired, googleServiceAccountJSONPath, zitadelKeysPath)
	assert.NoError(t, err)

	k8sClient.EXPECT().ApplySecret(&corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      secretName,
			Labels:    labels,
		},
		Type:       "Opaque",
		StringData: literalsS,
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
	dbClient := databasemock.NewMockClient(gomock.NewController(t))

	monitor := mntr.Monitor{Fields: map[string]interface{}{"component": "configuration"}}
	namespace := "test"
	cmName := "cm"
	secretName := "secret"
	consoleCMName := "console"
	secretVarsName := "vars"
	secretPasswordName := "passwords"
	getClientID := func() string { return "test" }
	certPath := "test"
	secretPath := "test"
	version := "test"
	users := map[string]string{
		"migration":    "migration",
		"management":   "management",
		"auth":         "auth",
		"authz":        "authz",
		"adminapi":     "adminapi",
		"notification": "notification",
		"eventstore":   "eventstore",
	}

	componentLabels := mocklabels.Component

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
		labels.MustForNameK8SMap(componentLabels, cmName),
		queried,
		desiredEmpty,
		&version,
		users,
		certPath,
		secretPath,
		googleServiceAccountJSONPath,
		zitadelKeysPath)

	SetSecretVars(
		t,
		k8sClient,
		namespace,
		secretVarsName,
		labels.MustForNameK8SMap(componentLabels, secretVarsName),
		desiredEmpty,
	)

	SetConsoleCM(
		k8sClient,
		namespace,
		consoleCMName,
		labels.MustForNameK8SMap(componentLabels, consoleCMName),
		getClientID,
		desiredEmpty,
	)

	SetSecrets(
		t,
		k8sClient,
		namespace,
		secretName,
		labels.MustForNameK8SMap(componentLabels, secretName),
		desiredEmpty,
	)

	SetSecretPasswords(
		k8sClient,
		namespace,
		secretPasswordName,
		labels.MustForNameK8SMap(componentLabels, secretPasswordName),
		users,
	)

	getQuery, _, _, err := AdaptFunc(
		monitor,
		componentLabels,
		namespace,
		desiredEmpty,
		cmName,
		certPath,
		secretName,
		secretPath,
		&version,
		consoleCMName,
		secretVarsName,
		secretPasswordName,
		dbClient,
		getClientID,
	)
	assert.NoError(t, err)
	query := getQuery(users)
	ensure, err := query(k8sClient, queried)
	assert.NoError(t, err)
	assert.NoError(t, ensure(k8sClient))

}

func TestConfiguration_AdaptFull(t *testing.T) {
	k8sClient := kubernetesmock.NewMockClientInt(gomock.NewController(t))
	dbClient := databasemock.NewMockClient(gomock.NewController(t))

	monitor := mntr.Monitor{Fields: map[string]interface{}{"component": "configuration"}}
	namespace := "test2"
	cmName := "cm2"
	secretName := "secret2"
	consoleCMName := "console2"
	secretVarsName := "vars2"
	secretPasswordName := "passwords2"
	getClientID := func() string { return "test2" }
	certPath := "test2"
	secretPath := "test2"
	version := "test"
	users := map[string]string{
		"migration":    "migration",
		"management":   "management",
		"auth":         "auth",
		"authz":        "authz",
		"adminapi":     "adminapi",
		"notification": "notification",
		"eventstore":   "eventstore",
	}

	componentLabels := mocklabels.Component

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
		labels.MustForNameK8SMap(componentLabels, cmName),
		queried,
		desiredFull,
		&version,
		users,
		certPath,
		secretPath,
		googleServiceAccountJSONPath,
		zitadelKeysPath)

	SetSecretVars(
		t,
		k8sClient,
		namespace,
		secretVarsName,
		labels.MustForNameK8SMap(componentLabels, secretVarsName),
		desiredFull,
	)

	SetConsoleCM(
		k8sClient,
		namespace,
		consoleCMName,
		labels.MustForNameK8SMap(componentLabels, consoleCMName),
		getClientID,
		desiredFull,
	)

	SetSecrets(
		t,
		k8sClient,
		namespace,
		secretName,
		labels.MustForNameK8SMap(componentLabels, secretName),
		desiredFull,
	)

	SetSecretPasswords(
		k8sClient,
		namespace,
		secretPasswordName,
		labels.MustForNameK8SMap(componentLabels, secretPasswordName),
		users,
	)

	getQuery, _, _, err := AdaptFunc(
		monitor,
		componentLabels,
		namespace,
		desiredFull,
		cmName,
		certPath,
		secretName,
		secretPath,
		&version,
		consoleCMName,
		secretVarsName,
		secretPasswordName,
		dbClient,
		getClientID,
	)

	assert.NoError(t, err)
	query := getQuery(users)
	ensure, err := query(k8sClient, queried)
	assert.NoError(t, err)
	assert.NoError(t, ensure(k8sClient))

}
