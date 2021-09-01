package managed

import (
	"github.com/caos/orbos/mntr"
	kubernetesmock "github.com/caos/orbos/pkg/kubernetes/mock"
	"github.com/caos/orbos/pkg/labels"
	"github.com/caos/orbos/pkg/secret"
	"github.com/caos/orbos/pkg/tree"
	"github.com/caos/zitadel/operator/zitadel/kinds/backups/bucket"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func getTreeWithDBAndBackup(t *testing.T, masterkey string, saJson string, backupName string) *tree.Tree {

	bucketDesired := getDesiredTree(t, masterkey, &bucket.DesiredV0{
		Common: tree.NewCommon("databases.caos.ch/BucketBackup", "v0", false),
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
		Common: tree.NewCommon("databases.caos.ch/CockroachDB", "v0", false),
		Spec: Spec{
			Verbose:         false,
			ReplicaCount:    1,
			StorageCapacity: "368Gi",
			StorageClass:    "testSC",
			NodeSelector:    map[string]string{},
			ClusterDns:      "testDns",
		},
	})
}

func TestManaged_AdaptClean(t *testing.T) {
	monitor := mntr.Monitor{}
	componentLabels := labels.MustForComponent(labels.MustForAPI(labels.MustForOperator("testProd", "testOp", "testVersion"), "testKind", "v0"), "database")
	namespace := "testNs"
	masterkey := "testMk"
	k8sClient := kubernetesmock.NewMockClientInt(gomock.NewController(t))
	saJson := "testSA"
	backupName := "testBucket"

	features := []string{Clean}
	SetClean(k8sClient, namespace, 1)
	//k8sClient.EXPECT().WaitUntilStatefulsetIsReady(namespace, SfsName, true, true, 60*time.Second).Times(1)

	desired := getTreeWithDBAndBackup(t, masterkey, saJson, backupName)

	query, _, _, _, _, _, err := Adapter(componentLabels, namespace, features, "")(monitor, desired, &tree.Tree{})
	assert.NoError(t, err)

	databases := []string{"test1", "test2"}
	users := []string{"test1", "test2"}
	queried := bucket.SetQueriedForDatabases(databases, users)
	ensure, err := query(k8sClient, queried)
	assert.NoError(t, err)
	assert.NotNil(t, ensure)

	assert.NoError(t, ensure(k8sClient))
}
