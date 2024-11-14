package gomap

import (
	"context"
	"testing"
	"time"

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

var testIndices = []testIndex{
	testIndexID,
	testIndexName,
}

type testObject struct {
	id    string
	names []string
}

func (o *testObject) Keys(index testIndex) []string {
	switch index {
	case testIndexID:
		return []string{o.id}
	case testIndexName:
		return o.names
	default:
		return nil
	}
}

func Test_mapCache_Get(t *testing.T) {
	c := NewCache[testIndex, string, *testObject](context.Background(), testIndices, cache.Config{
		MaxAge:     time.Second,
		LastUseAge: time.Second / 4,
		Log: &logging.Config{
			Level:     "debug",
			AddSource: true,
		},
	})
	obj := &testObject{
		id:    "id",
		names: []string{"foo", "bar"},
	}
	c.Set(context.Background(), obj)

	type args struct {
		index testIndex
		key   string
	}
	tests := []struct {
		name   string
		args   args
		want   *testObject
		wantOk bool
	}{
		{
			name: "ok",
			args: args{
				index: testIndexID,
				key:   "id",
			},
			want:   obj,
			wantOk: true,
		},
		{
			name: "miss",
			args: args{
				index: testIndexID,
				key:   "spanac",
			},
			want:   nil,
			wantOk: false,
		},
		{
			name: "unknown index",
			args: args{
				index: 99,
				key:   "id",
			},
			want:   nil,
			wantOk: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := c.Get(context.Background(), tt.args.index, tt.args.key)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantOk, ok)
		})
	}
}

func Test_mapCache_Invalidate(t *testing.T) {
	c := NewCache[testIndex, string, *testObject](context.Background(), testIndices, cache.Config{
		MaxAge:     time.Second,
		LastUseAge: time.Second / 4,
		Log: &logging.Config{
			Level:     "debug",
			AddSource: true,
		},
	})
	obj := &testObject{
		id:    "id",
		names: []string{"foo", "bar"},
	}
	c.Set(context.Background(), obj)
	err := c.Invalidate(context.Background(), testIndexName, "bar")
	require.NoError(t, err)
	got, ok := c.Get(context.Background(), testIndexID, "id")
	assert.Nil(t, got)
	assert.False(t, ok)
}

func Test_mapCache_Delete(t *testing.T) {
	c := NewCache[testIndex, string, *testObject](context.Background(), testIndices, cache.Config{
		MaxAge:     time.Second,
		LastUseAge: time.Second / 4,
		Log: &logging.Config{
			Level:     "debug",
			AddSource: true,
		},
	})
	obj := &testObject{
		id:    "id",
		names: []string{"foo", "bar"},
	}
	c.Set(context.Background(), obj)
	err := c.Delete(context.Background(), testIndexName, "bar")
	require.NoError(t, err)

	// Shouldn't find object by deleted name
	got, ok := c.Get(context.Background(), testIndexName, "bar")
	assert.Nil(t, got)
	assert.False(t, ok)

	// Should find object by other name
	got, ok = c.Get(context.Background(), testIndexName, "foo")
	assert.Equal(t, obj, got)
	assert.True(t, ok)

	// Should find object by id
	got, ok = c.Get(context.Background(), testIndexID, "id")
	assert.Equal(t, obj, got)
	assert.True(t, ok)
}

func Test_mapCache_Prune(t *testing.T) {
	c := NewCache[testIndex, string, *testObject](context.Background(), testIndices, cache.Config{
		MaxAge:     time.Second,
		LastUseAge: time.Second / 4,
		Log: &logging.Config{
			Level:     "debug",
			AddSource: true,
		},
	})

	objects := []*testObject{
		{
			id:    "id1",
			names: []string{"foo", "bar"},
		},
		{
			id:    "id2",
			names: []string{"hello"},
		},
	}
	for _, obj := range objects {
		c.Set(context.Background(), obj)
	}
	// invalidate one entry
	err := c.Invalidate(context.Background(), testIndexName, "bar")
	require.NoError(t, err)

	err = c.(cache.Pruner).Prune(context.Background())
	require.NoError(t, err)

	// Other object should still be found
	got, ok := c.Get(context.Background(), testIndexID, "id2")
	assert.Equal(t, objects[1], got)
	assert.True(t, ok)
}

func Test_mapCache_Truncate(t *testing.T) {
	c := NewCache[testIndex, string, *testObject](context.Background(), testIndices, cache.Config{
		MaxAge:     time.Second,
		LastUseAge: time.Second / 4,
		Log: &logging.Config{
			Level:     "debug",
			AddSource: true,
		},
	})
	objects := []*testObject{
		{
			id:    "id1",
			names: []string{"foo", "bar"},
		},
		{
			id:    "id2",
			names: []string{"hello"},
		},
	}
	for _, obj := range objects {
		c.Set(context.Background(), obj)
	}

	err := c.Truncate(context.Background())
	require.NoError(t, err)

	mc := c.(*mapCache[testIndex, string, *testObject])
	for _, index := range mc.indexMap {
		index.mutex.RLock()
		assert.Len(t, index.entries, 0)
		index.mutex.RUnlock()
	}
}

func Test_entry_isValid(t *testing.T) {
	type fields struct {
		created time.Time
		invalid bool
		lastUse time.Time
	}
	tests := []struct {
		name   string
		fields fields
		config *cache.Config
		want   bool
	}{
		{
			name: "invalid",
			fields: fields{
				created: time.Now(),
				invalid: true,
				lastUse: time.Now(),
			},
			config: &cache.Config{
				MaxAge:     time.Minute,
				LastUseAge: time.Second,
			},
			want: false,
		},
		{
			name: "max age exceeded",
			fields: fields{
				created: time.Now().Add(-(time.Minute + time.Second)),
				invalid: false,
				lastUse: time.Now(),
			},
			config: &cache.Config{
				MaxAge:     time.Minute,
				LastUseAge: time.Second,
			},
			want: false,
		},
		{
			name: "max age disabled",
			fields: fields{
				created: time.Now().Add(-(time.Minute + time.Second)),
				invalid: false,
				lastUse: time.Now(),
			},
			config: &cache.Config{
				LastUseAge: time.Second,
			},
			want: true,
		},
		{
			name: "last use age exceeded",
			fields: fields{
				created: time.Now().Add(-(time.Minute / 2)),
				invalid: false,
				lastUse: time.Now().Add(-(time.Second * 2)),
			},
			config: &cache.Config{
				MaxAge:     time.Minute,
				LastUseAge: time.Second,
			},
			want: false,
		},
		{
			name: "last use age disabled",
			fields: fields{
				created: time.Now().Add(-(time.Minute / 2)),
				invalid: false,
				lastUse: time.Now().Add(-(time.Second * 2)),
			},
			config: &cache.Config{
				MaxAge: time.Minute,
			},
			want: true,
		},
		{
			name: "valid",
			fields: fields{
				created: time.Now(),
				invalid: false,
				lastUse: time.Now(),
			},
			config: &cache.Config{
				MaxAge:     time.Minute,
				LastUseAge: time.Second,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &entry[any]{
				created: tt.fields.created,
			}
			e.invalid.Store(tt.fields.invalid)
			e.lastUse.Store(tt.fields.lastUse.UnixMicro())
			got := e.isValid(tt.config)
			assert.Equal(t, tt.want, got)
		})
	}
}
