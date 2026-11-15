package redis

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/cache"
)

type testIndex int

const (
	testIndexID testIndex = iota
	testIndexName
)

const (
	testDB = 99
)

var testIndices = []testIndex{
	testIndexID,
	testIndexName,
}

type testObject struct {
	ID   string
	Name []string
}

func (o *testObject) Keys(index testIndex) []string {
	switch index {
	case testIndexID:
		return []string{o.ID}
	case testIndexName:
		return o.Name
	default:
		return nil
	}
}

func Test_redisCache_set(t *testing.T) {
	type args struct {
		ctx   context.Context
		value *testObject
	}
	tests := []struct {
		name       string
		config     cache.Config
		args       args
		assertions func(t *testing.T, s *miniredis.Miniredis, objectID string)
		wantErr    error
	}{
		{
			name:   "ok",
			config: cache.Config{},
			args: args{
				ctx: context.Background(),
				value: &testObject{
					ID:   "one",
					Name: []string{"foo", "bar"},
				},
			},
			assertions: func(t *testing.T, s *miniredis.Miniredis, objectID string) {
				s.CheckGet(t, "VERSION:0:one", objectID)
				s.CheckGet(t, "VERSION:1:foo", objectID)
				s.CheckGet(t, "VERSION:1:bar", objectID)
				assert.Empty(t, s.HGet(objectID, "expiry"))
				assert.JSONEq(t, `{"ID":"one","Name":["foo","bar"]}`, s.HGet(objectID, "object"))
				assertKeysTTL(t, s, objectID, 0)
			},
		},
		{
			name: "with last use TTL",
			config: cache.Config{
				LastUseAge: time.Second,
			},
			args: args{
				ctx: context.Background(),
				value: &testObject{
					ID:   "one",
					Name: []string{"foo", "bar"},
				},
			},
			assertions: func(t *testing.T, s *miniredis.Miniredis, objectID string) {
				s.CheckGet(t, "VERSION:0:one", objectID)
				s.CheckGet(t, "VERSION:1:foo", objectID)
				s.CheckGet(t, "VERSION:1:bar", objectID)
				assert.Empty(t, s.HGet(objectID, "expiry"))
				assert.JSONEq(t, `{"ID":"one","Name":["foo","bar"]}`, s.HGet(objectID, "object"))
				assert.Positive(t, s.TTL(objectID))
				assertKeysTTL(t, s, objectID, time.Second)

				s.FastForward(2 * time.Second)
				v, err := s.Get(objectID)
				require.Error(t, err)
				assert.Empty(t, v)
				assert.Empty(t, s.Keys())
			},
		},
		{
			name: "with last use TTL and max age",
			config: cache.Config{
				MaxAge:     time.Minute,
				LastUseAge: time.Second,
			},
			args: args{
				ctx: context.Background(),
				value: &testObject{
					ID:   "one",
					Name: []string{"foo", "bar"},
				},
			},
			assertions: func(t *testing.T, s *miniredis.Miniredis, objectID string) {
				s.CheckGet(t, "VERSION:0:one", objectID)
				s.CheckGet(t, "VERSION:1:foo", objectID)
				s.CheckGet(t, "VERSION:1:bar", objectID)
				assert.NotEmpty(t, s.HGet(objectID, "expiry"))
				assert.JSONEq(t, `{"ID":"one","Name":["foo","bar"]}`, s.HGet(objectID, "object"))
				assert.Equal(t, time.Second, s.TTL(objectID))
				// the object is never returned after max age, so the keys expire then
				assertKeysTTL(t, s, objectID, time.Minute)

				s.FastForward(2 * time.Second)
				assert.False(t, s.Exists(objectID))
				s.FastForward(time.Minute)
				assert.Empty(t, s.Keys())
			},
		},
		{
			name: "with max age TTL",
			config: cache.Config{
				MaxAge: time.Minute,
			},
			args: args{
				ctx: context.Background(),
				value: &testObject{
					ID:   "one",
					Name: []string{"foo", "bar"},
				},
			},
			assertions: func(t *testing.T, s *miniredis.Miniredis, objectID string) {
				s.CheckGet(t, "VERSION:0:one", objectID)
				s.CheckGet(t, "VERSION:1:foo", objectID)
				s.CheckGet(t, "VERSION:1:bar", objectID)
				assert.Empty(t, s.HGet(objectID, "expiry"))
				assert.JSONEq(t, `{"ID":"one","Name":["foo","bar"]}`, s.HGet(objectID, "object"))
				assert.Positive(t, s.TTL(objectID))
				assertKeysTTL(t, s, objectID, time.Minute)

				s.FastForward(2 * time.Minute)
				v, err := s.Get(objectID)
				require.Error(t, err)
				assert.Empty(t, v)
				assert.Empty(t, s.Keys())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, server := prepareCache(t, tt.config)
			rc := c.(*redisCache[testIndex, string, *testObject])
			objectID, err := rc.set(tt.args.ctx, tt.args.value)
			require.ErrorIs(t, err, tt.wantErr)
			t.Log(rc.connector.HGetAll(context.Background(), objectID))
			tt.assertions(t, server, objectID)
		})
	}
}

