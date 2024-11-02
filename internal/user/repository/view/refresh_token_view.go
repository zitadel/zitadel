package view

import (
	"github.com/jinzhu/gorm"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/user/model"
	usr_model "github.com/zitadel/zitadel/internal/user/repository/view/model"
	"github.com/zitadel/zitadel/internal/view/repository"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func RefreshTokenByID(db *gorm.DB, table, tokenID, instanceID string) (*usr_model.RefreshTokenView, error) {
	token := new(usr_model.RefreshTokenView)
	query := repository.PrepareGetByQuery(table,
		&usr_model.RefreshTokenSearchQuery{Key: model.RefreshTokenSearchKeyRefreshTokenID, Method: domain.SearchMethodEquals, Value: tokenID},
		&usr_model.RefreshTokenSearchQuery{Key: model.RefreshTokenSearchKeyInstanceID, Method: domain.SearchMethodEquals, Value: instanceID},
	)
	err := query(db, token)
	if zerrors.IsNotFound(err) {
		return nil, zerrors.ThrowNotFound(nil, "VIEW-6ub3p", "Errors.RefreshToken.NotFound")
	}
	return token, err
}

func RefreshTokensByUserID(db *gorm.DB, table, userID, instanceID string) ([]*usr_model.RefreshTokenView, error) {
	tokens := make([]*usr_model.RefreshTokenView, 0)
	userIDQuery := &model.RefreshTokenSearchQuery{
		Key:    model.RefreshTokenSearchKeyUserID,
		Method: domain.SearchMethodEquals,
		Value:  userID,
	}
	instanceIDQuery := &model.RefreshTokenSearchQuery{
		Key:    model.RefreshTokenSearchKeyInstanceID,
		Method: domain.SearchMethodEquals,
		Value:  instanceID,
	}
	query := repository.PrepareSearchQuery(table, usr_model.RefreshTokenSearchRequest{
		Queries: []*model.RefreshTokenSearchQuery{userIDQuery, instanceIDQuery},
	})
	_, err := query(db, &tokens)
	return tokens, err
}

func SearchRefreshTokens(db *gorm.DB, table string, req *model.RefreshTokenSearchRequest) ([]*usr_model.RefreshTokenView, uint64, error) {
	tokens := make([]*usr_model.RefreshTokenView, 0)
	query := repository.PrepareSearchQuery(table, usr_model.RefreshTokenSearchRequest{Limit: req.Limit, Offset: req.Offset, Queries: req.Queries})
	count, err := query(db, &tokens)
	return tokens, count, err
}
