package migration

import (
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"strings"
	"testing"
)

func TestMigration_MigrationsEnvVars(t *testing.T) {

	envMigrationUser := "envmigration"
	migrationUser := "migration"
	envMigrationPW := "migration"
	userPasswordsSecret := "passwords"
	user1 := "test"
	users := []string{}

	baseEnv := baseEnvVars(envMigrationUser, envMigrationPW, migrationUser, userPasswordsSecret)

	equals := make([]corev1.EnvVar, 0)
	for _, v := range baseEnv {
		equals = append(equals, v)
	}

	envVars := migrationEnvVars(envMigrationUser, envMigrationPW, migrationUser, userPasswordsSecret, users)
	assert.ElementsMatch(t, envVars, equals)

	users = []string{user1}
	equals = make([]corev1.EnvVar, 0)
	for _, v := range baseEnv {
		equals = append(equals, v)
	}
	equals = append(equals, corev1.EnvVar{
		Name: "FLYWAY_PLACEHOLDERS_" + strings.ToUpper(user1) + "PASSWORD",
		ValueFrom: &corev1.EnvVarSource{
			SecretKeyRef: &corev1.SecretKeySelector{
				LocalObjectReference: corev1.LocalObjectReference{Name: userPasswordsSecret},
				Key:                  user1,
			},
		},
	})

	envVars = migrationEnvVars(envMigrationUser, envMigrationPW, migrationUser, userPasswordsSecret, users)
	assert.ElementsMatch(t, envVars, equals)

	user2 := "test2"
	users = []string{user1, user2}

	equals = make([]corev1.EnvVar, 0)
	for _, v := range baseEnv {
		equals = append(equals, v)
	}
	equals = append(equals, corev1.EnvVar{
		Name: "FLYWAY_PLACEHOLDERS_" + strings.ToUpper(user1) + "PASSWORD",
		ValueFrom: &corev1.EnvVarSource{
			SecretKeyRef: &corev1.SecretKeySelector{
				LocalObjectReference: corev1.LocalObjectReference{Name: userPasswordsSecret},
				Key:                  user1,
			},
		},
	})
	equals = append(equals, corev1.EnvVar{
		Name: "FLYWAY_PLACEHOLDERS_" + strings.ToUpper(user2) + "PASSWORD",
		ValueFrom: &corev1.EnvVarSource{
			SecretKeyRef: &corev1.SecretKeySelector{
				LocalObjectReference: corev1.LocalObjectReference{Name: userPasswordsSecret},
				Key:                  user2,
			},
		},
	})

	envVars = migrationEnvVars(envMigrationUser, envMigrationPW, migrationUser, userPasswordsSecret, users)
	assert.ElementsMatch(t, envVars, equals)

	user3 := "test3"
	users = []string{user1, user2, user3}

	equals = make([]corev1.EnvVar, 0)
	for _, v := range baseEnv {
		equals = append(equals, v)
	}
	equals = append(equals, corev1.EnvVar{
		Name: "FLYWAY_PLACEHOLDERS_" + strings.ToUpper(user1) + "PASSWORD",
		ValueFrom: &corev1.EnvVarSource{
			SecretKeyRef: &corev1.SecretKeySelector{
				LocalObjectReference: corev1.LocalObjectReference{Name: userPasswordsSecret},
				Key:                  user1,
			},
		},
	})
	equals = append(equals, corev1.EnvVar{
		Name: "FLYWAY_PLACEHOLDERS_" + strings.ToUpper(user2) + "PASSWORD",
		ValueFrom: &corev1.EnvVarSource{
			SecretKeyRef: &corev1.SecretKeySelector{
				LocalObjectReference: corev1.LocalObjectReference{Name: userPasswordsSecret},
				Key:                  user2,
			},
		},
	})
	equals = append(equals, corev1.EnvVar{
		Name: "FLYWAY_PLACEHOLDERS_" + strings.ToUpper(user3) + "PASSWORD",
		ValueFrom: &corev1.EnvVarSource{
			SecretKeyRef: &corev1.SecretKeySelector{
				LocalObjectReference: corev1.LocalObjectReference{Name: userPasswordsSecret},
				Key:                  user3,
			},
		},
	})

	envVars = migrationEnvVars(envMigrationUser, envMigrationPW, migrationUser, userPasswordsSecret, users)
	assert.ElementsMatch(t, envVars, equals)
}
