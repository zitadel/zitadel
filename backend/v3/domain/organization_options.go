package domain

import "github.com/zitadel/zitadel/backend/v3/storage/database"

var _ CreateOrganizationCommandOpts = (*withOrganizationID)(nil)

type withOrganizationID struct {
	id string
}

func WithOrganizationID(id string) *withOrganizationID {
	return &withOrganizationID{
		id: id,
	}
}

func (opt *withOrganizationID) applyOnCreateOrganizationCommand(cmd *CreateOrganizationCommand) {
	cmd.ID = opt.id
}

var _ UpdateOrganizationCommandOpts = (*withOrganizationName)(nil)

type withOrganizationName struct {
	name string
}

func WithOrganizationName(name string) *withOrganizationName {
	return &withOrganizationName{
		name: name,
	}
}

func (opt *withOrganizationName) applyOnUpdateOrganizationCommand(cmd *UpdateOrganizationCommand) {
	cmd.changes = append(cmd.changes, cmd.repo.SetName(opt.name))
}

var _ OrgsQueryOpts = (*orgByNameQueryOpt)(nil)

type orgByNameQueryOpt struct {
	name string
	op   database.TextOperation
}

func WithOrgByNameQuery(op database.TextOperation, name string) *orgByNameQueryOpt {
	return &orgByNameQueryOpt{
		name: name,
		op:   op,
	}
}

func (opt *orgByNameQueryOpt) applyOnOrgsQuery(query *OrgsQuery) {
	query.conditions = append(query.conditions, query.repo.NameCondition(opt.op, opt.name))
}

var _ OrgsQueryOpts = (*orgByDomainQueryOpt)(nil)

type orgByDomainQueryOpt struct {
	name string
	op   database.TextOperation
}

func WithOrgByDomainQuery(op database.TextOperation, name string) *orgByDomainQueryOpt {
	return &orgByDomainQueryOpt{
		name: name,
		op:   op,
	}
}

func (opt *orgByDomainQueryOpt) applyOnOrgsQuery(query *OrgsQuery) {
	query.conditions = append(query.conditions, query.domainRepo.DomainCondition(opt.op, opt.name))
}

var _ OrgsQueryOpts = (*orgByIDQueryOpt)(nil)

type orgByIDQueryOpt struct {
	id string
}

func WithOrgByIDQuery(id string) *orgByIDQueryOpt {
	return &orgByIDQueryOpt{
		id: id,
	}
}

func (opt *orgByIDQueryOpt) applyOnOrgsQuery(query *OrgsQuery) {
	query.conditions = append(query.conditions, query.repo.IDCondition(opt.id))
}

var _ OrgsQueryOpts = (*orgByIDQueryOpt)(nil)

type orgByStateQueryOpt struct {
	state OrgState
}

func WithOrgByStateQuery(state OrgState) *orgByStateQueryOpt {
	return &orgByStateQueryOpt{
		state: state,
	}
}

func (opt *orgByStateQueryOpt) applyOnOrgsQuery(query *OrgsQuery) {
	query.conditions = append(query.conditions, query.repo.StateCondition(opt.state))
}

var _ OrgsQueryOpts = (*orgByIDQueryOpt)(nil)

type orgQuerySortingColumnOpt struct {
	getColumn func(query *OrgsQuery) database.Column
}

func WithOrgQuerySortingColumn(getColumn func(query *OrgsQuery) database.Column) *orgQuerySortingColumnOpt {
	return &orgQuerySortingColumnOpt{
		getColumn: getColumn,
	}
}

func OrderOrgsByCreationDate(query *OrgsQuery) database.Column {
	return query.repo.CreatedAtColumn(true)
}

func OrderOrgsByName(query *OrgsQuery) database.Column {
	return query.repo.NameColumn(true)
}

func (opt *orgQuerySortingColumnOpt) applyOnOrgsQuery(query *OrgsQuery) {
	query.pagination.OrderColumns = append(query.pagination.OrderColumns, opt.getColumn(query))
}
