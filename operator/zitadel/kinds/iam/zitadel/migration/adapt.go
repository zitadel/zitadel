package migration

import (
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"

	"github.com/rakyll/statik/fs"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/kubernetes/resources/configmap"
	"github.com/caos/orbos/pkg/kubernetes/resources/job"
	"github.com/caos/orbos/pkg/labels"

	"github.com/caos/zitadel/operator"
	"github.com/caos/zitadel/operator/helpers"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/database"
	_ "github.com/caos/zitadel/statik"
)

const (
	migrationConfigmap = "migrate-db"
	migrationsPath     = "/migrate"
	rootUserInternal   = "root"
	rootUserPath       = "/certificates"
	envMigrationUser   = "FLYWAY_USER"
	envMigrationPW     = "FLYWAY_PASSWORD"
	jobNamePrefix      = "cockroachdb-cluster-migration-"
	createFile         = "create.sql"
	grantFile          = "grant.sql"
	deleteFile         = "delete.sql"
)

func AdaptFunc(
	monitor mntr.Monitor,
	componentLabels *labels.Component,
	namespace string,
	reason string,
	secretPasswordName string,
	migrationUser string,
	users []string,
	nodeselector map[string]string,
	tolerations []corev1.Toleration,
	customImageRegistry string,
) (
	operator.QueryFunc,
	operator.DestroyFunc,
	error,
) {
	internalMonitor := monitor.WithField("type", "migration")
	jobName := getJobName(reason)

	destroyCM, err := configmap.AdaptFuncToDestroy(namespace, migrationConfigmap)
	if err != nil {
		return nil, nil, err
	}

	destroyJ, err := job.AdaptFuncToDestroy(namespace, jobName)
	if err != nil {
		return nil, nil, err
	}

	destroyers := []operator.DestroyFunc{
		operator.ResourceDestroyToZitadelDestroy(destroyJ),
		operator.ResourceDestroyToZitadelDestroy(destroyCM),
	}

	return func(k8sClient kubernetes.ClientInt, queried map[string]interface{}) (operator.EnsureFunc, error) {
			dbCurrent, err := database.GetDatabaseInQueried(queried)
			if err != nil {
				return nil, err
			}
			dbHost := dbCurrent.Host
			dbPort := dbCurrent.Port

			allScripts := getMigrationFiles(monitor, "/cockroach/")

			nameLabels := labels.MustForNameK8SMap(componentLabels, jobName)
			jobDef := &batchv1.Job{
				ObjectMeta: metav1.ObjectMeta{
					Name:      jobName,
					Namespace: namespace,
					Labels:    nameLabels,
					Annotations: map[string]string{
						"migrationhash": getHash(allScripts),
					},
				},
				Spec: batchv1.JobSpec{
					Completions: helpers.PointerInt32(1),
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Annotations: map[string]string{
								"migrationhash": getHash(allScripts),
							},
						},
						Spec: corev1.PodSpec{
							SecurityContext: &corev1.PodSecurityContext{
								RunAsNonRoot: helpers.PointerBool(true),
							},
							NodeSelector:   nodeselector,
							Tolerations:    tolerations,
							InitContainers: getPreContainer(dbHost, dbPort, migrationUser, secretPasswordName, customImageRegistry),
							Containers: []corev1.Container{
								getMigrationContainer(dbHost, dbPort, migrationUser, secretPasswordName, users, customImageRegistry),
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
			queryCM, err := configmap.AdaptFuncToEnsure(namespace, migrationConfigmap, labels.MustForNameK8SMap(componentLabels, migrationConfigmap), allScriptsMap)
			if err != nil {
				return nil, err
			}
			queryJ, err := job.AdaptFuncToEnsure(jobDef)
			if err != nil {
				return nil, err
			}

			queriers := []operator.QueryFunc{
				operator.ResourceQueryToZitadelQuery(queryCM),
				operator.ResourceQueryToZitadelQuery(queryJ),
			}
			return operator.QueriersToEnsureFunc(internalMonitor, true, queriers, k8sClient, queried)
		},
		operator.DestroyersToDestroyFunc(internalMonitor, destroyers),
		nil
}

func baseEnvVars(envMigrationUser, envMigrationPW, migrationUser, userPasswordsSecret string) []corev1.EnvVar {
	envVars := []corev1.EnvVar{
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
	return envVars
}

func getHash(dataMap []migration) string {
	data, err := json.Marshal(dataMap)
	if err != nil {
		return ""
	}
	h := sha512.New()
	return base64.URLEncoding.EncodeToString(h.Sum(data))
}

type migration struct {
	Filename string
	Data     string
}

const migrationFileRegex = `(V|U)(\.|\d)+(__)(\w|\_|\ )+(\.sql)`

func getMigrationFiles(monitor mntr.Monitor, root string) []migration {
	migrations := make([]migration, 0)
	files := []string{}
	/*
		absPath, err := filepath.Abs(root)
		if err != nil {
			return migrations
		}*/

	statikFS, err := fs.New()
	if err != nil {
		monitor.Error(fmt.Errorf("failed to load migration files: %w", err))
		return migrations
	}
	err = fs.Walk(statikFS, root, func(path string, info os.FileInfo, err error) error {
		matched, err := regexp.MatchString(migrationFileRegex, info.Name())
		if err != nil {
			return err
		}
		if !info.IsDir() && matched {
			files = append(files, info.Name())
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	sort.Strings(files)

	for _, file := range files {

		fullName := filepath.Join(root, file)

		data, err := fs.ReadFile(statikFS, fullName)
		if err != nil || data == nil || len(data) == 0 {
			continue
		}
		migrations = append(migrations, migration{file, string(data)})
	}

	return migrations
}

func getJobName(reason string) string {
	return jobNamePrefix + reason
}
