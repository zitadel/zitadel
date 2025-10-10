package domain

import (
	"context"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/internal/zerrors"
	instance "github.com/zitadel/zitadel/pkg/grpc/instance/v2beta"
)

type ListInstancesQuery struct {
	Request *instance.ListInstancesRequest

	toReturn []*Instance
}

// Result implements Querier.
func (l *ListInstancesQuery) Result() []*Instance {
	return l.toReturn
}

var _ Querier[[]*Instance] = (*ListInstancesQuery)(nil)

func NewListInstancesCommand(inputRequest *instance.ListInstancesRequest) *ListInstancesQuery {
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

	instances, err := instanceRepo.List(ctx, pool, conds, sorting, limit, offset)
	if err != nil {
		return err
	}

	l.toReturn = instances

	return nil
}

func (l *ListInstancesQuery) Pagination() (database.QueryOption, database.QueryOption) {
	return database.WithLimit(l.Request.GetPagination().GetLimit()),
		database.WithOffset(uint32(l.Request.GetPagination().GetOffset()))
}

func (l *ListInstancesQuery) conditions(instanceRepo InstanceRepository, domainRepo InstanceDomainRepository) (database.QueryOption, error) {
	conditions := make([]database.Condition, len(l.Request.GetQueries()))
	for i, query := range l.Request.GetQueries() {
		switch assertedType := query.GetQuery().(type) {
		case *instance.Query_DomainQuery:
			domainConditions := make([]database.Condition, len(assertedType.DomainQuery.GetDomains()))
			for j, domain := range assertedType.DomainQuery.GetDomains() {
				domainConditions[j] = domainRepo.DomainCondition(database.TextOperationEqual, domain)
			}
			conditions[i] = instanceRepo.ExistsDomain(
				database.Or(domainConditions...),
			)
		case *instance.Query_IdQuery:
			idConditions := make([]database.Condition, len(assertedType.IdQuery.GetIds()))
			for j, id := range assertedType.IdQuery.GetIds() {
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
	case instance.FieldName_FIELD_NAME_CREATION_DATE:
		sortingCol = instanceRepo.CreatedAtColumn()
	case instance.FieldName_FIELD_NAME_ID:
		sortingCol = instanceRepo.IDColumn()
	case instance.FieldName_FIELD_NAME_NAME:
		sortingCol = instanceRepo.NameColumn()
	case instance.FieldName_FIELD_NAME_UNSPECIFIED:
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
	return nil
}
