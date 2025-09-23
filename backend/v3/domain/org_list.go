package domain

import (
	"context"
	"errors"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/pkg/grpc/object/v2"
	"github.com/zitadel/zitadel/pkg/grpc/org/v2"
	v2_org "github.com/zitadel/zitadel/pkg/grpc/org/v2"
)

type ListOrgsCommand struct {
	Request *org.ListOrganizationsRequest

	Result []*Organization
}

func NewListOrgsCommand(inputRequest *org.ListOrganizationsRequest) *ListOrgsCommand {
	return &ListOrgsCommand{
		Request: inputRequest,
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
	case org.OrganizationFieldName_ORGANIZATION_FIELD_NAME_NAME:
		sortingCol = orgRepo.NameColumn()
	case org.OrganizationFieldName_ORGANIZATION_FIELD_NAME_UNSPECIFIED:
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

		case *org.SearchQuery_DefaultQuery:
			conditions[i] = database.WithCondition(orgRepo.IDCondition(authz.GetInstance(ctx).DefaultOrganisationID()))
		case *org.SearchQuery_DomainQuery:
			method, err := l.OperationMapper(assertedType.DomainQuery.GetMethod())
			if err != nil {
				return nil, err
			}

			conditions[i] = database.WithCondition(
				database.And(
					orgRepo.InstanceIDCondition(authz.GetInstance(ctx).InstanceID()),
					orgRepo.ExistsDomain(domainRepo.DomainCondition(method, assertedType.DomainQuery.GetDomain())),
				),
			)
		case *org.SearchQuery_IdQuery:
			conditions[i] = database.WithCondition(orgRepo.IDCondition(assertedType.IdQuery.GetId()))
		case *org.SearchQuery_NameQuery:
			method, err := l.OperationMapper(assertedType.NameQuery.GetMethod())
			if err != nil {
				return nil, err
			}
			conditions[i] = database.WithCondition(orgRepo.NameCondition(method, assertedType.NameQuery.GetName()))
		case *org.SearchQuery_StateQuery:
			conditions[i] = database.WithCondition(orgRepo.StateCondition(OrgState(assertedType.StateQuery.GetState())))
		default:
			return nil, NewUnexpectedQueryTypeError("DOM-TCEzcr", assertedType)
		}
	}
	return conditions, nil
}

func (l *ListOrgsCommand) OperationMapper(queryOperation object.TextQueryMethod) (database.TextOperation, error) {
	switch queryOperation {
	case object.TextQueryMethod_TEXT_QUERY_METHOD_CONTAINS:
		return database.TextOperationContains, nil
	case object.TextQueryMethod_TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE:
		return database.TextOperationContainsIgnoreCase, nil
	case object.TextQueryMethod_TEXT_QUERY_METHOD_ENDS_WITH:
		return database.TextOperationEndsWith, nil
	case object.TextQueryMethod_TEXT_QUERY_METHOD_ENDS_WITH_IGNORE_CASE:
		return database.TextOperationEndsWithIgnoreCase, nil
	case object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS:
		return database.TextOperationEqual, nil
	case object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS_IGNORE_CASE:
		return database.TextOperationEqualIgnoreCase, nil
	case object.TextQueryMethod_TEXT_QUERY_METHOD_STARTS_WITH:
		return database.TextOperationStartsWith, nil
	case object.TextQueryMethod_TEXT_QUERY_METHOD_STARTS_WITH_IGNORE_CASE:
		return database.TextOperationStartsWithIgnoreCase, nil
	default:
		return 0, NewUnexpectedTextQueryOperationError("DOM-iBRBVe", queryOperation)
	}
}

// String implements Commander.
func (l *ListOrgsCommand) String() string {
	return "ListOrgsCommand"
}

// Validate implements Commander.
func (l *ListOrgsCommand) Validate() (err error) {
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
