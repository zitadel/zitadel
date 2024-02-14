package model

import (
	"time"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/user/model"
	es_model "github.com/zitadel/zitadel/internal/user/repository/eventsourcing/model"
	"github.com/zitadel/zitadel/internal/zerrors"
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
	CreationDate  time.Time `json:"-" gorm:"column:creation_date"`
	ChangeDate    time.Time `json:"-" gorm:"column:change_date"`
	ResourceOwner string    `json:"-" gorm:"column:resource_owner"`
	State         int32     `json:"-" gorm:"column:state"`
	UserAgentID   string    `json:"userAgentID" gorm:"column:user_agent_id;primary_key"`
	UserID        string    `json:"userID" gorm:"column:user_id;primary_key"`
	// As of https://github.com/zitadel/zitadel/pull/7199 the following 4 attributes
	// are not projected in the user session handler anymore
	// and are therefore annotated with a `gorm:"-"`.
	// They will be read from the corresponding projection directly.
	UserName                     string    `json:"-" gorm:"-"`
	LoginName                    string    `json:"-" gorm:"-"`
	DisplayName                  string    `json:"-" gorm:"-"`
	AvatarKey                    string    `json:"-" gorm:"-"`
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

func UserSessionFromEvent(event eventstore.Event) (*UserSessionView, error) {
	v := new(UserSessionView)
	if err := event.Unmarshal(v); err != nil {
		logging.Log("EVEN-lso9e").WithError(err).Error("could not unmarshal event data")
		return nil, zerrors.ThrowInternal(nil, "MODEL-sd325", "could not unmarshal data")
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

func (v *UserSessionView) AppendEvent(event eventstore.Event) error {
	v.Sequence = event.Sequence()
	v.ChangeDate = event.CreatedAt()
	switch event.Type() {
	case user.UserV1PasswordCheckSucceededType,
		user.HumanPasswordCheckSucceededType:
		v.PasswordVerification = event.CreatedAt()
		v.State = int32(domain.UserSessionStateActive)
	case user.UserIDPLoginCheckSucceededType:
		data := new(es_model.AuthRequest)
		err := data.SetData(event)
		if err != nil {
			return err
		}
		v.ExternalLoginVerification = event.CreatedAt()
		v.SelectedIDPConfigID = data.SelectedIDPConfigID
		v.State = int32(domain.UserSessionStateActive)
	case user.HumanPasswordlessTokenCheckSucceededType:
		v.PasswordlessVerification = event.CreatedAt()
		v.MultiFactorVerification = event.CreatedAt()
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
			v.setSecondFactorVerification(event.CreatedAt(), domain.MFATypeTOTP)
		}
	case user.UserV1MFAOTPCheckSucceededType,
		user.HumanMFAOTPCheckSucceededType:
		v.setSecondFactorVerification(event.CreatedAt(), domain.MFATypeTOTP)
	case user.HumanOTPSMSCheckSucceededType:
		data := new(es_model.OTPVerified)
		err := data.SetData(event)
		if err != nil {
			return err
		}
		if v.UserAgentID == data.UserAgentID {
			v.setSecondFactorVerification(event.CreatedAt(), domain.MFATypeOTPSMS)
		}
	case user.HumanOTPEmailCheckSucceededType:
		data := new(es_model.OTPVerified)
		err := data.SetData(event)
		if err != nil {
			return err
		}
		if v.UserAgentID == data.UserAgentID {
			v.setSecondFactorVerification(event.CreatedAt(), domain.MFATypeOTPEmail)
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
			v.setSecondFactorVerification(event.CreatedAt(), domain.MFATypeU2F)
		}
	case user.HumanU2FTokenCheckSucceededType:
		v.setSecondFactorVerification(event.CreatedAt(), domain.MFATypeU2F)
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
	}
	return nil
}

func (v *UserSessionView) setSecondFactorVerification(verificationTime time.Time, mfaType domain.MFAType) {
	v.SecondFactorVerification = verificationTime
	v.SecondFactorVerificationType = int32(mfaType)
	v.State = int32(domain.UserSessionStateActive)
}

func (v *UserSessionView) EventTypes() []eventstore.EventType {
	return []eventstore.EventType{
		user.UserV1PasswordCheckSucceededType,
		user.HumanPasswordCheckSucceededType,
		user.UserIDPLoginCheckSucceededType,
		user.HumanPasswordlessTokenCheckSucceededType,
		user.HumanPasswordlessTokenCheckFailedType,
		user.HumanPasswordlessTokenRemovedType,
		user.UserV1PasswordCheckFailedType,
		user.HumanPasswordCheckFailedType,
		user.UserV1PasswordChangedType,
		user.HumanPasswordChangedType,
		user.HumanMFAOTPVerifiedType,
		user.UserV1MFAOTPCheckSucceededType,
		user.HumanMFAOTPCheckSucceededType,
		user.UserV1MFAOTPCheckFailedType,
		user.UserV1MFAOTPRemovedType,
		user.HumanMFAOTPCheckFailedType,
		user.HumanMFAOTPRemovedType,
		user.HumanOTPSMSCheckSucceededType,
		user.HumanOTPSMSCheckFailedType,
		user.HumanOTPEmailCheckSucceededType,
		user.HumanOTPEmailCheckFailedType,
		user.HumanU2FTokenCheckFailedType,
		user.HumanU2FTokenRemovedType,
		user.HumanU2FTokenVerifiedType,
		user.HumanU2FTokenCheckSucceededType,
		user.UserV1SignedOutType,
		user.HumanSignedOutType,
		user.UserLockedType,
		user.UserDeactivatedType,
		user.UserIDPLinkRemovedType,
		user.UserIDPLinkCascadeRemovedType,
	}
}
