package runtime

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"

	"github.com/zitadel/zitadel/apps/cli/internal/auth"
	"github.com/zitadel/zitadel/apps/cli/internal/client"
	"github.com/zitadel/zitadel/apps/cli/internal/config"
	"github.com/zitadel/zitadel/apps/cli/internal/output"
)

// BuildCommands turns all registered CommandSpecs into Cobra commands and
// attaches them to the given root command.
func BuildCommands(root *cobra.Command, getCfg func() *config.Config, getOutput func() string) {
	allSpecs := AllSpecs()

	// Group specs by Group name.
	groups := make(map[string][]CommandSpec)
	groupDescs := make(map[string]string)
	for _, s := range allSpecs {
		groups[s.Group] = append(groups[s.Group], s)
		if s.GroupDescription != "" {
			groupDescs[s.Group] = s.GroupDescription
		}
	}

	for groupName, groupSpecs := range groups {
		groupCmd := &cobra.Command{
			Use:   groupName,
			Short: "Manage " + groupDescs[groupName],
		}
		for _, spec := range groupSpecs {
			spec := spec // capture for closures
			groupCmd.AddCommand(buildCommand(spec, getCfg, getOutput))
		}
		root.AddCommand(groupCmd)
	}
}

// buildCommand creates a single Cobra command from a CommandSpec.
func buildCommand(spec CommandSpec, getCfg func() *config.Config, getOutput func() string) *cobra.Command {
	// Resolve proto descriptors for the method.
	methodDesc, reqDesc, err := resolveMethod(spec.FullMethodName)
	if err != nil {
		// If the method can't be resolved, create a stub command that errors.
		return &cobra.Command{
			Use:   spec.Verb,
			Short: spec.Short,
			RunE: func(_ *cobra.Command, _ []string) error {
				return fmt.Errorf("command unavailable: %w", err)
			},
		}
	}

	if spec.HasOneofSubcmds && len(spec.OneofGroups) > 0 {
		return buildOneofCommand(spec, methodDesc, reqDesc, getCfg, getOutput)
	}
	return buildStandardCommand(spec, methodDesc, reqDesc, getCfg, getOutput)
}

// buildStandardCommand builds a leaf command (no oneof subcommands).
func buildStandardCommand(spec CommandSpec, methodDesc protoreflect.MethodDescriptor, reqDesc protoreflect.MessageDescriptor, getCfg func() *config.Config, getOutput func() string) *cobra.Command {
	// Track flag values: map from flag name → pointer to value.
	flagValues := make(map[string]interface{})

	use := spec.Verb
	for _, pa := range spec.PositionalArgs {
		use += " <" + pa.ProtoFieldName + ">"
	}

	cmd := &cobra.Command{
		Use:     use,
		Short:   spec.Short,
		Long:    spec.Long,
		Example: spec.Example,
		RunE: func(cmd *cobra.Command, args []string) error {
			req := dynamicpb.NewMessage(reqDesc)

			// Handle JSON input.
			useJSON, requestJSON, err := readJSONInput(cmd)
			if err != nil {
				return err
			}
			if useJSON {
				if err := protojson.Unmarshal(requestJSON, req); err != nil {
					return fmt.Errorf("parsing JSON input: %w", err)
				}
			}

			// Handle positional args.
			if err := applyPositionalArgs(req, reqDesc, spec.PositionalArgs, args, useJSON); err != nil {
				return err
			}

			// Handle flags (only when not using JSON input).
			if !useJSON {
				if err := checkRequiredFlags(cmd, reqDesc, spec); err != nil {
					return err
				}
				if err := ValidateAllFlags(flagValues); err != nil {
					return err
				}
				applyFlags(cmd, req, reqDesc, flagValues)
			}

			return execute(cmd, req, methodDesc, reqDesc, spec, getCfg, getOutput)
		},
	}

	// Register flags from proto descriptor.
	aliases := make(map[string][]string)
	registerFlags(cmd.Flags(), reqDesc, "", flagValues, spec.PositionalArgs, 0, true, aliases)
	registerAliases(cmd.Flags(), flagValues, aliases)

	// Register filter convenience flags.
	for _, fc := range spec.FilterConvenience {
		flagValues["_filter_"+fc.FlagName] = cmd.Flags().String(fc.FlagName, "", fc.Help)
	}

	return cmd
}

