package renderer

import (
	"io/ioutil"
	"net/http"
	"os"
	"text/template"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/i18n"

	"github.com/caos/logging"
	"golang.org/x/text/language"
)

const (
	TranslateFn = "t"

	templatesPath = "/templates"
)

type Renderer struct {
	Templates  map[string]*template.Template
	translator *i18n.Translator
}

func NewRenderer(dir http.FileSystem, tmplMapping map[string]string, funcs map[string]interface{}, translatorConfig i18n.TranslatorConfig) (*Renderer, error) {
	var err error
	r := new(Renderer)
	r.translator, err = i18n.NewTranslator(dir, translatorConfig)
	if err != nil {
		return nil, err
	}
	err = r.loadTemplates(dir, tmplMapping, funcs)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (r *Renderer) RenderTemplate(w http.ResponseWriter, req *http.Request, tmpl *template.Template, data interface{}, reqFuncs map[string]interface{}) {
	reqFuncs = r.registerTranslateFn(req, reqFuncs)
	if err := tmpl.Funcs(reqFuncs).Execute(w, data); err != nil {
		logging.Log("RENDE-lF8F6w").WithError(err).WithField("template", tmpl.Name).Error("error rendering template")
	}
}

func (r *Renderer) Localize(id string, args map[string]interface{}) string {
	return r.translator.Localize(id, args)
}

func (r *Renderer) LocalizeFromRequest(req *http.Request, id string, args map[string]interface{}) string {
	return r.translator.LocalizeFromRequest(req, id, args)
}
func (r *Renderer) Lang(req *http.Request) language.Tag {
	return r.translator.Lang(req)
}

func (r *Renderer) loadTemplates(dir http.FileSystem, tmplMapping map[string]string, funcs map[string]interface{}) error {
	funcs = r.registerTranslateFn(nil, funcs)
	funcs[TranslateFn] = func(id string, args ...interface{}) string {
		return id
	}
	templatesDir, err := dir.Open(templatesPath)
	if err != nil {
		return errors.ThrowNotFound(err, "RENDE-G3aea", "path not found")
	}
	defer templatesDir.Close()
	files, err := templatesDir.Readdir(0)
	if err != nil {
		return errors.ThrowNotFound(err, "RENDE-dfR33", "cannot read dir")
	}
	tmpl := template.New("")
	for _, file := range files {
		if err := r.addFileToTemplate(dir, tmpl, tmplMapping, funcs, file); err != nil {
			return errors.ThrowNotFound(err, "RENDE-dfTe1", "cannot append file to templates")
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
	content, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	tmpl, err = tmpl.New(file.Name()).Funcs(funcs).Parse(string(content))
	if err != nil {
		return err
	}
	return nil
}

func (r *Renderer) registerTranslateFn(req *http.Request, funcs map[string]interface{}) map[string]interface{} {
	if funcs == nil {
		funcs = make(map[string]interface{})
	}
	funcs[TranslateFn] = func(id string, args ...interface{}) string {
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
			return r.Localize(id, m)
		}
		return r.LocalizeFromRequest(req, id, m)
	}
	return funcs
}
