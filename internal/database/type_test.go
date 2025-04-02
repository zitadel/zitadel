package database

import (
	"database/sql"
	"database/sql/driver"
	"testing"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMap_Scan(t *testing.T) {
	type args struct {
		src []byte
	}
	type res[V any] struct {
		want Map[V]
		err  bool
	}
	type testCase[V any] struct {
		name string
		m    Map[V]
		args args
		res[V]
	}
	tests := []testCase[string]{
		{
			"nil",
			Map[string]{},
			args{src: nil},
			res[string]{
				want: Map[string]{},
				err:  false,
			},
		},
		{
			"null",
			Map[string]{},
			args{src: []byte("invalid")},
			res[string]{
				want: Map[string]{},
				err:  true,
			},
		},
		{
			"null",
			Map[string]{},
			args{src: nil},
			res[string]{
				want: Map[string]{},
			},
		},
		{
			"empty",
			Map[string]{},
			args{src: []byte(`{}`)},
			res[string]{
				want: Map[string]{},
			},
		},
		{
			"set",
			Map[string]{},
			args{src: []byte(`{"key": "value"}`)},
			res[string]{
				want: Map[string]{
					"key": "value",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.Scan(tt.args.src); (err != nil) != tt.res.err {
				t.Errorf("Scan() error = %v, wantErr %v", err, tt.res.err)
			}
			assert.Equal(t, tt.res.want, tt.m)
		})
	}
}

func TestMap_Value(t *testing.T) {
	type res struct {
		want driver.Value
		err  bool
	}
	type testCase[V any] struct {
		name string
		m    Map[V]
		res  res
	}
	tests := []testCase[string]{
		{
			"nil",
			nil,
			res{
				want: nil,
			},
		},
		{
			"empty",
			Map[string]{},
			res{
				want: nil,
			},
		},
		{
			"set",
			Map[string]{
				"key": "value",
			},
			res{
				want: driver.Value([]byte(`{"key":"value"}`)),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.Value()
			if tt.res.err {
				assert.Error(t, err)
			}
			if !tt.res.err {
				require.NoError(t, err)
				assert.Equalf(t, tt.res.want, got, "Value()")
			}
		})
	}
}

type typedInt int

func TestNumberArray_Scan(t *testing.T) {
	type args struct {
		src any
	}
	type res struct {
		want any
		err  bool
	}
	type testCase struct {
		name string
		m    sql.Scanner
		args args
		res  res
	}
	tests := []testCase{
		{
			name: "typedInt",
			m:    new(NumberArray[typedInt]),
			args: args{src: "{1,2}"},
			res: res{
				want: &NumberArray[typedInt]{1, 2},
			},
		},
		{
			name: "int8",
			m:    new(NumberArray[int8]),
			args: args{src: "{1,2}"},
			res: res{
				want: &NumberArray[int8]{1, 2},
			},
		},
		{
			name: "uint8",
			m:    new(NumberArray[uint8]),
			args: args{src: "{1,2}"},
			res: res{
				want: &NumberArray[uint8]{1, 2},
			},
		},
		{
			name: "int16",
			m:    new(NumberArray[int16]),
			args: args{src: "{1,2}"},
			res: res{
				want: &NumberArray[int16]{1, 2},
			},
		},
		{
			name: "uint16",
			m:    new(NumberArray[uint16]),
			args: args{src: "{1,2}"},
			res: res{
				want: &NumberArray[uint16]{1, 2},
			},
		},
		{
			name: "int32",
			m:    new(NumberArray[int32]),
			args: args{src: "{1,2}"},
			res: res{
				want: &NumberArray[int32]{1, 2},
			},
		},
		{
			name: "uint32",
			m:    new(NumberArray[uint32]),
			args: args{src: "{1,2}"},
			res: res{
				want: &NumberArray[uint32]{1, 2},
			},
		},
		{
			name: "int64",
			m:    new(NumberArray[int64]),
			args: args{src: "{1,2}"},
			res: res{
				want: &NumberArray[int64]{1, 2},
			},
		},
		{
			name: "uint64",
			m:    new(NumberArray[uint64]),
			args: args{src: "{1,2}"},
			res: res{
				want: &NumberArray[uint64]{1, 2},
			},
		},
		{
			name: "int",
			m:    new(NumberArray[int]),
			args: args{src: "{1,2}"},
			res: res{
				want: &NumberArray[int]{1, 2},
			},
		},
		{
			name: "uint",
			m:    new(NumberArray[uint]),
			args: args{src: "{1,2}"},
			res: res{
				want: &NumberArray[uint]{1, 2},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.Scan(tt.args.src); (err != nil) != tt.res.err {
				t.Errorf("Scan() error = %v, wantErr %v", err, tt.res.err)
			}

			assert.Equal(t, tt.res.want, tt.m)
		})
	}
}

type typedText string

func TestTextArray_Scan(t *testing.T) {
	type args struct {
		src any
	}
	type res struct {
		want sql.Scanner
		err  bool
	}
	type testCase struct {
		name string
		m    sql.Scanner
		args args
		res
	}
	tests := []testCase{
		{
			"string",
			new(TextArray[string]),
			args{src: "{asdf,fdas}"},
			res{
				want: &TextArray[string]{"asdf", "fdas"},
			},
		},
		{
			"typedText",
			new(TextArray[typedText]),
			args{src: "{asdf,fdas}"},
			res{
				want: &TextArray[typedText]{"asdf", "fdas"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.Scan(tt.args.src); (err != nil) != tt.res.err {
				t.Errorf("Scan() error = %v, wantErr %v", err, tt.res.err)
			}

			assert.Equal(t, tt.res.want, tt.m)
		})
	}
}

func TestTextArray_Value(t *testing.T) {
	type res struct {
		want driver.Value
		err  bool
	}
	type testCase struct {
		name string
		m    driver.Valuer
		res  res
	}
	tests := []testCase{
		{
			"empty",
			TextArray[string]{},
			res{
				want: nil,
			},
		},
		{
			"set",
			TextArray[string]{"a", "s", "d", "f"},
			res{
				want: driver.Value([]byte("{a,s,d,f}")),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.Value()
			if tt.res.err {
				assert.Error(t, err)
			}
			if !tt.res.err {
				require.NoError(t, err)
				assert.Equalf(t, tt.res.want, got, "Value()")
			}
		})
	}
}

type typedByte byte

func TestByteArray_Scan(t *testing.T) {
	wantedBytes := []byte("asdf")
	wantedTypedBytes := []typedByte("asdf")
	type args struct {
		src any
	}
	type res struct {
		want sql.Scanner
		err  bool
	}
	type testCase struct {
		name string
		m    sql.Scanner
		args args
		res
	}
	tests := []testCase{
		{
			"bytes",
			new(ByteArray[byte]),
			args{src: []byte("asdf")},
			res{
				want: (*ByteArray[byte])(&wantedBytes),
			},
		},
		{
			"typed",
			new(ByteArray[typedByte]),
			args{src: []byte("asdf")},
			res{
				want: (*ByteArray[typedByte])(&wantedTypedBytes),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.Scan(tt.args.src); (err != nil) != tt.res.err {
				t.Errorf("Scan() error = %v, wantErr %v", err, tt.res.err)
			}

			assert.Equal(t, tt.res.want, tt.m)
		})
	}
}

func TestByteArray_Value(t *testing.T) {
	type res struct {
		want driver.Value
		err  bool
	}
	type testCase struct {
		name string
		m    driver.Valuer
		res  res
	}
	tests := []testCase{
		{
			"empty",
			ByteArray[byte]{},
			res{
				want: nil,
			},
		},
		{
			"set",
			ByteArray[byte]([]byte("{\"type\": \"object\", \"$schema\": \"urn:zitadel:schema:v1\"}")),
			res{
				want: driver.Value([]byte("{\"type\": \"object\", \"$schema\": \"urn:zitadel:schema:v1\"}")),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.Value()
			if tt.res.err {
				assert.Error(t, err)
			}
			if !tt.res.err {
				require.NoError(t, err)
				assert.Equalf(t, tt.res.want, got, "Value()")
			}
		})
	}
}

func TestDuration_Scan(t *testing.T) {
	duration := Duration(10)
	type args struct {
		src any
	}
	type res struct {
		want sql.Scanner
		err  bool
	}
	type testCase[V ~string] struct {
		name string
		m    sql.Scanner
		args args
		res
	}
	tests := []testCase[string]{
		{
			name: "int64",
			m:    new(Duration),
			args: args{src: int64(duration)},
			res: res{
				want: &duration,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.Scan(tt.args.src); (err != nil) != tt.res.err {
				t.Errorf("Scan() error = %v, wantErr %v", err, tt.res.err)
			}

			assert.Equal(t, tt.res.want, tt.m)
		})
	}
}

func TestJSONArray_Scan(t *testing.T) {
	type args struct {
		src any
	}
	tests := []struct {
		name    string
		args    args
		want    *JSONArray[string]
		wantErr bool
	}{
		{
			name:    "nil",
			args:    args{src: nil},
			want:    new(JSONArray[string]),
			wantErr: false,
		},
		{
			name:    "zero bytes",
			args:    args{src: []byte("")},
			want:    new(JSONArray[string]),
			wantErr: false,
		},
		{
			name:    "empty",
			args:    args{src: []byte("[]")},
			want:    gu.Ptr(JSONArray[string]{}),
			wantErr: false,
		},
		{
			name:    "ok",
			args:    args{src: []byte("[\"a\", \"b\"]")},
			want:    gu.Ptr(JSONArray[string]{"a", "b"}),
			wantErr: false,
		},
		{
			name:    "json error",
			args:    args{src: []byte("{\"a\": \"b\"}")},
			want:    new(JSONArray[string]),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := new(JSONArray[string])
			err := got.Scan(tt.args.src)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestJSONArray_Value(t *testing.T) {
	tests := []struct {
		name string
		a    []string
		want driver.Value
	}{
		{
			name: "nil",
			a:    nil,
			want: nil,
		},
		{
			name: "empty",
			a:    []string{},
			want: []byte("[]"),
		},
		{
			name: "ok",
			a:    []string{"a", "b"},
			want: []byte("[\"a\",\"b\"]"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewJSONArray(tt.a).Value()
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
