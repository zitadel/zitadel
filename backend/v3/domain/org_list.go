package domain

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/zerrors"
	"github.com/zitadel/zitadel/pkg/grpc/object/v2"
	v2_org "github.com/zitadel/zitadel/pkg/grpc/org/v2"
	v2beta_org "github.com/zitadel/zitadel/pkg/grpc/org/v2beta"
)

type ListOrgsQuery struct {
	BaseCommand
	Request *v2_org.ListOrganizationsRequest

	result []*Organization
}

// Result implements [Querier].
func (l *ListOrgsQuery) Result() []*Organization {
	return l.result
}

func NewListOrgsQuery(inputRequest *v2_org.ListOrganizationsRequest) *ListOrgsQuery {
	return &ListOrgsQuery{
		BaseCommand: BaseCommand{},
		Request:     inputRequest,
	}
}

// Execute implements [Querier].
func (l *ListOrgsQuery) Execute(ctx context.Context, opts *InvokeOpts) (err error) {
	close, err := opts.EnsureTx(ctx)
	if err != nil {
		return err
	}
	defer func() { err = close(ctx, err) }()

	organizationRepo := opts.organizationRepo.LoadDomains()
	domainRepo := opts.organizationDomainRepo

	sorting := l.Sorting(organizationRepo)
	limit, pagination := l.Pagination()
	conditions, condErr := l.conditions(ctx, organizationRepo, domainRepo)
	if condErr != nil {
		err = condErr
		return err
	}

	l.result, err = organizationRepo.List(ctx, pool, conditions, sorting, limit, pagination)
	return err
}

func (l *ListOrgsQuery) Sorting(orgRepo OrganizationRepository) database.QueryOption {
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

func (l *ListOrgsQuery) Pagination() (database.QueryOption, database.QueryOption) {
	return database.WithLimit(l.Request.GetQuery().GetLimit()),
		database.WithOffset(uint32(l.Request.GetQuery().GetOffset()))
}

func (l *ListOrgsQuery) conditions(ctx context.Context, orgRepo OrganizationRepository, domainRepo OrganizationDomainRepository) (database.QueryOption, error) {
	conditions := make([]database.Condition, len(l.Request.GetQueries()))
	instance := authz.GetInstance(ctx)

	for i, query := range l.Request.GetQueries() {
		switch assertedType := query.GetQuery().(type) {

		case *v2_org.SearchQuery_DefaultQuery:
			conditions[i] = orgRepo.IDCondition(instance.DefaultOrganisationID())
		case *v2_org.SearchQuery_DomainQuery:
			method, err := l.TextOperationMapper(assertedType.DomainQuery.GetMethod())
			if err != nil {
				return nil, err
			}

			conditions[i] = database.And(
				orgRepo.InstanceIDCondition(instance.InstanceID()),
				orgRepo.ExistsDomain(domainRepo.DomainCondition(method, assertedType.DomainQuery.GetDomain())),
			)
		case *v2_org.SearchQuery_IdQuery:
			conditions[i] = orgRepo.IDCondition(assertedType.IdQuery.GetId())
		case *v2_org.SearchQuery_NameQuery:
			method, err := l.TextOperationMapper(assertedType.NameQuery.GetMethod())
			if err != nil {
				return nil, err
			}
			conditions[i] = orgRepo.NameCondition(method, assertedType.NameQuery.GetName())
		case *v2_org.SearchQuery_StateQuery:
			conditions[i] = orgRepo.StateCondition(OrgState(assertedType.StateQuery.GetState()))
		default:
			return nil, zerrors.ThrowInvalidArgument(NewUnexpectedQueryTypeError(assertedType), "DOM-TCEzcr", "List.Query.Invalid")
		}
	}

	return database.WithCondition(database.And(
		append(conditions, orgRepo.InstanceIDCondition(instance.InstanceID()))...,
	)), nil
}

// String implements [Querier].
func (l *ListOrgsQuery) String() string {
	return "ListOrgsQuery"
}

// Validate implements [Querier].
func (l *ListOrgsQuery) Validate(_ context.Context, _ *InvokeOpts) (err error) {
	return nil
}

var _ Querier[[]*Organization] = (*ListOrgsQuery)(nil)

func (l *ListOrgsQuery) ResultToGRPC() []*v2_org.Organization {
	toReturn := make([]*v2_org.Organization, len(l.result))

	for i, org := range l.result {
		toReturn[i] = l.orgToGRPC(org)
	}

	return toReturn
}

func (l *ListOrgsQuery) orgToGRPC(org *Organization) *v2_org.Organization {
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

// TODO(IAM-Marco): Remove in V5 (see https://github.com/zitadel/zitadel/issues/10877)
func (l *ListOrgsQuery) ResultToGRPCBeta() []*v2beta_org.Organization {
	toReturn := make([]*v2beta_org.Organization, len(l.result))

	for i, org := range l.result {
		toReturn[i] = l.orgToGRPCBeta(org)
	}

	return toReturn
}

// TODO(IAM-Marco): Remove in V5 (see https://github.com/zitadel/zitadel/issues/10877)
func (l *ListOrgsQuery) orgToGRPCBeta(org *Organization) *v2beta_org.Organization {
	return &v2beta_org.Organization{
		Id:           org.ID,
		ChangedDate:  timestamppb.New(org.UpdatedAt),
		CreationDate: timestamppb.New(org.CreatedAt),

		State:         v2beta_org.OrgState(org.State),
		Name:          org.Name,
		PrimaryDomain: org.PrimaryDomain(),
	}
}
