package services

import (
	"github.com/caos/orbos/mntr"
	kubernetesmock "github.com/caos/orbos/pkg/kubernetes/mock"
	"github.com/caos/orbos/pkg/labels"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"testing"
)

func TestService_Adapt1(t *testing.T) {
	k8sClient := kubernetesmock.NewMockClientInt(gomock.NewController(t))
	monitor := mntr.Monitor{}
	namespace := "testNs"
	componentLabels := labels.MustForComponent(labels.MustForAPI(labels.MustForOperator("testProd", "testOp", "testVersion"), "cockroachdb", "v0"), "testComponent")

	name := "testSvc"
	k8sLabels := map[string]string{
		"app.kubernetes.io/component":  "testComponent",
		"app.kubernetes.io/managed-by": "testOp",
		"app.kubernetes.io/name":       name,
		"app.kubernetes.io/part-of":    "testProd",
		"app.kubernetes.io/version":    "testVersion",
		"caos.ch/apiversion":           "v0",
		"caos.ch/kind":                 "cockroachdb",
	}
	nameLabels := labels.MustForName(componentLabels, name)
	publicName := "testPublic"
	k8sPublicLabels := map[string]string{
		"app.kubernetes.io/component":  "testComponent",
		"app.kubernetes.io/managed-by": "testOp",
		"app.kubernetes.io/name":       publicName,
		"app.kubernetes.io/part-of":    "testProd",
		"app.kubernetes.io/version":    "testVersion",
		"caos.ch/apiversion":           "v0",
		"caos.ch/kind":                 "cockroachdb",
		"orbos.ch/selectable":          "yes",
	}
	publicNameLabels := labels.MustForName(componentLabels, publicName)

	cdbName := "testCdbName"
	k8sCdbLabels := map[string]string{
		"app.kubernetes.io/component":  "testComponent",
		"app.kubernetes.io/managed-by": "testOp",
		"app.kubernetes.io/name":       cdbName,
		"app.kubernetes.io/part-of":    "testProd",
		"orbos.ch/selectable":          "yes",
	}
	cdbNameLabels := labels.MustForName(componentLabels, cdbName)

	cockroachPort := int32(25267)
	cockroachHttpPort := int32(8080)
	queried := map[string]interface{}{}

	k8sClient.EXPECT().ApplyService(&corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      publicName,
			Namespace: namespace,
			Labels:    k8sPublicLabels,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{Port: 26257, TargetPort: intstr.FromInt(int(cockroachPort)), Name: "grpc"},
				{Port: 8080, TargetPort: intstr.FromInt(int(cockroachHttpPort)), Name: "http"},
			},
			Selector:                 k8sCdbLabels,
			PublishNotReadyAddresses: false,
		},
	})

	k8sClient.EXPECT().ApplyService(&corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      publicName,
			Namespace: "default",
			Labels:    k8sPublicLabels,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{Port: 26257, TargetPort: intstr.FromInt(int(cockroachPort)), Name: "grpc"},
				{Port: 8080, TargetPort: intstr.FromInt(int(cockroachHttpPort)), Name: "http"},
			},
			Selector:                 k8sCdbLabels,
			PublishNotReadyAddresses: false,
		},
	})

	k8sClient.EXPECT().ApplyService(&corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    k8sLabels,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{Port: 26257, TargetPort: intstr.FromInt(int(cockroachPort)), Name: "grpc"},
				{Port: 8080, TargetPort: intstr.FromInt(int(cockroachHttpPort)), Name: "http"},
			},
			Selector:                 k8sCdbLabels,
			PublishNotReadyAddresses: true,
			ClusterIP:                "None",
		},
	})

	query, _, err := AdaptFunc(monitor, namespace, publicNameLabels, nameLabels, labels.DeriveNameSelector(cdbNameLabels, false), cockroachPort, cockroachHttpPort)
	assert.NoError(t, err)

	ensure, err := query(k8sClient, queried)
	assert.NoError(t, err)
	assert.NotNil(t, ensure)

	assert.NoError(t, ensure(k8sClient))
}

