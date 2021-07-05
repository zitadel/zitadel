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
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyLoginExternalUserDescription, "ExternalUserDescription", language.English,
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
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationOptionExternalLoginDescription, "ExternalLoginDescription", language.English,
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
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationUserDescriptionOrgRegister, "DescriptionOrgRegister", language.English,
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
						Title:                   "Title",
						Description:             "Description",
						TitleLinking:            "TitleLinking",
						DescriptionLinking:      "DescriptionLinking",
						LoginNameLabel:          "LoginNameLabel",
						UsernamePlaceholder:     "UsernamePlaceholder",
						LoginnamePlaceholder:    "LoginnamePlaceholder",
						RegisterButtonText:      "RegisterButtonText",
						NextButtonText:          "NextButtonText",
						ExternalUserDescription: "ExternalUserDescription",
						MustBeMemberOfOrg:       "MustBeMemberOfOrg",
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
						ExternalLoginDescription:           "ExternalLoginDescription",
					},
					RegistrationUser: domain.RegistrationUserScreenText{
						Title:                    "Title",
						Description:              "Description",
						DescriptionOrgRegister:   "DescriptionOrgRegister",
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
		{
			name: "custom login text remove all fields, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
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
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyLoginExternalUserDescription, "ExternalUserDescription", language.English,
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
						),
						eventFromEventPusher(
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
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationOptionExternalLoginDescription, "ExternalLoginDescription", language.English,
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
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationUserDescriptionOrgRegister, "DescriptionOrgRegister", language.English,
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
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeySelectAccountTitle, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeySelectAccountDescription, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeySelectAccountTitleLinkingProcess, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeySelectAccountDescriptionLinkingProcess, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeySelectAccountOtherUser, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeySelectAccountSessionStateActive, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeySelectAccountSessionStateInactive, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeySelectAccountUserMustBeMemberOfOrg, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyLoginTitle, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyLoginDescription, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyLoginTitleLinkingProcess, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyLoginDescriptionLinkingProcess, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyLoginNameLabel, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyLoginUsernamePlaceHolder, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyLoginLoginnamePlaceHolder, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyLoginRegisterButtonText, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyLoginNextButtonText, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyLoginExternalUserDescription, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyLoginUserMustBeMemberOfOrg, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordTitle, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordDescription, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordLabel, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordResetLinkText, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordBackButtonText, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordNextButtonText, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordMinLength, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordHasUppercase, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordHasLowercase, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordHasNumber, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordHasSymbol, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordConfirmation, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyUsernameChangeTitle, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyUsernameChangeDescription, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyUsernameChangeUsernameLabel, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyUsernameChangeCancelButtonText, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyUsernameChangeNextButtonText, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyUsernameChangeDoneTitle, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyUsernameChangeDoneDescription, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyUsernameChangeDoneNextButtonText, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitPasswordTitle, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitPasswordDescription, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitPasswordCodeLabel, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitPasswordNewPasswordLabel, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitPasswordNewPasswordConfirmLabel, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitPasswordNextButtonText, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitPasswordResendButtonText, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitPasswordDoneTitle, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitPasswordDoneDescription, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitPasswordDoneNextButtonText, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitPasswordDoneCancelButtonText, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyEmailVerificationTitle, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyEmailVerificationDescription, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyEmailVerificationCodeLabel, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyEmailVerificationNextButtonText, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyEmailVerificationResendButtonText, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyEmailVerificationDoneTitle, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyEmailVerificationDoneDescription, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyEmailVerificationDoneNextButtonText, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyEmailVerificationDoneCancelButtonText, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyEmailVerificationDoneLoginButtonText, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitializeUserTitle, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitializeUserDescription, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitializeUserCodeLabel, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitializeUserNewPasswordLabel, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitializeUserNewPasswordConfirmLabel, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitializeUserResendButtonText, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitializeUserNextButtonText, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitUserDoneTitle, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitUserDoneDescription, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitUserDoneCancelButtonText, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitUserDoneNextButtonText, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitMFAPromptTitle, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitMFAPromptDescription, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitMFAPromptOTPOption, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitMFAPromptU2FOption, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitMFAPromptSkipButtonText, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitMFAPromptNextButtonText, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitMFAOTPTitle, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitMFAOTPDescription, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitMFAOTPDescriptionOTP, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitMFAOTPSecretLabel, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitMFAOTPCodeLabel, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitMFAOTPNextButtonText, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitMFAOTPCancelButtonText, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitMFAU2FTitle, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitMFAU2FDescription, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitMFAU2FTokenNameLabel, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitMFAU2FRegisterTokenButtonText, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitMFAU2FNotSupported, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitMFAU2FErrorRetry, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitMFADoneTitle, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitMFADoneDescription, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitMFADoneCancelButtonText, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitMFADoneNextButtonText, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyMFAProvidersChooseOther, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyMFAProvidersOTP, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyMFAProvidersU2F, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyVerifyMFAOTPTitle, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyVerifyMFAOTPDescription, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyVerifyMFAOTPCodeLabel, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyVerifyMFAOTPNextButtonText, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyVerifyMFAU2FTitle, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyVerifyMFAU2FDescription, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyVerifyMFAU2FValidateTokenText, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyVerifyMFAU2FNotSupported, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyVerifyMFAU2FErrorRetry, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordlessTitle, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordlessDescription, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordlessLoginWithPwButtonText, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordlessValidateTokenButtonText, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordlessNotSupported, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordlessErrorRetry, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordChangeTitle, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordChangeDescription, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordChangeOldPasswordLabel, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordChangeNewPasswordLabel, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordChangeNewPasswordConfirmLabel, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordChangeCancelButtonText, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordChangeNextButtonText, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordChangeDoneTitle, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordChangeDoneDescription, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordChangeDoneNextButtonText, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordResetDoneTitle, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordResetDoneDescription, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordResetDoneNextButtonText, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationOptionTitle, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationOptionDescription, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationOptionUserNameButtonText, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationOptionExternalLoginDescription, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationUserTitle, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationUserDescription, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationUserDescriptionOrgRegister, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationUserFirstnameLabel, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationUserLastnameLabel, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationUserEmailLabel, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationUserUsernameLabel, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationUserLanguageLabel, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationUserGenderLabel, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationUserPasswordLabel, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationUserPasswordConfirmLabel, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationUserTOSAndPrivacyLabel, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationUserTOSConfirm, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationUserTOSLink, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationUserTOSLinkText, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationUserPrivacyConfirm, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationUserPrivacyLink, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationUserPrivacyLinkText, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationUserExternalLoginDescription, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationUserNextButtonText, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationUserBackButtonText, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegisterOrgTitle, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegisterOrgDescription, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegisterOrgOrgNameLabel, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegisterOrgFirstnameLabel, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegisterOrgLastnameLabel, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegisterOrgUsernameLabel, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegisterOrgEmailLabel, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegisterOrgPasswordLabel, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegisterOrgPasswordConfirmLabel, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegisterOrgTOSAndPrivacyLabel, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegisterOrgTOSConfirm, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegisterOrgTOSLink, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegisterOrgTOSLinkText, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegisterOrgPrivacyConfirm, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegisterOrgPrivacyLink, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegisterOrgPrivacyLinkText, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegisterOrgExternalLoginDescription, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegisterOrgSaveButtonText, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyLinkingUserDoneTitle, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyLinkingUserDoneDescription, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyLinkingUserDoneCancelButtonText, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyLinkingUserDoneNextButtonText, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyExternalNotFoundTitle, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyExternalNotFoundDescription, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyExternalNotFoundLinkButtonText, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyExternalNotFoundAutoRegisterButtonText, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeySuccessLoginTitle, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeySuccessLoginAutoRedirectDescription, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeySuccessLoginRedirectedDescription, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeySuccessLoginNextButtonText, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyLogoutDoneTitle, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyLogoutDoneDescription, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyLogoutDoneLoginButtonText, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyFooterTOS, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyFooterTOSLink, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyFooterPrivacy, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyFooterPrivacyLink, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyFooterHelp, language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyFooterHelpLink, language.English,
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
					Language:               language.English,
					SelectAccount:          domain.SelectAccountScreenText{},
					Login:                  domain.LoginScreenText{},
					Password:               domain.PasswordScreenText{},
					UsernameChange:         domain.UsernameChangeScreenText{},
					UsernameChangeDone:     domain.UsernameChangeDoneScreenText{},
					InitPassword:           domain.InitPasswordScreenText{},
					InitPasswordDone:       domain.InitPasswordDoneScreenText{},
					EmailVerification:      domain.EmailVerificationScreenText{},
					EmailVerificationDone:  domain.EmailVerificationDoneScreenText{},
					InitUser:               domain.InitializeUserScreenText{},
					InitUserDone:           domain.InitializeUserDoneScreenText{},
					InitMFAPrompt:          domain.InitMFAPromptScreenText{},
					InitMFAOTP:             domain.InitMFAOTPScreenText{},
					InitMFAU2F:             domain.InitMFAU2FScreenText{},
					InitMFADone:            domain.InitMFADoneScreenText{},
					MFAProvider:            domain.MFAProvidersText{},
					VerifyMFAOTP:           domain.VerifyMFAOTPScreenText{},
					VerifyMFAU2F:           domain.VerifyMFAU2FScreenText{},
					Passwordless:           domain.PasswordlessScreenText{},
					PasswordChange:         domain.PasswordChangeScreenText{},
					PasswordChangeDone:     domain.PasswordChangeDoneScreenText{},
					PasswordResetDone:      domain.PasswordResetDoneScreenText{},
					RegisterOption:         domain.RegistrationOptionScreenText{},
					RegistrationUser:       domain.RegistrationUserScreenText{},
					RegistrationOrg:        domain.RegistrationOrgScreenText{},
					LinkingUsersDone:       domain.LinkingUserDoneScreenText{},
					ExternalNotFoundOption: domain.ExternalUserNotFoundScreenText{},
					LoginSuccess:           domain.SuccessLoginScreenText{},
					LogoutDone:             domain.LogoutDoneScreenText{},
					Footer:                 domain.FooterText{},
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "custom login text set all fields, all texts removed, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
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
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyLoginExternalUserDescription, "ExternalUserDescription", language.English,
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
						),
						eventFromEventPusher(
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
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationOptionExternalLoginDescription, "ExternalLoginDescription", language.English,
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
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationUserDescriptionOrgRegister, "DescriptionOrgRegister", language.English,
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
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeySelectAccountTitle, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeySelectAccountDescription, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeySelectAccountTitleLinkingProcess, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeySelectAccountDescriptionLinkingProcess, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeySelectAccountOtherUser, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeySelectAccountSessionStateActive, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeySelectAccountSessionStateInactive, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeySelectAccountUserMustBeMemberOfOrg, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyLoginTitle, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyLoginDescription, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyLoginTitleLinkingProcess, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyLoginDescriptionLinkingProcess, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyLoginNameLabel, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyLoginUsernamePlaceHolder, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyLoginLoginnamePlaceHolder, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyLoginRegisterButtonText, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyLoginNextButtonText, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyLoginExternalUserDescription, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyLoginUserMustBeMemberOfOrg, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordTitle, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordDescription, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordLabel, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordResetLinkText, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordBackButtonText, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordNextButtonText, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordMinLength, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordHasUppercase, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordHasLowercase, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordHasNumber, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordHasSymbol, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordConfirmation, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyUsernameChangeTitle, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyUsernameChangeDescription, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyUsernameChangeUsernameLabel, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyUsernameChangeCancelButtonText, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyUsernameChangeNextButtonText, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyUsernameChangeDoneTitle, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyUsernameChangeDoneDescription, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyUsernameChangeDoneNextButtonText, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitPasswordTitle, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitPasswordDescription, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitPasswordCodeLabel, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitPasswordNewPasswordLabel, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitPasswordNewPasswordConfirmLabel, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitPasswordNextButtonText, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitPasswordResendButtonText, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitPasswordDoneTitle, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitPasswordDoneDescription, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitPasswordDoneNextButtonText, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitPasswordDoneCancelButtonText, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyEmailVerificationTitle, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyEmailVerificationDescription, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyEmailVerificationCodeLabel, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyEmailVerificationNextButtonText, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyEmailVerificationResendButtonText, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyEmailVerificationDoneTitle, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyEmailVerificationDoneDescription, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyEmailVerificationDoneNextButtonText, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyEmailVerificationDoneCancelButtonText, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyEmailVerificationDoneLoginButtonText, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitializeUserTitle, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitializeUserDescription, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitializeUserCodeLabel, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitializeUserNewPasswordLabel, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitializeUserNewPasswordConfirmLabel, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitializeUserResendButtonText, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitializeUserNextButtonText, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitUserDoneTitle, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitUserDoneDescription, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitUserDoneCancelButtonText, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitUserDoneNextButtonText, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitMFAPromptTitle, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitMFAPromptDescription, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitMFAPromptOTPOption, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitMFAPromptU2FOption, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitMFAPromptSkipButtonText, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitMFAPromptNextButtonText, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitMFAOTPTitle, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitMFAOTPDescription, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitMFAOTPDescriptionOTP, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitMFAOTPSecretLabel, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitMFAOTPCodeLabel, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitMFAOTPNextButtonText, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitMFAOTPCancelButtonText, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitMFAU2FTitle, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitMFAU2FDescription, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitMFAU2FTokenNameLabel, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitMFAU2FRegisterTokenButtonText, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitMFAU2FNotSupported, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitMFAU2FErrorRetry, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitMFADoneTitle, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitMFADoneDescription, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitMFADoneCancelButtonText, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyInitMFADoneNextButtonText, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyMFAProvidersChooseOther, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyMFAProvidersOTP, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyMFAProvidersU2F, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyVerifyMFAOTPTitle, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyVerifyMFAOTPDescription, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyVerifyMFAOTPCodeLabel, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyVerifyMFAOTPNextButtonText, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyVerifyMFAU2FTitle, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyVerifyMFAU2FDescription, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyVerifyMFAU2FValidateTokenText, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyVerifyMFAU2FNotSupported, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyVerifyMFAU2FErrorRetry, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordlessTitle, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordlessDescription, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordlessLoginWithPwButtonText, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordlessValidateTokenButtonText, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordlessNotSupported, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordlessErrorRetry, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordChangeTitle, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordChangeDescription, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordChangeOldPasswordLabel, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordChangeNewPasswordLabel, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordChangeNewPasswordConfirmLabel, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordChangeCancelButtonText, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordChangeNextButtonText, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordChangeDoneTitle, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordChangeDoneDescription, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordChangeDoneNextButtonText, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordResetDoneTitle, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordResetDoneDescription, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyPasswordResetDoneNextButtonText, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationOptionTitle, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationOptionDescription, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationOptionUserNameButtonText, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationOptionExternalLoginDescription, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationUserTitle, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationUserDescription, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationUserDescriptionOrgRegister, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationUserFirstnameLabel, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationUserLastnameLabel, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationUserEmailLabel, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationUserUsernameLabel, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationUserLanguageLabel, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationUserGenderLabel, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationUserPasswordLabel, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationUserPasswordConfirmLabel, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationUserTOSAndPrivacyLabel, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationUserTOSConfirm, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationUserTOSLink, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationUserTOSLinkText, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationUserPrivacyConfirm, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationUserPrivacyLink, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationUserPrivacyLinkText, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationUserExternalLoginDescription, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationUserNextButtonText, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationUserBackButtonText, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegisterOrgTitle, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegisterOrgDescription, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegisterOrgOrgNameLabel, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegisterOrgFirstnameLabel, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegisterOrgLastnameLabel, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegisterOrgUsernameLabel, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegisterOrgEmailLabel, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegisterOrgPasswordLabel, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegisterOrgPasswordConfirmLabel, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegisterOrgTOSAndPrivacyLabel, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegisterOrgTOSConfirm, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegisterOrgTOSLink, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegisterOrgTOSLinkText, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegisterOrgPrivacyConfirm, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegisterOrgPrivacyLink, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegisterOrgPrivacyLinkText, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegisterOrgExternalLoginDescription, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegisterOrgSaveButtonText, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyLinkingUserDoneTitle, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyLinkingUserDoneDescription, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyLinkingUserDoneCancelButtonText, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyLinkingUserDoneNextButtonText, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyExternalNotFoundTitle, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyExternalNotFoundDescription, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyExternalNotFoundLinkButtonText, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyExternalNotFoundAutoRegisterButtonText, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeySuccessLoginTitle, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeySuccessLoginAutoRedirectDescription, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeySuccessLoginRedirectedDescription, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeySuccessLoginNextButtonText, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyLogoutDoneTitle, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyLogoutDoneDescription, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyLogoutDoneLoginButtonText, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyFooterTOS, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyFooterTOSLink, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyFooterPrivacy, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyFooterPrivacyLink, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyFooterHelp, language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyFooterHelpLink, language.English,
							),
						),
					),
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
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyLoginExternalUserDescription, "ExternalUserDescription", language.English,
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
							),
							eventFromEventPusher(
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
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationOptionExternalLoginDescription, "ExternalLoginDescription", language.English,
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
									&org.NewAggregate("org1", "org1").Aggregate, domain.LoginCustomText, domain.LoginKeyRegistrationUserDescriptionOrgRegister, "DescriptionOrgRegister", language.English,
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
						Title:                   "Title",
						Description:             "Description",
						TitleLinking:            "TitleLinking",
						DescriptionLinking:      "DescriptionLinking",
						LoginNameLabel:          "LoginNameLabel",
						UsernamePlaceholder:     "UsernamePlaceholder",
						LoginnamePlaceholder:    "LoginnamePlaceholder",
						RegisterButtonText:      "RegisterButtonText",
						NextButtonText:          "NextButtonText",
						ExternalUserDescription: "ExternalUserDescription",
						MustBeMemberOfOrg:       "MustBeMemberOfOrg",
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
						ExternalLoginDescription:           "ExternalLoginDescription",
					},
					RegistrationUser: domain.RegistrationUserScreenText{
						Title:                    "Title",
						Description:              "Description",
						DescriptionOrgRegister:   "DescriptionOrgRegister",
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
