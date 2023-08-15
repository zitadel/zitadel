package model

import (
	"encoding/json"
	"time"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/user/model"
	es_model "github.com/zitadel/zitadel/internal/user/repository/eventsourcing/model"
)

const (
	UserSessionKeyUserAgentID   = "user_agent_id"
	UserSessionKeyUserID        = "user_id"
	UserSessionKeyState         = "state"
	UserSessionKeyResourceOwner = "resource_owner"
	UserSessionKeyInstanceID    = "instance_id"
	UserSessionKeyOwnerRemoved  = "owner_removed"
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
	AvatarKey                    string    `json:"-" gorm:"column:avatar_key"`
	SelectedIDPConfigID          string    `json:"selectedIDPConfigID" gorm:"column:selected_idp_config_id"`
	PasswordVerification         time.Time `json:"-" gorm:"column:password_verification"`
	PasswordlessVerification     time.Time `json:"-" gorm:"column:passwordless_verification"`
	ExternalLoginVerification    time.Time `json:"-" gorm:"column:external_login_verification"`
	SecondFactorVerification     time.Time `json:"-" gorm:"column:second_factor_verification"`
	SecondFactorVerificationType int32     `json:"-" gorm:"column:second_factor_verification_type"`
	MultiFactorVerification      time.Time `json:"-" gorm:"column:multi_factor_verification"`
	MultiFactorVerificationType  int32     `json:"-" gorm:"column:multi_factor_verification_type"`
	Sequence                     uint64    `json:"-" gorm:"column:sequence"`
	InstanceID                   string    `json:"instanceID" gorm:"column:instance_id;primary_key"`
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
		State:                        domain.UserSessionState(userSession.State),
		UserAgentID:                  userSession.UserAgentID,
		UserID:                       userSession.UserID,
		UserName:                     userSession.UserName,
		LoginName:                    userSession.LoginName,
		DisplayName:                  userSession.DisplayName,
		AvatarKey:                    userSession.AvatarKey,
		SelectedIDPConfigID:          userSession.SelectedIDPConfigID,
		PasswordVerification:         userSession.PasswordVerification,
		PasswordlessVerification:     userSession.PasswordlessVerification,
		ExternalLoginVerification:    userSession.ExternalLoginVerification,
		SecondFactorVerification:     userSession.SecondFactorVerification,
		SecondFactorVerificationType: domain.MFAType(userSession.SecondFactorVerificationType),
		MultiFactorVerification:      userSession.MultiFactorVerification,
		MultiFactorVerificationType:  domain.MFAType(userSession.MultiFactorVerificationType),
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
	switch eventstore.EventType(event.Type) {
	case user.UserV1PasswordCheckSucceededType,
		user.HumanPasswordCheckSucceededType:
		v.PasswordVerification = event.CreationDate
		v.State = int32(domain.UserSessionStateActive)
	case user.UserIDPLoginCheckSucceededType:
		data := new(es_model.AuthRequest)
		err := data.SetData(event)
		if err != nil {
			return err
		}
		v.ExternalLoginVerification = event.CreationDate
		v.SelectedIDPConfigID = data.SelectedIDPConfigID
		v.State = int32(domain.UserSessionStateActive)
	case user.HumanPasswordlessTokenCheckSucceededType:
		v.PasswordlessVerification = event.CreationDate
		v.MultiFactorVerification = event.CreationDate
		v.MultiFactorVerificationType = int32(domain.MFATypeU2FUserVerification)
		v.State = int32(domain.UserSessionStateActive)
	case user.HumanPasswordlessTokenCheckFailedType,
		user.HumanPasswordlessTokenRemovedType:
		v.PasswordlessVerification = time.Time{}
		v.MultiFactorVerification = time.Time{}
	case user.UserV1PasswordCheckFailedType,
		user.HumanPasswordCheckFailedType:
		v.PasswordVerification = time.Time{}
	case user.UserV1PasswordChangedType,
		user.HumanPasswordChangedType:
		data := new(es_model.PasswordChange)
		err := data.SetData(event)
		if err != nil {
			return err
		}
		if v.UserAgentID != data.UserAgentID {
			v.PasswordVerification = time.Time{}
		}
	case user.HumanMFAOTPVerifiedType:
		data := new(es_model.OTPVerified)
		err := data.SetData(event)
		if err != nil {
			return err
		}
		if v.UserAgentID == data.UserAgentID {
			v.setSecondFactorVerification(event.CreationDate, domain.MFATypeTOTP)
		}
	case user.UserV1MFAOTPCheckSucceededType,
		user.HumanMFAOTPCheckSucceededType:
		v.setSecondFactorVerification(event.CreationDate, domain.MFATypeTOTP)
	case user.HumanOTPSMSCheckSucceededType:
		data := new(es_model.OTPVerified)
		err := data.SetData(event)
		if err != nil {
			return err
		}
		if v.UserAgentID == data.UserAgentID {
			v.setSecondFactorVerification(event.CreationDate, domain.MFATypeOTPSMS)
		}
	case user.HumanOTPEmailCheckSucceededType:
		data := new(es_model.OTPVerified)
		err := data.SetData(event)
		if err != nil {
			return err
		}
		if v.UserAgentID == data.UserAgentID {
			v.setSecondFactorVerification(event.CreationDate, domain.MFATypeOTPEmail)
		}
	case user.UserV1MFAOTPCheckFailedType,
		user.UserV1MFAOTPRemovedType,
		user.HumanMFAOTPCheckFailedType,
		user.HumanMFAOTPRemovedType,
		user.HumanU2FTokenCheckFailedType,
		user.HumanU2FTokenRemovedType,
		user.HumanOTPSMSCheckFailedType,
		user.HumanOTPEmailCheckFailedType:
		v.SecondFactorVerification = time.Time{}
	case user.HumanU2FTokenVerifiedType:
		data := new(es_model.WebAuthNVerify)
		err := data.SetData(event)
		if err != nil {
			return err
		}
		if v.UserAgentID == data.UserAgentID {
			v.setSecondFactorVerification(event.CreationDate, domain.MFATypeU2F)
		}
	case user.HumanU2FTokenCheckSucceededType:
		v.setSecondFactorVerification(event.CreationDate, domain.MFATypeU2F)
	case user.UserV1SignedOutType,
		user.HumanSignedOutType,
		user.UserLockedType,
		user.UserDeactivatedType:
		v.PasswordlessVerification = time.Time{}
		v.PasswordVerification = time.Time{}
		v.SecondFactorVerification = time.Time{}
		v.SecondFactorVerificationType = int32(domain.MFALevelNotSetUp)
		v.MultiFactorVerification = time.Time{}
		v.MultiFactorVerificationType = int32(domain.MFALevelNotSetUp)
		v.ExternalLoginVerification = time.Time{}
		v.State = int32(domain.UserSessionStateTerminated)
	case user.UserIDPLinkRemovedType, user.UserIDPLinkCascadeRemovedType:
		v.ExternalLoginVerification = time.Time{}
		v.SelectedIDPConfigID = ""
	case user.HumanAvatarAddedType:
		key, err := avatarKeyFromEvent(event)
		if err != nil {
			return err
		}
		v.AvatarKey = key
	case user.HumanAvatarRemovedType:
		v.AvatarKey = ""
	}
	return nil
}

