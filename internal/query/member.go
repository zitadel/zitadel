package query

import (
	"time"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"

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
	return NewTextQuery(membershipUserID, value, TextEquals)
}
func NewMemberResourceOwnerSearchQuery(value string) (SearchQuery, error) {
	return NewTextQuery(membershipResourceOwner, value, TextEquals)
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
	Roles              database.TextArray[string]
	PreferredLoginName string
	Email              string
	FirstName          string
	LastName           string
	DisplayName        string
	AvatarURL          string
	UserType           domain.UserType
}
