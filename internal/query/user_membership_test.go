package query

import (
	"fmt"
	"testing"
)

func Test_prepare(t *testing.T) {
	b, _ := prepareMembershipsQuery()
	ro, _ := NewMembershipResourceOwnerQuery("ro")
	usr, _ := NewMembershipUserIDQuery("usr")
	q := &MembershipSearchQuery{
		SearchRequest: SearchRequest{
			Offset: 100,
			Limit:  100,
		},
		Queries: []SearchQuery{ro, usr},
	}

	stmt, args, err := q.toQuery(b).ToSql()
	fmt.Println(stmt, args, err)
}