func (v *UserSessionView) setSecondFactorVerification(verificationTime time.Time, mfaType domain.MFAType) {
	v.SecondFactorVerification = verificationTime
	v.SecondFactorVerificationType = int32(mfaType)
	v.State = int32(domain.UserSessionStateActive)
}

func avatarKeyFromEvent(event *models.Event) (string, error) {
	data := make(map[string]string)
	if err := json.Unmarshal(event.Data, &data); err != nil {
		logging.Log("EVEN-Sfew2").WithError(err).Error("could not unmarshal event data")
		return "", caos_errs.ThrowInternal(err, "MODEL-SFw2q", "could not unmarshal event")
	}
	return data["storeKey"], nil
}

func (v *UserSessionView) EventTypes() []models.EventType {
	return []models.EventType{
		models.EventType(user.UserV1PasswordCheckSucceededType),
		models.EventType(user.HumanPasswordCheckSucceededType),
		models.EventType(user.UserIDPLoginCheckSucceededType),
		models.EventType(user.HumanPasswordlessTokenCheckSucceededType),
		models.EventType(user.HumanPasswordlessTokenCheckFailedType),
		models.EventType(user.HumanPasswordlessTokenRemovedType),
		models.EventType(user.UserV1PasswordCheckFailedType),
		models.EventType(user.HumanPasswordCheckFailedType),
		models.EventType(user.UserV1PasswordChangedType),
		models.EventType(user.HumanPasswordChangedType),
		models.EventType(user.HumanMFAOTPVerifiedType),
		models.EventType(user.UserV1MFAOTPCheckSucceededType),
		models.EventType(user.HumanMFAOTPCheckSucceededType),
		models.EventType(user.UserV1MFAOTPCheckFailedType),
		models.EventType(user.UserV1MFAOTPRemovedType),
		models.EventType(user.HumanMFAOTPCheckFailedType),
		models.EventType(user.HumanMFAOTPRemovedType),
		models.EventType(user.HumanOTPSMSCheckSucceededType),
		models.EventType(user.HumanOTPSMSCheckFailedType),
		models.EventType(user.HumanOTPEmailCheckSucceededType),
		models.EventType(user.HumanOTPEmailCheckFailedType),
		models.EventType(user.HumanU2FTokenCheckFailedType),
		models.EventType(user.HumanU2FTokenRemovedType),
		models.EventType(user.HumanU2FTokenVerifiedType),
		models.EventType(user.HumanU2FTokenCheckSucceededType),
		models.EventType(user.UserV1SignedOutType),
		models.EventType(user.HumanSignedOutType),
		models.EventType(user.UserLockedType),
		models.EventType(user.UserDeactivatedType),
		models.EventType(user.UserIDPLinkRemovedType),
		models.EventType(user.UserIDPLinkCascadeRemovedType),
		models.EventType(user.HumanAvatarAddedType),
		models.EventType(user.HumanAvatarRemovedType),
	}
}
