package main

import (
	"strings"
	"unicode"

	openapiv2 "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

// humanizeRPCName generates a short description from the RPC name.
func humanizeRPCName(rpcName, resource string) string {
	// Split CamelCase, keeping consecutive uppercase as one word (acronyms).
	var words []string
	runes := []rune(rpcName)
	start := 0
	for i := 1; i < len(runes); i++ {
		if unicode.IsUpper(runes[i]) {
			// If previous was also upper and next (if exists) is also upper or end, continue the acronym.
			if unicode.IsUpper(runes[i-1]) && (i+1 >= len(runes) || unicode.IsUpper(runes[i+1])) {
				continue
			}
			// If previous was upper but next is lower, the current char starts a new word.
			if unicode.IsUpper(runes[i-1]) {
				// The previous run was an acronym, split before current char.
				words = append(words, string(runes[start:i]))
				start = i
				continue
			}
			words = append(words, string(runes[start:i]))
			start = i
		}
	}
	words = append(words, string(runes[start:]))

	// Lowercase all words except known acronyms, and join.
	for i := range words {
		words[i] = strings.ToLower(words[i])
	}
	result := strings.Join(words, " ")
	if len(result) > 0 {
		result = strings.ToUpper(result[:1]) + result[1:]
	}
	return result
}

// toKebab converts snake_case or CamelCase to kebab-case.
func toKebab(s string) string {
	// First handle snake_case.
	s = strings.ReplaceAll(s, "_", "-")
	var result strings.Builder
	runes := []rune(s)
	for i, r := range runes {
		if unicode.IsUpper(r) && i > 0 {
			prev := runes[i-1]
			if prev == '-' {
				// Already have a separator.
			} else if unicode.IsUpper(prev) {
				// In a run of uppercase. Only add dash if next char is lowercase (end of acronym).
				if i+1 < len(runes) && unicode.IsLower(runes[i+1]) {
					result.WriteRune('-')
				}
			} else {
				result.WriteRune('-')
			}
		}
		result.WriteRune(unicode.ToLower(r))
	}
	return result.String()
}

func titleCase(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func extractComment(loc protogen.CommentSet) string {
	leading := strings.TrimSpace(string(loc.Leading))
	if leading == "" {
		return ""
	}
	// Take just the first line.
	if idx := strings.IndexByte(leading, '\n'); idx >= 0 {
		leading = leading[:idx]
	}
	return strings.TrimSpace(leading)
}

// jsonFieldName returns the JSON field name for a proto field.
func jsonFieldName(field *protogen.Field) string {
	return string(field.Desc.JSONName())
}

// isMethodDeprecatedOpenAPI returns true when the RPC has deprecated:true in its
// grpc-gateway openapiv2_operation annotation (ZITADEL's deprecation convention).
func isMethodDeprecatedOpenAPI(method *protogen.Method) bool {
	opts := method.Desc.Options()
	if opts == nil {
		return false
	}
	ext := proto.GetExtension(opts, openapiv2.E_Openapiv2Operation)
	op, ok := ext.(*openapiv2.Operation)
	return ok && op != nil && op.Deprecated
}

// inferResourceName strips the verb prefix and service name to get the resource.
func inferResourceName(rpcName, serviceName string) string {
	// Remove common service suffix.
	resource := serviceName
	resource = strings.TrimSuffix(resource, "Service")
	return resource
}

// rpcNameToVerbAndSuffix maps a proto RPC name to a CLI verb and optional suffix.
func rpcNameToVerbAndSuffix(name, serviceName string) (verb, suffix string) {
	// Strip "Service" suffix from service name for matching
	resource := strings.TrimSuffix(serviceName, "Service")

	prefixes := []struct {
		prefix string
		verb   string
	}{
		{"List", "list"},
		{"Create", "create"},
		{"Add", "create"},
		{"Get", "get"},
		{"Update", "update"},
		{"Delete", "delete"},
		{"Remove", "remove"},
		{"Deactivate", "deactivate"},
		{"Activate", "activate"},
		{"Reactivate", "reactivate"},
		{"Generate", "generate"},
		{"Set", "set"},
		{"Verify", "verify"},
	}
	for _, p := range prefixes {
		if strings.HasPrefix(name, p.prefix) {
			remainder := name[len(p.prefix):]
			remainder = strings.TrimPrefix(remainder, resource+"s")
			remainder = strings.TrimPrefix(remainder, resource)
			if remainder != "" {
				suffix = toKebab(remainder)
			}
			return p.verb, suffix
		}
	}
	return "", ""
}

// isWellKnownProtoType returns true for google.protobuf.* messages.
func isWellKnownProtoType(msg *protogen.Message) bool {
	return strings.HasPrefix(string(msg.Desc.FullName()), "google.protobuf.")
}

// isFieldDeprecated returns true when the proto field/method is marked deprecated.
func isFieldDeprecated(field any) bool {
	var opts protoreflect.ProtoMessage
	switch f := field.(type) {
	case *protogen.Field:
		opts = f.Desc.Options()
	case *protogen.Method:
		opts = f.Desc.Options()
	}
	if opts == nil {
		return false
	}
	switch o := opts.(type) {
	case *descriptorpb.FieldOptions:
		return o.GetDeprecated()
	case *descriptorpb.MethodOptions:
		return o.GetDeprecated()
	}
	return false
}

// importAlias derives an import alias from a Go import path.
func importAlias(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) >= 2 {
		return parts[len(parts)-2] + parts[len(parts)-1]
	}
	return parts[len(parts)-1]
}
