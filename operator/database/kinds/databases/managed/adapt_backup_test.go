package managed

import (
	"testing"
	"time"

	"github.com/caos/orbos/mntr"
	kubernetesmock "github.com/caos/orbos/pkg/kubernetes/mock"
	"github.com/caos/orbos/pkg/labels"
	"github.com/caos/orbos/pkg/secret"
	"github.com/caos/orbos/pkg/tree"
	"github.com/caos/zitadel/operator/database/kinds/backups/bucket"
	"github.com/caos/zitadel/operator/database/kinds/backups/bucket/backup"
	"github.com/caos/zitadel/operator/database/kinds/backups/bucket/clean"
	"github.com/caos/zitadel/operator/database/kinds/backups/bucket/restore"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
)

func getTreeWithDBAndBackup(t *testing.T, masterkey string, saJson string, backupName string) *tree.Tree {

	bucketDesired := getDesiredTree(t, masterkey, &bucket.DesiredV0{
		Common: &tree.Common{
			Kind:    "databases.caos.ch/BucketBackup",
			Version: "v0",
		},
		Spec: &bucket.Spec{
			Verbose: true,
			Cron:    "testCron",
			Bucket:  "testBucket",
			ServiceAccountJSON: &secret.Secret{
				Value: saJson,
			},
		},
	})
	bucketDesiredKind, err := bucket.ParseDesiredV0(bucketDesired)
	assert.NoError(t, err)
	bucketDesired.Parsed = bucketDesiredKind

	return getDesiredTree(t, masterkey, &DesiredV0{
		Common: &tree.Common{
			Kind:    "databases.caos.ch/CockroachDB",
			Version: "v0",
		},
		Spec: Spec{
			Verbose:         false,
			ReplicaCount:    1,
			StorageCapacity: "368Gi",
			StorageClass:    "testSC",
			NodeSelector:    map[string]string{},
			ClusterDns:      "testDns",
			Backups:         map[string]*tree.Tree{backupName: bucketDesired},
		},
	})
}

func TestManaged_AdaptBucketBackup(t *testing.T) {
	monitor := mntr.Monitor{}
	componentLabels := labels.MustForComponent(labels.MustForAPI(labels.MustForOperator("testProd", "testOp", "testVersion"), "testKind", "v0"), "database")

	labels := map[string]string{
		"app.kubernetes.io/component":  "backup",
		"app.kubernetes.io/managed-by": "testOp",
		"app.kubernetes.io/name":       "backup-serviceaccountjson",
		"app.kubernetes.io/part-of":    "testProd",
		"app.kubernetes.io/version":    "testVersion",
		"caos.ch/apiversion":           "v0",
		"caos.ch/kind":                 "BucketBackup",
	}
	namespace := "testNs"
	timestamp := "testTs"
	nodeselector := map[string]string{"test": "test"}
	tolerations := []corev1.Toleration{}
	k8sClient := kubernetesmock.NewMockClientInt(gomock.NewController(t))
	backupName := "testBucket"
	saJson := "testSA"
	masterkey := "testMk"

	desired := getTreeWithDBAndBackup(t, masterkey, saJson, backupName)

	features := []string{backup.Normal}
	bucket.SetBackup(k8sClient, namespace, labels, saJson)
	k8sClient.EXPECT().WaitUntilStatefulsetIsReady(namespace, SfsName, true, true, 60*time.Second)

	query, _, _, _, _, _, err := Adapter(componentLabels, namespace, timestamp, nodeselector, tolerations, features)(monitor, desired, &tree.Tree{})
	assert.NoError(t, err)

	databases := []string{"test1", "test2"}
	queried := bucket.SetQueriedForDatabases(databases, []string{})
	ensure, err := query(k8sClient, queried)
	assert.NoError(t, err)
	assert.NotNil(t, ensure)

	assert.NoError(t, ensure(k8sClient))
}

