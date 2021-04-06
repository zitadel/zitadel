package grpc

import (
	"testing"

	"github.com/caos/orbos/pkg/labels"

	"github.com/caos/orbos/mntr"
	kubernetesmock "github.com/caos/orbos/pkg/kubernetes/mock"
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
	k8sClient.EXPECT().GetNamespacedCRDResource(group, version, kind, namespace, name).MinTimes(1).MaxTimes(1).Return(ret, nil)
}

func TestGrpc_Adapt(t *testing.T) {
	monitor := mntr.Monitor{}
	namespace := "test"
	url := "url"
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

	k8sClient.EXPECT().CheckCRD("mappings.getambassador.io").MinTimes(1).MaxTimes(1).Return(&apixv1beta1.CustomResourceDefinition{}, nil)

	group := "getambassador.io"
	version := "v2"
	kind := "Mapping"

	cors := map[string]interface{}{
		"origins":         "*",
		"methods":         "POST, GET, OPTIONS, DELETE, PUT",
		"headers":         "*",
		"credentials":     true,
		"exposed_headers": "*",
		"max_age":         "86400",
	}
	adminMName := labels.MustForName(componentLabels, AdminMName)
	adminM := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": group + "/" + version,
			"kind":       kind,
			"metadata": map[string]interface{}{
				"labels":    labels.MustK8sMap(adminMName),
				"name":      adminMName.Name(),
				"namespace": namespace,
			},
			"spec": map[string]interface{}{
				"connect_timeout_ms": 30000,
				"host":               ".",
				"prefix":             "/zitadel.admin.v1.AdminService/",
				"rewrite":            "",
				"service":            url,
				"timeout_ms":         30000,
				"cors":               cors,
				"grpc":               true,
			},
		},
	}
	SetReturnResourceVersion(k8sClient, group, version, kind, namespace, AdminMName, "")
	k8sClient.EXPECT().ApplyNamespacedCRDResource(group, version, kind, namespace, AdminMName, adminM).MinTimes(1).MaxTimes(1)

	authMName := labels.MustForName(componentLabels, AuthMName)
	authM := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": group + "/" + version,
			"kind":       kind,
			"metadata": map[string]interface{}{
				"labels":    labels.MustK8sMap(authMName),
				"name":      authMName.Name(),
				"namespace": namespace,
			},
			"spec": map[string]interface{}{
				"connect_timeout_ms": 30000,
				"host":               ".",
				"prefix":             "/zitadel.auth.v1.AuthService/",
				"rewrite":            "",
				"service":            url,
				"timeout_ms":         30000,
				"cors":               cors,
				"grpc":               true,
			},
		},
	}
	SetReturnResourceVersion(k8sClient, group, version, kind, namespace, AuthMName, "")
	k8sClient.EXPECT().ApplyNamespacedCRDResource(group, version, kind, namespace, AuthMName, authM).MinTimes(1).MaxTimes(1)

	mgmtMName := labels.MustForName(componentLabels, MgmtMName)
	mgmtM := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": group + "/" + version,
			"kind":       kind,
			"metadata": map[string]interface{}{
				"labels":    labels.MustK8sMap(mgmtMName),
				"name":      mgmtMName.Name(),
				"namespace": namespace,
			},
			"spec": map[string]interface{}{
				"connect_timeout_ms": 30000,
				"host":               ".",
				"prefix":             "/zitadel.management.v1.ManagementService/",
				"rewrite":            "",
				"service":            url,
				"timeout_ms":         30000,
				"cors":               cors,
				"grpc":               true,
			},
		},
	}
	SetReturnResourceVersion(k8sClient, group, version, kind, namespace, MgmtMName, "")
	k8sClient.EXPECT().ApplyNamespacedCRDResource(group, version, kind, namespace, MgmtMName, mgmtM).MinTimes(1).MaxTimes(1)

	query, _, err := AdaptFunc(monitor, componentLabels, namespace, url, dns)
	assert.NoError(t, err)
	queried := map[string]interface{}{}
	ensure, err := query(k8sClient, queried)
	assert.NoError(t, err)
	assert.NoError(t, ensure(k8sClient))
}

