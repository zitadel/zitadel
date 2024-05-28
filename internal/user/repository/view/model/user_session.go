package model

import (
	"database/sql"
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
	UserSessionKeyUserAgentID                  = "user_agent_id"
	UserSessionKeyUserID                       = "user_id"
	UserSessionKeyState                        = "state"
	UserSessionKeyResourceOwner                = "resource_owner"
	UserSessionKeyInstanceID                   = "instance_id"
	UserSessionKeyOwnerRemoved                 = "owner_removed"
	UserSessionKeyCreationDate                 = "creation_date"
	UserSessionKeyChangeDate                   = "change_date"
	UserSessionKeySequence                     = "sequence"
	UserSessionKeyPasswordVerification         = "password_verification"
	UserSessionKeySecondFactorVerification     = "second_factor_verification"
	UserSessionKeySecondFactorVerificationType = "second_factor_verification_type"
	UserSessionKeyMultiFactorVerification      = "multi_factor_verification"
	UserSessionKeyMultiFactorVerificationType  = "multi_factor_verification_type"
	UserSessionKeyPasswordlessVerification     = "passwordless_verification"
	UserSessionKeyExternalLoginVerification    = "external_login_verification"
	UserSessionKeySelectedIDPConfigID          = "selected_idp_config_id"
)

type UserSessionView struct {
	CreationDate  time.Time                         `json:"-" gorm:"column:creation_date"`
	ChangeDate    time.Time                         `json:"-" gorm:"column:change_date"`
	ResourceOwner string                            `json:"-" gorm:"column:resource_owner"`
	State         sql.Null[domain.UserSessionState] `json:"-" gorm:"column:state"`
	UserAgentID   string                            `json:"userAgentID" gorm:"column:user_agent_id;primary_key"`
	UserID        string                            `json:"userID" gorm:"column:user_id;primary_key"`
	// As of https://github.com/zitadel/zitadel/pull/7199 the following 4 attributes
	// are not projected in the user session handler anymore
	// and are therefore annotated with a `gorm:"-"`.
	// They will be read from the corresponding projection directly.
	UserName                     sql.NullString `json:"-" gorm:"-"`
	LoginName                    sql.NullString `json:"-" gorm:"-"`
	DisplayName                  sql.NullString `json:"-" gorm:"-"`
	AvatarKey                    sql.NullString `json:"-" gorm:"-"`
	SelectedIDPConfigID          sql.NullString `json:"selectedIDPConfigID" gorm:"column:selected_idp_config_id"`
	PasswordVerification         sql.NullTime   `json:"-" gorm:"column:password_verification"`
	PasswordlessVerification     sql.NullTime   `json:"-" gorm:"column:passwordless_verification"`
	ExternalLoginVerification    sql.NullTime   `json:"-" gorm:"column:external_login_verification"`
	SecondFactorVerification     sql.NullTime   `json:"-" gorm:"column:second_factor_verification"`
	SecondFactorVerificationType sql.NullInt32  `json:"-" gorm:"column:second_factor_verification_type"`
	MultiFactorVerification      sql.NullTime   `json:"-" gorm:"column:multi_factor_verification"`
	MultiFactorVerificationType  sql.NullInt32  `json:"-" gorm:"column:multi_factor_verification_type"`
	Sequence                     uint64         `json:"-" gorm:"column:sequence"`
	InstanceID                   string         `json:"instanceID" gorm:"column:instance_id;primary_key"`
}

type userAgentIDPayload struct {
	ID string `json:"userAgentID"`
}

func UserAgentIDFromEvent(event eventstore.Event) (string, error) {
	payload := new(userAgentIDPayload)
	if err := event.Unmarshal(payload); err != nil {
		logging.WithError(err).Error("could not unmarshal event data")
		return "", zerrors.ThrowInternal(nil, "MODEL-HJwk9", "could not unmarshal data")
	}
	return payload.ID, nil
}

