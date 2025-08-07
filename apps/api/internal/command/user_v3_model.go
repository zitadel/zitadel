package command

import (
	"bytes"
	"context"
	"encoding/json"
	"time"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	domain_schema "github.com/zitadel/zitadel/internal/domain/schema"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user/schemauser"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type UserV3WriteModel struct {
	eventstore.WriteModel

	PhoneWM bool
	EmailWM bool
	DataWM  bool

	SchemaID       string
	SchemaRevision uint64

	Email                    string
	IsEmailVerified          bool
	EmailVerifiedFailedCount int
	EmailCode                *VerifyCode

	Phone                    string
	IsPhoneVerified          bool
	PhoneVerifiedFailedCount int
	PhoneCode                *VerifyCode

	Data json.RawMessage

	Locked bool
	State  domain.UserState

	checkPermission      domain.PermissionCheck
	writePermissionCheck bool
}

func (wm *UserV3WriteModel) GetWriteModel() *eventstore.WriteModel {
	return &wm.WriteModel
}

type VerifyCode struct {
	Code         *crypto.CryptoValue
	CreationDate time.Time
	Expiry       time.Duration
}

func NewExistsUserV3WriteModel(resourceOwner, userID string, checkPermission domain.PermissionCheck) *UserV3WriteModel {
	return &UserV3WriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   userID,
			ResourceOwner: resourceOwner,
		},
		PhoneWM:         false,
		EmailWM:         false,
		DataWM:          false,
		checkPermission: checkPermission,
	}
}

func NewUserV3WriteModel(resourceOwner, userID string, checkPermission domain.PermissionCheck) *UserV3WriteModel {
	return &UserV3WriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   userID,
			ResourceOwner: resourceOwner,
		},
		PhoneWM:         true,
		EmailWM:         true,
		DataWM:          true,
		checkPermission: checkPermission,
	}
}

func NewUserV3EmailWriteModel(resourceOwner, userID string, checkPermission domain.PermissionCheck) *UserV3WriteModel {
	return &UserV3WriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   userID,
			ResourceOwner: resourceOwner,
		},
		EmailWM:         true,
		checkPermission: checkPermission,
	}
}

func NewUserV3PhoneWriteModel(resourceOwner, userID string, checkPermission domain.PermissionCheck) *UserV3WriteModel {
	return &UserV3WriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   userID,
			ResourceOwner: resourceOwner,
		},
		PhoneWM:         true,
		checkPermission: checkPermission,
	}
}

