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
		true,
		organizationID,
		&query.OrgMetadataSearchQueries{},
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
			value = mapBytesToByteArray(metadata.Value.Export())
		}
		list[i] = &domain.Metadata{
			Key:   metadata.Key,
			Value: value,
		}
	}

	return list
}

// mapBytesToByteArray is used for backwards compatibility of old metadata.push method
// converts the Javascript uint8 array which is exported as []interface{} to a []byte
func mapBytesToByteArray(i interface{}) []byte {
	bytes, ok := i.([]interface{})
	if !ok {
		return nil
	}
	value := make([]byte, len(bytes))
	for i, val := range bytes {
		b, ok := val.(int64)
		if !ok {
			return nil
		}
		value[i] = byte(b)
	}
	return value
}
