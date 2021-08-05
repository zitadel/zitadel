package command

import (
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/metadata"
)

type MetadataWriteModel struct {
	eventstore.WriteModel

	Key   string
	Value []byte
	State domain.MetadataState
}

func (wm *MetadataWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *metadata.SetEvent:
			if wm.Key != e.Key {
				continue
			}
			wm.Value = e.Value
			wm.State = domain.MetadataStateActive
		case *metadata.RemovedEvent:
			wm.State = domain.MetadataStateRemoved
		}
	}
	return wm.WriteModel.Reduce()
}

type MetadataListWriteModel struct {
	eventstore.WriteModel

	metaDataList map[string][]byte
}

func (wm *MetadataListWriteModel) Reduce() error {
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
