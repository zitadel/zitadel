package stmt

import "github.com/zitadel/zitadel/backend/v3/domain"

func UserIDCondition(id string) *TextCondition[string, domain.User] {
	return &TextCondition[string, domain.User]{
		condition: condition[string, domain.User, TextOperation]{
			field: userColumns[UserColumnID],
			op:    TextOperationEqual,
			value: id,
		},
	}
}

func UserUsernameCondition(op TextOperation, username string) *TextCondition[string, domain.User] {
	return &TextCondition[string, domain.User]{
		condition: condition[string, domain.User, TextOperation]{
			field: userColumns[UserColumnUsername],
			op:    op,
			value: username,
		},
	}
}
