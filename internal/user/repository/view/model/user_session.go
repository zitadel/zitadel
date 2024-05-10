package model

import (
	"time"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/user/model"
	es_model "github.com/zitadel/zitadel/internal/user/repository/eventsourcing/model"
	"github.com/zitadel/zitadel/internal/view/repository"
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
	UserAgentID string `json:"userAgentID" gorm:"column:user_agent_id;primary_key"`
	UserID      string `json:"userID" gorm:"column:user_id;primary_key"`
	InstanceID  string `json:"instanceID" gorm:"column:instance_id;primary_key"`

	// As of https://github.com/zitadel/zitadel/pull/7199 the following 4 attributes
	// are not projected in the user session handler anymore
	// and are therefore annotated with a `gorm:"-"`.
	// They will be read from the corresponding projection directly.
	UserName    string `json:"-" gorm:"-"`
	LoginName   string `json:"-" gorm:"-"`
	DisplayName string `json:"-" gorm:"-"`
	AvatarKey   string `json:"-" gorm:"-"`

	CreationDate                 repository.Field[time.Time] `json:"-" gorm:"column:creation_date"`
	ChangeDate                   repository.Field[time.Time] `json:"-" gorm:"column:change_date"`
	ResourceOwner                repository.Field[string]    `json:"-" gorm:"column:resource_owner"`
	State                        repository.Field[int32]     `json:"-" gorm:"column:state"`
	SelectedIDPConfigID          repository.Field[string]    `json:"selectedIDPConfigID" gorm:"column:selected_idp_config_id"`
	PasswordVerification         repository.Field[time.Time] `json:"-" gorm:"column:password_verification"`
	PasswordlessVerification     repository.Field[time.Time] `json:"-" gorm:"column:passwordless_verification"`
	ExternalLoginVerification    repository.Field[time.Time] `json:"-" gorm:"column:external_login_verification"`
	SecondFactorVerification     repository.Field[time.Time] `json:"-" gorm:"column:second_factor_verification"`
	SecondFactorVerificationType repository.Field[int32]     `json:"-" gorm:"column:second_factor_verification_type"`
	MultiFactorVerification      repository.Field[time.Time] `json:"-" gorm:"column:multi_factor_verification"`
	MultiFactorVerificationType  repository.Field[int32]     `json:"-" gorm:"column:multi_factor_verification_type"`
	Sequence                     repository.Field[uint64]    `json:"-" gorm:"column:sequence"`
}

func UserSessionFromEvent(event eventstore.Event) (*UserSessionView, error) {
	v := new(UserSessionView)
	if err := event.Unmarshal(v); err != nil {
		logging.WithError(err).Error("could not unmarshal event data")
		return nil, zerrors.ThrowInternal(nil, "MODEL-sd325", "could not unmarshal data")
	}
	return v, nil
}

