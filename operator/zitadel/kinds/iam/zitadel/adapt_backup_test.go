package zitadel

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	kubernetesmock "github.com/caos/orbos/pkg/kubernetes/mock"
	"github.com/caos/orbos/pkg/labels"
	"github.com/caos/orbos/pkg/secret"
	"github.com/caos/orbos/pkg/tree"
	"github.com/caos/zitadel/operator/zitadel/kinds/backups/bucket"
	"github.com/caos/zitadel/operator/zitadel/kinds/backups/bucket/backup"
	"github.com/caos/zitadel/operator/zitadel/kinds/backups/bucket/restore"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/configuration"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/database"
	databasemock "github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/database/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	corev1 "k8s.io/api/core/v1"
	"testing"
)

func getDesiredTree(t *testing.T, masterkey string, desired interface{}) *tree.Tree {
	secret.Masterkey = masterkey

	desiredTree := &tree.Tree{}
	data, err := yaml.Marshal(desired)
	assert.NoError(t, err)
	assert.NoError(t, yaml.Unmarshal(data, desiredTree))

	return desiredTree
}

func getTreeWithDBAndBackup(
	t *testing.T,
	masterkey string,
	saJson string,
	akid string,
	sak string,
	backupName string,
) *tree.Tree {

	bucketDesired := getDesiredTree(t, masterkey, &bucket.DesiredV0{
		Common: tree.NewCommon("zitadel.caos.ch/BucketBackup", "v0", false),
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
		Common: tree.NewCommon("zitadel.caos.ch/ZITADEL", "v0", false),
		Spec: &Spec{
			Verbose:      false,
			ReplicaCount: 1,
			NodeSelector: map[string]string{},
			Configuration: &configuration.Configuration{
				DNS: &configuration.DNS{
					Domain:        "test",
					TlsSecret:     "test-secret",
					ACMEAuthority: "none",
					Subdomains:    nil,
				},
				AssetStorage: &configuration.AssetStorage{
					Type:     "",
					Endpoint: "",
					AccessKeyID: &secret.Secret{
						Value: akid,
					},
					SecretAccessKey: &secret.Secret{
						Value: sak,
					},
					BucketPrefix: "",
					MultiDelete:  false,
				},
			},
			Backups: map[string]*tree.Tree{backupName: bucketDesired},
		},
	})
}

func getDbClient(
	t *testing.T,
	monitor mntr.Monitor,
	k8sClient kubernetes.ClientInt,
) database.Client {
	dbClient := databasemock.NewMockClient(gomock.NewController(t))
	host := "host"
	port := "port"
	httpPort := "80"

	dbClient.EXPECT().GetConnectionInfo(monitor, k8sClient).Return(host, port, httpPort, nil)

	return dbClient
}

