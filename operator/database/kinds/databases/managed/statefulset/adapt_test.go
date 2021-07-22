package statefulset

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes/k8s"
	kubernetesmock "github.com/caos/orbos/pkg/kubernetes/mock"
	"github.com/caos/orbos/pkg/labels"
	"github.com/caos/zitadel/operator/helpers"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"testing"
)

func TestStatefulset_JoinExec0(t *testing.T) {
	namespace := "testNs"
	name := "test"
	dbPort := 26257
	replicaCount := 0

	equals := "exec /cockroach/cockroach start --logtostderr --certs-dir /cockroach/cockroach-certs --advertise-host $(hostname -f) --http-addr 0.0.0.0 --join  --locality zone=testNs --cache 25% --max-sql-memory 25%"
	assert.Equal(t, equals, getJoinExec(namespace, name, dbPort, replicaCount))
}

func TestStatefulset_JoinExec1(t *testing.T) {
	namespace := "testNs2"
	name := "test2"
	dbPort := 26257
	replicaCount := 1

	equals := "exec /cockroach/cockroach start --logtostderr --certs-dir /cockroach/cockroach-certs --advertise-host $(hostname -f) --http-addr 0.0.0.0 --join test2-0.test2.testNs2:26257 --locality zone=testNs2 --cache 25% --max-sql-memory 25%"
	assert.Equal(t, equals, getJoinExec(namespace, name, dbPort, replicaCount))
}

func TestStatefulset_JoinExec2(t *testing.T) {
	namespace := "testNs"
	name := "test"
	dbPort := 23
	replicaCount := 2

	equals := "exec /cockroach/cockroach start --logtostderr --certs-dir /cockroach/cockroach-certs --advertise-host $(hostname -f) --http-addr 0.0.0.0 --join test-0.test.testNs:23,test-1.test.testNs:23 --locality zone=testNs --cache 25% --max-sql-memory 25%"
	assert.Equal(t, equals, getJoinExec(namespace, name, dbPort, replicaCount))
}

func TestStatefulset_Resources0(t *testing.T) {
	equals := corev1.ResourceRequirements{
		Requests: corev1.ResourceList{
			"cpu":    resource.MustParse("1"),
			"memory": resource.MustParse("6Gi"),
		},
		Limits: corev1.ResourceList{
			"cpu":    resource.MustParse("4"),
			"memory": resource.MustParse("8Gi"),
		},
	}

	assert.Equal(t, equals, getResources(nil))
}

func TestStatefulset_Resources1(t *testing.T) {
	res := &k8s.Resources{
		Requests: corev1.ResourceList{
			"cpu":    resource.MustParse("200m"),
			"memory": resource.MustParse("600Mi"),
		},
		Limits: corev1.ResourceList{
			"cpu":    resource.MustParse("500m"),
			"memory": resource.MustParse("126Mi"),
		},
	}

	equals := corev1.ResourceRequirements{
		Requests: corev1.ResourceList{
			"cpu":    resource.MustParse("200m"),
			"memory": resource.MustParse("600Mi"),
		},
		Limits: corev1.ResourceList{
			"cpu":    resource.MustParse("500m"),
			"memory": resource.MustParse("126Mi"),
		},
	}

	assert.Equal(t, equals, getResources(res))
}

func TestStatefulset_Resources2(t *testing.T) {
	res := &k8s.Resources{
		Requests: corev1.ResourceList{
			"cpu":    resource.MustParse("300m"),
			"memory": resource.MustParse("670Mi"),
		},
		Limits: corev1.ResourceList{
			"cpu":    resource.MustParse("600m"),
			"memory": resource.MustParse("256Mi"),
		},
	}

	equals := corev1.ResourceRequirements{
		Requests: corev1.ResourceList{
			"cpu":    resource.MustParse("300m"),
			"memory": resource.MustParse("670Mi"),
		},
		Limits: corev1.ResourceList{
			"cpu":    resource.MustParse("600m"),
			"memory": resource.MustParse("256Mi"),
		},
	}

	assert.Equal(t, equals, getResources(res))
}

