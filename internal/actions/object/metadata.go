package object

import (
	"context"
	"encoding/json"
	"time"

	"github.com/dop251/goja"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/actions"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
)

func OrgMetadataListFromQuery(c *actions.FieldConfig, orgMetadata *query.OrgMetadataList) goja.Value {
	result := &metadataList{
		Count:     orgMetadata.Count,
		Sequence:  orgMetadata.Sequence,
		Timestamp: orgMetadata.LastRun,
		Metadata:  make([]*metadata, len(orgMetadata.Metadata)),
	}

	for i, md := range orgMetadata.Metadata {
		result.Metadata[i] = &metadata{
			CreationDate:  md.CreationDate,
			ChangeDate:    md.ChangeDate,
			ResourceOwner: md.ResourceOwner,
			Sequence:      md.Sequence,
			Key:           md.Key,
			Value:         metadataByteArrayToValue(md.Value, c.Runtime),
		}
	}

	return c.Runtime.ToValue(result)
}

func UserMetadataListFromQuery(c *actions.FieldConfig, metadata *query.UserMetadataList) goja.Value {
	result := &userMetadataList{
		Count:     metadata.Count,
		Sequence:  metadata.Sequence,
		Timestamp: metadata.LastRun,
		Metadata:  make([]*userMetadata, len(metadata.Metadata)),
	}

	for i, md := range metadata.Metadata {
		result.Metadata[i] = &userMetadata{
			CreationDate:  md.CreationDate,
			ChangeDate:    md.ChangeDate,
			ResourceOwner: md.ResourceOwner,
			Sequence:      md.Sequence,
			Key:           md.Key,
			Value:         metadataByteArrayToValue(md.Value, c.Runtime),
		}
	}

	return c.Runtime.ToValue(result)
}

func UserMetadataListFromSlice(c *actions.FieldConfig, metadata []query.UserMetadata) goja.Value {
	result := &userMetadataList{
		// Count was the only field ever queried from the DB in the old implementation,
		// so Sequence and LastRun are omitted.
		Count:    uint64(len(metadata)),
		Metadata: make([]*userMetadata, len(metadata)),
	}
	for i, md := range metadata {
		result.Metadata[i] = &userMetadata{
			CreationDate:  md.CreationDate,
			ChangeDate:    md.ChangeDate,
			ResourceOwner: md.ResourceOwner,
			Sequence:      md.Sequence,
			Key:           md.Key,
			Value:         metadataByteArrayToValue(md.Value, c.Runtime),
		}
	}

	return c.Runtime.ToValue(result)
}

func GetOrganizationMetadata(ctx context.Context, queries *query.Queries, c *actions.FieldConfig, organizationID string) goja.Value {
	metadata, err := queries.SearchOrgMetadata(
		ctx,
		false,
		organizationID,
		&query.OrgMetadataSearchQueries{},
		false,
		false,
	)
	if err != nil {
		logging.WithError(err).Info("unable to get org metadata in action")
		panic(err)
	}
	return OrgMetadataListFromQuery(c, metadata)
}

func metadataByteArrayToValue(val []byte, runtime *goja.Runtime) goja.Value {
	var value interface{}
	if !json.Valid(val) {
		var err error
		val, err = json.Marshal(string(val))
		if err != nil {
			logging.WithError(err).Debug("unable to marshal unknown value")
			panic(err)
		}
	}
	err := json.Unmarshal(val, &value)
	if err != nil {
		logging.WithError(err).Debug("unable to unmarshal into map")
		panic(err)
	}
	return runtime.ToValue(value)
}

type metadataList struct {
	Count     uint64
	Sequence  uint64
	Timestamp time.Time
	Metadata  []*metadata
}

type metadata struct {
	CreationDate  time.Time
	ChangeDate    time.Time
	ResourceOwner string
	Sequence      uint64
	Key           string
	Value         goja.Value
}

type userMetadataList struct {
	Count     uint64
	Sequence  uint64
	Timestamp time.Time
	Metadata  []*userMetadata
}

type userMetadata struct {
	CreationDate  time.Time
	ChangeDate    time.Time
	ResourceOwner string
	Sequence      uint64
	Key           string
	Value         goja.Value
}

type MetadataList struct {
	metadata []*Metadata
}

type Metadata struct {
	Key string
	// Value is for exporting to javascript
	Value goja.Value
	// value is for mapping to [domain.Metadata]
	value []byte
}

func (md *MetadataList) AppendMetadataFunc(call goja.FunctionCall) goja.Value {
	if len(call.Arguments) != 2 {
		panic("exactly 2 (key, value) arguments expected")
	}

	value, err := json.Marshal(call.Arguments[1].Export())
	if err != nil {
		logging.WithError(err).Debug("unable to marshal")
		panic(err)
	}

	md.metadata = append(md.metadata,
		&Metadata{
			Key:   call.Arguments[0].Export().(string),
			Value: call.Arguments[1],
			value: value,
		})
	return nil
}

// AppendMetadataRawFunc appends a metadata entry storing the value as raw bytes
// without JSON encoding. In contrast to [MetadataList.AppendMetadataFunc], a string
// is stored as its plain UTF-8 bytes (e.g. `de` instead of `"de"`).
// Allowed values are strings and byte arrays (Uint8Array or an array of integers 0-255).
func (md *MetadataList) AppendMetadataRawFunc(call goja.FunctionCall) goja.Value {
	if len(call.Arguments) != 2 {
		panic("exactly 2 (key, value) arguments expected")
	}

	value := rawMetadataValue(call.Arguments[1].Export())
	if len(value) == 0 {
		panic("value must not be empty")
	}

	md.metadata = append(md.metadata,
		&Metadata{
			Key:   call.Arguments[0].Export().(string),
			Value: call.Arguments[1],
			value: value,
		})
	return nil
}

// rawMetadataValue converts an exported goja value to raw bytes.
// Strings are converted to their UTF-8 bytes, Uint8Array is exported as []byte directly
// and plain arrays must only contain integers in the range 0-255.
func rawMetadataValue(v interface{}) []byte {
	switch value := v.(type) {
	case string:
		return []byte(value)
	case []byte:
		return value
	case []interface{}:
		bytes := make([]byte, len(value))
		for i, item := range value {
			b, ok := item.(int64)
			if !ok || b < 0 || b > 255 {
				panic("array value must only contain integers between 0 and 255")
			}
			bytes[i] = byte(b)
		}
		return bytes
	default:
		panic("value must be a string or byte array")
	}
}

func (md *MetadataList) MetadataListFromDomain(runtime *goja.Runtime) interface{} {
	for i, metadata := range md.metadata {
		md.metadata[i].Value = metadataByteArrayToValue(metadata.value, runtime)
	}
	return &md.metadata
}

func MetadataListFromDomain(metadata []*domain.Metadata) *MetadataList {
	list := &MetadataList{metadata: make([]*Metadata, len(metadata))}

	for i, md := range metadata {
		list.metadata[i] = &Metadata{
			Key:   md.Key,
			value: md.Value,
		}
	}

	return list
}

func MetadataListToDomain(metadataList *MetadataList) []*domain.Metadata {
	if metadataList == nil {
		return nil
	}

	list := make([]*domain.Metadata, len(metadataList.metadata))
	for i, metadata := range metadataList.metadata {
		value := metadata.value
		if len(value) == 0 {
			var err error
			value, err = json.Marshal(metadata.Value.Export())
			if err != nil {
				logging.WithError(err).Debug("unable to marshal")
				panic(err)
			}
		}
		list[i] = &domain.Metadata{
			Key:   metadata.Key,
			Value: value,
		}
	}

	return list
}
