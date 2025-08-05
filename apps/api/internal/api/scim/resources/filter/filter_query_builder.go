package filter

import (
	"context"
	"strings"
	"time"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/scim/schemas"
	"github.com/zitadel/zitadel/internal/api/scim/serrors"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
)

// FieldPathMapping maps lowercase json field names of the resource to the matching column in the projection
type FieldPathMapping map[string]*QueryFieldInfo

// queryBuilder builds a query for a filter based on the visitor pattern
type queryBuilder struct {
	ctx              context.Context
	schema           schemas.ScimSchemaType
	fieldPathMapping FieldPathMapping

	// attrPathPrefixes prefixes of attributes that
	// should also take into account when resolving an attr path to a column.
	// This is used for "a[b eq 10]" expressions, when resolving b, a would be the prefix.
	attrPathPrefixStack []*AttrPath
}

type MappedQueryBuilderFunc func(ctx context.Context, compareValue *CompValue, op *CompareOp) (query.SearchQuery, error)

type QueryFieldInfo struct {
	Column           query.Column
	FieldType        FieldType
	CaseInsensitive  bool
	BuildMappedQuery MappedQueryBuilderFunc
}

type FieldType int

const (
	FieldTypeCustom FieldType = iota
	FieldTypeString
	FieldTypeNumber
	FieldTypeBoolean
	FieldTypeTimestamp
)

func (m FieldPathMapping) Resolve(path string) (*QueryFieldInfo, error) {
	info, ok := m[strings.ToLower(path)]
	if !ok {
		logging.WithFields("fieldPath", path).Info("scim filter: Invalid filter expression: unknown or unsupported field")
		return nil, serrors.ThrowInvalidFilter(zerrors.ThrowInvalidArgumentf(nil, "SCIM-FF433", "Invalid filter expression: unknown or unsupported field %s", path))
	}

	return info, nil
}

func (f *Filter) BuildQuery(ctx context.Context, schema schemas.ScimSchemaType, fieldPathColumnMapping FieldPathMapping) (query.SearchQuery, error) {
	builder := &queryBuilder{
		ctx:              ctx,
		schema:           schema,
		fieldPathMapping: fieldPathColumnMapping,
	}
	return builder.visitSegment(&f.Root)
}

func (b *queryBuilder) pushAttrPath(path *AttrPath) {
	b.attrPathPrefixStack = append(b.attrPathPrefixStack, path)
}

func (b *queryBuilder) popAttrPath() {
	b.attrPathPrefixStack = b.attrPathPrefixStack[:len(b.attrPathPrefixStack)-1]
}

func (b *queryBuilder) visitSegment(s *Segment) (query.SearchQuery, error) {
	return b.visitOr(&s.OrExp)
}

func (b *queryBuilder) visitOr(or *OrLogExp) (query.SearchQuery, error) {
	left, err := b.visitAnd(&or.Left)
	if err != nil {
		return nil, err
	}

	if or.Right == nil {
		return left, nil
	}

	right, err := b.visitOr(or.Right)
	if err != nil {
		return nil, err
	}

	// flatten nested or queries
	if rightOr, ok := right.(*query.OrQuery); ok {
		rightOr.Prepend(left)
		return rightOr, nil
	}

	return query.NewOrQuery(left, right)
}

func (b *queryBuilder) visitAnd(and *AndLogExp) (query.SearchQuery, error) {
	left, err := b.visitAtom(&and.Left)
	if err != nil {
		return nil, err
	}

	if and.Right == nil {
		return left, nil
	}

	right, err := b.visitAnd(and.Right)
	if err != nil {
		return nil, err
	}

	// flatten nested and queries
	if rightAnd, ok := right.(*query.AndQuery); ok {
		rightAnd.Prepend(left)
		return rightAnd, nil
	}

	return query.NewAndQuery(left, right)
}

func (b *queryBuilder) visitAtom(atom *ValueAtom) (query.SearchQuery, error) {
	switch {
	case atom.SubFilter != nil:
		return b.visitSegment(atom.SubFilter)
	case atom.Negation != nil:
		f, err := b.visitSegment(atom.Negation)
		if err != nil {
			return nil, err
		}

		return query.NewNotQuery(f)
	case atom.ValuePath != nil:
		return b.visitValuePath(atom.ValuePath)
	case atom.AttrExp != nil:
		return b.visitAttrExp(atom.AttrExp)
	}

	return nil, serrors.ThrowInvalidFilter(zerrors.ThrowInvalidArgument(nil, "SCIM-FF412", "Invalid filter expression"))
}

