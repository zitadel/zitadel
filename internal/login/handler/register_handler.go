package handler

import (
	"github.com/caos/zitadel/internal/auth_request/model"
	caos_errs "github.com/caos/zitadel/internal/errors"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"golang.org/x/text/language"
	"net/http"
)

const (
	tmplRegister = "register"
)

type registerFormData struct {
	Email     string `schema:"email"`
	Firstname string `schema:"firstname"`
	Lastname  string `schema:"lastname"`
	Nickname  string `schema:"nickname"`
	Language  string `schema:"language"`
	Gender    int32  `schema:"gender"`
	Password  string `schema:"password"`
	Password2 string `schema:"password2"`
}

type registerData struct {
	baseData
	registerFormData
}

func (l *Login) handleRegister(w http.ResponseWriter, r *http.Request) {
	data := new(registerFormData)
	authSession, err := l.getAuthRequestAndParseData(r, data)
	if err != nil {
		l.renderError(w, r, authSession, err)
		return
	}
	if data.Password != data.Password2 {
		err := caos_errs.ThrowInvalidArgument(nil, "VIEW-KaGue", "passwords dont match")
		l.renderRegister(w, r, authSession, data, err)
		return
	}
	//TODO: How to get ResourceOwner?
	_, err = l.authRepo.Register(r.Context(), data.toUserModel(), "GlobalResourceOwner")
	if err != nil {
		l.renderRegister(w, r, authSession, data, err)
		return
	}
	// authSession.UserSession.User.UserName = user.UserName //TODO: ?
	l.renderNextStep(w, r, authSession)
}

func (l *Login) renderRegister(w http.ResponseWriter, r *http.Request, authSession *model.AuthRequest, formData *registerFormData, err error) {
	var errType, errMessage string
	if err != nil {
		errMessage = err.Error()
	}
	if formData == nil {
		formData = new(registerFormData)
	}
	data := registerData{
		baseData:         l.getBaseData(r, authSession, "Register", errType, errMessage),
		registerFormData: *formData,
	}
	funcs := map[string]interface{}{
		"selectedLanguage": func(l string) bool {
			if formData == nil {
				return false
			}
			return formData.Language == l
		},
		"selectedGender": func(g int32) bool {
			if formData == nil {
				return false
			}
			return formData.Gender == g
		},
	}
	l.renderer.RenderTemplate(w, r, l.renderer.Templates[tmplRegister], data, funcs)
}

func (d registerFormData) toUserModel() *usr_model.User {
	return &usr_model.User{
		Profile: &usr_model.Profile{
			FirstName:         d.Firstname,
			LastName:          d.Lastname,
			NickName:          d.Nickname,
			PreferredLanguage: language.Make(d.Language),
			Gender:            usr_model.Gender(d.Gender),
		},
		Password: &usr_model.Password{
			SecretString: d.Password,
		},
		Email: &usr_model.Email{
			EmailAddress: d.Email,
		},
	}
}
