package domain

import (
	"context"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/internal/zerrors"
	instance_v2 "github.com/zitadel/zitadel/pkg/grpc/instance/v2"
)

type ListInstancesQuery struct {
	Request *instance_v2.ListInstancesRequest

	toReturn []*Instance
}

// Result implements Querier.
func (l *ListInstancesQuery) Result() []*Instance {
	return l.toReturn
}

var _ Querier[[]*Instance] = (*ListInstancesQuery)(nil)

func NewListInstancesCommand(inputRequest *instance_v2.ListInstancesRequest) *ListInstancesQuery {
	return &ListInstancesQuery{
		Request: inputRequest,
	}
}

// Execute implements [Querier].
func (l *ListInstancesQuery) Execute(ctx context.Context, opts *InvokeOpts) (err error) {
	instanceRepo := opts.instanceRepo.LoadDomains()
	domainRepo := opts.instanceDomainRepo

	sorting := l.Sorting(instanceRepo)
	limit, offset := l.Pagination()
	conds, err := l.conditions(instanceRepo, domainRepo)
	if err != nil {
		return err
	}

	instances, err := instanceRepo.List(ctx, opts.DB(), conds, sorting, limit, offset)
	if err != nil {
		return zerrors.ThrowInternal(err, "DOM-AIRPxN", "failed fetching instances")
	}

	l.toReturn = instances

	return nil
}

func (l *ListInstancesQuery) Pagination() (database.QueryOption, database.QueryOption) {
	return database.WithLimit(l.Request.GetPagination().GetLimit()),
		database.WithOffset(uint32(l.Request.GetPagination().GetOffset()))
}

func (l *ListInstancesQuery) conditions(instanceRepo InstanceRepository, domainRepo InstanceDomainRepository) (database.QueryOption, error) {
	conditions := make([]database.Condition, len(l.Request.GetFilters()))
	for i, query := range l.Request.GetFilters() {
		switch assertedType := query.GetFilter().(type) {
		case *instance_v2.Filter_CustomDomainsFilter:
			domainConditions := make([]database.Condition, len(assertedType.CustomDomainsFilter.GetDomains()))
			for j, domain := range assertedType.CustomDomainsFilter.GetDomains() {
				domainConditions[j] = domainRepo.DomainCondition(database.TextOperationEqual, domain)
			}
			conditions[i] = instanceRepo.ExistsDomain(
				database.Or(domainConditions...),
			)
		case *instance_v2.Filter_InIdsFilter:
			idConditions := make([]database.Condition, len(assertedType.InIdsFilter.GetIds()))
			for j, id := range assertedType.InIdsFilter.GetIds() {
				idConditions[j] = instanceRepo.IDCondition(id)
			}
			conditions[i] = database.Or(idConditions...)
		default:
			return nil, zerrors.ThrowInvalidArgument(NewUnexpectedQueryTypeError(assertedType), "DOM-AU99kR", "List.Query.Invalid")
		}
	}

	return database.WithCondition(database.And(conditions...)), nil
}

func (l *ListInstancesQuery) Sorting(instanceRepo InstanceRepository) database.QueryOption {
	var sortingCol database.Column

	switch l.Request.GetSortingColumn() {
	case instance_v2.FieldName_FIELD_NAME_CREATION_DATE:
		sortingCol = instanceRepo.CreatedAtColumn()
	case instance_v2.FieldName_FIELD_NAME_ID:
		sortingCol = instanceRepo.IDColumn()
	case instance_v2.FieldName_FIELD_NAME_NAME:
		sortingCol = instanceRepo.NameColumn()
	case instance_v2.FieldName_FIELD_NAME_UNSPECIFIED:
		fallthrough
	default:
		return func(opts *database.QueryOpts) {}
	}

	orderDirection := database.OrderDirectionDesc
	if l.Request.GetPagination().GetAsc() {
		orderDirection = database.OrderDirectionAsc
	}

	return database.WithOrderBy(orderDirection, sortingCol)
}

// String implements [Querier].
func (l *ListInstancesQuery) String() string {
	return "ListInstancesCommand"
}

// Validate implements [Querier].
func (l *ListInstancesQuery) Validate(ctx context.Context, opts *InvokeOpts) (err error) {
	// TODO(IAM-Marco): This is likely wrong, because it's checking the read permission of
	// the instance in context. In here, instead we should probably loop through all retrieved
	// instances and check their permissions.
	if authZErr := opts.Permissions.CheckInstancePermission(ctx, InstanceReadPermission); authZErr != nil {
		return zerrors.ThrowPermissionDenied(authZErr, "DOM-cuT6Ws", "permission denied")
	}

	return nil
}
