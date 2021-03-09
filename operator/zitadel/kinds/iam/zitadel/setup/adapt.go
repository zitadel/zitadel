package setup

import (
	"time"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/kubernetes/k8s"
	"github.com/caos/orbos/pkg/kubernetes/resources/job"
	"github.com/caos/orbos/pkg/labels"
	"github.com/caos/zitadel/operator"
	"github.com/caos/zitadel/operator/helpers"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/deployment"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	jobNamePrefix = "zitadel-setup-"
	containerName = "zitadel"
	rootSecret    = "client-root"
	dbSecrets     = "db-secrets"
	timeout       = 300 * time.Second
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
	users []string,
	migrationDone operator.EnsureFunc,
	configurationDone operator.EnsureFunc,
	getConfigurationHashes func(k8sClient kubernetes.ClientInt, queried map[string]interface{}) map[string]string,
) (
	operator.QueryFunc,
	operator.DestroyFunc,
	error,
) {
	internalMonitor := monitor.WithField("component", "setup")

	jobName := getJobName(reason)
	nameLabels := labels.MustForName(componentLabels, jobName)

	destroyJ, err := job.AdaptFuncToDestroy(jobName, namespace)
	if err != nil {
		return nil, nil, err
	}

	destroyers := []operator.DestroyFunc{
		operator.ResourceDestroyToZitadelDestroy(destroyJ),
	}

	return func(k8sClient kubernetes.ClientInt, queried map[string]interface{}) (operator.EnsureFunc, error) {

			jobDef := jobDef(
				nameLabels,
				users,
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
				componentLabels,
				nodeselector,
				tolerations,
			)

			hashes := getConfigurationHashes(k8sClient, queried)
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
				operator.EnsureFuncToQueryFunc(migrationDone),
				operator.EnsureFuncToQueryFunc(configurationDone),
				operator.ResourceQueryToZitadelQuery(query),
			}

			return operator.QueriersToEnsureFunc(internalMonitor, false, queriers, k8sClient, queried)
		},
		operator.DestroyersToDestroyFunc(internalMonitor, destroyers),
		nil
}

func jobDef(name *labels.Name, users []string, version *string, resources *k8s.Resources, cmName string, certPath string, secretName string, secretPath string, consoleCMName string, secretVarsName string, secretPasswordsName string, namespace string, componentLabels *labels.Component, nodeselector map[string]string, tolerations []corev1.Toleration) *batchv1.Job {
	initContainers := []corev1.Container{
		deployment.GetInitContainer(
			rootSecret,
			dbSecrets,
			users,
			deployment.RunAsUser,
		)}

	containers := []corev1.Container{
		deployment.GetContainer(
			containerName,
			*version,
			deployment.RunAsUser,
			true,
			deployment.GetResourcesFromDefault(resources),
			cmName,
			certPath,
			secretName,
			secretPath,
			consoleCMName,
			secretVarsName,
			secretPasswordsName,
			users,
			dbSecrets,
			"setup",
		)}

	volumes := deployment.GetVolumes(secretName, secretPasswordsName, consoleCMName, users)

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
					NodeSelector:    nodeselector,
					Tolerations:     tolerations,
					InitContainers:  initContainers,
					Containers:      containers,
					SecurityContext: &corev1.PodSecurityContext{},

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
