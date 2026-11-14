package resources

import (
	"encoding/json"
	"sort"
	"strings"

	"github.com/zitadel/zitadel/internal/api/scim/resources/filter"
	"github.com/zitadel/zitadel/internal/api/scim/resources/patch"
	"github.com/zitadel/zitadel/internal/api/scim/schemas"
	"github.com/zitadel/zitadel/internal/api/scim/serrors"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func groupMutations(operations patch.OperationCollection) ([]command.GroupMutation, error) {
	result := make([]command.GroupMutation, 0, len(operations))
	for _, op := range operations {
		mutations, err := parseGroupOperation(op)
		if err != nil {
			return nil, err
		}
		result = append(result, mutations...)
	}
	return result, nil
}

func parseGroupOperation(op *patch.Operation) ([]command.GroupMutation, error) {
	if op == nil {
		return nil, invalidGroupPatch("Missing operation")
	}
	kind := strings.ToLower(string(op.Operation))
	if kind != "add" && kind != "remove" && kind != "replace" {
		return nil, invalidGroupPatch("Unsupported operation")
	}
	if op.Path.IsZero() {
		return groupAttributeMutations(op, kind)
	}
	segments, err := op.Path.Segments(schemas.IdGroup)
	if err != nil {
		return nil, err
	}
	if len(segments) != 1 {
		return nil, invalidGroupPath()
	}
	var mutation command.GroupMutation
	switch strings.ToLower(segments[0]) {
	case "id", "meta", "schemas":
		// RFC 7644: read-only attributes in a request are ignored.
		return nil, nil
	case "displayname":
		if kind == "remove" || op.Path.ValuePath != nil {
			return nil, invalidGroupPatch("displayName is required")
		}
		var name string
		if err := json.Unmarshal(op.Value, &name); err != nil || strings.TrimSpace(name) == "" {
			return nil, invalidGroupPatch("displayName must be a non-empty string")
		}
		mutation.Name = &name
	case "externalid":
		if op.Path.ValuePath != nil {
			return nil, invalidGroupPath()
		}
		var value string
		if kind != "remove" {
			if err := json.Unmarshal(op.Value, &value); err != nil {
				return nil, invalidGroupPatch("externalId must be a string")
			}
		}
		mutation.ExternalID = &command.GroupExternalID{Value: value}
	case "members":
		mutation, err = groupMembershipMutation(op, kind)
		if err != nil {
			return nil, err
		}
	default:
		return nil, invalidGroupPath()
	}
	return []command.GroupMutation{mutation}, nil
}

func groupAttributeMutations(op *patch.Operation, kind string) ([]command.GroupMutation, error) {
	if kind == "remove" {
		return nil, invalidGroupPatch("Remove requires a path")
	}
	var values map[string]json.RawMessage
	if err := json.Unmarshal(op.Value, &values); err != nil || len(values) == 0 {
		return nil, invalidGroupPatch("Expected an attribute object")
	}
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	operations := make(patch.OperationCollection, 0, len(keys))
	for _, key := range keys {
		path, err := filter.ParsePath(key)
		if err != nil {
			return nil, err
		}
		if path.IsZero() {
			return nil, invalidGroupPath()
		}
		operations = append(operations, &patch.Operation{Operation: op.Operation, Path: path, Value: values[key]})
	}
	return groupMutations(operations)
}

func groupMembershipMutation(op *patch.Operation, kind string) (command.GroupMutation, error) {
	mutation := command.GroupMutation{}
	if op.Path.ValuePath != nil {
		if kind != "remove" {
			return mutation, invalidGroupPath()
		}
		id, err := groupMemberFilterID(op.Path.ValuePath)
		if err != nil {
			return mutation, err
		}
		mutation.RemoveMembers = []string{id}
		return mutation, nil
	}
	var members []*ScimGroupMember
	// Entra sends a value array on remove; an absent value removes all members.
	if len(op.Value) > 0 {
		if err := json.Unmarshal(op.Value, &members); err != nil || members == nil {
			return mutation, invalidGroupPatch("members must be an array")
		}
	} else if kind != "remove" {
		return mutation, invalidGroupPatch("members value is required")
	}
	ids, err := groupMemberIDs(members)
	if err != nil {
		return mutation, err
	}
	switch kind {
	case "add":
		mutation.AddMembers = ids
	case "replace":
		mutation.Members = &ids
	case "remove":
		if len(op.Value) == 0 {
			mutation.Members = &ids
		} else {
			mutation.RemoveMembers = ids
		}
	}
	return mutation, nil
}

func groupMemberFilterID(path *filter.ValuePathWithSubAttr) (string, error) {
	expression := path.ValuePath.ValFilter
	atom := expression.Left.Left
	if expression.Right != nil || expression.Left.Right != nil || atom.AttrExp == nil || atom.AttrExp.BinaryCondition == nil {
		return "", invalidGroupPath()
	}
	condition := atom.AttrExp.BinaryCondition
	if !strings.EqualFold(condition.Left.AttrName, "value") || condition.Left.SubAttr != nil || condition.Left.UrnAttributePrefix != nil || !condition.Operator.Equal || condition.Right.StringValue == nil || *condition.Right.StringValue == "" {
		return "", invalidGroupPath()
	}
	return *condition.Right.StringValue, nil
}

func groupMemberIDs(members []*ScimGroupMember) ([]string, error) {
	if err := validateMemberTypes(members); err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(members))
	for _, member := range members {
		ids = append(ids, member.Value)
	}
	return ids, nil
}

func invalidGroupPatch(message string) error {
	return serrors.ThrowInvalidValue(zerrors.ThrowInvalidArgument(nil, "SCIM-GRP8v", message))
}

func invalidGroupPath() error {
	return serrors.ThrowInvalidPath(zerrors.ThrowInvalidArgument(nil, "SCIM-GRP9p", "Supported group paths: displayName, externalId, members, members[value eq \"id\"]"))
}
