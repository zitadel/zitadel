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

const maxDomainNameLength = 253

type AddInstanceDomainCommand struct {
	InstanceID string     `json:"instance_id"`
	DomainName string     `json:"domain_name"`
	DType      DomainType `json:"domain_type"`
}

// RequiresTransaction implements [Transactional].
func (a *AddInstanceDomainCommand) RequiresTransaction() {}

func NewAddInstanceDomainCommand(instanceID, domainName string, domainType DomainType) *AddInstanceDomainCommand {
	return &AddInstanceDomainCommand{
		InstanceID: instanceID,
		DomainName: domainName,
		DType:      domainType,
	}
}

// Events implements [Commander].
func (a *AddInstanceDomainCommand) Events(ctx context.Context, _ *InvokeOpts) ([]eventstore.Command, error) {
	var toReturn eventstore.Command
	switch a.DType {
	case DomainTypeCustom:
		toReturn = instance.NewDomainAddedEvent(ctx, &instance.NewAggregate(a.InstanceID).Aggregate, a.DomainName, false)
	case DomainTypeTrusted:
		toReturn = instance.NewTrustedDomainAddedEvent(ctx, &instance.NewAggregate(a.InstanceID).Aggregate, a.DomainName)
	default:
		return nil, nil
	}

	return []eventstore.Command{
		toReturn,
	}, nil
}

// Execute implements [Commander].
func (a *AddInstanceDomainCommand) Execute(ctx context.Context, opts *InvokeOpts) (err error) {
	instanceDomainRepo := opts.instanceDomainRepo

	var isPrimary, isGenerated *bool
	if a.DType == DomainTypeCustom {
		isPrimary, isGenerated = gu.Ptr(false), gu.Ptr(false)
	}

	err = instanceDomainRepo.Add(ctx, opts.DB(), &AddInstanceDomain{
		InstanceID:  a.InstanceID,
		Domain:      a.DomainName,
		IsPrimary:   isPrimary,
		IsGenerated: isGenerated,
		Type:        a.DType,
	})
	if err != nil {
		return zerrors.ThrowInternal(err, "DOM-uSCVn3", "Errors.Instance.Domain.Add")
	}

	return nil
}

// String implements [Commander].
func (a *AddInstanceDomainCommand) String() string {
	return "AddInstanceDomainCommand"
}

// Validate implements [Commander].
func (a *AddInstanceDomainCommand) Validate(ctx context.Context, opts *InvokeOpts) (err error) {
	if a.DomainName = strings.TrimSpace(a.DomainName); a.DomainName == "" || len(a.DomainName) > maxDomainNameLength {
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
		return zerrors.ThrowPermissionDenied(authZErr, "DOM-c83vPX", "Errors.PermissionDenied")
	}

	domainRepo := opts.instanceDomainRepo
	_, err = domainRepo.Get(ctx, opts.DB(), database.WithCondition(domainRepo.DomainCondition(database.TextOperationEqual, a.DomainName)))
	if err == nil {
		return zerrors.ThrowAlreadyExists(nil, "DOM-CvQ8tf", "Errors.Instance.Domain.AlreadyExists")
	}

	if !errors.Is(err, &database.NoRowFoundError{}) {
		return zerrors.ThrowInternal(err, "DOM-LrTy2z", "Errors.Instance.Domain.Get")
	}

	return nil
}

var (
	_ Commander     = (*AddInstanceDomainCommand)(nil)
	_ Transactional = (*AddInstanceDomainCommand)(nil)
)
