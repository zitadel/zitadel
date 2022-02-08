package protocbase

import (
	"bytes"
	"fmt"
	"text/template"
	"time"

	"github.com/Masterminds/sprig"
	"github.com/golang/protobuf/proto"
	"github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway/descriptor"
	"golang.org/x/tools/imports"
)

var extensions = map[string]*proto.ExtensionDesc{}

type BaseTemplateData struct {
	Now  time.Time
	File *descriptor.File

	registry *descriptor.Registry
}

var templateFuncs = map[string]interface{}{
	"option": getOption,
	"duration": func(s string) interface{} {
		d, _ := time.ParseDuration(s)
		if d == 0 {
			return 0
		}
		return fmt.Sprintf("time.Duration(%d)", d.Nanoseconds())
	},
}

func RegisterTmplFunc(name string, f interface{}) {
	if _, found := templateFuncs[name]; found {
		panic(fmt.Sprintf("func with name %v is already registered", name))
	}

	templateFuncs[name] = f
}

func RegisterExtension(ext *proto.ExtensionDesc) {
	extensions[ext.Name] = ext
}

func GetBaseTemplateData(registry *descriptor.Registry, file *descriptor.File) *BaseTemplateData {
	return &BaseTemplateData{
		Now:      time.Now().UTC(),
		File:     file,
		registry: registry,
	}
}

func getOption(opts proto.Message, extName string) interface{} {
	extDesc := extensions[extName]

	if !proto.HasExtension(opts, extDesc) {
		return nil
	}

	ext, err := proto.GetExtension(opts, extDesc)
	if err != nil {
		panic(err)
	}

	return ext
}

func (data *BaseTemplateData) ResolveMsgType(msgType string) string {
	msg, err := data.registry.LookupMsg(data.File.GetPackage(), msgType)
	if err != nil {
		panic(err)
	}

	return msg.GoType(data.File.GoPkg.Path)
}

func (data *BaseTemplateData) ResolveFile(fileName string) *descriptor.File {
	file, err := data.registry.LookupFile(fileName)
	if err != nil {
		panic(err)
	}

	return file
}

func LoadTemplate(templateData []byte, err error) *template.Template {
	if err != nil {
		panic(err)
	}

	return template.Must(template.New("").
		Funcs(sprig.TxtFuncMap()).
		Funcs(templateFuncs).
		Parse(string(templateData)))
}

func GenerateFromTemplate(tmpl *template.Template, data interface{}) (string, error) {
	var tpl bytes.Buffer
	err := tmpl.Execute(&tpl, data)
	if err != nil {
		return "", err
	}

	tmplResult := tpl.Bytes()
	tmplResult, err = imports.Process(".", tmplResult, nil)
	return string(tmplResult), err
}

func GenerateFromBaseTemplate(tmpl *template.Template, registry *descriptor.Registry, file *descriptor.File) (string, error) {
	return GenerateFromTemplate(tmpl, GetBaseTemplateData(registry, file))
}
