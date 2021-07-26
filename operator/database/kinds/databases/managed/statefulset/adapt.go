package statefulset

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/kubernetes/k8s"
	"github.com/caos/orbos/pkg/kubernetes/resources"
	"github.com/caos/orbos/pkg/kubernetes/resources/statefulset"
	"github.com/caos/orbos/pkg/labels"

	"github.com/caos/zitadel/operator"
	"github.com/caos/zitadel/operator/helpers"
)

const (
	certPath            = "/cockroach/cockroach-certs"
	clientCertPath      = "/cockroach/cockroach-client-certs"
	datadirPath         = "/cockroach/cockroach-data"
	datadirInternal     = "datadir"
	certsInternal       = "certs"
	clientCertsInternal = "client-certs"
	defaultMode         = int32(256)
	nodeSecret          = "cockroachdb.node"
	rootSecret          = "cockroachdb.client.root"
)

type Affinity struct {
	key   string
	value string
}

type Affinitys []metav1.LabelSelectorRequirement

func (a Affinitys) Len() int           { return len(a) }
func (a Affinitys) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Affinitys) Less(i, j int) bool { return a[i].Key < a[j].Key }

func AdaptFunc(
	monitor mntr.Monitor,
	sfsSelectable *labels.Selectable,
	podSelector *labels.Selector,
	force bool,
	namespace string,
	image string,
	serviceAccountName string,
	replicaCount int,
	storageCapacity resource.Quantity,
	dbPort int32,
	httpPort int32,
	storageClass string,
	nodeSelector map[string]string,
	tolerations []corev1.Toleration,
	resourcesSFS *k8s.Resources,
) (
	resources.QueryFunc,
	resources.DestroyFunc,
	operator.EnsureFunc,
	operator.EnsureFunc,
	func(k8sClient kubernetes.ClientInt) ([]string, error),
	error,
) {
	internalMonitor := monitor.WithField("component", "statefulset")

	name := sfsSelectable.Name()
	k8sSelectable := labels.MustK8sMap(sfsSelectable)
	statefulsetDef := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    k8sSelectable,
		},
		Spec: appsv1.StatefulSetSpec{
			ServiceName: name,
			Replicas:    helpers.PointerInt32(int32(replicaCount)),
			Selector: &metav1.LabelSelector{
				MatchLabels: labels.MustK8sMap(podSelector),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: k8sSelectable,
				},
				Spec: corev1.PodSpec{
					NodeSelector:       nodeSelector,
					Tolerations:        tolerations,
					ServiceAccountName: serviceAccountName,
					Affinity:           getAffinity(k8sSelectable),
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
			PodManagementPolicy: appsv1.ParallelPodManagement,
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
							"storage": storageCapacity,
						},
					},
					StorageClassName: &storageClass,
				},
			}},
		},
	}

	query, err := statefulset.AdaptFuncToEnsure(statefulsetDef, force)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	destroy, err := statefulset.AdaptFuncToDestroy(namespace, name)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	wrapedQuery, wrapedDestroy, err := resources.WrapFuncs(internalMonitor, query, destroy)
	checkDBRunning := func(k8sClient kubernetes.ClientInt) error {
		internalMonitor.Info("waiting for statefulset to be running")
		if err := k8sClient.WaitUntilStatefulsetIsReady(namespace, name, true, false, 60*time.Second); err != nil {
			return fmt.Errorf("error while waiting for statefulset to be running: %w", err)
		}
		internalMonitor.Info("statefulset is running")
		return nil
	}

	checkDBNotReady := func(k8sClient kubernetes.ClientInt) error {
		internalMonitor.Info("checking for statefulset to not be ready")
		if err := k8sClient.WaitUntilStatefulsetIsReady(namespace, name, true, true, 1*time.Second); err != nil {
			internalMonitor.Info("statefulset is not ready")
			return nil
		}
		return errors.New("statefulset is ready")
	}

	ensureInit := func(k8sClient kubernetes.ClientInt) error {
		if err := checkDBRunning(k8sClient); err != nil {
			return err
		}

		if err := checkDBNotReady(k8sClient); err != nil {
			return nil
		}

		command := "/cockroach/cockroach init --certs-dir=" + clientCertPath + " --host=" + name + "-0." + name

		if err := k8sClient.ExecInPod(namespace, name+"-0", name, command); err != nil {
			return err
		}
		return nil
	}

	checkDBReady := func(k8sClient kubernetes.ClientInt) error {
		internalMonitor.Info("waiting for statefulset to be ready")
		if err := k8sClient.WaitUntilStatefulsetIsReady(namespace, name, true, true, 60*time.Second); err != nil {
			return fmt.Errorf("error while waiting for statefulset to be ready: %w", err)
		}
		internalMonitor.Info("statefulset is ready")
		return nil
	}

	getAllDBs := func(k8sClient kubernetes.ClientInt) ([]string, error) {
		if err := checkDBRunning(k8sClient); err != nil {
			return nil, err
		}

		if err := checkDBReady(k8sClient); err != nil {
			return nil, err
		}

		command := "/cockroach/cockroach sql --certs-dir=" + clientCertPath + " --host=" + name + "-0." + name + " -e 'SHOW DATABASES;'"

		databasesStr, err := k8sClient.ExecInPodWithOutput(namespace, name+"-0", name, command)
		if err != nil {
			return nil, err
		}
		databases := strings.Split(databasesStr, "\n")
		dbAndOwners := databases[1 : len(databases)-1]
		dbs := []string{}
		for _, dbAndOwner := range dbAndOwners {
			parts := strings.Split(dbAndOwner, "\t")
			if parts[1] != "node" {
				dbs = append(dbs, parts[0])
			}
		}
		return dbs, nil
	}

	return wrapedQuery, wrapedDestroy, ensureInit, checkDBReady, getAllDBs, err
}

