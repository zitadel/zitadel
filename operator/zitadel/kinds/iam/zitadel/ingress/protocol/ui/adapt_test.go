package ui

/*
import (
	"fmt"
	"testing"

	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/ingress/controllers/ambassador"

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
				"resourceVersion": resourceVersion,
			},
		},
	}
	k8sClient.EXPECT().GetNamespacedCRDResource(group, version, kind, namespace, name).Return(ret, nil)
}

func SetCheckCRD(k8sClient *kubernetesmock.MockClientInt) {
	k8sClient.EXPECT().CheckCRD("mappings.getambassador.io").Times(1).Return(&apixv1beta1.CustomResourceDefinition{}, true, nil)
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
				"labels":    labels.MustK8sMap(accountsLabels),
				"name":      accountsLabels.Name(),
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
	SetReturnResourceVersion(k8sClient, group, version, kind, namespace, accountsLabels.Name(), "")
	k8sClient.EXPECT().ApplyNamespacedCRDResource(group, version, kind, namespace, accountsLabels.Name(), accounts).Times(1)

	console := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": group + "/" + version,
			"kind":       kind,
			"metadata": map[string]interface{}{
				"labels":    labels.MustK8sMap(consoleLabels),
				"name":      consoleLabels.Name(),
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
	SetReturnResourceVersion(k8sClient, group, version, kind, namespace, consoleLabels.Name(), "")
	k8sClient.EXPECT().ApplyNamespacedCRDResource(group, version, kind, namespace, consoleLabels.Name(), console).Times(1)
}

func TestUi_Adapt(t *testing.T) {
	monitor := mntr.Monitor{}
	namespace := "test"
	service := "service"
	var port uint16 = 8080
	url := fmt.Sprintf("%s:%d", service, port)
	dns := &configuration.Ingress{
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
		url,
	)

	hostAdapter := ambassador.Adapt(".")
	query, _, err := AdaptFunc(
		monitor,
		componentLabels,
		namespace,
		service,
		port,
		dns.ControllerSpecifics,
		dns.TlsSecret,
		hostAdapter,
		hostAdapter,
	)
	assert.NoError(t, err)
	queried := map[string]interface{}{}
	ensure, err := query(k8sClient, queried)
	assert.NoError(t, err)
	assert.NoError(t, ensure(k8sClient))
}

func TestUi_Adapt2(t *testing.T) {
	monitor := mntr.Monitor{}
	namespace := "test"
	service := "service"
	var port uint16 = 8080
	accountsHost := "accounts.domain"
	consoleHost := "console.domain"
	accountsAdapter := ambassador.Adapt(accountsHost)
	consoleAdapter := ambassador.Adapt(consoleHost)

	url := fmt.Sprintf("%s:%d", service, port)
	dns := &configuration.Ingress{
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
				"labels":    labels.MustK8sMap(accountsName),
				"name":      accountsName.Name(),
				"namespace": namespace,
			},
			"spec": map[string]interface{}{
				"connect_timeout_ms": 30000,
				"host":               accountsHost,
				"prefix":             "/",
				"rewrite":            "/login/",
				"service":            url,
				"timeout_ms":         30000,
			},
		},
	}
	SetReturnResourceVersion(k8sClient, group, version, kind, namespace, AccountsName, "")
	k8sClient.EXPECT().ApplyNamespacedCRDResource(group, version, kind, namespace, AccountsName, accounts).Times(1)

	consoleName := labels.MustForName(componentLabels, ConsoleName)
	console := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": group + "/" + version,
			"kind":       kind,
			"metadata": map[string]interface{}{
				"labels":    labels.MustK8sMap(consoleName),
				"name":      consoleName.Name(),
				"namespace": namespace,
			},
			"spec": map[string]interface{}{
				"host":    consoleHost,
				"prefix":  "/",
				"rewrite": "/console/",
				"service": url,
			},
		},
	}
	SetReturnResourceVersion(k8sClient, group, version, kind, namespace, ConsoleName, "")
	k8sClient.EXPECT().ApplyNamespacedCRDResource(group, version, kind, namespace, ConsoleName, console).Times(1)

	query, _, err := AdaptFunc(
		monitor,
		componentLabels,
		namespace,
		service,
		port,
		dns.ControllerSpecifics,
		dns.TlsSecret,
		consoleAdapter,
		accountsAdapter,
	)
	assert.NoError(t, err)
	queried := map[string]interface{}{}
	ensure, err := query(k8sClient, queried)
	assert.NoError(t, err)
	assert.NoError(t, ensure(k8sClient))
}
*/
