package object

import (
	"encoding/json"
	"time"

	"github.com/dop251/goja"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/actions"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
)

func UserMetadataListFromQuery(c *actions.FieldConfig, metadata *query.UserMetadataList) goja.Value {
	result := &userMetadataList{
		Count:     metadata.Count,
		Sequence:  metadata.Sequence,
		Timestamp: metadata.Timestamp,
		Metadata:  make([]*userMetadata, len(metadata.Metadata)),
	}

	for i, md := range metadata.Metadata {
		var value interface{}
		if !json.Valid(md.Value) {
			var err error
			md.Value, err = json.Marshal(string(md.Value))
			if err != nil {
				logging.WithError(err).Debug("unable to marshal unknow value")
				panic(err)
			}
		}
		err := json.Unmarshal(md.Value, &value)
		if err != nil {
			logging.WithError(err).Debug("unable to unmarshal into map")
			panic(err)
		}
		result.Metadata[i] = &userMetadata{
			CreationDate:  md.CreationDate,
			ChangeDate:    md.ChangeDate,
			ResourceOwner: md.ResourceOwner,
			Sequence:      md.Sequence,
			Key:           md.Key,
			Value:         c.Runtime.ToValue(value),
		}
	}

	return c.Runtime.ToValue(result)
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
	Metadata []*Metadata
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

	md.Metadata = append(md.Metadata,
		&Metadata{
			Key:   call.Arguments[0].Export().(string),
			Value: call.Arguments[1],
			value: value,
		})
	return nil
}

func MetadataListToDomain(metadataList *MetadataList) []*domain.Metadata {
	if metadataList == nil {
		return nil
	}

	list := make([]*domain.Metadata, len(metadataList.Metadata))
	for i, metadata := range metadataList.Metadata {
		list[i] = &domain.Metadata{
			Key:   metadata.Key,
			Value: metadata.value,
		}
	}

	return list
}

func MetadataField(metadata *MetadataList) func(c *actions.FieldConfig) interface{} {
	return func(c *actions.FieldConfig) interface{} {
		for _, md := range metadata.Metadata {
			if json.Valid(md.value) {
				err := json.Unmarshal(md.value, &md.Value)
				if err != nil {
					panic(err)
				}
			}
		}

		return metadata.Metadata
	}
}

func MetadataListFromDomain(metadata []*domain.Metadata) *MetadataList {
	list := &MetadataList{Metadata: make([]*Metadata, len(metadata))}

	for i, md := range metadata {
		var val interface{}
		if json.Valid(md.Value) {
			err := json.Unmarshal(md.Value, &val)
			if err != nil {
				panic(err)
			}
		}

		list.Metadata[i] = &Metadata{
			Key:   md.Key,
			value: md.Value,
		}
	}

	return list
}
