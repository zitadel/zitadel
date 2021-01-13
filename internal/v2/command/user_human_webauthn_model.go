package command

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/user"
)

type HumanWebAuthNWriteModel struct {
	eventstore.WriteModel

	WebauthNTokenID string
	Challenge       string

	KeyID             []byte
	PublicKey         []byte
	AttestationType   string
	AAGUID            []byte
	SignCount         uint32
	WebAuthNTokenName string

	State    domain.WebAuthNState
	MFAState domain.MFAState
}

func NewHumanWebAuthNWriteModel(userID, wbAuthNTokenID, resourceOwner string) *HumanWebAuthNWriteModel {
	return &HumanWebAuthNWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   userID,
			ResourceOwner: resourceOwner,
		},
		WebauthNTokenID: wbAuthNTokenID,
	}
}

func (wm *HumanWebAuthNWriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *user.HumanWebAuthNAddedEvent:
			if wm.WebauthNTokenID == e.WebAuthNTokenID {
				wm.AppendEvents(e)
			}
		case *user.HumanWebAuthNRemovedEvent:
			if wm.WebauthNTokenID == e.WebAuthNTokenID {
				wm.AppendEvents(e)
			}
		case *user.UserRemovedEvent:
			wm.AppendEvents(e)
		}
	}
}

func (wm *HumanWebAuthNWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *user.HumanWebAuthNAddedEvent:
			wm.appendAddedEvent(e)
		case *user.HumanWebAuthNVerifiedEvent:
			wm.appendVerifiedEvent(e)
		case *user.HumanWebAuthNRemovedEvent:
			wm.State = domain.WebAuthNStateRemoved
		case *user.UserRemovedEvent:
			wm.State = domain.WebAuthNStateRemoved
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *HumanWebAuthNWriteModel) appendAddedEvent(e *user.HumanWebAuthNAddedEvent) {
	wm.WebauthNTokenID = e.WebAuthNTokenID
	wm.Challenge = e.Challenge
	wm.State = domain.WebAuthNStateActive
	wm.MFAState = domain.MFAStateNotReady
}

func (wm *HumanWebAuthNWriteModel) appendVerifiedEvent(e *user.HumanWebAuthNVerifiedEvent) {
	wm.KeyID = e.KeyID
	wm.PublicKey = e.PublicKey
	wm.AttestationType = e.AttestationType
	wm.AAGUID = e.AAGUID
	wm.SignCount = e.SignCount
	wm.WebAuthNTokenName = e.WebAuthNTokenName
	wm.MFAState = domain.MFAStateReady
}

func (wm *HumanWebAuthNWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, user.AggregateType).
		AggregateIDs(wm.AggregateID).
		ResourceOwner(wm.ResourceOwner)
}

type HumanU2FTokensReadModel struct {
	eventstore.WriteModel

	WebAuthNTokens []*HumanWebAuthNWriteModel
	UserState      domain.UserState
}

func NewHumanU2FTokensReadModel(userID, resourceOwner string) *HumanU2FTokensReadModel {
	return &HumanU2FTokensReadModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   userID,
			ResourceOwner: resourceOwner,
		},
	}
}

func (wm *HumanU2FTokensReadModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *user.HumanWebAuthNAddedEvent:
			wm.AppendEvents(e)
		case *user.HumanWebAuthNVerifiedEvent:
			wm.AppendEvents(e)
		case *user.HumanWebAuthNRemovedEvent:
			wm.AppendEvents(e)
		case *user.UserRemovedEvent:
			wm.AppendEvents(e)
		}
	}
}

