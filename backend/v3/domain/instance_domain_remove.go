package domain

import (
	"context"
	"errors"
	"strings"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type RemoveInstanceDomainCommand struct {
	InstanceID string     `json:"instance_id"`
	DomainName string     `json:"domain_name"`
	DType      DomainType `json:"domain_type"`
}

// RequiresTransaction implements [Transactional].
func (r *RemoveInstanceDomainCommand) RequiresTransaction() {}

func NewRemoveInstanceDomainCommand(instanceID, domainName string, domainType DomainType) *RemoveInstanceDomainCommand {
	return &RemoveInstanceDomainCommand{
		InstanceID: instanceID,
		DomainName: domainName,
		DType:      domainType,
	}
}

// Events implements [Commander].
func (r *RemoveInstanceDomainCommand) Events(ctx context.Context, _ *InvokeOpts) ([]eventstore.Command, error) {
	var toReturn eventstore.Command
	switch r.DType {
	case DomainTypeCustom:
		toReturn = instance.NewDomainRemovedEvent(ctx, &instance.NewAggregate(r.InstanceID).Aggregate, r.DomainName)
	case DomainTypeTrusted:
		toReturn = instance.NewTrustedDomainRemovedEvent(ctx, &instance.NewAggregate(r.InstanceID).Aggregate, r.DomainName)
	default:
		return nil, nil
	}
	return []eventstore.Command{toReturn}, nil
}

// Execute implements [Commander].
func (r *RemoveInstanceDomainCommand) Execute(ctx context.Context, opts *InvokeOpts) (err error) {
	instanceRepo := opts.instanceDomainRepo
	// TODO(IAM-Marco): I'm not sure if this is a mistake, but InstanceIDCondition must be passed, otherwise an error is returned.
	// But instance_id is not a primary key
	deletedRows, err := instanceRepo.Remove(ctx, opts.DB(), database.And(instanceRepo.InstanceIDCondition(r.InstanceID), instanceRepo.PrimaryKeyCondition(r.DomainName)))
	if err != nil {
		// TODO(IAM-Marco): Should we wrap err into zerrors.ThrowInternalError() ?
		return err
	}

	if deletedRows > 1 {
		err = zerrors.ThrowInternalf(nil, "DOM-XSCnJB", "expecting 1 row deleted, got %d", deletedRows)
		return err
	}

	if deletedRows < 1 {
		err = zerrors.ThrowNotFound(nil, "DOM-ZUteYg", "instance domain not found")
	}

	return err
}

// String implements [Commander].
func (r *RemoveInstanceDomainCommand) String() string {
	return "RemoveInstanceDomainCommand"
}

// Validate implements [Commander].
func (r *RemoveInstanceDomainCommand) Validate(ctx context.Context, opts *InvokeOpts) (err error) {
	if r.DomainName = strings.TrimSpace(r.DomainName); r.DomainName == "" || len(r.DomainName) > 253 {
		return zerrors.ThrowInvalidArgument(nil, "DOM-PLpYix", "Errors.Invalid.Argument")
	}

	if r.InstanceID = strings.TrimSpace(r.InstanceID); r.InstanceID == "" {
		return zerrors.ThrowInvalidArgument(nil, "DOM-VSsTTf", "Errors.Invalid.Argument")
	}

	if !allowedDomainCharacters.MatchString(r.DomainName) {
		return zerrors.ThrowInvalidArgument(nil, "DOM-dYD5I7", "Errors.Instance.Domain.InvalidCharacter")
	}

	// TODO(IAM-Marco): Do we want to restrict to the instance in context?
	if r.InstanceID != authz.GetInstance(ctx).InstanceID() {
		return zerrors.ThrowInvalidArgument(nil, "DOM-83FUdY", "Errors.Invalid.Argument")
	}

	if authZErr := opts.Permissions.CheckInstancePermission(ctx, DomainWritePermission); authZErr != nil {
		err = zerrors.ThrowPermissionDenied(authZErr, "DOM-eroxID", "permission denied")
		return err
	}

	domainRepo := opts.instanceDomainRepo
	d, err := domainRepo.Get(ctx, opts.DB(), database.WithCondition(domainRepo.DomainCondition(database.TextOperationEqual, r.DomainName)))
	if err != nil {
		if errors.Is(err, &database.NoRowFoundError{}) {
			return zerrors.ThrowNotFound(err, "DOM-nryNFt", "Errors.Instance.Domain.NotFound")
		}

		// TODO(IAM-Marco): Should we wrap err into zerrors.ThrowInternalError() ?
		return err
	}

	if d.Type == DomainTypeCustom && d.IsGenerated != nil && *d.IsGenerated {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-cSfCVG", "Errors.Instance.Domain.GeneratedNotRemovable")
	}

	return nil
}

var (
	_ Commander     = (*RemoveInstanceDomainCommand)(nil)
	_ Transactional = (*RemoveInstanceDomainCommand)(nil)
)