func TestGrpc_Adapt2(t *testing.T) {
	monitor := mntr.Monitor{}
	namespace := "test"
	url := "url"
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

	k8sClient.EXPECT().CheckCRD("mappings.getambassador.io").MinTimes(1).MaxTimes(1).Return(&apixv1beta1.CustomResourceDefinition{}, nil)

	group := "getambassador.io"
	version := "v2"
	kind := "Mapping"

	cors := map[string]interface{}{
		"origins":         "*",
		"methods":         "POST, GET, OPTIONS, DELETE, PUT",
		"headers":         "*",
		"credentials":     true,
		"exposed_headers": "*",
		"max_age":         "86400",
	}

	adminMName := labels.MustForName(componentLabels, AdminMName)
	adminM := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": group + "/" + version,
			"kind":       kind,
			"metadata": map[string]interface{}{
				"labels":    labels.MustK8sMap(adminMName),
				"name":      adminMName.Name(),
				"namespace": namespace,
			},
			"spec": map[string]interface{}{
				"connect_timeout_ms": 30000,
				"host":               "api.domain",
				"prefix":             "/zitadel.admin.v1.AdminService/",
				"rewrite":            "",
				"service":            url,
				"timeout_ms":         30000,
				"cors":               cors,
				"grpc":               true,
			},
		},
	}
	SetReturnResourceVersion(k8sClient, group, version, kind, namespace, AdminMName, "")
	k8sClient.EXPECT().ApplyNamespacedCRDResource(group, version, kind, namespace, AdminMName, adminM).MinTimes(1).MaxTimes(1)

	authMName := labels.MustForName(componentLabels, AuthMName)
	authM := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": group + "/" + version,
			"kind":       kind,
			"metadata": map[string]interface{}{
				"labels":    labels.MustK8sMap(authMName),
				"name":      authMName.Name(),
				"namespace": namespace,
			},
			"spec": map[string]interface{}{
				"connect_timeout_ms": 30000,
				"host":               "api.domain",
				"prefix":             "/zitadel.auth.v1.AuthService/",
				"rewrite":            "",
				"service":            url,
				"timeout_ms":         30000,
				"cors":               cors,
				"grpc":               true,
			},
		},
	}
	SetReturnResourceVersion(k8sClient, group, version, kind, namespace, AuthMName, "")
	k8sClient.EXPECT().ApplyNamespacedCRDResource(group, version, kind, namespace, AuthMName, authM).MinTimes(1).MaxTimes(1)

	mgmtMName := labels.MustForName(componentLabels, MgmtMName)
	mgmtM := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": group + "/" + version,
			"kind":       kind,
			"metadata": map[string]interface{}{
				"labels":    labels.MustK8sMap(mgmtMName),
				"name":      mgmtMName.Name(),
				"namespace": namespace,
			},
			"spec": map[string]interface{}{
				"connect_timeout_ms": 30000,
				"host":               "api.domain",
				"prefix":             "/zitadel.management.v1.ManagementService/",
				"rewrite":            "",
				"service":            url,
				"timeout_ms":         30000,
				"cors":               cors,
				"grpc":               true,
			},
		},
	}
	SetReturnResourceVersion(k8sClient, group, version, kind, namespace, MgmtMName, "")
	k8sClient.EXPECT().ApplyNamespacedCRDResource(group, version, kind, namespace, MgmtMName, mgmtM).MinTimes(1).MaxTimes(1)

	query, _, err := AdaptFunc(monitor, componentLabels, namespace, url, dns)
	assert.NoError(t, err)
	queried := map[string]interface{}{}
	ensure, err := query(k8sClient, queried)
	assert.NoError(t, err)
	assert.NoError(t, ensure(k8sClient))
}
