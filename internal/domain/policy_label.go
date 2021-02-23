package domain

import (
	"github.com/caos/zitadel/internal/eventstore/v1/models"
)

type LabelPolicy struct {
	models.ObjectRoot

	Default        bool
	PrimaryColor   string
	SecondaryColor string
}