// buildOneofCommand builds a parent command with variant subcommands.
func buildOneofCommand(spec CommandSpec, methodDesc protoreflect.MethodDescriptor, reqDesc protoreflect.MessageDescriptor, getCfg func() *config.Config, getOutput func() string) *cobra.Command {
	parentFlagValues := make(map[string]interface{})

	parent := &cobra.Command{
		Use:     spec.Verb,
		Short:   spec.Short,
		Long:    spec.Long,
		Example: spec.Example,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Parent only handles --from-json / --request-json.
			useJSON, requestJSON, err := readJSONInput(cmd)
			if err != nil {
				return err
			}
			if !useJSON {
				return cmd.Help()
			}
			req := dynamicpb.NewMessage(reqDesc)
			if err := protojson.Unmarshal(requestJSON, req); err != nil {
				return fmt.Errorf("parsing JSON input: %w", err)
			}
			return execute(cmd, req, methodDesc, reqDesc, spec, getCfg, getOutput)
		},
	}

	// Register common (non-oneof) flags on the parent as persistent flags.
	registerFlags(parent.PersistentFlags(), reqDesc, "", parentFlagValues, spec.PositionalArgs, 0, true, nil)

	// Build variant subcommands.
	if len(spec.OneofGroups) > 0 {
		og := spec.OneofGroups[0]
		oneofDesc := reqDesc.Oneofs().ByName(protoreflect.Name(og.ProtoOneofName))

		for _, variant := range og.Variants {
			variant := variant
			variantFieldDesc := reqDesc.Fields().ByName(protoreflect.Name(variant.ProtoFieldName))
			if variantFieldDesc == nil {
				continue
			}

			varFlagValues := make(map[string]interface{})

			use := variant.CliName
			// Add positional arg to variant's Use string.
			for _, pa := range spec.PositionalArgs {
				use += " <" + pa.ProtoFieldName + ">"
			}

			vc := &cobra.Command{
				Use:   use,
				Short: spec.Short + " (" + variant.CliName + ")",
				RunE: func(cmd *cobra.Command, args []string) error {
					req := dynamicpb.NewMessage(reqDesc)

					useJSON, requestJSON, err := readJSONInput(cmd)
					if err != nil {
						return err
					}
					if useJSON {
						if err := protojson.Unmarshal(requestJSON, req); err != nil {
							return fmt.Errorf("parsing JSON input: %w", err)
						}
					}

					// Apply positional args.
					if err := applyPositionalArgs(req, reqDesc, spec.PositionalArgs, args, useJSON); err != nil {
						return err
					}

					if !useJSON {
						// Apply common flags.
						if err := ValidateAllFlags(parentFlagValues); err != nil {
							return err
						}
						if err := ValidateAllFlags(varFlagValues); err != nil {
							return err
						}
						applyFlags(cmd, req, reqDesc, parentFlagValues)

						// Apply variant-specific fields.
						if err := applyVariantFlags(cmd, req, reqDesc, oneofDesc, variantFieldDesc, varFlagValues); err != nil {
							return err
						}
					}

					return execute(cmd, req, methodDesc, reqDesc, spec, getCfg, getOutput)
				},
			}

			// Register variant-specific flags.
			if variantFieldDesc.Kind() == protoreflect.MessageKind {
				variantMsgDesc := variantFieldDesc.Message()
				varAliases := make(map[string][]string)
				registerFlags(vc.Flags(), variantMsgDesc, "", varFlagValues, nil, 0, true, varAliases)
				registerAliases(vc.Flags(), varFlagValues, varAliases)
			}

			parent.AddCommand(vc)
		}
	}

	return parent
}

// resolveMethod looks up the proto method descriptor from the global registry.
func resolveMethod(fullMethodName string) (protoreflect.MethodDescriptor, protoreflect.MessageDescriptor, error) {
	// fullMethodName is "package.Service/Method"
	parts := strings.SplitN(fullMethodName, "/", 2)
	if len(parts) != 2 {
		return nil, nil, fmt.Errorf("invalid method name %q", fullMethodName)
	}
	serviceFQN := protoreflect.FullName(parts[0])
	methodName := protoreflect.Name(parts[1])

	// Look up the service descriptor.
	desc, err := protoregistry.GlobalFiles.FindDescriptorByName(serviceFQN)
	if err != nil {
		return nil, nil, fmt.Errorf("service %q not found in registry: %w", serviceFQN, err)
	}
	serviceDesc, ok := desc.(protoreflect.ServiceDescriptor)
	if !ok {
		return nil, nil, fmt.Errorf("%q is not a service", serviceFQN)
	}

	methodDesc := serviceDesc.Methods().ByName(methodName)
	if methodDesc == nil {
		return nil, nil, fmt.Errorf("method %q not found on service %q", methodName, serviceFQN)
	}

	return methodDesc, methodDesc.Input(), nil
}

