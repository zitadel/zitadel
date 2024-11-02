package domain

type Feature int

func (f Feature) Type() FeatureType {
	switch f {
	case FeatureUnspecified:
		return FeatureTypeUnspecified
	case FeatureLoginDefaultOrg:
		return FeatureTypeBoolean
	default:
		return FeatureTypeUnspecified
	}
}

const (
	FeatureTypeUnspecified FeatureType = iota
	FeatureTypeBoolean
)

type FeatureType int

const (
	FeatureUnspecified Feature = iota
	FeatureLoginDefaultOrg
)
