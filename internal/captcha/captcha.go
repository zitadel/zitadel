package captcha

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/zitadel/zitadel/internal/domain"
)

type Captcha interface {
	Initialize(config any) error
	GetToken(req *http.Request) (string, error)
	ValidateToken(token string) (bool, error)
}

var captchaFactory = map[domain.CaptchaType]func() Captcha{
	domain.CaptchaTypeReCaptcha: func() Captcha {
		return &RecaptchaV2{}
	},
}

func VerifyCaptcha(r *http.Request, authReq *domain.AuthRequest) error {
	var captchaClient Captcha

	factory, ok := captchaFactory[authReq.LoginPolicy.CaptchaType]
	if !ok {
		return errors.New("unsupported captcha type")
	}
	captchaClient = factory()

	err := captchaClient.Initialize(RecaptchaV2Config{
		SiteKey:   authReq.LoginPolicy.CaptchaSiteKey,
		SecretKey: authReq.LoginPolicy.CaptchaSecretKey,
	})
	if err != nil {
		return fmt.Errorf("failed to initialize recaptcha: %w", err)
	}

	token, err := captchaClient.GetToken(r)
	if err != nil {
		return fmt.Errorf("failed to get captcha token: %w", err)
	}

	ok, err = captchaClient.ValidateToken(token)
	if err != nil {
		return fmt.Errorf("captcha validation failed: %w", err)
	}
	if !ok {
		return errors.New("invalid captcha")
	}

	return nil
}
