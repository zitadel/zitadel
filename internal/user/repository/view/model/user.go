package model

import (
	"database/sql/driver"
	"encoding/json"
	"slices"
	"time"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	org_model "github.com/zitadel/zitadel/internal/org/model"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/user/model"
	es_model "github.com/zitadel/zitadel/internal/user/repository/eventsourcing/model"
	"github.com/zitadel/zitadel/internal/view/repository"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	UserKeyUserID             = "id"
	UserKeyUserName           = "user_name"
	UserKeyFirstName          = "first_name"
	UserKeyLastName           = "last_name"
	UserKeyNickName           = "nick_name"
	UserKeyDisplayName        = "display_name"
	UserKeyEmail              = "email"
	UserKeyState              = "user_state"
	UserKeyResourceOwner      = "resource_owner"
	UserKeyLoginNames         = "login_names"
	UserKeyPreferredLoginName = "preferred_login_name"
	UserKeyType               = "user_type"
	UserKeyInstanceID         = "instance_id"
	UserKeyOwnerRemoved       = "owner_removed"
)

type userType string

const (
	userTypeHuman   = "human"
	userTypeMachine = "machine"
)

type UserView struct {
	ID         string `json:"-" gorm:"column:id;primary_key"`
	InstanceID string `json:"instanceID" gorm:"column:instance_id;primary_key"`

	CreationDate       repository.Field[time.Time]                  `json:"-" gorm:"column:creation_date"`
	ChangeDate         repository.Field[time.Time]                  `json:"-" gorm:"column:change_date"`
	ResourceOwner      repository.Field[string]                     `json:"-" gorm:"column:resource_owner"`
	State              repository.Field[int32]                      `json:"-" gorm:"column:user_state"`
	LastLogin          repository.Field[time.Time]                  `json:"-" gorm:"column:last_login"`
	LoginNames         repository.Field[database.TextArray[string]] `json:"-" gorm:"column:login_names"`
	PreferredLoginName repository.Field[string]                     `json:"-" gorm:"column:preferred_login_name"`
	Sequence           repository.Field[uint64]                     `json:"-" gorm:"column:sequence"`
	Type               repository.Field[userType]                   `json:"-" gorm:"column:user_type"`
	UserName           repository.Field[string]                     `json:"userName" gorm:"column:user_name"`
	MachineView        repository.Field[*MachineView]
	HumanView          repository.Field[*HumanView]
}

type UserState int32

const (
	UserStateUnspecified UserState = iota
	UserStateActive
	UserStateInactive
	UserStateDeleted
	UserStateLocked
	UserStateSuspend
	UserStateInitial
)

