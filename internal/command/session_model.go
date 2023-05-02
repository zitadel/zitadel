package command

import (
	"bytes"
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/session"
)

type SessionWriteModel struct {
	eventstore.WriteModel

	Token             string
	UserID            string
	UserCheckedAt     time.Time
	PasswordCheckedAt time.Time
	Metadata          map[string][]byte
	State             domain.SessionState

	event *session.SetEvent
}

func NewSessionWriteModel(sessionID string, resourceOwner string) *SessionWriteModel {
	//var resourceOwner string //TODO: resourceowner?
	return &SessionWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   sessionID,
			ResourceOwner: resourceOwner,
		},
		Metadata: make(map[string][]byte),
	}
}

func (wm *SessionWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		//case *session.AddedEvent:
		//	wm.reduceAdded(e)
		case *session.SetEvent:
			wm.reduceSet(e)
		case *session.TerminateEvent:
			wm.reduceTerminate()
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *SessionWriteModel) Query() *eventstore.SearchQueryBuilder {
	query := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AddQuery().
		AggregateTypes(session.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			//session.AddedType,
			session.SetType,
			session.TerminateType,
		).
		Builder()

	if wm.ResourceOwner != "" {
		query.ResourceOwner(wm.ResourceOwner)
	}
	return query
}

func (wm *SessionWriteModel) reduceSet(e *session.SetEvent) {
	wm.State = domain.SessionStateActive
	wm.Token = e.Token
	if e.UserID != nil {
		wm.UserID = *e.UserID
	}
	if e.UserCheckedAt != nil {
		wm.UserCheckedAt = *e.UserCheckedAt
	}
	if e.PasswordCheckedAt != nil {
		wm.PasswordCheckedAt = *e.PasswordCheckedAt
	}
	if len(e.Metadata) != 0 {
		wm.Metadata = e.Metadata
	}
}

func (wm *SessionWriteModel) reduceTerminate() {
	wm.State = domain.SessionStateTerminated
}

func (wm *SessionWriteModel) UserChecked(ctx context.Context, userID string, checkedAt time.Time) {
	wm.setEvent(ctx).AddUserData(userID, checkedAt)
	wm.UserID = userID
}

func (wm *SessionWriteModel) PasswordChecked(ctx context.Context, checkedAt time.Time) {
	wm.setEvent(ctx).AddPasswordData(checkedAt)
}

func (wm *SessionWriteModel) SetToken(ctx context.Context) *session.SetEvent {
	wm.setEvent(ctx).SetToken(time.Now().String())
	return wm.event
}

func (wm *SessionWriteModel) setEvent(ctx context.Context) *session.SetEvent {
	if wm.event == nil {
		wm.event = session.NewSetEvent(ctx, &session.NewAggregate(wm.AggregateID, wm.ResourceOwner).Aggregate)
	}
	return wm.event
}

func (wm *SessionWriteModel) ChangeMetadata(ctx context.Context, metadata map[string][]byte) error {
	var changed bool
	for key, value := range metadata {
		currentValue, ok := wm.Metadata[key]
		//if !ok && len(value) != 0 {
		//	changed = true
		//}
		//changed = !bytes.Equal(currentValue, value)

		if len(value) != 0 {
			if !bytes.Equal(currentValue, value) {
				// if there's a value, just set / update the
				wm.Metadata[key] = value
				changed = true
			}
		} else {
			if ok {
				delete(wm.Metadata, key)
				changed = true
			}
			//_, ok := wm.Metadata[key]
			//if !ok {
			//	// do not allow empty values for not existing entries
			//	return caos_errs.ThrowInvalidArgument(nil, "SESSION-SDf4g", "metadata empty") //TODO: i18n
			//}
		}
	}
	if changed {
		wm.setEvent(ctx).AddMetadata(wm.Metadata)
	}
	return nil
}