func (wm *UserV3WriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *schemauser.CreatedEvent:
			wm.SchemaID = e.SchemaID
			wm.SchemaRevision = e.SchemaRevision
			wm.Data = e.Data
			wm.Locked = false

			wm.State = domain.UserStateActive
		case *schemauser.UpdatedEvent:
			if e.SchemaID != nil {
				wm.SchemaID = *e.SchemaID
			}
			if e.SchemaRevision != nil {
				wm.SchemaRevision = *e.SchemaRevision
			}
			if len(e.Data) > 0 {
				wm.Data = e.Data
			}
		case *schemauser.DeletedEvent:
			wm.State = domain.UserStateDeleted
		case *schemauser.EmailUpdatedEvent:
			wm.Email = string(e.EmailAddress)
			wm.IsEmailVerified = false
			wm.EmailVerifiedFailedCount = 0
			wm.EmailCode = nil
		case *schemauser.EmailCodeAddedEvent:
			wm.IsEmailVerified = false
			wm.EmailVerifiedFailedCount = 0
			wm.EmailCode = &VerifyCode{
				Code:         e.Code,
				CreationDate: e.CreationDate(),
				Expiry:       e.Expiry,
			}
		case *schemauser.EmailVerifiedEvent:
			wm.IsEmailVerified = true
			wm.EmailVerifiedFailedCount = 0
			wm.EmailCode = nil
		case *schemauser.EmailVerificationFailedEvent:
			wm.EmailVerifiedFailedCount += 1
		case *schemauser.PhoneUpdatedEvent:
			wm.Phone = string(e.PhoneNumber)
			wm.IsPhoneVerified = false
			wm.PhoneVerifiedFailedCount = 0
			wm.EmailCode = nil
		case *schemauser.PhoneCodeAddedEvent:
			wm.IsPhoneVerified = false
			wm.PhoneVerifiedFailedCount = 0
			wm.PhoneCode = &VerifyCode{
				Code:         e.Code,
				CreationDate: e.CreationDate(),
				Expiry:       e.Expiry,
			}
		case *schemauser.PhoneVerifiedEvent:
			wm.PhoneVerifiedFailedCount = 0
			wm.IsPhoneVerified = true
			wm.PhoneCode = nil
		case *schemauser.PhoneVerificationFailedEvent:
			wm.PhoneVerifiedFailedCount += 1
		case *schemauser.LockedEvent:
			wm.Locked = true
		case *schemauser.UnlockedEvent:
			wm.Locked = false
		case *schemauser.DeactivatedEvent:
			wm.State = domain.UserStateInactive
		case *schemauser.ActivatedEvent:
			wm.State = domain.UserStateActive
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *UserV3WriteModel) Query() *eventstore.SearchQueryBuilder {
	builder := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent)
	if wm.ResourceOwner != "" {
		builder = builder.ResourceOwner(wm.ResourceOwner)
	}
	eventtypes := []eventstore.EventType{
		schemauser.CreatedType,
		schemauser.DeletedType,
		schemauser.ActivatedType,
		schemauser.DeactivatedType,
		schemauser.LockedType,
		schemauser.UnlockedType,
	}
	if wm.DataWM {
		eventtypes = append(eventtypes,
			schemauser.UpdatedType,
		)
	}
	if wm.EmailWM {
		eventtypes = append(eventtypes,
			schemauser.EmailUpdatedType,
			schemauser.EmailVerifiedType,
			schemauser.EmailCodeAddedType,
			schemauser.EmailVerificationFailedType,
		)
	}
	if wm.PhoneWM {
		eventtypes = append(eventtypes,
			schemauser.PhoneUpdatedType,
			schemauser.PhoneVerifiedType,
			schemauser.PhoneCodeAddedType,
			schemauser.PhoneVerificationFailedType,
		)
	}
	return builder.AddQuery().
		AggregateTypes(schemauser.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(eventtypes...).Builder()
}

func (wm *UserV3WriteModel) NewCreated(
	ctx context.Context,
	schemaID string,
	schemaRevision uint64,
	data json.RawMessage,
	email *Email,
	phone *Phone,
	emailCode func(context.Context) (*EncryptedCode, error),
	phoneCode func(context.Context) (*EncryptedCode, string, error),
) (_ []eventstore.Command, codeEmail string, codePhone string, err error) {
	if err := wm.checkPermissionWrite(ctx, wm.ResourceOwner, wm.AggregateID); err != nil {
		return nil, "", "", err
	}
	if wm.Exists() {
		return nil, "", "", zerrors.ThrowPreconditionFailed(nil, "COMMAND-Nn8CRVlkeZ", "Errors.User.AlreadyExists")
	}
	events := []eventstore.Command{
		schemauser.NewCreatedEvent(ctx,
			UserV3AggregateFromWriteModel(&wm.WriteModel),
			schemaID, schemaRevision, data,
		),
	}
	if email != nil {
		emailEvents, plainCodeEmail, err := wm.NewEmailCreate(ctx,
			email,
			emailCode,
		)
		if err != nil {
			return nil, "", "", err
		}
		if plainCodeEmail != "" {
			codeEmail = plainCodeEmail
		}
		events = append(events, emailEvents...)
	}

	if phone != nil {
		phoneEvents, plainCodePhone, err := wm.NewPhoneCreate(ctx,
			phone,
			phoneCode,
		)
		if err != nil {
			return nil, "", "", err
		}
		if plainCodePhone != "" {
			codePhone = plainCodePhone
		}
		events = append(events, phoneEvents...)
	}

	return events, codeEmail, codePhone, nil
}

