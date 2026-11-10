package domain

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/muhlemmer/gu"

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

	DeleteTime *time.Time `json:"delete_date"`
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

	deletedRows, err := instanceRepo.Remove(ctx, opts.DB(), database.And(instanceRepo.InstanceIDCondition(r.InstanceID), instanceRepo.PrimaryKeyCondition(r.DomainName)))
	if err != nil {
		return zerrors.ThrowInternal(err, "DOM-KH7AuJ", "Errors.Instance.Domain.Remove")
	}

	if deletedRows > 1 {
		return zerrors.ThrowInternal(nil, "DOM-XSCnJB", "Errors.Instsance.Domain.DeleteMismatch")
	}

	if deletedRows < 1 {
		return nil
	}

	r.DeleteTime = gu.Ptr(time.Now())
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

	if r.InstanceID != authz.GetInstance(ctx).InstanceID() {
		return zerrors.ThrowInvalidArgument(nil, "DOM-83FUdY", "Errors.Invalid.Argument")
	}

	if authZErr := opts.Permissions.CheckInstancePermission(ctx, DomainWritePermission); authZErr != nil {
		return zerrors.ThrowPermissionDenied(authZErr, "DOM-eroxID", "Errors.PermissionDenied")
	}

	domainRepo := opts.instanceDomainRepo
	d, err := domainRepo.Get(ctx, opts.DB(), database.WithCondition(domainRepo.DomainCondition(database.TextOperationEqual, r.DomainName)))
	if err != nil {
		if errors.Is(err, &database.NoRowFoundError{}) {
			return nil
		}
		return zerrors.ThrowInternal(err, "DOM-Zvv1fi", "Errors.Instance.Domain.Get")
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
