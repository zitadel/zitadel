package backup

import (
	"github.com/caos/orbos/pkg/labels"
	"github.com/caos/zitadel/operator/helpers"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func getCronJob(
	namespace string,
	nameLabels *labels.Name,
	cron string,
	jobSpecDef batchv1.JobSpec,
) *v1beta1.CronJob {
	return &v1beta1.CronJob{
		ObjectMeta: v1.ObjectMeta{
			Name:      nameLabels.Name(),
			Namespace: namespace,
			Labels:    labels.MustK8sMap(nameLabels),
		},
		Spec: v1beta1.CronJobSpec{
			Schedule:          cron,
			ConcurrencyPolicy: v1beta1.ForbidConcurrent,
			JobTemplate: v1beta1.JobTemplateSpec{
				Spec: jobSpecDef,
			},
		},
	}
}

func getJob(
	namespace string,
	nameLabels *labels.Name,
	jobSpecDef batchv1.JobSpec,
) *batchv1.Job {
	return &batchv1.Job{
		ObjectMeta: v1.ObjectMeta{
			Name:      nameLabels.Name(),
			Namespace: namespace,
			Labels:    labels.MustK8sMap(nameLabels),
		},
		Spec: jobSpecDef,
	}
}

func getJobSpecDef(
	nodeselector map[string]string,
	tolerations []corev1.Toleration,
	saSecretName string,
	saSecretKey string,
	configSecretName string,
	configSecretKey string,
	backupName string,
	command string,
	image string,
) batchv1.JobSpec {
	return batchv1.JobSpec{
		Template: corev1.PodTemplateSpec{
			Spec: corev1.PodSpec{
				RestartPolicy: corev1.RestartPolicyNever,
				NodeSelector:  nodeselector,
				Tolerations:   tolerations,
				Containers: []corev1.Container{
					{},

					{
						Name:  backupName,
						Image: image,
						Command: []string{
							"/bin/bash",
							"-c",
							command,
						},
						VolumeMounts: []corev1.VolumeMount{{
							Name:      saInternalSecretName,
							SubPath:   saSecretKey,
							MountPath: saSecretPath,
						}, {
							Name:      configInternalSecretName,
							SubPath:   configSecretKey,
							MountPath: configSecretPath,
						}},
						ImagePullPolicy: corev1.PullAlways,
					}},
				Volumes: []corev1.Volume{{
					Name: saInternalSecretName,
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName:  saSecretName,
							DefaultMode: helpers.PointerInt32(defaultMode),
						},
					},
				}, {
					Name: configInternalSecretName,
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName:  configSecretName,
							DefaultMode: helpers.PointerInt32(defaultMode),
						},
					},
				}},
			},
		},
	}
}