func (b *queryBuilder) visitValuePath(path *ValuePath) (query.SearchQuery, error) {
	b.pushAttrPath(&path.AttrPath)
	defer b.popAttrPath()
	return b.visitOr(&path.ValFilter)
}

func (b *queryBuilder) visitAttrExp(exp *AttrExp) (query.SearchQuery, error) {
	switch {
	case exp.UnaryCondition != nil:
		return b.visitUnaryCondition(exp.UnaryCondition)
	case exp.BinaryCondition != nil:
		return b.visitBinaryCondition(exp.BinaryCondition)
	}

	return nil, serrors.ThrowInvalidFilter(zerrors.ThrowInvalidArgument(nil, "SCIM-FF413", "Invalid filter expression"))
}

func (b *queryBuilder) visitUnaryCondition(condition *UnaryCondition) (query.SearchQuery, error) {
	// only supported unary operator is present
	if !condition.Operator.Present {
		return nil, serrors.ThrowInvalidFilter(zerrors.ThrowInvalidArgument(nil, "SCIM-FF419", "Unknown unary filter operator"))
	}

	field, err := b.visitAttrPath(&condition.Left)
	if err != nil {
		return nil, err
	}

	if field.FieldType == FieldTypeCustom {
		return nil, serrors.ThrowInvalidFilter(zerrors.ThrowInvalidArgument(nil, "SCIM-FXX49", "Unsupported attribute for unary filter operator"))
	}

	return query.NewNotNullQuery(field.Column)
}

func (b *queryBuilder) visitBinaryCondition(condition *BinaryCondition) (query.SearchQuery, error) {
	left, err := b.visitAttrPath(&condition.Left)
	if err != nil {
		return nil, err
	}

	if condition.Operator.Equal && condition.Right.Null {
		return query.NewIsNullQuery(left.Column)
	}

	if condition.Operator.NotEqual && condition.Right.Null {
		return query.NewNotNullQuery(left.Column)
	}

	switch left.FieldType {
	case FieldTypeCustom:
		return left.BuildMappedQuery(b.ctx, &condition.Right, &condition.Operator)
	case FieldTypeTimestamp:
		return b.buildTimestampQuery(left, condition.Right, &condition.Operator)
	case FieldTypeString:
		return b.buildTextQuery(left, condition.Right, &condition.Operator)
	case FieldTypeNumber:
		return b.buildNumberQuery(left, condition.Right, &condition.Operator)
	case FieldTypeBoolean:
		return b.buildBooleanQuery(left, condition.Right, &condition.Operator)
	}

	logging.WithFields("fieldType", left.FieldType).Error("Unknown field type")
	return nil, serrors.ThrowInvalidFilter(zerrors.ThrowInvalidArgument(nil, "SCIM-FF417", "Unknown filter expression field type"))
}

func (b *queryBuilder) buildTimestampQuery(left *QueryFieldInfo, right CompValue, op *CompareOp) (query.SearchQuery, error) {
	if right.StringValue == nil {
		return nil, serrors.ThrowInvalidFilter(zerrors.ThrowInvalidArgument(nil, "SCIM-FF451", "Invalid filter expression: the compare value for a timestamp has to be a RFC3339 string"))
	}

	timestamp, err := time.Parse(time.RFC3339, *right.StringValue)
	if err != nil {
		return nil, serrors.ThrowInvalidFilter(zerrors.ThrowInvalidArgument(err, "SCIM-FF421", "Invalid filter expression: the compare value for a timestamp has to be a RFC3339 string"))
	}

	var comp query.TimestampComparison
	switch {
	case op.Equal:
		comp = query.TimestampEquals
	case op.GreaterThan:
		comp = query.TimestampGreater
	case op.GreaterThanOrEqual:
		comp = query.TimestampGreaterOrEquals
	case op.LessThan:
		comp = query.TimestampLess
	case op.LessThanOrEqual:
		comp = query.TimestampLessOrEquals
	default:
		return nil, serrors.ThrowInvalidFilter(zerrors.ThrowInvalidArgument(nil, "SCIM-FF422", "Invalid filter expression: unsupported comparison operator for timestamp fields"))
	}

	return query.NewTimestampQuery(left.Column, timestamp, comp)
}

