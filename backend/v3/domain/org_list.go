package domain

import (
	"context"
	"errors"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/pkg/grpc/object/v2"
	v2_org "github.com/zitadel/zitadel/pkg/grpc/org/v2"
	v2beta_org "github.com/zitadel/zitadel/pkg/grpc/org/v2beta"
)

type ListOrgsCommand struct {
	BaseCommand
	Request *v2_org.ListOrganizationsRequest

	Result []*Organization
}

func NewListOrgsCommand(inputRequest *v2_org.ListOrganizationsRequest) *ListOrgsCommand {
	return &ListOrgsCommand{
		BaseCommand: BaseCommand{},
		Request:     inputRequest,
	}
}

// Events implements Commander.
func (l *ListOrgsCommand) Events(ctx context.Context) []eventstore.Command {
	return nil
}

// Execute implements Commander.
func (l *ListOrgsCommand) Execute(ctx context.Context, opts *CommandOpts) (err error) {
	close, err := opts.EnsureClient(ctx)
	if err != nil {
		return err
	}
	defer func() {
		closeErr := close(ctx)
		if closeErr != nil {
			err = errors.Join(err, closeErr)
		}
	}()

	organizationRepo := opts.organizationRepo(pool)
	domainRepo := organizationRepo.Domains(true)

	sorting := l.Sorting(organizationRepo)
	limit, pagination := l.Pagination()
	conditions, condErr := l.conditions(ctx, organizationRepo, domainRepo)
	if condErr != nil {
		err = condErr
		return err
	}

	l.Result, err = organizationRepo.List(ctx, append(conditions, sorting, limit, pagination)...)
	return err
}

func (l *ListOrgsCommand) Sorting(orgRepo OrganizationRepository) database.QueryOption {
	var sortingCol database.Column
	switch l.Request.GetSortingColumn() {
	case v2_org.OrganizationFieldName_ORGANIZATION_FIELD_NAME_NAME:
		sortingCol = orgRepo.NameColumn()
	case v2_org.OrganizationFieldName_ORGANIZATION_FIELD_NAME_UNSPECIFIED:
		fallthrough
	default:
		return func(opts *database.QueryOpts) {}
	}

	orderDirection := database.OrderDirectionDesc
	if l.Request.GetQuery().GetAsc() {
		orderDirection = database.OrderDirectionAsc
	}

	return database.WithOrderBy(orderDirection, sortingCol)
}

func (l *ListOrgsCommand) Pagination() (database.QueryOption, database.QueryOption) {
	return database.WithLimit(l.Request.GetQuery().GetLimit()),
		database.WithOffset(uint32(l.Request.GetQuery().GetOffset()))
}

func (l *ListOrgsCommand) conditions(ctx context.Context, orgRepo OrganizationRepository, domainRepo OrganizationDomainRepository) ([]database.QueryOption, error) {
	conditions := make([]database.QueryOption, len(l.Request.GetQueries()))
	for i, query := range l.Request.GetQueries() {
		switch assertedType := query.GetQuery().(type) {

		case *v2_org.SearchQuery_DefaultQuery:
			conditions[i] = database.WithCondition(orgRepo.IDCondition(authz.GetInstance(ctx).DefaultOrganisationID()))
		case *v2_org.SearchQuery_DomainQuery:
			method, err := l.TextOperationMapper(assertedType.DomainQuery.GetMethod())
			if err != nil {
				return nil, err
			}

			conditions[i] = database.WithCondition(
				database.And(
					orgRepo.InstanceIDCondition(authz.GetInstance(ctx).InstanceID()),
					orgRepo.ExistsDomain(domainRepo.DomainCondition(method, assertedType.DomainQuery.GetDomain())),
				),
			)
		case *v2_org.SearchQuery_IdQuery:
			conditions[i] = database.WithCondition(orgRepo.IDCondition(assertedType.IdQuery.GetId()))
		case *v2_org.SearchQuery_NameQuery:
			method, err := l.TextOperationMapper(assertedType.NameQuery.GetMethod())
			if err != nil {
				return nil, err
			}
			conditions[i] = database.WithCondition(orgRepo.NameCondition(method, assertedType.NameQuery.GetName()))
		case *v2_org.SearchQuery_StateQuery:
			conditions[i] = database.WithCondition(orgRepo.StateCondition(OrgState(assertedType.StateQuery.GetState())))
		default:
			return nil, NewUnexpectedQueryTypeError("DOM-TCEzcr", assertedType)
		}
	}
	return conditions, nil
}

// String implements Commander.
func (l *ListOrgsCommand) String() string {
	return "ListOrgsCommand"
}

// Validate implements Commander.
func (l *ListOrgsCommand) Validate(_ context.Context, _ *CommandOpts) (err error) {
	if len(l.Request.GetQueries()) == 0 {
		return NewNoQueryCriteriaError("DOM-75lU26")
	}

	return nil
}

var _ Commander = (*ListOrgsCommand)(nil)

func (l *ListOrgsCommand) ResultToGRPC() []*v2_org.Organization {
	toReturn := make([]*v2_org.Organization, len(l.Result))

	for i, org := range l.Result {
		toReturn[i] = l.orgToGRPC(org)
	}

	return toReturn
}

func (l *ListOrgsCommand) orgToGRPC(org *Organization) *v2_org.Organization {
	return &v2_org.Organization{
		Id: org.ID,
		Details: &object.Details{
			ChangeDate:   timestamppb.New(org.UpdatedAt),
			CreationDate: timestamppb.New(org.CreatedAt),
		},
		State:         v2_org.OrganizationState(org.State),
		Name:          org.Name,
		PrimaryDomain: org.PrimaryDomain(),
	}
}

// TODO(IAM-Marco): Remove in V5
func (l *ListOrgsCommand) ResultToGRPCBeta() []*v2beta_org.Organization {
	toReturn := make([]*v2beta_org.Organization, len(l.Result))

	for i, org := range l.Result {
		toReturn[i] = l.orgToGRPCBeta(org)
	}

	return toReturn
}

// TODO(IAM-Marco): Remove in V5
func (l *ListOrgsCommand) orgToGRPCBeta(org *Organization) *v2beta_org.Organization {
	return &v2beta_org.Organization{
		Id:           org.ID,
		ChangedDate:  timestamppb.New(org.UpdatedAt),
		CreationDate: timestamppb.New(org.CreatedAt),

		State:         v2beta_org.OrgState(org.State),
		Name:          org.Name,
		PrimaryDomain: org.PrimaryDomain(),
	}
}
