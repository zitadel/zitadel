package hosts

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

func TestHosts_AdaptFunc(t *testing.T) {

	monitor := mntr.Monitor{}
	namespace := "test"
	labels := map[string]string{"test": "test"}
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

	k8sClient.EXPECT().CheckCRD("hosts.getambassador.io").Times(1).Return(&apixv1beta1.CustomResourceDefinition{}, nil)

	group := "getambassador.io"
	version := "v2"
	kind := "Host"

	issuerHost := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"kind":       kind,
			"apiVersion": group + "/" + version,
			"metadata": map[string]interface{}{
				"name":      IssuerHostName,
				"namespace": namespace,
				"labels":    labels,
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

	consoleHost := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"kind":       kind,
			"apiVersion": group + "/" + version,
			"metadata": map[string]interface{}{
				"name":      ConsoleHostName,
				"namespace": namespace,
				"labels":    labels,
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

	apiHost := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"kind":       kind,
			"apiVersion": group + "/" + version,
			"metadata": map[string]interface{}{
				"name":      ApiHostName,
				"namespace": namespace,
				"labels":    labels,
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

	accountsHost := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"kind":       kind,
			"apiVersion": group + "/" + version,
			"metadata": map[string]interface{}{
				"name":      AccountsHostName,
				"namespace": namespace,
				"labels":    labels,
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

	query, _, err := AdaptFunc(monitor, namespace, labels, dns)
	assert.NoError(t, err)
	queried := map[string]interface{}{}
	ensure, err := query(k8sClient, queried)
	assert.NoError(t, err)
	assert.NoError(t, ensure(k8sClient))
}

func TestHosts_AdaptFunc2(t *testing.T) {
	monitor := mntr.Monitor{}
	namespace := "test"
	labels := map[string]string{"test": "test"}
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

	k8sClient.EXPECT().CheckCRD("hosts.getambassador.io").Times(1).Return(&apixv1beta1.CustomResourceDefinition{}, nil)

	group := "getambassador.io"
	version := "v2"
	kind := "Host"

	issuerHost := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"kind":       kind,
			"apiVersion": group + "/" + version,
			"metadata": map[string]interface{}{
				"name":      IssuerHostName,
				"namespace": namespace,
				"labels":    labels,
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

	consoleHost := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"kind":       kind,
			"apiVersion": group + "/" + version,
			"metadata": map[string]interface{}{
				"name":      ConsoleHostName,
				"namespace": namespace,
				"labels":    labels,
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

	apiHost := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"kind":       kind,
			"apiVersion": group + "/" + version,
			"metadata": map[string]interface{}{
				"name":      ApiHostName,
				"namespace": namespace,
				"labels":    labels,
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

	accountsHost := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"kind":       kind,
			"apiVersion": group + "/" + version,
			"metadata": map[string]interface{}{
				"name":      AccountsHostName,
				"namespace": namespace,
				"labels":    labels,
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

	query, _, err := AdaptFunc(monitor, namespace, labels, dns)
	assert.NoError(t, err)
	queried := map[string]interface{}{}
	ensure, err := query(k8sClient, queried)
	assert.NoError(t, err)
	assert.NoError(t, ensure(k8sClient))
}
