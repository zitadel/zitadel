package command

import (
	"time"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user"
)

type HumanTOTPWriteModel struct {
	eventstore.WriteModel

	State            domain.MFAState
	Secret           *crypto.CryptoValue
	CheckFailedCount uint64
	UserLocked       bool
}

func NewHumanTOTPWriteModel(userID, resourceOwner string) *HumanTOTPWriteModel {
	return &HumanTOTPWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   userID,
			ResourceOwner: resourceOwner,
		},
	}
}

func (wm *HumanTOTPWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *user.HumanOTPAddedEvent:
			wm.Secret = e.Secret
			wm.State = domain.MFAStateNotReady
		case *user.HumanOTPVerifiedEvent:
			wm.State = domain.MFAStateReady
			wm.CheckFailedCount = 0
		case *user.HumanOTPCheckSucceededEvent:
			wm.CheckFailedCount = 0
		case *user.HumanOTPCheckFailedEvent:
			wm.CheckFailedCount++
		case *user.UserLockedEvent:
			wm.UserLocked = true
		case *user.UserUnlockedEvent:
			wm.CheckFailedCount = 0
			wm.UserLocked = false
		case *user.HumanOTPRemovedEvent:
			wm.State = domain.MFAStateRemoved
		case *user.UserRemovedEvent:
			wm.State = domain.MFAStateRemoved
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *HumanTOTPWriteModel) Query() *eventstore.SearchQueryBuilder {
	query := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AddQuery().
		AggregateTypes(user.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(user.HumanMFAOTPAddedType,
			user.HumanMFAOTPVerifiedType,
			user.HumanMFAOTPRemovedType,
			user.HumanMFAOTPCheckSucceededType,
			user.HumanMFAOTPCheckFailedType,
			user.UserLockedType,
			user.UserUnlockedType,
			user.UserRemovedType,
			user.UserV1MFAOTPAddedType,
			user.UserV1MFAOTPVerifiedType,
			user.UserV1MFAOTPRemovedType).
		Builder()

	if wm.ResourceOwner != "" {
		query.ResourceOwner(wm.ResourceOwner)
	}
	return query
}

type OTPWriteModel interface {
	OTPAdded() bool
	ResourceOwner() string
}

type OTPCodeWriteModel interface {
	OTPWriteModel
	CodeCreationDate() time.Time
	CodeExpiry() time.Duration
	Code() *crypto.CryptoValue
	CheckFailedCount() uint64
	UserLocked() bool
	eventstore.QueryReducer
}

type HumanOTPSMSWriteModel struct {
	eventstore.WriteModel

	phoneVerified bool
	otpAdded      bool
}

func (wm *HumanOTPSMSWriteModel) OTPAdded() bool {
	return wm.otpAdded
}

func (wm *HumanOTPSMSWriteModel) ResourceOwner() string {
	return wm.WriteModel.ResourceOwner
}

func NewHumanOTPSMSWriteModel(userID, resourceOwner string) *HumanOTPSMSWriteModel {
	return &HumanOTPSMSWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   userID,
			ResourceOwner: resourceOwner,
		},
	}
}

func (wm *HumanOTPSMSWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch event.(type) {
		case *user.HumanPhoneVerifiedEvent:
			wm.phoneVerified = true
		case *user.HumanOTPSMSAddedEvent:
			wm.otpAdded = true
		case *user.HumanOTPSMSRemovedEvent:
			wm.otpAdded = false
		case *user.HumanPhoneRemovedEvent,
			*user.UserRemovedEvent:
			wm.phoneVerified = false
			wm.otpAdded = false
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *HumanOTPSMSWriteModel) Query() *eventstore.SearchQueryBuilder {
	query := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AddQuery().
		AggregateTypes(user.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(user.HumanPhoneVerifiedType,
			user.HumanOTPSMSAddedType,
			user.HumanOTPSMSRemovedType,
			user.HumanPhoneRemovedType,
			user.UserRemovedType,
		).
		Builder()

	if wm.WriteModel.ResourceOwner != "" {
		query.ResourceOwner(wm.WriteModel.ResourceOwner)
	}
	return query
}

type HumanOTPSMSCodeWriteModel struct {
	*HumanOTPSMSWriteModel

	otpCode *OTPCode

	checkFailedCount uint64
	userLocked       bool
}

func (wm *HumanOTPSMSCodeWriteModel) CodeCreationDate() time.Time {
	if wm.otpCode == nil {
		return time.Time{}
	}
	return wm.otpCode.CreationDate
}

func (wm *HumanOTPSMSCodeWriteModel) CodeExpiry() time.Duration {
	if wm.otpCode == nil {
		return 0
	}
	return wm.otpCode.Expiry
}

func (wm *HumanOTPSMSCodeWriteModel) Code() *crypto.CryptoValue {
	if wm.otpCode == nil {
		return nil
	}
	return wm.otpCode.Code
}

func (wm *HumanOTPSMSCodeWriteModel) CheckFailedCount() uint64 {
	return wm.checkFailedCount
}

func (wm *HumanOTPSMSCodeWriteModel) UserLocked() bool {
	return wm.userLocked
}

func NewHumanOTPSMSCodeWriteModel(userID, resourceOwner string) *HumanOTPSMSCodeWriteModel {
	return &HumanOTPSMSCodeWriteModel{
		HumanOTPSMSWriteModel: NewHumanOTPSMSWriteModel(userID, resourceOwner),
	}
}

func (wm *HumanOTPSMSCodeWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *user.HumanOTPSMSCodeAddedEvent:
			wm.otpCode = &OTPCode{
				Code:         e.Code,
				CreationDate: e.CreationDate(),
				Expiry:       e.Expiry,
			}
		case *user.HumanOTPSMSCheckSucceededEvent:
			wm.checkFailedCount = 0
		case *user.HumanOTPSMSCheckFailedEvent:
			wm.checkFailedCount++
		case *user.UserLockedEvent:
			wm.userLocked = true
		case *user.UserUnlockedEvent:
			wm.checkFailedCount = 0
			wm.userLocked = false
		}
	}
	return wm.HumanOTPSMSWriteModel.Reduce()
}

