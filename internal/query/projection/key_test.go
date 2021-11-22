package projection

import (
	"testing"
	"time"

	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/repository/keypair"
)

func TestKeyProjection_reduces(t *testing.T) {
	type args struct {
		event func(t *testing.T) eventstore.EventReader
	}
	tests := []struct {
		name   string
		args   args
		reduce func(event eventstore.EventReader) (*handler.Statement, error)
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
			reduce: (&KeyProjection{}).reduceKeyPairAdded,
			want: wantReduce{
				projection:       KeyProjectionTable,
				aggregateType:    eventstore.AggregateType("key_pair"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO zitadel.projections.keys (id, is_private, creation_date, change_date, resource_owner, sequence, algorithm, use, expiry, key) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)",
							expectedArgs: []interface{}{
								"agg-id",
								true,
								anyArg{},
								anyArg{},
								"ro-id",
								uint64(15),
								"algorithm",
								domain.KeyUsageSigning,
								anyArg{},
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "alg",
									KeyID:      "keyID",
									Crypted:    []byte("privateKey"),
								},
							},
						},
						{
							expectedStmt: "INSERT INTO zitadel.projections.keys (id, is_private, creation_date, change_date, resource_owner, sequence, algorithm, use, expiry, key) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)",
							expectedArgs: []interface{}{
								"agg-id",
								false,
								anyArg{},
								anyArg{},
								"ro-id",
								uint64(15),
								"algorithm",
								domain.KeyUsageSigning,
								anyArg{},
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "alg",
									KeyID:      "keyID",
									Crypted:    []byte("publicKey"),
								},
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
			reduce: (&KeyProjection{}).reduceKeyPairAdded,
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
	return []byte(`{"algorithm": "algorithm", "usage": 0, "privateKey": {"key": {"cryptoType": 0, "algorithm": "alg", "keyID": "keyID", "crypted": "cHJpdmF0ZUtleQ=="}, "expiry": "` + t.Format(time.RFC3339) + `"}, "publicKey": {"key": {"cryptoType": 0, "algorithm": "alg", "keyID": "keyID", "crypted": "cHVibGljS2V5"}, "expiry": "` + t.Format(time.RFC3339) + `"}}`)
}
