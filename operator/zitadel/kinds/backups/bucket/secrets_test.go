package bucket

import (
	"testing"

	"github.com/caos/orbos/pkg/secret"
	"github.com/stretchr/testify/assert"
)

func TestBucket_getSecretsFull(t *testing.T) {
	secrets, existing := getSecretsMap(&desired)
	assert.Equal(t, desired.Spec.ServiceAccountJSON, secrets["serviceaccountjson"])
	assert.Equal(t, desired.Spec.ExistingServiceAccountJSON, existing["serviceaccountjson"])
}

func TestBucket_getSecretsEmpty(t *testing.T) {
	secrets, existing := getSecretsMap(&desiredWithoutSecret)
	assert.Equal(t, &secret.Secret{}, secrets["serviceaccountjson"])
	assert.Equal(t, &secret.Existing{}, existing["serviceaccountjson"])
}

func TestBucket_getSecretsNil(t *testing.T) {
	secrets, existing := getSecretsMap(&desiredNil)
	assert.Equal(t, &secret.Secret{}, secrets["serviceaccountjson"])
	assert.Equal(t, &secret.Existing{}, existing["serviceaccountjson"])
}
