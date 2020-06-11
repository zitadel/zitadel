package grpc

import (
	"encoding/json"

	"github.com/caos/logging"
	org_model "github.com/caos/zitadel/internal/org/model"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
)

func orgsFromModel(orgs []*org_model.Org) []*Org {
	orgList := make([]*Org, len(orgs))
	for i, org := range orgs {
		orgList[i] = orgFromModel(org)
	}
	return orgList
}

func orgFromModel(org *org_model.Org) *Org {
	creationDate, err := ptypes.TimestampProto(org.CreationDate)
	logging.Log("GRPC-GTHsZ").OnError(err).Debug("unable to get timestamp from time")

	changeDate, err := ptypes.TimestampProto(org.ChangeDate)
	logging.Log("GRPC-dVnoj").OnError(err).Debug("unable to get timestamp from time")

	return &Org{
		Domain:       org.Domain,
		ChangeDate:   changeDate,
		CreationDate: creationDate,
		Id:           org.AggregateID,
		Name:         org.Name,
		State:        orgStateFromModel(org.State),
	}
}

func orgFromView(org *org_model.OrgView) *Org {
	creationDate, err := ptypes.TimestampProto(org.CreationDate)
	logging.Log("GRPC-GTHsZ").OnError(err).Debug("unable to get timestamp from time")

	changeDate, err := ptypes.TimestampProto(org.ChangeDate)
	logging.Log("GRPC-dVnoj").OnError(err).Debug("unable to get timestamp from time")

	return &Org{
		Domain:       org.Domain,
		ChangeDate:   changeDate,
		CreationDate: creationDate,
		Id:           org.ID,
		Name:         org.Name,
		State:        orgStateFromModel(org.State),
	}
}

func orgStateFromModel(state org_model.OrgState) OrgState {
	switch state {
	case org_model.ORGSTATE_ACTIVE:
		return OrgState_ORGSTATE_ACTIVE
	case org_model.ORGSTATE_INACTIVE:
		return OrgState_ORGSTATE_INACTIVE
	default:
		return OrgState_ORGSTATE_UNSPECIFIED
	}
}

func orgChangesToResponse(response *org_model.OrgChanges, offset uint64, limit uint64) (_ *Changes) {
	return &Changes{
		Limit:   limit,
		Offset:  offset,
		Changes: orgChangesToMgtAPI(response),
	}
}

func orgChangesToMgtAPI(changes *org_model.OrgChanges) (_ []*Change) {
	result := make([]*Change, len(changes.Changes))

	for i, change := range changes.Changes {
		b, err := json.Marshal(change.Data)
		data := &structpb.Struct{}
		err = protojson.Unmarshal(b, data)
		if err != nil {
		}
		result[i] = &Change{
			ChangeDate: change.ChangeDate,
			EventType:  change.EventType,
			Sequence:   change.Sequence,
			Data:       data,
		}
	}

	return result
}
