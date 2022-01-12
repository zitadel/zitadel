package repository

import (
	"context"
	"time"

	"golang.org/x/text/language"

	org_model "github.com/caos/zitadel/internal/org/model"
)

type OrgRepository interface {
	Languages(ctx context.Context) ([]language.Tag, error)
	OrgChanges(ctx context.Context, id string, lastSequence uint64, limit uint64, sortAscending bool, auditLogRetention time.Duration) (*org_model.OrgChanges, error)

	GetOrgMemberRoles(isGlobal bool) []string
}
