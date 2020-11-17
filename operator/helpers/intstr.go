package helpers

import "k8s.io/apimachinery/pkg/util/intstr"

func IntToIntStr(value int) *intstr.IntOrString {
	v := intstr.FromInt(value)
	return &v
}

func StringToIntStr(value string) *intstr.IntOrString {
	v := intstr.FromString(value)
	return &v
}
