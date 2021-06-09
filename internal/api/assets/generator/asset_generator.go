package main

import (
	"flag"
	"io"
	"os"
	"text/template"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/config"
)

var (
	directory = flag.String("directory", "./", "working directory: asset.yaml must be in this directory, files will be generated into parent directory")
	assets    = flag.String("assets", "../../../../docs/docs/apis/assets/assets.md", "path where the assets.md will be generated")
)

func main() {
	flag.Parse()
	configFile := *directory + "asset.yaml"
	authz, err := os.OpenFile(*directory+"../authz.go", os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0755)
	logging.Log("ASSETS-Gn31f").OnError(err).Fatal("cannot open authz file")
	router, err := os.OpenFile(*directory+"../router.go", os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0755)
	logging.Log("ASSETS-ABen3").OnError(err).Fatal("cannot open router file")
	docs, err := os.OpenFile(*assets, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0755)
	logging.Log("ASSETS-Dfvsd").OnError(err).Fatal("cannot open docs file")
	GenerateAssetHandler(configFile, authz, router, docs)
}

type Method struct {
	Path        string
	Feature     string
	HasDarkMode bool
	Handlers    []Handler
}

type Handler struct {
	Name       string
	Comment    string
	Type       HandlerType
	Permission string
}

func (a Handler) Method() string {
	if a.Type == MethodTypeUpload {
		return "POST"
	}
	return "GET"
}

func (a Handler) PathSuffix() string {
	if a.Type == MethodTypePreview {
		return "/_preview"
	}
	return ""
}

func (a Handler) MethodReturn() string {
	if a.Type == MethodTypeUpload {
		return "Uploader"
	}
	if a.Type == MethodTypeDownload {
		return "Downloader"
	}
	if a.Type == MethodTypePreview {
		return "Downloader"
	}
	return ""
}

func (a Handler) HandlerType() string {
	if a.Type == MethodTypeUpload {
		return "UploadHandleFunc"
	}
	if a.Type == MethodTypeDownload {
		return "DownloadHandleFunc"
	}
	if a.Type == MethodTypePreview {
		return "DownloadHandleFunc"
	}
	return ""
}

type HandlerType string

const (
	MethodTypeUpload   = "upload"
	MethodTypeDownload = "download"
	MethodTypePreview  = "preview"
)

type Services map[string]Service

type Service struct {
	Prefix  string
	Methods map[string]Method
}

func GenerateAssetHandler(configFilePath string, authz, router, docs io.Writer) {
	conf := new(struct {
		Services Services
	})
	err := config.Read(conf, configFilePath)
	logging.Log("ASSETS-Dgbn4").OnError(err).Fatal("cannot read config")
	tmplAuthz, err := template.New("").Parse(authzTmpl)
	logging.Log("ASSETS-BGbbg").OnError(err).Fatal("cannot parse authz template")
	tmplRouter, err := template.New("").Parse(routerTmpl)
	logging.Log("ASSETS-gh4rq").OnError(err).Fatal("cannot parse router template")
	tmplDocs, err := template.New("").Parse(docsTmpl)
	logging.Log("ASSETS-FGdgs").OnError(err).Fatal("cannot parse docs template")
	data := &struct {
		GoPkgName string
		Name      string
		Prefix    string
		Services  Services
	}{
		GoPkgName: "assets",
		Name:      "AssetsService",
		Prefix:    "/assets/v1",
		Services:  conf.Services,
	}
	err = tmplAuthz.Execute(authz, data)
	logging.Log("ASSETS-BHngj").OnError(err).Fatal("cannot generate authz")
	err = tmplRouter.Execute(router, data)
	logging.Log("ASSETS-Bfd41").OnError(err).Fatal("cannot generate router")
	err = tmplDocs.Execute(docs, data)
	logging.Log("ASSETS-Bfd41").OnError(err).Fatal("cannot generate docs")
}