func TestService_Adapt2(t *testing.T) {
	k8sClient := kubernetesmock.NewMockClientInt(gomock.NewController(t))
	monitor := mntr.Monitor{}
	namespace := "testNs2"
	componentLabels := labels.MustForComponent(labels.MustForAPI(labels.MustForOperator("testProd2", "testOp2", "testVersion2"), "cockroachdb", "v0"), "testComponent2")

	name := "testSvc2"
	k8sLabels := map[string]string{
		"app.kubernetes.io/component":  "testComponent2",
		"app.kubernetes.io/managed-by": "testOp2",
		"app.kubernetes.io/name":       name,
		"app.kubernetes.io/part-of":    "testProd2",
		"app.kubernetes.io/version":    "testVersion2",
		"caos.ch/apiversion":           "v0",
		"caos.ch/kind":                 "cockroachdb",
	}
	nameLabels := labels.MustForName(componentLabels, name)
	publicName := "testPublic2"
	k8sPublicLabels := map[string]string{
		"app.kubernetes.io/component":  "testComponent2",
		"app.kubernetes.io/managed-by": "testOp2",
		"app.kubernetes.io/name":       publicName,
		"app.kubernetes.io/part-of":    "testProd2",
		"app.kubernetes.io/version":    "testVersion2",
		"caos.ch/apiversion":           "v0",
		"caos.ch/kind":                 "cockroachdb",
		"orbos.ch/selectable":          "yes",
	}
	publicNameLabels := labels.MustForName(componentLabels, publicName)

	cdbName := "testCdbName2"
	k8sCdbLabels := map[string]string{
		"app.kubernetes.io/component":  "testComponent2",
		"app.kubernetes.io/managed-by": "testOp2",
		"app.kubernetes.io/name":       cdbName,
		"app.kubernetes.io/part-of":    "testProd2",
		"orbos.ch/selectable":          "yes",
	}
	cdbNameLabels := labels.MustForName(componentLabels, cdbName)
	cockroachPort := int32(23)
	cockroachHttpPort := int32(24)
	queried := map[string]interface{}{}

	k8sClient.EXPECT().ApplyService(&corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      publicName,
			Namespace: namespace,
			Labels:    k8sPublicLabels,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{Port: 26257, TargetPort: intstr.FromInt(int(cockroachPort)), Name: "grpc"},
				{Port: 8080, TargetPort: intstr.FromInt(int(cockroachHttpPort)), Name: "http"},
			},
			Selector:                 k8sCdbLabels,
			PublishNotReadyAddresses: false,
		},
	})

	k8sClient.EXPECT().ApplyService(&corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      publicName,
			Namespace: "default",
			Labels:    k8sPublicLabels,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{Port: 26257, TargetPort: intstr.FromInt(int(cockroachPort)), Name: "grpc"},
				{Port: 8080, TargetPort: intstr.FromInt(int(cockroachHttpPort)), Name: "http"},
			},
			Selector:                 k8sCdbLabels,
			PublishNotReadyAddresses: false,
		},
	})

	k8sClient.EXPECT().ApplyService(&corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    k8sLabels,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{Port: 26257, TargetPort: intstr.FromInt(int(cockroachPort)), Name: "grpc"},
				{Port: 8080, TargetPort: intstr.FromInt(int(cockroachHttpPort)), Name: "http"},
			},
			Selector:                 k8sCdbLabels,
			PublishNotReadyAddresses: true,
			ClusterIP:                "None",
		},
	})

	query, _, err := AdaptFunc(monitor, namespace, publicNameLabels, nameLabels, labels.DeriveNameSelector(cdbNameLabels, false), cockroachPort, cockroachHttpPort)
	assert.NoError(t, err)

	ensure, err := query(k8sClient, queried)
	assert.NoError(t, err)
	assert.NotNil(t, ensure)

	assert.NoError(t, ensure(k8sClient))
}