func UserSessionToModel(userSession *UserSessionView) *model.UserSessionView {
	return &model.UserSessionView{
		ChangeDate:                   userSession.ChangeDate,
		CreationDate:                 userSession.CreationDate,
		ResourceOwner:                userSession.ResourceOwner,
		State:                        userSession.State.V,
		UserAgentID:                  userSession.UserAgentID,
		UserID:                       userSession.UserID,
		UserName:                     userSession.UserName.String,
		LoginName:                    userSession.LoginName.String,
		DisplayName:                  userSession.DisplayName.String,
		AvatarKey:                    userSession.AvatarKey.String,
		SelectedIDPConfigID:          userSession.SelectedIDPConfigID.String,
		PasswordVerification:         userSession.PasswordVerification.Time,
		PasswordlessVerification:     userSession.PasswordlessVerification.Time,
		ExternalLoginVerification:    userSession.ExternalLoginVerification.Time,
		SecondFactorVerification:     userSession.SecondFactorVerification.Time,
		SecondFactorVerificationType: domain.MFAType(userSession.SecondFactorVerificationType.Int32),
		MultiFactorVerification:      userSession.MultiFactorVerification.Time,
		MultiFactorVerificationType:  domain.MFAType(userSession.MultiFactorVerificationType.Int32),
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
	// in case anything needs to be change here check if the Reduce function needs the change as well
	v.Sequence = event.Sequence()
	v.ChangeDate = event.CreatedAt()
	switch event.Type() {
	case user.UserV1PasswordCheckSucceededType,
		user.HumanPasswordCheckSucceededType:
		v.PasswordVerification = sql.NullTime{Time: event.CreatedAt(), Valid: true}
		v.State.V = domain.UserSessionStateActive
	case user.UserIDPLoginCheckSucceededType:
		data := new(es_model.AuthRequest)
		err := data.SetData(event)
		if err != nil {
			return err
		}
		v.ExternalLoginVerification = sql.NullTime{Time: event.CreatedAt(), Valid: true}
		v.SelectedIDPConfigID = sql.NullString{String: data.SelectedIDPConfigID, Valid: true}
		v.State.V = domain.UserSessionStateActive
	case user.HumanPasswordlessTokenCheckSucceededType:
		v.PasswordlessVerification = sql.NullTime{Time: event.CreatedAt(), Valid: true}
		v.MultiFactorVerification = sql.NullTime{Time: event.CreatedAt(), Valid: true}
		v.MultiFactorVerificationType = sql.NullInt32{Int32: int32(domain.MFATypeU2FUserVerification)}
		v.State.V = domain.UserSessionStateActive
	case user.HumanPasswordlessTokenCheckFailedType,
		user.HumanPasswordlessTokenRemovedType:
		v.PasswordlessVerification = sql.NullTime{Time: time.Time{}, Valid: true}
		v.MultiFactorVerification = sql.NullTime{Time: time.Time{}, Valid: true}
	case user.UserV1PasswordCheckFailedType,
		user.HumanPasswordCheckFailedType:
		v.PasswordVerification = sql.NullTime{Time: time.Time{}, Valid: true}
	case user.UserV1PasswordChangedType,
		user.HumanPasswordChangedType:
		data := new(es_model.PasswordChange)
		err := data.SetData(event)
		if err != nil {
			return err
		}
		if v.UserAgentID != data.UserAgentID {
			v.PasswordVerification = sql.NullTime{Time: time.Time{}, Valid: true}
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
		v.SecondFactorVerification = sql.NullTime{Time: time.Time{}, Valid: true}
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
		v.PasswordlessVerification = sql.NullTime{Time: time.Time{}, Valid: true}
		v.PasswordVerification = sql.NullTime{Time: time.Time{}, Valid: true}
		v.SecondFactorVerification = sql.NullTime{Time: time.Time{}, Valid: true}
		v.SecondFactorVerificationType = sql.NullInt32{Int32: int32(domain.MFALevelNotSetUp)}
		v.MultiFactorVerification = sql.NullTime{Time: time.Time{}, Valid: true}
		v.MultiFactorVerificationType = sql.NullInt32{Int32: int32(domain.MFALevelNotSetUp)}
		v.ExternalLoginVerification = sql.NullTime{Time: time.Time{}, Valid: true}
		v.State.V = domain.UserSessionStateTerminated
	case user.UserIDPLinkRemovedType, user.UserIDPLinkCascadeRemovedType:
		v.ExternalLoginVerification = sql.NullTime{Time: time.Time{}, Valid: true}
		v.SelectedIDPConfigID = sql.NullString{String: "", Valid: true}
	}
	return nil
}

func (v *UserSessionView) setSecondFactorVerification(verificationTime time.Time, mfaType domain.MFAType) {
	v.SecondFactorVerification = sql.NullTime{Time: verificationTime, Valid: true}
	v.SecondFactorVerificationType = sql.NullInt32{Int32: int32(mfaType)}
	v.State.V = domain.UserSessionStateActive
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
