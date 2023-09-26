//go:generate enumer -type Feature

package domain

type Feature int

const (
	FeatureUnspecified Feature = iota
	FeatureLoginDefaultOrg
)
