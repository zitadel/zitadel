package feature

type Defaults struct {
	FeatureA bool
}

func (d *Defaults) ToInstance() *Instance {
	return &Instance{
		FeatureA: d.FeatureA,
	}
}

type System struct {
}

type Instance struct {
	FeatureA bool `json:"feature_a,omitempty"`
}

type Organization struct {
}

type User struct {
}
