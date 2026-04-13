package v2

import (
	"strings"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/pkg/grpc/session/v2"
)

type sessionCommand interface {
	domain.CheckUserParent
	domain.CheckPasswordParent
}

func checksToCommands[P sessionCommand](parent P, checks *session.Checks) []domain.Commander {
	cmds := make([]domain.Commander, 0, 8)
	if checks.GetUser() != nil {
		cmds = append(cmds, userCheckToCommand(parent, checks.GetUser()))
	}
	if checks.GetPassword() != nil {
		cmds = append(cmds, passwordCheckToCommand(parent, checks.GetPassword()))
	}
	// if checks.GetWebAuthN() != nil {
	// 	cmds = append(cmds, webAuthNCheckToCommand(parent, checks.GetWebAuthN()))
	// }
	// if checks.GetIdpIntent() != nil {
	// 	cmds = append(cmds, idpIntentCheckToCommand(parent, checks.GetIdpIntent()))
	// }
	// if checks.GetTotp() != nil {
	// 	cmds = append(cmds, totpCheckToCommand(parent, checks.GetTotp()))
	// }
	// if checks.GetOtpSms() != nil {
	// 	cmds = append(cmds, otpSMSCheckToCommand(parent, checks.GetOtpSms()))
	// }
	// if checks.GetOtpEmail() != nil {
	// 	cmds = append(cmds, otpEmailCheckToCommand(parent, checks.GetOtpEmail()))
	// }
	// if checks.GetRecoveryCode() != nil {
	// 	cmds = append(cmds, recoveryCodeCheckToCommand(parent, checks.GetRecoveryCode()))
	// }
	return cmds
}

func userCheckToCommand[P sessionCommand](parent P, check *session.CheckUser) domain.Commander {
	var userID, loginName *string
	switch t := check.GetSearch().(type) {
	case *session.CheckUser_UserId:
		if trimmed := strings.TrimSpace(t.UserId); trimmed != "" {
			userID = &trimmed
		}
	case *session.CheckUser_LoginName:
		if trimmed := strings.TrimSpace(t.LoginName); trimmed != "" {
			loginName = &trimmed
		}
	}
	return domain.NewCheckUserCommand(parent, userID, loginName)
}

func passwordCheckToCommand[P sessionCommand](parent P, check *session.CheckPassword) domain.Commander {
	return domain.NewCheckPasswordCommand(parent, check.GetPassword())
}
