package command

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

type InstanceSecretGeneratorConfigWriteModel struct {
	eventstore.WriteModel

	GeneratorType       domain.SecretGeneratorType
	Length              uint
	Expiry              time.Duration
	IncludeLowerLetters bool
	IncludeUpperLetters bool
	IncludeDigits       bool
	IncludeSymbols      bool
	State               domain.SecretGeneratorState
}

func NewInstanceSecretGeneratorConfigWriteModel(ctx context.Context, GeneratorType domain.SecretGeneratorType) *InstanceSecretGeneratorConfigWriteModel {
	return &InstanceSecretGeneratorConfigWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   authz.GetInstance(ctx).InstanceID(),
			ResourceOwner: authz.GetInstance(ctx).InstanceID(),
			InstanceID:    authz.GetInstance(ctx).InstanceID(),
		},
		GeneratorType: GeneratorType,
	}
}

func (wm *InstanceSecretGeneratorConfigWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *instance.SecretGeneratorAddedEvent:
			if wm.GeneratorType != e.GeneratorType {
				continue
			}
			wm.Length = e.Length
			wm.Expiry = e.Expiry
			wm.IncludeLowerLetters = e.IncludeLowerLetters
			wm.IncludeUpperLetters = e.IncludeUpperLetters
			wm.IncludeDigits = e.IncludeDigits
			wm.IncludeSymbols = e.IncludeSymbols
			wm.State = domain.SecretGeneratorStateActive
		case *instance.SecretGeneratorChangedEvent:
			if wm.GeneratorType != e.GeneratorType {
				continue
			}
			if e.Length != nil {
				wm.Length = *e.Length
			}
			if e.Expiry != nil {
				wm.Expiry = *e.Expiry
			}
			if e.IncludeUpperLetters != nil {
				wm.IncludeUpperLetters = *e.IncludeUpperLetters
			}
			if e.IncludeLowerLetters != nil {
				wm.IncludeLowerLetters = *e.IncludeLowerLetters
			}
			if e.IncludeDigits != nil {
				wm.IncludeDigits = *e.IncludeDigits
			}
			if e.IncludeSymbols != nil {
				wm.IncludeSymbols = *e.IncludeSymbols
			}
		case *instance.SecretGeneratorRemovedEvent:
			if wm.GeneratorType != e.GeneratorType {
				continue
			}
			wm.State = domain.SecretGeneratorStateRemoved
			wm.Length = 0
			wm.Expiry = 0
			wm.IncludeLowerLetters = false
			wm.IncludeUpperLetters = false
			wm.IncludeDigits = false
			wm.IncludeSymbols = false
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *InstanceSecretGeneratorConfigWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(instance.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			instance.SecretGeneratorAddedEventType,
			instance.SecretGeneratorChangedEventType,
			instance.SecretGeneratorRemovedEventType).
		Builder()
}

func (wm *InstanceSecretGeneratorConfigWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	generatorType domain.SecretGeneratorType,
	length uint,
	expiry time.Duration,
	includeLowerLetters,
	includeUpperLetters,
	includeDigits,
	includeSymbols bool,
) (*instance.SecretGeneratorChangedEvent, bool, error) {
	changes := make([]instance.SecretGeneratorChanges, 0)
	var err error

	if wm.Length != length {
		changes = append(changes, instance.ChangeSecretGeneratorLength(length))
	}
	if wm.Expiry != expiry {
		changes = append(changes, instance.ChangeSecretGeneratorExpiry(expiry))
	}
	if wm.IncludeLowerLetters != includeLowerLetters {
		changes = append(changes, instance.ChangeSecretGeneratorIncludeLowerLetters(includeLowerLetters))
	}
	if wm.IncludeUpperLetters != includeUpperLetters {
		changes = append(changes, instance.ChangeSecretGeneratorIncludeUpperLetters(includeUpperLetters))
	}
	if wm.IncludeDigits != includeDigits {
		changes = append(changes, instance.ChangeSecretGeneratorIncludeDigits(includeDigits))
	}
	if wm.IncludeSymbols != includeSymbols {
		changes = append(changes, instance.ChangeSecretGeneratorIncludeSymbols(includeSymbols))
	}
	if len(changes) == 0 {
		return nil, false, nil
	}
	changeEvent, err := instance.NewSecretGeneratorChangeEvent(ctx, aggregate, generatorType, changes)
	if err != nil {
		return nil, false, err
	}
	return changeEvent, true, nil
}
