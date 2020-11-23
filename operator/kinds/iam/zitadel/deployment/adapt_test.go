package deployment

import (
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
	"testing"
)

func TestDeployment_Adapt(t *testing.T) {
	monitor := mntr.Monitor{}
	version := "test"
	namespace := "test"
	labels := map[string]string{"test": "test"}
	replicaCount := 1
	nodeSelector := map[string]string{"test": "test"}
	secretVarsName := "testVars"
	secretPasswordsName := "testPasswords"
	secretPath := "testSecretPath"
	certPath := "testCert"
	secretName := "testSecret"
	consoleCMName := "testConsoleCM"
	cmName := "testCM"
	users := []string{"test"}
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
			Name:        deployName,
			Namespace:   namespace,
			Labels:      labels,
			Annotations: annotations,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: helpers.PointerInt32(int32(replicaCount)),
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
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
					Labels:      labels,
					Annotations: annotations,
				},
				Spec: corev1.PodSpec{
					NodeSelector: nodeSelector,
					Tolerations:  nil,
					Affinity:     nil,
					InitContainers: []corev1.Container{
						getInitContainer(
							rootSecret,
							dbSecrets,
							users,
							runAsUser,
						),
					},
					Containers: []corev1.Container{
						getContainer(
							containerName,
							version,
							runAsUser,
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
						),
					},
					Volumes: getVolumes(
						secretName,
						secretPasswordsName,
						consoleCMName,
						users,
					),
				},
			},
		},
	}
	k8sClient.EXPECT().ApplyDeployment(deploymentDef).Times(1)

	getConfigurationHashes := func(k8sClient kubernetes.ClientInt, queried map[string]interface{}) map[string]string {
		return map[string]string{"testHash": "test"}
	}
	migrationDone := func(k8sClient kubernetes.ClientInt) error {
		return nil
	}
	configurationDone := func(k8sClient kubernetes.ClientInt) error {
		return nil
	}

	query, _, _, _, _, err := AdaptFunc(
		monitor,
		version,
		namespace,
		labels,
		replicaCount,
		nil,
		cmName,
		certPath,
		secretName,
		secretPath,
		consoleCMName,
		secretVarsName,
		secretPasswordsName,
		users,
		nodeSelector,
		nil,
		resources,
		migrationDone,
		configurationDone,
		getConfigurationHashes,
	)
	assert.NoError(t, err)
	queried := map[string]interface{}{}
	ensure, err := query(k8sClient, queried)
	assert.NoError(t, err)
	assert.NoError(t, ensure(k8sClient))
}
