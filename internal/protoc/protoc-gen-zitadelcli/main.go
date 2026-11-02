package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

var (
	//go:embed cmd_dynamic.go.tmpl
	cmdTemplate []byte
)

func main() {
	input, _ := io.ReadAll(os.Stdin)
	var req pluginpb.CodeGeneratorRequest
	if err := proto.Unmarshal(input, &req); err != nil {
		panic(err)
	}

	opts := protogen.Options{}
	plugin, err := opts.New(&req)
	if err != nil {
		panic(err)
	}
	plugin.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

	// Use dynamic (metadata-only) template by default.
	cmdTmpl := loadTemplate("cmd", cmdTemplate)

	for _, file := range plugin.Files {
		if !file.Generate {
			continue
		}
		pkgName := string(file.Desc.Package())
		cfg, ok := v2ServiceFilter[pkgName]
		if !ok {
			continue
		}

		for _, service := range file.Services {
			sd := buildServiceData(service, file, cfg)
			if len(sd.Methods) == 0 {
				continue
			}

			// Generate the command file for this service.
			var buf bytes.Buffer
			if err := cmdTmpl.Execute(&buf, &sd); err != nil {
				panic(fmt.Sprintf("executing cmd template for %s: %v", sd.ServiceName, err))
			}
			outFile := plugin.NewGeneratedFile(
				fmt.Sprintf("cmd_%s.go", cfg.resourceName),
				protogen.GoImportPath(sd.GoImportPath),
			)
			outFile.Write(buf.Bytes())
		}
	}

	out, err := proto.Marshal(plugin.Response())
	if err != nil {
		panic(err)
	}
	if _, err := os.Stdout.Write(out); err != nil {
		panic(err)
	}
}

func loadTemplate(name string, data []byte) *template.Template {
	funcMap := template.FuncMap{
		"kebab":  toKebab,
		"lower":  strings.ToLower,
		"title":  titleCase,
		"quote":  func(s string) string { return fmt.Sprintf("%q", s) },
		"join":   strings.Join,
		"repeat": strings.Repeat,
		"contains": strings.Contains,
		"trimSuffix": strings.TrimSuffix,
		"dict": func(values ...interface{}) (map[string]interface{}, error) {
			if len(values)%2 != 0 {
				return nil, fmt.Errorf("invalid dict call")
			}
			dict := make(map[string]interface{}, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, fmt.Errorf("dict keys must be strings")
				}
				dict[key] = values[i+1]
			}
			return dict, nil
		},
	}
	return template.Must(template.New(name).Funcs(funcMap).Parse(string(data)))
}