type HumanView struct {
	FirstName                repository.Field[string]         `json:"firstName" gorm:"column:first_name"`
	LastName                 repository.Field[string]         `json:"lastName" gorm:"column:last_name"`
	NickName                 repository.Field[string]         `json:"nickName" gorm:"column:nick_name"`
	DisplayName              repository.Field[string]         `json:"displayName" gorm:"column:display_name"`
	PreferredLanguage        repository.Field[string]         `json:"preferredLanguage" gorm:"column:preferred_language"`
	Gender                   repository.Field[int32]          `json:"gender" gorm:"column:gender"`
	AvatarKey                repository.Field[string]         `json:"storeKey" gorm:"column:avatar_key"`
	Email                    repository.Field[string]         `json:"email" gorm:"column:email"`
	IsEmailVerified          repository.Field[bool]           `json:"-" gorm:"column:is_email_verified"`
	Phone                    repository.Field[string]         `json:"phone" gorm:"column:phone"`
	IsPhoneVerified          repository.Field[bool]           `json:"-" gorm:"column:is_phone_verified"`
	Country                  repository.Field[string]         `json:"country" gorm:"column:country"`
	Locality                 repository.Field[string]         `json:"locality" gorm:"column:locality"`
	PostalCode               repository.Field[string]         `json:"postalCode" gorm:"column:postal_code"`
	Region                   repository.Field[string]         `json:"region" gorm:"column:region"`
	StreetAddress            repository.Field[string]         `json:"streetAddress" gorm:"column:street_address"`
	OTPState                 repository.Field[int32]          `json:"-" gorm:"column:otp_state"`
	OTPSMSAdded              repository.Field[bool]           `json:"-" gorm:"column:otp_sms_added"`
	OTPEmailAdded            repository.Field[bool]           `json:"-" gorm:"column:otp_email_added"`
	U2FTokens                repository.Field[WebAuthNTokens] `json:"-" gorm:"column:u2f_tokens"`
	MFAMaxSetUp              repository.Field[int32]          `json:"-" gorm:"column:mfa_max_set_up"`
	MFAInitSkipped           repository.Field[time.Time]      `json:"-" gorm:"column:mfa_init_skipped"`
	InitRequired             repository.Field[bool]           `json:"-" gorm:"column:init_required"`
	PasswordlessInitRequired repository.Field[bool]           `json:"-" gorm:"column:passwordless_init_required"`
	PasswordInitRequired     repository.Field[bool]           `json:"-" gorm:"column:password_init_required"`
	PasswordSet              repository.Field[bool]           `json:"-" gorm:"column:password_set"`
	PasswordChangeRequired   repository.Field[bool]           `json:"-" gorm:"column:password_change_required"`
	UsernameChangeRequired   repository.Field[bool]           `json:"-" gorm:"column:username_change_required"`
	PasswordChanged          repository.Field[time.Time]      `json:"-" gorm:"column:password_change"`
	PasswordlessTokens       repository.Field[WebAuthNTokens] `json:"-" gorm:"column:passwordless_tokens"`
}

type WebAuthNTokens []*WebAuthNView

type WebAuthNView struct {
	ID    repository.Field[string] `json:"webAuthNTokenId"`
	Name  repository.Field[string] `json:"webAuthNTokenName,omitempty"`
	State repository.Field[int32]  `json:"state,omitempty"`
}

func (t WebAuthNTokens) Value() (driver.Value, error) {
	if t == nil {
		return nil, nil
	}
	return json.Marshal(&t)
}

func (t *WebAuthNTokens) Scan(src interface{}) error {
	if b, ok := src.([]byte); ok {
		return json.Unmarshal(b, t)
	}
	if s, ok := src.(string); ok {
		return json.Unmarshal([]byte(s), t)
	}
	return nil
}

func (h *HumanView) IsZero() bool {
	return h == nil || h.FirstName.Value() == ""
}

type MachineView struct {
	Name        repository.Field[string] `json:"name" gorm:"column:machine_name"`
	Description repository.Field[string] `json:"description" gorm:"column:machine_description"`
}

func (m *MachineView) IsZero() bool {
	return m == nil || m.Name.Value() == ""
}

