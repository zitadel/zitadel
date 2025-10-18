package domain

import (
	"context"
	"strings"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
	instance "github.com/zitadel/zitadel/pkg/grpc/instance/v2beta"
)

type GetInstanceCommand struct {
	ID               string `json:"id"`
	ReturnedInstance *Instance
}

func NewGetInstanceCommand(instanceID string) *GetInstanceCommand {
	return &GetInstanceCommand{ID: instanceID}
}

// Events implements Commander.
func (g *GetInstanceCommand) Events(ctx context.Context, opts *CommandOpts) ([]eventstore.Command, error) {
	return nil, nil
}

// Execute implements Commander.
func (g *GetInstanceCommand) Execute(ctx context.Context, opts *CommandOpts) (err error) {
	closeFunc, err := opts.EnsureTx(ctx)
	if err != nil {
		return err
	}

	defer func() { err = closeFunc(ctx, err) }()

	instanceRepo := opts.instanceRepo.LoadDomains()

	instance, err := instanceRepo.Get(ctx, pool, database.WithCondition(instanceRepo.IDCondition(g.ID)))
	if err != nil {
		return err
	}

	g.ReturnedInstance = instance
	return nil
}

// String implements Commander.
func (g *GetInstanceCommand) String() string {
	return "GetInstanceCommand"
}

// Validate implements Commander.
func (g *GetInstanceCommand) Validate(ctx context.Context, opts *CommandOpts) (err error) {
	instanceID := strings.TrimSpace(g.ID)
	if instanceID == "" {
		return zerrors.ThrowInvalidArgument(nil, "DOM-32a0o2", "invalid instance ID")
	}
	if instanceID != authz.GetInstance(ctx).InstanceID() {
		return zerrors.ThrowPermissionDenied(nil, "DOM-n0SvVB", "input instance ID doesn't match context instance")
	}

	if authZErr := opts.Permissions.CheckInstancePermission(ctx, InstanceReadPermission); authZErr != nil {
		return zerrors.ThrowPermissionDenied(authZErr, "DOM-Uq6b00", "permission denied")
	}

	return nil
}

func (g *GetInstanceCommand) ResultToGRPC() *instance.Instance {
	inst := g.ReturnedInstance
	return &instance.Instance{
		Id:           inst.ID,
		ChangeDate:   timestamppb.New(inst.UpdatedAt),
		CreationDate: timestamppb.New(inst.CreatedAt),
		State:        instance.State_STATE_RUNNING, // TODO(IAM-Marco): Not sure what to put here
		Name:         inst.Name,
		Version:      "",
		Domains:      g.domainsToGRPC(),
	}
}

func (g *GetInstanceCommand) domainsToGRPC() []*instance.Domain {
	toReturn := make([]*instance.Domain, len(g.ReturnedInstance.Domains))
	for i, domain := range g.ReturnedInstance.Domains {
		isGenerated := domain.IsGenerated != nil && *domain.IsGenerated
		isPrimary := domain.IsPrimary != nil && *domain.IsPrimary
		toReturn[i] = &instance.Domain{
			InstanceId:   domain.InstanceID,
			CreationDate: timestamppb.New(domain.CreatedAt),
			Domain:       domain.Domain,
			Primary:      isPrimary,
			Generated:    isGenerated,
		}
	}

	return toReturn
}

var _ Commander = (*GetInstanceCommand)(nil)
