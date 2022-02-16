package command

import (
	"context"
	"time"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/iam"
)

type IAMSecretGeneratorConfigWriteModel struct {
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

func NewIAMSecretGeneratorConfigWriteModel(GeneratorType domain.SecretGeneratorType) *IAMSecretGeneratorConfigWriteModel {
	return &IAMSecretGeneratorConfigWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   domain.IAMID,
			ResourceOwner: domain.IAMID,
		},
		GeneratorType: GeneratorType,
	}
}

func (wm *IAMSecretGeneratorConfigWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *iam.SecretGeneratorAddedEvent:
			if wm.GeneratorType != e.GeneratorType {
				continue
			}
			wm.Length = e.Length
			wm.Expiry = e.Expiry
			wm.IncludeLowerLetters = e.IncludeLowerLetters
			wm.IncludeUpperLetters = e.IncludeUpperLetters
			wm.IncludeDigits = e.IncludeDigits
			wm.IncludeSymbols = e.IncludeDigits
			wm.State = domain.SecretGeneratorStateActive
		case *iam.SecretGeneratorChangedEvent:
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
		case *iam.SecretGeneratorRemovedEvent:
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

func (wm *IAMSecretGeneratorConfigWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(iam.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			iam.SecretGeneratorAddedEventType,
			iam.SecretGeneratorChangedEventType,
			iam.SecretGeneratorRemovedEventType).
		Builder()
}

func (wm *IAMSecretGeneratorConfigWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	generatorType domain.SecretGeneratorType,
	length uint,
	expiry time.Duration,
	includeLowerLetters,
	includeUpperLetters,
	includeDigits,
	includeSymbols bool,
) (*iam.SecretGeneratorChangedEvent, bool, error) {
	changes := make([]iam.SecretGeneratorChanges, 0)
	var err error

	if wm.Length != length {
		changes = append(changes, iam.ChangeSecretGeneratorLength(length))
	}
	if wm.Expiry != expiry {
		changes = append(changes, iam.ChangeSecretGeneratorExpiry(expiry))
	}
	if wm.IncludeLowerLetters != includeLowerLetters {
		changes = append(changes, iam.ChangeSecretGeneratorIncludeLowerLetters(includeLowerLetters))
	}
	if wm.IncludeUpperLetters != includeUpperLetters {
		changes = append(changes, iam.ChangeSecretGeneratorIncludeUpperLetters(includeUpperLetters))
	}
	if wm.IncludeDigits != includeDigits {
		changes = append(changes, iam.ChangeSecretGeneratorIncludeDigits(includeDigits))
	}
	if wm.IncludeSymbols != includeSymbols {
		changes = append(changes, iam.ChangeSecretGeneratorIncludeSymbols(includeSymbols))
	}
	if len(changes) == 0 {
		return nil, false, nil
	}
	changeEvent, err := iam.NewSecretGeneratorChangeEvent(ctx, aggregate, generatorType, changes)
	if err != nil {
		return nil, false, err
	}
	return changeEvent, true, nil
}