const authzTmpl = `package {{.GoPkgName}}

import (
	"github.com/caos/zitadel/internal/api/authz"
)

/**
 * {{.Name}}
 */

{{ $prefix := .Prefix }}
var {{.Name}}_AuthMethods = authz.MethodMapping {
    {{ range $service := .Services}}
	{{ range $method := .Methods}}
	{{ range $handler := .Handlers}}
    {{ if (or $method.Feature $handler.Permission) }}
    	"{{$handler.Method}}:{{$prefix}}{{$service.Prefix}}{{$method.Path}}{{$handler.PathSuffix}}": authz.Option{
               Permission: "{{$handler.Permission}}",
               Feature:    "{{$method.Feature}}",
        },
	{{ if $method.HasDarkMode }}
		"{{$handler.Method}}:{{$prefix}}{{$service.Prefix}}{{$method.Path}}/dark{{$handler.PathSuffix}}": authz.Option{
               Permission: "{{$handler.Permission}}",
               Feature:    "{{$method.Feature}}",
        },
	{{end}}
	{{end}}
    {{end}}
    {{end}}
    {{end}}
}
`

const routerTmpl = `package {{.GoPkgName}}

import (
	"github.com/gorilla/mux"

	http_mw "github.com/caos/zitadel/internal/api/http/middleware"
	"github.com/caos/zitadel/internal/command"
	"github.com/caos/zitadel/internal/static"
)

type {{.Name}} interface {
	AuthInterceptor() *http_mw.AuthInterceptor
	Commands() *command.Commands
	ErrorHandler() ErrorHandler
	Storage() static.Storage
    
	{{ range $service := .Services}}
	{{ range $methodName, $method := .Methods}}
	{{ range $handler := .Handlers}}
	{{$handler.Name}}{{$methodName}}() {{if $handler.MethodReturn}}{{$handler.MethodReturn}}{{end}}
	{{ if $method.HasDarkMode }}
	{{$handler.Name}}{{$methodName}}Dark() {{if $handler.MethodReturn}}{{$handler.MethodReturn}}{{end}}
	{{ end }}
    {{ end }}
	{{ end }}
	{{ end }}
}

func RegisterRoutes(router *mux.Router, s {{.Name}}) {

	router.Use(s.AuthInterceptor().Handler)

	{{ range $service := .Services}}
	{{ range $methodName, $method := .Methods}}
	{{ range $handler := .Handlers}}
	router.Path("{{$service.Prefix}}{{$method.Path}}{{$handler.PathSuffix}}").Methods("{{$handler.Method}}").HandlerFunc({{if $handler.HandlerType}}{{$handler.HandlerType}}(s, {{end}}s.{{$handler.Name}}{{$methodName}}(){{if $handler.HandlerType}}){{end}})	
	{{ if $method.HasDarkMode }}
	router.Path("{{$service.Prefix}}{{$method.Path}}/dark{{$handler.PathSuffix}}").Methods("{{$handler.Method}}").HandlerFunc({{if $handler.HandlerType}}{{$handler.HandlerType}}(s, {{end}}s.{{$handler.Name}}{{$methodName}}Dark(){{if $handler.HandlerType}}){{end}})
    {{ end }}
	{{ end }}
	{{ end }}
	{{ end }}
}
`

const docsTmpl = `---
title: zitadel/admin.proto
---

## {{.Name}}

	{{ range $service := .Services}}
	{{ range $methodName, $method := .Methods}}
	{{ range $handler := .Handlers}}

### {{$handler.Name}}{{$methodName}}()

> {{$handler.Name}}{{$methodName}}()

{{$handler.Method}}: {{$service.Prefix}}{{$method.Path}}{{$handler.PathSuffix}}
{{ if $method.HasDarkMode }}
### {{$handler.Name}}{{$methodName}}()

> {{$handler.Name}}{{$methodName}}Dark()

{{$handler.Method}}: {{$service.Prefix}}{{$method.Path}}/dark{{$handler.PathSuffix}}
 	{{ end }}
 	{{ end }}
	{{ end }}
	{{ end }}
`