func UserToModel(user *UserView) *model.UserView {
	userView := &model.UserView{
		ID:                 user.ID,
		UserName:           user.UserName.Value(),
		ChangeDate:         user.ChangeDate.Value(),
		CreationDate:       user.CreationDate.Value(),
		ResourceOwner:      user.ResourceOwner.Value(),
		State:              model.UserState(user.State.Value()),
		LastLogin:          user.LastLogin.Value(),
		PreferredLoginName: user.PreferredLoginName.Value(),
		LoginNames:         user.LoginNames.Value(),
		Sequence:           user.Sequence.Value(),
	}
	if !user.HumanView.Value().IsZero() {
		userView.HumanView = &model.HumanView{
			PasswordSet:              user.HumanView.Value().PasswordSet.Value(),
			PasswordInitRequired:     user.HumanView.Value().PasswordInitRequired.Value(),
			PasswordChangeRequired:   user.HumanView.Value().PasswordChangeRequired.Value(),
			PasswordChanged:          user.HumanView.Value().PasswordChanged.Value(),
			PasswordlessTokens:       WebauthnTokensToModel(user.HumanView.Value().PasswordlessTokens.Value()),
			U2FTokens:                WebauthnTokensToModel(user.HumanView.Value().U2FTokens.Value()),
			FirstName:                user.HumanView.Value().FirstName.Value(),
			LastName:                 user.HumanView.Value().LastName.Value(),
			NickName:                 user.HumanView.Value().NickName.Value(),
			DisplayName:              user.HumanView.Value().DisplayName.Value(),
			AvatarKey:                user.HumanView.Value().AvatarKey.Value(),
			PreferredLanguage:        user.HumanView.Value().PreferredLanguage.Value(),
			Gender:                   model.Gender(user.HumanView.Value().Gender.Value()),
			Email:                    user.HumanView.Value().Email.Value(),
			IsEmailVerified:          user.HumanView.Value().IsEmailVerified.Value(),
			Phone:                    user.HumanView.Value().Phone.Value(),
			IsPhoneVerified:          user.HumanView.Value().IsPhoneVerified.Value(),
			Country:                  user.HumanView.Value().Country.Value(),
			Locality:                 user.HumanView.Value().Locality.Value(),
			PostalCode:               user.HumanView.Value().PostalCode.Value(),
			Region:                   user.HumanView.Value().Region.Value(),
			StreetAddress:            user.HumanView.Value().StreetAddress.Value(),
			OTPState:                 model.MFAState(user.HumanView.Value().OTPState.Value()),
			OTPSMSAdded:              user.HumanView.Value().OTPSMSAdded.Value(),
			OTPEmailAdded:            user.HumanView.Value().OTPEmailAdded.Value(),
			MFAMaxSetUp:              domain.MFALevel(user.HumanView.Value().MFAMaxSetUp.Value()),
			MFAInitSkipped:           user.HumanView.Value().MFAInitSkipped.Value(),
			InitRequired:             user.HumanView.Value().InitRequired.Value(),
			PasswordlessInitRequired: user.HumanView.Value().PasswordlessInitRequired.Value(),
		}
	}

	if !user.MachineView.Value().IsZero() {
		userView.MachineView = &model.MachineView{
			Description: user.MachineView.Value().Description.Value(),
			Name:        user.MachineView.Value().Name.Value(),
		}
	}
	return userView
}

func WebauthnTokensToModel(tokens []*WebAuthNView) []*model.WebAuthNView {
	if tokens == nil {
		return nil
	}
	result := make([]*model.WebAuthNView, len(tokens))
	for i, t := range tokens {
		result[i] = WebauthnTokenToModel(t)
	}
	return result
}

func WebauthnTokenToModel(token *WebAuthNView) *model.WebAuthNView {
	return &model.WebAuthNView{
		TokenID: token.ID.Value(),
		Name:    token.Name.Value(),
		State:   model.MFAState(token.State.Value()),
	}
}

func (u *UserView) GenerateLoginName(domain string, appendDomain bool) string {
	if !appendDomain {
		return u.UserName.Value()
	}
	return u.UserName.Value() + "@" + domain
}

func (u *UserView) SetLoginNames(userLoginMustBeDomain bool, domains []*org_model.OrgDomain) {
	loginNames := make([]string, 0, len(domains))

	for _, d := range domains {
		if d.Verified {
			loginNames = append(loginNames, u.GenerateLoginName(d.Domain, true))
		}
	}
	if !userLoginMustBeDomain {
		loginNames = append(loginNames, u.GenerateLoginName(u.UserName.Value(), true))
	}

	u.LoginNames.Set(loginNames)
}