func TestStatefulset_Adapt1(t *testing.T) {
	k8sClient := kubernetesmock.NewMockClientInt(gomock.NewController(t))
	monitor := mntr.Monitor{}
	namespace := "testNs"
	name := "test"
	image := "cockroach"
	nameLabels := labels.MustForName(labels.MustForComponent(labels.MustForAPI(labels.MustForOperator("testProd", "testOp", "testVersion"), "cockroachdb", "v0"), "testComponent"), name)
	k8sSelectableLabels := map[string]string{
		"app.kubernetes.io/component":  "testComponent",
		"app.kubernetes.io/managed-by": "testOp",
		"app.kubernetes.io/name":       name,
		"app.kubernetes.io/part-of":    "testProd",
		"app.kubernetes.io/version":    "testVersion",
		"caos.ch/apiversion":           "v0",
		"caos.ch/kind":                 "cockroachdb",
		"orbos.ch/selectable":          "yes",
	}
	k8sSelectorLabels := map[string]string{
		"app.kubernetes.io/component":  "testComponent",
		"app.kubernetes.io/managed-by": "testOp",
		"app.kubernetes.io/name":       name,
		"app.kubernetes.io/part-of":    "testProd",
		"orbos.ch/selectable":          "yes",
	}
	selector := labels.DeriveNameSelector(nameLabels, false)
	selectable := labels.AsSelectable(nameLabels)

	serviceAccountName := "testSA"
	replicaCount := 1
	storageCapacity := "20Gi"
	dbPort := int32(26257)
	httpPort := int32(8080)
	storageClass := "testSC"
	nodeSelector := map[string]string{}
	tolerations := []corev1.Toleration{}
	resourcesSFS := &k8s.Resources{}

	quantity, err := resource.ParseQuantity(storageCapacity)
	assert.NoError(t, err)

	sfs := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    k8sSelectableLabels,
		},
		Spec: appsv1.StatefulSetSpec{
			ServiceName: name,
			Replicas:    helpers.PointerInt32(int32(replicaCount)),
			Selector: &metav1.LabelSelector{
				MatchLabels: k8sSelectorLabels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: k8sSelectableLabels,
				},
				Spec: corev1.PodSpec{
					NodeSelector:       nodeSelector,
					Tolerations:        tolerations,
					ServiceAccountName: serviceAccountName,
					Affinity:           getAffinity(k8sSelectableLabels),
					Containers: []corev1.Container{{
						Name:            name,
						Image:           image,
						ImagePullPolicy: "IfNotPresent",
						Ports: []corev1.ContainerPort{
							{ContainerPort: dbPort, Name: "grpc"},
							{ContainerPort: httpPort, Name: "http"},
						},
						LivenessProbe: &corev1.Probe{
							Handler: corev1.Handler{
								HTTPGet: &corev1.HTTPGetAction{
									Path:   "/health",
									Port:   intstr.Parse("http"),
									Scheme: "HTTPS",
								},
							},
							InitialDelaySeconds: 30,
							PeriodSeconds:       5,
						},
						ReadinessProbe: &corev1.Probe{
							Handler: corev1.Handler{
								HTTPGet: &corev1.HTTPGetAction{
									Path:   "/health?ready=1",
									Port:   intstr.Parse("http"),
									Scheme: "HTTPS",
								},
							},
							InitialDelaySeconds: 10,
							PeriodSeconds:       5,
							FailureThreshold:    2,
						},
						VolumeMounts: []corev1.VolumeMount{{
							Name:      datadirInternal,
							MountPath: datadirPath,
						}, {
							Name:      certsInternal,
							MountPath: certPath,
						}, {
							Name:      clientCertsInternal,
							MountPath: clientCertPath,
						}},
						Env: []corev1.EnvVar{{
							Name:  "COCKROACH_CHANNEL",
							Value: "kubernetes-multiregion",
						}},
						Command: []string{
							"/bin/bash",
							"-ecx",
							getJoinExec(
								namespace,
								name,
								int(dbPort),
								replicaCount,
							),
						},
						Resources: getResources(resourcesSFS),
					}},
					Volumes: []corev1.Volume{{
						Name: datadirInternal,
						VolumeSource: corev1.VolumeSource{
							PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
								ClaimName: datadirInternal,
							},
						},
					}, {
						Name: certsInternal,
						VolumeSource: corev1.VolumeSource{
							Secret: &corev1.SecretVolumeSource{
								SecretName:  nodeSecret,
								DefaultMode: helpers.PointerInt32(defaultMode),
							},
						},
					}, {
						Name: clientCertsInternal,
						VolumeSource: corev1.VolumeSource{
							Secret: &corev1.SecretVolumeSource{
								SecretName:  rootSecret,
								DefaultMode: helpers.PointerInt32(defaultMode),
							},
						},
					}},
				},
			},
			PodManagementPolicy: appsv1.PodManagementPolicyType("Parallel"),
			UpdateStrategy: appsv1.StatefulSetUpdateStrategy{
				Type: "RollingUpdate",
			},
			VolumeClaimTemplates: []corev1.PersistentVolumeClaim{{
				ObjectMeta: metav1.ObjectMeta{
					Name: datadirInternal,
				},
				Spec: corev1.PersistentVolumeClaimSpec{
					AccessModes: []corev1.PersistentVolumeAccessMode{
						corev1.PersistentVolumeAccessMode("ReadWriteOnce"),
					},
					Resources: corev1.ResourceRequirements{
						Requests: corev1.ResourceList{
							"storage": quantity,
						},
					},
					StorageClassName: &storageClass,
				},
			}},
		},
	}

	k8sClient.EXPECT().ApplyStatefulSet(sfs, false)

	query, _, _, _, _, err := AdaptFunc(
		monitor,
		selectable,
		selector,
		false,
		namespace,
		image,
		serviceAccountName,
		replicaCount,
		storageCapacity,
		dbPort,
		httpPort,
		storageClass,
		nodeSelector,
		tolerations,
		resourcesSFS,
	)
	assert.NoError(t, err)

	ensure, err := query(k8sClient)
	assert.NoError(t, err)
	assert.NotNil(t, ensure)

	assert.NoError(t, ensure(k8sClient))
}