// maxFlagDepth is the maximum recursion depth for expanding message fields to flags.
const maxFlagDepth = 3

// wellKnownLeafTypes are proto messages that should be treated as leaf string fields, not recursed into.
var wellKnownLeafTypes = map[protoreflect.FullName]bool{
	"google.protobuf.Timestamp": true,
	"google.protobuf.Duration":  true,
	"google.protobuf.Struct":    true,
	"google.protobuf.Value":     true,
	"google.protobuf.Any":       true,
}

// registerFlags walks the message descriptor and registers Cobra flags for scalar fields.
// It recurses into message fields up to maxFlagDepth with prefixed flag names.
// parentRequired controls whether [REQUIRED] markers propagate: children of optional
// parent messages are not marked required even if the proto says they are.
// aliasTracker collects leaf names for alias registration; pass nil to disable.
func registerFlags(flags *pflagSet, msgDesc protoreflect.MessageDescriptor, prefix string, values map[string]interface{}, skipFields []PosArg, depth int, parentRequired bool, aliasTracker map[string][]string) {
	skipSet := make(map[string]bool)
	for _, pa := range skipFields {
		skipSet[pa.ProtoFieldName] = true
	}

	for i := 0; i < msgDesc.Fields().Len(); i++ {
		fd := msgDesc.Fields().Get(i)

		// Skip fields that are positional args.
		if skipSet[string(fd.Name())] {
			continue
		}

		// Skip oneof fields at the top level (handled as subcommands).
		// At deeper levels, promote oneof alternatives as regular flags.
		if depth == 0 && fd.ContainingOneof() != nil && !fd.HasOptionalKeyword() {
			continue
		}

		leafName := toKebabCase(string(fd.Name()))
		flagName := leafName
		if prefix != "" {
			// De-duplication: if prefix ends with the leaf name, don't repeat.
			// e.g., prefix="email", leaf="email" → "email" not "email-email"
			if prefix == leafName || strings.HasSuffix(prefix, "-"+leafName) {
				flagName = prefix
			} else {
				flagName = prefix + "-" + leafName
			}
		}

		help := "Set " + strings.ReplaceAll(leafName, "-", " ")
		// Only show [REQUIRED] if the field is required AND all parent messages are also required.
		if parentRequired && isRequiredField(fd) {
			help += " [REQUIRED]"
		}

		switch fd.Kind() {
		case protoreflect.StringKind:
			if fd.IsList() {
				val := flags.StringSlice(flagName, nil, help)
				values[flagName] = val
			} else {
				val := flags.String(flagName, "", help)
				values[flagName] = val
			}
		case protoreflect.BoolKind:
			val := flags.Bool(flagName, false, help)
			values[flagName] = val
		case protoreflect.Int32Kind, protoreflect.Sint32Kind:
			val := flags.Int32(flagName, 0, help)
			values[flagName] = val
		case protoreflect.Uint32Kind:
			val := flags.Uint32(flagName, 0, help)
			values[flagName] = val
		case protoreflect.Int64Kind, protoreflect.Sint64Kind:
			val := flags.Int64(flagName, 0, help)
			values[flagName] = val
		case protoreflect.Uint64Kind:
			val := flags.Uint64(flagName, 0, help)
			values[flagName] = val
		case protoreflect.EnumKind:
			// Present enum as string flag with valid values in help.
			enumDesc := fd.Enum()
			var valNames []string
			for j := 0; j < enumDesc.Values().Len(); j++ {
				name := string(enumDesc.Values().Get(j).Name())
				if !strings.HasSuffix(name, "UNSPECIFIED") {
					valNames = append(valNames, name)
				}
			}
			if fd.IsList() {
				help += " (values: " + strings.Join(valNames, ", ") + ")"
				val := flags.StringSlice(flagName, nil, help)
				values[flagName] = val
			} else {
				help += " (one of: " + strings.Join(valNames, ", ") + ")"
				val := flags.String(flagName, "", help)
				values[flagName] = val
			}
		case protoreflect.MessageKind:
			// Well-known types: treat as leaf string fields.
			if wellKnownLeafTypes[fd.Message().FullName()] {
				val := flags.String(flagName, "", help+" (as string)")
				values[flagName] = val
				continue
			}
			// Recurse into message fields up to maxFlagDepth.
			// Children inherit required-ness only if this message field is also required.
			if depth < maxFlagDepth {
				childRequired := parentRequired && isRequiredField(fd)
				registerFlags(flags, fd.Message(), flagName, values, nil, depth+1, childRequired, aliasTracker)
			}
			continue // don't track aliases for message containers
		}

		// Track for alias registration: leafName → list of canonical flagNames that use it.
		if aliasTracker != nil && prefix != "" && flagName != leafName {
			aliasTracker[leafName] = append(aliasTracker[leafName], flagName)
		}
	}
}

