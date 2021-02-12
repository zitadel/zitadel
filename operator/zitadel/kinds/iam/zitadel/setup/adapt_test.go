package setup

import (
	"testing"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/kubernetes/k8s"
	kubernetesmock "github.com/caos/orbos/pkg/kubernetes/mock"
	"github.com/caos/orbos/pkg/labels"
	"github.com/caos/orbos/pkg/labels/mocklabels"
	"github.com/caos/zitadel/operator/helpers"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/database"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/deployment"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	macherrs "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func TestSetup_AdaptFunc(t *testing.T) {
	monitor := mntr.Monitor{}
	client := kubernetesmock.NewMockClientInt(gomock.NewController(t))
	namespace := "test"
	reason := "test"
	users := []string{"test"}
	nodeselector := map[string]string{"test": "test"}
	tolerations := []corev1.Toleration{}
	dbHost := "test"
	dbPort := "test"

	version := "test"
	secretVarsName := "testVars"
	secretPasswordsName := "testPasswords"
	secretPath := "testSecretPath"
	certPath := "testCert"
	secretName := "testSecret"
	consoleCMName := "testConsoleCM"
	cmName := "testCM"
	annotations := map[string]string{"testHash": "test"}

	componentLabels := mocklabels.Component

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

	initContainers := []corev1.Container{deployment.GetInitContainer(rootSecret, dbSecrets, users, deployment.RunAsUser)}
	containers := []corev1.Container{deployment.GetContainer(
		containerName,
		version,
		deployment.RunAsUser,
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
		"setup",
	)}
	volumes := deployment.GetVolumes(secretName, secretPasswordsName, consoleCMName, users)

	jobName := labels.MustForName(componentLabels, jobNamePrefix+reason)
	jobDef := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:        jobName.Name(),
			Namespace:   namespace,
			Labels:      labels.MustK8sMap(jobName),
			Annotations: annotations,
		},
		Spec: batchv1.JobSpec{
			Completions: helpers.PointerInt32(1),
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: annotations,
				},
				Spec: corev1.PodSpec{
					NodeSelector:    nodeselector,
					Tolerations:     tolerations,
					InitContainers:  initContainers,
					Containers:      containers,
					SecurityContext: &corev1.PodSecurityContext{},

					RestartPolicy:                 "Never",
					DNSPolicy:                     "ClusterFirst",
					SchedulerName:                 "default-scheduler",
					TerminationGracePeriodSeconds: helpers.PointerInt64(30),
					Volumes:                       volumes,
				},
			},
		},
	}

	client.EXPECT().ApplyJob(jobDef).Times(1)
	client.EXPECT().GetJob(namespace, getJobName(reason)).Times(1).Return(nil, macherrs.NewNotFound(schema.GroupResource{"batch", "jobs"}, jobNamePrefix+reason))

	getConfigurationHashes := func(k8sClient kubernetes.ClientInt, queried map[string]interface{}) map[string]string {
		return map[string]string{"testHash": "test"}
	}
	migrationDone := func(k8sClient kubernetes.ClientInt) error {
		return nil
	}
	configurationDone := func(k8sClient kubernetes.ClientInt) error {
		return nil
	}

	query, _, err := AdaptFunc(
		monitor,
		componentLabels,
		namespace,
		reason,
		nodeselector,
		tolerations,
		resources,
		&version,
		cmName,
		certPath,
		secretName,
		secretPath,
		consoleCMName,
		secretVarsName,
		secretPasswordsName,
		users,
		migrationDone,
		configurationDone,
		getConfigurationHashes,
	)

	queried := map[string]interface{}{}
	database.SetDatabaseInQueried(queried, &database.Current{
		Host: dbHost,
		Port: dbPort,
	})

	assert.NoError(t, err)
	ensure, err := query(client, queried)
	assert.NoError(t, err)
	assert.NoError(t, ensure(client))
}
