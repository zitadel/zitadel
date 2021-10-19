package backups

import (
	"fmt"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/labels"
	"github.com/caos/orbos/pkg/secret"
	"github.com/caos/orbos/pkg/tree"
	"github.com/caos/zitadel/operator"
	"github.com/caos/zitadel/operator/database/kinds/backups/bucket"
	"github.com/caos/zitadel/operator/database/kinds/backups/s3"
	corev1 "k8s.io/api/core/v1"
)

func Adapt(
	monitor mntr.Monitor,
	desiredTree *tree.Tree,
	currentTree *tree.Tree,
	name string,
	namespace string,
	componentLabels *labels.Component,
	checkDBReady operator.EnsureFunc,
	timestamp string,
	nodeselector map[string]string,
	tolerations []corev1.Toleration,
	version string,
	dbURL string,
	dbPort int32,
	features []string,
	customImageRegistry string,
) (
	operator.QueryFunc,
	operator.DestroyFunc,
	operator.ConfigureFunc,
	map[string]*secret.Secret,
	map[string]*secret.Existing,
	bool,
	error,
) {
	switch desiredTree.Common.Kind {
	case "databases.caos.ch/BucketBackup":
		return bucket.AdaptFunc(
			name,
			namespace,
			labels.MustForComponent(
				labels.MustReplaceAPI(
					labels.GetAPIFromComponent(componentLabels),
					"BucketBackup",
					desiredTree.Common.Version(),
				),
				"backup"),
			checkDBReady,
			timestamp,
			nodeselector,
			tolerations,
			version,
			dbURL,
			dbPort,
			features,
			customImageRegistry,
		)(monitor, desiredTree, currentTree)
	case "databases.caos.ch/S3Backup":
		return s3.AdaptFunc(
			name,
			namespace,
			labels.MustForComponent(
				labels.MustReplaceAPI(
					labels.GetAPIFromComponent(componentLabels),
					"S3Backup",
					desiredTree.Common.Version(),
				),
				"backup"),
			checkDBReady,
			timestamp,
			nodeselector,
			tolerations,
			version,
			dbURL,
			dbPort,
			features,
			customImageRegistry,
		)(monitor, desiredTree, currentTree)
	default:
		return nil, nil, nil, nil, nil, false, mntr.ToUserError(fmt.Errorf("unknown database kind %s", desiredTree.Common.Kind))
	}
}

func GetBackupList(
	monitor mntr.Monitor,
	k8sClient kubernetes.ClientInt,
	name string,
	desiredTree *tree.Tree,
) (
	[]string,
	error,
) {
	switch desiredTree.Common.Kind {
	case "databases.caos.ch/BucketBackup":
		return bucket.BackupList()(monitor, k8sClient, name, desiredTree)
	case "databases.caos.ch/S3Backup":
		return s3.BackupList()(monitor, k8sClient, name, desiredTree)
	default:
		return nil, mntr.ToUserError(fmt.Errorf("unknown database kind %s", desiredTree.Common.Kind))
	}
}
