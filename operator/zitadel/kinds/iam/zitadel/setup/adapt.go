package setup

import (
	"fmt"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/kubernetes/k8s"
	"github.com/caos/orbos/pkg/kubernetes/resources/job"
	"github.com/caos/orbos/pkg/labels"
	"github.com/caos/zitadel/operator"
	"github.com/caos/zitadel/operator/helpers"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/deployment"
	"github.com/caos/zitadel/pkg/databases/db"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	jobNamePrefix = "zitadel-setup-"
	containerName = "zitadel"
)

func AdaptFunc(
	monitor mntr.Monitor,
	componentLabels *labels.Component,
	namespace string,
	reason string,
	nodeselector map[string]string,
	tolerations []corev1.Toleration,
	resources *k8s.Resources,
	version *string,
	cmName string,
	certPath string,
	secretName string,
	secretPath string,
	consoleCMName string,
	secretVarsName string,
	secretPasswordsName string,
	customImageRegistry string,
	dbConn db.Connection,
) (
	func(
		getConfigurationHashes func(k8sClient kubernetes.ClientInt, queried map[string]interface{}) (map[string]string, error),
	) operator.QueryFunc,
	operator.DestroyFunc,
	error,
) {
	internalMonitor := monitor.WithField("component", "setup")

	jobName := getJobName(reason)
	nameLabels := labels.MustForName(componentLabels, jobName)

	destroyJ, err := job.AdaptFuncToDestroy(namespace, jobName)
	if err != nil {
		return nil, nil, err
	}

	destroyers := []operator.DestroyFunc{
		operator.ResourceDestroyToZitadelDestroy(destroyJ),
	}

	return func(
			getConfigurationHashes func(k8sClient kubernetes.ClientInt, queried map[string]interface{}) (map[string]string, error),
		) operator.QueryFunc {
			return func(k8sClient kubernetes.ClientInt, queried map[string]interface{}) (operator.EnsureFunc, error) {
				jobDef := jobDef(
					nameLabels,
					version,
					resources,
					cmName,
					certPath,
					secretName,
					secretPath,
					consoleCMName,
					secretVarsName,
					secretPasswordsName,
					namespace,
					nodeselector,
					tolerations,
					customImageRegistry,
					dbConn,
				)

				hashes, err := getConfigurationHashes(k8sClient, queried)
				if err != nil {
					return nil, err
				}
				if hashes != nil && len(hashes) != 0 {
					for k, v := range hashes {
						jobDef.Annotations[k] = v
						jobDef.Spec.Template.Annotations[k] = v
					}
				}

				query, err := job.AdaptFuncToEnsure(jobDef)
				if err != nil {
					return nil, err
				}

				queriers := []operator.QueryFunc{
					operator.ResourceQueryToZitadelQuery(query),
				}

				return operator.QueriersToEnsureFunc(internalMonitor, false, queriers, k8sClient, queried)
			}
		},
		operator.DestroyersToDestroyFunc(internalMonitor, destroyers),
		nil

}

func jobDef(
	name *labels.Name,
	version *string,
	resources *k8s.Resources,
	cmName string,
	certPath string,
	secretName string,
	secretPath string,
	consoleCMName string,
	secretVarsName string,
	secretPasswordsName string,
	namespace string,
	nodeselector map[string]string,
	tolerations []corev1.Toleration,
	customImageRegistry string,
	dbConn db.Connection,
) *batchv1.Job {

	volumes := deployment.GetVolumes(secretName, secretPasswordsName, consoleCMName)
	var initContainers []corev1.Container
	var chownedVolumeMount *corev1.VolumeMount

	ssl := dbConn.SSL()
	if ssl != nil {
		chownedVolumeMount = &corev1.VolumeMount{
			Name:      "chowned-certs",
			MountPath: certPath,
		}

		srcVolume, destVolume, chownCertsContainer := db.InitChownCerts(ssl, customImageRegistry, fmt.Sprintf("%d:%d", deployment.RunAsUser, deployment.RunAsUser), corev1.VolumeMount{
			Name:      "certs",
			MountPath: "/certificates",
		}, *chownedVolumeMount)

		volumes = append(volumes, srcVolume, destVolume)
		initContainers = []corev1.Container{chownCertsContainer}
	}

	containers := []corev1.Container{
		deployment.GetContainer(
			containerName,
			*version,
			deployment.RunAsUser,
			true,
			deployment.GetResourcesFromDefault(resources),
			cmName,
			secretName,
			secretPath,
			consoleCMName,
			secretVarsName,
			chownedVolumeMount,
			"setup",
			customImageRegistry,
			dbConn,
		)}

	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name.Name(),
			Namespace:   namespace,
			Labels:      labels.MustK8sMap(name),
			Annotations: map[string]string{},
		},
		Spec: batchv1.JobSpec{
			Completions: helpers.PointerInt32(1),
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{},
				},
				Spec: corev1.PodSpec{
					NodeSelector:                  nodeselector,
					Tolerations:                   tolerations,
					InitContainers:                initContainers,
					Containers:                    containers,
					RestartPolicy:                 "Never",
					DNSPolicy:                     "ClusterFirst",
					SchedulerName:                 "default-scheduler",
					TerminationGracePeriodSeconds: helpers.PointerInt64(30),
					Volumes:                       volumes,
				},
			},
		},
	}
}

func getJobName(reason string) string {
	return jobNamePrefix + reason
}
