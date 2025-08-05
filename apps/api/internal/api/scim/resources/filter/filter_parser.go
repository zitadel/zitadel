package filter

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/scim/schemas"
	"github.com/zitadel/zitadel/internal/api/scim/serrors"
	"github.com/zitadel/zitadel/internal/zerrors"
)

// Filter The scim v2 filter
// Separation between FilterSegment and Filter is required
// due to the UnmarshalText method, which is used by the schema parser
// as well as the participle parser but should do different things here.
type Filter struct {
	Root Segment
}

// Segment The root ast node for the filter grammar
// according to the filter ABNF of https://datatracker.ietf.org/doc/html/rfc7644#section-3.4.2.2
// FILTER    = attrExp / logExp / valuePath / *1"not" "(" FILTER ")"
// to reduce lookahead needs and reduce stack depth of the parser,
// always match log expressions with optional operators
type Segment struct {
	OrExp OrLogExp `parser:"@@"`
}

// OrLogExp The logical expression according to the filter ABNF
// separated in OrLogExp and AndLogExp to simplify parser stack depth and precedence
// logExp    = FILTER SP ("and" / "or") SP FILTER
type OrLogExp struct {
	Left  AndLogExp `parser:"@@"`
	Right *OrLogExp `parser:"(Whitespace 'or' Whitespace @@)?"`
}

type AndLogExp struct {
	Left  ValueAtom  `parser:"@@"`
	Right *AndLogExp `parser:"(Whitespace 'and' Whitespace @@)?"`
}

type ValueAtom struct {
	SubFilter *Segment   `parser:"'(' @@ ')' |"`
	Negation  *Segment   `parser:"'not' '(' @@ ')' |"`
	ValuePath *ValuePath `parser:"@@ |"`
	AttrExp   *AttrExp   `parser:"@@"`
}

// ValuePath The value path according to the filter ABNF
// valuePath = attrPath "[" valFilter "]"
// instead of a separate valFilter entity the LogExp
// is used to simplify parsing.
type ValuePath struct {
	AttrPath  AttrPath `parser:"@@"`
	ValFilter OrLogExp `parser:"'[' @@ ']'"`
}

// AttrExp The attribute expression according to the filter ABNF
// attrExp   = (attrPath SP "pr") / (attrPath SP compareOp SP compValue)
type AttrExp struct {
	UnaryCondition  *UnaryCondition  `parser:"@@ |"`
	BinaryCondition *BinaryCondition `parser:"@@"`
}

type UnaryCondition struct {
	Left     AttrPath               `parser:"@@ Whitespace"`
	Operator UnaryConditionOperator `parser:"@@"`
}

type UnaryConditionOperator struct {
	Present bool `parser:"@'pr'"`
}

type BinaryCondition struct {
	Left     AttrPath  `parser:"@@ Whitespace"`
	Operator CompareOp `parser:"@@ Whitespace"`
	Right    CompValue `parser:"@@"`
}

// CompareOp according to the scim filter ABNF
// compareOp = "eq" / "ne" / "co" /
// "sw" / "ew" /
// "gt" / "lt" /
// "ge" / "le"
type CompareOp struct {
	Equal              bool `parser:"@'eq' |"`
	NotEqual           bool `parser:"@'ne' |"`
	Contains           bool `parser:"@'co' |"`
	StartsWith         bool `parser:"@'sw' |"`
	EndsWith           bool `parser:"@'ew' |"`
	GreaterThan        bool `parser:"@'gt' |"`
	GreaterThanOrEqual bool `parser:"@'ge' |"`
	LessThan           bool `parser:"@'lt' |"`
	LessThanOrEqual    bool `parser:"@'le'"`
}

// CompValue the compare value according to the scim filter ABNF
// compValue = false / null / true / number / string
type CompValue struct {
	Null         bool     `parser:"@'null' |"`
	BooleanTrue  bool     `parser:"@'true' |"`
	BooleanFalse bool     `parser:"@'false' |"`
	Int          *int     `parser:"@Int |"`
	Float        *float64 `parser:"@Float |"`
	StringValue  *string  `parser:"@String"`
}

// AttrPath the attribute path according to the scim filter ABNF
// [URI ":"] AttrName *1subAttr
type AttrPath struct {
	UrnAttributePrefix *string `parser:"(@UrnAttributePrefix)?"`
	AttrName           string  `parser:"@AttrName"`
	SubAttr            *string `parser:"('.' @AttrName)?"`
}

const (
	maxInputLength = 1000
)

var (
	scimFilterLexer = lexer.MustSimple([]lexer.SimpleRule{
		// simplified version of RFC8141, last part isn't matched as in scim this is the attribute name
		// urn is additionally verified after parsing, use a more relaxed matching here
		{Name: "UrnAttributePrefix", Pattern: `urn:([\w()+,.=@;$_!*'%/?#-]+:)+`},
		{Name: "Float", Pattern: `[-+]?\d*\.\d+`},
		{Name: "Int", Pattern: `[-+]?\d+`},
		{Name: "Parenthesis", Pattern: `\(|\)|\[|\]`},
		{Name: "Punctuation", Pattern: `\.`},
		{Name: "String", Pattern: `"(\\"|[^"])*"`},
		{Name: "AttrName", Pattern: `[a-zA-Z][\w-]*`}, // AttrName according to the scim ABNF
		{Name: "Whitespace", Pattern: `[ \t\n\r]+`},
	})

	scimFilterParser = buildParser[Segment]()
)

