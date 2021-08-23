package backup

import (
	"testing"

	"github.com/caos/orbos/mntr"
	kubernetesmock "github.com/caos/orbos/pkg/kubernetes/mock"
	"github.com/caos/orbos/pkg/labels"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	macherrs "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func TestBackup_AdaptInstantBackup1(t *testing.T) {
	client := kubernetesmock.NewMockClientInt(gomock.NewController(t))

	features := []string{Instant}
	monitor := mntr.Monitor{}
	namespace := "testNs"

	bucketName := "testBucket"
	cron := "testCron"
	timestamp := "test"
	nodeselector := map[string]string{"test": "test"}
	tolerations := []corev1.Toleration{
		{Key: "testKey", Operator: "testOp"}}
	backupName := "testName"
	image := "testImage"
	saSecretName := "testSaSecretName"
	saSecretKey := "testSaSecretKey"
	configSecretName := "testConfigSecretName"
	configSecretKey := "testConfigSecretKey"
	jobName := GetJobName(backupName)
	componentLabels := labels.MustForComponent(labels.MustForAPI(labels.MustForOperator("testProd2", "testOp2", "testVersion2"), "testKind2", "testVersion2"), "testComponent")
	nameLabels := labels.MustForName(componentLabels, jobName)

	jobDef := getJob(
		namespace,
		nameLabels,
		getJobSpecDef(
			nodeselector,
			tolerations,
			saSecretName,
			saSecretKey,
			configSecretName,
			configSecretKey,
			backupName,
			getBackupCommand(
				timestamp,
				bucketName,
				backupName,
			),
			image,
		),
	)

	client.EXPECT().ApplyJob(jobDef).Times(1).Return(nil)
	client.EXPECT().GetJob(jobDef.Namespace, jobDef.Name).Times(1).Return(nil, macherrs.NewNotFound(schema.GroupResource{"batch", "jobs"}, jobName))

	query, _, err := AdaptFunc(
		monitor,
		backupName,
		namespace,
		componentLabels,
		bucketName,
		cron,
		saSecretName,
		saSecretKey,
		configSecretName,
		configSecretKey,
		timestamp,
		nodeselector,
		tolerations,
		features,
		image,
	)

	assert.NoError(t, err)
	queried := map[string]interface{}{}
	ensure, err := query(client, queried)
	assert.NoError(t, err)
	assert.NoError(t, ensure(client))
}

func TestBackup_AdaptInstantBackup2(t *testing.T) {
	client := kubernetesmock.NewMockClientInt(gomock.NewController(t))

	features := []string{Instant}
	monitor := mntr.Monitor{}
	namespace := "testNs2"
	bucketName := "testBucket2"
	cron := "testCron2"
	timestamp := "test2"
	nodeselector := map[string]string{"test2": "test2"}
	tolerations := []corev1.Toleration{
		{Key: "testKey2", Operator: "testOp2"}}
	backupName := "testName2"
	image := "testImage2"
	saSecretName := "testSaSecretName2"
	saSecretKey := "testSaSecretKey2"
	configSecretName := "testConfigSecretName2"
	configSecretKey := "testConfigSecretKey2"
	jobName := GetJobName(backupName)
	componentLabels := labels.MustForComponent(labels.MustForAPI(labels.MustForOperator("testProd2", "testOp2", "testVersion2"), "testKind2", "testVersion2"), "testComponent")
	nameLabels := labels.MustForName(componentLabels, jobName)

	jobDef := getJob(
		namespace,
		nameLabels,
		getJobSpecDef(
			nodeselector,
			tolerations,
			saSecretName,
			saSecretKey,
			configSecretName,
			configSecretKey,
			backupName,
			getBackupCommand(
				timestamp,
				bucketName,
				backupName,
			),
			image,
		),
	)

	client.EXPECT().ApplyJob(jobDef).Times(1).Return(nil)
	client.EXPECT().GetJob(jobDef.Namespace, jobDef.Name).Times(1).Return(nil, macherrs.NewNotFound(schema.GroupResource{"batch", "jobs"}, jobName))

	query, _, err := AdaptFunc(
		monitor,
		backupName,
		namespace,
		componentLabels,
		bucketName,
		cron,
		saSecretName,
		saSecretKey,
		configSecretName,
		configSecretKey,
		timestamp,
		nodeselector,
		tolerations,
		features,
		image,
	)

	assert.NoError(t, err)
	queried := map[string]interface{}{}
	ensure, err := query(client, queried)
	assert.NoError(t, err)
	assert.NoError(t, ensure(client))
}

