package rbac

import (
	"github.com/caos/orbos/mntr"
	kubernetesmock "github.com/caos/orbos/pkg/kubernetes/mock"
	"github.com/caos/orbos/pkg/labels"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestRbac_Adapt1(t *testing.T) {
	k8sClient := kubernetesmock.NewMockClientInt(gomock.NewController(t))
	monitor := mntr.Monitor{}
	namespace := "testNs"
	name := "testName"
	k8sLabels := map[string]string{
		"app.kubernetes.io/component":  "testComponent",
		"app.kubernetes.io/managed-by": "testOp",
		"app.kubernetes.io/name":       name,
		"app.kubernetes.io/part-of":    "testProd",
		"app.kubernetes.io/version":    "testVersion",
		"caos.ch/apiversion":           "v0",
		"caos.ch/kind":                 "cockroachdb",
	}
	nameLabels := labels.MustForName(labels.MustForComponent(labels.MustForAPI(labels.MustForOperator("testProd", "testOp", "testVersion"), "cockroachdb", "v0"), "testComponent"), name)

	queried := map[string]interface{}{}

	k8sClient.EXPECT().ApplyServiceAccount(&corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    k8sLabels,
		}})

	k8sClient.EXPECT().ApplyRole(&rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    k8sLabels,
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"secrets"},
				Verbs:     []string{"create", "get"},
			},
		},
	})
	k8sClient.EXPECT().ApplyClusterRole(&rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: k8sLabels,
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{"certificates.k8s.io"},
				Resources: []string{"certificatesigningrequests"},
				Verbs:     []string{"create", "get", "watch"},
			},
		},
	})
	k8sClient.EXPECT().ApplyRoleBinding(&rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    k8sLabels,
		},
		Subjects: []rbacv1.Subject{{
			Kind:      "ServiceAccount",
			Name:      name,
			Namespace: namespace,
		}},
		RoleRef: rbacv1.RoleRef{
			Name:     name,
			Kind:     "Role",
			APIGroup: "rbac.authorization.k8s.io",
		},
	})

	k8sClient.EXPECT().ApplyClusterRoleBinding(&rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: k8sLabels,
		},
		Subjects: []rbacv1.Subject{{
			Kind:      "ServiceAccount",
			Name:      name,
			Namespace: namespace,
		}},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Name:     name,
			Kind:     "ClusterRole",
		},
	})

	query, _, err := AdaptFunc(monitor, namespace, nameLabels)
	assert.NoError(t, err)

	ensure, err := query(k8sClient, queried)
	assert.NoError(t, err)
	assert.NotNil(t, ensure)

	assert.NoError(t, ensure(k8sClient))
}

func TestRbac_Adapt2(t *testing.T) {
	k8sClient := kubernetesmock.NewMockClientInt(gomock.NewController(t))
	monitor := mntr.Monitor{}
	namespace := "testNs2"
	name := "testName2"
	k8sLabels := map[string]string{
		"app.kubernetes.io/component":  "testComponent2",
		"app.kubernetes.io/managed-by": "testOp2",
		"app.kubernetes.io/name":       name,
		"app.kubernetes.io/part-of":    "testProd2",
		"app.kubernetes.io/version":    "testVersion2",
		"caos.ch/apiversion":           "v0",
		"caos.ch/kind":                 "cockroachdb",
	}
	nameLabels := labels.MustForName(labels.MustForComponent(labels.MustForAPI(labels.MustForOperator("testProd2", "testOp2", "testVersion2"), "cockroachdb", "v0"), "testComponent2"), name)

	queried := map[string]interface{}{}

	k8sClient.EXPECT().ApplyServiceAccount(&corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    k8sLabels,
		}})

	k8sClient.EXPECT().ApplyRole(&rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    k8sLabels,
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"secrets"},
				Verbs:     []string{"create", "get"},
			},
		},
	})
	k8sClient.EXPECT().ApplyClusterRole(&rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: k8sLabels,
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{"certificates.k8s.io"},
				Resources: []string{"certificatesigningrequests"},
				Verbs:     []string{"create", "get", "watch"},
			},
		},
	})
	k8sClient.EXPECT().ApplyRoleBinding(&rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    k8sLabels,
		},
		Subjects: []rbacv1.Subject{{
			Kind:      "ServiceAccount",
			Name:      name,
			Namespace: namespace,
		}},
		RoleRef: rbacv1.RoleRef{
			Name:     name,
			Kind:     "Role",
			APIGroup: "rbac.authorization.k8s.io",
		},
	})

	k8sClient.EXPECT().ApplyClusterRoleBinding(&rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: k8sLabels,
		},
		Subjects: []rbacv1.Subject{{
			Kind:      "ServiceAccount",
			Name:      name,
			Namespace: namespace,
		}},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Name:     name,
			Kind:     "ClusterRole",
		},
	})

	query, _, err := AdaptFunc(monitor, namespace, nameLabels)
	assert.NoError(t, err)

	ensure, err := query(k8sClient, queried)
	assert.NoError(t, err)
	assert.NotNil(t, ensure)

	assert.NoError(t, ensure(k8sClient))
}
