package ui

import (
	"github.com/caos/orbos/mntr"
	kubernetesmock "github.com/caos/orbos/pkg/kubernetes/mock"
	"github.com/caos/zitadel/operator/kinds/iam/zitadel/configuration"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	apixv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"testing"
)

func SetReturnResourceVersion(
	k8sClient *kubernetesmock.MockClientInt,
	group,
	version,
	kind,
	namespace,
	name string,
	resourceVersion string,
) {
	ret := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"metadata": map[string]interface{}{
				"resourceVersion": resourceVersion,
			},
		},
	}
	k8sClient.EXPECT().GetNamespacedCRDResource(group, version, kind, namespace, name).Return(ret, nil)
}

func SetCheckCRD(k8sClient *kubernetesmock.MockClientInt) {
	k8sClient.EXPECT().CheckCRD("mappings.getambassador.io").Times(1).Return(&apixv1beta1.CustomResourceDefinition{}, nil)
}

func SetMappingsEmpty(
	k8sClient *kubernetesmock.MockClientInt,
	namespace string,
	labels map[string]string,
	url string,
) {
	group := "getambassador.io"
	version := "v2"
	kind := "Mapping"

	accounts := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": group + "/" + version,
			"kind":       kind,
			"metadata": map[string]interface{}{
				"labels":    labels,
				"name":      AccountsName,
				"namespace": namespace,
			},
			"spec": map[string]interface{}{
				"connect_timeout_ms": 30000,
				"host":               ".",
				"prefix":             "/",
				"rewrite":            "/login/",
				"service":            url,
				"timeout_ms":         30000,
			},
		},
	}
	SetReturnResourceVersion(k8sClient, group, version, kind, namespace, AccountsName, "")
	k8sClient.EXPECT().ApplyNamespacedCRDResource(group, version, kind, namespace, AccountsName, accounts).Times(1)

	console := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": group + "/" + version,
			"kind":       kind,
			"metadata": map[string]interface{}{
				"labels":    labels,
				"name":      ConsoleName,
				"namespace": namespace,
			},
			"spec": map[string]interface{}{
				"host":    ".",
				"prefix":  "/",
				"rewrite": "/console/",
				"service": url,
			},
		},
	}
	SetReturnResourceVersion(k8sClient, group, version, kind, namespace, ConsoleName, "")
	k8sClient.EXPECT().ApplyNamespacedCRDResource(group, version, kind, namespace, ConsoleName, console).Times(1)
}

func TestUi_Adapt(t *testing.T) {
	monitor := mntr.Monitor{}
	namespace := "test"
	labels := map[string]string{"test": "test"}
	uiURL := "url"
	dns := &configuration.DNS{
		Domain:    "",
		TlsSecret: "",
		Subdomains: &configuration.Subdomains{
			Accounts: "",
			API:      "",
			Console:  "",
			Issuer:   "",
		},
	}
	k8sClient := kubernetesmock.NewMockClientInt(gomock.NewController(t))

	SetCheckCRD(k8sClient)
	SetMappingsEmpty(k8sClient, namespace, labels, uiURL)

	query, _, err := AdaptFunc(monitor, namespace, labels, uiURL, dns)
	assert.NoError(t, err)
	queried := map[string]interface{}{}
	ensure, err := query(k8sClient, queried)
	assert.NoError(t, err)
	assert.NoError(t, ensure(k8sClient))
}

func TestUi_Adapt2(t *testing.T) {
	monitor := mntr.Monitor{}
	namespace := "test"
	labels := map[string]string{"test": "test"}
	uiURL := "url"
	dns := &configuration.DNS{
		Domain:    "domain",
		TlsSecret: "tls",
		Subdomains: &configuration.Subdomains{
			Accounts: "accounts",
			API:      "api",
			Console:  "console",
			Issuer:   "issuer",
		},
	}
	k8sClient := kubernetesmock.NewMockClientInt(gomock.NewController(t))

	SetCheckCRD(k8sClient)

	group := "getambassador.io"
	version := "v2"
	kind := "Mapping"

	accounts := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": group + "/" + version,
			"kind":       kind,
			"metadata": map[string]interface{}{
				"labels":    labels,
				"name":      AccountsName,
				"namespace": namespace,
			},
			"spec": map[string]interface{}{
				"connect_timeout_ms": 30000,
				"host":               "accounts.domain",
				"prefix":             "/",
				"rewrite":            "/login/",
				"service":            uiURL,
				"timeout_ms":         30000,
			},
		},
	}
	SetReturnResourceVersion(k8sClient, group, version, kind, namespace, AccountsName, "")
	k8sClient.EXPECT().ApplyNamespacedCRDResource(group, version, kind, namespace, AccountsName, accounts).Times(1)

	console := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": group + "/" + version,
			"kind":       kind,
			"metadata": map[string]interface{}{
				"labels":    labels,
				"name":      ConsoleName,
				"namespace": namespace,
			},
			"spec": map[string]interface{}{
				"host":    "console.domain",
				"prefix":  "/",
				"rewrite": "/console/",
				"service": uiURL,
			},
		},
	}
	SetReturnResourceVersion(k8sClient, group, version, kind, namespace, ConsoleName, "")
	k8sClient.EXPECT().ApplyNamespacedCRDResource(group, version, kind, namespace, ConsoleName, console).Times(1)

	query, _, err := AdaptFunc(monitor, namespace, labels, uiURL, dns)
	assert.NoError(t, err)
	queried := map[string]interface{}{}
	ensure, err := query(k8sClient, queried)
	assert.NoError(t, err)
	assert.NoError(t, ensure(k8sClient))
}