func (u *UserView) AppendEvent(event eventstore.Event) (err error) {
	u.ChangeDate.Set(event.CreatedAt())
	u.Sequence.Set(event.Sequence())
	switch event.Type() {
	case user.MachineAddedEventType:
		u.CreationDate.Set(event.CreatedAt())
		u.setRootData(event)
		u.Type.Set(userTypeMachine)
		err = u.setData(event)
		if err != nil {
			return err
		}
	case user.UserV1AddedType,
		user.UserV1RegisteredType,
		user.HumanRegisteredType,
		user.HumanAddedType:
		u.CreationDate.Set(event.CreatedAt())
		u.setRootData(event)
		u.Type.Set(userTypeHuman)
		err = u.setData(event)
		if err != nil {
			return err
		}
		err = u.setPasswordData(event)
	case user.UserRemovedType:
		u.State.Set(int32(model.UserStateDeleted))
	case user.UserV1PasswordChangedType,
		user.HumanPasswordChangedType:
		err = u.setPasswordData(event)
	case user.HumanPasswordlessTokenAddedType:
		err = u.addPasswordlessToken(event)
	case user.HumanPasswordlessTokenVerifiedType:
		err = u.updatePasswordlessToken(event)
	case user.HumanPasswordlessTokenRemovedType:
		err = u.removePasswordlessToken(event)
	case user.UserV1ProfileChangedType,
		user.HumanProfileChangedType,
		user.UserV1AddressChangedType,
		user.HumanAddressChangedType,
		user.MachineChangedEventType:
		err = u.setData(event)
	case user.UserDomainClaimedType:
		if u.HumanView.Value() != nil {
			u.HumanView.Value().UsernameChangeRequired.Set(true)
		}
		err = u.setData(event)
	case user.UserUserNameChangedType:
		if u.HumanView.Value() != nil {
			u.HumanView.Value().UsernameChangeRequired.Set(false)
		}
		err = u.setData(event)
	case user.UserV1EmailChangedType,
		user.HumanEmailChangedType:
		u.HumanView.Value().IsEmailVerified.Set(false)
		err = u.setData(event)
	case user.UserV1EmailVerifiedType,
		user.HumanEmailVerifiedType:
		u.HumanView.Value().IsEmailVerified.Set(true)
	case user.UserV1PhoneChangedType,
		user.HumanPhoneChangedType:
		u.HumanView.Value().IsPhoneVerified.Set(false)
		err = u.setData(event)
	case user.UserV1PhoneVerifiedType,
		user.HumanPhoneVerifiedType:
		u.HumanView.Value().IsPhoneVerified.Set(true)
	case user.UserV1PhoneRemovedType,
		user.HumanPhoneRemovedType:
		u.HumanView.Value().Phone.Set("")
		u.HumanView.Value().IsPhoneVerified.Set(false)
		u.HumanView.Value().OTPSMSAdded.Set(false)
		u.HumanView.Value().MFAInitSkipped.Set(time.Time{})
	case user.UserDeactivatedType:
		u.State.Set(int32(model.UserStateInactive))
	case user.UserReactivatedType,
		user.UserUnlockedType:
		u.State.Set(int32(model.UserStateActive))
	case user.UserLockedType:
		u.State.Set(int32(model.UserStateLocked))
	case user.UserV1MFAOTPAddedType,
		user.HumanMFAOTPAddedType:
		if u.HumanView.Value() == nil {
			logging.WithFields("event_sequence", event.Sequence, "aggregate_id", event.Aggregate().ID, "instance", event.Aggregate().InstanceID).Warn("event is ignored because human not exists")
			return zerrors.ThrowInvalidArgument(nil, "MODEL-p2BXx", "event ignored: human not exists")
		}
		u.HumanView.Value().OTPState.Set(int32(model.MFAStateNotReady))
	case user.UserV1MFAOTPVerifiedType,
		user.HumanMFAOTPVerifiedType:
		if u.HumanView.Value() == nil {
			logging.WithFields("event_sequence", event.Sequence, "aggregate_id", event.Aggregate().ID, "instance", event.Aggregate().InstanceID).Warn("event is ignored because human not exists")
			return zerrors.ThrowInvalidArgument(nil, "MODEL-o6Lcq", "event ignored: human not exists")
		}
		u.HumanView.Value().OTPState.Set(int32(model.MFAStateReady))
		u.HumanView.Value().MFAInitSkipped.Set(time.Time{})
	case user.UserV1MFAOTPRemovedType,
		user.HumanMFAOTPRemovedType:
		u.HumanView.Value().OTPState.Set(int32(model.MFAStateUnspecified))
	case user.HumanOTPSMSAddedType:
		u.HumanView.Value().OTPSMSAdded.Set(true)
	case user.HumanOTPSMSRemovedType:
		u.HumanView.Value().OTPSMSAdded.Set(false)
		u.HumanView.Value().MFAInitSkipped.Set(time.Time{})
	case user.HumanOTPEmailAddedType:
		u.HumanView.Value().OTPEmailAdded.Set(true)
	case user.HumanOTPEmailRemovedType:
		u.HumanView.Value().OTPEmailAdded.Set(false)
		u.HumanView.Value().MFAInitSkipped.Set(time.Time{})
	case user.HumanU2FTokenAddedType:
		err = u.addU2FToken(event)
	case user.HumanU2FTokenVerifiedType:
		err = u.updateU2FToken(event)
		if err != nil {
			return err
		}
		u.HumanView.Value().MFAInitSkipped.Set(time.Time{})
	case user.HumanU2FTokenRemovedType:
		err = u.removeU2FToken(event)
	case user.UserV1MFAInitSkippedType,
		user.HumanMFAInitSkippedType:
		u.HumanView.Value().MFAInitSkipped.Set(event.CreatedAt())
	case user.UserV1InitialCodeAddedType,
		user.HumanInitialCodeAddedType:
		u.HumanView.Value().InitRequired.Set(true)
	case user.UserV1InitializedCheckSucceededType,
		user.HumanInitializedCheckSucceededType:
		u.HumanView.Value().InitRequired.Set(false)
	case user.HumanAvatarAddedType:
		err = u.setData(event)
	case user.HumanAvatarRemovedType:
		u.HumanView.Value().AvatarKey.Set("")
	case user.HumanPasswordlessInitCodeAddedType,
		user.HumanPasswordlessInitCodeRequestedType:
		if u.HumanView.Value() == nil {
			logging.WithFields("event_sequence", event.Sequence, "aggregate_id", event.Aggregate().ID, "instance", event.Aggregate().InstanceID).Warn("event is ignored because human not exists")
			return zerrors.ThrowInvalidArgument(nil, "MODEL-MbyC0", "event ignored: human not exists")
		}
		if !u.HumanView.Value().PasswordSet.Value() {
			u.HumanView.Value().PasswordlessInitRequired.Set(true)
			u.HumanView.Value().PasswordInitRequired.Set(false)
		}
	}
	u.ComputeObject()
	return err
}

