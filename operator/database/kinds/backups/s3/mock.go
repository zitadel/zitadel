package s3

import (
	kubernetesmock "github.com/caos/orbos/pkg/kubernetes/mock"
	coreBackup "github.com/caos/zitadel/operator/database/kinds/backups/core"
	"github.com/caos/zitadel/operator/database/kinds/databases/core"
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
	labelsSecret map[string]string,
	akid, sak, st string,
) {
	k8sClient.EXPECT().ApplySecret(&corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      coreBackup.GetSecretName(backupName),
			Namespace: namespace,
			Labels:    labelsSecret,
		},
		StringData: map[string]string{
			accessKeyIDKey:     akid,
			secretAccessKeyKey: sak,
			sessionTokenKey:    st,
		},
		Type: "Opaque",
	}).MinTimes(1).MaxTimes(1).Return(nil)

	k8sClient.EXPECT().ApplyJob(gomock.Any()).Times(1).Return(nil)
	k8sClient.EXPECT().GetJob(namespace, coreBackup.GetBackupJobName(backupName)).Times(1).Return(nil, macherrs.NewNotFound(schema.GroupResource{"batch", "jobs"}, coreBackup.GetBackupJobName(backupName)))
	k8sClient.EXPECT().WaitUntilJobCompleted(namespace, coreBackup.GetBackupJobName(backupName), gomock.Any()).Times(1).Return(nil)
	k8sClient.EXPECT().DeleteJob(namespace, coreBackup.GetBackupJobName(backupName)).Times(1).Return(nil)
}

func SetBackup(
	k8sClient *kubernetesmock.MockClientInt,
	namespace string,
	backupName string,
	labelsSecret map[string]string,
	akid, sak, st string,
) {
	k8sClient.EXPECT().ApplySecret(&corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      coreBackup.GetSecretName(backupName),
			Namespace: namespace,
			Labels:    labelsSecret,
		},
		StringData: map[string]string{
			accessKeyIDKey:     akid,
			secretAccessKeyKey: sak,
			sessionTokenKey:    st,
		},
		Type: "Opaque",
	}).MinTimes(1).MaxTimes(1).Return(nil)

	k8sClient.EXPECT().ApplyCronJob(gomock.Any()).Times(1).Return(nil)
}

func SetRestore(
	k8sClient *kubernetesmock.MockClientInt,
	namespace string,
	backupName string,
	labelsSecret map[string]string,
	akid, sak, st string,
) {
	k8sClient.EXPECT().ApplySecret(&corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      coreBackup.GetSecretName(backupName),
			Namespace: namespace,
			Labels:    labelsSecret,
		},
		StringData: map[string]string{
			accessKeyIDKey:     akid,
			secretAccessKeyKey: sak,
			sessionTokenKey:    st,
		},
		Type: "Opaque",
	}).MinTimes(1).MaxTimes(1).Return(nil)

	k8sClient.EXPECT().ApplyJob(gomock.Any()).Times(1).Return(nil)
	k8sClient.EXPECT().GetJob(namespace, coreBackup.GetRestoreJobName(backupName)).Times(1).Return(nil, macherrs.NewNotFound(schema.GroupResource{"batch", "jobs"}, coreBackup.GetRestoreJobName(backupName)))
	k8sClient.EXPECT().WaitUntilJobCompleted(namespace, coreBackup.GetRestoreJobName(backupName), gomock.Any()).Times(1).Return(nil)
	k8sClient.EXPECT().DeleteJob(namespace, coreBackup.GetRestoreJobName(backupName)).Times(1).Return(nil)
}
