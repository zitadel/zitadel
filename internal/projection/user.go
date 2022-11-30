package projection

import "github.com/zitadel/zitadel/internal/eventstore"

var _ Projection = (*User)(nil)

type User struct{}

func (u *User) Reduce(events []eventstore.Event) {}

func (u *User) SearchQuery() *eventstore.SearchQueryBuilder {
	return nil
}