func (u *UserView) setRootData(event eventstore.Event) {
	u.ID = event.Aggregate().ID
	u.ResourceOwner.Set(event.Aggregate().ResourceOwner)
	u.InstanceID = event.Aggregate().InstanceID
}

func (u *UserView) setData(event eventstore.Event) error {
	if err := event.Unmarshal(u); err != nil {
		logging.Log("MODEL-lso9e").WithError(err).Error("could not unmarshal event data")
		return zerrors.ThrowInternal(nil, "MODEL-8iows", "could not unmarshal data")
	}
	return nil
}

func (u *UserView) setPasswordData(event eventstore.Event) error {
	password := new(es_model.Password)
	if err := event.Unmarshal(password); err != nil {
		logging.Log("MODEL-sdw4r").WithError(err).Error("could not unmarshal event data")
		return zerrors.ThrowInternal(nil, "MODEL-6jhsw", "could not unmarshal data")
	}
	u.HumanView.Value().PasswordSet.Set(password.Secret != nil || password.EncodedHash != "")
	u.HumanView.Value().PasswordInitRequired.Set(!u.HumanView.Value().PasswordSet.Value())
	u.HumanView.Value().PasswordChangeRequired.Set(password.ChangeRequired)
	u.HumanView.Value().PasswordChanged.Set(event.CreatedAt())
	return nil
}

