// +kubebuilder:object:generate=true
// +groupName=caos.ch
package v1

import (
	"github.com/caos/orbos/pkg/tree"
	orbz "github.com/caos/zitadel/operator/zitadel/kinds/orb"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

var (
	// GroupVersion is group version used to register these objects
	GroupVersion = schema.GroupVersion{Group: "caos.ch", Version: "v1"}

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuilder = &scheme.Builder{GroupVersion: GroupVersion}

	// AddToScheme adds the types in this group-version to the given scheme.
	AddToScheme = SchemeBuilder.AddToScheme
)

// +kubebuilder:storageversion
// +kubebuilder:object:root=true
// +kubebuilder:crd=Zitadel
type Zitadel struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   Spec   `json:"spec,omitempty"`
	Status Status `json:"status,omitempty"`
}

type Status struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

type Spec struct {
	Common *tree.Common `json:",inline" yaml:",inline"`
	Spec   *orbz.Spec   `json:"spec" yaml:"spec"`
	IAM    *Empty       `json:"iam" yaml:"iam"`
}

type Empty struct{}

// +kubebuilder:object:root=true
type ZitadelList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Zitadel `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Zitadel{}, &ZitadelList{})
}
