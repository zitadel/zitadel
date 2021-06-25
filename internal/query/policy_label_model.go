package query

import (
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/policy"
)

type LabelPolicyReadModel struct {
	eventstore.ReadModel

	PrimaryColor        string
	BackgroundColor     string
	WarnColor           string
	FontColor           string
	PrimaryColorDark    string
	BackgroundColorDark string
	WarnColorDark       string
	FontColorDark       string
	IsActive            bool
}

func (rm *LabelPolicyReadModel) Reduce() error {
	for _, event := range rm.Events {
		switch e := event.(type) {
		case *policy.LabelPolicyAddedEvent:
			rm.PrimaryColor = e.PrimaryColor
			rm.BackgroundColor = e.BackgroundColor
			rm.FontColor = e.FontColor
			rm.WarnColor = e.WarnColor
			rm.PrimaryColorDark = e.PrimaryColorDark
			rm.BackgroundColorDark = e.BackgroundColorDark
			rm.FontColorDark = e.FontColorDark
			rm.WarnColorDark = e.WarnColorDark
			rm.IsActive = true
		case *policy.LabelPolicyChangedEvent:
			if e.PrimaryColor != nil {
				rm.PrimaryColor = *e.PrimaryColor
			}
			if e.BackgroundColor != nil {
				rm.BackgroundColor = *e.BackgroundColor
			}
			if e.WarnColor != nil {
				rm.WarnColor = *e.WarnColor
			}
			if e.FontColor != nil {
				rm.FontColor = *e.FontColor
			}
			if e.PrimaryColorDark != nil {
				rm.PrimaryColorDark = *e.PrimaryColorDark
			}
			if e.BackgroundColorDark != nil {
				rm.BackgroundColorDark = *e.BackgroundColorDark
			}
			if e.WarnColorDark != nil {
				rm.WarnColorDark = *e.WarnColorDark
			}
			if e.FontColorDark != nil {
				rm.FontColorDark = *e.FontColorDark
			}
		case *policy.LabelPolicyRemovedEvent:
			rm.IsActive = false
		}
	}
	return rm.ReadModel.Reduce()
}
