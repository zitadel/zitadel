package domain

import (
	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/eventstore/v1/models"
)

type CustomText struct {
	models.ObjectRoot

	State    CustomTextState
	Default  bool
	Key      string
	Language language.Tag
	Text     string
}

type CustomTextState int32

const (
	CustomTextStateUnspecified CustomTextState = iota
	CustomTextStateActive
	CustomTextStateRemoved

	customTextStateCount
)

func (m *CustomText) IsValid() bool {
	return m.Key != "" && m.Language != language.Und && m.Text != ""
}
