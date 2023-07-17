package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	org_model "github.com/zitadel/zitadel/internal/org/model"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/user/model"
	es_model "github.com/zitadel/zitadel/internal/user/repository/eventsourcing/model"
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
	ID                 string               `json:"-" gorm:"column:id;primary_key"`
	CreationDate       time.Time            `json:"-" gorm:"column:creation_date"`
	ChangeDate         time.Time            `json:"-" gorm:"column:change_date"`
	ResourceOwner      string               `json:"-" gorm:"column:resource_owner"`
	State              int32                `json:"-" gorm:"column:user_state"`
	LastLogin          time.Time            `json:"-" gorm:"column:last_login"`
	LoginNames         database.StringArray `json:"-" gorm:"column:login_names"`
	PreferredLoginName string               `json:"-" gorm:"column:preferred_login_name"`
	Sequence           uint64               `json:"-" gorm:"column:sequence"`
	Type               userType             `json:"-" gorm:"column:user_type"`
	UserName           string               `json:"userName" gorm:"column:user_name"`
	InstanceID         string               `json:"instanceID" gorm:"column:instance_id;primary_key"`
	*MachineView
	*HumanView
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
	FirstName                string         `json:"firstName" gorm:"column:first_name"`
	LastName                 string         `json:"lastName" gorm:"column:last_name"`
	NickName                 string         `json:"nickName" gorm:"column:nick_name"`
	DisplayName              string         `json:"displayName" gorm:"column:display_name"`
	PreferredLanguage        string         `json:"preferredLanguage" gorm:"column:preferred_language"`
	Gender                   int32          `json:"gender" gorm:"column:gender"`
	AvatarKey                string         `json:"storeKey" gorm:"column:avatar_key"`
	Email                    string         `json:"email" gorm:"column:email"`
	IsEmailVerified          bool           `json:"-" gorm:"column:is_email_verified"`
	Phone                    string         `json:"phone" gorm:"column:phone"`
	IsPhoneVerified          bool           `json:"-" gorm:"column:is_phone_verified"`
	Country                  string         `json:"country" gorm:"column:country"`
	Locality                 string         `json:"locality" gorm:"column:locality"`
	PostalCode               string         `json:"postalCode" gorm:"column:postal_code"`
	Region                   string         `json:"region" gorm:"column:region"`
	StreetAddress            string         `json:"streetAddress" gorm:"column:street_address"`
	OTPState                 int32          `json:"-" gorm:"column:otp_state"`
	U2FTokens                WebAuthNTokens `json:"-" gorm:"column:u2f_tokens"`
	MFAMaxSetUp              int32          `json:"-" gorm:"column:mfa_max_set_up"`
	MFAInitSkipped           time.Time      `json:"-" gorm:"column:mfa_init_skipped"`
	InitRequired             bool           `json:"-" gorm:"column:init_required"`
	PasswordlessInitRequired bool           `json:"-" gorm:"column:passwordless_init_required"`
	PasswordInitRequired     bool           `json:"-" gorm:"column:password_init_required"`
	PasswordSet              bool           `json:"-" gorm:"column:password_set"`
	PasswordChangeRequired   bool           `json:"-" gorm:"column:password_change_required"`
	UsernameChangeRequired   bool           `json:"-" gorm:"column:username_change_required"`
	PasswordChanged          time.Time      `json:"-" gorm:"column:password_change"`
	PasswordlessTokens       WebAuthNTokens `json:"-" gorm:"column:passwordless_tokens"`
}

type WebAuthNTokens []*WebAuthNView

type WebAuthNView struct {
	ID    string `json:"webAuthNTokenId"`
	Name  string `json:"webAuthNTokenName,omitempty"`
	State int32  `json:"state,omitempty"`
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
	return h == nil || h.FirstName == ""
}

type MachineView struct {
	Name        string `json:"name" gorm:"column:machine_name"`
	Description string `json:"description" gorm:"column:machine_description"`
}

func (m *MachineView) IsZero() bool {
	return m == nil || m.Name == ""
}

