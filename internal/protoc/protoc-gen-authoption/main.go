package main

import (
	base "github.com/zitadel/zitadel/internal/protoc/protoc-base"
	"github.com/zitadel/zitadel/internal/protoc/protoc-gen-authoption/authoption"
)

const (
	fileName = "%v.pb.authoptions.go"
)

func main() {
	base.RegisterExtension(authoption.E_AuthOption)
	base.RunWithBaseTemplate(fileName, base.LoadTemplate(templatesAuth_method_mappingGoTmplBytes()))
}