func getJoinExec(namespace string, name string, dbPort int, replicaCount int) string {
	joinList := make([]string, 0)
	for i := 0; i < replicaCount; i++ {
		joinList = append(joinList, fmt.Sprintf("%s-%d.%s.%s:%d", name, i, name, namespace, dbPort))
	}
	joinListStr := strings.Join(joinList, ",")
	locality := "zone=" + namespace

	return "exec /cockroach/cockroach start --logtostderr --certs-dir " + certPath + " --advertise-host $(hostname -f) --http-addr 0.0.0.0 --join " + joinListStr + " --locality " + locality + " --cache 25% --max-sql-memory 25%"
}

func getResources(resourcesSFS *k8s.Resources) corev1.ResourceRequirements {
	internalResources := corev1.ResourceRequirements{
		Requests: corev1.ResourceList{
			"cpu":    resource.MustParse("1"),
			"memory": resource.MustParse("6Gi"),
		},
		Limits: corev1.ResourceList{
			"cpu":    resource.MustParse("4"),
			"memory": resource.MustParse("8Gi"),
		},
	}

	if resourcesSFS != nil {
		internalResources = corev1.ResourceRequirements{}
		if resourcesSFS.Requests != nil {
			internalResources.Requests = resourcesSFS.Requests
		}
		if resourcesSFS.Limits != nil {
			internalResources.Limits = resourcesSFS.Limits
		}
	}

	return internalResources
}

func getAffinity(labels map[string]string) *corev1.Affinity {
	affinity := Affinitys{}
	for k, v := range labels {
		affinity = append(affinity, metav1.LabelSelectorRequirement{
			Key:      k,
			Operator: metav1.LabelSelectorOpIn,
			Values: []string{
				v,
			}})
	}
	sort.Sort(affinity)

	return &corev1.Affinity{
		PodAntiAffinity: &corev1.PodAntiAffinity{
			RequiredDuringSchedulingIgnoredDuringExecution: []corev1.PodAffinityTerm{{
				LabelSelector: &metav1.LabelSelector{
					MatchExpressions: affinity,
				},
				TopologyKey: "kubernetes.io/hostname",
			}},
		},
	}
}
