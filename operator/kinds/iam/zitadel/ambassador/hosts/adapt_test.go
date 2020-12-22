package hosts

import (
	"testing"

	"github.com/caos/orbos/mntr"
	kubernetesmock "github.com/caos/orbos/pkg/kubernetes/mock"
	"github.com/caos/orbos/pkg/labels"
	"github.com/caos/orbos/pkg/labels/mocklabels"
	"github.com/caos/zitadel/operator/kinds/iam/zitadel/configuration"
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

func TestHosts_AdaptFunc(t *testing.T) {

	monitor := mntr.Monitor{}
	namespace := "test"
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

	componentLabels := mocklabels.Component

	k8sClient := kubernetesmock.NewMockClientInt(gomock.NewController(t))

	k8sClient.EXPECT().CheckCRD("hosts.getambassador.io").Times(1).Return(&apixv1beta1.CustomResourceDefinition{}, nil)

	group := "getambassador.io"
	version := "v2"
	kind := "Host"

	issuerHostName := labels.MustForName(componentLabels, IssuerHostName)
	issuerHost := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"kind":       kind,
			"apiVersion": group + "/" + version,
			"metadata": map[string]interface{}{
				"name":      issuerHostName.Name(),
				"namespace": namespace,
				"labels":    labels.MustK8sMap(issuerHostName),
				"annotations": map[string]interface{}{
					"aes_res_changed": "true",
				},
			},
			"spec": map[string]interface{}{
				"hostname": ".",
				"acmeProvider": map[string]interface{}{
					"authority": "none",
				},
				"ambassadorId": []string{
					"default",
				},
				"selector": map[string]interface{}{
					"matchLabels": map[string]interface{}{
						"hostname": ".",
					},
				},
				"tlsSecret": map[string]interface{}{
					"name": "",
				},
			},
		}}

	SetReturnResourceVersion(k8sClient, group, version, kind, namespace, IssuerHostName, "")
	k8sClient.EXPECT().ApplyNamespacedCRDResource(group, version, kind, namespace, IssuerHostName, issuerHost).Times(1)

	consoleHostName := labels.MustForName(componentLabels, ConsoleHostName)
	consoleHost := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"kind":       kind,
			"apiVersion": group + "/" + version,
			"metadata": map[string]interface{}{
				"name":      consoleHostName.Name(),
				"namespace": namespace,
				"labels":    labels.MustK8sMap(consoleHostName),
				"annotations": map[string]interface{}{
					"aes_res_changed": "true",
				},
			},
			"spec": map[string]interface{}{
				"hostname": ".",
				"acmeProvider": map[string]interface{}{
					"authority": "none",
				},
				"ambassadorId": []string{
					"default",
				},
				"selector": map[string]interface{}{
					"matchLabels": map[string]interface{}{
						"hostname": ".",
					},
				},
				"tlsSecret": map[string]interface{}{
					"name": "",
				},
			},
		}}

	SetReturnResourceVersion(k8sClient, group, version, kind, namespace, ConsoleHostName, "")
	k8sClient.EXPECT().ApplyNamespacedCRDResource(group, version, kind, namespace, ConsoleHostName, consoleHost).Times(1)

	apiHostName := labels.MustForName(componentLabels, ApiHostName)
	apiHost := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"kind":       kind,
			"apiVersion": group + "/" + version,
			"metadata": map[string]interface{}{
				"name":      apiHostName.Name(),
				"namespace": namespace,
				"labels":    labels.MustK8sMap(apiHostName),
				"annotations": map[string]interface{}{
					"aes_res_changed": "true",
				},
			},
			"spec": map[string]interface{}{
				"hostname": ".",
				"acmeProvider": map[string]interface{}{
					"authority": "none",
				},
				"ambassadorId": []string{
					"default",
				},
				"selector": map[string]interface{}{
					"matchLabels": map[string]interface{}{
						"hostname": ".",
					},
				},
				"tlsSecret": map[string]interface{}{
					"name": "",
				},
			},
		}}

	SetReturnResourceVersion(k8sClient, group, version, kind, namespace, ApiHostName, "")
	k8sClient.EXPECT().ApplyNamespacedCRDResource(group, version, kind, namespace, ApiHostName, apiHost).Times(1)

	accountsHostName := labels.MustForName(componentLabels, AccountsHostName)
	accountsHost := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"kind":       kind,
			"apiVersion": group + "/" + version,
			"metadata": map[string]interface{}{
				"name":      accountsHostName.Name(),
				"namespace": namespace,
				"labels":    labels.MustK8sMap(accountsHostName),
				"annotations": map[string]interface{}{
					"aes_res_changed": "true",
				},
			},
			"spec": map[string]interface{}{
				"hostname": ".",
				"acmeProvider": map[string]interface{}{
					"authority": "none",
				},
				"ambassadorId": []string{
					"default",
				},
				"selector": map[string]interface{}{
					"matchLabels": map[string]interface{}{
						"hostname": ".",
					},
				},
				"tlsSecret": map[string]interface{}{
					"name": "",
				},
			},
		}}

	SetReturnResourceVersion(k8sClient, group, version, kind, namespace, AccountsHostName, "")
	k8sClient.EXPECT().ApplyNamespacedCRDResource(group, version, kind, namespace, AccountsHostName, accountsHost).Times(1)

	query, _, err := AdaptFunc(monitor, componentLabels, namespace, dns)
	assert.NoError(t, err)
	queried := map[string]interface{}{}
	ensure, err := query(k8sClient, queried)
	assert.NoError(t, err)
	assert.NoError(t, ensure(k8sClient))
}

