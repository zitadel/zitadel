package view

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/user/repository/view/model"
	"github.com/zitadel/zitadel/internal/zerrors"
)

//go:embed user_by_id.sql
var userByIDQuery string

func UserByID(ctx context.Context, db *gorm.DB, userID, instanceID string) (*model.UserView, error) {
	user := new(model.UserView)

	query := db.Raw(userByIDQuery, instanceID, userID)

	tx := query.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	defer func() {
		if err := tx.Commit().Error; err != nil {
			logging.OnError(err).Info("commit failed")
		}
		tx.RollbackUnlessCommitted()
	}()

	err := tx.Scan(user).Error
	if err == nil {
		user.SetEmptyUserType()
		return user, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, zerrors.ThrowNotFound(err, "VIEW-hodc6", "Errors.User.NotFound")
	}
	logging.WithError(err).Warn("unable to get user by id")
	return nil, zerrors.ThrowInternal(err, "VIEW-qJBg9", "unable to get user by id")
}
