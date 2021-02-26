package protocol

import (
	"testing"

	"github.com/caos/zitadel/operator"

	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/ingress/protocol/core"

	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/configuration"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/labels/mocklabels"
)

func TestAdaptFuncCover(t *testing.T) {
	AdaptFunc(
		mntr.Monitor{},
		mocklabels.Component,
		"",
		"",
		0,
		"",
		0,
		"string",
		0,
		&configuration.Ingress{
			Subdomains: &configuration.Subdomains{},
		},
		map[string]string{},
		func(virtualHost string) core.PathAdapter {
			return func(arguments core.PathArguments) (queryFunc operator.QueryFunc, destroyFunc operator.DestroyFunc, err error) {
				return nil, nil, nil
			}
		},
	)
}

/*
import (
	"testing"

	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/ingress/controllers/ambassador"

	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/ingress/protocol/grpc"

	"github.com/caos/orbos/pkg/labels/mocklabels"

	"github.com/caos/orbos/mntr"
	kubernetesmock "github.com/caos/orbos/pkg/kubernetes/mock"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/configuration"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/ingress/protocol/http"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/ingress/protocol/ui"
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

func SetMappingsUI(
	k8sClient *kubernetesmock.MockClientInt,
	namespace string,
) {
	group := "getambassador.io"
	version := "v2"
	kind := "Mapping"
	k8sClient.EXPECT().CheckCRD("mappings.getambassador.io").Times(2).Return(&apixv1beta1.CustomResourceDefinition{}, true, nil)

	SetReturnResourceVersion(k8sClient, group, version, kind, namespace, ui.AccountsName, "")
	k8sClient.EXPECT().ApplyNamespacedCRDResource(group, version, kind, namespace, ui.AccountsName, gomock.Any()).Times(1)
	SetReturnResourceVersion(k8sClient, group, version, kind, namespace, ui.ConsoleName, "")
	k8sClient.EXPECT().ApplyNamespacedCRDResource(group, version, kind, namespace, ui.ConsoleName, gomock.Any()).Times(1)
}

func SetMappingsHTTP(
	k8sClient *kubernetesmock.MockClientInt,
	namespace string,
) {
	group := "getambassador.io"
	version := "v2"
	kind := "Mapping"
	k8sClient.EXPECT().CheckCRD("mappings.getambassador.io").Times(1).Return(&apixv1beta1.CustomResourceDefinition{}, true, nil)

	SetReturnResourceVersion(k8sClient, group, version, kind, namespace, http.AdminRName, "")
	k8sClient.EXPECT().ApplyNamespacedCRDResource(group, version, kind, namespace, http.AdminRName, gomock.Any()).Times(1)
	SetReturnResourceVersion(k8sClient, group, version, kind, namespace, http.AuthorizeName, "")
	k8sClient.EXPECT().ApplyNamespacedCRDResource(group, version, kind, namespace, http.AuthorizeName, gomock.Any()).Times(1)
	SetReturnResourceVersion(k8sClient, group, version, kind, namespace, http.AuthRName, "")
	k8sClient.EXPECT().ApplyNamespacedCRDResource(group, version, kind, namespace, http.AuthRName, gomock.Any()).Times(1)
	SetReturnResourceVersion(k8sClient, group, version, kind, namespace, http.EndsessionName, "")
	k8sClient.EXPECT().ApplyNamespacedCRDResource(group, version, kind, namespace, http.EndsessionName, gomock.Any()).Times(1)
	SetReturnResourceVersion(k8sClient, group, version, kind, namespace, http.IssuerName, "")
	k8sClient.EXPECT().ApplyNamespacedCRDResource(group, version, kind, namespace, http.IssuerName, gomock.Any()).Times(1)
	SetReturnResourceVersion(k8sClient, group, version, kind, namespace, http.MgmtName, "")
	k8sClient.EXPECT().ApplyNamespacedCRDResource(group, version, kind, namespace, http.MgmtName, gomock.Any()).Times(1)
	SetReturnResourceVersion(k8sClient, group, version, kind, namespace, http.OauthName, "")
	k8sClient.EXPECT().ApplyNamespacedCRDResource(group, version, kind, namespace, http.OauthName, gomock.Any()).Times(1)
}

func SetMappingsGRPC(
	k8sClient *kubernetesmock.MockClientInt,
	namespace string,
) {
	group := "getambassador.io"
	version := "v2"
	kind := "Mapping"
	k8sClient.EXPECT().CheckCRD("mappings.getambassador.io").Times(1).Return(&apixv1beta1.CustomResourceDefinition{}, true, nil)

	SetReturnResourceVersion(k8sClient, group, version, kind, namespace, grpc.AdminIName, "")
	k8sClient.EXPECT().ApplyNamespacedCRDResource(group, version, kind, namespace, grpc.AdminIName, gomock.Any()).Times(1)
	SetReturnResourceVersion(k8sClient, group, version, kind, namespace, grpc.AuthIName, "")
	k8sClient.EXPECT().ApplyNamespacedCRDResource(group, version, kind, namespace, grpc.AuthIName, gomock.Any()).Times(1)
	SetReturnResourceVersion(k8sClient, group, version, kind, namespace, grpc.MgmtIName, "")
	k8sClient.EXPECT().ApplyNamespacedCRDResource(group, version, kind, namespace, grpc.MgmtIName, gomock.Any()).Times(1)
}

func TestAmbassador_Adapt(t *testing.T) {

	monitor := mntr.Monitor{}
	namespace := "test"
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
	svc := "svc"
	var port uint16 = 8080
	k8sClient := kubernetesmock.NewMockClientInt(gomock.NewController(t))

	SetMappingsUI(k8sClient, namespace)
	SetMappingsHTTP(k8sClient, namespace)
	SetMappingsGRPC(k8sClient, namespace)

	query, _, err := AdaptFunc(
		monitor,
		mocklabels.Component,
		namespace,
		svc,
		port,
		svc,
		port,
		svc,
		port,
		dns,
		nil,
		ambassador.Adapt,
	)
	assert.NoError(t, err)
	queried := map[string]interface{}{}
	ensure, err := query(k8sClient, queried)
	assert.NoError(t, err)
	assert.NoError(t, ensure(k8sClient))
}
*/