func UserToModel(user *UserView) *model.UserView {
	userView := &model.UserView{
		ID:                 user.ID,
		UserName:           user.UserName,
		ChangeDate:         user.ChangeDate,
		CreationDate:       user.CreationDate,
		ResourceOwner:      user.ResourceOwner,
		State:              model.UserState(user.State),
		LastLogin:          user.LastLogin,
		PreferredLoginName: user.PreferredLoginName,
		LoginNames:         user.LoginNames,
		Sequence:           user.Sequence,
	}
	if !user.HumanView.IsZero() {
		userView.HumanView = &model.HumanView{
			PasswordSet:              user.PasswordSet,
			PasswordInitRequired:     user.PasswordInitRequired,
			PasswordChangeRequired:   user.PasswordChangeRequired,
			PasswordChanged:          user.PasswordChanged,
			PasswordlessTokens:       WebauthnTokensToModel(user.PasswordlessTokens),
			U2FTokens:                WebauthnTokensToModel(user.U2FTokens),
			FirstName:                user.FirstName,
			LastName:                 user.LastName,
			NickName:                 user.NickName,
			DisplayName:              user.DisplayName,
			AvatarKey:                user.AvatarKey,
			PreferredLanguage:        user.PreferredLanguage,
			Gender:                   model.Gender(user.Gender),
			Email:                    user.Email,
			IsEmailVerified:          user.IsEmailVerified,
			Phone:                    user.Phone,
			IsPhoneVerified:          user.IsPhoneVerified,
			Country:                  user.Country,
			Locality:                 user.Locality,
			PostalCode:               user.PostalCode,
			Region:                   user.Region,
			StreetAddress:            user.StreetAddress,
			OTPState:                 model.MFAState(user.OTPState),
			MFAMaxSetUp:              domain.MFALevel(user.MFAMaxSetUp),
			MFAInitSkipped:           user.MFAInitSkipped,
			InitRequired:             user.InitRequired,
			PasswordlessInitRequired: user.PasswordlessInitRequired,
		}
	}

	if !user.MachineView.IsZero() {
		userView.MachineView = &model.MachineView{
			Description: user.MachineView.Description,
			Name:        user.MachineView.Name,
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
		TokenID: token.ID,
		Name:    token.Name,
		State:   model.MFAState(token.State),
	}
}

func (u *UserView) GenerateLoginName(domain string, appendDomain bool) string {
	if !appendDomain {
		return u.UserName
	}
	return u.UserName + "@" + domain
}

func (u *UserView) SetLoginNames(userLoginMustBeDomain bool, domains []*org_model.OrgDomain) {
	u.LoginNames = make([]string, 0, len(domains))
	for _, d := range domains {
		if d.Verified {
			u.LoginNames = append(u.LoginNames, u.GenerateLoginName(d.Domain, true))
		}
	}
	if !userLoginMustBeDomain {
		u.LoginNames = append(u.LoginNames, u.GenerateLoginName(u.UserName, true))
	}
}

func (u *UserView) AppendEvent(event *models.Event) (err error) {
	u.ChangeDate = event.CreationDate
	u.Sequence = event.Sequence
	switch eventstore.EventType(event.Type) {
	case user.MachineAddedEventType:
		u.CreationDate = event.CreationDate
		u.setRootData(event)
		u.Type = userTypeMachine
		err = u.setData(event)
		if err != nil {
			return err
		}
	case user.UserV1AddedType,
		user.UserV1RegisteredType,
		user.HumanRegisteredType,
		user.HumanAddedType:
		u.CreationDate = event.CreationDate
		u.setRootData(event)
		u.Type = userTypeHuman
		err = u.setData(event)
		if err != nil {
			return err
		}
		err = u.setPasswordData(event)
	case user.UserRemovedType:
		u.State = int32(model.UserStateDeleted)
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
		if u.HumanView != nil {
			u.HumanView.UsernameChangeRequired = true
		}
		err = u.setData(event)
	case user.UserUserNameChangedType:
		if u.HumanView != nil {
			u.HumanView.UsernameChangeRequired = false
		}
		err = u.setData(event)
	case user.UserV1EmailChangedType,
		user.HumanEmailChangedType:
		u.IsEmailVerified = false
		err = u.setData(event)
	case user.UserV1EmailVerifiedType,
		user.HumanEmailVerifiedType:
		u.IsEmailVerified = true
	case user.UserV1PhoneChangedType,
		user.HumanPhoneChangedType:
		u.IsPhoneVerified = false
		err = u.setData(event)
	case user.UserV1PhoneVerifiedType,
		user.HumanPhoneVerifiedType:
		u.IsPhoneVerified = true
	case user.UserV1PhoneRemovedType,
		user.HumanPhoneRemovedType:
		u.Phone = ""
		u.IsPhoneVerified = false
	case user.UserDeactivatedType:
		u.State = int32(model.UserStateInactive)
	case user.UserReactivatedType,
		user.UserUnlockedType:
		u.State = int32(model.UserStateActive)
	case user.UserLockedType:
		u.State = int32(model.UserStateLocked)
	case user.UserV1MFAOTPAddedType,
		user.HumanMFAOTPAddedType:
		if u.HumanView == nil {
			logging.WithFields("sequence", event.Sequence, "instance", event.InstanceID).Warn("event is ignored because human not exists")
			return errors.ThrowInvalidArgument(nil, "MODEL-p2BXx", "event ignored: human not exists")
		}
		u.OTPState = int32(model.MFAStateNotReady)
	case user.UserV1MFAOTPVerifiedType,
		user.HumanMFAOTPVerifiedType:
		if u.HumanView == nil {
			logging.WithFields("sequence", event.Sequence, "instance", event.InstanceID).Warn("event is ignored because human not exists")
			return errors.ThrowInvalidArgument(nil, "MODEL-o6Lcq", "event ignored: human not exists")
		}
		u.OTPState = int32(model.MFAStateReady)
		u.MFAInitSkipped = time.Time{}
	case user.UserV1MFAOTPRemovedType,
		user.HumanMFAOTPRemovedType:
		u.OTPState = int32(model.MFAStateUnspecified)
	case user.HumanU2FTokenAddedType:
		err = u.addU2FToken(event)
	case user.HumanU2FTokenVerifiedType:
		err = u.updateU2FToken(event)
		if err != nil {
			return err
		}
		u.MFAInitSkipped = time.Time{}
	case user.HumanU2FTokenRemovedType:
		err = u.removeU2FToken(event)
	case user.UserV1MFAInitSkippedType,
		user.HumanMFAInitSkippedType:
		u.MFAInitSkipped = event.CreationDate
	case user.UserV1InitialCodeAddedType,
		user.HumanInitialCodeAddedType:
		u.InitRequired = true
	case user.UserV1InitializedCheckSucceededType,
		user.HumanInitializedCheckSucceededType:
		u.InitRequired = false
	case user.HumanAvatarAddedType:
		err = u.setData(event)
	case user.HumanAvatarRemovedType:
		u.AvatarKey = ""
	case user.HumanPasswordlessInitCodeAddedType,
		user.HumanPasswordlessInitCodeRequestedType:
		if !u.PasswordSet {
			u.PasswordlessInitRequired = true
			u.PasswordInitRequired = false
		}
	}
	u.ComputeObject()
	return err
}

func (u *UserView) setRootData(event *models.Event) {
	u.ID = event.AggregateID
	u.ResourceOwner = event.ResourceOwner
	u.InstanceID = event.InstanceID
}

func (u *UserView) setData(event *models.Event) error {
	if err := json.Unmarshal(event.Data, u); err != nil {
		logging.Log("MODEL-lso9e").WithError(err).Error("could not unmarshal event data")
		return errors.ThrowInternal(nil, "MODEL-8iows", "could not unmarshal data")
	}
	return nil
}

func (u *UserView) setPasswordData(event *models.Event) error {
	password := new(es_model.Password)
	if err := json.Unmarshal(event.Data, password); err != nil {
		logging.Log("MODEL-sdw4r").WithError(err).Error("could not unmarshal event data")
		return errors.ThrowInternal(nil, "MODEL-6jhsw", "could not unmarshal data")
	}
	u.PasswordSet = password.Secret != nil || password.EncodedHash != ""
	u.PasswordInitRequired = !u.PasswordSet
	u.PasswordChangeRequired = password.ChangeRequired
	u.PasswordChanged = event.CreationDate
	return nil
}

func (u *UserView) addPasswordlessToken(event *models.Event) error {
	token, err := webAuthNViewFromEvent(event)
	if err != nil {
		return err
	}
	for i, t := range u.PasswordlessTokens {
		if t.State == int32(model.MFAStateNotReady) {
			u.PasswordlessTokens[i].ID = token.ID
			return nil
		}
	}
	token.State = int32(model.MFAStateNotReady)
	u.PasswordlessTokens = append(u.PasswordlessTokens, token)
	return nil
}

func (u *UserView) updatePasswordlessToken(event *models.Event) error {
	token, err := webAuthNViewFromEvent(event)
	if err != nil {
		return err
	}
	for i, t := range u.PasswordlessTokens {
		if t.ID == token.ID {
			u.PasswordlessTokens[i].Name = token.Name
			u.PasswordlessTokens[i].State = int32(model.MFAStateReady)
			return nil
		}
	}
	return nil
}

func (u *UserView) removePasswordlessToken(event *models.Event) error {
	token, err := webAuthNViewFromEvent(event)
	if err != nil {
		return err
	}
	for i, t := range u.PasswordlessTokens {
		if t.ID == token.ID {
			u.PasswordlessTokens[i] = u.PasswordlessTokens[len(u.PasswordlessTokens)-1]
			u.PasswordlessTokens[len(u.PasswordlessTokens)-1] = nil
			u.PasswordlessTokens = u.PasswordlessTokens[:len(u.PasswordlessTokens)-1]
			return nil
		}
	}
	return nil
}

func (u *UserView) addU2FToken(event *models.Event) error {
	token, err := webAuthNViewFromEvent(event)
	if err != nil {
		return err
	}
	for i, t := range u.U2FTokens {
		if t.State == int32(model.MFAStateNotReady) {
			u.U2FTokens[i].ID = token.ID
			return nil
		}
	}
	token.State = int32(model.MFAStateNotReady)
	u.U2FTokens = append(u.U2FTokens, token)
	return nil
}

func (u *UserView) updateU2FToken(event *models.Event) error {
	token, err := webAuthNViewFromEvent(event)
	if err != nil {
		return err
	}
	for i, t := range u.U2FTokens {
		if t.ID == token.ID {
			u.U2FTokens[i].Name = token.Name
			u.U2FTokens[i].State = int32(model.MFAStateReady)
			return nil
		}
	}
	return nil
}

func (u *UserView) removeU2FToken(event *models.Event) error {
	token, err := webAuthNViewFromEvent(event)
	if err != nil {
		return err
	}
	for i := len(u.U2FTokens) - 1; i >= 0; i-- {
		if u.U2FTokens[i].ID == token.ID {
			u.U2FTokens[i] = u.U2FTokens[len(u.U2FTokens)-1]
			u.U2FTokens[len(u.U2FTokens)-1] = nil
			u.U2FTokens = u.U2FTokens[:len(u.U2FTokens)-1]
		}
	}
	return nil
}

func webAuthNViewFromEvent(event *models.Event) (*WebAuthNView, error) {
	token := new(WebAuthNView)
	err := json.Unmarshal(event.Data, token)
	if err != nil {
		return nil, errors.ThrowInternal(err, "MODEL-FSaq1", "could not unmarshal data")
	}
	return token, err
}

func (u *UserView) ComputeObject() {
	if !u.MachineView.IsZero() {
		if u.State == int32(model.UserStateUnspecified) {
			u.State = int32(model.UserStateActive)
		}
		return
	}
	if u.State == int32(model.UserStateUnspecified) || u.State == int32(model.UserStateInitial) {
		if u.IsEmailVerified {
			u.State = int32(model.UserStateActive)
		} else {
			u.State = int32(model.UserStateInitial)
		}
	}
	u.ComputeMFAMaxSetUp()
}

func (u *UserView) ComputeMFAMaxSetUp() {
	for _, token := range u.PasswordlessTokens {
		if token.State == int32(model.MFAStateReady) {
			u.MFAMaxSetUp = int32(domain.MFALevelMultiFactor)
			u.PasswordlessInitRequired = false
			return
		}
	}
	for _, token := range u.U2FTokens {
		if token.State == int32(model.MFAStateReady) {
			u.MFAMaxSetUp = int32(domain.MFALevelSecondFactor)
			return
		}
	}
	if u.OTPState == int32(model.MFAStateReady) {
		u.MFAMaxSetUp = int32(domain.MFALevelSecondFactor)
		return
	}
	u.MFAMaxSetUp = int32(domain.MFALevelNotSetUp)
}

func (u *UserView) SetEmptyUserType() {
	if u.MachineView != nil && u.MachineView.Name == "" {
		u.MachineView = nil
	} else {
		u.HumanView = nil
	}
}
