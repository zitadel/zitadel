package bucket

import (
	"github.com/caos/orbos/pkg/secret"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBucket_getSecretsFull(t *testing.T) {
	secrets := getSecretsMap(&desired)
	assert.Equal(t, desired.Spec.ServiceAccountJSON, secrets["serviceaccountjson"])
}

func TestBucket_getSecretsEmpty(t *testing.T) {
	secrets := getSecretsMap(&desiredWithoutSecret)
	assert.Equal(t, &secret.Secret{}, secrets["serviceaccountjson"])
}

func TestBucket_getSecretsNil(t *testing.T) {
	secrets := getSecretsMap(&desiredNil)
	assert.Equal(t, &secret.Secret{}, secrets["serviceaccountjson"])
}
