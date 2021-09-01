package restore

import (
	"testing"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	kubernetesmock "github.com/caos/orbos/pkg/kubernetes/mock"
	"github.com/caos/orbos/pkg/labels"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	macherrs "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func TestBackup_Adapt1(t *testing.T) {
	client := kubernetesmock.NewMockClientInt(gomock.NewController(t))

	monitor := mntr.Monitor{}
	namespace := "testNs"
	nodeselector := map[string]string{"test": "test"}
	tolerations := []corev1.Toleration{
		{Key: "testKey", Operator: "testOp"}}
	timestamp := "testTs"
	backupName := "testName"
	bucketName := "testBucket"
	image := "testImage"
	secretName := "testSecretName"
	saSecretKey := "testKey"
	akidKey := "testAKID"
	sakKey := "testSAK"
	endpoint := "testEndpoint"
	prefix := "testPrefix"
	dbURL := "testDB"
	dbPort := int32(80)
	jobName := GetJobName(backupName)
	componentLabels := labels.MustForComponent(labels.MustForAPI(labels.MustForOperator("testProd", "testOp", "testVersion"), "testKind", "testVersion"), "testComponent")
	nameLabels := labels.MustForName(componentLabels, jobName)

	checkDBReady := func(k8sClient kubernetes.ClientInt) error {
		return nil
	}

	jobDef := getJob(
		namespace,
		nameLabels,
		nodeselector,
		tolerations,
		secretName,
		saSecretKey,
		akidKey,
		sakKey,
		getCommand(
			timestamp,
			bucketName,
			backupName,
			certPath,
			saSecretPath,
			dbURL,
			dbPort,
			endpoint,
			prefix,
		),
		image,
	)

	client.EXPECT().ApplyJob(jobDef).Times(1).Return(nil)
	client.EXPECT().GetJob(jobDef.Namespace, jobDef.Name).Times(1).Return(nil, macherrs.NewNotFound(schema.GroupResource{"batch", "jobs"}, jobName))

	query, _, err := AdaptFunc(
		monitor,
		backupName,
		namespace,
		componentLabels,
		bucketName,
		timestamp,
		nodeselector,
		tolerations,
		checkDBReady,
		secretName,
		saSecretKey,
		akidKey,
		sakKey,
		dbURL,
		dbPort,
		image,
		endpoint,
		prefix,
	)

	assert.NoError(t, err)
	queried := map[string]interface{}{}
	ensure, err := query(client, queried)
	assert.NoError(t, err)
	assert.NoError(t, ensure(client))
}

func TestBackup_Adapt2(t *testing.T) {
	client := kubernetesmock.NewMockClientInt(gomock.NewController(t))

	monitor := mntr.Monitor{}
	namespace := "testNs2"
	nodeselector := map[string]string{"test2": "test2"}
	tolerations := []corev1.Toleration{
		{Key: "testKey2", Operator: "testOp2"}}
	timestamp := "testTs"
	backupName := "testName2"
	bucketName := "testBucket2"
	image := "testImage2"
	secretName := "testSecretName2"
	saSecretKey := "testKey2"
	akidKey := "testAKID2"
	sakKey := "testSAK2"
	endpoint := "testEndpoint2"
	prefix := "testPrefix2"
	dbURL := "testDB"
	dbPort := int32(80)
	jobName := GetJobName(backupName)
	componentLabels := labels.MustForComponent(labels.MustForAPI(labels.MustForOperator("testProd2", "testOp2", "testVersion2"), "testKind2", "testVersion2"), "testComponent2")
	nameLabels := labels.MustForName(componentLabels, jobName)

	checkDBReady := func(k8sClient kubernetes.ClientInt) error {
		return nil
	}

	jobDef := getJob(
		namespace,
		nameLabels,
		nodeselector,
		tolerations,
		secretName,
		saSecretKey,
		akidKey,
		sakKey,
		getCommand(
			timestamp,
			bucketName,
			backupName,
			certPath,
			saSecretPath,
			dbURL,
			dbPort,
			endpoint,
			prefix,
		),
		image,
	)

	client.EXPECT().ApplyJob(jobDef).Times(1).Return(nil)
	client.EXPECT().GetJob(jobDef.Namespace, jobDef.Name).Times(1).Return(nil, macherrs.NewNotFound(schema.GroupResource{"batch", "jobs"}, jobName))

	query, _, err := AdaptFunc(
		monitor,
		backupName,
		namespace,
		componentLabels,
		bucketName,
		timestamp,
		nodeselector,
		tolerations,
		checkDBReady,
		secretName,
		saSecretKey,
		akidKey,
		sakKey,
		dbURL,
		dbPort,
		image,
		endpoint,
		prefix,
	)

	assert.NoError(t, err)
	queried := map[string]interface{}{}
	ensure, err := query(client, queried)
	assert.NoError(t, err)
	assert.NoError(t, ensure(client))
}