// registerAliases creates hidden alias flags for unique leaf names.
// If a leaf name (e.g., "given-name") maps to exactly one canonical flag
// (e.g., "profile-given-name"), register a hidden alias pointing to the same value.
func registerAliases(flags *pflagSet, values map[string]interface{}, aliasTracker map[string][]string) {
	for leafName, canonicals := range aliasTracker {
		if len(canonicals) != 1 {
			continue // ambiguous — skip
		}
		canonical := canonicals[0]
		if flags.Lookup(leafName) != nil {
			continue // already registered (e.g., a top-level field with the same name)
		}
		val, ok := values[canonical]
		if !ok {
			continue
		}
		// Register alias pointing to the same value.
		switch v := val.(type) {
		case *string:
			flags.StringVar(v, leafName, "", "Alias for --"+canonical)
		case *bool:
			flags.BoolVar(v, leafName, false, "Alias for --"+canonical)
		case *int32:
			flags.Int32Var(v, leafName, 0, "Alias for --"+canonical)
		case *uint32:
			flags.Uint32Var(v, leafName, 0, "Alias for --"+canonical)
		case *int64:
			flags.Int64Var(v, leafName, 0, "Alias for --"+canonical)
		case *uint64:
			flags.Uint64Var(v, leafName, 0, "Alias for --"+canonical)
		}
		// Map the alias to the same value for applyFlags lookup.
		values[leafName] = val
		// Mark alias as hidden so --help isn't cluttered.
		if f := flags.Lookup(leafName); f != nil {
			f.Hidden = true
		}
	}
}

// pflagSet is a type alias for the pflag FlagSet.
type pflagSet = pflag.FlagSet

// applyFlags reads set flag values and applies them to the dynamic message.
// It recurses into message fields to set nested sub-messages from prefixed flags.
func applyFlags(cmd *cobra.Command, msg *dynamicpb.Message, msgDesc protoreflect.MessageDescriptor, values map[string]interface{}) {
	applyFlagsWithPrefix(cmd, msg, msgDesc, "", values, 0)
}

func applyFlagsWithPrefix(cmd *cobra.Command, msg *dynamicpb.Message, msgDesc protoreflect.MessageDescriptor, prefix string, values map[string]interface{}, depth int) {
	for i := 0; i < msgDesc.Fields().Len(); i++ {
		fd := msgDesc.Fields().Get(i)

		// Skip oneof fields at top level only (subcommands); promote at deeper levels.
		if fd.ContainingOneof() != nil && !fd.HasOptionalKeyword() {
			// In applyFlags we always process nested oneofs — the depth=0 skip
			// is only in registerFlags. Here we need to process them.
		}

		leafName := toKebabCase(string(fd.Name()))
		flagName := leafName
		if prefix != "" {
			if prefix == leafName || strings.HasSuffix(prefix, "-"+leafName) {
				flagName = prefix
			} else {
				flagName = prefix + "-" + leafName
			}
		}

		if fd.Kind() == protoreflect.MessageKind && !wellKnownLeafTypes[fd.Message().FullName()] && depth < maxFlagDepth {
			// Recurse: build sub-message if any sub-flags were set.
			subMsgDesc := fd.Message()
			subMsg := dynamicpb.NewMessage(subMsgDesc)
			anySet := applyFlagsCheckSet(cmd, subMsg, subMsgDesc, flagName, values, depth+1)
			if anySet {
				msg.Set(fd, protoreflect.ValueOfMessage(subMsg))
			}
			continue
		}

		// Check if this flag or its aliases were changed.
		if !isFlagChanged(cmd, flagName, leafName) {
			continue
		}

		// Look up value by canonical name first, then alias.
		val, ok := values[flagName]
		if !ok {
			val, ok = values[leafName]
		}
		if !ok {
			continue
		}

		applyFlagValue(msg, fd, val)
	}
}

