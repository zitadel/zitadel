package protocbase

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway/descriptor"
)

type GeneratorFunc func(target string, registry *descriptor.Registry, file *descriptor.File) (string, string, error)

type ProtocGenerator interface {
	Generate(target string, registry *descriptor.Registry, file *descriptor.File) (string, string, error)
}

func (f GeneratorFunc) Generate(target string, registry *descriptor.Registry, file *descriptor.File) (string, string, error) {
	return f(target, registry, file) //TODO: in my opinion we should use file.GoPkg here analog https://github.com/grpc-ecosystem/grpc-gateway/blob/0cc2680a4990244dcc7602bad34fef935310c0e8/protoc-gen-grpc-gateway/internal/gengateway/generator.go#L111
}

func parseReq(r io.Reader) (*plugin.CodeGeneratorRequest, error) {
	glog.V(1).Info("Parsing code generator request")

	input, err := ioutil.ReadAll(r)

	if err != nil {
		glog.Errorf("Failed to read code generator request: %v", err)
		return nil, err
	}

	req := &plugin.CodeGeneratorRequest{}

	if err = proto.Unmarshal(input, req); err != nil {
		glog.Errorf("Failed to unmarshal code generator request: %v", err)
		return nil, err
	}

	glog.V(1).Info("Parsed code generator request")

	return req, nil
}

func RunWithBaseTemplate(targetFileNameFmt string, tmpl *template.Template) {
	Run(GeneratorFunc(func(target string, registry *descriptor.Registry, file *descriptor.File) (string, string, error) {
		fileName := fmt.Sprintf(targetFileNameFmt, strings.Split(target, ".")[0])
		fContent, err := GenerateFromBaseTemplate(tmpl, registry, file)
		return fileName, fContent, err
	}))
}

func Run(generator ProtocGenerator) {
	flag.Parse()
	defer glog.Flush()

	req, err := parseReq(os.Stdin)
	if err != nil {
		glog.Fatal(err)
	}

	registry := descriptor.NewRegistry()
	registry.SetAllowDeleteBody(true)
	if err = registry.Load(req); err != nil {
		glog.Fatal(err)
	}

	var result []*plugin.CodeGeneratorResponse_File

	for _, t := range req.FileToGenerate {
		file, err := registry.LookupFile(t)
		if err != nil {
			EmitError(err)
			return
		}

		fName, fContent, err := generator.Generate(t, registry, file)
		if err != nil {
			EmitError(err)
			return
		}

		result = append(result, &plugin.CodeGeneratorResponse_File{
			Name:    &fName,
			Content: &fContent,
		})
	}

	EmitFiles(result)
}

func EmitFiles(out []*plugin.CodeGeneratorResponse_File) {
	EmitResp(&plugin.CodeGeneratorResponse{File: out})
}

func EmitError(err error) {
	EmitResp(&plugin.CodeGeneratorResponse{Error: proto.String(err.Error())})
}

func EmitResp(resp *plugin.CodeGeneratorResponse) {
	buf, err := proto.Marshal(resp)
	if err != nil {
		glog.Fatal(err)
	}
	if _, err := os.Stdout.Write(buf); err != nil {
		glog.Fatal(err)
	}
}
