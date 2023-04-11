package query

import (
	"time"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query/projection"

	sq "github.com/Masterminds/squirrel"
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
	return NewTextQuery(memberUserID, value, TextEquals)
}
func NewMemberResourceOwnerSearchQuery(value string) (SearchQuery, error) {
	return NewTextQuery(memberResourceOwner, value, TextEquals)
}

type Members struct {
	SearchResponse
	Members []*Member
}

type Member struct {
	CreationDate  time.Time
	ChangeDate    time.Time
	Sequence      uint64
	ResourceOwner string

	UserID             string
	Roles              database.StringArray
	PreferredLoginName string
	Email              string
	FirstName          string
	LastName           string
	DisplayName        string
	AvatarURL          string
	UserType           domain.UserType
}

var (
	memberTableAlias = table{
		name:          "members",
		alias:         "members",
		instanceIDCol: projection.MemberInstanceID,
	}
	memberUserID = Column{
		name:  projection.MemberUserIDCol,
		table: memberTableAlias,
	}
	memberResourceOwner = Column{
		name:  projection.MemberResourceOwner,
		table: memberTableAlias,
	}
)
