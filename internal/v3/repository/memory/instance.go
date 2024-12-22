package memory

import (
	"context"
	"errors"
	"sync"
	"time"

	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/v3/instance"
	"github.com/zitadel/zitadel/internal/v3/storage"
)

var _ instance.InstanceStorage = (*InstanceMemory)(nil)

type InstanceMemory struct {
	instances map[string]*memoryInstance
	mu        sync.RWMutex
}

func NewInstanceMemory() *InstanceMemory {
	return &InstanceMemory{
		instances: make(map[string]*memoryInstance),
	}
}

// WriteInstanceAdded implements instance.InstanceStorage.
func (i *InstanceMemory) WriteInstanceAdded(ctx context.Context, tx storage.Transaction, instance *instance.AddInstanceRequest) error {
	defaultLanguage, err := language.Parse(instance.DefaultLanguage)
	if err != nil {
		return err
	}

	if instance.CreatedAt.IsZero() {
		instance.CreatedAt = time.Now()
	}

	i.mu.Lock()

	if i.instances[instance.ID] != nil {
		return errors.New("instance already exists")
	}

	i.instances[instance.ID] = &memoryInstance{
		id:              instance.ID,
		name:            instance.InstanceName,
		customDomain:    instance.CustomDomain,
		defaultLanguage: defaultLanguage,
	}

	tx.OnCommit(func(ctx context.Context) error {
		i.mu.Unlock()
		return nil
	})

	tx.OnRollback(func(ctx context.Context) error {
		delete(i.instances, instance.ID)
		i.mu.Unlock()
		return nil
	})

	return nil
}

type memoryInstance struct {
	id              string
	name            string
	customDomain    string
	defaultLanguage language.Tag
}
