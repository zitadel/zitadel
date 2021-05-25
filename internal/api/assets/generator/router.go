package main

import (
	"io"
	"os"
	"text/template"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/config"
)

func main() {
	configFile := "./asset.yaml"
	output := os.Stdout
	output2 := output
	output, err := os.OpenFile("../authz.go", os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0755)
	logging.Log("ASSETS-DAg42").OnError(err).Fatal("cannot read config")
	output2, err = os.OpenFile("../router.go", os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0755)
	logging.Log("ASSETS-DAg42").OnError(err).Fatal("cannot read config")
	GenerateAssetHandler(configFile, output, output2)
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

func GenerateAssetHandler(configFilePath string, output io.Writer, output2 io.Writer) {
	conf := new(struct {
		Services Services
	})
	err := config.Read(conf, configFilePath)
	logging.Log("ASSETS-DAg42").OnError(err).Fatal("cannot read config")
	t, err := template.New("").Parse(authzTmpl)
	logging.Log("ASSETS-DAg42").OnError(err).Fatal("cannot read config")
	t2, err := template.New("").Parse(routerTmpl)
	logging.Log("ASSETS-DAg42").OnError(err).Fatal("cannot read config")
	err = t.Execute(output, struct {
		GoPkgName string
		Name      string
		Prefix    string
		Services  Services
	}{
		GoPkgName: "assets",
		Name:      "AssetsService",
		Prefix:    "/assets/v1",
		Services:  conf.Services,
	})
	logging.Log("ASSETS-DAg42").OnError(err).Fatal("cannot read config")
	err = t2.Execute(output2, struct {
		GoPkgName string
		Name      string
		Services  Services
	}{
		GoPkgName: "assets",
		Name:      "AssetsService",
		Services:  conf.Services,
	})
	logging.Log("ASSETS-DAg42").OnError(err).Fatal("cannot read config")
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
