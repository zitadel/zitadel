package command

import (
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
)

type HostedLoginTranslationWriteModel struct {
	eventstore.WriteModel
	Language    language.Tag
	Translation map[string]any
	Level       string
	LevelID     string
}

func NewHostedLoginTranslationWriteModel(resourceID string) *HostedLoginTranslationWriteModel {
	return &HostedLoginTranslationWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   resourceID,
			ResourceOwner: resourceID,
		},
	}
}

func (wm *HostedLoginTranslationWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *org.HostedLoginTranslationSetEvent:
			wm.Language = e.Language
			wm.Translation = e.Translation
			wm.Level = e.Level
			wm.LevelID = e.Aggregate().ID
		case *instance.HostedLoginTranslationSetEvent:
			wm.Language = e.Language
			wm.Translation = e.Translation
			wm.Level = e.Level
			wm.LevelID = e.Aggregate().ID
		}
	}

	return wm.WriteModel.Reduce()
}
