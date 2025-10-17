package domain

import (
	"context"
	"errors"

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

// Events implements Commander.
//
// TODO(IAM-Marco): Finish implementation when policies, org settings, idp links and entities repositories
// are implemented
func (d *DeleteOrgCommand) Events(ctx context.Context, opts *InvokeOpts) ([]eventstore.Command, error) {
	closeFunc, err := opts.EnsureTx(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { err = closeFunc(ctx, err) }()

	usernames := []string{}
	// userRepo := opts.usersRepo(pool)
	// users, err := userRepo.List(ctx, database.WithCondition(opts.organizationRepo(pool).IDCondition(d.ID)))
	// if err != nil {
	// 	return nil, err
	// }
	// for _, u := range users {
	// 	usernames = append(usernames, u.UserID)
	// }

	// domainPolicyRepo := opts.domainPolicyRepo(pool)
	// policy, err := domainPolicyRepo.Get(ctx, instanceID, d.ID)
	// if err != nil {
	// 	return nil, err
	// }

	// orgSettingsRepo := opts.organizationSettingsRepo(pool)
	// orgSettings, err := orgSettingsRepo.Get(ctx, d.ID)
	// if err != nil {
	// 	return nil, err
	// }
	// areUsernamesOrganizationScoped := policy.UserLoginMustBeDomain || orgSettings.UsernamesUnique
	areUsernamesOrganizationScoped := false

	domainNames := make([]string, len(d.Domains))
	for i, domain := range d.Domains {
		domainNames[i] = domain.Domain
	}

	externalIDPLinks := []*internal_domain.UserIDPLink{}
	// idpLinksRepo := opts.idpLinksRepo(pool)
	// idpLinks, err := idpLinksRepo.List(ctx, d.ID)
	// if err != nil {
	// 	return nil, err
	// }
	// for _, link := range idpLinks {
	// 	// Convert repo to internal_domain
	// 	externalIDPLinks = append(externalIDPLinks, link)
	// }

	samlEntityIDs := []string{}
	// entityIDsRepo := opts.entityIDsRepo(pool)
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
			&org.NewAggregate(d.ID).Aggregate,
			d.OrganizationName,
			usernames,
			areUsernamesOrganizationScoped,
			domainNames,
			externalIDPLinks,
			samlEntityIDs,
		),
	}, nil
}

// Execute implements Commander.
func (d *DeleteOrgCommand) Execute(ctx context.Context, opts *InvokeOpts) (err error) {
	closeFunc, err := opts.EnsureTx(ctx)
	if err != nil {
		return err
	}

	defer func() { err = closeFunc(ctx, err) }()
	instance := authz.GetInstance(ctx)

	orgRepo := opts.organizationRepo

	orgToDelete, err := orgRepo.Get(ctx, pool, database.WithCondition(database.And(
		orgRepo.IDCondition(d.ID),
		orgRepo.InstanceIDCondition(instance.InstanceID()),
	)))
	if err != nil {
		return err
	}
	d.OrganizationName = orgToDelete.Name
	d.Domains = orgToDelete.Domains

	deletedRows, err := orgRepo.Delete(ctx, pool,
		database.And(
			orgRepo.IDCondition(d.ID),
			orgRepo.InstanceIDCondition(authz.GetInstance(ctx).InstanceID()),
		),
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

// String implements Commander.
func (d *DeleteOrgCommand) String() string {
	return "DeleteOrgCommand"
}

// Validate implements Commander.
func (d *DeleteOrgCommand) Validate(ctx context.Context, opts *InvokeOpts) (err error) {
	instance := authz.GetInstance(ctx)

	if d.ID == instance.DefaultOrganisationID() {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-LCkE69", "Errors.Org.DefaultOrgNotDeletable")
	}

	closeFunc, err := opts.EnsureTx(ctx)
	if err != nil {
		return err
	}

	defer func() { err = closeFunc(ctx, err) }()

	// Check if the ZITADEL project exists on the input organization
	projectRepo := opts.projectRepo
	_, getErr := projectRepo.Get(ctx, pool,
		database.WithCondition(database.And(
			projectRepo.IDCondition(instance.ProjectID()),
			projectRepo.OrganizationIDCondition(d.ID),
			projectRepo.InstanceIDCondition(instance.InstanceID()),
		)),
	)
	if getErr == nil {
		return zerrors.ThrowPreconditionFailed(nil, "DOM-X7YXxC", "Errors.Org.ZitadelOrgNotDeletable")
	}
	// "database.NoRowFoundError" error means the project does not exist, return other errors in case it's not that
	if !errors.Is(getErr, &database.NoRowFoundError{}) {
		err = getErr
		return err
	}

	orgRepo := opts.organizationRepo
	_, errGetOrg := orgRepo.Get(ctx, pool,
		database.WithCondition(database.And(
			orgRepo.IDCondition(d.ID),
			orgRepo.InstanceIDCondition(instance.InstanceID()),
		)))
	if errGetOrg != nil {
		if errors.Is(errGetOrg, &database.NoRowFoundError{}) {
			err = zerrors.ThrowNotFound(errGetOrg, "DOM-8KYOH3", "Errors.Org.NotFound")
		}
		return err
	}

	return err
}

var _ Commander = (*DeleteOrgCommand)(nil)
