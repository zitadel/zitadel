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

type ListInstanceTrustedDomainsQuery struct {
	Request *instance_v2.ListTrustedDomainsRequest

	toReturn []*InstanceDomain
}

var _ Querier[[]*InstanceDomain] = (*ListInstanceTrustedDomainsQuery)(nil)

func NewListInstanceTrustedDomainsQuery(inputRequest *instance_v2.ListTrustedDomainsRequest) *ListInstanceTrustedDomainsQuery {
	return &ListInstanceTrustedDomainsQuery{
		Request: inputRequest,
	}
}

// Result implements [Querier].
func (l *ListInstanceTrustedDomainsQuery) Result() []*InstanceDomain {
	return l.toReturn
}

// Execute implements [Querier].
func (l *ListInstanceTrustedDomainsQuery) Execute(ctx context.Context, opts *InvokeOpts) (err error) {
	domainRepo := opts.instanceDomainRepo

	sorting := l.Sorting(domainRepo)
	limit, offset := l.Pagination()
	conds, err := l.conditions(domainRepo)
	if err != nil {
		return err
	}

	instances, err := domainRepo.List(ctx, opts.DB(), conds, sorting, limit, offset)
	if err != nil {
		return err
	}

	l.toReturn = instances

	return nil
}

func (l *ListInstanceTrustedDomainsQuery) Pagination() (database.QueryOption, database.QueryOption) {
	return database.WithLimit(l.Request.GetPagination().GetLimit()),
		database.WithOffset(uint32(l.Request.GetPagination().GetOffset()))
}

func (l *ListInstanceTrustedDomainsQuery) conditions(domainRepo InstanceDomainRepository) (database.QueryOption, error) {
	conditions := make([]database.Condition, len(l.Request.GetFilters()))
	for i, query := range l.Request.GetFilters() {
		switch assertedType := query.GetFilter().(type) {

		case *instance_v2.TrustedDomainFilter_DomainFilter:
			method, err := object.TextQueryMethodToTextOperation(assertedType.DomainFilter.GetMethod())
			if err != nil {
				return nil, err
			}
			conditions[i] = domainRepo.DomainCondition(
				method,
				assertedType.DomainFilter.GetDomain(),
			)

		default:
			return nil, zerrors.ThrowInvalidArgument(NewUnexpectedQueryTypeError(assertedType), "DOM-qMUBMr", "List.Query.Invalid")
		}
	}

	conditions = append(conditions, domainRepo.TypeCondition(DomainTypeTrusted))
	return database.WithCondition(database.And(conditions...)), nil
}

func (l *ListInstanceTrustedDomainsQuery) Sorting(domainRepo InstanceDomainRepository) database.QueryOption {
	var sortingCol database.Column

	switch l.Request.GetSortingColumn() {
	case instance_v2.TrustedDomainFieldName_TRUSTED_DOMAIN_FIELD_NAME_CREATION_DATE:
		sortingCol = domainRepo.CreatedAtColumn()

	case instance_v2.TrustedDomainFieldName_TRUSTED_DOMAIN_FIELD_NAME_DOMAIN:
		sortingCol = domainRepo.DomainColumn()
	case instance_v2.TrustedDomainFieldName_TRUSTED_DOMAIN_FIELD_NAME_UNSPECIFIED:
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
func (l *ListInstanceTrustedDomainsQuery) String() string {
	return "ListInstanceTrustedDomainsQuery"
}

// Validate implements [Querier].
func (l *ListInstanceTrustedDomainsQuery) Validate(ctx context.Context, opts *InvokeOpts) (err error) {
	// TODO(IAM-Marco): These validations are copied from custom domains.
	// The eventstore side has no validation in place. To decide if we want to keep this or remove it
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

func (l *ListInstanceTrustedDomainsQuery) checkDomainPerms(ctx context.Context, opts *InvokeOpts) error {
	if authZErr := opts.Permissions.CheckInstancePermission(ctx, DomainReadPermission); authZErr != nil {
		return zerrors.ThrowPermissionDenied(authZErr, "DOM-RyCEyr", "permission denied")
	}
	return nil
}
