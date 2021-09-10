package action

import (
	object_grpc "github.com/caos/zitadel/internal/api/grpc/object"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/query"
	action_pb "github.com/caos/zitadel/pkg/grpc/action"
	"google.golang.org/protobuf/types/known/durationpb"
)

func FlowTypeToDomain(flowType action_pb.FlowType) domain.FlowType {
	switch flowType { //TODO: converter
	case action_pb.FlowType_FLOW_TYPE_EXTERNAL_AUTHENTICATION:
		return domain.FlowTypeExternalAuthentication
	case action_pb.FlowType_FLOW_TYPE_EXTERNAL_REGISTRATION:
		return domain.FlowTypeExternalRegistration
	default:
		return domain.FlowTypeUnspecified
	}
}

func TriggerTypeToDomain(triggerType action_pb.TriggerType) domain.TriggerType {
	switch triggerType { //TODO: converter
	case action_pb.TriggerType_TRIGGER_TYPE_POST_AUTHENTICATION:
		return domain.TriggerTypePostAuthentication
	default:
		return domain.TriggerTypeUnspecified
	}
}

func FlowToPb(flow *query.Flow) *action_pb.Flow {
	return &action_pb.Flow{
		Type:    0,
		Details: nil,
		State:   0,
		//Triggers: TriggersToPb(flow.TriggerActions),
		TriggerActions: TriggerActionsToPb(flow.TriggerActions),
	}
}

//func TriggersToPb(triggers map[domain.TriggerType][]query.Action) map[int32]*action_pb.TriggerAction {
//	triggerActions := make(map[int32]*action_pb.TriggerAction)
//	for triggerType, actions := range triggers {
//		triggerActions[int32(triggerType)] = TriggerActionsToPb(actions)
//	}
//	return triggerActions
//}

func TriggerActionToPb(trigger domain.TriggerType, actions []*query.Action) *action_pb.TriggerAction {
	return &action_pb.TriggerAction{
		TriggerType: TriggerTypeToPb(trigger),
		Actions:     ActionsToPb(actions),
	}
}

func TriggerTypeToPb(triggerType domain.TriggerType) action_pb.TriggerType {
	switch triggerType { //TODO: converter
	case domain.TriggerTypePostAuthentication:
		return action_pb.TriggerType_TRIGGER_TYPE_POST_AUTHENTICATION
	default:
		return action_pb.TriggerType_TRIGGER_TYPE_UNSPECIFIED
	}
}

func TriggerActionsToPb(triggers map[domain.TriggerType][]*query.Action) []*action_pb.TriggerAction {
	list := make([]*action_pb.TriggerAction, 0)
	for trigger, actions := range triggers {
		list = append(list, TriggerActionToPb(trigger, actions))
	}
	return list
	//return &action_pb.TriggerAction{
	//	Actions: m,
	//}
}

func ActionsToPb(actions []*query.Action) []*action_pb.Action {
	list := make([]*action_pb.Action, len(actions))
	for i, action := range actions {
		list[i] = ActionToPb(action)
	}
	return list
}

func ActionToPb(action *query.Action) *action_pb.Action {
	return &action_pb.Action{
		Id:            action.ID,
		Details:       nil,
		State:         0,
		Name:          action.Name,
		Script:        action.Script,
		Timeout:       durationpb.New(action.Timeout),
		AllowedToFail: action.AllowedToFail,
	}
}

func ActionNameQuery(q *action_pb.ActionNameQuery) (query.SearchQuery, error) {
	return query.NewActionNameSearchQuery(object_grpc.TextMethodToQuery(q.Method), q.Name)
}
func ActionStateQuery(q *action_pb.ActionStateQuery) (query.SearchQuery, error) {
	return query.NewActionStateSearchQuery(ActionStateToDomain(q.State))
}

func ActionStateToDomain(state action_pb.ActionState) domain.ActionState {
	switch state {
	case action_pb.ActionState_ACTION_STATE_ACTIVE:
		return domain.ActionStateActive
	case action_pb.ActionState_ACTION_STATE_INACTIVE:
		return domain.ActionStateInactive
	default:
		return domain.ActionStateUnspecified
	}
}
