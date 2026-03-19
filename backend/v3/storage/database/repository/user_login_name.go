package repository

import "github.com/zitadel/zitadel/backend/v3/storage/database"

type userLoginName struct{}

func (u userLoginName) qualifiedTableName() string {
	return "zitadel.login_names"
}

func (u userLoginName) unqualifiedTableName() string {
	return "login_names"
}

// -------------------------------------------------------------
// conditions
// -------------------------------------------------------------

// LoginNameCondition implements [domain.UserRepository].
func (u userLoginName) LoginNameCondition(op database.TextOperation, loginName string) database.Condition {
	return database.NewTextCondition(u.loginNameColumn(), op, loginName)
}

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

func (u userLoginName) loginNameColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "login_name")
}

func (u userLoginName) instanceIDColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "instance_id")
}

func (u userLoginName) userIDColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "user_id")
}