func (wm *HumanOTPSMSCodeWriteModel) Query() *eventstore.SearchQueryBuilder {
	query := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AddQuery().
		AggregateTypes(user.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			user.HumanOTPSMSCodeAddedType,
			user.HumanOTPSMSCheckSucceededType,
			user.HumanOTPSMSCheckFailedType,
			user.UserLockedType,
			user.UserUnlockedType,
			user.HumanPhoneVerifiedType,
			user.HumanOTPSMSAddedType,
			user.HumanOTPSMSRemovedType,
			user.HumanPhoneRemovedType,
			user.UserRemovedType,
		).
		Builder()

	if wm.WriteModel.ResourceOwner != "" {
		query.ResourceOwner(wm.WriteModel.ResourceOwner)
	}
	return query
}

type HumanOTPEmailWriteModel struct {
	eventstore.WriteModel

	emailVerified bool
	otpAdded      bool
}

func (wm *HumanOTPEmailWriteModel) OTPAdded() bool {
	return wm.otpAdded
}

func (wm *HumanOTPEmailWriteModel) ResourceOwner() string {
	return wm.WriteModel.ResourceOwner
}

func NewHumanOTPEmailWriteModel(userID, resourceOwner string) *HumanOTPEmailWriteModel {
	return &HumanOTPEmailWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   userID,
			ResourceOwner: resourceOwner,
		},
	}
}

func (wm *HumanOTPEmailWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch event.(type) {
		case *user.HumanEmailVerifiedEvent:
			wm.emailVerified = true
		case *user.HumanOTPEmailAddedEvent:
			wm.otpAdded = true
		case *user.HumanOTPEmailRemovedEvent:
			wm.otpAdded = false
		case *user.UserRemovedEvent:
			wm.emailVerified = false
			wm.otpAdded = false
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *HumanOTPEmailWriteModel) Query() *eventstore.SearchQueryBuilder {
	query := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AddQuery().
		AggregateTypes(user.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			user.HumanEmailVerifiedType,
			user.HumanOTPEmailAddedType,
			user.HumanOTPEmailRemovedType,
			user.UserRemovedType,
		).
		Builder()

	if wm.WriteModel.ResourceOwner != "" {
		query.ResourceOwner(wm.WriteModel.ResourceOwner)
	}
	return query
}

type HumanOTPEmailCodeWriteModel struct {
	*HumanOTPEmailWriteModel

	otpCode *OTPCode

	checkFailedCount uint64
	userLocked       bool
}

func (wm *HumanOTPEmailCodeWriteModel) CodeCreationDate() time.Time {
	if wm.otpCode == nil {
		return time.Time{}
	}
	return wm.otpCode.CreationDate
}

func (wm *HumanOTPEmailCodeWriteModel) CodeExpiry() time.Duration {
	if wm.otpCode == nil {
		return 0
	}
	return wm.otpCode.Expiry
}

func (wm *HumanOTPEmailCodeWriteModel) Code() *crypto.CryptoValue {
	if wm.otpCode == nil {
		return nil
	}
	return wm.otpCode.Code
}

func (wm *HumanOTPEmailCodeWriteModel) CheckFailedCount() uint64 {
	return wm.checkFailedCount
}

func (wm *HumanOTPEmailCodeWriteModel) UserLocked() bool {
	return wm.userLocked
}

func NewHumanOTPEmailCodeWriteModel(userID, resourceOwner string) *HumanOTPEmailCodeWriteModel {
	return &HumanOTPEmailCodeWriteModel{
		HumanOTPEmailWriteModel: NewHumanOTPEmailWriteModel(userID, resourceOwner),
	}
}

func (wm *HumanOTPEmailCodeWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *user.HumanOTPEmailCodeAddedEvent:
			wm.otpCode = &OTPCode{
				Code:         e.Code,
				CreationDate: e.CreationDate(),
				Expiry:       e.Expiry,
			}
		case *user.HumanOTPEmailCheckSucceededEvent:
			wm.checkFailedCount = 0
		case *user.HumanOTPEmailCheckFailedEvent:
			wm.checkFailedCount++
		case *user.UserLockedEvent:
			wm.userLocked = true
		case *user.UserUnlockedEvent:
			wm.checkFailedCount = 0
			wm.userLocked = false
		}
	}
	return wm.HumanOTPEmailWriteModel.Reduce()
}

func (wm *HumanOTPEmailCodeWriteModel) Query() *eventstore.SearchQueryBuilder {
	query := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AddQuery().
		AggregateTypes(user.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			user.HumanOTPEmailCodeAddedType,
			user.HumanOTPEmailCheckSucceededType,
			user.HumanOTPEmailCheckFailedType,
			user.UserLockedType,
			user.UserUnlockedType,
			user.HumanEmailVerifiedType,
			user.HumanOTPEmailAddedType,
			user.HumanOTPEmailRemovedType,
			user.UserRemovedType,
		).
		Builder()

	if wm.WriteModel.ResourceOwner != "" {
		query.ResourceOwner(wm.WriteModel.ResourceOwner)
	}
	return query
}
