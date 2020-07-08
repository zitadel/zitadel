package html

import (
	"net/http"
	"path"
	"text/template"

	"github.com/caos/logging"
	"golang.org/x/text/language"
)

const (
	TranslateFn = "t"
)

type Renderer struct {
	Templates  map[string]*template.Template
	translator *Translator
}

func NewRenderer(templatesDir string, tmplMapping map[string]string, funcs map[string]interface{}, translatorConfig TranslatorConfig) (*Renderer, error) {
	var err error
	r := new(Renderer)
	r.translator, err = NewTranslator(translatorConfig)
	if err != nil {
		return nil, err
	}
	r.loadTemplates(templatesDir, tmplMapping, funcs)
	return r, nil
}

func (r *Renderer) RenderTemplate(w http.ResponseWriter, req *http.Request, tmpl *template.Template, data interface{}, reqFuncs map[string]interface{}) {
	reqFuncs = r.registerTranslateFn(req, reqFuncs)
	err := tmpl.Funcs(reqFuncs).Execute(w, data)
	logging.LogWithFields("HTML-lF8F6w", "template", tmpl.Name).OnError(err).Error("error rendering template")
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

func (r *Renderer) loadTemplates(templatesDir string, tmplMapping map[string]string, funcs map[string]interface{}) {
	funcs = r.registerTranslateFn(nil, funcs)
	funcs[TranslateFn] = func(id string, args ...interface{}) string {
		return id
	}
	tmpls := template.Must(template.New("").Funcs(funcs).ParseGlob(path.Join(templatesDir, "*.html")))
	r.Templates = make(map[string]*template.Template, len(tmplMapping))
	for name, file := range tmplMapping {
		r.Templates[name] = tmpls.Lookup(file)
	}
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