func (wm *UserV3WriteModel) getSchemaRoleForWrite(ctx context.Context, resourceOwner, userID string) (domain_schema.Role, error) {
	if userID == authz.GetCtxData(ctx).UserID {
		return domain_schema.RoleSelf, nil
	}
	if err := wm.checkPermission(ctx, domain.PermissionUserWrite, resourceOwner, userID); err != nil {
		return domain_schema.RoleUnspecified, err
	}
	return domain_schema.RoleOwner, nil
}

func (wm *UserV3WriteModel) validateData(ctx context.Context, data []byte, schemaWM *UserSchemaWriteModel) (string, uint64, error) {
	// get role for permission check in schema through extension
	role, err := wm.getSchemaRoleForWrite(ctx, wm.ResourceOwner, wm.AggregateID)
	if err != nil {
		return "", 0, err
	}

	schema, err := domain_schema.NewSchema(role, bytes.NewReader(schemaWM.Schema))
	if err != nil {
		return "", 0, err
	}

	// if data not changed but a new schema or revision should be used
	if data == nil {
		data = wm.Data
	}
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return "", 0, zerrors.ThrowInvalidArgument(nil, "COMMAND-7o3ZGxtXUz", "Errors.User.Invalid")
	}

	if err := schema.Validate(v); err != nil {
		return "", 0, zerrors.ThrowPreconditionFailed(nil, "COMMAND-SlKXqLSeL6", "Errors.UserSchema.Data.Invalid")
	}
	return schemaWM.AggregateID, schemaWM.SchemaRevision, nil
}

func (wm *UserV3WriteModel) NewUpdate(
	ctx context.Context,
	schemaWM *UserSchemaWriteModel,
	user *SchemaUser,
	email *Email,
	phone *Phone,
	emailCode func(context.Context) (*EncryptedCode, error),
	phoneCode func(context.Context) (*EncryptedCode, string, error),
) (_ []eventstore.Command, codeEmail string, codePhone string, err error) {
	if err := wm.checkPermissionWrite(ctx, wm.ResourceOwner, wm.AggregateID); err != nil {
		return nil, "", "", err
	}
	if !wm.Exists() {
		return nil, "", "", zerrors.ThrowPreconditionFailed(nil, "COMMAND-Nn8CRVlkeZ", "Errors.User.NotFound")
	}
	events := make([]eventstore.Command, 0)
	if user != nil {
		schemaID, schemaRevision, err := wm.validateData(ctx, user.Data, schemaWM)
		if err != nil {
			return nil, "", "", err
		}
		userEvents := wm.newUpdatedEvents(ctx,
			schemaID,
			schemaRevision,
			user.Data,
		)
		events = append(events, userEvents...)
	}
	if email != nil {
		emailEvents, plainCodeEmail, err := wm.NewEmailUpdate(ctx,
			email,
			emailCode,
		)
		if err != nil {
			return nil, "", "", err
		}
		if plainCodeEmail != "" {
			codeEmail = plainCodeEmail
		}
		events = append(events, emailEvents...)
	}

	if phone != nil {
		phoneEvents, plainCodePhone, err := wm.NewPhoneCreate(ctx,
			phone,
			phoneCode,
		)
		if err != nil {
			return nil, "", "", err
		}
		if plainCodePhone != "" {
			codePhone = plainCodePhone
		}
		events = append(events, phoneEvents...)
	}

	return events, codeEmail, codePhone, nil
}

