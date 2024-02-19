package target

import (
	"strings"

	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	UniqueExecution = "execution"
	DuplicateTarget = "Errors.Execution.AlreadyExists"

	conditionRequest  = "request"
	conditionResponse = "response"
	conditionFunction = "function"
	conditionEvent    = "event"
)

func ConditionRequestService(name string) string {
	return strings.Join([]string{conditionRequest, name}, "-")
}
func ConditionRequestMethod(name string) string {
	return strings.Join([]string{conditionRequest, name}, "-")
}
func ConditionRequestAll() string {
	return conditionRequest
}
func ConditionResponseService(name string) string {
	return strings.Join([]string{conditionResponse, name}, "-")
}
func ConditionResponseMethod(name string) string {
	return strings.Join([]string{conditionResponse, name}, "-")
}
func ConditionResponseAll() string {
	return conditionResponse
}
func ConditionFunction(name string) string {
	return strings.Join([]string{conditionFunction, name}, "-")
}
func ConditionEvent(name string) string {
	return strings.Join([]string{conditionEvent, name}, "-")
}
func ConditionEventGroup(name string) string {
	return strings.Join([]string{conditionEvent, name}, "-")
}
func ConditionEventAll() string {
	return conditionEvent
}

func NewAddUniqueConstraint(condition string) *eventstore.UniqueConstraint {
	return eventstore.NewAddEventUniqueConstraint(
		UniqueExecution,
		condition,
		DuplicateTarget,
	)
}

func NewRemoveUniqueConstraint(condition string) *eventstore.UniqueConstraint {
	return eventstore.NewRemoveUniqueConstraint(
		UniqueExecution,
		condition,
	)
}