func (u *UserView) addPasswordlessToken(event eventstore.Event) error {
	token, err := webAuthNViewFromEvent(event)
	if err != nil {
		return err
	}
	for i, t := range u.HumanView.Value().PasswordlessTokens.Value() {
		if t.State.Value() == int32(model.MFAStateNotReady) {
			u.HumanView.Value().PasswordlessTokens.Value()[i].ID.Set(token.ID.Value())
			return nil
		}
	}
	token.State.Set(int32(model.MFAStateNotReady))
	u.HumanView.Value().PasswordlessTokens.Set(append(u.HumanView.Value().PasswordlessTokens.Value(), token))
	return nil
}

func (u *UserView) updatePasswordlessToken(event eventstore.Event) error {
	token, err := webAuthNViewFromEvent(event)
	if err != nil {
		return err
	}
	for i, t := range u.HumanView.Value().PasswordlessTokens.Value() {
		if t.ID == token.ID {
			u.HumanView.Value().PasswordlessTokens.Value()[i].Name.Set(token.Name.Value())
			u.HumanView.Value().PasswordlessTokens.Value()[i].State.Set(int32(model.MFAStateReady))
			return nil
		}
	}
	return nil
}

func (u *UserView) removePasswordlessToken(event eventstore.Event) error {
	token, err := webAuthNViewFromEvent(event)
	if err != nil {
		return err
	}
	for i, t := range u.HumanView.Value().PasswordlessTokens.Value() {
		if t.ID == token.ID {
			u.HumanView.Value().PasswordlessTokens.Set(slices.Delete(u.HumanView.Value().PasswordlessTokens.Value(), i, i+1))
			return nil
		}
	}
	return nil
}

func (u *UserView) addU2FToken(event eventstore.Event) error {
	token, err := webAuthNViewFromEvent(event)
	if err != nil {
		return err
	}
	for i, t := range u.HumanView.Value().U2FTokens.Value() {
		if t.State.Value() == int32(model.MFAStateNotReady) {
			u.HumanView.Value().U2FTokens.Value()[i].ID.Set(token.ID.Value())
			return nil
		}
	}
	token.State.Set(int32(model.MFAStateNotReady))
	u.HumanView.Value().U2FTokens.Set(append(u.HumanView.Value().U2FTokens.Value(), token))
	return nil
}

func (u *UserView) updateU2FToken(event eventstore.Event) error {
	token, err := webAuthNViewFromEvent(event)
	if err != nil {
		return err
	}
	for i, t := range u.HumanView.Value().U2FTokens.Value() {
		if t.ID == token.ID {
			u.HumanView.Value().U2FTokens.Value()[i].Name.Set(token.Name.Value())
			u.HumanView.Value().U2FTokens.Value()[i].State.Set(int32(model.MFAStateReady))
			return nil
		}
	}
	return nil
}

func (u *UserView) removeU2FToken(event eventstore.Event) error {
	token, err := webAuthNViewFromEvent(event)
	if err != nil {
		return err
	}
	for i := len(u.HumanView.Value().U2FTokens.Value()) - 1; i >= 0; i-- {
		if u.HumanView.Value().U2FTokens.Value()[i].ID == token.ID {
			u.HumanView.Value().U2FTokens.Set(slices.Delete(u.HumanView.Value().U2FTokens.Value(), i, i+1))
		}
	}
	return nil
}

func webAuthNViewFromEvent(event eventstore.Event) (*WebAuthNView, error) {
	token := new(WebAuthNView)
	err := event.Unmarshal(token)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "MODEL-FSaq1", "could not unmarshal data")
	}
	return token, err
}

func (u *UserView) ComputeObject() {
	if !u.MachineView.Value().IsZero() {
		if u.State.Value() == int32(model.UserStateUnspecified) {
			u.State.Set(int32(model.UserStateActive))
		}
		return
	}
	if u.State.Value() == int32(model.UserStateUnspecified) || u.State.Value() == int32(model.UserStateInitial) {
		if u.HumanView.Value().IsEmailVerified.Value() {
			u.State.Set(int32(model.UserStateActive))
		} else {
			u.State.Set(int32(model.UserStateInitial))
		}
	}
	u.ComputeMFAMaxSetUp()
}

