package handler

import (
	caos_errs "github.com/caos/zitadel/internal/errors"
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
	authSession, err := l.getAuthSessionAndParseData(r, data)
	if err != nil {
		l.renderError(w, r, authSession, err)
		return
	}
	if data.Password != data.Password2 {
		err := caos_errs.ThrowInvalidArgument(nil, "VIEW-KaGue", "passwords dont match")
		l.renderRegister(w, r, authSession, data, err)
		return
	}
	_, err = l.service.Auth.RegisterUser(r.Context(), data.toProto())
	if err != nil {
		l.renderRegister(w, r, authSession, data, err)
		return
	}
	// authSession.UserSession.User.UserName = user.UserName //TODO: ?
	l.renderNextStep(w, r, authSession)
}

func (l *Login) renderRegister(w http.ResponseWriter, r *http.Request, authSession *model.AuthSession, formData *registerFormData, err error) {
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

func (d registerFormData) toProto() *auth_api.RegisterUserRequest { //TODO: refactoring?
	return &auth_api.RegisterUserRequest{
		Email:             d.Email,
		FirstName:         d.Firstname,
		LastName:          d.Lastname,
		NickName:          d.Nickname,
		PreferredLanguage: d.Language,
		Gender:            auth_api.Gender(d.Gender),
		Password:          d.Password,
	}
}
