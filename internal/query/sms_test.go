package query

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/domain"
	errs "github.com/caos/zitadel/internal/errors"
)

var (
	expectedSMSConfigQuery = regexp.QuoteMeta(`SELECT zitadel.projections.sms_configs.id,` +
		` zitadel.projections.sms_configs.aggregate_id,` +
		` zitadel.projections.sms_configs.creation_date,` +
		` zitadel.projections.sms_configs.change_date,` +
		` zitadel.projections.sms_configs.resource_owner,` +
		` zitadel.projections.sms_configs.state,` +
		` zitadel.projections.sms_configs.sequence,` +

		// twilio config
		` zitadel.projections.sms_configs_twilio.sms_id,` +
		` zitadel.projections.sms_configs_twilio.sid,` +
		` zitadel.projections.sms_configs_twilio.token,` +
		` zitadel.projections.sms_configs_twilio.sender_name` +
		` FROM zitadel.projections.sms_configs` +
		` LEFT JOIN zitadel.projections.sms_configs_twilio ON zitadel.projections.sms_configs.id = zitadel.projections.sms_configs_twilio.sms_id`)
	expectedSMSConfigsQuery = regexp.QuoteMeta(`SELECT zitadel.projections.sms_configs.id,` +
		` zitadel.projections.sms_configs.aggregate_id,` +
		` zitadel.projections.sms_configs.creation_date,` +
		` zitadel.projections.sms_configs.change_date,` +
		` zitadel.projections.sms_configs.resource_owner,` +
		` zitadel.projections.sms_configs.state,` +
		` zitadel.projections.sms_configs.sequence,` +

		// twilio config
		` zitadel.projections.sms_configs_twilio.sms_id,` +
		` zitadel.projections.sms_configs_twilio.sid,` +
		` zitadel.projections.sms_configs_twilio.token,` +
		` zitadel.projections.sms_configs_twilio.sender_name,` +
		` COUNT(*) OVER ()` +
		` FROM zitadel.projections.sms_configs` +
		` LEFT JOIN zitadel.projections.sms_configs_twilio ON zitadel.projections.sms_configs.id = zitadel.projections.sms_configs_twilio.sms_id`)

	smsConfigCols = []string{
		"id",
		"aggregate_id",
		"creation_date",
		"change_date",
		"resource_owner",
		"state",
		"sequence",
		// twilio config
		"sms_id",
		"sid",
		"token",
		"sender-name",
	}
	smsConfigsCols = append(smsConfigCols, "count")
)

