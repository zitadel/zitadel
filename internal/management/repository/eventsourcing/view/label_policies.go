package view

// import (
// 	"github.com/caos/zitadel/internal/eventstore/v1/models"
// 	"github.com/caos/zitadel/internal/iam/repository/view"
// 	"github.com/caos/zitadel/internal/iam/repository/view/model"
// 	global_view "github.com/caos/zitadel/internal/view/repository"
// )

// const (
// 	labelPolicyTable = "management.label_policies"
// )

// func (v *View) StylingByAggregateIDAndState(aggregateID string, state int32) (*model.LabelPolicyView, error) {
// 	return view.GetStylingByAggregateIDAndState(v.Db, labelPolicyTable, aggregateID, state)
// }

// func (v *View) PutStyling(policy *model.LabelPolicyView, event *models.Event) error {
// 	err := view.PutStyling(v.Db, labelPolicyTable, policy)
// 	if err != nil {
// 		return err
// 	}
// 	return v.ProcessedStylingSequence(event)
// }

// func (v *View) ProcessedStylingSequence(event *models.Event) error {
// 	return v.saveCurrentSequence(labelPolicyTable, event)
// }

// func (v *View) UpdateStylingSpoolerRunTimestamp() error {
// 	return v.updateSpoolerRunSequence(labelPolicyTable)
// }

// func (v *View) GetLatestStylingFailedEvent(sequence uint64) (*global_view.FailedEvent, error) {
// 	return v.latestFailedEvent(labelPolicyTable, sequence)
// }

// func (v *View) ProcessedStylingFailedEvent(failedEvent *global_view.FailedEvent) error {
// 	return v.saveFailedEvent(failedEvent)
// }
