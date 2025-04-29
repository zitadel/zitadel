package domain

import (
	"context"
	"time"
)

type Instance struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	DeletedAt time.Time `json:"-"`
}

// Keys implements the [cache.Entry].
func (i *Instance) Keys(index string) (key []string) {
	// TODO: Return the correct keys for the instance cache, e.g., i.ID, i.Domain
	return []string{}
}

type InstanceRepository interface {
	ByID(ctx context.Context, id string) (*Instance, error)
	Create(ctx context.Context, instance *Instance) error
	On(id string) InstanceOperation
}

type InstanceOperation interface {
	AdminRepository
	Update(ctx context.Context, instance *Instance) error
	Delete(ctx context.Context) error
}

type CreateInstance struct {
	Name string `json:"name"`
}