func TestManaged_AdaptBucketInstantBackup(t *testing.T) {
	monitor := mntr.Monitor{}
	componentLabels := labels.MustForComponent(labels.MustForAPI(labels.MustForOperator("testProd", "testOp", "testVersion"), "testKind", "v0"), "database")
	labels := map[string]string{
		"app.kubernetes.io/component":  "backup",
		"app.kubernetes.io/managed-by": "testOp",
		"app.kubernetes.io/name":       "backup-serviceaccountjson",
		"app.kubernetes.io/part-of":    "testProd",
		"app.kubernetes.io/version":    "testVersion",
		"caos.ch/apiversion":           "v0",
		"caos.ch/kind":                 "BucketBackup",
	}
	namespace := "testNs"
	timestamp := "testTs"
	nodeselector := map[string]string{"test": "test"}
	tolerations := []corev1.Toleration{}
	masterkey := "testMk"
	k8sClient := kubernetesmock.NewMockClientInt(gomock.NewController(t))
	saJson := "testSA"
	backupName := "testBucket"

	features := []string{backup.Instant}
	bucket.SetInstantBackup(k8sClient, namespace, backupName, labels, saJson)
	k8sClient.EXPECT().WaitUntilStatefulsetIsReady(namespace, SfsName, true, true, 60*time.Second)

	desired := getTreeWithDBAndBackup(t, masterkey, saJson, backupName)

	query, _, _, _, _, _, err := Adapter(componentLabels, namespace, timestamp, nodeselector, tolerations, features)(monitor, desired, &tree.Tree{})
	assert.NoError(t, err)

	databases := []string{"test1", "test2"}
	queried := bucket.SetQueriedForDatabases(databases, []string{})
	ensure, err := query(k8sClient, queried)
	assert.NoError(t, err)
	assert.NotNil(t, ensure)

	assert.NoError(t, ensure(k8sClient))
}

func TestManaged_AdaptBucketCleanAndRestore(t *testing.T) {
	monitor := mntr.Monitor{}
	componentLabels := labels.MustForComponent(labels.MustForAPI(labels.MustForOperator("testProd", "testOp", "testVersion"), "testKind", "v0"), "database")
	labels := map[string]string{
		"app.kubernetes.io/component":  "backup",
		"app.kubernetes.io/managed-by": "testOp",
		"app.kubernetes.io/name":       "backup-serviceaccountjson",
		"app.kubernetes.io/part-of":    "testProd",
		"app.kubernetes.io/version":    "testVersion",
		"caos.ch/apiversion":           "v0",
		"caos.ch/kind":                 "BucketBackup",
	}
	namespace := "testNs"
	timestamp := "testTs"
	nodeselector := map[string]string{"test": "test"}
	tolerations := []corev1.Toleration{}
	masterkey := "testMk"
	k8sClient := kubernetesmock.NewMockClientInt(gomock.NewController(t))
	saJson := "testSA"
	backupName := "testBucket"

	features := []string{restore.Instant, clean.Instant}
	bucket.SetRestore(k8sClient, namespace, backupName, labels, saJson)
	SetClean(k8sClient, namespace, 1)
	k8sClient.EXPECT().WaitUntilStatefulsetIsReady(namespace, SfsName, true, true, 60*time.Second).Times(1)

	desired := getTreeWithDBAndBackup(t, masterkey, saJson, backupName)

	query, _, _, _, _, _, err := Adapter(componentLabels, namespace, timestamp, nodeselector, tolerations, features)(monitor, desired, &tree.Tree{})
	assert.NoError(t, err)

	databases := []string{"test1", "test2"}
	users := []string{"test1", "test2"}
	queried := bucket.SetQueriedForDatabases(databases, users)
	ensure, err := query(k8sClient, queried)
	assert.NoError(t, err)
	assert.NotNil(t, ensure)

	assert.NoError(t, ensure(k8sClient))
}
