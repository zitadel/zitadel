package command

import "github.com/zitadel/zitadel/internal/eventstore"

type HumanProfileWriteModels struct {
	eventstore.WriteModel
	Profiles map[string]*HumanProfileWriteModel
}

func NewHumanProfileWriteModels() *HumanProfileWriteModels {
	return &HumanProfileWriteModels{
		Profiles: make(map[string]*HumanProfileWriteModel),
	}
}

func (wm *HumanProfileWriteModels) Reduce() error {
	for _, event := range wm.Events {
		agg := event.Aggregate()
		profile, ok := wm.Profiles[agg.ID]
		if !ok {
			profile = NewHumanProfileWriteModel(agg.ID, agg.ResourceOwner)
			wm.Profiles[agg.ID] = profile
		}
		profile.AppendEvents(event)
	}
	for _, profile := range wm.Profiles {
		if err := profile.Reduce(); err != nil {
			return err
		}
	}
	return nil
}

func (wm *HumanProfileWriteModels) Query() *eventstore.SearchQueryBuilder {
	return humanProfileWriteModelQuery("", "")
}
