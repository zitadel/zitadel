package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/repository/org"
)

func TestCommandSide_SetCustomOrgLoginText(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		resourceOwner string
		config        *domain.CustomLoginText
	}
	type res struct {
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "no resourceowner, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:    context.Background(),
				config: &domain.CustomLoginText{},
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "invalid custom login text, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				config:        &domain.CustomLoginText{},
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "custom login text set all fields, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeySelectAccountTitle, "Title", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeySelectAccountDescription, "Description", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeySelectAccountTitleLinkingProcess, "TitleLinking", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeySelectAccountDescriptionLinkingProcess, "DescriptionLinking", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeySelectAccountOtherUser, "OtherUser", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeySelectAccountSessionStateActive, "SessionState0", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeySelectAccountSessionStateInactive, "SessionState1", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeySelectAccountUserMustBeMemberOfOrg, "MustBeMemberOfOrg", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyLoginTitle, "Title", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyLoginDescription, "Description", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyLoginTitleLinkingProcess, "TitleLinking", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyLoginDescriptionLinkingProcess, "DescriptionLinking", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyLoginNameLabel, "LoginNameLabel", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyLoginUsernamePlaceHolder, "UsernamePlaceholder", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyLoginLoginnamePlaceHolder, "LoginnamePlaceholder", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyLoginRegisterButtonText, "RegisterButtonText", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyLoginNextButtonText, "NextButtonText", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyLoginExternalUserDescription, "ExternalLogin", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyLoginUserMustBeMemberOfOrg, "MustBeMemberOfOrg", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordTitle, "Title", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordDescription, "Description", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordLabel, "PasswordLabel", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordResetLinkText, "ResetLinkText", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordBackButtonText, "BackButtonText", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordNextButtonText, "NextButtonText", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordMinLength, "MinLength", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordHasUppercase, "HasUppercase", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordHasLowercase, "HasLowercase", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordHasNumber, "HasNumber", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordHasSymbol, "HasSymbol", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordConfirmation, "Confirmation", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyUsernameChangeTitle, "Title", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyUsernameChangeDescription, "Description", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyUsernameChangeUsernameLabel, "UsernameLabel", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyUsernameChangeCancelButtonText, "CancelButtonText", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyUsernameChangeNextButtonText, "NextButtonText", language.English,
								),
							), eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyUsernameChangeDoneTitle, "Title", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyUsernameChangeDoneDescription, "Description", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyUsernameChangeDoneNextButtonText, "NextButtonText", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitPasswordTitle, "Title", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitPasswordDescription, "Description", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitPasswordCodeLabel, "CodeLabel", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitPasswordNewPasswordLabel, "NewPasswordLabel", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitPasswordNewPasswordConfirmLabel, "NewPasswordConfirmLabel", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitPasswordNextButtonText, "NextButtonText", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitPasswordResendButtonText, "ResendButtonText", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitPasswordDoneTitle, "Title", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitPasswordDoneDescription, "Description", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitPasswordDoneNextButtonText, "NextButtonText", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitPasswordDoneCancelButtonText, "CancelButtonText", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyEmailVerificationTitle, "Title", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyEmailVerificationDescription, "Description", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyEmailVerificationCodeLabel, "CodeLabel", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyEmailVerificationNextButtonText, "NextButtonText", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyEmailVerificationResendButtonText, "ResendButtonText", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyEmailVerificationDoneTitle, "Title", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyEmailVerificationDoneDescription, "Description", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyEmailVerificationDoneNextButtonText, "NextButtonText", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyEmailVerificationDoneCancelButtonText, "CancelButtonText", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyEmailVerificationDoneLoginButtonText, "LoginButtonText", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitializeUserTitle, "Title", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitializeUserDescription, "Description", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitializeUserCodeLabel, "CodeLabel", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitializeUserNewPasswordLabel, "NewPasswordLabel", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitializeUserNewPasswordConfirmLabel, "NewPasswordConfirmLabel", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitializeUserResendButtonText, "ResendButtonText", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitializeUserNextButtonText, "NextButtonText", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitUserDoneTitle, "Title", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitUserDoneDescription, "Description", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitUserDoneCancelButtonText, "CancelButtonText", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitUserDoneNextButtonText, "NextButtonText", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitMFAPromptTitle, "Title", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitMFAPromptDescription, "Description", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitMFAPromptOTPOption, "Provider0", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitMFAPromptU2FOption, "Provider1", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitMFAPromptSkipButtonText, "SkipButtonText", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitMFAPromptNextButtonText, "NextButtonText", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitMFAOTPTitle, "Title", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitMFAOTPDescription, "Description", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitMFAOTPDescriptionOTP, "OTPDescription", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitMFAOTPSecretLabel, "SecretLabel", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitMFAOTPCodeLabel, "CodeLabel", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitMFAOTPNextButtonText, "NextButtonText", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitMFAOTPCancelButtonText, "CancelButtonText", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitMFAU2FTitle, "Title", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitMFAU2FDescription, "Description", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitMFAU2FTokenNameLabel, "TokenNameLabel", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitMFAU2FRegisterTokenButtonText, "RegisterTokenButtonText", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitMFAU2FNotSupported, "NotSupported", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitMFAU2FErrorRetry, "ErrorRetry", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitMFADoneTitle, "Title", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitMFADoneDescription, "Description", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitMFADoneCancelButtonText, "CancelButtonText", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitMFADoneNextButtonText, "NextButtonText", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyMFAProvidersChooseOther, "ChooseOther", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyMFAProvidersOTP, "Provider0", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyMFAProvidersU2F, "Provider1", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyVerifyMFAOTPTitle, "Title", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyVerifyMFAOTPDescription, "Description", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyVerifyMFAOTPCodeLabel, "CodeLabel", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyVerifyMFAOTPNextButtonText, "NextButtonText", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyVerifyMFAU2FTitle, "Title", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyVerifyMFAU2FDescription, "Description", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyVerifyMFAU2FValidateTokenText, "ValidateTokenButtonText", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyVerifyMFAU2FNotSupported, "NotSupported", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyVerifyMFAU2FErrorRetry, "ErrorRetry", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordlessTitle, "Title", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordlessDescription, "Description", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordlessLoginWithPwButtonText, "LoginWithPwButtonText", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordlessValidateTokenButtonText, "ValidateTokenButtonText", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordlessNotSupported, "NotSupported", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordlessErrorRetry, "ErrorRetry", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordChangeTitle, "Title", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordChangeDescription, "Description", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordChangeOldPasswordLabel, "OldPasswordLabel", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordChangeNewPasswordLabel, "NewPasswordLabel", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordChangeNewPasswordConfirmLabel, "NewPasswordConfirmLabel", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordChangeCancelButtonText, "CancelButtonText", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordChangeNextButtonText, "NextButtonText", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordChangeDoneTitle, "Title", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordChangeDoneDescription, "Description", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordChangeDoneNextButtonText, "NextButtonText", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordResetDoneTitle, "Title", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordResetDoneDescription, "Description", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordResetDoneNextButtonText, "NextButtonText", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationOptionTitle, "Title", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationOptionDescription, "Description", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationOptionUserNameButtonText, "RegisterUsernamePasswordButtonText", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationUserTitle, "Title", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationUserDescription, "Description", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationUserFirstnameLabel, "FirstnameLabel", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationUserLastnameLabel, "LastnameLabel", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationUserEmailLabel, "EmailLabel", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationUserUsernameLabel, "UsernameLabel", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationUserLanguageLabel, "LanguageLabel", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationUserGenderLabel, "GenderLabel", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationUserPasswordLabel, "PasswordLabel", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationUserPasswordConfirmLabel, "PasswordConfirmLabel", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationUserTOSAndPrivacyLabel, "TOSAndPrivacyLabel", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationUserTOSConfirm, "TOSConfirm", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationUserTOSLink, "TOSLink", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationUserTOSLinkText, "TOSLinkText", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationUserPrivacyConfirm, "PrivacyConfirm", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationUserPrivacyLink, "PrivacyLink", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationUserPrivacyLinkText, "PrivacyLinkText", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationUserExternalLoginDescription, "ExternalLoginDescription", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationUserNextButtonText, "NextButtonText", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationUserBackButtonText, "BackButtonText", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegisterOrgTitle, "Title", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegisterOrgDescription, "Description", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegisterOrgOrgNameLabel, "OrgNameLabel", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegisterOrgFirstnameLabel, "FirstnameLabel", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegisterOrgLastnameLabel, "LastnameLabel", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegisterOrgUsernameLabel, "UsernameLabel", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegisterOrgEmailLabel, "EmailLabel", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegisterOrgPasswordLabel, "PasswordLabel", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegisterOrgPasswordConfirmLabel, "PasswordConfirmLabel", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegisterOrgTOSAndPrivacyLabel, "TOSAndPrivacyLabel", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegisterOrgTOSConfirm, "TOSConfirm", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegisterOrgTOSLink, "TOSLink", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegisterOrgTOSLinkText, "TOSLinkText", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegisterOrgPrivacyConfirm, "PrivacyConfirm", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegisterOrgPrivacyLink, "PrivacyLink", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegisterOrgPrivacyLinkText, "PrivacyLinkText", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegisterOrgExternalLoginDescription, "ExternalLoginDescription", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegisterOrgSaveButtonText, "SaveButtonText", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyLinkingUserDoneTitle, "Title", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyLinkingUserDoneDescription, "Description", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyLinkingUserDoneCancelButtonText, "CancelButtonText", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyLinkingUserDoneNextButtonText, "NextButtonText", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyExternalNotFoundTitle, "Title", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyExternalNotFoundDescription, "Description", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyExternalNotFoundLinkButtonText, "LinkButtonText", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyExternalNotFoundAutoRegisterButtonText, "AutoRegisterButtonText", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeySuccessLoginTitle, "Title", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeySuccessLoginAutoRedirectDescription, "AutoRedirectDescription", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeySuccessLoginRedirectedDescription, "RedirectedDescription", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeySuccessLoginNextButtonText, "NextButtonText", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyLogoutDoneTitle, "Title", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyLogoutDoneDescription, "Description", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyLogoutDoneLoginButtonText, "LoginButtonText", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyFooterTOS, "TOS", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyFooterTOSLink, "TOSLink", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyFooterPrivacy, "PrivacyPolicy", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyFooterPrivacyLink, "PrivacyPolicyLink", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyFooterHelp, "Help", language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyFooterHelpLink, "HelpLink", language.English,
								),
							),
						},
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				config: &domain.CustomLoginText{
					Language: language.English,
					SelectAccount: domain.SelectAccountScreenText{
						Title:              "Title",
						Description:        "Description",
						TitleLinking:       "TitleLinking",
						DescriptionLinking: "DescriptionLinking",
						OtherUser:          "OtherUser",
						SessionState0:      "SessionState0",
						SessionState1:      "SessionState1",
						MustBeMemberOfOrg:  "MustBeMemberOfOrg",
					},
					Login: domain.LoginScreenText{
						Title:                "Title",
						Description:          "Description",
						TitleLinking:         "TitleLinking",
						DescriptionLinking:   "DescriptionLinking",
						LoginNameLabel:       "LoginNameLabel",
						UsernamePlaceholder:  "UsernamePlaceholder",
						LoginnamePlaceholder: "LoginnamePlaceholder",
						RegisterButtonText:   "RegisterButtonText",
						NextButtonText:       "NextButtonText",
						ExternalLogin:        "ExternalLogin",
						MustBeMemberOfOrg:    "MustBeMemberOfOrg",
					},
					Password: domain.PasswordScreenText{
						Title:          "Title",
						Description:    "Description",
						PasswordLabel:  "PasswordLabel",
						ResetLinkText:  "ResetLinkText",
						BackButtonText: "BackButtonText",
						NextButtonText: "NextButtonText",
						MinLength:      "MinLength",
						HasUppercase:   "HasUppercase",
						HasLowercase:   "HasLowercase",
						HasNumber:      "HasNumber",
						HasSymbol:      "HasSymbol",
						Confirmation:   "Confirmation",
					},
					UsernameChange: domain.UsernameChangeScreenText{
						Title:            "Title",
						Description:      "Description",
						UsernameLabel:    "UsernameLabel",
						CancelButtonText: "CancelButtonText",
						NextButtonText:   "NextButtonText",
					},
					UsernameChangeDone: domain.UsernameChangeDoneScreenText{
						Title:          "Title",
						Description:    "Description",
						NextButtonText: "NextButtonText",
					},
					InitPassword: domain.InitPasswordScreenText{
						Title:                   "Title",
						Description:             "Description",
						CodeLabel:               "CodeLabel",
						NewPasswordLabel:        "NewPasswordLabel",
						NewPasswordConfirmLabel: "NewPasswordConfirmLabel",
						NextButtonText:          "NextButtonText",
						ResendButtonText:        "ResendButtonText",
					},
					InitPasswordDone: domain.InitPasswordDoneScreenText{
						Title:            "Title",
						Description:      "Description",
						NextButtonText:   "NextButtonText",
						CancelButtonText: "CancelButtonText",
					},
					EmailVerification: domain.EmailVerificationScreenText{
						Title:            "Title",
						Description:      "Description",
						CodeLabel:        "CodeLabel",
						NextButtonText:   "NextButtonText",
						ResendButtonText: "ResendButtonText",
					},
					EmailVerificationDone: domain.EmailVerificationDoneScreenText{
						Title:            "Title",
						Description:      "Description",
						NextButtonText:   "NextButtonText",
						CancelButtonText: "CancelButtonText",
						LoginButtonText:  "LoginButtonText",
					},
					InitUser: domain.InitializeUserScreenText{
						Title:                   "Title",
						Description:             "Description",
						CodeLabel:               "CodeLabel",
						NewPasswordLabel:        "NewPasswordLabel",
						NewPasswordConfirmLabel: "NewPasswordConfirmLabel",
						ResendButtonText:        "ResendButtonText",
						NextButtonText:          "NextButtonText",
					},
					InitUserDone: domain.InitializeUserDoneScreenText{
						Title:            "Title",
						Description:      "Description",
						CancelButtonText: "CancelButtonText",
						NextButtonText:   "NextButtonText",
					},
					InitMFAPrompt: domain.InitMFAPromptScreenText{
						Title:          "Title",
						Description:    "Description",
						Provider0:      "Provider0",
						Provider1:      "Provider1",
						SkipButtonText: "SkipButtonText",
						NextButtonText: "NextButtonText",
					},
					InitMFAOTP: domain.InitMFAOTPScreenText{
						Title:            "Title",
						Description:      "Description",
						OTPDescription:   "OTPDescription",
						SecretLabel:      "SecretLabel",
						CodeLabel:        "CodeLabel",
						NextButtonText:   "NextButtonText",
						CancelButtonText: "CancelButtonText",
					},
					InitMFAU2F: domain.InitMFAU2FScreenText{
						Title:                   "Title",
						Description:             "Description",
						TokenNameLabel:          "TokenNameLabel",
						RegisterTokenButtonText: "RegisterTokenButtonText",
						NotSupported:            "NotSupported",
						ErrorRetry:              "ErrorRetry",
					},
					InitMFADone: domain.InitMFADoneScreenText{
						Title:            "Title",
						Description:      "Description",
						CancelButtonText: "CancelButtonText",
						NextButtonText:   "NextButtonText",
					},
					MFAProvider: domain.MFAProvidersText{
						ChooseOther: "ChooseOther",
						Provider0:   "Provider0",
						Provider1:   "Provider1",
					},
					VerifyMFAOTP: domain.VerifyMFAOTPScreenText{
						Title:          "Title",
						Description:    "Description",
						CodeLabel:      "CodeLabel",
						NextButtonText: "NextButtonText",
					},
					VerifyMFAU2F: domain.VerifyMFAU2FScreenText{
						Title:                   "Title",
						Description:             "Description",
						ValidateTokenButtonText: "ValidateTokenButtonText",
						NotSupported:            "NotSupported",
						ErrorRetry:              "ErrorRetry",
					},
					Passwordless: domain.PasswordlessScreenText{
						Title:                   "Title",
						Description:             "Description",
						LoginWithPwButtonText:   "LoginWithPwButtonText",
						ValidateTokenButtonText: "ValidateTokenButtonText",
						NotSupported:            "NotSupported",
						ErrorRetry:              "ErrorRetry",
					},
					PasswordChange: domain.PasswordChangeScreenText{
						Title:                   "Title",
						Description:             "Description",
						OldPasswordLabel:        "OldPasswordLabel",
						NewPasswordLabel:        "NewPasswordLabel",
						NewPasswordConfirmLabel: "NewPasswordConfirmLabel",
						CancelButtonText:        "CancelButtonText",
						NextButtonText:          "NextButtonText",
					},
					PasswordChangeDone: domain.PasswordChangeDoneScreenText{
						Title:          "Title",
						Description:    "Description",
						NextButtonText: "NextButtonText",
					},
					PasswordResetDone: domain.PasswordResetDoneScreenText{
						Title:          "Title",
						Description:    "Description",
						NextButtonText: "NextButtonText",
					},
					RegisterOption: domain.RegistrationOptionScreenText{
						Title:                              "Title",
						Description:                        "Description",
						RegisterUsernamePasswordButtonText: "RegisterUsernamePasswordButtonText",
						ExternalLoginDescription:           "",
					},
					RegistrationUser: domain.RegistrationUserScreenText{
						Title:                    "Title",
						Description:              "Description",
						FirstnameLabel:           "FirstnameLabel",
						LastnameLabel:            "LastnameLabel",
						EmailLabel:               "EmailLabel",
						UsernameLabel:            "UsernameLabel",
						LanguageLabel:            "LanguageLabel",
						GenderLabel:              "GenderLabel",
						PasswordLabel:            "PasswordLabel",
						PasswordConfirmLabel:     "PasswordConfirmLabel",
						TOSAndPrivacyLabel:       "TOSAndPrivacyLabel",
						TOSConfirm:               "TOSConfirm",
						TOSLink:                  "TOSLink",
						TOSLinkText:              "TOSLinkText",
						PrivacyConfirm:           "PrivacyConfirm",
						PrivacyLink:              "PrivacyLink",
						PrivacyLinkText:          "PrivacyLinkText",
						ExternalLoginDescription: "ExternalLoginDescription",
						NextButtonText:           "NextButtonText",
						BackButtonText:           "BackButtonText",
					},
					RegistrationOrg: domain.RegistrationOrgScreenText{
						Title:                    "Title",
						Description:              "Description",
						OrgNameLabel:             "OrgNameLabel",
						FirstnameLabel:           "FirstnameLabel",
						LastnameLabel:            "LastnameLabel",
						UsernameLabel:            "UsernameLabel",
						EmailLabel:               "EmailLabel",
						PasswordLabel:            "PasswordLabel",
						PasswordConfirmLabel:     "PasswordConfirmLabel",
						TOSAndPrivacyLabel:       "TOSAndPrivacyLabel",
						TOSConfirm:               "TOSConfirm",
						TOSLink:                  "TOSLink",
						TOSLinkText:              "TOSLinkText",
						PrivacyConfirm:           "PrivacyConfirm",
						PrivacyLink:              "PrivacyLink",
						PrivacyLinkText:          "PrivacyLinkText",
						ExternalLoginDescription: "ExternalLoginDescription",
						SaveButtonText:           "SaveButtonText",
					},
					LinkingUsersDone: domain.LinkingUserDoneScreenText{
						Title:            "Title",
						Description:      "Description",
						CancelButtonText: "CancelButtonText",
						NextButtonText:   "NextButtonText",
					},
					ExternalNotFoundOption: domain.ExternalUserNotFoundScreenText{
						Title:                  "Title",
						Description:            "Description",
						LinkButtonText:         "LinkButtonText",
						AutoRegisterButtonText: "AutoRegisterButtonText",
					},
					LoginSuccess: domain.SuccessLoginScreenText{
						Title:                   "Title",
						AutoRedirectDescription: "AutoRedirectDescription",
						RedirectedDescription:   "RedirectedDescription",
						NextButtonText:          "NextButtonText",
					},
					LogoutDone: domain.LogoutDoneScreenText{
						Title:           "Title",
						Description:     "Description",
						LoginButtonText: "LoginButtonText",
					},
					Footer: domain.FooterText{
						TOS:               "TOS",
						TOSLink:           "TOSLink",
						PrivacyPolicy:     "PrivacyPolicy",
						PrivacyPolicyLink: "PrivacyPolicyLink",
						Help:              "Help",
						HelpLink:          "HelpLink",
					},
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.SetOrgLoginText(tt.args.ctx, tt.args.resourceOwner, tt.args.config)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}
