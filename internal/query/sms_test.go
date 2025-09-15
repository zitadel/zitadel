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
	expectedSMSConfigQuery = regexp.QuoteMeta(`SELECT projections.sms_configs3.id,` +
		` projections.sms_configs3.aggregate_id,` +
		` projections.sms_configs3.creation_date,` +
		` projections.sms_configs3.change_date,` +
		` projections.sms_configs3.resource_owner,` +
		` projections.sms_configs3.state,` +
		` projections.sms_configs3.sequence,` +
		` projections.sms_configs3.description,` +

		// twilio config
		` projections.sms_configs3_twilio.sms_id,` +
		` projections.sms_configs3_twilio.sid,` +
		` projections.sms_configs3_twilio.token,` +
		` projections.sms_configs3_twilio.sender_number,` +
		` projections.sms_configs3_twilio.verify_service_sid,` +

		// http config
		` projections.sms_configs3_http.sms_id,` +
		` projections.sms_configs3_http.endpoint,` +
		` projections.sms_configs3_http.signing_key` +
		` FROM projections.sms_configs3` +
		` LEFT JOIN projections.sms_configs3_twilio ON projections.sms_configs3.id = projections.sms_configs3_twilio.sms_id AND projections.sms_configs3.instance_id = projections.sms_configs3_twilio.instance_id` +
		` LEFT JOIN projections.sms_configs3_http ON projections.sms_configs3.id = projections.sms_configs3_http.sms_id AND projections.sms_configs3.instance_id = projections.sms_configs3_http.instance_id`)
	expectedSMSConfigsQuery = regexp.QuoteMeta(`SELECT projections.sms_configs3.id,` +
		` projections.sms_configs3.aggregate_id,` +
		` projections.sms_configs3.creation_date,` +
		` projections.sms_configs3.change_date,` +
		` projections.sms_configs3.resource_owner,` +
		` projections.sms_configs3.state,` +
		` projections.sms_configs3.sequence,` +
		` projections.sms_configs3.description,` +

		// twilio config
		` projections.sms_configs3_twilio.sms_id,` +
		` projections.sms_configs3_twilio.sid,` +
		` projections.sms_configs3_twilio.token,` +
		` projections.sms_configs3_twilio.sender_number,` +
		` projections.sms_configs3_twilio.verify_service_sid,` +

		// http config
		` projections.sms_configs3_http.sms_id,` +
		` projections.sms_configs3_http.endpoint,` +
		` projections.sms_configs3_http.signing_key,` +
		` COUNT(*) OVER ()` +
		` FROM projections.sms_configs3` +
		` LEFT JOIN projections.sms_configs3_twilio ON projections.sms_configs3.id = projections.sms_configs3_twilio.sms_id AND projections.sms_configs3.instance_id = projections.sms_configs3_twilio.instance_id` +
		` LEFT JOIN projections.sms_configs3_http ON projections.sms_configs3.id = projections.sms_configs3_http.sms_id AND projections.sms_configs3.instance_id = projections.sms_configs3_http.instance_id`)

	smsConfigCols = []string{
		"id",
		"aggregate_id",
		"creation_date",
		"change_date",
		"resource_owner",
		"state",
		"sequence",
		"description",
		// twilio config
		"sms_id",
		"sid",
		"token",
		"sender-number",
		"verify_sid",
		// http config
		"sms_id",
		"endpoint",
		"signing_key",
	}
	smsConfigsCols = append(smsConfigCols, "count")
)

