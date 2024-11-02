package view

import (
	"github.com/jinzhu/gorm"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/user/model"
	usr_model "github.com/zitadel/zitadel/internal/user/repository/view/model"
	"github.com/zitadel/zitadel/internal/view/repository"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TokenByIDs(db *gorm.DB, table, tokenID, userID, instanceID string) (*usr_model.TokenView, error) {
	token := new(usr_model.TokenView)
	query := repository.PrepareGetByQuery(table,
		&usr_model.TokenSearchQuery{Key: model.TokenSearchKeyTokenID, Method: domain.SearchMethodEquals, Value: tokenID},
		&usr_model.TokenSearchQuery{Key: model.TokenSearchKeyUserID, Method: domain.SearchMethodEquals, Value: userID},
		&usr_model.TokenSearchQuery{Key: model.TokenSearchKeyInstanceID, Method: domain.SearchMethodEquals, Value: instanceID},
	)
	err := query(db, token)
	if zerrors.IsNotFound(err) {
		return nil, zerrors.ThrowNotFound(nil, "VIEW-6ub3p", "Errors.Token.NotFound")
	}
	return token, err
}

func TokensByUserID(db *gorm.DB, table, userID, instanceID string) ([]*usr_model.TokenView, error) {
	tokens := make([]*usr_model.TokenView, 0)
	userIDQuery := &model.TokenSearchQuery{
		Key:    model.TokenSearchKeyUserID,
		Method: domain.SearchMethodEquals,
		Value:  userID,
	}
	instanceIDQuery := &model.TokenSearchQuery{
		Key:    model.TokenSearchKeyInstanceID,
		Method: domain.SearchMethodEquals,
		Value:  instanceID,
	}
	expirationQuery := &model.TokenSearchQuery{
		Key:    model.TokenSearchKeyExpiration,
		Method: domain.SearchMethodGreaterThan,
		Value:  "now()",
	}
	query := repository.PrepareSearchQuery(table, usr_model.TokenSearchRequest{
		Queries: []*model.TokenSearchQuery{userIDQuery, instanceIDQuery, expirationQuery},
	})
	_, err := query(db, &tokens)
	return tokens, err
}