func Test_SMSConfigssPrepare(t *testing.T) {
	type want struct {
		sqlExpectations sqlExpectation
		err             checkErr
	}
	tests := []struct {
		name    string
		prepare interface{}
		want    want
		object  interface{}
	}{
		{
			name:    "prepareSMSConfigsQuery no result",
			prepare: prepareSMSConfigsQuery,
			want: want{
				sqlExpectations: mockQueries(
					expectedSMSConfigsQuery,
					nil,
					nil,
				),
			},
			object: &SMSConfigs{Configs: []*SMSConfig{}},
		},
		{
			name:    "prepareSMSQuery twilio config",
			prepare: prepareSMSConfigsQuery,
			want: want{
				sqlExpectations: mockQueries(
					expectedSMSConfigsQuery,
					smsConfigsCols,
					[][]driver.Value{
						{
							"sms-id",
							"agg-id",
							testNow,
							testNow,
							"ro",
							domain.SMSConfigStateInactive,
							uint64(20211109),
							// twilio config
							"sms-id",
							"sid",
							&crypto.CryptoValue{},
							"sender-name",
						},
					},
				),
			},
			object: &SMSConfigs{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				Configs: []*SMSConfig{
					{
						ID:            "sms-id",
						AggregateID:   "agg-id",
						CreationDate:  testNow,
						ChangeDate:    testNow,
						ResourceOwner: "ro",
						State:         domain.SMSConfigStateInactive,
						Sequence:      20211109,
						TwilioConfig: &Twilio{
							SID:        "sid",
							Token:      &crypto.CryptoValue{},
							SenderName: "sender-name",
						},
					},
				},
			},
		},
		{
			name:    "prepareSMSConfigsQuery multiple result",
			prepare: prepareSMSConfigsQuery,
			want: want{
				sqlExpectations: mockQueries(
					expectedSMSConfigsQuery,
					smsConfigsCols,
					[][]driver.Value{
						{
							"sms-id",
							"agg-id",
							testNow,
							testNow,
							"ro",
							domain.SMSConfigStateInactive,
							uint64(20211109),
							// twilio config
							"sms-id",
							"sid",
							&crypto.CryptoValue{},
							"sender-name",
						},
						{
							"sms-id2",
							"agg-id",
							testNow,
							testNow,
							"ro",
							domain.SMSConfigStateInactive,
							uint64(20211109),
							// twilio config
							"sms-id2",
							"sid2",
							&crypto.CryptoValue{},
							"sender-name2",
						},
					},
				),
			},
			object: &SMSConfigs{
				SearchResponse: SearchResponse{
					Count: 2,
				},
				Configs: []*SMSConfig{
					{
						ID:            "sms-id",
						AggregateID:   "agg-id",
						CreationDate:  testNow,
						ChangeDate:    testNow,
						ResourceOwner: "ro",
						State:         domain.SMSConfigStateInactive,
						Sequence:      20211109,
						TwilioConfig: &Twilio{
							SID:        "sid",
							Token:      &crypto.CryptoValue{},
							SenderName: "sender-name",
						},
					},
					{
						ID:            "sms-id2",
						AggregateID:   "agg-id",
						CreationDate:  testNow,
						ChangeDate:    testNow,
						ResourceOwner: "ro",
						State:         domain.SMSConfigStateInactive,
						Sequence:      20211109,
						TwilioConfig: &Twilio{
							SID:        "sid2",
							Token:      &crypto.CryptoValue{},
							SenderName: "sender-name2",
						},
					},
				},
			},
		},
		{
			name:    "prepareSMSConfigsQuery sql err",
			prepare: prepareSMSConfigsQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					expectedSMSConfigsQuery,
					sql.ErrConnDone,
				),
				err: func(err error) (error, bool) {
					if !errors.Is(err, sql.ErrConnDone) {
						return fmt.Errorf("err should be sql.ErrConnDone got: %w", err), false
					}
					return nil, true
				},
			},
			object: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err)
		})
	}
}

func Test_SMSConfigPrepare(t *testing.T) {
	type want struct {
		sqlExpectations sqlExpectation
		err             checkErr
	}
	tests := []struct {
		name    string
		prepare interface{}
		want    want
		object  interface{}
	}{
		{
			name:    "prepareSMSConfigQuery no result",
			prepare: prepareSMSConfigQuery,
			want: want{
				sqlExpectations: mockQueries(
					expectedSMSConfigQuery,
					nil,
					nil,
				),
				err: func(err error) (error, bool) {
					if !errs.IsNotFound(err) {
						return fmt.Errorf("err should be zitadel.NotFoundError got: %w", err), false
					}
					return nil, true
				},
			},
			object: (*SMSConfig)(nil),
		},
		{
			name:    "prepareSMSConfigQuery found",
			prepare: prepareSMSConfigQuery,
			want: want{
				sqlExpectations: mockQuery(
					expectedSMSConfigQuery,
					smsConfigCols,
					[]driver.Value{
						"sms-id",
						"agg-id",
						testNow,
						testNow,
						"ro",
						domain.SMSConfigStateInactive,
						uint64(20211109),
						// twilio config
						"sms-id",
						"sid",
						&crypto.CryptoValue{},
						"sender-name",
					},
				),
			},
			object: &SMSConfig{
				ID:            "sms-id",
				AggregateID:   "agg-id",
				CreationDate:  testNow,
				ChangeDate:    testNow,
				ResourceOwner: "ro",
				State:         domain.SMSConfigStateInactive,
				Sequence:      20211109,
				TwilioConfig: &Twilio{
					SID:        "sid",
					SenderName: "sender-name",
					Token:      &crypto.CryptoValue{},
				},
			},
		},
		{
			name:    "prepareSMSConfigQuery sql err",
			prepare: prepareSMSConfigQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					expectedSMSConfigQuery,
					sql.ErrConnDone,
				),
				err: func(err error) (error, bool) {
					if !errors.Is(err, sql.ErrConnDone) {
						return fmt.Errorf("err should be sql.ErrConnDone got: %w", err), false
					}
					return nil, true
				},
			},
			object: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err)
		})
	}
}
