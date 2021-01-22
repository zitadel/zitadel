package bucket

import (
	"github.com/caos/orbos/pkg/secret"
	"github.com/caos/orbos/pkg/tree"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"testing"
)

const (
	masterkey  = "testMk"
	cron       = "testCron"
	bucketName = "testBucket"
	saJson     = "testSa"
	yamlFile   = `kind: databases.caos.ch/BucketBackup
version: v0
spec:
    verbose: true
    cron: testCron
    bucket: testBucket
    serviceAccountJSON:
        encryption: AES256
        encoding: Base64
        value: luyAqtopzwLcaIhJj7KhWmbUsA7cQg==
`

	yamlFileWithoutSecret = `kind: databases.caos.ch/BucketBackup
version: v0
spec:
    verbose: true
    cron: testCron
    bucket: testBucket
`
	yamlEmpty = `kind: databases.caos.ch/BucketBackup
version: v0`
)

var (
	desired = DesiredV0{
		Common: &tree.Common{
			Kind:    "databases.caos.ch/BucketBackup",
			Version: "v0",
		},
		Spec: &Spec{
			Verbose: true,
			Cron:    cron,
			Bucket:  bucketName,
			ServiceAccountJSON: &secret.Secret{
				Value:      saJson,
				Encryption: "AES256",
				Encoding:   "Base64",
			},
		},
	}
	desiredWithoutSecret = DesiredV0{
		Common: &tree.Common{
			Kind:    "databases.caos.ch/BucketBackup",
			Version: "v0",
		},
		Spec: &Spec{
			Verbose: true,
			Cron:    cron,
			Bucket:  bucketName,
		},
	}
	desiredEmpty = DesiredV0{
		Common: &tree.Common{
			Kind:    "databases.caos.ch/BucketBackup",
			Version: "v0",
		},
		Spec: &Spec{
			Verbose: false,
			Cron:    "",
			Bucket:  "",
			ServiceAccountJSON: &secret.Secret{
				Value: "",
			},
		},
	}

	desiredNil = DesiredV0{
		Common: &tree.Common{
			Kind:    "databases.caos.ch/BucketBackup",
			Version: "v0",
		},
	}
)

func marshalYaml(t *testing.T, masterkey string, struc *DesiredV0) []byte {
	secret.Masterkey = masterkey
	data, err := yaml.Marshal(struc)
	assert.NoError(t, err)
	return data
}

func unmarshalYaml(t *testing.T, masterkey string, yamlFile []byte) *tree.Tree {
	secret.Masterkey = masterkey
	desiredTree := &tree.Tree{}
	assert.NoError(t, yaml.Unmarshal(yamlFile, desiredTree))
	return desiredTree
}

func getDesiredTree(t *testing.T, masterkey string, desired *DesiredV0) *tree.Tree {
	return unmarshalYaml(t, masterkey, marshalYaml(t, masterkey, desired))
}

func TestBucket_DesiredParse(t *testing.T) {
	assert.Equal(t, yamlFileWithoutSecret, string(marshalYaml(t, masterkey, &desiredWithoutSecret)))

	desiredTree := unmarshalYaml(t, masterkey, []byte(yamlFile))
	desiredKind, err := ParseDesiredV0(desiredTree)
	assert.NoError(t, err)
	assert.Equal(t, &desired, desiredKind)
}

func TestBucket_DesiredNotZero(t *testing.T) {
	desiredTree := unmarshalYaml(t, masterkey, []byte(yamlFile))
	desiredKind, err := ParseDesiredV0(desiredTree)
	assert.NoError(t, err)
	assert.False(t, desiredKind.Spec.IsZero())
}

func TestBucket_DesiredZero(t *testing.T) {
	desiredTree := unmarshalYaml(t, masterkey, []byte(yamlEmpty))
	desiredKind, err := ParseDesiredV0(desiredTree)
	assert.NoError(t, err)
	assert.True(t, desiredKind.Spec.IsZero())
}