func (wm *UserV3WriteModel) newUpdatedEvents(
	ctx context.Context,
	schemaID string,
	schemaRevision uint64,
	data json.RawMessage,
) []eventstore.Command {
	changes := make([]schemauser.Changes, 0)
	if wm.SchemaID != schemaID {
		changes = append(changes, schemauser.ChangeSchemaID(schemaID))
	}
	if wm.SchemaRevision != schemaRevision {
		changes = append(changes, schemauser.ChangeSchemaRevision(schemaRevision))
	}
	if data != nil && !bytes.Equal(wm.Data, data) {
		changes = append(changes, schemauser.ChangeData(data))
	}
	if len(changes) == 0 {
		return nil
	}
	return []eventstore.Command{schemauser.NewUpdatedEvent(ctx, UserV3AggregateFromWriteModel(&wm.WriteModel), changes)}
}

func (wm *UserV3WriteModel) NewDelete(
	ctx context.Context,
) (_ []eventstore.Command, err error) {
	if !wm.Exists() {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-syHyCsGmvM", "Errors.User.NotFound")
	}
	if err := wm.checkPermissionDelete(ctx, wm.ResourceOwner, wm.AggregateID); err != nil {
		return nil, err
	}
	return []eventstore.Command{schemauser.NewDeletedEvent(ctx, UserV3AggregateFromWriteModel(&wm.WriteModel))}, nil

}

func UserV3AggregateFromWriteModel(wm *eventstore.WriteModel) *eventstore.Aggregate {
	return &eventstore.Aggregate{
		ID:            wm.AggregateID,
		Type:          schemauser.AggregateType,
		ResourceOwner: wm.ResourceOwner,
		InstanceID:    wm.InstanceID,
		Version:       schemauser.AggregateVersion,
	}
}

func (wm *UserV3WriteModel) Exists() bool {
	return wm.State != domain.UserStateDeleted && wm.State != domain.UserStateUnspecified
}

func (wm *UserV3WriteModel) checkPermissionWrite(
	ctx context.Context,
	resourceOwner string,
	userID string,
) error {
	if wm.writePermissionCheck {
		return nil
	}
	if userID != "" && userID == authz.GetCtxData(ctx).UserID {
		return nil
	}
	if err := wm.checkPermission(ctx, domain.PermissionUserWrite, resourceOwner, userID); err != nil {
		return err
	}
	wm.writePermissionCheck = true
	return nil
}

func (wm *UserV3WriteModel) checkPermissionDelete(
	ctx context.Context,
	resourceOwner string,
	userID string,
) error {
	if userID != "" && userID == authz.GetCtxData(ctx).UserID {
		return nil
	}
	return wm.checkPermission(ctx, domain.PermissionUserDelete, resourceOwner, userID)
}

func (wm *UserV3WriteModel) NewEmailCreate(
	ctx context.Context,
	email *Email,
	code func(context.Context) (*EncryptedCode, error),
) (_ []eventstore.Command, plainCode string, err error) {
	if err := wm.checkPermissionWrite(ctx, wm.ResourceOwner, wm.AggregateID); err != nil {
		return nil, "", err
	}
	if email == nil || wm.Email == string(email.Address) {
		return nil, "", nil
	}
	events := []eventstore.Command{
		schemauser.NewEmailUpdatedEvent(ctx,
			UserV3AggregateFromWriteModel(&wm.WriteModel),
			email.Address,
		),
	}
	if email.Verified {
		events = append(events, wm.newEmailVerifiedEvent(ctx))
	} else {
		codeEvent, code, err := wm.newEmailCodeAddedEvent(ctx, code, email.URLTemplate, email.ReturnCode)
		if err != nil {
			return nil, "", err
		}
		events = append(events, codeEvent)
		if code != "" {
			plainCode = code
		}
	}
	return events, plainCode, nil
}

func (wm *UserV3WriteModel) NewEmailUpdate(
	ctx context.Context,
	email *Email,
	code func(context.Context) (*EncryptedCode, error),
) (_ []eventstore.Command, plainCode string, err error) {
	if !wm.EmailWM {
		return nil, "", nil
	}
	if !wm.Exists() {
		return nil, "", zerrors.ThrowNotFound(nil, "COMMAND-nJ0TQFuRmP", "Errors.User.NotFound")
	}
	return wm.NewEmailCreate(ctx, email, code)
}

