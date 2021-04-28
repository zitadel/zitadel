package managed

import (
	kubernetesmock "github.com/caos/orbos/pkg/kubernetes/mock"
	"github.com/golang/mock/gomock"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

func SetClean(
	k8sClient *kubernetesmock.MockClientInt,
	namespace string,
	replicas int,
) {
	k8sClient.EXPECT().ScaleStatefulset(namespace, gomock.Any(), 0).Return(nil)
	k8sClient.EXPECT().ListPersistentVolumeClaims(namespace).Return(&core.PersistentVolumeClaimList{
		Items: []core.PersistentVolumeClaim{
			{ObjectMeta: metav1.ObjectMeta{
				Name: "datadir-cockroachdb-0",
			}},
		},
	}, nil)
	k8sClient.EXPECT().ScaleStatefulset(namespace, gomock.Any(), 1).Return(nil)
	k8sClient.EXPECT().DeletePersistentVolumeClaim(namespace, gomock.Any(), gomock.Any()).Times(replicas).Return(nil)
	k8sClient.EXPECT().WaitUntilStatefulsetIsReady(namespace, gomock.Any(), true, false, gomock.Any())
	k8sClient.EXPECT().WaitUntilStatefulsetIsReady(namespace, gomock.Any(), true, true, time.Second*1)
	k8sClient.EXPECT().WaitUntilStatefulsetIsReady(namespace, gomock.Any(), true, true, gomock.Any())

	/*
		k8sClient.EXPECT().ApplyJob(gomock.Any()).Times(1).Return(nil)
		k8sClient.EXPECT().GetJob(namespace, clean.GetJobName(backupName)).Times(1).Return(nil, macherrs.NewNotFound(schema.GroupResource{"batch", "jobs"}, clean.GetJobName(backupName)))
		k8sClient.EXPECT().WaitUntilJobCompleted(namespace, clean.GetJobName(backupName), gomock.Any()).Times(1).Return(nil)
		k8sClient.EXPECT().DeleteJob(namespace, clean.GetJobName(backupName)).Times(1).Return(nil)*/
}
