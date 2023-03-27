package actions

type ctxConfig struct {
	FieldConfig
}

type contextFields func(*ctxConfig)

func SetContextFields(opts ...FieldOption) contextFields {
	return func(p *ctxConfig) {
		if p.fields == nil {
			p.fields = fields{}
		}

		for _, opt := range opts {
			opt(&p.FieldConfig)
		}
	}
}