func (wm *UserV3WriteModel) NewEmailVerify(
	ctx context.Context,
	verify func(creationDate time.Time, expiry time.Duration, cryptoCode *crypto.CryptoValue) error,
) ([]eventstore.Command, error) {
	if !wm.EmailWM {
		return nil, nil
	}
	if !wm.Exists() {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-qbGyMPvjvj", "Errors.User.NotFound")
	}
	if err := wm.checkPermissionWrite(ctx, wm.ResourceOwner, wm.AggregateID); err != nil {
		return nil, err
	}
	if wm.EmailCode == nil {
		return nil, nil
	}
	if err := verify(wm.EmailCode.CreationDate, wm.EmailCode.Expiry, wm.EmailCode.Code); err != nil {
		return nil, err
	}
	return []eventstore.Command{wm.newEmailVerifiedEvent(ctx)}, nil
}

func (wm *UserV3WriteModel) newEmailVerifiedEvent(
	ctx context.Context,
) *schemauser.EmailVerifiedEvent {
	return schemauser.NewEmailVerifiedEvent(ctx, UserV3AggregateFromWriteModel(&wm.WriteModel))
}

func (wm *UserV3WriteModel) NewResendEmailCode(
	ctx context.Context,
	code func(context.Context) (*EncryptedCode, error),
	urlTemplate string,
	isReturnCode bool,
) (_ []eventstore.Command, plainCode string, err error) {
	if !wm.EmailWM {
		return nil, "", nil
	}
	if !wm.Exists() {
		return nil, "", zerrors.ThrowNotFound(nil, "COMMAND-EajeF6ypOV", "Errors.User.NotFound")
	}
	if err := wm.checkPermissionWrite(ctx, wm.ResourceOwner, wm.AggregateID); err != nil {
		return nil, "", err
	}
	if wm.EmailCode == nil {
		return nil, "", zerrors.ThrowPreconditionFailed(err, "COMMAND-QRkNTBwF8q", "Errors.User.Code.Empty")
	}
	event, plainCode, err := wm.newEmailCodeAddedEvent(ctx, code, urlTemplate, isReturnCode)
	if err != nil {
		return nil, "", err
	}
	return []eventstore.Command{event}, plainCode, nil
}

func (wm *UserV3WriteModel) newEmailCodeAddedEvent(
	ctx context.Context,
	code func(context.Context) (*EncryptedCode, error),
	urlTemplate string,
	isReturnCode bool,
) (_ *schemauser.EmailCodeAddedEvent, plainCode string, err error) {
	cryptoCode, err := code(ctx)
	if err != nil {
		return nil, "", err
	}
	if isReturnCode {
		plainCode = cryptoCode.Plain
	}
	return schemauser.NewEmailCodeAddedEvent(ctx,
		UserV3AggregateFromWriteModel(&wm.WriteModel),
		cryptoCode.Crypted,
		cryptoCode.Expiry,
		urlTemplate,
		isReturnCode,
	), plainCode, nil
}

func (wm *UserV3WriteModel) NewPhoneCreate(
	ctx context.Context,
	phone *Phone,
	code func(context.Context) (*EncryptedCode, string, error),
) (_ []eventstore.Command, plainCode string, err error) {
	if err := wm.checkPermissionWrite(ctx, wm.ResourceOwner, wm.AggregateID); err != nil {
		return nil, "", err
	}
	if phone == nil || wm.Phone == string(phone.Number) {
		return nil, "", nil
	}
	events := []eventstore.Command{
		schemauser.NewPhoneUpdatedEvent(ctx,
			UserV3AggregateFromWriteModel(&wm.WriteModel),
			phone.Number,
		),
	}
	if phone.Verified {
		events = append(events, wm.newPhoneVerifiedEvent(ctx))
	} else {
		codeEvent, code, err := wm.newPhoneCodeAddedEvent(ctx, code, phone.ReturnCode)
		if err != nil {
			return nil, "", err
		}
		events = append(events, codeEvent)
		if code != "" {
			plainCode = code
		}
	}
	return events, plainCode, nil
}

