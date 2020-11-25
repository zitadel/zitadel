package model

import (
	"encoding/json"
	"github.com/caos/zitadel/pkg/grpc/auth"
	"testing"

	es_models "github.com/caos/zitadel/internal/eventstore/models"
)

func TestAppendMFAU2FAddedEvent(t *testing.T) {
	type args struct {
		user  *Human
		u2f   *WebAuthNToken
		event *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *Human
	}{
		{
			name: "append user u2f event",
			args: args{
				user:  &Human{},
				u2f:   &WebAuthNToken{WebauthNTokenID: "WebauthNTokenID", Challenge: "Challenge"},
				event: &es_models.Event{},
			},
			result: &Human{
				U2FTokens: []*WebAuthNToken{
					{WebauthNTokenID: "WebauthNTokenID", Challenge: "Challenge", State: int32(auth.MFAState_MFASTATE_NOT_READY)},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.u2f != nil {
				data, _ := json.Marshal(tt.args.u2f)
				tt.args.event.Data = data
			}
			tt.args.user.appendU2FAddedEvent(tt.args.event)
			if tt.args.user.U2FTokens[0].State != tt.result.U2FTokens[0].State {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.U2FTokens[0].State, tt.args.user.U2FTokens[0].State)
			}
		})
	}
}

func TestAppendMFAU2FVerifyEvent(t *testing.T) {
	type args struct {
		user  *Human
		u2f   *WebAuthNVerify
		event *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *Human
	}{
		{
			name: "append u2f verify event",
			args: args{
				user: &Human{
					U2FTokens: []*WebAuthNToken{
						{WebauthNTokenID: "WebauthNTokenID", Challenge: "Challenge", State: int32(auth.MFAState_MFASTATE_NOT_READY)},
					},
				},
				u2f:   &WebAuthNVerify{WebAuthNTokenID: "WebauthNTokenID", KeyID: []byte("KeyID"), PublicKey: []byte("PublicKey"), AttestationType: "AttestationType", AAGUID: []byte("AAGUID"), SignCount: 1},
				event: &es_models.Event{},
			},
			result: &Human{
				U2FTokens: []*WebAuthNToken{
					{
						WebauthNTokenID: "WebauthNTokenID",
						Challenge:       "Challenge",
						State:           int32(auth.MFAState_MFASTATE_READY),
						KeyID:           []byte("KeyID"),
						PublicKey:       []byte("PublicKey"),
						AttestationType: "AttestationType",
						AAGUID:          []byte("AAGUID"),
						SignCount:       1,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.u2f != nil {
				data, _ := json.Marshal(tt.args.u2f)
				tt.args.event.Data = data
			}
			tt.args.user.appendU2FVerifiedEvent(tt.args.event)
			if tt.args.user.U2FTokens[0].State != tt.result.U2FTokens[0].State {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.U2FTokens[0].State, tt.args.user.U2FTokens[0].State)
			}
			if tt.args.user.U2FTokens[0].AttestationType != tt.result.U2FTokens[0].AttestationType {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.U2FTokens[0].AttestationType, tt.args.user.U2FTokens[0].AttestationType)
			}
		})
	}
}

func TestAppendMFAU2FRemoveEvent(t *testing.T) {
	type args struct {
		user  *Human
		u2f   *WebAuthNTokenID
		event *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *Human
	}{
		{
			name: "append u2f remove event",
			args: args{
				user: &Human{
					U2FTokens: []*WebAuthNToken{
						{
							WebauthNTokenID: "WebauthNTokenID",
							Challenge:       "Challenge",
							State:           int32(auth.MFAState_MFASTATE_NOT_READY),
							KeyID:           []byte("KeyID"),
							PublicKey:       []byte("PublicKey"),
							AttestationType: "AttestationType",
							AAGUID:          []byte("AAGUID"),
							SignCount:       1,
						},
					},
				},
				u2f:   &WebAuthNTokenID{WebauthNTokenID: "WebauthNTokenID"},
				event: &es_models.Event{},
			},
			result: &Human{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.u2f != nil {
				data, _ := json.Marshal(tt.args.u2f)
				tt.args.event.Data = data
			}
			tt.args.user.appendU2FRemovedEvent(tt.args.event)
			if len(tt.args.user.U2FTokens) != 0 {
				t.Errorf("got wrong result: actual: %v ", tt.result.U2FTokens)
			}
		})
	}
}
