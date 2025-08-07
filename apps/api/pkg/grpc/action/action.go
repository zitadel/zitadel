package action

import "github.com/zitadel/zitadel/internal/api/grpc/server/middleware"

func (f *Flow) Localizers() []middleware.Localizer {
	if f == nil {
		return nil
	}

	localizers := make([]middleware.Localizer, 0, len(f.TriggerActions)+1)
	localizers = append(localizers, f.Type.Localizers()...)
	for _, action := range f.TriggerActions {
		localizers = append(localizers, action.Localizers()...)
	}

	return localizers
}

func (t *FlowType) Localizers() []middleware.Localizer {
	if t == nil {
		return nil
	}

	return []middleware.Localizer{t.Name}
}

func (t *TriggerType) Localizers() []middleware.Localizer {
	if t == nil {
		return nil
	}

	return []middleware.Localizer{t.Name}
}

func (ta *TriggerAction) Localizers() []middleware.Localizer {
	if ta == nil {
		return nil
	}

	return ta.TriggerType.Localizers()
}
