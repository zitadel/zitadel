package proto

import (
	"bytes"
	"encoding/json"

	"github.com/golang/protobuf/jsonpb"
	pb_struct "github.com/golang/protobuf/ptypes/struct"

	"github.com/caos/logging"
)

func MustToPBStruct(object interface{}) *pb_struct.Struct {
	s, err := ToPBStruct(object)
	logging.Log("PROTO-7Aa3t").OnError(err).Panic("unable to map object to pb-struct")

	return s
}

func BytesToPBStruct(b []byte) (*pb_struct.Struct, error) {
	fields := new(pb_struct.Struct)
	err := jsonpb.Unmarshal(bytes.NewReader(b), fields)
	return fields, err
}

func ToPBStruct(object interface{}) (*pb_struct.Struct, error) {
	fields := new(pb_struct.Struct)

	marshalled, err := json.Marshal(object)
	if err != nil {
		return nil, err
	}

	err = jsonpb.Unmarshal(bytes.NewReader(marshalled), fields)

	return fields, err
}

func MustFromPBStruct(object interface{}, s *pb_struct.Struct) {
	err := FromPBStruct(object, s)
	logging.Log("PROTO-WeMYY").OnError(err).Panic("unable to map pb-struct into object")
}

func FromPBStruct(object interface{}, s *pb_struct.Struct) error {
	marshaller := new(jsonpb.Marshaler)
	jsonString, err := marshaller.MarshalToString(s)
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(jsonString), object)
}