func TestBackup_AdaptBackup1(t *testing.T) {
	client := kubernetesmock.NewMockClientInt(gomock.NewController(t))

	features := []string{Normal}
	monitor := mntr.Monitor{}
	namespace := "testNs"
	bucketName := "testBucket"
	cron := "testCron"
	timestamp := "test"
	nodeselector := map[string]string{"test": "test"}
	tolerations := []corev1.Toleration{
		{Key: "testKey", Operator: "testOp"}}
	backupName := "testName"
	image := "testImage"
	saSecretName := "testSaSecretName"
	saSecretKey := "testSaSecretKey"
	configSecretName := "testConfigSecretName"
	configSecretKey := "testConfigSecretKey"
	jobName := GetJobName(backupName)
	componentLabels := labels.MustForComponent(labels.MustForAPI(labels.MustForOperator("testProd2", "testOp2", "testVersion2"), "testKind2", "testVersion2"), "testComponent")
	nameLabels := labels.MustForName(componentLabels, jobName)

	jobDef := getCronJob(
		namespace,
		nameLabels,
		cron,
		getJobSpecDef(
			nodeselector,
			tolerations,
			saSecretName,
			saSecretKey,
			configSecretName,
			configSecretKey,
			backupName,
			getBackupCommand(
				timestamp,
				bucketName,
				backupName,
			),
			image,
		),
	)

	client.EXPECT().ApplyCronJob(jobDef).Times(1).Return(nil)

	query, _, err := AdaptFunc(
		monitor,
		backupName,
		namespace,
		componentLabels,
		bucketName,
		cron,
		saSecretName,
		saSecretKey,
		configSecretName,
		configSecretKey,
		timestamp,
		nodeselector,
		tolerations,
		features,
		image,
	)

	assert.NoError(t, err)
	queried := map[string]interface{}{}
	ensure, err := query(client, queried)
	assert.NoError(t, err)
	assert.NoError(t, ensure(client))
}

func TestBackup_AdaptBackup2(t *testing.T) {
	client := kubernetesmock.NewMockClientInt(gomock.NewController(t))

	features := []string{Normal}
	monitor := mntr.Monitor{}
	namespace := "testNs2"
	bucketName := "testBucket2"
	cron := "testCron2"
	timestamp := "test2"
	nodeselector := map[string]string{"test2": "test2"}
	tolerations := []corev1.Toleration{
		{Key: "testKey2", Operator: "testOp2"}}
	backupName := "testName2"
	image := "testImage2"
	saSecretName := "testSaSecretName2"
	saSecretKey := "testSaSecretKey2"
	configSecretName := "testConfigSecretName2"
	configSecretKey := "testConfigSecretKey2"
	jobName := GetJobName(backupName)
	componentLabels := labels.MustForComponent(labels.MustForAPI(labels.MustForOperator("testProd2", "testOp2", "testVersion2"), "testKind2", "testVersion2"), "testComponent")
	nameLabels := labels.MustForName(componentLabels, jobName)

	jobDef := getCronJob(
		namespace,
		nameLabels,
		cron,
		getJobSpecDef(
			nodeselector,
			tolerations,
			saSecretName,
			saSecretKey,
			configSecretName,
			configSecretKey,
			backupName,
			getBackupCommand(
				timestamp,
				bucketName,
				backupName,
			),
			image,
		),
	)

	client.EXPECT().ApplyCronJob(jobDef).Times(1).Return(nil)

	query, _, err := AdaptFunc(
		monitor,
		backupName,
		namespace,
		componentLabels,
		bucketName,
		cron,
		saSecretName,
		saSecretKey,
		configSecretName,
		configSecretKey,
		timestamp,
		nodeselector,
		tolerations,
		features,
		image,
	)

	assert.NoError(t, err)
	queried := map[string]interface{}{}
	ensure, err := query(client, queried)
	assert.NoError(t, err)
	assert.NoError(t, ensure(client))
}
