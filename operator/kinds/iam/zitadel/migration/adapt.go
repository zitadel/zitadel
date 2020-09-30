package migration

import (
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"github.com/caos/zitadel/operator/kinds/iam/zitadel/database"
	"strings"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/kubernetes/resources/configmap"
	"github.com/caos/orbos/pkg/kubernetes/resources/job"
	"github.com/caos/zitadel/operator"
	"github.com/caos/zitadel/operator/kinds/iam/zitadel/migration/scripts"
	"github.com/pkg/errors"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func AdaptFunc(
	monitor mntr.Monitor,
	namespace string,
	reason string,
	labels map[string]string,
	secretPasswordName string,
	migrationUser string,
	users []string,
	nodeselector map[string]string,
	tolerations []corev1.Toleration,
	repoURL, repoKey string,
) (
	operator.QueryFunc,
	operator.DestroyFunc,
	operator.EnsureFunc,
	operator.EnsureFunc,
	error,
) {
	internalMonitor := monitor.WithField("component", "migration")

	migrationConfigmap := "migrate-db"
	migrationsPath := "/migrate"
	rootUserInternal := "root"
	rootUserPath := "/certificates"
	defaultMode := int32(0400)
	envMigrationUser := "FLYWAY_USER"
	envMigrationPW := "FLYWAY_PASSWORD"
	jobName := "cockroachdb-cluster-migration-" + reason
	createFile := "create.sql"
	grantFile := "grant.sql"
	deleteFile := "delete.sql"

	destroyCM, err := configmap.AdaptFuncToDestroy(namespace, migrationConfigmap)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	destroyJ, err := job.AdaptFuncToDestroy(jobName, namespace)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	destroyers := []operator.DestroyFunc{
		operator.ResourceDestroyToZitadelDestroy(destroyJ),
		operator.ResourceDestroyToZitadelDestroy(destroyCM),
	}

	return func(k8sClient *kubernetes.Client, queried map[string]interface{}) (operator.EnsureFunc, error) {
			dbHost, dbPort, err := database.GetConnectionInfo(monitor, k8sClient, repoURL, repoKey)
			if err != nil {
				return nil, err
			}

			internalLabels := make(map[string]string, 0)
			for k, v := range labels {
				internalLabels[k] = v
			}
			internalLabels["app.kubernetes.io/component"] = "migration"

			allScripts := scripts.GetAll()
			gracePeriod := int64(30)
			completions := int32(1)

			jobDef := &batchv1.Job{
				ObjectMeta: metav1.ObjectMeta{
					Name:      jobName,
					Namespace: namespace,
					Labels:    internalLabels,
					Annotations: map[string]string{
						"migrationhash": getHash(allScripts),
					},
				},
				Spec: batchv1.JobSpec{
					Completions: &completions,
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							NodeSelector:    nodeselector,
							Tolerations:     tolerations,
							SecurityContext: &corev1.PodSecurityContext{},
							InitContainers: []corev1.Container{
								{
									Name:  "check-db-ready",
									Image: "postgres:9.6.17",
									Command: []string{
										"sh",
										"-c",
										"until pg_isready -h " + dbHost + " -p " + dbPort + "; do echo waiting for database; sleep 2; done;",
									},
									TerminationMessagePath:   corev1.TerminationMessagePathDefault,
									TerminationMessagePolicy: "File",
									ImagePullPolicy:          "IfNotPresent",
								},
								{
									Name:  "create-flyway-user",
									Image: "cockroachdb/cockroach:v20.1.5",
									Env:   baseEnvVars(envMigrationUser, envMigrationPW, migrationUser, secretPasswordName),
									VolumeMounts: []corev1.VolumeMount{{
										Name:      rootUserInternal,
										MountPath: rootUserPath,
									}},
									Command: []string{"/bin/bash", "-c", "--"},
									Args: []string{
										strings.Join([]string{
											createUserCommand(envMigrationUser, envMigrationPW, createFile),
											grantUserCommand(envMigrationUser, grantFile),
											"cockroach.sh sql --certs-dir=/certificates --host=" + dbHost + ":" + dbPort + " -e \"$(cat " + createFile + ")\" -e \"$(cat " + grantFile + ")\";",
										},
											";"),
									},
									TerminationMessagePath:   corev1.TerminationMessagePathDefault,
									TerminationMessagePolicy: "File",
									ImagePullPolicy:          "IfNotPresent",
								},
								{
									Name:  "db-migration",
									Image: "flyway/flyway:6.5.0",
									Args: []string{
										"-url=jdbc:postgresql://" + dbHost + ":" + dbPort + "/defaultdb?&sslmode=verify-full&ssl=true&sslrootcert=" + rootUserPath + "/ca.crt&sslfactory=org.postgresql.ssl.NonValidatingFactory",
										"-locations=filesystem:" + migrationsPath,
										"migrate",
									},
									Env: migrationEnvVars(envMigrationUser, envMigrationPW, migrationUser, secretPasswordName, users),
									VolumeMounts: []corev1.VolumeMount{{
										Name:      migrationConfigmap,
										MountPath: migrationsPath,
									}, {
										Name:      rootUserInternal,
										MountPath: rootUserPath,
									}},
									TerminationMessagePath:   corev1.TerminationMessagePathDefault,
									TerminationMessagePolicy: "File",
									ImagePullPolicy:          "IfNotPresent",
								},
							},
							Containers: []corev1.Container{
								{
									Name:    "delete-flyway-user",
									Image:   "cockroachdb/cockroach:v20.1.5",
									Command: []string{"/bin/bash", "-c", "--"},
									Args: []string{
										strings.Join([]string{
											deleteUserCommand(envMigrationUser, deleteFile),
											"cockroach.sh sql --certs-dir=/certificates --host=" + dbHost + ":" + dbPort + " -e \"$(cat " + deleteFile + ")\";",
										}, ";"),
									},
									Env: baseEnvVars(envMigrationUser, envMigrationPW, migrationUser, secretPasswordName),
									VolumeMounts: []corev1.VolumeMount{{
										Name:      rootUserInternal,
										MountPath: rootUserPath,
									}},
									TerminationMessagePath:   corev1.TerminationMessagePathDefault,
									TerminationMessagePolicy: "File",
									ImagePullPolicy:          "IfNotPresent",
								},
							},
							RestartPolicy:                 "Never",
							DNSPolicy:                     "ClusterFirst",
							SchedulerName:                 "default-scheduler",
							TerminationGracePeriodSeconds: &gracePeriod,
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
										DefaultMode: &defaultMode,
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

			queryCM, err := configmap.AdaptFuncToEnsure(namespace, migrationConfigmap, labels, allScripts)
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
		func(k8sClient *kubernetes.Client) error {
			internalMonitor.Info("waiting for migration to be completed")
			if err := k8sClient.WaitUntilJobCompleted(namespace, jobName, 60); err != nil {
				internalMonitor.Error(errors.Wrap(err, "error while waiting for migration to be completed"))
				return err
			}
			internalMonitor.Info("migration is completed")
			return nil
		},
		func(k8sClient *kubernetes.Client) error {
			internalMonitor.Info("cleanup migration job")
			if err := k8sClient.DeleteJob(namespace, jobName); err != nil {
				internalMonitor.Error(errors.Wrap(err, "error during job deletion"))
				return err
			}
			internalMonitor.Info("migration cleanup is completed")
			return nil
		},
		nil
}

func createUserCommand(user, pw, file string) string {
	return strings.Join([]string{
		"echo -n 'CREATE USER IF NOT EXISTS ' > " + file,
		"echo -n ${" + user + "} >> " + file,
		"echo -n ';' >> " + file,
		"echo -n 'ALTER USER ' >> " + file,
		"echo -n ${" + user + "} >> " + file,
		"echo -n ' WITH PASSWORD ' >> " + file,
		"echo -n ${" + pw + "} >> " + file,
		"echo -n ';' >> " + file,
	}, ";")
}

func grantUserCommand(user, file string) string {
	return strings.Join([]string{
		"echo -n 'GRANT admin TO ' > " + file,
		"echo -n ${" + user + "} >> " + file,
		"echo -n ' WITH ADMIN OPTION;'  >> " + file,
	}, ";")
}
func deleteUserCommand(user, file string) string {
	return strings.Join([]string{
		"echo -n 'DROP USER IF EXISTS ' > " + file,
		"echo -n ${" + user + "} >> " + file,
		"echo -n ';' >> " + file,
	}, ";")
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

func migrationEnvVars(envMigrationUser, envMigrationPW, migrationUser, userPasswordsSecret string, users []string) []corev1.EnvVar {
	envVars := baseEnvVars(envMigrationUser, envMigrationPW, migrationUser, userPasswordsSecret)

	migrationEnvVars := make([]corev1.EnvVar, 0)
	for _, v := range envVars {
		migrationEnvVars = append(migrationEnvVars, v)
	}
	for _, user := range users {
		migrationEnvVars = append(migrationEnvVars, corev1.EnvVar{
			Name: "FLYWAY_PLACEHOLDERS_" + strings.ToUpper(user) + "PASSWORD",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{Name: userPasswordsSecret},
					Key:                  user,
				},
			},
		})
	}
	return migrationEnvVars
}

func getHash(dataMap map[string]string) string {
	data, err := json.Marshal(dataMap)
	if err != nil {
		return ""
	}
	h := sha512.New()
	return base64.URLEncoding.EncodeToString(h.Sum(data))
}

//func getHash(values map[string]string) string {
//	scriptsStr := ""
//	for k, v := range values {
//		if scriptsStr == "" {
//			scriptsStr = k + ": " + v
//		} else {
//			scriptsStr = scriptsStr + "," + k + ": " + v
//		}
//	}
//
//	h := sha512.New()
//	_, err := h.Write([]byte(scriptsStr))
//	if err != nil {
//		return ""
//	}
//	hash := h.Sum(nil)
//	return base64.URLEncoding.EncodeToString(h.Sum(hash))
//}
