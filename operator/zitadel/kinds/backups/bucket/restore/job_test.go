package restore

import (
	"testing"

	"github.com/caos/orbos/pkg/labels"
	"github.com/caos/zitadel/operator/helpers"
	"github.com/stretchr/testify/assert"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestBackup_Job1(t *testing.T) {
	nodeselector := map[string]string{"test": "test"}
	tolerations := []corev1.Toleration{
		{Key: "testKey", Operator: "testOp"}}
	command := "test"
	backupSecretName := "testSecret"
	saSecretKey := "testSaKey"
	assetAKIDKey := "testAkidKey"
	assetSAKKey := "testSakKey"
	jobName := "testJob"
	namespace := "testNs"
	image := "testImage"
	k8sLabels := map[string]string{
		"app.kubernetes.io/component":  "testComponent",
		"app.kubernetes.io/managed-by": "testOp",
		"app.kubernetes.io/name":       jobName,
		"app.kubernetes.io/part-of":    "testProd",
		"app.kubernetes.io/version":    "testOpVersion",
		"caos.ch/apiversion":           "testVersion",
		"caos.ch/kind":                 "testKind"}
	componentLabels := labels.MustForComponent(labels.MustForAPI(labels.MustForOperator("testProd", "testOp", "testOpVersion"), "testKind", "testVersion"), "testComponent")
	nameLabels := labels.MustForName(componentLabels, jobName)

	equals :=
		&batchv1.Job{
			ObjectMeta: metav1.ObjectMeta{
				Name:      jobName,
				Namespace: namespace,
				Labels:    k8sLabels,
			},
			Spec: batchv1.JobSpec{
				Template: corev1.PodTemplateSpec{
					Spec: corev1.PodSpec{
						RestartPolicy: corev1.RestartPolicyNever,
						NodeSelector:  nodeselector,
						Tolerations:   tolerations,
						Containers: []corev1.Container{{
							Name:  jobName,
							Image: image,
							Command: []string{
								"/bin/bash",
								"-c",
								command,
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      certsInternalSecretName,
									MountPath: certPath,
								}, {
									Name:      saInternalSecretName,
									SubPath:   saSecretKey,
									MountPath: saSecretPath,
								}, {
									Name:      akidInternalSecretName,
									SubPath:   assetAKIDKey,
									MountPath: akidSecretPath,
								}, {
									Name:      sakInternalSecretName,
									SubPath:   assetSAKKey,
									MountPath: sakSecretPath,
								},
							},
							ImagePullPolicy: corev1.PullAlways,
						}},
						Volumes: []corev1.Volume{
							{
								Name: certsInternalSecretName,
								VolumeSource: corev1.VolumeSource{
									Secret: &corev1.SecretVolumeSource{
										SecretName:  rootSecretName,
										DefaultMode: helpers.PointerInt32(defaultMode),
									},
								},
							}, {
								Name: saInternalSecretName,
								VolumeSource: corev1.VolumeSource{
									Secret: &corev1.SecretVolumeSource{
										SecretName:  backupSecretName,
										DefaultMode: helpers.PointerInt32(defaultMode),
									},
								},
							}, {
								Name: akidInternalSecretName,
								VolumeSource: corev1.VolumeSource{
									Secret: &corev1.SecretVolumeSource{
										SecretName:  backupSecretName,
										DefaultMode: helpers.PointerInt32(defaultMode),
									},
								},
							}, {
								Name: sakInternalSecretName,
								VolumeSource: corev1.VolumeSource{
									Secret: &corev1.SecretVolumeSource{
										SecretName:  backupSecretName,
										DefaultMode: helpers.PointerInt32(defaultMode),
									},
								},
							},
						},
					},
				},
			},
		}

	assert.Equal(t, equals, getJob(
		namespace,
		nameLabels,
		nodeselector,
		tolerations,
		backupSecretName,
		saSecretKey,
		assetAKIDKey,
		assetSAKKey,
		command,
		image,
	))
}

func TestBackup_Job2(t *testing.T) {
	nodeselector := map[string]string{"test2": "test2"}
	tolerations := []corev1.Toleration{
		{Key: "testKey2", Operator: "testOp2"}}
	command := "test2"
	backupSecretName := "testSecret2"
	saSecretKey := "testSaKey2"
	assetAKIDKey := "testAkidKey2"
	assetSAKKey := "testSakKey2"
	jobName := "testJob2"
	namespace := "testNs2"
	testImage := "testImage2"
	k8sLabels := map[string]string{
		"app.kubernetes.io/component":  "testComponent2",
		"app.kubernetes.io/managed-by": "testOp2",
		"app.kubernetes.io/name":       jobName,
		"app.kubernetes.io/part-of":    "testProd2",
		"app.kubernetes.io/version":    "testOpVersion2",
		"caos.ch/apiversion":           "testVersion2",
		"caos.ch/kind":                 "testKind2"}
	componentLabels := labels.MustForComponent(labels.MustForAPI(labels.MustForOperator("testProd2", "testOp2", "testOpVersion2"), "testKind2", "testVersion2"), "testComponent2")
	nameLabels := labels.MustForName(componentLabels, jobName)

	equals :=
		&batchv1.Job{
			ObjectMeta: metav1.ObjectMeta{
				Name:      jobName,
				Namespace: namespace,
				Labels:    k8sLabels,
			},
			Spec: batchv1.JobSpec{
				Template: corev1.PodTemplateSpec{
					Spec: corev1.PodSpec{
						RestartPolicy: corev1.RestartPolicyNever,
						NodeSelector:  nodeselector,
						Tolerations:   tolerations,
						Containers: []corev1.Container{{
							Name:  jobName,
							Image: testImage,
							Command: []string{
								"/bin/bash",
								"-c",
								command,
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      certsInternalSecretName,
									MountPath: certPath,
								}, {
									Name:      saInternalSecretName,
									SubPath:   saSecretKey,
									MountPath: saSecretPath,
								}, {
									Name:      akidInternalSecretName,
									SubPath:   assetAKIDKey,
									MountPath: akidSecretPath,
								}, {
									Name:      sakInternalSecretName,
									SubPath:   assetSAKKey,
									MountPath: sakSecretPath,
								},
							},
							ImagePullPolicy: corev1.PullAlways,
						}},
						Volumes: []corev1.Volume{
							{
								Name: certsInternalSecretName,
								VolumeSource: corev1.VolumeSource{
									Secret: &corev1.SecretVolumeSource{
										SecretName:  rootSecretName,
										DefaultMode: helpers.PointerInt32(defaultMode),
									},
								},
							}, {
								Name: saInternalSecretName,
								VolumeSource: corev1.VolumeSource{
									Secret: &corev1.SecretVolumeSource{
										SecretName:  backupSecretName,
										DefaultMode: helpers.PointerInt32(defaultMode),
									},
								},
							}, {
								Name: akidInternalSecretName,
								VolumeSource: corev1.VolumeSource{
									Secret: &corev1.SecretVolumeSource{
										SecretName:  backupSecretName,
										DefaultMode: helpers.PointerInt32(defaultMode),
									},
								},
							}, {
								Name: sakInternalSecretName,
								VolumeSource: corev1.VolumeSource{
									Secret: &corev1.SecretVolumeSource{
										SecretName:  backupSecretName,
										DefaultMode: helpers.PointerInt32(defaultMode),
									},
								},
							},
						},
					},
				},
			},
		}

	assert.Equal(t, equals, getJob(
		namespace,
		nameLabels,
		nodeselector,
		tolerations,
		backupSecretName,
		saSecretKey,
		assetAKIDKey,
		assetSAKKey,
		command,
		testImage,
	))
}
