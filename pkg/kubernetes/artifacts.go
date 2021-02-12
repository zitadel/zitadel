package kubernetes

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/labels"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	rbac "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	mach "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func EnsureZitadelOperatorArtifacts(
	monitor mntr.Monitor,
	apiLabels *labels.API,
	client kubernetes.ClientInt,
	version string,
	nodeselector map[string]string,
	tolerations []core.Toleration,
	imageRegistry string,
	gitops bool,
) error {

	monitor.WithFields(map[string]interface{}{
		"zitadel": version,
	}).Debug("Ensuring zitadel artifacts")

	nameLabels := labels.MustForName(labels.MustForComponent(apiLabels, "operator"), "zitadel-operator")
	k8sNameLabels := labels.MustK8sMap(nameLabels)
	k8sPodSelector := labels.MustK8sMap(labels.DeriveNameSelector(nameLabels, false))

	if version == "" {
		return nil
	}

	if err := client.ApplyServiceAccount(&core.ServiceAccount{
		ObjectMeta: mach.ObjectMeta{
			Name:      nameLabels.Name(),
			Namespace: "caos-system",
		},
	}); err != nil {
		return err
	}

	if err := client.ApplyClusterRole(&rbac.ClusterRole{
		ObjectMeta: mach.ObjectMeta{
			Name:   nameLabels.Name(),
			Labels: k8sNameLabels,
		},
		Rules: []rbac.PolicyRule{{
			APIGroups: []string{"*"},
			Resources: []string{"*"},
			Verbs:     []string{"*"},
		}},
	}); err != nil {
		return err
	}

	if err := client.ApplyClusterRoleBinding(&rbac.ClusterRoleBinding{
		ObjectMeta: mach.ObjectMeta{
			Name:   nameLabels.Name(),
			Labels: k8sNameLabels,
		},

		RoleRef: rbac.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     nameLabels.Name(),
		},
		Subjects: []rbac.Subject{{
			Kind:      "ServiceAccount",
			Name:      nameLabels.Name(),
			Namespace: "caos-system",
		}},
	}); err != nil {
		return err
	}

	if !gitops {
		crd := `apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  annotations:
    controller-gen.kubebuilder.io/version: v0.2.2
  creationTimestamp: null
  name: zitadels.caos.ch
spec:
  group: caos.ch
  names:
    kind: Zitadel
    listKind: ZitadelList
    plural: zitadels
    singular: zitadel
  scope: ""
  validation:
    openAPIV3Schema:
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          properties:
            database:
              type: object
            kind:
              type: string
            spec:
              properties:
                customImageRegistry:
                  description: 'Use this registry to pull the zitadel operator image
                    from @default: ghcr.io'
                  type: string
                databaseCrd:
                  properties:
                    name:
                      type: string
                    namespace:
                      type: string
                  type: object
                nodeSelector:
                  additionalProperties:
                    type: string
                  type: object
                selfReconciling:
                  type: boolean
                tolerations:
                  items:
                    description: The pod this Toleration is attached to tolerates
                      any taint that matches the triple <key,value,effect> using the
                      matching operator <operator>.
                    properties:
                      effect:
                        description: Effect indicates the taint effect to match. Empty
                          means match all taint effects. When specified, allowed values
                          are NoSchedule, PreferNoSchedule and NoExecute.
                        type: string
                      key:
                        description: Key is the taint key that the toleration applies
                          to. Empty means match all taint keys. If the key is empty,
                          operator must be Exists; this combination means to match
                          all values and all keys.
                        type: string
                      operator:
                        description: Operator represents a key's relationship to the
                          value. Valid operators are Exists and Equal. Defaults to
                          Equal. Exists is equivalent to wildcard for value, so that
                          a pod can tolerate all taints of a particular category.
                        type: string
                      tolerationSeconds:
                        description: TolerationSeconds represents the period of time
                          the toleration (which must be of effect NoExecute, otherwise
                          this field is ignored) tolerates the taint. By default,
                          it is not set, which means tolerate the taint forever (do
                          not evict). Zero and negative values will be treated as
                          0 (evict immediately) by the system.
                        format: int64
                        type: integer
                      value:
                        description: Value is the taint value the toleration matches
                          to. If the operator is Exists, the value should be empty,
                          otherwise just a regular string.
                        type: string
                    type: object
                  type: array
                verbose:
                  type: boolean
                version:
                  type: string
              required:
              - selfReconciling
              - verbose
              type: object
            version:
              type: string
          required:
          - database
          - kind
          - spec
          - version
          type: object
        status:
          type: object
      type: object
  version: v1
  versions:
  - name: v1
    served: true
    storage: true
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []`

		crdDefinition := &unstructured.Unstructured{}
		if err := yaml.Unmarshal([]byte(crd), &crdDefinition.Object); err != nil {
			return err
		}

		if err := client.ApplyCRDResource(
			crdDefinition,
		); err != nil {
			return err
		}
		monitor.WithFields(map[string]interface{}{
			"version": version,
		}).Debug("Database Operator crd ensured")
	}

	cmd := []string{"/zitadelctl", "operator", "-f", "/secrets/orbconfig"}
	if gitops {
		cmd = append(cmd, "--gitops")
	}

	deployment := &apps.Deployment{
		ObjectMeta: mach.ObjectMeta{
			Name:      nameLabels.Name(),
			Namespace: "caos-system",
			Labels:    k8sNameLabels,
		},
		Spec: apps.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &mach.LabelSelector{
				MatchLabels: k8sPodSelector,
			},
			Template: core.PodTemplateSpec{
				ObjectMeta: mach.ObjectMeta{
					Labels: labels.MustK8sMap(labels.AsSelectable(nameLabels)),
				},
				Spec: core.PodSpec{
					ServiceAccountName: nameLabels.Name(),
					Containers: []core.Container{{
						Name:            "zitadel",
						ImagePullPolicy: core.PullIfNotPresent,
						Image:           fmt.Sprintf("%s/caos/zitadel-operator:%s", imageRegistry, version),
						Command:         cmd,
						Args:            []string{},
						Ports: []core.ContainerPort{{
							Name:          "metrics",
							ContainerPort: 2112,
							Protocol:      "TCP",
						}},
						VolumeMounts: []core.VolumeMount{{
							Name:      "orbconfig",
							ReadOnly:  true,
							MountPath: "/secrets",
						}},
						Resources: core.ResourceRequirements{
							Limits: core.ResourceList{
								"cpu":    resource.MustParse("500m"),
								"memory": resource.MustParse("500Mi"),
							},
							Requests: core.ResourceList{
								"cpu":    resource.MustParse("250m"),
								"memory": resource.MustParse("250Mi"),
							},
						},
					}},
					NodeSelector: nodeselector,
					Tolerations:  tolerations,
					Volumes: []core.Volume{{
						Name: "orbconfig",
						VolumeSource: core.VolumeSource{
							Secret: &core.SecretVolumeSource{
								SecretName: "caos",
							},
						},
					}},
					TerminationGracePeriodSeconds: int64Ptr(10),
				},
			},
		},
	}
	if err := client.ApplyDeployment(deployment, true); err != nil {
		return err
	}
	monitor.WithFields(map[string]interface{}{
		"version": version,
	}).Debug("Zitadel Operator deployment ensured")

	return nil
}

