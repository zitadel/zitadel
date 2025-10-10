package domain

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
	instance "github.com/zitadel/zitadel/pkg/grpc/instance/v2beta"
)

type ListInstancesCommand struct {
	BaseCommand
	Request *instance.ListInstancesRequest

	Result []*Instance
}

var _ Commander = (*ListInstancesCommand)(nil)

func NewListInstancesCommand(inputRequest *instance.ListInstancesRequest) *ListInstancesCommand {
	return &ListInstancesCommand{
		BaseCommand: BaseCommand{},
		Request:     inputRequest,
	}
}

// Events implements Commander.
func (l *ListInstancesCommand) Events(ctx context.Context, opts *CommandOpts) ([]eventstore.Command, error) {
	return nil, nil
}

// Execute implements Commander.
func (l *ListInstancesCommand) Execute(ctx context.Context, opts *CommandOpts) (err error) {
	close, err := opts.EnsureTx(ctx)
	if err != nil {
		return err
	}
	defer func() { err = close(ctx, err) }()

	instanceRepo := opts.instanceRepo.LoadDomains()
	domainRepo := opts.instanceDomainRepo

	sorting := l.Sorting(instanceRepo)
	limit, offset := l.Pagination(l.Request.GetPagination().GetLimit(), l.Request.GetPagination().GetOffset())
	conds, err := l.conditions(instanceRepo, domainRepo)
	if err != nil {
		return err
	}

	instances, err := instanceRepo.List(ctx, pool, conds, sorting, limit, offset)
	if err != nil {
		return err
	}

	l.Result = instances

	return nil
}

func (l *ListInstancesCommand) conditions(instanceRepo InstanceRepository, domainRepo InstanceDomainRepository) (database.QueryOption, error) {
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

func (l *ListInstancesCommand) Sorting(instanceRepo InstanceRepository) database.QueryOption {
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

// String implements Commander.
func (l *ListInstancesCommand) String() string {
	return "ListInstancesCommand"
}

// Validate implements Commander.
func (l *ListInstancesCommand) Validate(ctx context.Context, opts *CommandOpts) (err error) {
	return nil
}

func (l *ListInstancesCommand) ResultToGRPC() []*instance.Instance {
	toReturn := make([]*instance.Instance, len(l.Result))

	for i, inst := range l.Result {
		toReturn[i] = &instance.Instance{
			Id:           inst.ID,
			ChangeDate:   timestamppb.New(inst.UpdatedAt),
			CreationDate: timestamppb.New(inst.CreatedAt),
			State:        instance.State_STATE_RUNNING,
			Name:         inst.Name,
			Domains:      l.domainsToGRPC(inst.Domains),
		}
	}
	return toReturn
}

func (l *ListInstancesCommand) domainsToGRPC(domains []*InstanceDomain) []*instance.Domain {
	toReturn := make([]*instance.Domain, len(domains))
	for i, domain := range domains {
		isPrimary := domain.IsPrimary != nil && *domain.IsPrimary
		isGenerated := domain.IsGenerated != nil && *domain.IsGenerated
		toReturn[i] = &instance.Domain{
			InstanceId:   domain.InstanceID,
			CreationDate: timestamppb.New(domain.CreatedAt),
			Domain:       domain.Domain,
			Primary:      isPrimary,
			Generated:    isGenerated,
		}
	}

	return toReturn
}