// applyFlagsCheckSet is like applyFlagsWithPrefix but returns whether any flag was set.
func applyFlagsCheckSet(cmd *cobra.Command, msg *dynamicpb.Message, msgDesc protoreflect.MessageDescriptor, prefix string, values map[string]interface{}, depth int) bool {
	anySet := false
	for i := 0; i < msgDesc.Fields().Len(); i++ {
		fd := msgDesc.Fields().Get(i)

		// Process oneof fields at all depths (nested oneofs are promoted to flags).

		leafName := toKebabCase(string(fd.Name()))
		flagName := leafName
		if prefix != "" {
			if prefix == leafName || strings.HasSuffix(prefix, "-"+leafName) {
				flagName = prefix
			} else {
				flagName = prefix + "-" + leafName
			}
		}

		if fd.Kind() == protoreflect.MessageKind && !wellKnownLeafTypes[fd.Message().FullName()] && depth < maxFlagDepth {
			subMsgDesc := fd.Message()
			subMsg := dynamicpb.NewMessage(subMsgDesc)
			if applyFlagsCheckSet(cmd, subMsg, subMsgDesc, flagName, values, depth+1) {
				msg.Set(fd, protoreflect.ValueOfMessage(subMsg))
				anySet = true
			}
			continue
		}

		if !isFlagChanged(cmd, flagName, leafName) {
			continue
		}

		val, ok := values[flagName]
		if !ok {
			val, ok = values[leafName]
		}
		if !ok {
			continue
		}

		applyFlagValue(msg, fd, val)
		anySet = true
	}
	return anySet
}

// isFlagChanged checks if a flag or its alias was explicitly set.
func isFlagChanged(cmd *cobra.Command, flagName, leafName string) bool {
	if cmd.Flags().Changed(flagName) {
		return true
	}
	if leafName != flagName && cmd.Flags().Changed(leafName) {
		return true
	}
	return false
}

// applyFlagValue sets a single field on a dynamic message from a flag pointer value.
func applyFlagValue(msg *dynamicpb.Message, fd protoreflect.FieldDescriptor, val interface{}) {
	switch fd.Kind() {
	case protoreflect.StringKind:
		if fd.IsList() {
			if ptr, ok := val.(*[]string); ok && ptr != nil {
				list := msg.Mutable(fd).List()
				for _, s := range *ptr {
					list.Append(protoreflect.ValueOfString(s))
				}
			}
		} else {
			if ptr, ok := val.(*string); ok && *ptr != "" {
				msg.Set(fd, protoreflect.ValueOfString(*ptr))
			}
		}
	case protoreflect.BoolKind:
		if ptr, ok := val.(*bool); ok {
			msg.Set(fd, protoreflect.ValueOfBool(*ptr))
		}
	case protoreflect.Int32Kind, protoreflect.Sint32Kind:
		if ptr, ok := val.(*int32); ok && *ptr != 0 {
			msg.Set(fd, protoreflect.ValueOfInt32(*ptr))
		}
	case protoreflect.Uint32Kind:
		if ptr, ok := val.(*uint32); ok && *ptr != 0 {
			msg.Set(fd, protoreflect.ValueOfUint32(*ptr))
		}
	case protoreflect.Int64Kind, protoreflect.Sint64Kind:
		if ptr, ok := val.(*int64); ok && *ptr != 0 {
			msg.Set(fd, protoreflect.ValueOfInt64(*ptr))
		}
	case protoreflect.Uint64Kind:
		if ptr, ok := val.(*uint64); ok && *ptr != 0 {
			msg.Set(fd, protoreflect.ValueOfUint64(*ptr))
		}
	case protoreflect.EnumKind:
		if fd.IsList() {
			// Repeated enum: expect *[]string.
			if ptr, ok := val.(*[]string); ok && ptr != nil && len(*ptr) > 0 {
				list := msg.Mutable(fd).List()
				for _, name := range *ptr {
					enumVal := fd.Enum().Values().ByName(protoreflect.Name(name))
					if enumVal != nil {
						list.Append(protoreflect.ValueOfEnum(enumVal.Number()))
					}
				}
			}
		} else if ptr, ok := val.(*string); ok && *ptr != "" {
			// Single enum: expect *string.
			enumVal := fd.Enum().Values().ByName(protoreflect.Name(*ptr))
			if enumVal != nil {
				msg.Set(fd, protoreflect.ValueOfEnum(enumVal.Number()))
			}
		}
	}
}

