package projection

import "github.com/zitadel/zitadel/internal/eventstore"

var _ Projection = (*LoginNames)(nil)

type LoginNames struct {
	LoginNames []*LoginName
	userID     string
	instanceID string
}

func NewLoginNames(userID, instanceID string) *LoginNames {
	return &LoginNames{
		userID:     userID,
		instanceID: instanceID,
	}
}

func (u *LoginNames) Reduce(events []eventstore.Event) {}

func (u *LoginNames) SearchQuery() *eventstore.SearchQueryBuilder {
	return nil
}

type LoginName struct {
	Name      string
	IsPrimary bool
}