func (u *UserView) ComputeMFAMaxSetUp() {
	for _, token := range u.HumanView.Value().PasswordlessTokens.Value() {
		if token.State.Value() == int32(model.MFAStateReady) {
			u.HumanView.Value().MFAMaxSetUp.Set(int32(domain.MFALevelMultiFactor))
			u.HumanView.Value().PasswordlessInitRequired.Set(false)
			return
		}
	}
	for _, token := range u.HumanView.Value().U2FTokens.Value() {
		if token.State.Value() == int32(model.MFAStateReady) {
			u.HumanView.Value().MFAMaxSetUp.Set(int32(domain.MFALevelSecondFactor))
			return
		}
	}
	if u.HumanView.Value().OTPState.Value() == int32(model.MFAStateReady) ||
		u.HumanView.Value().OTPSMSAdded.Value() || u.HumanView.Value().OTPEmailAdded.Value() {
		u.HumanView.Value().MFAMaxSetUp.Set(int32(domain.MFALevelSecondFactor))
		return
	}
	u.HumanView.Value().MFAMaxSetUp.Set(int32(domain.MFALevelNotSetUp))
}

func (u *UserView) SetEmptyUserType() {
	if u.MachineView.Value() != nil && u.MachineView.Value().Name.Value() == "" {
		u.MachineView.Set(nil)
	} else {
		u.HumanView.Set(nil)
	}
}

func (u *UserView) EventTypes() []eventstore.EventType {
	return []eventstore.EventType{
		user.MachineAddedEventType,
		user.UserV1AddedType,
		user.UserV1RegisteredType,
		user.HumanRegisteredType,
		user.HumanAddedType,
		user.UserRemovedType,
		user.UserV1PasswordChangedType,
		user.HumanPasswordChangedType,
		user.HumanPasswordlessTokenAddedType,
		user.HumanPasswordlessTokenVerifiedType,
		user.HumanPasswordlessTokenRemovedType,
		user.UserV1ProfileChangedType,
		user.HumanProfileChangedType,
		user.UserV1AddressChangedType,
		user.HumanAddressChangedType,
		user.MachineChangedEventType,
		user.UserDomainClaimedType,
		user.UserUserNameChangedType,
		user.UserV1EmailChangedType,
		user.HumanEmailChangedType,
		user.UserV1EmailVerifiedType,
		user.HumanEmailVerifiedType,
		user.UserV1PhoneChangedType,
		user.HumanPhoneChangedType,
		user.UserV1PhoneVerifiedType,
		user.HumanPhoneVerifiedType,
		user.UserV1PhoneRemovedType,
		user.HumanPhoneRemovedType,
		user.UserDeactivatedType,
		user.UserReactivatedType,
		user.UserUnlockedType,
		user.UserLockedType,
		user.UserV1MFAOTPAddedType,
		user.HumanMFAOTPAddedType,
		user.UserV1MFAOTPVerifiedType,
		user.HumanMFAOTPVerifiedType,
		user.UserV1MFAOTPRemovedType,
		user.HumanMFAOTPRemovedType,
		user.HumanOTPSMSAddedType,
		user.HumanOTPSMSRemovedType,
		user.HumanOTPEmailAddedType,
		user.HumanOTPEmailRemovedType,
		user.HumanU2FTokenAddedType,
		user.HumanU2FTokenVerifiedType,
		user.HumanU2FTokenRemovedType,
		user.UserV1MFAInitSkippedType,
		user.HumanMFAInitSkippedType,
		user.UserV1InitialCodeAddedType,
		user.HumanInitialCodeAddedType,
		user.UserV1InitializedCheckSucceededType,
		user.HumanInitializedCheckSucceededType,
		user.HumanAvatarAddedType,
		user.HumanAvatarRemovedType,
		user.HumanPasswordlessInitCodeAddedType,
		user.HumanPasswordlessInitCodeRequestedType,
	}
}