func TestBackup_AdaptBucketBackup(t *testing.T) {
	monitor := mntr.Monitor{}
	apiLabels := labels.MustForAPI(labels.MustForOperator("testProd", "testOp", "testVersion"), "testKind", "v0")

	labels := map[string]string{
		"app.kubernetes.io/component":  "backup",
		"app.kubernetes.io/managed-by": "testOp",
		"app.kubernetes.io/name":       "backup-accounts",
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
	akid := "testAkid"
	sak := "testSak"
	masterkey := "testMk"
	version := "test"
	action := "testAction"

	desired := getTreeWithDBAndBackup(t, masterkey, saJson, akid, sak, backupName)

	features := []string{backup.Normal}
	bucket.SetBackup(k8sClient, namespace, labels, saJson, akid, sak)
	//k8sClient.EXPECT().WaitUntilStatefulsetIsReady(namespace, managed.SfsName, true, true, 60*time.Second)

	query, _, _, _, _, _, err := AdaptFunc(apiLabels, nodeselector, tolerations, getDbClient(t, monitor, k8sClient), namespace, action, &version, features, "", timestamp)(monitor, desired, &tree.Tree{})
	assert.NoError(t, err)

	databases := []string{"test1", "test2"}
	queried := bucket.SetQueriedForDatabases(databases, []string{})
	ensure, err := query(k8sClient, queried)
	assert.NoError(t, err)
	assert.NotNil(t, ensure)

	assert.NoError(t, ensure(k8sClient))
}

func TestBackup_AdaptBucketInstantBackup(t *testing.T) {
	monitor := mntr.Monitor{}
	apiLabels := labels.MustForAPI(labels.MustForOperator("testProd", "testOp", "testVersion"), "testKind", "v0")
	labels := map[string]string{
		"app.kubernetes.io/component":  "backup",
		"app.kubernetes.io/managed-by": "testOp",
		"app.kubernetes.io/name":       "backup-accounts",
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
	akid := "testAkid"
	sak := "testSak"
	backupName := "testBucket"
	version := "test"
	action := "testAction"

	features := []string{backup.Instant}
	bucket.SetInstantBackup(k8sClient, namespace, backupName, labels, saJson, akid, sak)
	//k8sClient.EXPECT().WaitUntilStatefulsetIsReady(namespace, managed.SfsName, true, true, 60*time.Second)

	desired := getTreeWithDBAndBackup(t, masterkey, saJson, akid, sak, backupName)

	query, _, _, _, _, _, err := AdaptFunc(apiLabels, nodeselector, tolerations, getDbClient(t, monitor, k8sClient), namespace, action, &version, features, "", timestamp)(monitor, desired, &tree.Tree{})
	assert.NoError(t, err)

	databases := []string{"test1", "test2"}
	queried := bucket.SetQueriedForDatabases(databases, []string{})
	ensure, err := query(k8sClient, queried)
	assert.NoError(t, err)
	assert.NotNil(t, ensure)

	assert.NoError(t, ensure(k8sClient))
}

func TestBackup_AdaptBucketRestore(t *testing.T) {
	monitor := mntr.Monitor{}
	apiLabels := labels.MustForAPI(labels.MustForOperator("testProd", "testOp", "testVersion"), "testKind", "v0")
	labels := map[string]string{
		"app.kubernetes.io/component":  "backup",
		"app.kubernetes.io/managed-by": "testOp",
		"app.kubernetes.io/name":       "backup-accounts",
		"app.kubernetes.io/part-of":    "testProd",
		"app.kubernetes.io/version":    "testVersion",
		"caos.ch/apiversion":           "v0",
		"caos.ch/kind":                 "BucketBackup",
	}
	namespace := "testNs"
	timestamp := "testTs"
	nodeselector := map[string]string{"test": "test"}
	tolerations := []corev1.Toleration{}
	version := "testVersion"
	masterkey := "testMk"
	k8sClient := kubernetesmock.NewMockClientInt(gomock.NewController(t))
	saJson := "testSA"
	akid := "testAkid"
	sak := "testSak"
	backupName := "testBucket"
	action := "testAction"

	features := []string{restore.Instant}

	bucket.SetRestore(k8sClient, namespace, backupName, labels, saJson, akid, sak)
	//k8sClient.EXPECT().WaitUntilStatefulsetIsReady(namespace, managed.SfsName, true, true, 60*time.Second).Times(1)

	desired := getTreeWithDBAndBackup(t, masterkey, saJson, akid, sak, backupName)

	query, _, _, _, _, _, err := AdaptFunc(apiLabels, nodeselector, tolerations, getDbClient(t, monitor, k8sClient), namespace, action, &version, features, "", timestamp)(monitor, desired, &tree.Tree{})
	assert.NoError(t, err)

	databases := []string{"test1", "test2"}
	users := []string{"test1", "test2"}
	queried := bucket.SetQueriedForDatabases(databases, users)
	ensure, err := query(k8sClient, queried)
	assert.NoError(t, err)
	assert.NotNil(t, ensure)

	assert.NoError(t, ensure(k8sClient))
}
