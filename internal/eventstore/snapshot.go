package eventstore

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type SnapshotType string

func NewSnapshotType[T any](object T) SnapshotType {
	return SnapshotType(fmt.Sprintf("%T", object))
}

type SnapshotBase struct {
	InstanceID   string
	SnapshotType SnapshotType
	AggregateID  string
	Position     float64
	ChangeDate   time.Time
}

type SnapshotData struct {
	SnapshotBase
	Payload []byte
}

type Snapshotter interface {
	SetSnapshot(ctx context.Context, snapshot *SnapshotData) error
	GetSnapshot(ctx context.Context, instanceID string, typ SnapshotType, aggregateID string) (*SnapshotData, error)
}

func (es *Eventstore) SetSnapshot(ctx context.Context, snapshot *SnapshotData) error {
	return es.snapshots.SetSnapshot(ctx, snapshot)
}

func (es *Eventstore) GetSnapshot(ctx context.Context, instanceID string, typ SnapshotType, aggregateID string) (*SnapshotData, error) {
	return es.snapshots.GetSnapshot(ctx, instanceID, typ, aggregateID)
}

type Snapshot[T any] struct {
	SnapshotBase
	Object T
}

// NewSnapshot returns an empty snapshot ready for Get.
func NewSnapshot[T any](aggregate *Aggregate) *Snapshot[T] {
	var object T
	return &Snapshot[T]{
		SnapshotBase: SnapshotBase{
			InstanceID:   aggregate.InstanceID,
			SnapshotType: NewSnapshotType(object),
			AggregateID:  aggregate.ID,
		},
	}
}

// SnapshotFromWriteModel returns a populated snapshot ready for Set.
func SnapshotFromWriteModel[T any](model *WriteModel, object T) *Snapshot[T] {
	return &Snapshot[T]{
		SnapshotBase: SnapshotBase{
			InstanceID:   model.InstanceID,
			SnapshotType: NewSnapshotType(object),
			AggregateID:  model.AggregateID,
			Position:     model.Position,
			ChangeDate:   model.ChangeDate,
		},
		Object: object,
	}
}

// SnapshotFromReadModel returns a populated snapshot ready for Set.
func SnapshotFromReadModel[T any](model *ReadModel, object T) *Snapshot[T] {
	return &Snapshot[T]{
		SnapshotBase: SnapshotBase{
			InstanceID:   model.InstanceID,
			SnapshotType: NewSnapshotType(object),
			AggregateID:  model.AggregateID,
			Position:     model.Position,
			ChangeDate:   model.ChangeDate,
		},
		Object: object,
	}
}

func (s *Snapshot[T]) Set(ctx context.Context, repo Snapshotter) (err error) {
	payload, err := json.Marshal(s.Object)
	if err != nil {
		return err
	}
	return repo.SetSnapshot(ctx, &SnapshotData{
		SnapshotBase: s.SnapshotBase,
		Payload:      payload,
	})
}

func (s *Snapshot[T]) Get(ctx context.Context, repo Snapshotter) (err error) {
	data, err := repo.GetSnapshot(ctx, s.InstanceID, NewSnapshotType(s.Object), s.AggregateID)
	if err != nil {
		return err
	}
	s.SnapshotBase = data.SnapshotBase
	if len(data.Payload) == 0 {
		return nil
	}
	if err = json.Unmarshal(data.Payload, &s.Object); err != nil {
		return err
	}
	return nil
}