func Test_SMSConfigsPrepare(t *testing.T) {
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
							"description",
							// twilio config
							"sms-id",
							"sid",
							&crypto.CryptoValue{},
							"sender-number",
							"",
							// http config
							nil,
							nil,
							nil,
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
						Description:   "description",
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
			name:    "prepareSMSQuery http config",
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
							"description",
							// twilio config
							nil,
							nil,
							nil,
							nil,
							nil,
							// http config
							"sms-id",
							"endpoint",
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "alg",
								KeyID:      "encKey",
								Crypted:    []byte("crypted"),
							},
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
						Description:   "description",
						HTTPConfig: &HTTP{
							Endpoint: "endpoint",
							signingKey: &crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "alg",
								KeyID:      "encKey",
								Crypted:    []byte("crypted"),
							},
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
							domain.SMSConfigStateActive,
							uint64(20211109),
							"description",
							// twilio config
							"sms-id",
							"sid",
							&crypto.CryptoValue{},
							"sender-number",
							"verify-service-sid",
							// http config
							nil,
							nil,
							nil,
						},
						{
							"sms-id2",
							"agg-id",
							testNow,
							testNow,
							"ro",
							domain.SMSConfigStateInactive,
							uint64(20211109),
							"description",
							// twilio config
							"sms-id2",
							"sid2",
							&crypto.CryptoValue{},
							"sender-number2",
							"verify-service-sid2",
							// http config
							nil,
							nil,
							nil,
						},
						{
							"sms-id3",
							"agg-id",
							testNow,
							testNow,
							"ro",
							domain.SMSConfigStateInactive,
							uint64(20211109),
							"description",
							// twilio config
							nil,
							nil,
							nil,
							nil,
							nil,
							// http config
							"sms-id3",
							"endpoint3",
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "alg",
								KeyID:      "encKey",
								Crypted:    []byte("crypted"),
							},
						},
					},
				),
			},
			object: &SMSConfigs{
				SearchResponse: SearchResponse{
					Count: 3,
				},
				Configs: []*SMSConfig{
					{
						ID:            "sms-id",
						AggregateID:   "agg-id",
						CreationDate:  testNow,
						ChangeDate:    testNow,
						ResourceOwner: "ro",
						State:         domain.SMSConfigStateActive,
						Sequence:      20211109,
						Description:   "description",
						TwilioConfig: &Twilio{
							SID:              "sid",
							Token:            &crypto.CryptoValue{},
							SenderNumber:     "sender-number",
							VerifyServiceSID: "verify-service-sid",
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
						Description:   "description",
						TwilioConfig: &Twilio{
							SID:              "sid2",
							Token:            &crypto.CryptoValue{},
							SenderNumber:     "sender-number2",
							VerifyServiceSID: "verify-service-sid2",
						},
					},
					{
						ID:            "sms-id3",
						AggregateID:   "agg-id",
						CreationDate:  testNow,
						ChangeDate:    testNow,
						ResourceOwner: "ro",
						State:         domain.SMSConfigStateInactive,
						Sequence:      20211109,
						Description:   "description",
						HTTPConfig: &HTTP{
							Endpoint: "endpoint3",
							signingKey: &crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "alg",
								KeyID:      "encKey",
								Crypted:    []byte("crypted"),
							},
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
			name:    "prepareSMSConfigQuery, twilio, found",
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
						domain.SMSConfigStateActive,
						uint64(20211109),
						"description",
						// twilio config
						"sms-id",
						"sid",
						&crypto.CryptoValue{},
						"sender-number",
						"verify-service-sid",
						// http config
						nil,
						nil,
						nil,
					},
				),
			},
			object: &SMSConfig{
				ID:            "sms-id",
				AggregateID:   "agg-id",
				CreationDate:  testNow,
				ChangeDate:    testNow,
				ResourceOwner: "ro",
				State:         domain.SMSConfigStateActive,
				Sequence:      20211109,
				Description:   "description",
				TwilioConfig: &Twilio{
					SID:              "sid",
					SenderNumber:     "sender-number",
					Token:            &crypto.CryptoValue{},
					VerifyServiceSID: "verify-service-sid",
				},
			},
		},
		{
			name:    "prepareSMSConfigQuery, http, found",
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
						"description",
						// twilio config
						nil,
						nil,
						nil,
						nil,
						nil,
						// http config
						"sms-id",
						"endpoint",
						&crypto.CryptoValue{
							CryptoType: crypto.TypeEncryption,
							Algorithm:  "alg",
							KeyID:      "encKey",
							Crypted:    []byte("crypted"),
						},
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
				Description:   "description",
				HTTPConfig: &HTTP{
					Endpoint: "endpoint",
					signingKey: &crypto.CryptoValue{
						CryptoType: crypto.TypeEncryption,
						Algorithm:  "alg",
						KeyID:      "encKey",
						Crypted:    []byte("crypted"),
					},
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
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err)
		})
	}
}
