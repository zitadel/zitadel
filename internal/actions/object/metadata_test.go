package object

import (
	"testing"

	"github.com/dop251/goja"
	"github.com/stretchr/testify/assert"

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
