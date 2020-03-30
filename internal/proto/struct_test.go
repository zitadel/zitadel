package proto

import (
	"testing"

	pb_struct "github.com/golang/protobuf/ptypes/struct"
)

func Test_ToPBStruct(t *testing.T) {
	type obj struct {
		ID  string
		Seq uint64
	}
	type args struct {
		obj obj
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		length  int
		result  obj
	}{
		{
			name: "to pb stuct",
			args: args{
				obj: obj{ID: "ID", Seq: 12345},
			},
			wantErr: false,
			length:  2,
			result:  obj{ID: "ID", Seq: 12345},
		},
		{
			name: "empty struct",
			args: args{
				obj: obj{},
			},
			wantErr: false,
			length:  2,
			result:  obj{ID: "", Seq: 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields, err := ToPBStruct(tt.args.obj)
			if tt.wantErr && err == nil {
				t.Errorf("got wrong result, should get err: actual: %v ", err)
			}
			if !tt.wantErr && len(fields.Fields) != tt.length {
				t.Errorf("got wrong result length, expecting: %v, actual: %v ", tt.length, len(fields.Fields))
			}
			if !tt.wantErr && tt.result.ID != fields.Fields["ID"].GetStringValue() {
				t.Errorf("got wrong result, ID should be same: expecting: %v, actual: %v ", tt.result.ID, fields.Fields["ID"].GetStringValue())
			}
			if !tt.wantErr && int(tt.result.Seq) != int(fields.Fields["Seq"].GetNumberValue()) {
				t.Errorf("got wrong result, Seq should be same: expecting: %v, actual: %v ", tt.result.Seq, fields.Fields["Seq"].GetStringValue())
			}
		})
	}
}

func Test_FromPBStruct(t *testing.T) {
	type obj struct {
		ID  string
		Seq uint64
	}
	type args struct {
		obj    *obj
		fields *pb_struct.Struct
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		result  obj
	}{
		{
			name: "from pb stuct",
			args: args{
				obj: &obj{},
				fields: &pb_struct.Struct{Fields: map[string]*pb_struct.Value{
					"ID":  &pb_struct.Value{Kind: &pb_struct.Value_StringValue{StringValue: "ID"}},
					"Seq": &pb_struct.Value{Kind: &pb_struct.Value_NumberValue{NumberValue: 12345}},
				},
				},
			},
			wantErr: false,
			result:  obj{ID: "ID", Seq: 12345},
		},
		{
			name: "no fields",
			args: args{
				obj: &obj{},
				fields: &pb_struct.Struct{Fields: map[string]*pb_struct.Value{},
				},
			},
			wantErr: false,
			result:  obj{ID: "", Seq: 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := FromPBStruct(tt.args.obj, tt.args.fields)
			if tt.wantErr && err == nil {
				t.Errorf("got wrong result, should get err: actual: %v ", err)
			}
			if !tt.wantErr && tt.result.ID != tt.args.obj.ID {
				t.Errorf("got wrong result, ID should be same: expecting: %v, actual: %v ", tt.result.ID, tt.args.obj.ID)
			}
			if !tt.wantErr && int(tt.result.Seq) != int(tt.args.obj.Seq) {
				t.Errorf("got wrong result, Seq should be same: expecting: %v, actual: %v ", tt.result.Seq, tt.args.obj.Seq)
			}
		})
	}
}
