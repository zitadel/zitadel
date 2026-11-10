package renderer

import (
	"context"
	"html/template"
	"io"
	"net/http"
	"os"

	"github.com/zitadel/logging"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/i18n"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	TranslateFn = "t"

	templatesPath = "/templates"
)

type Renderer struct {
	Templates  map[string]*template.Template
	cookieName string
}

func NewRenderer(tmplMapping map[string]string, funcs map[string]interface{}, cookieName string) (*Renderer, error) {
	var err error
	r := &Renderer{cookieName: cookieName}
	err = r.loadTemplates(i18n.LoadFilesystem(i18n.LOGIN), nil, tmplMapping, funcs)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (r *Renderer) RenderTemplate(w http.ResponseWriter, req *http.Request, translator *i18n.Translator, tmpl *template.Template, data interface{}, reqFuncs map[string]interface{}) {
	reqFuncs = r.registerTranslateFn(req, translator, reqFuncs)
	if err := tmpl.Funcs(reqFuncs).Execute(w, data); err != nil {
		logging.Log("RENDE-lF8F6w").WithError(err).WithField("template", tmpl.Name).Error("error rendering template")
	}
}

func (r *Renderer) NewTranslator(ctx context.Context, allowedLanguages []language.Tag) *i18n.Translator {
	return i18n.NewLoginTranslator(authz.GetInstance(ctx).DefaultLanguage(), allowedLanguages, r.cookieName)
}

func (r *Renderer) Localize(translator *i18n.Translator, id string, args map[string]interface{}) string {
	if translator == nil {
		return ""
	}
	return translator.Localize(id, args)
}

func (r *Renderer) AddMessages(translator *i18n.Translator, tag language.Tag, messages ...i18n.Message) error {
	if translator == nil {
		return nil
	}
	return translator.AddMessages(tag, messages...)
}

func (r *Renderer) LocalizeFromRequest(translator *i18n.Translator, req *http.Request, id string, args map[string]interface{}) string {
	if translator == nil {
		return ""
	}
	return translator.LocalizeFromRequest(req, id, args)
}
func (r *Renderer) ReqLang(translator *i18n.Translator, req *http.Request) language.Tag {
	if translator == nil {
		return language.Und
	}
	return translator.Lang(req)
}

func (r *Renderer) loadTemplates(dir http.FileSystem, translator *i18n.Translator, tmplMapping map[string]string, funcs map[string]interface{}) error {
	funcs = r.registerTranslateFn(nil, translator, funcs)
	funcs[TranslateFn] = func(id string, args ...interface{}) string {
		return id
	}
	templatesDir, err := dir.Open(templatesPath)
	if err != nil {
		return zerrors.ThrowNotFound(err, "RENDE-G3aea", "path not found")
	}
	defer templatesDir.Close()
	files, err := templatesDir.Readdir(0)
	if err != nil {
		return zerrors.ThrowNotFound(err, "RENDE-dfR33", "cannot read dir")
	}
	tmpl := template.New("")
	for _, file := range files {
		if err := r.addFileToTemplate(dir, tmpl, tmplMapping, funcs, file); err != nil {
			return zerrors.ThrowNotFound(err, "RENDE-dfTe1", "cannot append file to templates")
		}
	}
	r.Templates = make(map[string]*template.Template, len(tmplMapping))
	for name, file := range tmplMapping {
		r.Templates[name] = tmpl.Lookup(file)
	}
	return nil
}

func (r *Renderer) addFileToTemplate(dir http.FileSystem, tmpl *template.Template, tmplMapping map[string]string, funcs map[string]interface{}, file os.FileInfo) error {
	f, err := dir.Open(templatesPath + "/" + file.Name())
	if err != nil {
		return err
	}
	defer f.Close()
	content, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	tmpl, err = tmpl.New(file.Name()).Funcs(funcs).Parse(string(content))

	return err
}

func (r *Renderer) registerTranslateFn(req *http.Request, translator *i18n.Translator, funcs map[string]interface{}) map[string]interface{} {
	if funcs == nil {
		funcs = make(map[string]interface{})
	}
	if translator == nil {
		return funcs
	}
	funcs[TranslateFn] = func(id string, args ...interface{}) template.HTML {
		m := map[string]interface{}{}
		var key string
		for i, arg := range args {
			if i%2 == 0 {
				key = arg.(string)
				continue
			}
			m[key] = arg
		}
		if r == nil {
			return template.HTML(r.Localize(translator, id, m))
		}
		return template.HTML(r.LocalizeFromRequest(translator, req, id, m))
	}
	return funcs
}
