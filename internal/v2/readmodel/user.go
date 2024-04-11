package readmodel

import (
	"time"

	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/v2/user"
)

type User struct {
	ID                 string
	CreationDate       time.Time
	ChangeDate         time.Time
	ResourceOwner      string
	Sequence           uint64
	State              domain.UserState
	Type               domain.UserType
	Username           string
	LoginNames         database.TextArray[string]
	PreferredLoginName string
	Human              *Human
	Machine            *Machine
}

func NewUser(id string) *User {
	return &User{
		ID: id,
	}
}

func (rm *User) Filter() []*eventstore.Filter {
	return []*eventstore.Filter{
		eventstore.NewFilter(
			eventstore.AppendAggregateFilter(
				user.AggregateType,
				eventstore.AggregateID(rm.ID),
			),
		),
	}
}

func (rm *User) Reduce(events ...*eventstore.Event[eventstore.StoragePayload]) error {
	for _, event := range events {
		switch event.Type {
		case "user.human.added", "user.added":
			added, err := user.A
			if err != nil {
				return err
			}
			rm.Name = added.Payload.Name
			rm.Owner = event.Aggregate.Owner
			rm.CreationDate = event.CreatedAt
		}
	}
}

type Human struct {
	FirstName              string
	LastName               string
	NickName               string
	DisplayName            string
	AvatarKey              string
	PreferredLanguage      language.Tag
	Gender                 domain.Gender
	Email                  domain.EmailAddress
	IsEmailVerified        bool
	Phone                  domain.PhoneNumber
	IsPhoneVerified        bool
	PasswordChangeRequired bool
}

type Machine struct {
	Name            string
	Description     string
	EncodedSecret   string
	AccessTokenType domain.OIDCTokenType
}