// applyVariantFlags applies variant-specific flags and sets the oneof field.
func applyVariantFlags(cmd *cobra.Command, req *dynamicpb.Message, reqDesc protoreflect.MessageDescriptor, oneofDesc protoreflect.OneofDescriptor, variantFD protoreflect.FieldDescriptor, varFlagValues map[string]interface{}) error {
	switch variantFD.Kind() {
	case protoreflect.BoolKind:
		// Scalar bool variant: just set it to true.
		req.Set(variantFD, protoreflect.ValueOfBool(true))
	case protoreflect.StringKind:
		// Scalar string variant handled via positional args in the variant's Use.
		// Already handled by applyPositionalArgs or the caller.
	case protoreflect.MessageKind:
		// Message variant: build the sub-message from flags.
		variantMsgDesc := variantFD.Message()
		variantMsg := dynamicpb.NewMessage(variantMsgDesc)
		applyFlags(cmd, variantMsg, variantMsgDesc, varFlagValues)
		req.Set(variantFD, protoreflect.ValueOfMessage(variantMsg))
	}
	_ = oneofDesc // used implicitly via variantFD
	_ = reqDesc
	return nil
}

// applyPositionalArgs maps positional arguments to proto fields.
func applyPositionalArgs(req *dynamicpb.Message, reqDesc protoreflect.MessageDescriptor, posArgs []PosArg, args []string, fromJSON bool) error {
	for i, pa := range posArgs {
		fd := reqDesc.Fields().ByName(protoreflect.Name(pa.ProtoFieldName))
		if fd == nil {
			continue
		}
		if i < len(args) {
			req.Set(fd, protoreflect.ValueOfString(args[i]))
		} else if !fromJSON && pa.Required {
			// Check if already set from JSON.
			if req.Get(fd).String() == "" {
				return fmt.Errorf("missing required argument <%s>", pa.ProtoFieldName)
			}
		}
	}
	return nil
}

// checkRequiredFlags validates that required fields have been provided.
// Fields marked with google.api.field_behavior = REQUIRED in the proto
// definitions must have non-zero values when not using JSON input.
func checkRequiredFlags(cmd *cobra.Command, reqDesc protoreflect.MessageDescriptor, spec CommandSpec) error {
	skipSet := make(map[string]bool)
	for _, pa := range spec.PositionalArgs {
		skipSet[pa.ProtoFieldName] = true
	}

	return checkRequiredFieldsRecursive(cmd, reqDesc, "", skipSet, 0, true)
}

// checkRequiredFieldsRecursive walks message fields and checks that REQUIRED ones are set.
// parentRequired controls whether enforcement propagates: children of optional parent
// messages are not checked even if they have REQUIRED annotations.
func checkRequiredFieldsRecursive(cmd *cobra.Command, msgDesc protoreflect.MessageDescriptor, prefix string, skipSet map[string]bool, depth int, parentRequired bool) error {
	if depth > maxFlagDepth {
		return nil
	}
	for i := 0; i < msgDesc.Fields().Len(); i++ {
		fd := msgDesc.Fields().Get(i)

		if skipSet[string(fd.Name())] {
			continue
		}
		if fd.ContainingOneof() != nil && !fd.HasOptionalKeyword() {
			continue
		}

		leafName := toKebabCase(string(fd.Name()))
		flagName := leafName
		if prefix != "" {
			if prefix == leafName || strings.HasSuffix(prefix, "-"+leafName) {
				flagName = prefix
			} else {
				flagName = prefix + "-" + leafName
			}
		}

		if fd.Kind() == protoreflect.MessageKind && !wellKnownLeafTypes[fd.Message().FullName()] {
			// Only recurse with required enforcement if this field itself is required.
			childRequired := parentRequired && isRequiredField(fd)
			if err := checkRequiredFieldsRecursive(cmd, fd.Message(), flagName, nil, depth+1, childRequired); err != nil {
				return err
			}
			continue
		}

		if parentRequired && isRequiredField(fd) && !isFlagChanged(cmd, flagName, leafName) {
			return fmt.Errorf("missing required flag --%s", flagName)
		}
	}
	return nil
}

