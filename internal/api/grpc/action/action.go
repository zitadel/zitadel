package action

import (
	object_grpc "github.com/caos/zitadel/internal/api/grpc/object"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/query"
	action_pb "github.com/caos/zitadel/pkg/grpc/action"
	"google.golang.org/protobuf/types/known/durationpb"
)

func FlowTypeToDomain(flowType action_pb.FlowType) domain.FlowType {
	switch flowType {
	case action_pb.FlowType_FLOW_TYPE_EXTERNAL_AUTHENTICATION:
		return domain.FlowTypeExternalAuthentication
	default:
		return domain.FlowTypeUnspecified
	}
}

func TriggerTypeToDomain(triggerType action_pb.TriggerType) domain.TriggerType {
	switch triggerType {
	case action_pb.TriggerType_TRIGGER_TYPE_POST_AUTHENTICATION:
		return domain.TriggerTypePostAuthentication
	case action_pb.TriggerType_TRIGGER_TYPE_PRE_CREATION:
		return domain.TriggerTypePreCreation
	case action_pb.TriggerType_TRIGGER_TYPE_POST_CREATION:
		return domain.TriggerTypePostCreation
	default:
		return domain.TriggerTypeUnspecified
	}
}

func FlowToPb(flow *query.Flow) *action_pb.Flow {
	return &action_pb.Flow{
		Type:           FlowTypeToPb(flow.Type),
		Details:        object_grpc.ChangeToDetailsPb(flow.Sequence, flow.ChangeDate, flow.ResourceOwner),
		State:          action_pb.FlowState_FLOW_STATE_ACTIVE, //TODO: state in next release
		TriggerActions: TriggerActionsToPb(flow.TriggerActions),
	}
}

func TriggerActionToPb(trigger domain.TriggerType, actions []*query.Action) *action_pb.TriggerAction {
	return &action_pb.TriggerAction{
		TriggerType: TriggerTypeToPb(trigger),
		Actions:     ActionsToPb(actions),
	}
}

func FlowTypeToPb(flowType domain.FlowType) action_pb.FlowType {
	switch flowType {
	case domain.FlowTypeExternalAuthentication:
		return action_pb.FlowType_FLOW_TYPE_EXTERNAL_AUTHENTICATION
	default:
		return action_pb.FlowType_FLOW_TYPE_UNSPECIFIED
	}
}

func TriggerTypeToPb(triggerType domain.TriggerType) action_pb.TriggerType {
	switch triggerType {
	case domain.TriggerTypePostAuthentication:
		return action_pb.TriggerType_TRIGGER_TYPE_POST_AUTHENTICATION
	case domain.TriggerTypePreCreation:
		return action_pb.TriggerType_TRIGGER_TYPE_PRE_CREATION
	case domain.TriggerTypePostCreation:
		return action_pb.TriggerType_TRIGGER_TYPE_POST_CREATION
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
		Details:       object_grpc.ChangeToDetailsPb(action.Sequence, action.ChangeDate, action.ResourceOwner),
		State:         ActionStateToPb(action.State),
		Name:          action.Name,
		Script:        action.Script,
		Timeout:       durationpb.New(action.Timeout),
		AllowedToFail: action.AllowedToFail,
	}
}

func ActionStateToPb(state domain.ActionState) action_pb.ActionState {
	switch state {
	case domain.ActionStateActive:
		return action_pb.ActionState_ACTION_STATE_ACTIVE
	case domain.ActionStateInactive:
		return action_pb.ActionState_ACTION_STATE_INACTIVE
	default:
		return action_pb.ActionState_ACTION_STATE_UNSPECIFIED
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
