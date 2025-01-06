package consistent

import (
	"context"
	"errors"
	"reflect"
	"slices"
	"strings"
	"sync"

	"github.com/zitadel/zitadel/backend/internal/port"
	"github.com/zitadel/zitadel/backend/internal/port/storage"
)

var (
	_ port.Executor[port.Object] = (*MapRepository[port.Object])(nil)
	_ port.Querier[port.Object]  = (*MapRepository[port.Object])(nil)
	_ port.Client[port.Object]   = (*MapRepository[port.Object])(nil)
	// _ port.Pool[any]        = (*MapRepository[any])(nil)

	_ port.Executor[port.Object]    = (*MapTransaction[port.Object])(nil)
	_ port.Querier[port.Object]     = (*MapTransaction[port.Object])(nil)
	_ port.Transaction[port.Object] = (*MapTransaction[port.Object])(nil)
)
var pkFields = &sync.Map{}

type MapRepository[T port.Object] struct {
	objects []*row[T]
}

type row[T port.Object] struct {
	object *T
	mu     *sync.RWMutex
}

func (m *MapRepository[T]) Begin(ctx context.Context) (port.Transaction[T], error) {
	return &MapTransaction[T]{
		repo: m,
	}, nil
}

func (m *MapRepository[T]) Release(ctx context.Context) error {
	return nil
}

// Create implements [port.Executor]
func (m *MapRepository[T]) Create(ctx context.Context, object T) (err error) {
	tx, err := m.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.End(ctx, err)

	return tx.Create(ctx, object)
}

// Delete implements [port.Executor]
func (m *MapRepository[T]) Delete(ctx context.Context, filters []port.Filter) (err error) {
	tx, err := m.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.End(ctx, err)

	return tx.Delete(ctx, filters)
}

// Update implements [port.Executor]
func (m *MapRepository[T]) Update(ctx context.Context, columns []*port.Column, filters []port.Filter) (err error) {
	tx, err := m.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.End(ctx, err)

	return tx.Update(ctx, columns, filters)
}

// Get implements [port.Querier]
func (m *MapRepository[T]) Get(ctx context.Context, filters []port.Filter) (o T, err error) {
	tx, err := m.Begin(ctx)
	if err != nil {
		return o, err
	}
	defer tx.End(ctx, err)

	return tx.Get(ctx, filters)
}

// List implements [port.Querier]
func (m *MapRepository[T]) List(ctx context.Context, filters []port.Filter) (_ []T, err error) {
	tx, err := m.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.End(ctx, err)

	return tx.List(ctx, filters)
}

type MapTransaction[T port.Object] struct {
	repo    *MapRepository[T]
	changes []*change[T]
}

type change[T port.Object] struct {
	row     *row[T]
	changed *T
	typ     changeType
}

type changeType uint8

const (
	createChangeType changeType = iota
	updateChangeType
	deleteChangeType
)

func (m *MapTransaction[T]) End(ctx context.Context, gotErr error) error {
	if ctx.Err() == nil && gotErr == nil {
		return m.commit()
	}
	return m.rollback()
}

func (m *MapTransaction[T]) commit() error {
	for _, change := range m.changes {
		switch change.typ {
		case createChangeType:
			change.row.object = change.changed
		case deleteChangeType:
			m.repo.objects = slices.DeleteFunc(m.repo.objects, func(r *row[T]) bool {
				return r == change.row
			})
			change.row.mu.Unlock()
			change.row = nil
			continue
		case updateChangeType:
			*change.row.object = *change.changed
		default:
			return errors.New("unknown change type")
		}
		change.row.mu.Unlock()
	}
	return nil
}

func (m *MapTransaction[T]) rollback() error {
	for _, change := range m.changes {
		if change.changed == nil {
			// TODO: remove from repo.objects
		}
		change.row.mu.Unlock()
		change = nil
	}
	m.changes = nil
	return nil
}

// Create implements port.Executor.
func (m *MapTransaction[T]) Create(ctx context.Context, object T) error {
	_, err := m.get(ctx, pkFilter(object))
	if err == nil {
		return storage.ErrAlreadyExists
	}

	mu := &sync.RWMutex{}
	mu.Lock()
	m.repo.objects = append(m.repo.objects, &row[T]{
		object: &object,
		mu:     mu,
	})
	m.changes = append(m.changes, &change[T]{
		row: m.repo.objects[len(m.repo.objects)-1],
		typ: createChangeType,
	})
	return nil
}

// Delete implements port.Executor.
func (m *MapTransaction[T]) Delete(ctx context.Context, filters []port.Filter) error {
	rows, err := m.list(ctx, filters)
	if err != nil {
		return err
	}

	for _, row := range rows {
		row.mu.Lock()
		m.changes = append(m.changes, &change[T]{
			row: row,
			typ: deleteChangeType,
		})
	}

	return nil
}

