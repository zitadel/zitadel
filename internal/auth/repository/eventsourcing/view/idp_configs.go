package view

import (
	"github.com/caos/zitadel/internal/errors"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/iam/repository/view"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/view/model"
	global_view "github.com/caos/zitadel/internal/view/repository"
	"time"
)

const (
	idpConfigTable = "auth.idp_configs"
)

func (v *View) IDPConfigByID(idpID string) (*iam_es_model.IDPConfigView, error) {
	return view.IDPByID(v.Db, idpConfigTable, idpID)
}

func (v *View) GetIDPConfigsByAggregateID(aggregateID string) ([]*iam_es_model.IDPConfigView, error) {
	return view.GetIDPConfigsByAggregateID(v.Db, idpConfigTable, aggregateID)
}

func (v *View) SearchIDPConfigs(request *iam_model.IDPConfigSearchRequest) ([]*iam_es_model.IDPConfigView, uint64, error) {
	return view.SearchIDPs(v.Db, idpConfigTable, request)
}

func (v *View) PutIDPConfig(idp *iam_es_model.IDPConfigView, sequence uint64, eventTimestamp time.Time) error {
	err := view.PutIDP(v.Db, idpConfigTable, idp)
	if err != nil {
		return err
	}
	return v.ProcessedIDPConfigSequence(sequence, eventTimestamp)
}

func (v *View) DeleteIDPConfig(idpID string, eventSequence uint64, eventTimestamp time.Time) error {
	err := view.DeleteIDP(v.Db, idpConfigTable, idpID)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return v.ProcessedIDPConfigSequence(eventSequence, eventTimestamp)
}

func (v *View) GetLatestIDPConfigSequence() (*global_view.CurrentSequence, error) {
	return v.latestSequence(idpConfigTable)
}

func (v *View) ProcessedIDPConfigSequence(eventSequence uint64, eventTimestamp time.Time) error {
	return v.saveCurrentSequence(idpConfigTable, eventSequence, eventTimestamp)
}

func (v *View) UpdateIDPConfigSpoolerRunTimestamp() error {
	return v.updateSpoolerRunSequence(idpConfigTable)
}

func (v *View) GetLatestIDPConfigFailedEvent(sequence uint64) (*global_view.FailedEvent, error) {
	return v.latestFailedEvent(idpConfigTable, sequence)
}

func (v *View) ProcessedIDPConfigFailedEvent(failedEvent *global_view.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}