func TestStatefulset_Adapt2(t *testing.T) {
	k8sClient := kubernetesmock.NewMockClientInt(gomock.NewController(t))
	monitor := mntr.Monitor{}
	namespace := "testNs2"
	name := "test2"
	image := "cockroach2"

	nameLabels := labels.MustForName(labels.MustForComponent(labels.MustForAPI(labels.MustForOperator("testProd2", "testOp2", "testVersion2"), "cockroachdb", "v0"), "testComponent2"), name)
	k8sSelectableLabels := map[string]string{
		"app.kubernetes.io/component":  "testComponent2",
		"app.kubernetes.io/managed-by": "testOp2",
		"app.kubernetes.io/name":       name,
		"app.kubernetes.io/part-of":    "testProd2",
		"app.kubernetes.io/version":    "testVersion2",
		"caos.ch/apiversion":           "v0",
		"caos.ch/kind":                 "cockroachdb",
		"orbos.ch/selectable":          "yes",
	}
	k8sSelectorLabels := map[string]string{
		"app.kubernetes.io/component":  "testComponent2",
		"app.kubernetes.io/managed-by": "testOp2",
		"app.kubernetes.io/name":       name,
		"app.kubernetes.io/part-of":    "testProd2",
		"orbos.ch/selectable":          "yes",
	}
	selector := labels.DeriveNameSelector(nameLabels, false)
	selectable := labels.AsSelectable(nameLabels)

	serviceAccountName := "testSA2"
	replicaCount := 2
	storageCapacity := "40Gi"
	dbPort := int32(23)
	httpPort := int32(24)
	storageClass := "testSC2"
	nodeSelector := map[string]string{}
	tolerations := []corev1.Toleration{}
	resourcesSFS := &k8s.Resources{}

	quantity, err := resource.ParseQuantity(storageCapacity)
	assert.NoError(t, err)

	sfs := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    k8sSelectableLabels,
		},
		Spec: appsv1.StatefulSetSpec{
			ServiceName: name,
			Replicas:    helpers.PointerInt32(int32(replicaCount)),
			Selector: &metav1.LabelSelector{
				MatchLabels: k8sSelectorLabels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: k8sSelectableLabels,
				},
				Spec: corev1.PodSpec{
					NodeSelector:       nodeSelector,
					Tolerations:        tolerations,
					ServiceAccountName: serviceAccountName,
					Affinity:           getAffinity(k8sSelectableLabels),
					Containers: []corev1.Container{{
						Name:            name,
						Image:           image,
						ImagePullPolicy: "IfNotPresent",
						Ports: []corev1.ContainerPort{
							{ContainerPort: dbPort, Name: "grpc"},
							{ContainerPort: httpPort, Name: "http"},
						},
						LivenessProbe: &corev1.Probe{
							Handler: corev1.Handler{
								HTTPGet: &corev1.HTTPGetAction{
									Path:   "/health",
									Port:   intstr.Parse("http"),
									Scheme: "HTTPS",
								},
							},
							InitialDelaySeconds: 30,
							PeriodSeconds:       5,
						},
						ReadinessProbe: &corev1.Probe{
							Handler: corev1.Handler{
								HTTPGet: &corev1.HTTPGetAction{
									Path:   "/health?ready=1",
									Port:   intstr.Parse("http"),
									Scheme: "HTTPS",
								},
							},
							InitialDelaySeconds: 10,
							PeriodSeconds:       5,
							FailureThreshold:    2,
						},
						VolumeMounts: []corev1.VolumeMount{{
							Name:      datadirInternal,
							MountPath: datadirPath,
						}, {
							Name:      certsInternal,
							MountPath: certPath,
						}, {
							Name:      clientCertsInternal,
							MountPath: clientCertPath,
						}},
						Env: []corev1.EnvVar{{
							Name:  "COCKROACH_CHANNEL",
							Value: "kubernetes-multiregion",
						}},
						Command: []string{
							"/bin/bash",
							"-ecx",
							getJoinExec(
								namespace,
								name,
								int(dbPort),
								replicaCount,
							),
						},
						Resources: getResources(resourcesSFS),
					}},
					Volumes: []corev1.Volume{{
						Name: datadirInternal,
						VolumeSource: corev1.VolumeSource{
							PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
								ClaimName: datadirInternal,
							},
						},
					}, {
						Name: certsInternal,
						VolumeSource: corev1.VolumeSource{
							Secret: &corev1.SecretVolumeSource{
								SecretName:  nodeSecret,
								DefaultMode: helpers.PointerInt32(defaultMode),
							},
						},
					}, {
						Name: clientCertsInternal,
						VolumeSource: corev1.VolumeSource{
							Secret: &corev1.SecretVolumeSource{
								SecretName:  rootSecret,
								DefaultMode: helpers.PointerInt32(defaultMode),
							},
						},
					}},
				},
			},
			PodManagementPolicy: appsv1.PodManagementPolicyType("Parallel"),
			UpdateStrategy: appsv1.StatefulSetUpdateStrategy{
				Type: "RollingUpdate",
			},
			VolumeClaimTemplates: []corev1.PersistentVolumeClaim{{
				ObjectMeta: metav1.ObjectMeta{
					Name: datadirInternal,
				},
				Spec: corev1.PersistentVolumeClaimSpec{
					AccessModes: []corev1.PersistentVolumeAccessMode{
						corev1.PersistentVolumeAccessMode("ReadWriteOnce"),
					},
					Resources: corev1.ResourceRequirements{
						Requests: corev1.ResourceList{
							"storage": quantity,
						},
					},
					StorageClassName: &storageClass,
				},
			}},
		},
	}

	k8sClient.EXPECT().ApplyStatefulSet(sfs, false)

	query, _, _, _, _, err := AdaptFunc(
		monitor,
		selectable,
		selector,
		false,
		namespace,
		image,
		serviceAccountName,
		replicaCount,
		storageCapacity,
		dbPort,
		httpPort,
		storageClass,
		nodeSelector,
		tolerations,
		resourcesSFS,
	)
	assert.NoError(t, err)

	ensure, err := query(k8sClient)
	assert.NoError(t, err)
	assert.NotNil(t, ensure)

	assert.NoError(t, ensure(k8sClient))

}
