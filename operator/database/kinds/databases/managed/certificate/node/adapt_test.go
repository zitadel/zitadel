package node

import (
	"github.com/caos/orbos/mntr"
	kubernetesmock "github.com/caos/orbos/pkg/kubernetes/mock"
	"github.com/caos/orbos/pkg/labels"
	"github.com/caos/orbos/pkg/tree"
	"github.com/caos/zitadel/operator/database/kinds/databases/core"
	coremock "github.com/caos/zitadel/operator/database/kinds/databases/core/mock"
	"github.com/caos/zitadel/operator/database/kinds/databases/managed/certificate/pem"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

const (
	caPem = `-----BEGIN CERTIFICATE-----
MIIE9TCCAt2gAwIBAgICB+MwDQYJKoZIhvcNAQELBQAwKzESMBAGA1UEChMJQ29j
a3JvYWNoMRUwEwYDVQQDEwxDb2Nrcm9hY2ggQ0EwHhcNMjAxMjAxMTAyMjA0WhcN
MzAxMjAxMTAyMjA0WjArMRIwEAYDVQQKEwlDb2Nrcm9hY2gxFTATBgNVBAMTDENv
Y2tyb2FjaCBDQTCCAiIwDQYJKoZIhvcNAQEBBQADggIPADCCAgoCggIBAOja5IXJ
GUY9sFgyvWkav+O74gzcv8y69uJzSOx0piOP+sfpZWVeGEjqO4JgdcS5NPMrT5Tb
aJ52CjiOdlHVTyR87i5JmOIvA2qA4dWQmSbX7AQ8r9ptYDe9xMn+qFegcbr4YxCz
K9mmbDZUhlLO7cz3QV6nvRxGFWbzffo8BXZnOUCAOyOHrbnPpLumnfZlL5BckdtY
pS7jAlUpKSMBTK4AHcmrouFsNKHqlUopYXeJFdg9g1F0DuCnVP9x7+XcUW8dAVut
Q7Jswy+++GAXs6mPVsYFLXUSYNyW+Bfl/jKwx0XTQx/6iyNpK0XtAzjBZFjXwmot
0mODkqnfE3BB4lXxZ5knomAQEGSScUhCUb9upbF4uJJF27xr/kIkwtWxMGpCXds0
IxI+wNRCenhfFZEIQCzri0zn6WdN8b/gbv1BErNcccYwolPUv1oUgYbzowbQ5O2D
aLQPqO1VAZiZHLxb787bRywpCl33VZ1ptMHi2ogKjcsh4DQ9SsRj+rU/Tk5lyk7G
FHteyHcq12TGpz9/CQYMacl8yeRRHfNO3Rq0jFTYeD4+ZdVBPKeuTHyXGzy2T157
pgqMFzwqxlNYPzpuz7xsZBExJzCtcomB8fMlsCJnxV/kuMTTsWsPrRc7hsqzBCq3
plfYT8a7EBCVmJmDu+6mMVh+A9M2zZCqtV57AgMBAAGjIzAhMA4GA1UdDwEB/wQE
AwICpDAPBgNVHRMBAf8EBTADAQH/MA0GCSqGSIb3DQEBCwUAA4ICAQDRmkBYDk/F
lYTgMaL7yhRAowDR8mtOl06Oba/MCDDsmBvfOVIqY9Sa9sJtb78Uvecv9WAE4qpb
PNSaIvNmujmGQrBxtGwQdeHZvtXQ1lKAGBM+ukfz/NFU+6rwOdTuLgS5jTaiM4r1
uMoIG6S5z/brWGHAw8Qv2Fq05dSLnevmPWK903V/RnwmCKT9as5Qpfkaie2PNory
euxVGeGolxzgbSps3ZJlDSSymQBl10iJYYdGIsgcw5FHcCdpS4AutdpQbIFuCGk2
CHTcTTa4KMaVfRm/gKm1gh8UnuVDLBKQS8+1CHFYUnih8ozBNUyBo5r5L4BHJK/n
f0gnrYqaxtqPAyHkfo0PRc50HlAQRS/gW8ESv2PQmcSWlDggEzXt5MqFIYOfEXWE
gtC7Ct1P7gzIRxolYVsNSgBR62sJM1PUa39E5v0QJsmuM51txHuw/PGqkkBlEwx9
4u3FFXD9Pslg/g1l50Ek5O3/TWMta4Q+NagDDaOdkb4LRQr522ySgGIinBw00Xr+
Ics2L42rP1LQyaE7+BljDEwmGaT7mrl8QGZR9uv/9mbtk2+q7wAOQzcMSFFACfwx
KdKENmK0GpQZ6HaeDZO7wxGiAyFKhTNpd39XNOEhaKpWLSUxnZTLGlMIzvSyMGCj
PV3xm+Xg1Xvki7f/qC1djysMcaokoJzNdg==
-----END CERTIFICATE-----
`
	caPrivPem = `-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----
`
)

func TestNode_AdaptWithCA(t *testing.T) {
	monitor := mntr.Monitor{}
	namespace := "testNs"
	clusterDns := "testDns"
	nameLabels := labels.MustForName(labels.MustForComponent(labels.MustForAPI(labels.MustForOperator("testProd", "testOp", "testVersion"), "cockroachdb", "v0"), "cockroachdb"), "testNode")

	nodeLabels := map[string]string{
		"app.kubernetes.io/component":  "cockroachdb",
		"app.kubernetes.io/managed-by": "testOp",
		"app.kubernetes.io/name":       "testNode",
		"app.kubernetes.io/part-of":    "testProd",
		"orbos.ch/selectable":          "yes",
	}
	dbCurrent := coremock.NewMockDatabaseCurrent(gomock.NewController(t))
	k8sClient := kubernetesmock.NewMockClientInt(gomock.NewController(t))

	secretList := &corev1.SecretList{
		Items: []corev1.Secret{},
	}
	ca, err := pem.DecodeCertificate([]byte(caPem))
	assert.NoError(t, err)
	caPriv, err := pem.DecodeKey([]byte(caPrivPem))
	assert.NoError(t, err)

	k8sClient.EXPECT().ListSecrets(namespace, nodeLabels).Times(1).Return(secretList, nil)
	dbCurrent.EXPECT().GetCertificate().Times(1).Return(ca)
	dbCurrent.EXPECT().GetCertificateKey().Times(1).Return(caPriv)
	dbCurrent.EXPECT().SetCertificate(ca).Times(1)
	dbCurrent.EXPECT().SetCertificateKey(caPriv).Times(1)

	queried := map[string]interface{}{}
	current := &tree.Tree{
		Parsed: dbCurrent,
	}

	core.SetQueriedForDatabase(queried, current)

	query, _, err := AdaptFunc(monitor, namespace, nameLabels, clusterDns, true)
	assert.NoError(t, err)

	ensure, err := query(k8sClient, queried)
	assert.NoError(t, err)
	assert.NotNil(t, ensure)

	assert.NoError(t, ensure(k8sClient))
}

func TestNode_AdaptWithoutCA(t *testing.T) {
	monitor := mntr.Monitor{}
	namespace := "testNs"
	clusterDns := "testDns"
	nameLabels := labels.MustForName(labels.MustForComponent(labels.MustForAPI(labels.MustForOperator("testProd", "testOp", "testVersion"), "cockroachdb", "v0"), "cockroachdb"), "testNode")

	nodeLabels := map[string]string{
		"app.kubernetes.io/component":  "cockroachdb",
		"app.kubernetes.io/managed-by": "testOp",
		"app.kubernetes.io/name":       "testNode",
		"app.kubernetes.io/part-of":    "testProd",
		"orbos.ch/selectable":          "yes",
	}
	dbCurrent := coremock.NewMockDatabaseCurrent(gomock.NewController(t))
	k8sClient := kubernetesmock.NewMockClientInt(gomock.NewController(t))

	secretList := &corev1.SecretList{
		Items: []corev1.Secret{},
	}

	k8sClient.EXPECT().ListSecrets(namespace, nodeLabels).Times(1).Return(secretList, nil)
	dbCurrent.EXPECT().GetCertificate().Times(1).Return(nil)
	dbCurrent.EXPECT().GetCertificateKey().Times(1).Return(nil)
	dbCurrent.EXPECT().SetCertificate(gomock.Any()).Times(1)
	dbCurrent.EXPECT().SetCertificateKey(gomock.Any()).Times(1)
	k8sClient.EXPECT().ApplySecret(gomock.Any()).Times(1).Return(nil)

	queried := map[string]interface{}{}
	current := &tree.Tree{
		Parsed: dbCurrent,
	}

	core.SetQueriedForDatabase(queried, current)

	query, _, err := AdaptFunc(monitor, namespace, nameLabels, clusterDns, true)
	assert.NoError(t, err)

	ensure, err := query(k8sClient, queried)
	assert.NoError(t, err)
	assert.NotNil(t, ensure)

	assert.NoError(t, ensure(k8sClient))
}

func TestNode_AdaptAlreadyExisting(t *testing.T) {
	monitor := mntr.Monitor{}
	namespace := "testNs"
	clusterDns := "testDns"
	nameLabels := labels.MustForName(labels.MustForComponent(labels.MustForAPI(labels.MustForOperator("testProd", "testOp", "testVersion"), "cockroachdb", "v0"), "cockroachdb"), "testNode")
	nodeLabels := map[string]string{
		"app.kubernetes.io/component":  "cockroachdb",
		"app.kubernetes.io/managed-by": "testOp",
		"app.kubernetes.io/name":       "testNode",
		"app.kubernetes.io/part-of":    "testProd",
		"orbos.ch/selectable":          "yes",
	}
	dbCurrent := coremock.NewMockDatabaseCurrent(gomock.NewController(t))
	k8sClient := kubernetesmock.NewMockClientInt(gomock.NewController(t))

	secretList := &corev1.SecretList{
		Items: []corev1.Secret{{
			ObjectMeta: metav1.ObjectMeta{},
			Data: map[string][]byte{
				caCertKey:    []byte(caPem),
				caPrivKeyKey: []byte(caPrivPem),
			},
			Type: "Opaque",
		}},
	}

	caCert, err := pem.DecodeCertificate([]byte(caPem))
	assert.NoError(t, err)
	caPrivKey, err := pem.DecodeKey([]byte(caPrivPem))
	assert.NoError(t, err)

	dbCurrent.EXPECT().SetCertificate(caCert).Times(1)
	dbCurrent.EXPECT().SetCertificateKey(caPrivKey).Times(1)
	k8sClient.EXPECT().ListSecrets(namespace, nodeLabels).Times(1).Return(secretList, nil)

	queried := map[string]interface{}{}
	current := &tree.Tree{
		Parsed: dbCurrent,
	}

	core.SetQueriedForDatabase(queried, current)

	query, _, err := AdaptFunc(monitor, namespace, nameLabels, clusterDns, true)
	assert.NoError(t, err)

	ensure, err := query(k8sClient, queried)
	assert.NoError(t, err)
	assert.NotNil(t, ensure)

	assert.NoError(t, ensure(k8sClient))

}