func TestHosts_AdaptFunc2(t *testing.T) {
	monitor := mntr.Monitor{}
	namespace := "test"
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

	componentLabels := mocklabels.Component

	k8sClient := kubernetesmock.NewMockClientInt(gomock.NewController(t))

	k8sClient.EXPECT().CheckCRD("hosts.getambassador.io").Times(1).Return(&apixv1beta1.CustomResourceDefinition{}, nil)

	group := "getambassador.io"
	version := "v2"
	kind := "Host"

	issuerHostName := labels.MustForName(componentLabels, IssuerHostName)
	issuerHost := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"kind":       kind,
			"apiVersion": group + "/" + version,
			"metadata": map[string]interface{}{
				"name":      issuerHostName.Name(),
				"namespace": namespace,
				"labels":    labels.MustK8sMap(issuerHostName),
				"annotations": map[string]interface{}{
					"aes_res_changed": "true",
				},
			},
			"spec": map[string]interface{}{
				"hostname": "issuer.domain",
				"acmeProvider": map[string]interface{}{
					"authority": "none",
				},
				"ambassadorId": []string{
					"default",
				},
				"selector": map[string]interface{}{
					"matchLabels": map[string]interface{}{
						"hostname": "issuer.domain",
					},
				},
				"tlsSecret": map[string]interface{}{
					"name": "tls",
				},
			},
		}}

	SetReturnResourceVersion(k8sClient, group, version, kind, namespace, IssuerHostName, "")
	k8sClient.EXPECT().ApplyNamespacedCRDResource(group, version, kind, namespace, IssuerHostName, issuerHost).Times(1)

	consoleHostName := labels.MustForName(componentLabels, ConsoleHostName)
	consoleHost := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"kind":       kind,
			"apiVersion": group + "/" + version,
			"metadata": map[string]interface{}{
				"name":      consoleHostName.Name(),
				"namespace": namespace,
				"labels":    labels.MustK8sMap(consoleHostName),
				"annotations": map[string]interface{}{
					"aes_res_changed": "true",
				},
			},
			"spec": map[string]interface{}{
				"hostname": "console.domain",
				"acmeProvider": map[string]interface{}{
					"authority": "none",
				},
				"ambassadorId": []string{
					"default",
				},
				"selector": map[string]interface{}{
					"matchLabels": map[string]interface{}{
						"hostname": "console.domain",
					},
				},
				"tlsSecret": map[string]interface{}{
					"name": "tls",
				},
			},
		}}

	SetReturnResourceVersion(k8sClient, group, version, kind, namespace, ConsoleHostName, "")
	k8sClient.EXPECT().ApplyNamespacedCRDResource(group, version, kind, namespace, ConsoleHostName, consoleHost).Times(1)

	apiHostName := labels.MustForName(componentLabels, ApiHostName)
	apiHost := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"kind":       kind,
			"apiVersion": group + "/" + version,
			"metadata": map[string]interface{}{
				"name":      apiHostName.Name(),
				"namespace": namespace,
				"labels":    labels.MustK8sMap(apiHostName),
				"annotations": map[string]interface{}{
					"aes_res_changed": "true",
				},
			},
			"spec": map[string]interface{}{
				"hostname": "api.domain",
				"acmeProvider": map[string]interface{}{
					"authority": "none",
				},
				"ambassadorId": []string{
					"default",
				},
				"selector": map[string]interface{}{
					"matchLabels": map[string]interface{}{
						"hostname": "api.domain",
					},
				},
				"tlsSecret": map[string]interface{}{
					"name": "tls",
				},
			},
		}}

	SetReturnResourceVersion(k8sClient, group, version, kind, namespace, ApiHostName, "")
	k8sClient.EXPECT().ApplyNamespacedCRDResource(group, version, kind, namespace, ApiHostName, apiHost).Times(1)

	accountsHostName := labels.MustForName(componentLabels, AccountsHostName)
	accountsHost := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"kind":       kind,
			"apiVersion": group + "/" + version,
			"metadata": map[string]interface{}{
				"name":      accountsHostName.Name(),
				"namespace": namespace,
				"labels":    labels.MustK8sMap(accountsHostName),
				"annotations": map[string]interface{}{
					"aes_res_changed": "true",
				},
			},
			"spec": map[string]interface{}{
				"hostname": "accounts.domain",
				"acmeProvider": map[string]interface{}{
					"authority": "none",
				},
				"ambassadorId": []string{
					"default",
				},
				"selector": map[string]interface{}{
					"matchLabels": map[string]interface{}{
						"hostname": "accounts.domain",
					},
				},
				"tlsSecret": map[string]interface{}{
					"name": "tls",
				},
			},
		}}

	SetReturnResourceVersion(k8sClient, group, version, kind, namespace, AccountsHostName, "")
	k8sClient.EXPECT().ApplyNamespacedCRDResource(group, version, kind, namespace, AccountsHostName, accountsHost).Times(1)

	query, _, err := AdaptFunc(monitor, componentLabels, namespace, dns)
	assert.NoError(t, err)
	queried := map[string]interface{}{}
	ensure, err := query(k8sClient, queried)
	assert.NoError(t, err)
	assert.NoError(t, ensure(k8sClient))
}
