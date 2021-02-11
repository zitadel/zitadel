package grpc

import (
	"fmt"
	"testing"

	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/ingress/controllers/ambassador"

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
	k8sClient.EXPECT().GetNamespacedCRDResource(group, version, kind, namespace, name).Return(ret, nil)
}

func TestGrpc_Adapt(t *testing.T) {
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

	componentLabels := mocklabels.Component

	k8sClient := kubernetesmock.NewMockClientInt(gomock.NewController(t))

	k8sClient.EXPECT().CheckCRD("mappings.getambassador.io").Times(1).Return(&apixv1beta1.CustomResourceDefinition{}, nil)

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
	adminMName := labels.MustForName(componentLabels, AdminIName)
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
				"prefix":             "/caos.zitadel.admin.api.v1.AdminService/",
				"rewrite":            "/caos.zitadel.admin.api.v1.AdminService/",
				"service":            url,
				"timeout_ms":         30000,
				"cors":               cors,
				"grpc":               true,
			},
		},
	}
	SetReturnResourceVersion(k8sClient, group, version, kind, namespace, AdminIName, "")
	k8sClient.EXPECT().ApplyNamespacedCRDResource(group, version, kind, namespace, AdminIName, adminM).Times(1)

	authMName := labels.MustForName(componentLabels, AuthIName)
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
				"prefix":             "/caos.zitadel.auth.api.v1.AuthService/",
				"rewrite":            "/caos.zitadel.auth.api.v1.AuthService/",
				"service":            url,
				"timeout_ms":         30000,
				"cors":               cors,
				"grpc":               true,
			},
		},
	}
	SetReturnResourceVersion(k8sClient, group, version, kind, namespace, AuthIName, "")
	k8sClient.EXPECT().ApplyNamespacedCRDResource(group, version, kind, namespace, AuthIName, authM).Times(1)

	mgmtMName := labels.MustForName(componentLabels, MgmtIName)
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
				"prefix":             "/caos.zitadel.management.api.v1.ManagementService/",
				"rewrite":            "/caos.zitadel.management.api.v1.ManagementService/",
				"service":            url,
				"timeout_ms":         30000,
				"cors":               cors,
				"grpc":               true,
			},
		},
	}
	SetReturnResourceVersion(k8sClient, group, version, kind, namespace, MgmtIName, "")
	k8sClient.EXPECT().ApplyNamespacedCRDResource(group, version, kind, namespace, MgmtIName, mgmtM).Times(1)

	query, _, err := AdaptFunc(monitor, componentLabels, namespace, "", service, port, dns, make(map[string]interface{}), ambassador.QueryMappingFunc, ambassador.DestroyMapping)
	assert.NoError(t, err)
	queried := map[string]interface{}{}
	ensure, err := query(k8sClient, queried)
	assert.NoError(t, err)
	assert.NoError(t, ensure(k8sClient))
}

func TestGrpc_Adapt2(t *testing.T) {
	monitor := mntr.Monitor{}
	namespace := "test"
	service := "service"
	var port uint16 = 8080
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

	componentLabels := mocklabels.Component

	k8sClient := kubernetesmock.NewMockClientInt(gomock.NewController(t))

	k8sClient.EXPECT().CheckCRD("mappings.getambassador.io").Times(1).Return(&apixv1beta1.CustomResourceDefinition{}, nil)

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

	adminMName := labels.MustForName(componentLabels, AdminIName)
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
				"prefix":             "/caos.zitadel.admin.api.v1.AdminService/",
				"rewrite":            "/caos.zitadel.admin.api.v1.AdminService/",
				"service":            url,
				"timeout_ms":         30000,
				"cors":               cors,
				"grpc":               true,
			},
		},
	}
	SetReturnResourceVersion(k8sClient, group, version, kind, namespace, AdminIName, "")
	k8sClient.EXPECT().ApplyNamespacedCRDResource(group, version, kind, namespace, AdminIName, adminM).Times(1)

	authMName := labels.MustForName(componentLabels, AuthIName)
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
				"prefix":             "/caos.zitadel.auth.api.v1.AuthService/",
				"rewrite":            "/caos.zitadel.auth.api.v1.AuthService/",
				"service":            url,
				"timeout_ms":         30000,
				"cors":               cors,
				"grpc":               true,
			},
		},
	}
	SetReturnResourceVersion(k8sClient, group, version, kind, namespace, AuthIName, "")
	k8sClient.EXPECT().ApplyNamespacedCRDResource(group, version, kind, namespace, AuthIName, authM).Times(1)

	mgmtMName := labels.MustForName(componentLabels, MgmtIName)
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
				"prefix":             "/caos.zitadel.management.api.v1.ManagementService/",
				"rewrite":            "/caos.zitadel.management.api.v1.ManagementService/",
				"service":            url,
				"timeout_ms":         30000,
				"cors":               cors,
				"grpc":               true,
			},
		},
	}
	SetReturnResourceVersion(k8sClient, group, version, kind, namespace, MgmtIName, "")
	k8sClient.EXPECT().ApplyNamespacedCRDResource(group, version, kind, namespace, MgmtIName, mgmtM).Times(1)

	query, _, err := AdaptFunc(monitor, componentLabels, namespace, "", service, port, dns, make(map[string]interface{}), ambassador.QueryMappingFunc, ambassador.DestroyMapping)
	assert.NoError(t, err)
	queried := map[string]interface{}{}
	ensure, err := query(k8sClient, queried)
	assert.NoError(t, err)
	assert.NoError(t, ensure(k8sClient))
}
