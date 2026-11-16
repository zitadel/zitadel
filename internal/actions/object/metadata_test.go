package object

import (
	"testing"

	"github.com/dop251/goja"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/domain"
)

func TestMetadataListToDomain(t *testing.T) {
	type args struct {
		metadataList *MetadataList
	}
	tests := []struct {
		name string
		args args
		want []*domain.Metadata
	}{
		{
			name: "nil",
			args: args{metadataList: nil},
			want: nil,
		},
		{
			name: "empty",
			args: args{metadataList: &MetadataList{}},
			want: []*domain.Metadata{},
		},
		{
			name: "from mapped value",
			args: args{metadataList: &MetadataList{
				metadata: []*Metadata{
					{
						Key:   "key1",
						value: []byte("value1"),
					},
				},
			}},
			want: []*domain.Metadata{
				{
					Key:   "key1",
					Value: []byte("value1"),
				},
			},
		},
		{
			name: "from goja value string",
			args: args{metadataList: &MetadataList{
				metadata: []*Metadata{
					{
						Key:   "key1",
						Value: (&goja.Runtime{}).ToValue("value1"),
					},
				},
			}},
			want: []*domain.Metadata{
				{
					Key:   "key1",
					Value: []byte(`"value1"`),
				},
			},
		},
		{
			name: "from goja value int",
			args: args{metadataList: &MetadataList{
				metadata: []*Metadata{
					{
						Key:   "key1",
						Value: (&goja.Runtime{}).ToValue(1),
					},
				},
			}},
			want: []*domain.Metadata{
				{
					Key:   "key1",
					Value: []byte("1"),
				},
			},
		},
		{
			name: "from goja value float",
			args: args{metadataList: &MetadataList{
				metadata: []*Metadata{
					{
						Key:   "key1",
						Value: (&goja.Runtime{}).ToValue(1.2),
					},
				},
			}},
			want: []*domain.Metadata{
				{
					Key:   "key1",
					Value: []byte("1.2"),
				},
			},
		},
		{
			name: "from goja value bool",
			args: args{metadataList: &MetadataList{
				metadata: []*Metadata{
					{
						Key:   "key1",
						Value: (&goja.Runtime{}).ToValue(true),
					},
				},
			}},
			want: []*domain.Metadata{
				{
					Key:   "key1",
					Value: []byte("true"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, MetadataListToDomain(tt.args.metadataList))
		})
	}
}

func TestAppendMetadataRawFunc(t *testing.T) {
	tests := []struct {
		name   string
		script string
		want   []*domain.Metadata
	}{
		{
			name:   "string is stored as raw bytes without quotes",
			script: `appendMetadataRaw("locale", "de")`,
			want: []*domain.Metadata{
				{
					Key:   "locale",
					Value: []byte("de"),
				},
			},
		},
		{
			name:   "json string is stored verbatim",
			script: `appendMetadataRaw("tokenPayload", '{"acr":"phrh"}')`,
			want: []*domain.Metadata{
				{
					Key:   "tokenPayload",
					Value: []byte(`{"acr":"phrh"}`),
				},
			},
		},
		{
			name:   "uint8 array is stored as raw bytes",
			script: `appendMetadataRaw("locale", new Uint8Array([100, 101]))`,
			want: []*domain.Metadata{
				{
					Key:   "locale",
					Value: []byte("de"),
				},
			},
		},
		{
			name:   "integer array is stored as raw bytes",
			script: `appendMetadataRaw("locale", [100, 101])`,
			want: []*domain.Metadata{
				{
					Key:   "locale",
					Value: []byte("de"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runtime := goja.New()
			metadataList := &MetadataList{}
			require.NoError(t, runtime.Set("appendMetadataRaw", metadataList.AppendMetadataRawFunc))
			_, err := runtime.RunString(tt.script)
			require.NoError(t, err)
			assert.Equal(t, tt.want, MetadataListToDomain(metadataList))
		})
	}
}

func TestAppendMetadataRawFuncPanics(t *testing.T) {
	runtime := goja.New()
	tests := []struct {
		name string
		args []goja.Value
	}{
		{
			name: "wrong argument count",
			args: []goja.Value{runtime.ToValue("locale")},
		},
		{
			name: "empty string value",
			args: []goja.Value{runtime.ToValue("locale"), runtime.ToValue("")},
		},
		{
			name: "empty array value",
			args: []goja.Value{runtime.ToValue("locale"), runtime.ToValue([]interface{}{})},
		},
		{
			name: "number value",
			args: []goja.Value{runtime.ToValue("locale"), runtime.ToValue(1)},
		},
		{
			name: "object value",
			args: []goja.Value{runtime.ToValue("locale"), runtime.ToValue(map[string]interface{}{"locale": "de"})},
		},
		{
			name: "array value out of byte range",
			args: []goja.Value{runtime.ToValue("locale"), runtime.ToValue([]interface{}{int64(300)})},
		},
		{
			name: "array value with non integer",
			args: []goja.Value{runtime.ToValue("locale"), runtime.ToValue([]interface{}{1.5})},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metadataList := &MetadataList{}
			assert.Panics(t, func() {
				metadataList.AppendMetadataRawFunc(goja.FunctionCall{Arguments: tt.args})
			})
		})
	}
}

// TestAppendMetadataFuncJSONEncodes locks in that the existing appendMetadata
// still JSON-encodes its value (a string is stored quoted), in contrast to appendMetadataRaw.
func TestAppendMetadataFuncJSONEncodes(t *testing.T) {
	runtime := goja.New()
	metadataList := &MetadataList{}
	require.NoError(t, runtime.Set("appendMetadata", metadataList.AppendMetadataFunc))
	require.NoError(t, runtime.Set("appendMetadataRaw", metadataList.AppendMetadataRawFunc))
	_, err := runtime.RunString(`appendMetadata("locale", "de"); appendMetadataRaw("locale", "de")`)
	require.NoError(t, err)
	assert.Equal(t, []*domain.Metadata{
		{
			Key:   "locale",
			Value: []byte(`"de"`),
		},
		{
			Key:   "locale",
			Value: []byte(`de`),
		},
	}, MetadataListToDomain(metadataList))
}
