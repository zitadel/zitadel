package model

import (
	"encoding/json"
	"time"

	"github.com/caos/logging"

	req_model "github.com/caos/zitadel/internal/auth_request/model"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/user/model"
	es_model "github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
)

const (
	UserSessionKeyUserAgentID   = "user_agent_id"
	UserSessionKeyUserID        = "user_id"
	UserSessionKeyState         = "state"
	UserSessionKeyResourceOwner = "resource_owner"
)

type UserSessionView struct {
	CreationDate                 time.Time `json:"-" gorm:"column:creation_date"`
	ChangeDate                   time.Time `json:"-" gorm:"column:change_date"`
	ResourceOwner                string    `json:"-" gorm:"column:resource_owner"`
	State                        int32     `json:"-" gorm:"column:state"`
	UserAgentID                  string    `json:"userAgentID" gorm:"column:user_agent_id;primary_key"`
	UserID                       string    `json:"userID" gorm:"column:user_id;primary_key"`
	UserName                     string    `json:"-" gorm:"column:user_name"`
	LoginName                    string    `json:"-" gorm:"column:login_name"`
	DisplayName                  string    `json:"-" gorm:"column:user_display_name"`
	SelectedIDPConfigID          string    `json:"selectedIDPConfigID" gorm:"column:selected_idp_config_id"`
	PasswordVerification         time.Time `json:"-" gorm:"column:password_verification"`
	PasswordlessVerification     time.Time `json:"-" gorm:"column:passwordless_verification"`
	ExternalLoginVerification    time.Time `json:"-" gorm:"column:external_login_verification"`
	SecondFactorVerification     time.Time `json:"-" gorm:"column:second_factor_verification"`
	SecondFactorVerificationType int32     `json:"-" gorm:"column:second_factor_verification_type"`
	MultiFactorVerification      time.Time `json:"-" gorm:"column:multi_factor_verification"`
	MultiFactorVerificationType  int32     `json:"-" gorm:"column:multi_factor_verification_type"`
	Sequence                     uint64    `json:"-" gorm:"column:sequence"`
}

func UserSessionFromEvent(event *models.Event) (*UserSessionView, error) {
	v := new(UserSessionView)
	if err := json.Unmarshal(event.Data, v); err != nil {
		logging.Log("EVEN-lso9e").WithError(err).Error("could not unmarshal event data")
		return nil, caos_errs.ThrowInternal(nil, "MODEL-sd325", "could not unmarshal data")
	}
	return v, nil
}

func UserSessionToModel(userSession *UserSessionView) *model.UserSessionView {
	return &model.UserSessionView{
		ChangeDate:                   userSession.ChangeDate,
		CreationDate:                 userSession.CreationDate,
		ResourceOwner:                userSession.ResourceOwner,
		State:                        req_model.UserSessionState(userSession.State),
		UserAgentID:                  userSession.UserAgentID,
		UserID:                       userSession.UserID,
		UserName:                     userSession.UserName,
		LoginName:                    userSession.LoginName,
		DisplayName:                  userSession.DisplayName,
		SelectedIDPConfigID:          userSession.SelectedIDPConfigID,
		PasswordVerification:         userSession.PasswordVerification,
		PasswordlessVerification:     userSession.PasswordlessVerification,
		ExternalLoginVerification:    userSession.ExternalLoginVerification,
		SecondFactorVerification:     userSession.SecondFactorVerification,
		SecondFactorVerificationType: req_model.MFAType(userSession.SecondFactorVerificationType),
		MultiFactorVerification:      userSession.MultiFactorVerification,
		MultiFactorVerificationType:  req_model.MFAType(userSession.MultiFactorVerificationType),
		Sequence:                     userSession.Sequence,
	}
}

func UserSessionsToModel(userSessions []*UserSessionView) []*model.UserSessionView {
	result := make([]*model.UserSessionView, len(userSessions))
	for i, s := range userSessions {
		result[i] = UserSessionToModel(s)
	}
	return result
}

