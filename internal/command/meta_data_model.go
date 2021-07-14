package command

import (
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/metadata"
)

type MetaDataWriteModel struct {
	eventstore.WriteModel

	Key   string
	Value string
	State domain.MetaDataState
}

func (wm *MetaDataWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *metadata.SetEvent:
			if wm.Key != e.Key {
				continue
			}
			wm.Value = e.Value
			wm.State = domain.MetaDataStateActive
		case *metadata.RemovedEvent:
			wm.State = domain.MetaDataStateRemoved
		}
	}
	return wm.WriteModel.Reduce()
}

type MetaDataListWriteModel struct {
	eventstore.WriteModel

	metaDataList map[string]string
}

func (wm *MetaDataListWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *metadata.SetEvent:
			wm.metaDataList[e.Key] = e.Value
		case *metadata.RemovedEvent:
			delete(wm.metaDataList, e.Key)
		}
	}
	return wm.WriteModel.Reduce()
}
