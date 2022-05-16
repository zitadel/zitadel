package projection

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/repository/keypair"
)

func TestKeyProjection_reduces(t *testing.T) {
	type args struct {
		event func(t *testing.T) eventstore.Event
	}
	tests := []struct {
		name   string
		args   args
		reduce func(event eventstore.Event) (*handler.Statement, error)
		want   wantReduce
	}{
		{
			name: "reduceKeyPairAdded",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(keypair.AddedEventType),
					keypair.AggregateType,
					keypairAddedEventData(time.Now().Add(time.Hour)),
				), keypair.AddedEventMapper),
			},
			reduce: (&keyProjection{encryptionAlgorithm: crypto.CreateMockEncryptionAlg(gomock.NewController(t))}).reduceKeyPairAdded,
			want: wantReduce{
				projection:       KeyProjectionTable,
				aggregateType:    eventstore.AggregateType("key_pair"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO zitadel.projections.keys (id, creation_date, change_date, resource_owner, sequence, algorithm, use) VALUES ($1, $2, $3, $4, $5, $6, $7)",
							expectedArgs: []interface{}{
								"agg-id",
								anyArg{},
								anyArg{},
								"ro-id",
								uint64(15),
								"algorithm",
								domain.KeyUsageSigning,
							},
						},
						{
							expectedStmt: "INSERT INTO zitadel.projections.keys_private (id, expiry, key) VALUES ($1, $2, $3)",
							expectedArgs: []interface{}{
								"agg-id",
								anyArg{},
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("privateKey"),
								},
							},
						},
						{
							expectedStmt: "INSERT INTO zitadel.projections.keys_public (id, expiry, key) VALUES ($1, $2, $3)",
							expectedArgs: []interface{}{
								"agg-id",
								anyArg{},
								[]byte("publicKey"),
							},
						},
					},
				},
			},
		},
		{
			name: "reduceKeyPairAdded expired",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(keypair.AddedEventType),
					keypair.AggregateType,
					keypairAddedEventData(time.Now().Add(-time.Hour)),
				), keypair.AddedEventMapper),
			},
			reduce: (&keyProjection{}).reduceKeyPairAdded,
			want: wantReduce{
				projection:       KeyProjectionTable,
				aggregateType:    eventstore.AggregateType("key_pair"),
				sequence:         15,
				previousSequence: 10,
				executer:         &testExecuter{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := baseEvent(t)
			got, err := tt.reduce(event)
			if !errors.IsErrorInvalidArgument(err) {
				t.Errorf("no wrong event mapping: %v, got: %v", err, got)
			}

			event = tt.args.event(t)
			got, err = tt.reduce(event)
			assertReduce(t, got, err, tt.want)
		})
	}
}

func keypairAddedEventData(t time.Time) []byte {
	return []byte(`{"algorithm": "algorithm", "usage": 0, "privateKey": {"key": {"cryptoType": 0, "algorithm": "enc", "keyID": "id", "crypted": "cHJpdmF0ZUtleQ=="}, "expiry": "` + t.Format(time.RFC3339) + `"}, "publicKey": {"key": {"cryptoType": 0, "algorithm": "enc", "keyID": "id", "crypted": "cHVibGljS2V5"}, "expiry": "` + t.Format(time.RFC3339) + `"}}`)
}