func ScaleZitadelOperator(
	monitor mntr.Monitor,
	client *kubernetes.Client,
	replicaCount int,
) error {
	monitor.Debug("Scaling zitadel-operator")
	return client.ScaleDeployment("caos-system", "zitadel-operator", replicaCount)
}

func int32Ptr(i int32) *int32 { return &i }
func int64Ptr(i int64) *int64 { return &i }

func toNameLabels(apiLabels *labels.API, operatorName string) *labels.Name {
	return labels.MustForName(labels.MustForComponent(apiLabels, "operator"), operatorName)
}

func EnsureDatabaseArtifacts(
	monitor mntr.Monitor,
	apiLabels *labels.API,
	client kubernetes.ClientInt,
	version string,
	nodeselector map[string]string,
	tolerations []core.Toleration,
	imageRegistry string,
	gitops bool) error {

	monitor.WithFields(map[string]interface{}{
		"database": version,
	}).Debug("Ensuring database artifacts")

	if version == "" {
		return nil
	}

	nameLabels := toNameLabels(apiLabels, "database-operator")
	k8sNameLabels := labels.MustK8sMap(nameLabels)

	if err := client.ApplyServiceAccount(&core.ServiceAccount{
		ObjectMeta: mach.ObjectMeta{
			Name:      "database-operator",
			Namespace: "caos-system",
			Labels:    k8sNameLabels,
		},
	}); err != nil {
		return err
	}

	if err := client.ApplyClusterRole(&rbac.ClusterRole{
		ObjectMeta: mach.ObjectMeta{
			Name:   "database-operator-clusterrole",
			Labels: k8sNameLabels,
		},
		Rules: []rbac.PolicyRule{{
			APIGroups: []string{"*"},
			Resources: []string{"*"},
			Verbs:     []string{"*"},
		}},
	}); err != nil {
		return err
	}

	if err := client.ApplyClusterRoleBinding(&rbac.ClusterRoleBinding{
		ObjectMeta: mach.ObjectMeta{
			Name:   "database-operator-clusterrolebinding",
			Labels: k8sNameLabels,
		},

		RoleRef: rbac.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     "database-operator-clusterrole",
		},
		Subjects: []rbac.Subject{{
			Kind:      "ServiceAccount",
			Name:      "database-operator",
			Namespace: "caos-system",
		}},
	}); err != nil {
		return err
	}

	if !gitops {
		crd := `apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  annotations:
    controller-gen.kubebuilder.io/version: v0.2.2
  creationTimestamp: null
  name: databases.caos.ch
spec:
  group: caos.ch
  names:
    kind: Database
    listKind: DatabaseList
    plural: databases
    singular: database
  scope: ""
  validation:
    openAPIV3Schema:
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          properties:
            database:
              type: object
            kind:
              type: string
            spec:
              properties:
                customImageRegistry:
                  description: 'Use this registry to pull the Database operator image
                    from @default: ghcr.io'
                  type: string
                nodeSelector:
                  additionalProperties:
                    type: string
                  type: object
                selfReconciling:
                  type: boolean
                tolerations:
                  items:
                    description: The pod this Toleration is attached to tolerates
                      any taint that matches the triple <key,value,effect> using the
                      matching operator <operator>.
                    properties:
                      effect:
                        description: Effect indicates the taint effect to match. Empty
                          means match all taint effects. When specified, allowed values
                          are NoSchedule, PreferNoSchedule and NoExecute.
                        type: string
                      key:
                        description: Key is the taint key that the toleration applies
                          to. Empty means match all taint keys. If the key is empty,
                          operator must be Exists; this combination means to match
                          all values and all keys.
                        type: string
                      operator:
                        description: Operator represents a key's relationship to the
                          value. Valid operators are Exists and Equal. Defaults to
                          Equal. Exists is equivalent to wildcard for value, so that
                          a pod can tolerate all taints of a particular category.
                        type: string
                      tolerationSeconds:
                        description: TolerationSeconds represents the period of time
                          the toleration (which must be of effect NoExecute, otherwise
                          this field is ignored) tolerates the taint. By default,
                          it is not set, which means tolerate the taint forever (do
                          not evict). Zero and negative values will be treated as
                          0 (evict immediately) by the system.
                        format: int64
                        type: integer
                      value:
                        description: Value is the taint value the toleration matches
                          to. If the operator is Exists, the value should be empty,
                          otherwise just a regular string.
                        type: string
                    type: object
                  type: array
                verbose:
                  type: boolean
                version:
                  type: string
              required:
              - selfReconciling
              - verbose
              type: object
            version:
              type: string
          required:
          - database
          - kind
          - spec
          - version
          type: object
        status:
          type: object
      type: object
  version: v1
  versions:
  - name: v1
    served: true
    storage: true
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []`

		crdDefinition := &unstructured.Unstructured{}
		if err := yaml.Unmarshal([]byte(crd), &crdDefinition.Object); err != nil {
			return err
		}

		if err := client.ApplyCRDResource(
			crdDefinition,
		); err != nil {
			return err
		}
		monitor.WithFields(map[string]interface{}{
			"version": version,
		}).Debug("Database Operator crd ensured")
	}

	cmd := []string{"/zitadelctl", "database", "-f", "/secrets/orbconfig"}
	if gitops {
		cmd = append(cmd, "--gitops")
	}

	deployment := &apps.Deployment{
		ObjectMeta: mach.ObjectMeta{
			Name:      "database-operator",
			Namespace: "caos-system",
			Labels:    k8sNameLabels,
		},
		Spec: apps.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &mach.LabelSelector{
				MatchLabels: labels.MustK8sMap(labels.DeriveNameSelector(nameLabels, false)),
			},
			Template: core.PodTemplateSpec{
				ObjectMeta: mach.ObjectMeta{
					Labels: labels.MustK8sMap(labels.AsSelectable(nameLabels)),
				},
				Spec: core.PodSpec{
					ServiceAccountName: "database-operator",
					Containers: []core.Container{{
						Name:            "database",
						ImagePullPolicy: core.PullIfNotPresent,
						Image:           fmt.Sprintf("%s/caos/zitadel-operator:%s", imageRegistry, version),
						Command:         cmd,
						Args:            []string{},
						Ports: []core.ContainerPort{{
							Name:          "metrics",
							ContainerPort: 2112,
							Protocol:      "TCP",
						}},
						VolumeMounts: []core.VolumeMount{{
							Name:      "orbconfig",
							ReadOnly:  true,
							MountPath: "/secrets",
						}},
						Resources: core.ResourceRequirements{
							Limits: core.ResourceList{
								"cpu":    resource.MustParse("500m"),
								"memory": resource.MustParse("500Mi"),
							},
							Requests: core.ResourceList{
								"cpu":    resource.MustParse("250m"),
								"memory": resource.MustParse("250Mi"),
							},
						},
					}},
					NodeSelector: nodeselector,
					Tolerations:  tolerations,
					Volumes: []core.Volume{{
						Name: "orbconfig",
						VolumeSource: core.VolumeSource{
							Secret: &core.SecretVolumeSource{
								SecretName: "caos",
							},
						},
					}},
					TerminationGracePeriodSeconds: int64Ptr(10),
				},
			},
		},
	}

	if err := client.ApplyDeployment(deployment, true); err != nil {
		return err
	}
	monitor.WithFields(map[string]interface{}{
		"version": version,
	}).Debug("Database Operator deployment ensured")

	return nil
}
