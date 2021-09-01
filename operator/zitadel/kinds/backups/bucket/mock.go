package bucket

import (
	kubernetesmock "github.com/caos/orbos/pkg/kubernetes/mock"
	"github.com/caos/zitadel/operator/database/kinds/databases/core"
	"github.com/caos/zitadel/operator/zitadel/kinds/backups/bucket/backup"
	"github.com/caos/zitadel/operator/zitadel/kinds/backups/bucket/restore"
	"github.com/golang/mock/gomock"
	corev1 "k8s.io/api/core/v1"
	macherrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func SetQueriedForDatabases(databases, users []string) map[string]interface{} {
	queried := map[string]interface{}{}
	core.SetQueriedForDatabaseDBList(queried, databases, users)

	return queried
}

func SetInstantBackup(
	k8sClient *kubernetesmock.MockClientInt,
	namespace string,
	backupName string,
	labels map[string]string,
	saJson string,
	akid string,
	sak string,
) {
	k8sClient.EXPECT().ApplySecret(&corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      backupSecretName,
			Namespace: namespace,
			Labels:    labels,
		},
		StringData: map[string]string{
			saSecretKey:  saJson,
			assetAKIDKey: akid,
			assetSAKKey:  sak,
		},
		Type: "Opaque",
	}).MinTimes(1).MaxTimes(1).Return(nil)

	k8sClient.EXPECT().ApplyJob(gomock.Any()).Times(1).Return(nil)
	k8sClient.EXPECT().GetJob(namespace, backup.GetJobName(backupName)).Times(1).Return(nil, macherrs.NewNotFound(schema.GroupResource{"batch", "jobs"}, backup.GetJobName(backupName)))
	k8sClient.EXPECT().WaitUntilJobCompleted(namespace, backup.GetJobName(backupName), gomock.Any()).Times(1).Return(nil)
	k8sClient.EXPECT().DeleteJob(namespace, backup.GetJobName(backupName)).Times(1).Return(nil)
}

func SetBackup(
	k8sClient *kubernetesmock.MockClientInt,
	namespace string,
	labels map[string]string,
	saJson string,
	akid string,
	sak string,
) {
	k8sClient.EXPECT().ApplySecret(&corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      backupSecretName,
			Namespace: namespace,
			Labels:    labels,
		},
		StringData: map[string]string{
			saSecretKey:  saJson,
			assetAKIDKey: akid,
			assetSAKKey:  sak,
		},
		Type: "Opaque",
	}).MinTimes(1).MaxTimes(1).Return(nil)
	k8sClient.EXPECT().ApplyCronJob(gomock.Any()).Times(1).Return(nil)
}

func SetRestore(
	k8sClient *kubernetesmock.MockClientInt,
	namespace string,
	backupName string,
	labels map[string]string,
	saJson string,
	akid string,
	sak string,
) {
	k8sClient.EXPECT().ApplySecret(&corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      backupSecretName,
			Namespace: namespace,
			Labels:    labels,
		},
		StringData: map[string]string{
			saSecretKey:  saJson,
			assetAKIDKey: akid,
			assetSAKKey:  sak,
		},
		Type: "Opaque",
	}).MinTimes(1).MaxTimes(1).Return(nil)

	k8sClient.EXPECT().ApplyJob(gomock.Any()).Times(1).Return(nil)
	k8sClient.EXPECT().GetJob(namespace, restore.GetJobName(backupName)).Times(1).Return(nil, macherrs.NewNotFound(schema.GroupResource{"batch", "jobs"}, restore.GetJobName(backupName)))
	k8sClient.EXPECT().WaitUntilJobCompleted(namespace, restore.GetJobName(backupName), gomock.Any()).Times(1).Return(nil)
	k8sClient.EXPECT().DeleteJob(namespace, restore.GetJobName(backupName)).Times(1).Return(nil)
}
