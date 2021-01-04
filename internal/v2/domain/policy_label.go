package domain

import (
	"github.com/caos/zitadel/internal/eventstore/models"
)

type LabelPolicy struct {
	models.ObjectRoot

	Default        bool
	PrimaryColor   string
	SecondaryColor string
}
