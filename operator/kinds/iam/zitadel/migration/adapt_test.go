package migration

import (
	"testing"

	"github.com/caos/orbos/pkg/labels"

	"github.com/caos/orbos/pkg/labels/mocklabels"

	"github.com/caos/orbos/mntr"
	kubernetesmock "github.com/caos/orbos/pkg/kubernetes/mock"
	"github.com/caos/zitadel/operator/helpers"
	"github.com/caos/zitadel/operator/kinds/iam/zitadel/database"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	macherrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func TestMigration_BaseEnvVars(t *testing.T) {
	envMigrationUser := "envmigration"
	migrationUser := "migration"
	envMigrationPW := "migration"
	userPasswordsSecret := "passwords"

	equals := []corev1.EnvVar{
		{
			Name:  envMigrationUser,
			Value: migrationUser,
		}, {
			Name: envMigrationPW,
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{Name: userPasswordsSecret},
					Key:                  migrationUser,
				},
			},
		},
	}

	envVars := baseEnvVars(envMigrationUser, envMigrationPW, migrationUser, userPasswordsSecret)
	assert.ElementsMatch(t, envVars, equals)
}

func TestMigration_GetMigrationFiles(t *testing.T) {
	equals := []migration{
		{
			Filename: "V1.1__test.sql",
			Data:     "test",
		},

		{
			Filename: "V1.2__test2.sql",
			Data:     "test2",
		},
	}

	files := getMigrationFiles("./testfiles")
	assert.ElementsMatch(t, equals, files)
}

func TestMigration_AdaptFunc(t *testing.T) {
	client := kubernetesmock.NewMockClientInt(gomock.NewController(t))
	namespace := "test"
	reason := "test"
	secretPasswordName := "test"
	migrationUser := "migration"
	users := []string{"test"}
	nodeselector := map[string]string{"test": "test"}
	tolerations := []corev1.Toleration{}
	dbHost := "test"
	dbPort := "test"
	localMigrationsPath := "./testfiles"

	allScripts := getMigrationFiles(localMigrationsPath)

	componentLabels := mocklabels.Component
	jobName := labels.MustForName(componentLabels, jobNamePrefix+reason)

	jobDef := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName.Name(),
			Namespace: namespace,
			Labels:    labels.MustK8sMap(jobName),
			Annotations: map[string]string{
				"migrationhash": getHash(allScripts),
			},
		},
		Spec: batchv1.JobSpec{
			Completions: helpers.PointerInt32(1),
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					NodeSelector:    nodeselector,
					Tolerations:     tolerations,
					SecurityContext: &corev1.PodSecurityContext{},
					InitContainers:  getPreContainer(dbHost, dbPort, migrationUser, secretPasswordName),
					Containers: []corev1.Container{
						getMigrationContainer(dbHost, dbPort, migrationUser, secretPasswordName, users),
					},
					RestartPolicy:                 "Never",
					DNSPolicy:                     "ClusterFirst",
					SchedulerName:                 "default-scheduler",
					TerminationGracePeriodSeconds: helpers.PointerInt64(30),
					Volumes: []corev1.Volume{{
						Name: migrationConfigmap,
						VolumeSource: corev1.VolumeSource{
							ConfigMap: &corev1.ConfigMapVolumeSource{
								LocalObjectReference: corev1.LocalObjectReference{Name: migrationConfigmap},
							},
						},
					}, {
						Name: rootUserInternal,
						VolumeSource: corev1.VolumeSource{
							Secret: &corev1.SecretVolumeSource{
								SecretName:  "cockroachdb.client.root",
								DefaultMode: helpers.PointerInt32(0400),
							},
						},
					}, {
						Name: secretPasswordName,
						VolumeSource: corev1.VolumeSource{
							Secret: &corev1.SecretVolumeSource{
								SecretName: secretPasswordName,
							},
						},
					}},
				},
			},
		},
	}

	allScriptsMap := make(map[string]string)
	for _, script := range allScripts {
		allScriptsMap[script.Filename] = script.Data
	}

	cmName := labels.MustForName(componentLabels, migrationConfigmap)
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cmName.Name(),
			Namespace: namespace,
			Labels:    labels.MustK8sMap(cmName),
		},
		Data: allScriptsMap,
	}
	client.EXPECT().ApplyJob(jobDef).Times(1)
	client.EXPECT().GetJob(namespace, getJobName(reason)).Times(1).Return(nil, macherrs.NewNotFound(schema.GroupResource{"batch", "jobs"}, jobNamePrefix+reason))
	client.EXPECT().ApplyConfigmap(cm).Times(1)

	query, _, err := AdaptFunc(
		mntr.Monitor{},
		componentLabels,
		namespace,
		reason,
		secretPasswordName,
		migrationUser,
		users,
		nodeselector,
		tolerations,
		localMigrationsPath,
	)

	queried := map[string]interface{}{}
	database.SetDatabaseInQueried(queried, &database.Current{
		Host: dbHost,
		Port: dbPort,
	})

	assert.NoError(t, err)
	ensure, err := query(client, queried)
	assert.NoError(t, err)
	assert.NoError(t, ensure(client))

}