func Test_redisCache_Get(t *testing.T) {
	type args struct {
		ctx   context.Context
		index testIndex
		key   string
	}
	tests := []struct {
		name        string
		config      cache.Config
		preparation func(t *testing.T, c cache.Cache[testIndex, string, *testObject], s *miniredis.Miniredis)
		args        args
		want        *testObject
		wantOK      bool
		assertions  func(t *testing.T, s *miniredis.Miniredis)
	}{
		{
			name:   "connection error",
			config: cache.Config{},
			preparation: func(_ *testing.T, _ cache.Cache[testIndex, string, *testObject], s *miniredis.Miniredis) {
				s.RequireAuth("foobar")
			},
			args: args{
				ctx:   context.Background(),
				index: testIndexName,
				key:   "foo",
			},
			wantOK: false,
		},
		{
			name:   "get by ID",
			config: cache.Config{},
			preparation: func(t *testing.T, c cache.Cache[testIndex, string, *testObject], s *miniredis.Miniredis) {
				c.Set(context.Background(), &testObject{
					ID:   "one",
					Name: []string{"foo", "bar"},
				})
			},
			args: args{
				ctx:   context.Background(),
				index: testIndexID,
				key:   "one",
			},
			want: &testObject{
				ID:   "one",
				Name: []string{"foo", "bar"},
			},
			wantOK: true,
		},
		{
			name:   "get by name",
			config: cache.Config{},
			preparation: func(t *testing.T, c cache.Cache[testIndex, string, *testObject], s *miniredis.Miniredis) {
				c.Set(context.Background(), &testObject{
					ID:   "one",
					Name: []string{"foo", "bar"},
				})
			},
			args: args{
				ctx:   context.Background(),
				index: testIndexName,
				key:   "foo",
			},
			want: &testObject{
				ID:   "one",
				Name: []string{"foo", "bar"},
			},
			wantOK: true,
		},
		{
			name: "usage timeout",
			config: cache.Config{
				LastUseAge: time.Minute,
			},
			preparation: func(t *testing.T, c cache.Cache[testIndex, string, *testObject], s *miniredis.Miniredis) {
				c.Set(context.Background(), &testObject{
					ID:   "one",
					Name: []string{"foo", "bar"},
				})
				_, ok := c.Get(context.Background(), testIndexID, "one")
				require.True(t, ok)
				s.FastForward(2 * time.Minute)
				assert.Empty(t, s.Keys())
			},
			args: args{
				ctx:   context.Background(),
				index: testIndexName,
				key:   "foo",
			},
			want:   nil,
			wantOK: false,
		},
		{
			name: "usage TTL reset on all keys",
			config: cache.Config{
				LastUseAge: time.Minute,
			},
			preparation: func(t *testing.T, c cache.Cache[testIndex, string, *testObject], s *miniredis.Miniredis) {
				c.Set(context.Background(), &testObject{
					ID:   "one",
					Name: []string{"foo", "bar"},
				})
				s.FastForward(40 * time.Second)
				_, ok := c.Get(context.Background(), testIndexID, "one")
				require.True(t, ok)
				s.FastForward(40 * time.Second)
			},
			args: args{
				ctx:   context.Background(),
				index: testIndexName,
				key:   "foo",
			},
			want: &testObject{
				ID:   "one",
				Name: []string{"foo", "bar"},
			},
			wantOK: true,
			assertions: func(t *testing.T, s *miniredis.Miniredis) {
				// object, key set and the three index keys
				keys := s.Keys()
				assert.Len(t, keys, 5)
				for _, key := range keys {
					assert.Equal(t, time.Minute, s.TTL(key), key)
				}
			},
		},
		{
			name: "usage TTL reset skips keys of a newer object",
			config: cache.Config{
				LastUseAge: time.Minute,
			},
			preparation: func(t *testing.T, c cache.Cache[testIndex, string, *testObject], s *miniredis.Miniredis) {
				c.Set(context.Background(), &testObject{
					ID:   "one",
					Name: []string{"foo", "bar"},
				})
				s.FastForward(30 * time.Second)
				// "bar" now points to the newer object
				c.Set(context.Background(), &testObject{
					ID:   "two",
					Name: []string{"bar"},
				})
				s.FastForward(20 * time.Second)
			},
			args: args{
				ctx:   context.Background(),
				index: testIndexID,
				key:   "one",
			},
			want: &testObject{
				ID:   "one",
				Name: []string{"foo", "bar"},
			},
			wantOK: true,
			assertions: func(t *testing.T, s *miniredis.Miniredis) {
				assert.Equal(t, time.Minute, s.TTL("VERSION:1:foo"))
				assert.Equal(t, 40*time.Second, s.TTL("VERSION:1:bar"))
			},
		},
		{
			name: "usage timeout with max age removes the keys",
			config: cache.Config{
				MaxAge:     time.Minute,
				LastUseAge: time.Second,
			},
			preparation: func(t *testing.T, c cache.Cache[testIndex, string, *testObject], s *miniredis.Miniredis) {
				c.Set(context.Background(), &testObject{
					ID:   "one",
					Name: []string{"foo", "bar"},
				})
				// the object expired, its keys only expire at max age
				s.FastForward(2 * time.Second)
				require.Len(t, s.Keys(), 4)
			},
			args: args{
				ctx:   context.Background(),
				index: testIndexName,
				key:   "foo",
			},
			want:   nil,
			wantOK: false,
			assertions: func(t *testing.T, s *miniredis.Miniredis) {
				assert.Empty(t, s.Keys())
			},
		},
		{
			name: "max age timeout",
			config: cache.Config{
				MaxAge: time.Minute,
			},
			preparation: func(t *testing.T, c cache.Cache[testIndex, string, *testObject], s *miniredis.Miniredis) {
				c.Set(context.Background(), &testObject{
					ID:   "one",
					Name: []string{"foo", "bar"},
				})
				_, ok := c.Get(context.Background(), testIndexID, "one")
				require.True(t, ok)
				s.FastForward(2 * time.Minute)
				assert.Empty(t, s.Keys())
			},
			args: args{
				ctx:   context.Background(),
				index: testIndexName,
				key:   "foo",
			},
			want:   nil,
			wantOK: false,
		},
		{
			name:   "not found",
			config: cache.Config{},
			preparation: func(t *testing.T, c cache.Cache[testIndex, string, *testObject], s *miniredis.Miniredis) {
				c.Set(context.Background(), &testObject{
					ID:   "one",
					Name: []string{"foo", "bar"},
				})
			},
			args: args{
				ctx:   context.Background(),
				index: testIndexName,
				key:   "spanac",
			},
			wantOK: false,
		},
		{
			name:   "json decode error",
			config: cache.Config{},
			preparation: func(t *testing.T, c cache.Cache[testIndex, string, *testObject], s *miniredis.Miniredis) {
				c.Set(context.Background(), &testObject{
					ID:   "one",
					Name: []string{"foo", "bar"},
				})
				objectID, err := s.Get(c.(*redisCache[testIndex, string, *testObject]).redisIndexKeys(testIndexID, "one")[0])
				require.NoError(t, err)
				s.HSet(objectID, "object", "~~~")
			},
			args: args{
				ctx:   context.Background(),
				index: testIndexID,
				key:   "one",
			},
			wantOK: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, server := prepareCache(t, tt.config)
			tt.preparation(t, c, server)
			t.Log(server.Keys())

			got, ok := c.Get(tt.args.ctx, tt.args.index, tt.args.key)
			require.Equal(t, tt.wantOK, ok)
			assert.Equal(t, tt.want, got)
			if tt.assertions != nil {
				tt.assertions(t, server)
			}
		})
	}
}

