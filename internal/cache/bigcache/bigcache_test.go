package bigcache

import (
	a_cache "github.com/allegro/bigcache"
	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	"reflect"
	"testing"
)

type TestStruct struct {
	Test string
}

func getBigCacheMock() *Bigcache {
	cache, _ := a_cache.NewBigCache(a_cache.DefaultConfig(2000))
	return &Bigcache{cache: cache}
}

func TestSet(t *testing.T) {
	type args struct {
		cache *Bigcache
		key   string
		value *TestStruct
	}
	type res struct {
		result  *TestStruct
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "set cache no err",
			args: args{
				cache: getBigCacheMock(),
				key:   "KEY",
				value: &TestStruct{Test: "Test"},
			},
			res: res{
				result: &TestStruct{},
			},
		},
		{
			name: "key empty",
			args: args{
				cache: getBigCacheMock(),
				key:   "",
				value: &TestStruct{Test: "Test"},
			},
			res: res{
				errFunc: errors.IsErrorInvalidArgument,
			},
		},
		{
			name: "set cache nil value",
			args: args{
				cache: getBigCacheMock(),
				key:   "KEY",
			},
			res: res{
				errFunc: errors.IsErrorInvalidArgument,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.cache.Set(tt.args.key, tt.args.value)

			if tt.res.errFunc == nil && err != nil {
				t.Errorf("got wrong result should not get err: %v ", err)
			}

			if tt.res.errFunc == nil {
				tt.args.cache.Get(tt.args.key, tt.res.result)
				if tt.res.result == nil {
					t.Errorf("got wrong result should get result: %v ", err)
				}
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestGet(t *testing.T) {
	type args struct {
		event    []*es_models.Event
		cache    *Bigcache
		key      string
		setValue *TestStruct
		getValue *TestStruct
	}
	type res struct {
		result  *TestStruct
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "get cache no err",
			args: args{
				cache:    getBigCacheMock(),
				key:      "KEY",
				setValue: &TestStruct{Test: "Test"},
				getValue: &TestStruct{Test: "Test"},
			},
			res: res{
				result: &TestStruct{Test: "Test"},
			},
		},
		{
			name: "get cache no key",
			args: args{
				cache:    getBigCacheMock(),
				setValue: &TestStruct{Test: "Test"},
				getValue: &TestStruct{Test: "Test"},
			},
			res: res{
				errFunc: errors.IsErrorInvalidArgument,
			},
		},
		{
			name: "get cache no value",
			args: args{
				cache:    getBigCacheMock(),
				key:      "KEY",
				setValue: &TestStruct{Test: "Test"},
			},
			res: res{
				errFunc: errors.IsErrorInvalidArgument,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.cache.Set("KEY", tt.args.setValue)
			if err != nil {
				t.Errorf("something went wrong")
			}

			err = tt.args.cache.Get(tt.args.key, tt.args.getValue)

			if tt.res.errFunc == nil && err != nil {
				t.Errorf("got wrong result should not get err: %v ", err)
			}

			if tt.res.errFunc == nil && !reflect.DeepEqual(tt.args.getValue, tt.res.result) {
				t.Errorf("got wrong result expected: %v actual: %v", tt.res.result, tt.args.getValue)
			}

			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	type args struct {
		event    []*es_models.Event
		cache    *Bigcache
		key      string
		setValue *TestStruct
		getValue *TestStruct
	}
	type res struct {
		result  *TestStruct
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "delete cache no err",
			args: args{
				cache:    getBigCacheMock(),
				key:      "KEY",
				setValue: &TestStruct{Test: "Test"},
			},
			res: res{},
		},
		{
			name: "get cache no key",
			args: args{
				cache:    getBigCacheMock(),
				setValue: &TestStruct{Test: "Test"},
				getValue: &TestStruct{Test: "Test"},
			},
			res: res{
				errFunc: errors.IsErrorInvalidArgument,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.cache.Set("KEY", tt.args.setValue)
			if err != nil {
				t.Errorf("something went wrong")
			}

			err = tt.args.cache.Delete(tt.args.key)

			if tt.res.errFunc == nil && err != nil {
				t.Errorf("got wrong result should not get err: %v ", err)
			}

			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}