func UserSessionToModel(userSession *UserSessionView) *model.UserSessionView {
	return &model.UserSessionView{
		UserAgentID:                  userSession.UserAgentID,
		UserID:                       userSession.UserID,
		UserName:                     userSession.UserName,
		LoginName:                    userSession.LoginName,
		DisplayName:                  userSession.DisplayName,
		AvatarKey:                    userSession.AvatarKey,
		ChangeDate:                   userSession.ChangeDate.Value(),
		CreationDate:                 userSession.CreationDate.Value(),
		ResourceOwner:                userSession.ResourceOwner.Value(),
		State:                        domain.UserSessionState(userSession.State.Value()),
		SelectedIDPConfigID:          userSession.SelectedIDPConfigID.Value(),
		PasswordVerification:         userSession.PasswordVerification.Value(),
		PasswordlessVerification:     userSession.PasswordlessVerification.Value(),
		ExternalLoginVerification:    userSession.ExternalLoginVerification.Value(),
		SecondFactorVerification:     userSession.SecondFactorVerification.Value(),
		SecondFactorVerificationType: domain.MFAType(userSession.SecondFactorVerificationType.Value()),
		MultiFactorVerification:      userSession.MultiFactorVerification.Value(),
		MultiFactorVerificationType:  domain.MFAType(userSession.MultiFactorVerificationType.Value()),
		Sequence:                     userSession.Sequence.Value(),
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
	v.Sequence.Set(event.Sequence())
	v.ChangeDate.Set(event.CreatedAt())
	switch event.Type() {
	case user.UserV1PasswordCheckSucceededType,
		user.HumanPasswordCheckSucceededType:
		v.PasswordVerification.Set(event.CreatedAt())
		v.State.Set(int32(domain.UserSessionStateActive))
	case user.UserIDPLoginCheckSucceededType:
		data := new(es_model.AuthRequest)
		err := data.SetData(event)
		if err != nil {
			return err
		}
		v.ExternalLoginVerification.Set(event.CreatedAt())
		v.SelectedIDPConfigID.Set(data.SelectedIDPConfigID)
		v.State.Set(int32(domain.UserSessionStateActive))
	case user.HumanPasswordlessTokenCheckSucceededType:
		v.PasswordlessVerification.Set(event.CreatedAt())
		v.MultiFactorVerification.Set(event.CreatedAt())
		v.MultiFactorVerificationType.Set(int32(domain.MFATypeU2FUserVerification))
		v.State.Set(int32(domain.UserSessionStateActive))
	case user.HumanPasswordlessTokenCheckFailedType,
		user.HumanPasswordlessTokenRemovedType:
		v.PasswordlessVerification.Set(time.Time{})
		v.MultiFactorVerification.Set(time.Time{})
	case user.UserV1PasswordCheckFailedType,
		user.HumanPasswordCheckFailedType:
		v.PasswordVerification.Set(time.Time{})
	case user.UserV1PasswordChangedType,
		user.HumanPasswordChangedType:
		data := new(es_model.PasswordChange)
		err := data.SetData(event)
		if err != nil {
			return err
		}
		if v.UserAgentID != data.UserAgentID {
			v.PasswordVerification.Set(time.Time{})
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
		v.SecondFactorVerification.Set(time.Time{})
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
		v.PasswordlessVerification.Set(time.Time{})
		v.PasswordVerification.Set(time.Time{})
		v.SecondFactorVerification.Set(time.Time{})
		v.SecondFactorVerificationType.Set(int32(domain.MFALevelNotSetUp))
		v.MultiFactorVerification.Set(time.Time{})
		v.MultiFactorVerificationType.Set(int32(domain.MFALevelNotSetUp))
		v.ExternalLoginVerification.Set(time.Time{})
		v.State.Set(int32(domain.UserSessionStateTerminated))
	case user.UserIDPLinkRemovedType, user.UserIDPLinkCascadeRemovedType:
		v.ExternalLoginVerification.Set(time.Time{})
		v.SelectedIDPConfigID.Set("")
	}
	return nil
}

func (v *UserSessionView) setSecondFactorVerification(verificationTime time.Time, mfaType domain.MFAType) {
	v.SecondFactorVerification.Set(verificationTime)
	v.SecondFactorVerificationType.Set(int32(mfaType))
	v.State.Set(int32(domain.UserSessionStateActive))
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

func (v *UserSessionView) PKColumns() []handler.Column {
	return []handler.Column{
		handler.NewCol("user_agent_id", v.UserAgentID),
		handler.NewCol("user_id", v.UserID),
		handler.NewCol("instance_id", v.InstanceID),
	}
}

func (v *UserSessionView) PKConditions() []handler.Condition {
	return []handler.Condition{
		handler.NewCond("user_agent_id", v.UserAgentID),
		handler.NewCond("user_id", v.UserID),
		handler.NewCond("instance_id", v.InstanceID),
	}
}

func (v *UserSessionView) Changes() []handler.Column {
	changes := make([]handler.Column, 0, 12)

	if v.CreationDate.DidChange() {
		changes = append(changes, handler.NewCol("creation_date", v.CreationDate.Value()))
	}
	if v.ChangeDate.DidChange() {
		changes = append(changes, handler.NewCol("change_date", v.ChangeDate.Value()))
	}
	if v.ResourceOwner.DidChange() {
		changes = append(changes, handler.NewCol("resource_owner", v.ResourceOwner.Value()))
	}
	if v.State.DidChange() {
		changes = append(changes, handler.NewCol("state", v.State.Value()))
	}
	if v.SelectedIDPConfigID.DidChange() {
		changes = append(changes, handler.NewCol("selected_idp_config_id", v.SelectedIDPConfigID.Value()))
	}
	if v.PasswordVerification.DidChange() {
		changes = append(changes, handler.NewCol("password_verification", v.PasswordVerification.Value()))
	}
	if v.PasswordlessVerification.DidChange() {
		changes = append(changes, handler.NewCol("passwordless_verification", v.PasswordlessVerification.Value()))
	}
	if v.ExternalLoginVerification.DidChange() {
		changes = append(changes, handler.NewCol("external_login_verification", v.ExternalLoginVerification.Value()))
	}
	if v.SecondFactorVerification.DidChange() {
		changes = append(changes, handler.NewCol("second_factor_verification", v.SecondFactorVerification.Value()))
	}
	if v.SecondFactorVerificationType.DidChange() {
		changes = append(changes, handler.NewCol("second_factor_verification_type", v.SecondFactorVerificationType.Value()))
	}
	if v.MultiFactorVerification.DidChange() {
		changes = append(changes, handler.NewCol("multi_factor_verification", v.MultiFactorVerification.Value()))
	}
	if v.MultiFactorVerificationType.DidChange() {
		changes = append(changes, handler.NewCol("multi_factor_verification_type", v.MultiFactorVerificationType.Value()))
	}
	if v.Sequence.DidChange() {
		changes = append(changes, handler.NewCol("sequence", v.Sequence.Value()))
	}

	return changes
}
