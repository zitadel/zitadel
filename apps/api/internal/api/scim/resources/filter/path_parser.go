package filter

import (
	"encoding/json"
	"fmt"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/scim/schemas"
	"github.com/zitadel/zitadel/internal/api/scim/serrors"
)

// Path The root ast node for the path grammar
// PATH = attrPath / valuePath [subAttr]
type Path struct {
	ValuePath *ValuePathWithSubAttr `parser:"@@ |"`
	AttrPath  *AttrPath             `parser:"@@"`
}

type ValuePathWithSubAttr struct {
	ValuePath ValuePath `parser:"@@"`
	SubAttr   *string   `parser:"('.' @AttrName)?"`
}

var scimPathParser = buildParser[Path]()

func ParsePath(path string) (*Path, error) {
	if path == "" {
		return nil, nil
	}

	if len(path) > maxInputLength {
		logging.WithFields("len", len(path)).Infof("scim: path exceeds maximum allowed length: %d", maxInputLength)
		return nil, serrors.ThrowInvalidFilter(fmt.Errorf("path exceeds maximum allowed length: %d", maxInputLength))
	}

	parsedPath, err := scimPathParser.ParseString("", path)
	if err != nil {
		logging.WithError(err).Info("scim: failed to parse path")
		return nil, serrors.ThrowInvalidFilter(err)
	}

	return parsedPath, nil
}

func (p *Path) UnmarshalJSON(data []byte) error {
	var rawPath string
	if err := json.Unmarshal(data, &rawPath); err != nil {
		return err
	}

	if rawPath == "" {
		return nil
	}

	parsedPath, err := ParsePath(rawPath)
	if err != nil {
		return err
	}

	*p = *parsedPath
	return nil
}

func (p *Path) String() string {
	if p.ValuePath != nil {
		return p.ValuePath.String()
	}

	return p.AttrPath.String()
}

func (p *Path) IsZero() bool {
	return p == nil || *p == Path{}
}

func (p *Path) Segments(schema schemas.ScimSchemaType) ([]string, error) {
	if p.ValuePath != nil {
		return p.ValuePath.Segments(schema)
	}

	if err := p.AttrPath.validateSchema(schema); err != nil {
		return nil, err
	}
	return p.AttrPath.Segments(), nil
}

func (v *ValuePathWithSubAttr) String() string {
	if v.SubAttr != nil {
		return v.ValuePath.String() + "." + *v.SubAttr
	}

	return v.ValuePath.String()
}

func (v *ValuePathWithSubAttr) Segments(schema schemas.ScimSchemaType) ([]string, error) {
	if err := v.ValuePath.AttrPath.validateSchema(schema); err != nil {
		return nil, err
	}

	segments := v.ValuePath.AttrPath.Segments()
	if v.SubAttr != nil {
		segments = append(segments, *v.SubAttr)
	}

	return segments, nil
}
