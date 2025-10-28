package domain

import (
	"context"
	"errors"
	"regexp"
	"strings"

	"github.com/muhlemmer/gu"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var allowedDomainCharacters = regexp.MustCompile(`^[a-zA-Z0-9\\.\\-]+$`)

type AddInstanceDomainCommand struct {
	InstanceID string `json:"instance_id"`
	DomainName string `json:"domain_name"`
}

// RequiresTransaction implements [Transactional].
func (a *AddInstanceDomainCommand) RequiresTransaction() {}

func NewAddInstanceDomainCommand(instanceID, domainName string) *AddInstanceDomainCommand {
	return &AddInstanceDomainCommand{
		InstanceID: instanceID,
		DomainName: domainName,
	}
}

// Events implements [Commander].
func (a *AddInstanceDomainCommand) Events(ctx context.Context, _ *InvokeOpts) ([]eventstore.Command, error) {
	return []eventstore.Command{
		instance.NewDomainAddedEvent(ctx, &instance.NewAggregate(a.InstanceID).Aggregate, a.DomainName, false),
	}, nil
}

// Execute implements [Commander].
func (a *AddInstanceDomainCommand) Execute(ctx context.Context, opts *InvokeOpts) (err error) {
	instanceRepo := opts.instanceDomainRepo
	err = instanceRepo.Add(ctx, opts.DB(), &AddInstanceDomain{
		InstanceID:  a.InstanceID,
		Domain:      a.DomainName,
		IsPrimary:   gu.Ptr(false),
		IsGenerated: gu.Ptr(false),
		Type:        DomainTypeCustom,
	})
	if err != nil {
		// TODO(IAM-Marco): Should we wrap err into zerrors.ThrowInternalError() ?
		return err
	}

	return nil
}

// String implements [Commander].
func (a *AddInstanceDomainCommand) String() string {
	return "AddInstanceDomainCommand"
}

// Validate implements [Commander].
func (a *AddInstanceDomainCommand) Validate(ctx context.Context, opts *InvokeOpts) (err error) {
	if a.DomainName = strings.TrimSpace(a.DomainName); a.DomainName == "" {
		return zerrors.ThrowInvalidArgument(nil, "DOM-jieuM8", "Errors.Invalid.Argument")
	}

	if a.InstanceID = strings.TrimSpace(a.InstanceID); a.InstanceID == "" {
		return zerrors.ThrowInvalidArgument(nil, "DOM-YaUBp5", "Errors.Invalid.Argument")
	}

	if !allowedDomainCharacters.MatchString(a.DomainName) {
		return zerrors.ThrowInvalidArgument(nil, "DOM-98VcSQ", "Errors.Instance.Domain.InvalidCharacter")
	}

	// TODO(IAM-Marco): Do we want to restrict to the instance in context?
	if a.InstanceID != authz.GetInstance(ctx).InstanceID() {
		return zerrors.ThrowInvalidArgument(nil, "DOM-x01cai", "Errors.Invalid.Argument")
	}

	if authZErr := opts.Permissions.CheckInstancePermission(ctx, DomainWritePermission); authZErr != nil {
		err = zerrors.ThrowPermissionDenied(authZErr, "DOM-c83vPX", "permission denied")
		return err
	}

	domainRepo := opts.instanceDomainRepo
	_, err = domainRepo.Get(ctx, opts.DB(), database.WithCondition(domainRepo.DomainCondition(database.TextOperationEqual, a.DomainName)))
	if err == nil {
		return zerrors.ThrowAlreadyExists(nil, "DOM-CvQ8tf", "Errors.Instance.Domain.AlreadyExists")
	}

	if !errors.Is(err, &database.NoRowFoundError{}) {
		return err
	}

	return nil
}

var (
	_ Commander     = (*AddInstanceDomainCommand)(nil)
	_ Transactional = (*AddInstanceDomainCommand)(nil)
)