func buildParser[T any]() *participle.Parser[T] {
	return participle.MustBuild[T](
		participle.Lexer(scimFilterLexer),
		participle.Unquote("String"),
		// Keyword literals are matched case-insensitive (according to https://datatracker.ietf.org/doc/html/rfc7644#section-3.4.2.2)
		// Keywords are a subset of AttrName
		participle.CaseInsensitive("AttrName"),
		participle.Elide("Whitespace"),
		participle.UseLookahead(participle.MaxLookahead),
	)
}

func (f *Filter) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		*f = Filter{}
		return nil
	}

	parsedFilter, err := ParseFilter(string(text))
	if err != nil {
		return err
	}

	*f = *parsedFilter
	return nil
}

func (f *Filter) UnmarshalJSON(data []byte) error {
	var rawFilter string
	if err := json.Unmarshal(data, &rawFilter); err != nil {
		return err
	}

	return f.UnmarshalText([]byte(rawFilter))
}

func (f *Filter) IsZero() bool {
	return f == nil || *f == Filter{}
}

func ParseFilter(filter string) (*Filter, error) {
	if filter == "" {
		return nil, nil
	}

	if len(filter) > maxInputLength {
		logging.WithFields("len", len(filter)).Infof("scim: filter exceeds maximum allowed length: %d", maxInputLength)
		return nil, serrors.ThrowInvalidFilter(zerrors.ThrowInvalidArgumentf(nil, "SCIM-filt13", "filter exceeds maximum allowed length: %d", maxInputLength))
	}

	parsedFilter, err := scimFilterParser.ParseString("", filter)
	if err != nil {
		logging.WithError(err).Info("scim: failed to parse filter")
		return nil, serrors.ThrowInvalidFilter(zerrors.ThrowInvalidArgument(err, "SCIM-filt14", "failed to parse filter"))
	}

	return &Filter{Root: *parsedFilter}, nil
}

func (f *Filter) String() string {
	return f.Root.String()
}

func (f *Segment) String() string {
	return f.OrExp.String()
}

func (o *OrLogExp) String() string {
	if o.Right == nil {
		return o.Left.String()
	}

	return "((" + o.Left.String() + ") or (" + o.Right.String() + "))"
}

func (a *AndLogExp) String() string {
	if a.Right == nil {
		return a.Left.String()
	}

	return "((" + a.Left.String() + ") and (" + a.Right.String() + "))"
}

func (a *ValueAtom) String() string {
	switch {
	case a.SubFilter != nil:
		return "(" + a.SubFilter.String() + ")"
	case a.Negation != nil:
		return "not(" + a.Negation.String() + ")"
	case a.ValuePath != nil:
		return a.ValuePath.String()
	}

	return a.AttrExp.String()
}

func (v *ValuePath) String() string {
	return v.AttrPath.String() + "[" + v.ValFilter.String() + "]"
}

func (a *AttrExp) String() string {
	if a.UnaryCondition != nil {
		return a.UnaryCondition.String()
	}

	return a.BinaryCondition.String()
}

func (u *UnaryCondition) String() string {
	return u.Left.String() + " " + u.Operator.String()
}

func (u *UnaryConditionOperator) String() string {
	return "pr"
}

func (b *BinaryCondition) String() string {
	return b.Left.String() + " " + b.Operator.String() + " " + b.Right.String()
}

func (c *CompareOp) String() string {
	switch {
	case c.Equal:
		return "eq"
	case c.NotEqual:
		return "ne"
	case c.Contains:
		return "co"
	case c.StartsWith:
		return "sw"
	case c.EndsWith:
		return "ew"
	case c.GreaterThan:
		return "gt"
	case c.GreaterThanOrEqual:
		return "ge"
	case c.LessThan:
		return "lt"
	case c.LessThanOrEqual:
		return "le"
	}

	return "<unknown CompareOp>"
}

func (c *CompValue) String() string {
	switch {
	case c.Null:
		return "null"
	case c.BooleanTrue:
		return "true"
	case c.BooleanFalse:
		return "false"
	case c.Int != nil:
		return strconv.Itoa(*c.Int)
	case c.Float != nil:
		return strconv.FormatFloat(*c.Float, 'f', -1, 64)
	case c.StringValue != nil:
		return "\"" + *c.StringValue + "\""
	}
	return "<unknown CompValue>"
}

func (a *AttrPath) String() string {
	var s = ""
	if a.UrnAttributePrefix != nil {
		s += *a.UrnAttributePrefix
	}

	s += a.AttrName

	if a.SubAttr != nil {
		s += "." + *a.SubAttr
	}

	return s
}

func (a *AttrPath) validateSchema(expectedSchema schemas.ScimSchemaType) error {
	if a.UrnAttributePrefix == nil || *a.UrnAttributePrefix == string(expectedSchema)+":" {
		return nil
	}

	logging.WithFields("urnPrefix", *a.UrnAttributePrefix).Info("scim filter: Invalid filter expression: unknown urn attribute prefix")
	return serrors.ThrowInvalidFilter(zerrors.ThrowInvalidArgument(nil, "SCIM-FF431", "Invalid filter expression: unknown urn attribute prefix"))
}

func (a *AttrPath) Segments() []string {
	// user lower, since attribute names in scim are always case-insensitive
	if a.SubAttr != nil {
		return []string{strings.ToLower(a.AttrName), strings.ToLower(*a.SubAttr)}
	}

	return []string{strings.ToLower(a.AttrName)}
}

func (a *AttrPath) FieldPath() string {
	return strings.Join(a.Segments(), ".")
}
