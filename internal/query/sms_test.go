package query

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var (
	expectedSMSConfigQuery = regexp.QuoteMeta(`SELECT projections.sms_configs2.id,` +
		` projections.sms_configs2.aggregate_id,` +
		` projections.sms_configs2.creation_date,` +
		` projections.sms_configs2.change_date,` +
		` projections.sms_configs2.resource_owner,` +
		` projections.sms_configs2.state,` +
		` projections.sms_configs2.sequence,` +

		// twilio config
		` projections.sms_configs2_twilio.sms_id,` +
		` projections.sms_configs2_twilio.sid,` +
		` projections.sms_configs2_twilio.token,` +
		` projections.sms_configs2_twilio.sender_number` +
		` FROM projections.sms_configs2` +
		` LEFT JOIN projections.sms_configs2_twilio ON projections.sms_configs2.id = projections.sms_configs2_twilio.sms_id AND projections.sms_configs2.instance_id = projections.sms_configs2_twilio.instance_id` +
		` AS OF SYSTEM TIME '-1 ms'`)
	expectedSMSConfigsQuery = regexp.QuoteMeta(`SELECT projections.sms_configs2.id,` +
		` projections.sms_configs2.aggregate_id,` +
		` projections.sms_configs2.creation_date,` +
		` projections.sms_configs2.change_date,` +
		` projections.sms_configs2.resource_owner,` +
		` projections.sms_configs2.state,` +
		` projections.sms_configs2.sequence,` +

		// twilio config
		` projections.sms_configs2_twilio.sms_id,` +
		` projections.sms_configs2_twilio.sid,` +
		` projections.sms_configs2_twilio.token,` +
		` projections.sms_configs2_twilio.sender_number,` +
		` COUNT(*) OVER ()` +
		` FROM projections.sms_configs2` +
		` LEFT JOIN projections.sms_configs2_twilio ON projections.sms_configs2.id = projections.sms_configs2_twilio.sms_id AND projections.sms_configs2.instance_id = projections.sms_configs2_twilio.instance_id` +
		` AS OF SYSTEM TIME '-1 ms'`)

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
		"sender-number",
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
							"sender-number",
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
							SID:              "sid",
							Token:            &crypto.CryptoValue{},
							SenderNumber:     "sender-number",
							VerifyServiceSID: "",
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
							"sender-number",
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
							"sender-number2",
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
							SID:              "sid",
							Token:            &crypto.CryptoValue{},
							SenderNumber:     "sender-number",
							VerifyServiceSID: "",
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
							SID:              "sid2",
							Token:            &crypto.CryptoValue{},
							SenderNumber:     "sender-number2",
							VerifyServiceSID: "",
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
			object: (*SMSConfigs)(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err, defaultPrepareArgs...)
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
				sqlExpectations: mockQueriesScanErr(
					expectedSMSConfigQuery,
					nil,
					nil,
				),
				err: func(err error) (error, bool) {
					if !zerrors.IsNotFound(err) {
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
						"sender-number",
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
					SID:              "sid",
					SenderNumber:     "sender-number",
					Token:            &crypto.CryptoValue{},
					VerifyServiceSID: "",
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
			object: (*SMSConfig)(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err, defaultPrepareArgs...)
		})
	}
}