func (b *queryBuilder) buildNumberQuery(left *QueryFieldInfo, right CompValue, op *CompareOp) (query.SearchQuery, error) {
	if right.Int == nil && right.Float == nil {
		return nil, serrors.ThrowInvalidFilter(zerrors.ThrowInvalidArgument(nil, "SCIM-FF423", "Invalid filter expression: unsupported comparison value for numeric fields"))
	}

	var comp query.NumberComparison
	switch {
	case op.Equal:
		comp = query.NumberEquals
	case op.NotEqual:
		comp = query.NumberNotEquals
	case op.GreaterThan:
		comp = query.NumberGreater
	case op.GreaterThanOrEqual:
		comp = query.NumberGreaterOrEqual
	case op.LessThan:
		comp = query.NumberLess
	case op.LessThanOrEqual:
		comp = query.NumberLessOrEqual
	default:
		return nil, serrors.ThrowInvalidFilter(zerrors.ThrowInvalidArgument(nil, "SCIM-FF424", "Invalid filter expression: unsupported comparison operator for number fields"))
	}

	var value interface{}
	if right.Int != nil {
		value = *right.Int
	} else {
		value = *right.Float
	}
	return query.NewNumberQuery(left.Column, value, comp)
}

func (b *queryBuilder) buildBooleanQuery(field *QueryFieldInfo, right CompValue, op *CompareOp) (query.SearchQuery, error) {
	if !right.BooleanTrue && !right.BooleanFalse {
		return nil, serrors.ThrowInvalidFilter(zerrors.ThrowInvalidArgument(nil, "SCIM-FF428", "Invalid filter expression: unsupported comparison value for boolean field"))
	}

	if !op.Equal && !op.NotEqual {
		return nil, serrors.ThrowInvalidFilter(zerrors.ThrowInvalidArgument(nil, "SCIM-FF427", "Invalid filter expression: unsupported comparison operator for boolean field"))
	}

	return query.NewBoolQuery(field.Column, (op.Equal && right.BooleanTrue) || (op.NotEqual && right.BooleanFalse))
}

func (b *queryBuilder) buildTextQuery(field *QueryFieldInfo, right CompValue, op *CompareOp) (query.SearchQuery, error) {
	if right.StringValue == nil {
		return nil, serrors.ThrowInvalidFilter(zerrors.ThrowInvalidArgument(nil, "SCIM-FF429", "Invalid filter expression: unsupported comparison value for text field"))
	}

	var comp query.TextComparison
	if field.CaseInsensitive {
		switch {
		case op.Equal:
			comp = query.TextEqualsIgnoreCase
		case op.NotEqual:
			comp = query.TextNotEqualsIgnoreCase
		case op.Contains:
			comp = query.TextContainsIgnoreCase
		case op.StartsWith:
			comp = query.TextStartsWithIgnoreCase
		case op.EndsWith:
			comp = query.TextEndsWithIgnoreCase
		default:
			return nil, serrors.ThrowInvalidFilter(zerrors.ThrowInvalidArgument(nil, "SCIM-FF529", "Invalid filter expression: unsupported comparison operator for text fields"))
		}
	} else {
		switch {
		case op.Equal:
			comp = query.TextEquals
		case op.NotEqual:
			comp = query.TextNotEquals
		case op.Contains:
			comp = query.TextContains
		case op.StartsWith:
			comp = query.TextStartsWith
		case op.EndsWith:
			comp = query.TextEndsWith
		default:
			return nil, serrors.ThrowInvalidFilter(zerrors.ThrowInvalidArgument(nil, "SCIM-FF425", "Invalid filter expression: unsupported comparison operator for text fields"))
		}
	}

	return query.NewTextQuery(field.Column, *right.StringValue, comp)
}

func (b *queryBuilder) visitAttrPath(attrPath *AttrPath) (*QueryFieldInfo, error) {
	b.pushAttrPath(attrPath)
	defer b.popAttrPath()

	field, err := b.reduceAttrPaths(b.attrPathPrefixStack)
	if err != nil {
		return nil, err
	}

	return b.fieldPathMapping.Resolve(field)
}

// reduceAttrPaths reduces a slice of AttrPath
// to a simple urn + fieldPath combination.
// The urn is ensured to be unique across all segments and either to be empty or to match the schema of the builder.
// The resulting fieldPath is in the form of a.b.c with a minimum of one path segment.
func (b *queryBuilder) reduceAttrPaths(attrPaths []*AttrPath) (fieldPath string, err error) {
	if len(attrPaths) == 0 {
		err = serrors.ThrowInvalidFilter(zerrors.ThrowInvalidArgument(nil, "SCIM-FF431", "Invalid filter expression: unknown urn attribute prefix"))
		return fieldPath, err
	}

	sb := strings.Builder{}

	for _, p := range attrPaths {
		if err = p.validateSchema(b.schema); err != nil {
			return
		}

		sb.WriteString(p.FieldPath())
		sb.WriteRune('.')
	}

	fieldPath = sb.String()
	fieldPath = strings.TrimRight(fieldPath, ".") // trim very last '.'
	return fieldPath, err
}
