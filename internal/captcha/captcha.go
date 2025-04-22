package captcha

import (
	"net/http"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/zerrors"
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
		return zerrors.ThrowError(nil, "CAPTCHA-aCh1D", "Errors.Captcha.UnexpectedType")
	}
	captchaClient = factory()

	switch authReq.LoginPolicy.CaptchaType {
	case domain.CaptchaTypeReCaptcha:
		err := captchaClient.Initialize(RecaptchaV2Config{
			SiteKey:   authReq.LoginPolicy.CaptchaSiteKey,
			SecretKey: authReq.LoginPolicy.CaptchaSecretKey,
		})
		if err != nil {
			return zerrors.ThrowInvalidArgument(err, "CAPTCHA-ip9Oh", "Errors.Captcha.InitializationFailed")
		}
	case domain.CaptchaTypeDisabled:
	default:
		return zerrors.ThrowError(nil, "CAPTCHA-aCh1D", "Errors.Captcha.UnexpectedType")
	}

	token, err := captchaClient.GetToken(r)
	if err != nil {
		return zerrors.ThrowInvalidArgument(err, "CAPTCHA-Ien0e", "Errors.Captcha.MissingToken")
	}

	ok, err = captchaClient.ValidateToken(token)
	if err != nil {
		return zerrors.ThrowInternal(err, "CAPTCHA-Zohl8", "Errors.Captcha.ValidationFailed")
	}
	if !ok {
		return zerrors.ThrowError(err, "CAPTCHA-oW7ib", "Errors.Captcha.InvalidCaptcha")
	}

	return nil
}