// readJSONInput reads --from-json / --request-json input.
func readJSONInput(cmd *cobra.Command) (useJSON bool, data []byte, err error) {
	requestJSON, _ := cmd.Flags().GetString("request-json")
	fromJSON, _ := cmd.Flags().GetBool("from-json")

	if requestJSON != "" {
		return true, []byte(requestJSON), nil
	}
	if fromJSON {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) != 0 {
			return false, nil, fmt.Errorf("--from-json specified but no data piped to stdin")
		}
		data, err = io.ReadAll(os.Stdin)
		if err != nil {
			return false, nil, fmt.Errorf("reading stdin: %w", err)
		}
		return true, data, nil
	}
	return false, nil, nil
}

// execute is the common execution path for all commands.
func execute(cmd *cobra.Command, req *dynamicpb.Message, methodDesc protoreflect.MethodDescriptor, reqDesc protoreflect.MessageDescriptor, spec CommandSpec, getCfg func() *config.Config, getOutput func() string) error {
	if spec.Deprecated {
		fmt.Fprintln(os.Stderr, "Warning: this command is deprecated and may be removed in a future release.")
	}

	// Dry-run: print the request as JSON.
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	if dryRun {
		reqJSON, err := protojson.MarshalOptions{Indent: "  "}.Marshal(req)
		if err != nil {
			return fmt.Errorf("marshalling request: %w", err)
		}
		envelope := map[string]json.RawMessage{
			"method":  json.RawMessage(fmt.Sprintf("%q", spec.FullMethodName)),
			"request": json.RawMessage(reqJSON),
		}
		out, _ := json.MarshalIndent(envelope, "", "  ")
		fmt.Println(string(out))
		return nil
	}

	// Safety: confirm destructive operations before calling the API.
	if isDestructiveVerb(spec.Verb) {
		yesFlag, _ := cmd.Flags().GetBool("yes")
		if err := confirmDestructive(cmd, spec, req, reqDesc, yesFlag); err != nil {
			return err
		}
	}

	// Real execution: call the API.
	cfg := getCfg()
	actx, _, err := config.ActiveCtx(cfg)
	if err != nil {
		return err
	}
	tokenSource, err := auth.TokenSource(cmd.Context(), actx)
	if err != nil {
		return err
	}
	httpClient := client.New(tokenSource)
	baseURL := client.InstanceURL(actx.Instance)

	// Make the ConnectRPC call using the generic transport.
	output.StartSpinner(fmt.Sprintf("Calling %s...", spec.FullMethodName))
	respMsg, err := callConnect(cmd.Context(), httpClient, baseURL, methodDesc, req)
	output.StopSpinner()
	if err != nil {
		return err
	}

	// Output formatting.
	quietFlag, _ := cmd.Flags().GetBool("quiet")
	if quietFlag {
		return nil
	}
	fieldsFlag, _ := cmd.Flags().GetString("fields")
	return formatOutput(getOutput(), respMsg, spec, fieldsFlag)
}

// formatOutput renders the response based on the output mode (json, table, describe).
func formatOutput(mode string, resp proto.Message, spec CommandSpec, fields string) error {
	if mode == "json" {
		data, err := protojson.MarshalOptions{Indent: "  "}.Marshal(resp)
		if err != nil {
			return fmt.Errorf("marshalling response: %w", err)
		}

		// Apply client-side field filtering if --fields is specified.
		if fields != "" {
			filtered, err := filterJSONFields(data, fields)
			if err == nil {
				data = filtered
			}
		}

		_, err = os.Stdout.Write(append(data, '\n'))
		return err
	}

	// Table mode: use for list commands (explicit or default) when columns are defined.
	if len(spec.TableColumns) > 0 && (mode == "table" || (mode == "" && spec.IsListMethod)) {
		return renderTable(resp, spec)
	}

	// Describe mode: default for get/single-resource commands, or explicit --output describe.
	output.Describe(resp)
	return nil
}

// renderTable extracts fields from the response and renders them as a table.
func renderTable(resp proto.Message, spec CommandSpec) error {
	headers := make([]string, len(spec.TableColumns))
	for i, col := range spec.TableColumns {
		headers[i] = col.Header
	}

	refMsg := resp.ProtoReflect()
	if spec.IsListMethod && spec.ListFieldName != "" {
		// List method: iterate over the repeated field.
		listFD := refMsg.Descriptor().Fields().ByName(protoreflect.Name(spec.ListFieldName))
		if listFD == nil {
			output.Table(headers, nil)
			return nil
		}
		list := refMsg.Get(listFD).List()
		rows := make([][]string, list.Len())
		for i := 0; i < list.Len(); i++ {
			itemMsg := list.Get(i).Message()
			rows[i] = extractRow(itemMsg, spec.TableColumns)
		}
		output.Table(headers, rows)
		return nil
	}

	// Single result — optionally unwrap a nested field.
	targetMsg := refMsg
	if spec.ResponseUnwrapField != "" {
		fd := refMsg.Descriptor().Fields().ByName(protoreflect.Name(spec.ResponseUnwrapField))
		if fd != nil && fd.Kind() == protoreflect.MessageKind {
			targetMsg = refMsg.Get(fd).Message()
		}
	}

	row := extractRow(targetMsg, spec.TableColumns)
	output.Table(headers, [][]string{row})
	return nil
}

