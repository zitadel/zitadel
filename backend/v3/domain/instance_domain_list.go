package domain

import (
	"context"
	"strings"

	"github.com/zitadel/zitadel/backend/v3/api/object"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/zerrors"
	instance_v2 "github.com/zitadel/zitadel/pkg/grpc/instance/v2"
)

type ListInstanceDomainsQuery struct {
	Request *instance_v2.ListCustomDomainsRequest

	toReturn []*InstanceDomain
}

var _ Querier[[]*InstanceDomain] = (*ListInstanceDomainsQuery)(nil)

func NewListInstanceDomainsQuery(inputRequest *instance_v2.ListCustomDomainsRequest) *ListInstanceDomainsQuery {
	return &ListInstanceDomainsQuery{
		Request: inputRequest,
	}
}

// Result implements [Querier].
func (l *ListInstanceDomainsQuery) Result() []*InstanceDomain {
	return l.toReturn
}

// Execute implements [Querier].
func (l *ListInstanceDomainsQuery) Execute(ctx context.Context, opts *InvokeOpts) (err error) {
	domainRepo := opts.instanceDomainRepo

	sorting := l.Sorting(domainRepo)
	limit, offset := l.Pagination()
	conds, err := l.conditions(domainRepo)
	if err != nil {
		return err
	}

	instances, err := domainRepo.List(ctx, opts.DB(), conds, sorting, limit, offset)
	if err != nil {
		return zerrors.ThrowInternal(err, "DOM-ubaPNU", "failed fetching instance domains")
	}

	l.toReturn = instances

	return nil
}

func (l *ListInstanceDomainsQuery) Pagination() (database.QueryOption, database.QueryOption) {
	return database.WithLimit(l.Request.GetPagination().GetLimit()),
		database.WithOffset(uint32(l.Request.GetPagination().GetOffset()))
}

func (l *ListInstanceDomainsQuery) conditions(domainRepo InstanceDomainRepository) (database.QueryOption, error) {
	conditions := make([]database.Condition, len(l.Request.GetFilters()))
	for i, query := range l.Request.GetFilters() {
		switch assertedType := query.GetFilter().(type) {
		case *instance_v2.CustomDomainFilter_DomainFilter:
			method, err := object.TextQueryMethodToTextOperation(assertedType.DomainFilter.GetMethod())
			if err != nil {
				return nil, err
			}
			conditions[i] = domainRepo.DomainCondition(
				method,
				assertedType.DomainFilter.GetDomain(),
			)
		case *instance_v2.CustomDomainFilter_GeneratedFilter:
			conditions[i] = domainRepo.IsGeneratedCondition(assertedType.GeneratedFilter)
		case *instance_v2.CustomDomainFilter_PrimaryFilter:
			conditions[i] = domainRepo.IsPrimaryCondition(assertedType.PrimaryFilter)
		default:
			return nil, zerrors.ThrowInvalidArgument(NewUnexpectedQueryTypeError(assertedType), "DOM-CjM93P", "List.Query.Invalid")
		}
	}

	conditions = append(conditions, domainRepo.TypeCondition(DomainTypeCustom))
	return database.WithCondition(database.And(conditions...)), nil
}

func (l *ListInstanceDomainsQuery) Sorting(domainRepo InstanceDomainRepository) database.QueryOption {
	var sortingCol database.Column

	switch l.Request.GetSortingColumn() {

	case instance_v2.DomainFieldName_DOMAIN_FIELD_NAME_CREATION_DATE:
		sortingCol = domainRepo.CreatedAtColumn()
	case instance_v2.DomainFieldName_DOMAIN_FIELD_NAME_DOMAIN:
		sortingCol = domainRepo.DomainColumn()
	case instance_v2.DomainFieldName_DOMAIN_FIELD_NAME_GENERATED:
		sortingCol = domainRepo.IsGeneratedColumn()
	case instance_v2.DomainFieldName_DOMAIN_FIELD_NAME_PRIMARY:
		sortingCol = domainRepo.IsPrimaryColumn()
	case instance_v2.DomainFieldName_DOMAIN_FIELD_NAME_UNSPECIFIED:
		fallthrough
	default:
		return func(opts *database.QueryOpts) {}
	}

	if l.Request.GetPagination().GetAsc() {
		return database.WithOrderByAscending(sortingCol)
	}

	return database.WithOrderByDescending(sortingCol)
}

// String implements [Querier].
func (l *ListInstanceDomainsQuery) String() string {
	return "ListInstanceDomainsQuery"
}

// Validate implements [Querier].
func (l *ListInstanceDomainsQuery) Validate(ctx context.Context, opts *InvokeOpts) (err error) {
	instanceID := strings.TrimSpace(l.Request.GetInstanceId())
	if instanceID == "" || instanceID == authz.GetInstance(ctx).InstanceID() {
		return l.checkDomainPerms(ctx, opts)
	}

	// TODO(IAM-Marco): This is wrong, as it should check the permission for the input instance and not the one in context
	if authZErr := opts.Permissions.CheckInstancePermission(ctx, InstanceReadPermission); authZErr != nil {
		return zerrors.ThrowPermissionDenied(authZErr, "DOM-yN7oCp", "permission denied")
	}

	return l.checkDomainPerms(ctx, opts)
}

func (l *ListInstanceDomainsQuery) checkDomainPerms(ctx context.Context, opts *InvokeOpts) error {
	if authZErr := opts.Permissions.CheckInstancePermission(ctx, DomainReadPermission); authZErr != nil {
		return zerrors.ThrowPermissionDenied(authZErr, "DOM-RyCEyr", "permission denied")
	}
	return nil
}
