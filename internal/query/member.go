package query

import (
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
)

type MembersQuery struct {
	SearchRequest
	Queries []SearchQuery
}

func (q *MembersQuery) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	query = q.SearchRequest.toQuery(query)
	for _, q := range q.Queries {
		query = q.toQuery(query)
	}
	return query
}

func NewMemberEmailSearchQuery(method TextComparison, value string) (SearchQuery, error) {
	return NewTextQuery(HumanEmailCol, value, method)
}

func NewMemberFirstNameSearchQuery(method TextComparison, value string) (SearchQuery, error) {
	return NewTextQuery(HumanFirstNameCol, value, method)
}

func NewMemberLastNameSearchQuery(method TextComparison, value string) (SearchQuery, error) {
	return NewTextQuery(HumanLastNameCol, value, method)
}

func NewMemberUserIDSearchQuery(value string) (SearchQuery, error) {
	return NewTextQuery(MembershipUserID, value, TextEquals)
}

func NewMemberInUserIDsSearchQuery(ids []string) (SearchQuery, error) {
	list := make([]interface{}, len(ids))
	for i, value := range ids {
		list[i] = value
	}
	return NewListQuery(MembershipUserID, list, ListIn)
}

func NewMemberResourceOwnerSearchQuery(value string) (SearchQuery, error) {
	return NewTextQuery(membershipResourceOwner, value, TextEquals)
}

type Members struct {
	SearchResponse
	Members []*Member
}

type Member struct {
	CreationDate       time.Time
	ChangeDate         time.Time
	Sequence           uint64
	ResourceOwner      string
	UserResourceOwner  string
	UserID             string
	Roles              database.TextArray[string]
	PreferredLoginName string
	Email              string
	FirstName          string
	LastName           string
	DisplayName        string
	AvatarURL          string
	UserType           domain.UserType
}

func prepareInstanceMembers() string {
	builder := sq.Select(
		InstanceMemberCreationDate.identifier(),
		InstanceMemberChangeDate.identifier(),
		InstanceMemberSequence.identifier(),
		InstanceMemberResourceOwner.identifier(),
		InstanceMemberUserResourceOwner.identifier(),
		InstanceMemberUserID.identifier(),
		InstanceMemberRoles.identifier(),
		LoginNameNameCol.identifier(),
		HumanEmailCol.identifier(),
		HumanFirstNameCol.identifier(),
		HumanLastNameCol.identifier(),
		HumanDisplayNameCol.identifier(),
		MachineNameCol.identifier(),
		HumanAvatarURLCol.identifier(),
		UserTypeCol.identifier(),
		countColumn.identifier(),
	).From(instanceMemberTable.identifier()).
		LeftJoin(join(HumanUserIDCol, InstanceMemberUserID)).
		LeftJoin(join(MachineUserIDCol, InstanceMemberUserID)).
		LeftJoin(join(UserIDCol, InstanceMemberUserID)).
		LeftJoin(join(LoginNameUserIDCol, InstanceMemberUserID)).
		Where(
			sq.Eq{LoginNameIsPrimaryCol.identifier(): true},
		).PlaceholderFormat(sq.Dollar)

	stmt, _ := builder.MustSql()
	return stmt
}

func prepareOrgMembers() string {
	builder := sq.Select(
		OrgMemberCreationDate.identifier(),
		OrgMemberChangeDate.identifier(),
		OrgMemberSequence.identifier(),
		OrgMemberResourceOwner.identifier(),
		OrgMemberUserResourceOwner.identifier(),
		OrgMemberUserID.identifier(),
		OrgMemberRoles.identifier(),
		LoginNameNameCol.identifier(),
		HumanEmailCol.identifier(),
		HumanFirstNameCol.identifier(),
		HumanLastNameCol.identifier(),
		HumanDisplayNameCol.identifier(),
		MachineNameCol.identifier(),
		HumanAvatarURLCol.identifier(),
		UserTypeCol.identifier(),
		countColumn.identifier(),
	).From(orgMemberTable.identifier()).
		LeftJoin(join(HumanUserIDCol, OrgMemberUserID)).
		LeftJoin(join(MachineUserIDCol, OrgMemberUserID)).
		LeftJoin(join(UserIDCol, OrgMemberUserID)).
		LeftJoin(join(LoginNameUserIDCol, OrgMemberUserID)).
		Where(
			sq.Eq{LoginNameIsPrimaryCol.identifier(): true},
		).PlaceholderFormat(sq.Dollar)

	stmt, _ := builder.MustSql()
	return stmt
}
