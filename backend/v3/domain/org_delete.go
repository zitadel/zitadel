package domain

import (
	"context"
	"errors"
	"strings"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/internal/api/authz"
	internal_domain "github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type DeleteOrgCommand struct {
	OrganizationName string `json:"organization_name"`
	Domains          []*OrganizationDomain
	ID               string `json:"id"`
}

func NewDeleteOrgCommand(organizationID string) *DeleteOrgCommand {
	return &DeleteOrgCommand{ID: organizationID}
}

// RequiresTransaction implements [Transactional].
func (cmd *DeleteOrgCommand) RequiresTransaction() {}

// Events implements [Commander].
//
// TODO(IAM-Marco): Finish implementation when policies, org settings, idp links and entities repositories
// are implemented
func (cmd *DeleteOrgCommand) Events(ctx context.Context, opts *InvokeOpts) ([]eventstore.Command, error) {
	usernames := []string{}
	// userRepo := opts.usersRepo(opts.DB())
	// users, err := userRepo.List(ctx, database.WithCondition(opts.organizationRepo(opts.DB()).IDCondition(d.ID)))
	// if err != nil {
	// 	return nil, err
	// }
	// for _, u := range users {
	// 	usernames = append(usernames, u.UserID)
	// }

	// domainPolicyRepo := opts.domainPolicyRepo(opts.DB())
	// policy, err := domainPolicyRepo.Get(ctx, instanceID, d.ID)
	// if err != nil {
	// 	return nil, err
	// }

	// orgSettingsRepo := opts.organizationSettingsRepo(opts.DB())
	// orgSettings, err := orgSettingsRepo.Get(ctx, d.ID)
	// if err != nil {
	// 	return nil, err
	// }
	// areUsernamesOrganizationScoped := policy.UserLoginMustBeDomain || orgSettings.UsernamesUnique
	areUsernamesOrganizationScoped := false

	domainNames := make([]string, len(cmd.Domains))
	for i, domain := range cmd.Domains {
		domainNames[i] = domain.Domain
	}

	externalIDPLinks := []*internal_domain.UserIDPLink{}
	// idpLinksRepo := opts.idpLinksRepo(opts.DB())
	// idpLinks, err := idpLinksRepo.List(ctx, d.ID)
	// if err != nil {
	// 	return nil, err
	// }
	// for _, link := range idpLinks {
	// 	// Convert repo to internal_domain
	// 	externalIDPLinks = append(externalIDPLinks, link)
	// }

	samlEntityIDs := []string{}
	// entityIDsRepo := opts.entityIDsRepo(opts.DB())
	// entityIDs, err := entityIDsRepo.List(ctx, d.ID)
	// if err != nil {
	// 	return nil, err
	// }
	// for _, ei := range entityIDs {
	// 	samlEntityIDs = append(samlEntityIDs, ei.ID)
	// }

	return []eventstore.Command{
		org.NewOrgRemovedEvent(
			ctx,
			&org.NewAggregate(cmd.ID).Aggregate,
			cmd.OrganizationName,
			usernames,
			areUsernamesOrganizationScoped,
			domainNames,
			externalIDPLinks,
			samlEntityIDs,
		),
	}, nil
}

// Execute implements [Commander].
func (cmd *DeleteOrgCommand) Execute(ctx context.Context, opts *InvokeOpts) (err error) {
	instance := authz.GetInstance(ctx)

	orgRepo := opts.organizationRepo.LoadDomains()

	orgToDelete, err := orgRepo.Get(ctx, opts.DB(), database.WithCondition(
		orgRepo.PrimaryKeyCondition(instance.InstanceID(), cmd.ID),
	))
	if err != nil {
		return err
	}
	cmd.OrganizationName = orgToDelete.Name
	cmd.Domains = orgToDelete.Domains

	deletedRows, err := orgRepo.Delete(ctx, opts.DB(),
		orgRepo.PrimaryKeyCondition(authz.GetInstance(ctx).InstanceID(), cmd.ID),
	)
	if err != nil {
		return err
	}

	if deletedRows > 1 {
		err = zerrors.ThrowInternalf(nil, "DOM-5cE9u6", "expecting 1 row deleted, got %d", deletedRows)
		return err
	}

	if deletedRows < 1 {
		err = zerrors.ThrowNotFound(nil, "DOM-ur6Qyv", "organization not found")
	}
	return err
}

// String implements [Commander].
func (DeleteOrgCommand) String() string {
	return "DeleteOrgCommand"
}

// Validate implements [Commander].
func (cmd *DeleteOrgCommand) Validate(ctx context.Context, opts *InvokeOpts) (err error) {
	instance := authz.GetInstance(ctx)

	cmd.ID = strings.TrimSpace(cmd.ID)
	if cmd.ID == instance.DefaultOrganisationID() {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-LCkE69", "Errors.Org.DefaultOrgNotDeletable")
	}

	// Check if the ZITADEL project exists on the input organization
	projectRepo := opts.projectRepo
	_, getErr := projectRepo.Get(ctx, opts.DB(),
		database.WithCondition(database.And(
			projectRepo.IDCondition(instance.ProjectID()),
			projectRepo.OrganizationIDCondition(cmd.ID),
			projectRepo.InstanceIDCondition(instance.InstanceID()),
		)),
	)
	if getErr == nil {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-X7YXxC", "Errors.Org.ZitadelOrgNotDeletable")
	}
	// [database.NoRowFoundError] error means the project does not exist, return other errors in case it's not that
	if !errors.Is(getErr, &database.NoRowFoundError{}) {
		err = getErr
		return err
	}

	orgRepo := opts.organizationRepo
	_, errGetOrg := orgRepo.Get(ctx, opts.DB(),
		database.WithCondition(
			orgRepo.PrimaryKeyCondition(instance.InstanceID(), cmd.ID),
		))
	if errGetOrg != nil {
		if errors.Is(errGetOrg, &database.NoRowFoundError{}) {
			err = zerrors.ThrowNotFound(errGetOrg, "DOM-8KYOH3", "Errors.Org.NotFound")
		}
		return err
	}

	return err
}

var (
	_ Commander     = (*DeleteOrgCommand)(nil)
	_ Transactional = (*DeleteOrgCommand)(nil)
)
