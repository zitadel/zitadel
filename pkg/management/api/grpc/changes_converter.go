package grpc

import (
	chg_model "github.com/caos/zitadel/internal/changes/model"
)

func changesToResponse(response *chg_model.Changes, offset uint64, limit uint64) (_ *Changes) {
	return &Changes{
		Limit:   limit,
		Offset:  offset,
		Changes: changesToMgtAPI(response),
	}
}

func changesToMgtAPI(changes *chg_model.Changes) (_ []*Change) {
	result := make([]*Change, len(changes.Changes))

	for i, change := range changes.Changes {
		result[i] = &Change{
			ChangeDate: change.ChangeDate,
			EventType:  change.EventType,
			Sequence:   change.Sequence,
		}
	}

	return result
}