func Test_redisCache_Invalidate(t *testing.T) {
	type args struct {
		ctx   context.Context
		index testIndex
		key   []string
	}
	tests := []struct {
		name        string
		config      cache.Config
		preparation func(t *testing.T, c cache.Cache[testIndex, string, *testObject], s *miniredis.Miniredis)
		assertions  func(t *testing.T, c cache.Cache[testIndex, string, *testObject])
		args        args
		wantErr     bool
	}{
		{
			name:   "connection error",
			config: cache.Config{},
			preparation: func(_ *testing.T, _ cache.Cache[testIndex, string, *testObject], s *miniredis.Miniredis) {
				s.RequireAuth("foobar")
			},
			args: args{
				ctx:   context.Background(),
				index: testIndexName,
				key:   []string{"foo"},
			},
			wantErr: true,
		},
		{
			name:   "no keys, noop",
			config: cache.Config{},
			preparation: func(t *testing.T, c cache.Cache[testIndex, string, *testObject], s *miniredis.Miniredis) {
				c.Set(context.Background(), &testObject{
					ID:   "one",
					Name: []string{"foo", "bar"},
				})
			},
			args: args{
				ctx:   context.Background(),
				index: testIndexID,
				key:   []string{},
			},
			wantErr: false,
		},
		{
			name:   "invalidate by ID",
			config: cache.Config{},
			preparation: func(t *testing.T, c cache.Cache[testIndex, string, *testObject], s *miniredis.Miniredis) {
				c.Set(context.Background(), &testObject{
					ID:   "one",
					Name: []string{"foo", "bar"},
				})
			},
			assertions: func(t *testing.T, c cache.Cache[testIndex, string, *testObject]) {
				obj, ok := c.Get(context.Background(), testIndexID, "one")
				assert.False(t, ok)
				assert.Nil(t, obj)
				obj, ok = c.Get(context.Background(), testIndexName, "foo")
				assert.False(t, ok)
				assert.Nil(t, obj)
			},
			args: args{
				ctx:   context.Background(),
				index: testIndexID,
				key:   []string{"one"},
			},
			wantErr: false,
		},
		{
			name:   "invalidate by name",
			config: cache.Config{},
			preparation: func(t *testing.T, c cache.Cache[testIndex, string, *testObject], s *miniredis.Miniredis) {
				c.Set(context.Background(), &testObject{
					ID:   "one",
					Name: []string{"foo", "bar"},
				})
			},
			assertions: func(t *testing.T, c cache.Cache[testIndex, string, *testObject]) {
				obj, ok := c.Get(context.Background(), testIndexID, "one")
				assert.False(t, ok)
				assert.Nil(t, obj)
				obj, ok = c.Get(context.Background(), testIndexName, "foo")
				assert.False(t, ok)
				assert.Nil(t, obj)
			},
			args: args{
				ctx:   context.Background(),
				index: testIndexName,
				key:   []string{"foo"},
			},
			wantErr: false,
		},
		{
			name: "invalidate after timeout",
			config: cache.Config{
				LastUseAge: time.Minute,
			},
			preparation: func(t *testing.T, c cache.Cache[testIndex, string, *testObject], s *miniredis.Miniredis) {
				c.Set(context.Background(), &testObject{
					ID:   "one",
					Name: []string{"foo", "bar"},
				})
				_, ok := c.Get(context.Background(), testIndexID, "one")
				require.True(t, ok)
				s.FastForward(2 * time.Minute)
			},
			assertions: func(t *testing.T, c cache.Cache[testIndex, string, *testObject]) {
				obj, ok := c.Get(context.Background(), testIndexID, "one")
				assert.False(t, ok)
				assert.Nil(t, obj)
				obj, ok = c.Get(context.Background(), testIndexName, "foo")
				assert.False(t, ok)
				assert.Nil(t, obj)
			},
			args: args{
				ctx:   context.Background(),
				index: testIndexName,
				key:   []string{"foo"},
			},
			wantErr: false,
		},
		{
			name:   "not found",
			config: cache.Config{},
			preparation: func(t *testing.T, c cache.Cache[testIndex, string, *testObject], s *miniredis.Miniredis) {
				c.Set(context.Background(), &testObject{
					ID:   "one",
					Name: []string{"foo", "bar"},
				})
			},
			assertions: func(t *testing.T, c cache.Cache[testIndex, string, *testObject]) {
				obj, ok := c.Get(context.Background(), testIndexID, "one")
				assert.True(t, ok)
				assert.NotNil(t, obj)
				obj, ok = c.Get(context.Background(), testIndexName, "foo")
				assert.True(t, ok)
				assert.NotNil(t, obj)
			},
			args: args{
				ctx:   context.Background(),
				index: testIndexName,
				key:   []string{"spanac"},
			},
			wantErr: false,
		},
		{
			name:   "invalidate multiple keys, first not found",
			config: cache.Config{},
			preparation: func(t *testing.T, c cache.Cache[testIndex, string, *testObject], s *miniredis.Miniredis) {
				c.Set(context.Background(), &testObject{
					ID:   "one",
					Name: []string{"foo", "bar"},
				})
			},
			assertions: func(t *testing.T, c cache.Cache[testIndex, string, *testObject]) {
				obj, ok := c.Get(context.Background(), testIndexID, "one")
				assert.False(t, ok)
				assert.Nil(t, obj)
			},
			args: args{
				ctx:   context.Background(),
				index: testIndexID,
				key:   []string{"spanac", "one"},
			},
			wantErr: false,
		},
		{
			name: "invalidate multiple keys, first expired",
			config: cache.Config{
				LastUseAge: time.Minute,
			},
			preparation: func(t *testing.T, c cache.Cache[testIndex, string, *testObject], s *miniredis.Miniredis) {
				c.Set(context.Background(), &testObject{
					ID:   "one",
					Name: []string{"foo"},
				})
				s.FastForward(40 * time.Second)
				c.Set(context.Background(), &testObject{
					ID:   "two",
					Name: []string{"bar"},
				})
				s.FastForward(30 * time.Second)
			},
			assertions: func(t *testing.T, c cache.Cache[testIndex, string, *testObject]) {
				obj, ok := c.Get(context.Background(), testIndexID, "two")
				assert.False(t, ok)
				assert.Nil(t, obj)
			},
			args: args{
				ctx:   context.Background(),
				index: testIndexID,
				key:   []string{"one", "two"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, server := prepareCache(t, tt.config)
			tt.preparation(t, c, server)
			t.Log(server.Keys())

			err := c.Invalidate(tt.args.ctx, tt.args.index, tt.args.key...)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			if tt.assertions != nil {
				tt.assertions(t, c)
			}
		})
	}
}

func Test_redisCache_Delete(t *testing.T) {
	type args struct {
		ctx   context.Context
		index testIndex
		key   []string
	}
	tests := []struct {
		name        string
		config      cache.Config
		preparation func(t *testing.T, c cache.Cache[testIndex, string, *testObject], s *miniredis.Miniredis)
		assertions  func(t *testing.T, c cache.Cache[testIndex, string, *testObject])
		args        args
		wantErr     bool
	}{
		{
			name:   "connection error",
			config: cache.Config{},
			preparation: func(_ *testing.T, _ cache.Cache[testIndex, string, *testObject], s *miniredis.Miniredis) {
				s.RequireAuth("foobar")
			},
			args: args{
				ctx:   context.Background(),
				index: testIndexName,
				key:   []string{"foo"},
			},
			wantErr: true,
		},
		{
			name:   "no keys, noop",
			config: cache.Config{},
			preparation: func(t *testing.T, c cache.Cache[testIndex, string, *testObject], s *miniredis.Miniredis) {
				c.Set(context.Background(), &testObject{
					ID:   "one",
					Name: []string{"foo", "bar"},
				})
			},
			args: args{
				ctx:   context.Background(),
				index: testIndexID,
				key:   []string{},
			},
			wantErr: false,
		},
		{
			name:   "delete ID",
			config: cache.Config{},
			preparation: func(t *testing.T, c cache.Cache[testIndex, string, *testObject], s *miniredis.Miniredis) {
				c.Set(context.Background(), &testObject{
					ID:   "one",
					Name: []string{"foo", "bar"},
				})
			},
			assertions: func(t *testing.T, c cache.Cache[testIndex, string, *testObject]) {
				obj, ok := c.Get(context.Background(), testIndexID, "one")
				assert.False(t, ok)
				assert.Nil(t, obj)
				// Get be name should still work
				obj, ok = c.Get(context.Background(), testIndexName, "foo")
				assert.True(t, ok)
				assert.NotNil(t, obj)
			},
			args: args{
				ctx:   context.Background(),
				index: testIndexID,
				key:   []string{"one"},
			},
			wantErr: false,
		},
		{
			name:   "delete name",
			config: cache.Config{},
			preparation: func(t *testing.T, c cache.Cache[testIndex, string, *testObject], s *miniredis.Miniredis) {
				c.Set(context.Background(), &testObject{
					ID:   "one",
					Name: []string{"foo", "bar"},
				})
			},
			assertions: func(t *testing.T, c cache.Cache[testIndex, string, *testObject]) {
				// get by ID should still work
				obj, ok := c.Get(context.Background(), testIndexID, "one")
				assert.True(t, ok)
				assert.NotNil(t, obj)
				obj, ok = c.Get(context.Background(), testIndexName, "foo")
				assert.False(t, ok)
				assert.Nil(t, obj)
			},
			args: args{
				ctx:   context.Background(),
				index: testIndexName,
				key:   []string{"foo"},
			},
			wantErr: false,
		},
		{
			name:   "not found",
			config: cache.Config{},
			preparation: func(t *testing.T, c cache.Cache[testIndex, string, *testObject], s *miniredis.Miniredis) {
				c.Set(context.Background(), &testObject{
					ID:   "one",
					Name: []string{"foo", "bar"},
				})
			},
			assertions: func(t *testing.T, c cache.Cache[testIndex, string, *testObject]) {
				obj, ok := c.Get(context.Background(), testIndexID, "one")
				assert.True(t, ok)
				assert.NotNil(t, obj)
				obj, ok = c.Get(context.Background(), testIndexName, "foo")
				assert.True(t, ok)
				assert.NotNil(t, obj)
			},
			args: args{
				ctx:   context.Background(),
				index: testIndexName,
				key:   []string{"spanac"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, server := prepareCache(t, tt.config)
			tt.preparation(t, c, server)
			t.Log(server.Keys())

			err := c.Delete(tt.args.ctx, tt.args.index, tt.args.key...)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			if tt.assertions != nil {
				tt.assertions(t, c)
			}
		})
	}
}

func Test_redisCache_Truncate(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name        string
		config      cache.Config
		preparation func(t *testing.T, c cache.Cache[testIndex, string, *testObject], s *miniredis.Miniredis)
		assertions  func(t *testing.T, c cache.Cache[testIndex, string, *testObject])
		args        args
		wantErr     bool
	}{
		{
			name:   "connection error",
			config: cache.Config{},
			preparation: func(_ *testing.T, _ cache.Cache[testIndex, string, *testObject], s *miniredis.Miniredis) {
				s.RequireAuth("foobar")
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
		},
		{
			name:   "ok",
			config: cache.Config{},
			preparation: func(t *testing.T, c cache.Cache[testIndex, string, *testObject], s *miniredis.Miniredis) {
				c.Set(context.Background(), &testObject{
					ID:   "one",
					Name: []string{"foo", "bar"},
				})
				c.Set(context.Background(), &testObject{
					ID:   "two",
					Name: []string{"Hello", "World"},
				})
			},
			assertions: func(t *testing.T, c cache.Cache[testIndex, string, *testObject]) {
				obj, ok := c.Get(context.Background(), testIndexID, "one")
				assert.False(t, ok)
				assert.Nil(t, obj)
				obj, ok = c.Get(context.Background(), testIndexName, "World")
				assert.False(t, ok)
				assert.Nil(t, obj)
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, server := prepareCache(t, tt.config)
			tt.preparation(t, c, server)
			t.Log(server.Keys())

			err := c.Truncate(tt.args.ctx)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			if tt.assertions != nil {
				tt.assertions(t, c)
			}
		})
	}
}

// assertKeysTTL asserts that the index keys and the key set of the object exist with the given TTL.
func assertKeysTTL(t *testing.T, s *miniredis.Miniredis, objectID string, ttl time.Duration) {
	t.Helper()
	for _, key := range []string{"VERSION:0:one", "VERSION:1:foo", "VERSION:1:bar", objectID + "-keys"} {
		assert.True(t, s.Exists(key), key)
		assert.Equal(t, ttl, s.TTL(key), key)
	}
}

func prepareCache(t *testing.T, conf cache.Config, options ...func(*Config)) (cache.Cache[testIndex, string, *testObject], *miniredis.Miniredis) {
	conf.Log = &logging.Config{
		Level:     "debug",
		AddSource: true,
	}
	server := miniredis.RunT(t)
	server.Select(testDB)

	connConfig := Config{
		Enabled:          true,
		Network:          "tcp",
		Addr:             server.Addr(),
		DisableIndentity: true,
	}
	for _, option := range options {
		option(&connConfig)
	}
	connector, err := NewConnector(connConfig)
	require.NoError(t, err)
	t.Cleanup(func() {
		connector.Close()
		server.Close()
	})
	c := NewCache[testIndex, string, *testObject](conf, "VERSION", connector, testDB, testIndices)
	return c, server
}

func withCircuitBreakerOption(cb *CBConfig) func(*Config) {
	return func(c *Config) {
		c.CircuitBreaker = cb
	}
}
