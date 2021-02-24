package ui

import (
	"testing"

	"github.com/caos/orbos/mntr"
	kubernetesmock "github.com/caos/orbos/pkg/kubernetes/mock"
	"github.com/caos/orbos/pkg/labels"
	"github.com/caos/orbos/pkg/labels/mocklabels"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/configuration"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	apixv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
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
				"labels":          map[string]string{},
				"annotations":     map[string]string{},
				"resourceVersion": resourceVersion,
			},
		},
	}
	k8sClient.EXPECT().GetNamespacedCRDResource(group, version, kind, namespace, name).MinTimes(1).MaxTimes(1).Return(ret, nil)
}

func SetCheckCRD(k8sClient *kubernetesmock.MockClientInt) {
	k8sClient.EXPECT().CheckCRD("mappings.getambassador.io").MinTimes(1).MaxTimes(1).Return(&apixv1beta1.CustomResourceDefinition{}, nil)
}

func SetMappingsEmpty(
	k8sClient *kubernetesmock.MockClientInt,
	namespace string,
	accountsLabels *labels.Name,
	consoleLabels *labels.Name,
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
				"labels":      labels.MustK8sMap(accountsLabels),
				"name":        accountsLabels.Name(),
				"namespace":   namespace,
				"annotations": map[string]string{},
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
	SetReturnResourceVersion(k8sClient, group, version, kind, namespace, accountsLabels.Name(), "")
	k8sClient.EXPECT().ApplyNamespacedCRDResource(group, version, kind, namespace, accountsLabels.Name(), accounts).MinTimes(1).MaxTimes(1)

	console := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": group + "/" + version,
			"kind":       kind,
			"metadata": map[string]interface{}{
				"labels":      labels.MustK8sMap(consoleLabels),
				"name":        consoleLabels.Name(),
				"namespace":   namespace,
				"annotations": map[string]string{},
			},
			"spec": map[string]interface{}{
				"host":    ".",
				"prefix":  "/",
				"rewrite": "/console/",
				"service": url,
			},
		},
	}
	SetReturnResourceVersion(k8sClient, group, version, kind, namespace, consoleLabels.Name(), "")
	k8sClient.EXPECT().ApplyNamespacedCRDResource(group, version, kind, namespace, consoleLabels.Name(), console).MinTimes(1).MaxTimes(1)
}

func TestUi_Adapt(t *testing.T) {
	monitor := mntr.Monitor{}
	namespace := "test"
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

	componentLabels := mocklabels.Component

	SetCheckCRD(k8sClient)
	SetMappingsEmpty(
		k8sClient,
		namespace,
		labels.MustForName(componentLabels, AccountsName),
		labels.MustForName(componentLabels, ConsoleName),
		uiURL,
	)

	query, _, err := AdaptFunc(monitor, componentLabels, namespace, uiURL, dns)
	assert.NoError(t, err)
	queried := map[string]interface{}{}
	ensure, err := query(k8sClient, queried)
	assert.NoError(t, err)
	assert.NoError(t, ensure(k8sClient))
}

func TestUi_Adapt2(t *testing.T) {
	monitor := mntr.Monitor{}
	namespace := "test"
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

	componentLabels := mocklabels.Component

	accountsName := labels.MustForName(componentLabels, AccountsName)
	accounts := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": group + "/" + version,
			"kind":       kind,
			"metadata": map[string]interface{}{
				"labels":      labels.MustK8sMap(accountsName),
				"name":        accountsName.Name(),
				"namespace":   namespace,
				"annotations": map[string]string{},
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
	k8sClient.EXPECT().ApplyNamespacedCRDResource(group, version, kind, namespace, AccountsName, accounts).MinTimes(1).MaxTimes(1)

	consoleName := labels.MustForName(componentLabels, ConsoleName)
	console := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": group + "/" + version,
			"kind":       kind,
			"metadata": map[string]interface{}{
				"labels":      labels.MustK8sMap(consoleName),
				"name":        consoleName.Name(),
				"namespace":   namespace,
				"annotations": map[string]string{},
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
	k8sClient.EXPECT().ApplyNamespacedCRDResource(group, version, kind, namespace, ConsoleName, console).MinTimes(1).MaxTimes(1)

	query, _, err := AdaptFunc(monitor, componentLabels, namespace, uiURL, dns)
	assert.NoError(t, err)
	queried := map[string]interface{}{}
	ensure, err := query(k8sClient, queried)
	assert.NoError(t, err)
	assert.NoError(t, ensure(k8sClient))
}
