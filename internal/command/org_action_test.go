package command

import (
	"context"
	"testing"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/id"
	"github.com/stretchr/testify/assert"
)

func TestCommands_AddAction(t *testing.T) {
	type fields struct {
		eventstore  *eventstore.Eventstore
		idGenerator id.Generator
	}
	type args struct {
		ctx           context.Context
		addAction     *domain.Action
		resourceOwner string
	}
	type res struct {
		id      string
		details *domain.ObjectDetails
		err     func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"no name, error",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx: context.Background(),
				addAction: &domain.Action{
					Script: "test()",
				},
				resourceOwner: "org1",
			},
			res{
				err: errors.IsErrorInvalidArgument,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.fields.eventstore,
				//static:                      tt.fields.static,
				idGenerator: tt.fields.idGenerator,
				//iamDomain:                   tt.fields.iamDomain,
				//zitadelRoles:                tt.fields.zitadelRoles,
				//idpConfigSecretCrypto:       tt.fields.idpConfigSecretCrypto,
				//userPasswordAlg:             tt.fields.userPasswordAlg,
				//initializeUserCode:          tt.fields.initializeUserCode,
				//emailVerificationCode:       tt.fields.emailVerificationCode,
				//phoneVerificationCode:       tt.fields.phoneVerificationCode,
				//passwordVerificationCode:    tt.fields.passwordVerificationCode,
				//passwordlessInitCode:        tt.fields.passwordlessInitCode,
				//machineKeyAlg:               tt.fields.machineKeyAlg,
				//machineKeySize:              tt.fields.machineKeySize,
				//applicationKeySize:          tt.fields.applicationKeySize,
				//applicationSecretGenerator:  tt.fields.applicationSecretGenerator,
				//domainVerificationAlg:       tt.fields.domainVerificationAlg,
				//domainVerificationGenerator: tt.fields.domainVerificationGenerator,
				//domainVerificationValidator: tt.fields.domainVerificationValidator,
				//multifactors:                tt.fields.multifactors,
				//webauthn:                    tt.fields.webauthn,
				//keySize:                     tt.fields.keySize,
				//keyAlgorithm:                tt.fields.keyAlgorithm,
				//privateKeyLifetime:          tt.fields.privateKeyLifetime,
				//publicKeyLifetime:           tt.fields.publicKeyLifetime,
				//tokenVerifier:               tt.fields.tokenVerifier,
			}
			id, details, err := c.AddAction(tt.args.ctx, tt.args.addAction, tt.args.resourceOwner)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.id, id)
				assert.Equal(t, tt.res.details, details)
			}
		})
	}
}
