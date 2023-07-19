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

	"github.com/zitadel/zitadel/internal/protoc/protoc-gen-authoption/authoption"
)

var (
	//go:embed auth_method_mapping.go.tmpl
	authTemplate []byte
)

type authMethods struct {
	GoPackageName    string
	ProtoPackageName string
	ServiceName      string
	AuthOptions      []authOption
}

type authOption struct {
	Name           string
	Permission     string
	CheckFieldName string
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

	authTemp := loadTemplate(authTemplate)

	for _, file := range plugin.Files {

		var buf bytes.Buffer

		var methods authMethods
		for _, service := range file.Services {
			methods.ServiceName = service.GoName
			methods.GoPackageName = string(file.GoPackageName)
			methods.ProtoPackageName = *file.Proto.Package
			for _, method := range service.Methods {
				if options := method.Desc.Options().(*descriptorpb.MethodOptions); options != nil {
					ext := proto.GetExtension(options, authoption.E_AuthOption).(*authoption.AuthOption)
					if ext != nil {
						methods.AuthOptions = append(methods.AuthOptions, authOption{Name: string(method.Desc.Name()), Permission: ext.Permission, CheckFieldName: ext.CheckFieldName})
					}
				}
			}
		}
		if len(methods.AuthOptions) > 0 {
			authTemp.Execute(&buf, &methods)

			filename := file.GeneratedFilenamePrefix + ".pb.authoptions.go"
			file := plugin.NewGeneratedFile(filename, ".")

			file.Write(buf.Bytes())
		}
	}

	// Generate a response from our plugin and marshall as protobuf
	stdout := plugin.Response()
	out, err := proto.Marshal(stdout)
	if err != nil {
		panic(err)
	}

	// Write the response to stdout, to be picked up by protoc
	fmt.Fprintf(os.Stdout, string(out))
}

func loadTemplate(templateData []byte) *template.Template {
	return template.Must(template.New("").
		Parse(string(templateData)))
}