func (wm *HumanU2FTokensReadModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *user.HumanWebAuthNAddedEvent:
			token := &HumanWebAuthNWriteModel{}
			token.appendAddedEvent(e)
			wm.WebAuthNTokens = append(wm.WebAuthNTokens, token)
		case *user.HumanWebAuthNVerifiedEvent:
			idx, token := wm.WebAuthNTokenByID(e.WebAuthNTokenID)
			if idx < 0 {
				continue
			}
			token.appendVerifiedEvent(e)
		case *user.HumanWebAuthNRemovedEvent:
			idx, _ := wm.WebAuthNTokenByID(e.WebAuthNTokenID)
			if idx < 0 {
				continue
			}
			copy(wm.WebAuthNTokens[idx:], wm.WebAuthNTokens[idx+1:])
			wm.WebAuthNTokens[len(wm.WebAuthNTokens)-1] = nil
			wm.WebAuthNTokens = wm.WebAuthNTokens[:len(wm.WebAuthNTokens)-1]
		case *user.UserRemovedEvent:
			wm.UserState = domain.UserStateDeleted
		}
	}
	return wm.WriteModel.Reduce()
}

func (rm *HumanU2FTokensReadModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, user.AggregateType).
		AggregateIDs(rm.AggregateID).
		ResourceOwner(rm.ResourceOwner).
		EventTypes(
			user.HumanU2FTokenAddedType,
			user.HumanU2FTokenVerifiedType,
			user.HumanU2FTokenRemovedType,
			user.UserV1MFAOTPRemovedType)

}

func (wm *HumanU2FTokensReadModel) WebAuthNTokenByID(id string) (idx int, token *HumanWebAuthNWriteModel) {
	for idx, token = range wm.WebAuthNTokens {
		if token.WebauthNTokenID == id {
			return idx, token
		}
	}
	return -1, nil
}

type HumanPasswordlessTokensReadModel struct {
	eventstore.WriteModel

	WebAuthNTokens []*HumanWebAuthNWriteModel
	UserState      domain.UserState
}

func NewHumanPasswordlessTokensReadModel(userID, resourceOwner string) *HumanPasswordlessTokensReadModel {
	return &HumanPasswordlessTokensReadModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   userID,
			ResourceOwner: resourceOwner,
		},
	}
}

func (wm *HumanPasswordlessTokensReadModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *user.HumanWebAuthNAddedEvent:
			wm.AppendEvents(e)
		case *user.HumanWebAuthNVerifiedEvent:
			wm.AppendEvents(e)
		case *user.HumanWebAuthNRemovedEvent:
			wm.AppendEvents(e)
		case *user.UserRemovedEvent:
			wm.AppendEvents(e)
		}
	}
}

func (wm *HumanPasswordlessTokensReadModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *user.HumanWebAuthNAddedEvent:
			token := &HumanWebAuthNWriteModel{}
			token.appendAddedEvent(e)
			wm.WebAuthNTokens = append(wm.WebAuthNTokens, token)
		case *user.HumanWebAuthNVerifiedEvent:
			idx, token := wm.WebAuthNTokenByID(e.WebAuthNTokenID)
			if idx < 0 {
				continue
			}
			token.appendVerifiedEvent(e)
		case *user.HumanWebAuthNRemovedEvent:
			idx, _ := wm.WebAuthNTokenByID(e.WebAuthNTokenID)
			if idx < 0 {
				continue
			}
			copy(wm.WebAuthNTokens[idx:], wm.WebAuthNTokens[idx+1:])
			wm.WebAuthNTokens[len(wm.WebAuthNTokens)-1] = nil
			wm.WebAuthNTokens = wm.WebAuthNTokens[:len(wm.WebAuthNTokens)-1]
		case *user.UserRemovedEvent:
			wm.UserState = domain.UserStateDeleted
		}
	}
	return wm.WriteModel.Reduce()
}

func (rm *HumanPasswordlessTokensReadModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, user.AggregateType).
		AggregateIDs(rm.AggregateID).
		ResourceOwner(rm.ResourceOwner).
		EventTypes(
			user.HumanPasswordlessTokenAddedType,
			user.HumanPasswordlessTokenVerifiedType,
			user.HumanPasswordlessTokenRemovedType,
			user.UserV1MFAOTPRemovedType)

}

func (wm *HumanPasswordlessTokensReadModel) WebAuthNTokenByID(id string) (idx int, token *HumanWebAuthNWriteModel) {
	for idx, token = range wm.WebAuthNTokens {
		if token.WebauthNTokenID == id {
			return idx, token
		}
	}
	return -1, nil
}
