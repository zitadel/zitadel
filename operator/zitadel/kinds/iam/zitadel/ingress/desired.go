package ingress

// TODO: Validate
type Spec struct {
	Controller          string                 `yaml:"controller"`
	ControllerSpecifics map[string]interface{} `yaml:"controllerSpecifics,omitempty"`
}
