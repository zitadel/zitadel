package model

import (
	"encoding/json"
	"time"

	"github.com/caos/logging"

	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/user/model"
	es_model "github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
)

const (
	UserSessionKeySessionID     = "id"
	UserSessionKeyUserAgentID   = "user_agent_id"
	UserSessionKeyUserID        = "user_id"
	UserSessionKeyState         = "session_state"
	UserSessionKeyResourceOwner = "resource_owner"
)

type UserSessionView struct {
	ID            string    `json:"-" gorm:"column:id;primary_key"`
	CreationDate  time.Time `json:"-" gorm:"column:creation_date"`
	ChangeDate    time.Time `json:"-" gorm:"column:change_date"`
	ResourceOwner string    `json:"-" gorm:"column:resource_owner"`
	//State                   int32     `json:"-" gorm:"column:user_state"`
	UserAgentID             string    `json:"userAgentID" gorm:"column:user_agent_id"`
	UserID                  string    `json:"userID" gorm:"column:user_id"`
	UserName                string    `json:"userName" gorm:"column:user_name"`
	PasswordVerification    time.Time `json:"-" gorm:"column:password_verification"`
	MfaSoftwareVerification time.Time `json:"-" gorm:"column:mfa_software_verification"`
	MfaHardwareVerification time.Time `json:"-" gorm:"column:mfa_hardware_verification"`
	Sequence                uint64    `json:"-" gorm:"column:sequence"`
}

func UserSessionFromEvent(event *models.Event) (*UserSessionView, error) {
	v := new(UserSessionView)
	if err := json.Unmarshal(event.Data, v); err != nil {
		logging.Log("EVEN-lso9e").WithError(err).Error("could not unmarshal event data")
		return nil, caos_errs.ThrowInternal(nil, "MODEL-sd325", "could not unmarshal data")
	}
	return v, nil
}

func UserSessionFromModel(userSession *model.UserSessionView) *UserSessionView {
	return &UserSessionView{
		ID:            userSession.ID,
		ChangeDate:    userSession.ChangeDate,
		CreationDate:  userSession.CreationDate,
		ResourceOwner: userSession.ResourceOwner,
		//State:                   int32(userSession.State),
		UserAgentID:             userSession.UserAgentID,
		UserID:                  userSession.UserID,
		UserName:                userSession.UserName,
		PasswordVerification:    userSession.PasswordVerification,
		MfaSoftwareVerification: userSession.MfaSoftwareVerification,
		MfaHardwareVerification: userSession.MfaHardwareVerification,
		Sequence:                userSession.Sequence,
	}
}

func UserSessionToModel(userSession *UserSessionView) *model.UserSessionView {
	return &model.UserSessionView{
		ID:            userSession.ID,
		ChangeDate:    userSession.ChangeDate,
		CreationDate:  userSession.CreationDate,
		ResourceOwner: userSession.ResourceOwner,
		//State:                   model.UserSessionState(userSession.State),
		UserAgentID:             userSession.UserAgentID,
		UserID:                  userSession.UserID,
		UserName:                userSession.UserName,
		PasswordVerification:    userSession.PasswordVerification,
		MfaSoftwareVerification: userSession.MfaSoftwareVerification,
		MfaHardwareVerification: userSession.MfaHardwareVerification,
		Sequence:                userSession.Sequence,
	}
}

func UserSessionsToModel(userSessions []*UserSessionView) []*model.UserSessionView {
	result := make([]*model.UserSessionView, len(userSessions))
	for i, s := range userSessions {
		result[i] = UserSessionToModel(s)
	}
	return result
}

func (v *UserSessionView) AppendEvent(event *models.Event) (err error) {
	v.ChangeDate = event.CreationDate
	//p.
	switch event.Type {
	case es_model.UserPasswordCheckSucceeded:
		v.PasswordVerification = event.CreationDate
	case es_model.UserPasswordCheckFailed,
		es_model.UserPasswordChanged:
		v.PasswordVerification = time.Time{}
	case es_model.MfaOtpCheckSucceeded:
		v.MfaSoftwareVerification = event.CreationDate
	case es_model.MfaOtpCheckFailed,
		es_model.MfaOtpRemoved:
		v.MfaSoftwareVerification = time.Time{}
	}
	return err
}