func (wm *UserV3WriteModel) NewPhoneUpdate(
	ctx context.Context,
	phone *Phone,
	code func(context.Context) (*EncryptedCode, string, error),
) (_ []eventstore.Command, plainCode string, err error) {
	if !wm.PhoneWM {
		return nil, "", nil
	}
	if !wm.Exists() {
		return nil, "", zerrors.ThrowNotFound(nil, "COMMAND-b33QAVgel6", "Errors.User.NotFound")
	}
	return wm.NewPhoneCreate(ctx, phone, code)
}

func (wm *UserV3WriteModel) NewPhoneVerify(
	ctx context.Context,
	verify func(creationDate time.Time, expiry time.Duration, cryptoCode *crypto.CryptoValue) error,
) ([]eventstore.Command, error) {
	if !wm.PhoneWM {
		return nil, nil
	}
	if !wm.Exists() {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-bx2OLtgGNS", "Errors.User.NotFound")
	}
	if err := wm.checkPermissionWrite(ctx, wm.ResourceOwner, wm.AggregateID); err != nil {
		return nil, err
	}
	if wm.PhoneCode == nil {
		return nil, nil
	}
	if err := verify(wm.PhoneCode.CreationDate, wm.PhoneCode.Expiry, wm.PhoneCode.Code); err != nil {
		return nil, err
	}
	return []eventstore.Command{wm.newPhoneVerifiedEvent(ctx)}, nil
}

func (wm *UserV3WriteModel) newPhoneVerifiedEvent(
	ctx context.Context,
) *schemauser.PhoneVerifiedEvent {
	return schemauser.NewPhoneVerifiedEvent(ctx, UserV3AggregateFromWriteModel(&wm.WriteModel))
}

func (wm *UserV3WriteModel) NewResendPhoneCode(
	ctx context.Context,
	code func(context.Context) (*EncryptedCode, string, error),
	isReturnCode bool,
) (_ []eventstore.Command, plainCode string, err error) {
	if !wm.PhoneWM {
		return nil, "", nil
	}
	if !wm.Exists() {
		return nil, "", zerrors.ThrowNotFound(nil, "COMMAND-z8Bu9vuL9s", "Errors.User.NotFound")
	}
	if err := wm.checkPermissionWrite(ctx, wm.ResourceOwner, wm.AggregateID); err != nil {
		return nil, "", err
	}
	if wm.PhoneCode == nil {
		return nil, "", zerrors.ThrowPreconditionFailed(err, "COMMAND-fEsHdqECzb", "Errors.User.Code.Empty")
	}
	event, plainCode, err := wm.newPhoneCodeAddedEvent(ctx, code, isReturnCode)
	if err != nil {
		return nil, "", err
	}
	return []eventstore.Command{event}, plainCode, nil
}

func (wm *UserV3WriteModel) newPhoneCodeAddedEvent(
	ctx context.Context,
	code func(context.Context) (*EncryptedCode, string, error),
	isReturnCode bool,
) (_ *schemauser.PhoneCodeAddedEvent, plainCode string, err error) {
	cryptoCode, generatorID, err := code(ctx)
	if err != nil {
		return nil, "", err
	}
	if isReturnCode {
		plainCode = cryptoCode.Plain
	}
	return schemauser.NewPhoneCodeAddedEvent(ctx,
		UserV3AggregateFromWriteModel(&wm.WriteModel),
		cryptoCode.CryptedCode(),
		cryptoCode.CodeExpiry(),
		isReturnCode,
		generatorID,
	), plainCode, nil
}