// Update implements port.Executor.
func (m *MapTransaction[T]) Update(ctx context.Context, columns []*port.Column, filters []port.Filter) error {
	rows, err := m.list(ctx, filters)
	if err != nil {
		return err
	}
	for _, row := range rows {
		row.mu.Lock()
		changed, err := update(row, columns)
		if err != nil {
			return err
		}
		m.changes = append(m.changes, &change[T]{
			row:     row,
			changed: changed,
			typ:     updateChangeType,
		})
	}
	return nil
}

// Get implements [port.Querier]
func (m *MapTransaction[T]) Get(ctx context.Context, filters []port.Filter) (o T, err error) {
	r, err := m.get(ctx, filters)
	if err != nil {
		return o, err
	}
	r.mu.RLock()

	return o, nil
}

func (m *MapTransaction[T]) get(ctx context.Context, filters []port.Filter) (o *row[T], err error) {
	rows, err := m.list(ctx, filters)
	if err != nil {
		return nil, err
	}
	if len(rows) > 1 {
		return nil, errors.New("multiple rows returned")
	}
	return rows[0], nil
}

// List implements [port.Querier]
func (m *MapTransaction[T]) List(ctx context.Context, filters []port.Filter) ([]T, error) {
	rows, err := m.list(ctx, filters)
	if err != nil {
		return nil, err
	}
	objects := make([]T, len(rows))
	for i, row := range rows {
		row.mu.RLock()
		objects[i] = *row.object
	}
	return objects, nil
}

func (m *MapTransaction[T]) list(ctx context.Context, filters []port.Filter) (res []*row[T], err error) {
	res = slices.Clone(m.repo.objects)
	for _, filter := range filters {
		res = slices.DeleteFunc(res, func(r *row[T]) bool {
			v := reflect.ValueOf(r.object)
			definition := getDefinition[T]()
			return v.Field(definition.columnIndexes[filter.Column().Name]).Interface() != filter.Column().Value
		})
		res = slices.Clip(res)
	}
	return res, nil
}

func pkFilter[T port.Object](o T) []port.Filter {
	definition := getDefinition[T]()

	v := reflect.ValueOf(o)
	filters := make([]port.Filter, len(definition.pkFieldIndexes))

	for _, idx := range definition.pkFieldIndexes {
		filters = append(filters, port.NewEqualFilter(
			&port.Column{
				Name:  definition.columnNames[idx],
				Value: v.Field(idx).Interface(),
			},
		))
	}

	return filters
}

type structDefinition struct {
	columnNames    []string
	columnIndexes  map[string]int
	pkFieldIndexes []int
	fkFieldIndexes []int
}

func getDefinition[T port.Object]() *structDefinition {
	definition, ok := pkFields.Load(reflect.TypeFor[T]())
	if !ok {
		definition = registerPkFields[T]()
	}
	return definition.(*structDefinition)
}

func registerPkFields[T port.Object]() *structDefinition {
	t := reflect.TypeFor[T]()

	definition := &structDefinition{
		columnIndexes:  make(map[string]int, t.NumField()),
		columnNames:    make([]string, t.NumField()),
		pkFieldIndexes: make([]int, 0, t.NumField()),
		fkFieldIndexes: make([]int, 0, t.NumField()),
	}

	for i := 0; i < t.NumField(); i++ {
		tags := t.Field(i).Tag.Get("consistent")
		fields := strings.Split(tags, ",")
		if fields[0] == "" {
			fields[0] = t.Field(i).Name
		}
		definition.columnIndexes[fields[0]] = i
		definition.columnNames[i] = fields[0]

		for _, field := range fields[1:] {
			switch field {
			case "pk":
				definition.pkFieldIndexes = append(definition.pkFieldIndexes, i)
			case "fk":
				definition.fkFieldIndexes = append(definition.fkFieldIndexes, i)
			}
		}
	}
	pkFields.Store(t, definition)
	return definition
}

func update[T port.Object](r *row[T], columns []*port.Column) (o *T, err error) {
	v := reflect.Zero(reflect.TypeOf(r.object))
	current := reflect.ValueOf(r.object)
	definition := getDefinition[T]()

fields:
	for i := 0; i < v.NumField(); i++ {
		for _, column := range columns {
			if definition.columnIndexes[column.Name] != i {
				continue
			}
			v.Field(i).Set(reflect.ValueOf(column.Value))
			continue fields
		}
		v.Field(i).Set(current.Field(i))
	}

	return r.object, nil
}