func (v *UserSessionView) AppendEvent(event *models.Event) error {
	v.Sequence = event.Sequence
	v.ChangeDate = event.CreationDate
	switch event.Type {
	case es_model.UserPasswordCheckSucceeded,
		es_model.HumanPasswordCheckSucceeded:
		v.PasswordVerification = event.CreationDate
		v.State = int32(req_model.UserSessionStateActive)
	case es_model.HumanExternalLoginCheckSucceeded:
		data := new(es_model.AuthRequest)
		err := data.SetData(event)
		if err != nil {
			return err
		}
		v.ExternalLoginVerification = event.CreationDate
		v.SelectedIDPConfigID = data.SelectedIDPConfigID
		v.State = int32(req_model.UserSessionStateActive)
	case es_model.HumanPasswordlessTokenCheckSucceeded:
		v.PasswordlessVerification = event.CreationDate
		v.MultiFactorVerification = event.CreationDate
		v.MultiFactorVerificationType = int32(req_model.MFATypeU2FUserVerification)
		v.State = int32(req_model.UserSessionStateActive)
	case es_model.HumanPasswordlessTokenCheckFailed,
		es_model.HumanPasswordlessTokenRemoved:
		v.PasswordlessVerification = time.Time{}
		v.MultiFactorVerification = time.Time{}
		v.State = int32(req_model.UserSessionStateInitiated)
	case es_model.UserPasswordCheckFailed,
		es_model.HumanPasswordCheckFailed:
		v.PasswordVerification = time.Time{}
		v.State = int32(req_model.UserSessionStateInitiated)
	case es_model.UserPasswordChanged,
		es_model.HumanPasswordChanged:
		data := new(es_model.PasswordChange)
		err := data.SetData(event)
		if err != nil {
			return err
		}
		if v.UserAgentID != data.UserAgentID {
			v.PasswordVerification = time.Time{}
			v.State = int32(req_model.UserSessionStateInitiated)
		}
	case es_model.MFAOTPVerified,
		es_model.HumanMFAOTPVerified:
		data := new(es_model.OTPVerified)
		err := data.SetData(event)
		if err != nil {
			return err
		}
		if v.UserAgentID == data.UserAgentID {
			v.setSecondFactorVerification(event.CreationDate, req_model.MFATypeOTP)
		}
	case es_model.MFAOTPCheckSucceeded,
		es_model.HumanMFAOTPCheckSucceeded:
		v.setSecondFactorVerification(event.CreationDate, req_model.MFATypeOTP)
	case es_model.MFAOTPCheckFailed,
		es_model.MFAOTPRemoved,
		es_model.HumanMFAOTPCheckFailed,
		es_model.HumanMFAOTPRemoved,
		es_model.HumanMFAU2FTokenCheckFailed,
		es_model.HumanMFAU2FTokenRemoved:
		v.SecondFactorVerification = time.Time{}
		v.State = int32(req_model.UserSessionStateInitiated)
	case es_model.HumanMFAU2FTokenVerified:
		data := new(es_model.WebAuthNVerify)
		err := data.SetData(event)
		if err != nil {
			return err
		}
		if v.UserAgentID == data.UserAgentID {
			v.setSecondFactorVerification(event.CreationDate, req_model.MFATypeU2F)
		}
	case es_model.HumanMFAU2FTokenCheckSucceeded:
		v.setSecondFactorVerification(event.CreationDate, req_model.MFATypeU2F)
	case es_model.SignedOut,
		es_model.HumanSignedOut,
		es_model.UserLocked,
		es_model.UserDeactivated,
		es_model.UserRemoved:
		v.PasswordlessVerification = time.Time{}
		v.PasswordVerification = time.Time{}
		v.SecondFactorVerification = time.Time{}
		v.SecondFactorVerificationType = int32(req_model.MFALevelNotSetUp)
		v.MultiFactorVerification = time.Time{}
		v.MultiFactorVerificationType = int32(req_model.MFALevelNotSetUp)
		v.ExternalLoginVerification = time.Time{}
		v.State = int32(req_model.UserSessionStateTerminated)
	case es_model.HumanExternalIDPRemoved,
		es_model.HumanExternalIDPCascadeRemoved:
		v.ExternalLoginVerification = time.Time{}
		v.SelectedIDPConfigID = ""
		v.State = int32(req_model.UserSessionStateTerminated)
	}
	return nil
}

func (v *UserSessionView) setSecondFactorVerification(verificationTime time.Time, mfaType req_model.MFAType) {
	v.SecondFactorVerification = verificationTime
	v.SecondFactorVerificationType = int32(mfaType)
	v.State = int32(req_model.UserSessionStateActive)
}
