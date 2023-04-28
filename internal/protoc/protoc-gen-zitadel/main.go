package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"os"
	"text/template"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/zitadel/zitadel/internal/protoc/protoc-gen-zitadel/zitadel"
)

var (
	//go:embed zitadel.pb.go.tmpl
	zitadelTemplate []byte
)

type authMethods struct {
	GoPackageName       string
	ProtoPackageName    string
	ServiceName         string
	AuthOptions         []authOption
	AuthContext         []authContext
	CustomHTTPResponses []httpResponse
}

type authOption struct {
	Name           string
	Permission     string
	CheckFieldName string
}

type authContext struct {
	Name      string
	OrgMethod string
}

type httpResponse struct {
	Name string
	Code int32
}

func main() {

	input, _ := io.ReadAll(os.Stdin)
	var req pluginpb.CodeGeneratorRequest
	err := proto.Unmarshal(input, &req)
	if err != nil {
		panic(err)
	}

	opts := protogen.Options{}
	plugin, err := opts.New(&req)
	if err != nil {
		panic(err)
	}
	plugin.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

	authTemp := loadTemplate(zitadelTemplate)

	for _, file := range plugin.Files {

		var buf bytes.Buffer

		var methods authMethods
		for _, service := range file.Services {
			methods.ServiceName = service.GoName
			methods.GoPackageName = string(file.GoPackageName)
			methods.ProtoPackageName = *file.Proto.Package
			for _, method := range service.Methods {
				if options := method.Desc.Options().(*descriptorpb.MethodOptions); options != nil {
					ext := proto.GetExtension(options, zitadel.E_Options).(*zitadel.Options)
					if ext == nil {
						continue
					}
					if ext.AuthOption != nil {
						methods.AuthOptions = append(methods.AuthOptions, authOption{Name: string(method.Desc.Name()), Permission: ext.AuthOption.Permission /*CheckFieldName: authExt.CheckFieldName*/})
						if ext.AuthOption.OrgField != "" {
							orgMethod := buildAuthContextField(method.Input.Fields, ext.AuthOption.OrgField)
							if orgMethod != "" {
								methods.AuthContext = append(methods.AuthContext, authContext{Name: string(method.Input.Desc.Name()), OrgMethod: orgMethod})
							}
						}
					}
					if ext.HttpResponse != nil {
						methods.CustomHTTPResponses = append(methods.CustomHTTPResponses, httpResponse{Name: string(method.Output.Desc.Name()), Code: ext.HttpResponse.SuccessCode})
					}
				}
			}
		}
		if len(methods.AuthOptions) > 0 {
			err = authTemp.Execute(&buf, &methods)
			if err != nil {
				panic(err)
			}

			filename := file.GeneratedFilenamePrefix + ".pb.zitadel.go"
			file := plugin.NewGeneratedFile(filename, ".")

			_, err = file.Write(buf.Bytes())
			if err != nil {
				panic(err)
			}
		}
	}

	// Generate a response from our plugin and marshall as protobuf
	stdout := plugin.Response()
	out, err := proto.Marshal(stdout)
	if err != nil {
		panic(err)
	}

	// Write the response to stdout, to be picked up by protoc
	_, err = fmt.Fprint(os.Stdout, string(out))
	if err != nil {
		panic(err)
	}
}

func loadTemplate(templateData []byte) *template.Template {
	return template.Must(template.New("").
		Parse(string(templateData)))
}

func buildAuthContextField(fields []*protogen.Field, fieldName string) string {
	for _, field := range fields {
		if string(field.Desc.Name()) == fieldName {
			return ".Get" + field.GoName + "()"
		}
	}
	return ""
}
