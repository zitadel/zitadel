package actions

type apiConfig struct {
	FieldConfig
}

type apiFields func(*apiConfig)

func WithAPIFields(opts ...FieldOption) apiFields {
	return func(p *apiConfig) {
		if p.fields == nil {
			p.fields = fields{}
		}

		for _, opt := range opts {
			opt(&p.FieldConfig)
		}
	}
}
