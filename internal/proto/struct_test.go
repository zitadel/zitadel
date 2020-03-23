package proto

import (
	"testing"

	pb_struct "github.com/golang/protobuf/ptypes/struct"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToPBStruct(t *testing.T) {
	obj := struct {
		ID   string
		Name string
		Seq  uint64
	}{
		ID:   "asdf",
		Name: "ueli",
		Seq:  208582075,
	}
	fields, err := ToPBStruct(&obj)
	require.NoError(t, err)
	require.Len(t, fields.Fields, 3)

	assert.Equal(t, obj.ID, fields.Fields["ID"].GetStringValue())
	assert.Equal(t, int(obj.Seq), int(fields.Fields["Seq"].GetNumberValue()))
	assert.Equal(t, obj.Name, fields.Fields["Name"].GetStringValue())
}

func TestFromPBStruct(t *testing.T) {
	name := "ueli"
	id := "asdf"
	seq := float64(208582075)
	s := &pb_struct.Struct{Fields: map[string]*pb_struct.Value{
		"ID":   &pb_struct.Value{Kind: &pb_struct.Value_StringValue{StringValue: id}},
		"Name": &pb_struct.Value{Kind: &pb_struct.Value_StringValue{StringValue: name}},
		"Seq":  &pb_struct.Value{Kind: &pb_struct.Value_NumberValue{NumberValue: seq}},
	}}

	obj := struct {
		ID   string
		Name string
		Seq  uint64
	}{}

	err := FromPBStruct(&obj, s)
	require.NoError(t, err)

	assert.Equal(t, id, obj.ID)
	assert.Equal(t, name, obj.Name)
	assert.Equal(t, int(seq), int(obj.Seq))
}
