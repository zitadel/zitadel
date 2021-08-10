package deployment

import (
	"testing"

	"github.com/caos/orbos/pkg/labels/mocklabels"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/kubernetes/k8s"
	kubernetesmock "github.com/caos/orbos/pkg/kubernetes/mock"
	"github.com/caos/zitadel/operator/helpers"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestDeployment_Adapt(t *testing.T) {
	monitor := mntr.Monitor{}
	imageVersion := "test"
	namespace := "test"

	replicaCount := 1
	nodeSelector := map[string]string{"test": "test"}
	secretVarsName := "testVars"
	secretPasswordsName := "testPasswords"
	secretPath := "testSecretPath"
	certPath := "testCert"
	secretName := "testSecret"
	consoleCMName := "testConsoleCM"
	cmName := "testCM"
	customImageRegistry := ""

	usersMap := map[string]string{"test": "test"}
	users := []string{}
	for _, user := range usersMap {
		users = append(users, user)
	}
	annotations := map[string]string{"testHash": "test"}

	k8sClient := kubernetesmock.NewMockClientInt(gomock.NewController(t))

	resources := &k8s.Resources{
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("2"),
			corev1.ResourceMemory: resource.MustParse("2Gi"),
		},
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("500m"),
			corev1.ResourceMemory: resource.MustParse("512Mi"),
		},
	}

	deploymentDef := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        mocklabels.NameVal,
			Namespace:   namespace,
			Labels:      mocklabels.NameMap,
			Annotations: annotations,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: helpers.PointerInt32(int32(replicaCount)),
			Selector: &metav1.LabelSelector{
				MatchLabels: mocklabels.ClosedNameSelectorMap,
			},
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.RollingUpdateDeploymentStrategyType,
				RollingUpdate: &appsv1.RollingUpdateDeployment{
					MaxUnavailable: helpers.IntToIntStr(1),
					MaxSurge:       helpers.IntToIntStr(1),
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      mocklabels.SelectableMap,
					Annotations: annotations,
				},
				Spec: corev1.PodSpec{
					NodeSelector: nodeSelector,
					Tolerations:  nil,
					Affinity:     nil,
					InitContainers: []corev1.Container{
						GetInitContainer(
							"zitadel",
							rootSecret,
							dbSecrets,
							users,
							RunAsUser,
							customImageRegistry,
							imageVersion,
						),
					},
					Containers: []corev1.Container{
						GetContainer(
							containerName,
							imageVersion,
							RunAsUser,
							true,
							resources,
							cmName,
							certPath,
							secretName,
							secretPath,
							consoleCMName,
							secretVarsName,
							secretPasswordsName,
							users,
							dbSecrets,
							"start",
							customImageRegistry,
						),
					},
					Volumes: GetVolumes(
						secretName,
						secretPasswordsName,
						consoleCMName,
						users,
					),
				},
			},
		},
	}
	k8sClient.EXPECT().ApplyDeployment(deploymentDef, false).Times(1)

	getConfigurationHashes := func(k8sClient kubernetes.ClientInt, queried map[string]interface{}, necessaryUsers map[string]string) (map[string]string, error) {
		return map[string]string{"testHash": "test"}, nil
	}
	migrationDone := func(k8sClient kubernetes.ClientInt) error {
		return nil
	}
	configurationDone := func(k8sClient kubernetes.ClientInt) error {
		return nil
	}
	setupDone := func(k8sClient kubernetes.ClientInt) error {
		return nil
	}

	getQuery, _, err := AdaptFunc(
		monitor,
		mocklabels.Name,
		mocklabels.ClosedNameSelector,
		false,
		&imageVersion,
		namespace,
		replicaCount,
		nil,
		cmName,
		certPath,
		secretName,
		secretPath,
		consoleCMName,
		secretVarsName,
		secretPasswordsName,
		nodeSelector,
		nil,
		resources,
		migrationDone,
		configurationDone,
		setupDone,
		customImageRegistry,
	)
	assert.NoError(t, err)
	queried := map[string]interface{}{}
	query := getQuery(usersMap, getConfigurationHashes)
	ensure, err := query(k8sClient, queried)
	assert.NoError(t, err)
	assert.NoError(t, ensure(k8sClient))
}
