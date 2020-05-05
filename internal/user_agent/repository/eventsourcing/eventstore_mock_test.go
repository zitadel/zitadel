package eventsourcing

//
//import (
//	"encoding/json"
//	"github.com/caos/zitadel/internal/eventstore/mock"
//	es_models "github.com/caos/zitadel/internal/eventstore/models"
//	"github.com/caos/zitadel/internal/project/model"
//	"github.com/golang/mock/gomock"
//)
//
//func GetMockProjectByIDOK(ctrl *gomock.Controller) *ProjectEventstore {
//	data, _ := json.Marshal(Project{Name: "Name"})
//	events := []*es_models.Event{
//		&es_models.Event{AggregateID: "ID", Sequence: 1, Type: model.ProjectAdded, Data: data},
//	}
//	mockEs := mock.NewMockEventstore(ctrl)
//	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
//	return &ProjectEventstore{Eventstore: mockEs}
//}
//
//func GetMockProjectByIDNoEvents(ctrl *gomock.Controller) *ProjectEventstore {
//	events := []*es_models.Event{}
//	mockEs := mock.NewMockEventstore(ctrl)
//	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
//	return &ProjectEventstore{Eventstore: mockEs}
//}
//
//func GetMockManipulateProject(ctrl *gomock.Controller) *ProjectEventstore {
//	mockEs := mock.NewMockEventstore(ctrl)
//	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("TEST"))
//	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).Return(nil)
//	return &ProjectEventstore{Eventstore: mockEs}
//}