// extractRow extracts column values from a message using dot-separated field paths.
func extractRow(msg protoreflect.Message, cols []ColumnSpec) []string {
	row := make([]string, len(cols))
	for i, col := range cols {
		row[i] = resolveFieldPath(msg, col)
	}
	return row
}

// resolveFieldPath traverses a dot-separated path on a message to get the display value.
func resolveFieldPath(msg protoreflect.Message, col ColumnSpec) string {
	parts := strings.Split(col.FieldPath, ".")
	current := msg
	for j, part := range parts {
		fd := current.Descriptor().Fields().ByName(protoreflect.Name(part))
		if fd == nil {
			return ""
		}
		if j < len(parts)-1 {
			// Intermediate: must be a message.
			if fd.Kind() != protoreflect.MessageKind {
				return ""
			}
			current = current.Get(fd).Message()
			continue
		}
		// Leaf field.
		val := current.Get(fd)
		if col.IsTimestamp && fd.Kind() == protoreflect.MessageKind {
			tsMsg := val.Message()
			secsFd := tsMsg.Descriptor().Fields().ByName("seconds")
			nanosFd := tsMsg.Descriptor().Fields().ByName("nanos")
			if secsFd == nil {
				return ""
			}
			secs := tsMsg.Get(secsFd).Int()
			if secs == 0 {
				return ""
			}
			var nanos int32
			if nanosFd != nil {
				nanos = int32(tsMsg.Get(nanosFd).Int())
			}
			return time.Unix(secs, int64(nanos)).UTC().Format(time.RFC3339)
		}
		if col.IsEnum && fd.Kind() == protoreflect.EnumKind {
			enumVal := fd.Enum().Values().ByNumber(val.Enum())
			if enumVal != nil {
				return string(enumVal.Name())
			}
		}
		return fmt.Sprint(val.Interface())
	}
	return ""
}

// isRequiredField checks if a proto field has the REQUIRED field_behavior annotation.
// It reads the google.api.field_behavior extension (field number 1052) from the field options.
// REQUIRED is enum value 2 in google.api.FieldBehavior.
func isRequiredField(fd protoreflect.FieldDescriptor) bool {
	opts, ok := fd.Options().(*descriptorpb.FieldOptions)
	if !ok || opts == nil {
		return false
	}

	// Use proto reflection to read the field_behavior extension.
	// Extension field number 1052 = google.api.field_behavior (repeated enum).
	found := false
	optsMsg := opts.ProtoReflect()
	optsMsg.Range(func(extFD protoreflect.FieldDescriptor, val protoreflect.Value) bool {
		if extFD.Number() == 1052 && extFD.IsList() {
			list := val.List()
			for i := 0; i < list.Len(); i++ {
				// REQUIRED = 2 in google.api.FieldBehavior
				if list.Get(i).Enum() == 2 {
					found = true
					return false
				}
			}
		}
		return true
	})

	return found
}

// toKebabCase converts snake_case to kebab-case.
func toKebabCase(s string) string {
	return strings.ReplaceAll(s, "_", "-")
}

// filterJSONFields filters a JSON object to only include the named fields.
// fields is a comma-separated list of top-level field names.
func filterJSONFields(data []byte, fields string) ([]byte, error) {
	var full map[string]json.RawMessage
	if err := json.Unmarshal(data, &full); err != nil {
		return nil, err
	}

	wanted := make(map[string]bool)
	for _, f := range strings.Split(fields, ",") {
		f = strings.TrimSpace(f)
		if f != "" {
			wanted[f] = true
		}
	}

	filtered := make(map[string]json.RawMessage)
	for k, v := range full {
		if wanted[k] {
			filtered[k] = v
		}
	}

	return json.MarshalIndent(filtered, "", "  ")
}
