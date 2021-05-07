package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/caos/logging"
	"github.com/lib/pq"

	req_model "github.com/caos/zitadel/internal/auth_request/model"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	org_model "github.com/caos/zitadel/internal/org/model"
	"github.com/caos/zitadel/internal/user/model"
	es_model "github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
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
)

type userType string

const (
	userTypeHuman   = "human"
	userTypeMachine = "machine"
)

type UserView struct {
	ID                 string         `json:"-" gorm:"column:id;primary_key"`
	CreationDate       time.Time      `json:"-" gorm:"column:creation_date"`
	ChangeDate         time.Time      `json:"-" gorm:"column:change_date"`
	ResourceOwner      string         `json:"-" gorm:"column:resource_owner"`
	State              int32          `json:"-" gorm:"column:user_state"`
	LastLogin          time.Time      `json:"-" gorm:"column:last_login"`
	LoginNames         pq.StringArray `json:"-" gorm:"column:login_names"`
	PreferredLoginName string         `json:"-" gorm:"column:preferred_login_name"`
	Sequence           uint64         `json:"-" gorm:"column:sequence"`
	Type               userType       `json:"-" gorm:"column:user_type"`
	UserName           string         `json:"userName" gorm:"column:user_name"`
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
	FirstName         string         `json:"firstName" gorm:"column:first_name"`
	LastName          string         `json:"lastName" gorm:"column:last_name"`
	NickName          string         `json:"nickName" gorm:"column:nick_name"`
	DisplayName       string         `json:"displayName" gorm:"column:display_name"`
	PreferredLanguage string         `json:"preferredLanguage" gorm:"column:preferred_language"`
	Gender            int32          `json:"gender" gorm:"column:gender"`
	AvatarKey         string         `json:"storeKey" gorm:"column:avatar"`
	Email             string         `json:"email" gorm:"column:email"`
	IsEmailVerified   bool           `json:"-" gorm:"column:is_email_verified"`
	Phone             string         `json:"phone" gorm:"column:phone"`
	IsPhoneVerified   bool           `json:"-" gorm:"column:is_phone_verified"`
	Country           string         `json:"country" gorm:"column:country"`
	Locality          string         `json:"locality" gorm:"column:locality"`
	PostalCode        string         `json:"postalCode" gorm:"column:postal_code"`
	Region            string         `json:"region" gorm:"column:region"`
	StreetAddress     string         `json:"streetAddress" gorm:"column:street_address"`
	OTPState          int32          `json:"-" gorm:"column:otp_state"`
	U2FTokens         WebAuthNTokens `json:"-" gorm:"column:u2f_tokens"`
	MFAMaxSetUp       int32          `json:"-" gorm:"column:mfa_max_set_up"`
	MFAInitSkipped    time.Time      `json:"-" gorm:"column:mfa_init_skipped"`
	InitRequired      bool           `json:"-" gorm:"column:init_required"`

	PasswordSet            bool           `json:"-" gorm:"column:password_set"`
	PasswordChangeRequired bool           `json:"-" gorm:"column:password_change_required"`
	UsernameChangeRequired bool           `json:"-" gorm:"column:username_change_required"`
	PasswordChanged        time.Time      `json:"-" gorm:"column:password_change"`
	PasswordlessTokens     WebAuthNTokens `json:"-" gorm:"column:passwordless_tokens"`
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
			PasswordSet:            user.PasswordSet,
			PasswordChangeRequired: user.PasswordChangeRequired,
			PasswordChanged:        user.PasswordChanged,
			PasswordlessTokens:     WebauthnTokensToModel(user.PasswordlessTokens),
			U2FTokens:              WebauthnTokensToModel(user.U2FTokens),
			FirstName:              user.FirstName,
			LastName:               user.LastName,
			NickName:               user.NickName,
			DisplayName:            user.DisplayName,
			AvatarKey:              user.AvatarKey,
			PreferredLanguage:      user.PreferredLanguage,
			Gender:                 model.Gender(user.Gender),
			Email:                  user.Email,
			IsEmailVerified:        user.IsEmailVerified,
			Phone:                  user.Phone,
			IsPhoneVerified:        user.IsPhoneVerified,
			Country:                user.Country,
			Locality:               user.Locality,
			PostalCode:             user.PostalCode,
			Region:                 user.Region,
			StreetAddress:          user.StreetAddress,
			OTPState:               model.MFAState(user.OTPState),
			MFAMaxSetUp:            req_model.MFALevel(user.MFAMaxSetUp),
			MFAInitSkipped:         user.MFAInitSkipped,
			InitRequired:           user.InitRequired,
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

func UsersToModel(users []*UserView) []*model.UserView {
	result := make([]*model.UserView, len(users))
	for i, p := range users {
		result[i] = UserToModel(p)
	}
	return result
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

func (u *UserView) SetLoginNames(policy *iam_model.OrgIAMPolicy, domains []*org_model.OrgDomain) {
	loginNames := make([]string, 0)
	for _, d := range domains {
		if d.Verified {
			loginNames = append(loginNames, u.GenerateLoginName(d.Domain, true))
		}
	}
	if !policy.UserLoginMustBeDomain {
		loginNames = append(loginNames, u.UserName)
	}
	u.LoginNames = loginNames
}

func (u *UserView) AppendEvent(event *models.Event) (err error) {
	u.ChangeDate = event.CreationDate
	u.Sequence = event.Sequence
	switch event.Type {
	case es_model.MachineAdded:
		u.CreationDate = event.CreationDate
		u.setRootData(event)
		u.Type = userTypeMachine
		err = u.setData(event)
		if err != nil {
			return err
		}
	case es_model.UserAdded,
		es_model.UserRegistered,
		es_model.HumanRegistered,
		es_model.HumanAdded:
		u.CreationDate = event.CreationDate
		u.setRootData(event)
		u.Type = userTypeHuman
		err = u.setData(event)
		if err != nil {
			return err
		}
		err = u.setPasswordData(event)
	case es_model.UserRemoved:
		u.State = int32(model.UserStateDeleted)
	case es_model.UserPasswordChanged,
		es_model.HumanPasswordChanged:
		err = u.setPasswordData(event)
	case es_model.HumanPasswordlessTokenAdded:
		err = u.addPasswordlessToken(event)
	case es_model.HumanPasswordlessTokenVerified:
		err = u.updatePasswordlessToken(event)
	case es_model.HumanPasswordlessTokenRemoved:
		err = u.removePasswordlessToken(event)
	case es_model.UserProfileChanged,
		es_model.HumanProfileChanged,
		es_model.UserAddressChanged,
		es_model.HumanAddressChanged,
		es_model.MachineChanged:
		err = u.setData(event)
	case es_model.DomainClaimed:
		if u.HumanView != nil {
			u.HumanView.UsernameChangeRequired = true
		}
		err = u.setData(event)
	case es_model.UserUserNameChanged:
		if u.HumanView != nil {
			u.HumanView.UsernameChangeRequired = false
		}
		err = u.setData(event)
	case es_model.UserEmailChanged,
		es_model.HumanEmailChanged:
		u.IsEmailVerified = false
		err = u.setData(event)
	case es_model.UserEmailVerified,
		es_model.HumanEmailVerified:
		u.IsEmailVerified = true
	case es_model.UserPhoneChanged,
		es_model.HumanPhoneChanged:
		u.IsPhoneVerified = false
		err = u.setData(event)
	case es_model.UserPhoneVerified,
		es_model.HumanPhoneVerified:
		u.IsPhoneVerified = true
	case es_model.UserPhoneRemoved,
		es_model.HumanPhoneRemoved:
		u.Phone = ""
		u.IsPhoneVerified = false
	case es_model.UserDeactivated:
		u.State = int32(model.UserStateInactive)
	case es_model.UserReactivated,
		es_model.UserUnlocked:
		u.State = int32(model.UserStateActive)
	case es_model.UserLocked:
		u.State = int32(model.UserStateLocked)
	case es_model.MFAOTPAdded,
		es_model.HumanMFAOTPAdded:
		u.OTPState = int32(model.MFAStateNotReady)
	case es_model.MFAOTPVerified,
		es_model.HumanMFAOTPVerified:
		u.OTPState = int32(model.MFAStateReady)
		u.MFAInitSkipped = time.Time{}
	case es_model.MFAOTPRemoved,
		es_model.HumanMFAOTPRemoved:
		u.OTPState = int32(model.MFAStateUnspecified)
	case es_model.HumanMFAU2FTokenAdded:
		err = u.addU2FToken(event)
	case es_model.HumanMFAU2FTokenVerified:
		err = u.updateU2FToken(event)
		if err != nil {
			return err
		}
		u.MFAInitSkipped = time.Time{}
	case es_model.HumanMFAU2FTokenRemoved:
		err = u.removeU2FToken(event)
	case es_model.MFAInitSkipped,
		es_model.HumanMFAInitSkipped:
		u.MFAInitSkipped = event.CreationDate
	case es_model.InitializedUserCodeAdded,
		es_model.InitializedHumanCodeAdded:
		u.InitRequired = true
	case es_model.InitializedUserCheckSucceeded,
		es_model.InitializedHumanCheckSucceeded:
		u.InitRequired = false
	case es_model.HumanAvatarAdded:
		u.setData(event)
	case es_model.HumanAvatarRemoved:
		u.AvatarKey = ""
	}
	u.ComputeObject()
	return err
}

func (u *UserView) setRootData(event *models.Event) {
	u.ID = event.AggregateID
	u.ResourceOwner = event.ResourceOwner
}

func (u *UserView) setData(event *models.Event) error {
	if err := json.Unmarshal(event.Data, u); err != nil {
		logging.Log("MODEL-lso9e").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(nil, "MODEL-8iows", "could not unmarshal data")
	}
	return nil
}

func (u *UserView) setPasswordData(event *models.Event) error {
	password := new(es_model.Password)
	if err := json.Unmarshal(event.Data, password); err != nil {
		logging.Log("MODEL-sdw4r").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(nil, "MODEL-6jhsw", "could not unmarshal data")
	}
	u.PasswordSet = password.Secret != nil
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
		return nil, caos_errs.ThrowInternal(err, "MODEL-FSaq1", "could not unmarshal data")
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
			u.MFAMaxSetUp = int32(req_model.MFALevelMultiFactor)
			return
		}
	}
	for _, token := range u.U2FTokens {
		if token.State == int32(model.MFAStateReady) {
			u.MFAMaxSetUp = int32(req_model.MFALevelSecondFactor)
			return
		}
	}
	if u.OTPState == int32(model.MFAStateReady) {
		u.MFAMaxSetUp = int32(req_model.MFALevelSecondFactor)
		return
	}
	u.MFAMaxSetUp = int32(req_model.MFALevelNotSetUp)
}

func (u *UserView) SetEmptyUserType() {
	if u.MachineView != nil && u.MachineView.Name == "" {
		u.MachineView = nil
	} else {
		u.HumanView = nil
	}
}
