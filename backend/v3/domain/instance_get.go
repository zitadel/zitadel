package domain

import (
	"context"
	"errors"
	"strings"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type GetInstanceQuery struct {
	ID               string `json:"id"`
	returnedInstance *Instance
}

// Result implements [Querier].
func (g *GetInstanceQuery) Result() *Instance {
	return g.returnedInstance
}

func NewGetInstanceCommand(instanceID string) *GetInstanceQuery {
	return &GetInstanceQuery{ID: instanceID}
}

// Execute implements [Executor].
func (g *GetInstanceQuery) Execute(ctx context.Context, opts *InvokeOpts) (err error) {
	instanceRepo := opts.instanceRepo.LoadDomains()

	instance, err := instanceRepo.Get(ctx, opts.DB(), database.WithCondition(instanceRepo.IDCondition(g.ID)))
	if err != nil {
		if errors.Is(err, &database.NoRowFoundError{}) {
			return zerrors.ThrowNotFound(err, "DOM-QVrUwc", "instance not found")
		}
		return zerrors.ThrowInternal(err, "DOM-lvsRce", "failed fetching instance")
	}

	g.returnedInstance = instance
	return nil
}

// String implements [Executor].
func (g *GetInstanceQuery) String() string {
	return "GetInstanceCommand"
}

// Validate implements [Validator].
func (g *GetInstanceQuery) Validate(ctx context.Context, opts *InvokeOpts) (err error) {
	g.ID = strings.TrimSpace(g.ID)
	if g.ID == "" {
		return zerrors.ThrowInvalidArgument(nil, "DOM-32a0o2", "invalid instance ID")
	}
	if g.ID != authz.GetInstance(ctx).InstanceID() {
		return zerrors.ThrowPermissionDenied(nil, "DOM-n0SvVB", "input instance ID doesn't match context instance")
	}

	if authZErr := opts.Permissions.CheckInstancePermission(ctx, InstanceReadPermission); authZErr != nil {
		return zerrors.ThrowPermissionDenied(authZErr, "DOM-Uq6b00", "permission denied")
	}

	return nil
}

var _ Querier[*Instance] = (*GetInstanceQuery)(nil)
