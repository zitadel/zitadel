package main

import (
	"bytes"
	"flag"
	"os"
	"path/filepath"
	"testing"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

var update = flag.Bool("update", false, "update .golden files")

func TestGolden(t *testing.T) {
	// 1. Setup a synthetic FileDescriptorProto
	fd := &descriptorpb.FileDescriptorProto{
		Name:    proto.String("test.proto"),
		Package: proto.String("zitadel.user.v2"), // Must be in v2ServiceFilter
		Service: []*descriptorpb.ServiceDescriptorProto{
			{
				Name: proto.String("UserService"),
				Method: []*descriptorpb.MethodDescriptorProto{
					{
						Name:       proto.String("GetUserByID"),
						InputType:  proto.String(".zitadel.user.v2.GetUserByIDRequest"),
						OutputType: proto.String(".zitadel.user.v2.GetUserByIDResponse"),
					},
				},
			},
		},
		MessageType: []*descriptorpb.DescriptorProto{
			{
				Name: proto.String("GetUserByIDRequest"),
				Field: []*descriptorpb.FieldDescriptorProto{
					stringField("user_id", 1),
				},
			},
			{
				Name: proto.String("GetUserByIDResponse"),
				Field: []*descriptorpb.FieldDescriptorProto{
					messageField("user", 1, ".zitadel.user.v2.User"),
				},
			},
			{
				Name: proto.String("User"),
				Field: []*descriptorpb.FieldDescriptorProto{
					stringField("id", 1),
					stringField("username", 2),
				},
			},
		},
	}

	plugin := buildPlugin(t, fd)
	file := plugin.Files[0]
	service := file.Services[0]
	
	cfg := v2ServiceFilter[string(file.Desc.Package())]
	sd := buildServiceData(service, file, cfg)
	
	// 2. Execute template
	cmdTmpl := loadTemplate("cmd", cmdTemplate)
	var buf bytes.Buffer
	if err := cmdTmpl.Execute(&buf, &sd); err != nil {
		t.Fatalf("executing cmd template: %v", err)
	}
	
	got := buf.Bytes()
	goldenPath := filepath.Join("testdata", "golden", "cmd_users.go.golden")
	
	if *update {
		if err := os.MkdirAll(filepath.Dir(goldenPath), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(goldenPath, got, 0644); err != nil {
			t.Fatal(err)
		}
	}
	
	want, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("reading golden file: %v. Run with -update to create it.", err)
	}
	
	if !bytes.Equal(got, want) {
		t.Errorf("generated output does not match golden file %s", goldenPath)
		// We could use a diff tool here for better output
	}
}
