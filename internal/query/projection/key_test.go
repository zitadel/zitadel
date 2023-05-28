package projection

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
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
				event: getEvent(
					testEvent(
						keypair.AddedEventType,
						keypair.AggregateType,
						keypairAddedEventData(domain.KeyUsageSigning, time.Now().Add(time.Hour)),
					), keypair.AddedEventMapper),
			},
			reduce: (&keyProjection{encryptionAlgorithm: crypto.CreateMockEncryptionAlg(gomock.NewController(t))}).reduceKeyPairAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("key_pair"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.keys4 (id, creation_date, change_date, resource_owner, instance_id, sequence, algorithm, use) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
							expectedArgs: []interface{}{
								"agg-id",
								anyArg{},
								anyArg{},
								"ro-id",
								"instance-id",
								uint64(15),
								"algorithm",
								domain.KeyUsageSigning,
							},
						},
						{
							expectedStmt: "INSERT INTO projections.keys4_private (id, instance_id, expiry, key) VALUES ($1, $2, $3, $4)",
							expectedArgs: []interface{}{
								"agg-id",
								"instance-id",
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
							expectedStmt: "INSERT INTO projections.keys4_public (id, instance_id, expiry, key) VALUES ($1, $2, $3, $4)",
							expectedArgs: []interface{}{
								"agg-id",
								"instance-id",
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
				event: getEvent(
					testEvent(
						keypair.AddedEventType,
						keypair.AggregateType,
						keypairAddedEventData(domain.KeyUsageSigning, time.Now().Add(-time.Hour)),
					), keypair.AddedEventMapper),
			},
			reduce: (&keyProjection{}).reduceKeyPairAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("key_pair"),
				sequence:      15,
				executer:      &testExecuter{},
			},
		},
		{
			name: "instance reduceInstanceRemoved",
			args: args{
				event: getEvent(
					testEvent(
						instance.InstanceRemovedEventType,
						instance.AggregateType,
						nil,
					), instance.InstanceRemovedEventMapper),
			},
			reduce: reduceInstanceRemovedHelper(KeyColumnInstanceID),
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.keys4 WHERE (instance_id = $1)",
							expectedArgs: []interface{}{
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceCertificateAdded",
			args: args{
				event: getEvent(
					testEvent(
						keypair.AddedCertificateEventType,
						keypair.AggregateType,
						certificateAddedEventData(domain.KeyUsageSAMLMetadataSigning, time.Now().Add(time.Hour)),
					), keypair.AddedCertificateEventMapper),
			},
			reduce: (&keyProjection{certEncryptionAlgorithm: crypto.CreateMockEncryptionAlg(gomock.NewController(t))}).reduceCertificateAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("key_pair"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.keys4_certificate (id, instance_id, expiry, certificate) VALUES ($1, $2, $3, $4)",
							expectedArgs: []interface{}{
								"agg-id",
								"instance-id",
								anyArg{},
								[]byte("privateKey"),
							},
						},
					},
				},
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
			assertReduce(t, got, err, KeyProjectionTable, tt.want)
		})
	}
}

func keypairAddedEventData(usage domain.KeyUsage, t time.Time) []byte {
	return []byte(`{"algorithm": "algorithm", "usage": ` + fmt.Sprintf("%d", usage) + `, "privateKey": {"key": {"cryptoType": 0, "algorithm": "enc", "keyID": "id", "crypted": "cHJpdmF0ZUtleQ=="}, "expiry": "` + t.Format(time.RFC3339) + `"}, "publicKey": {"key": {"cryptoType": 0, "algorithm": "enc", "keyID": "id", "crypted": "cHVibGljS2V5"}, "expiry": "` + t.Format(time.RFC3339) + `"}}`)
}

func certificateAddedEventData(usage domain.KeyUsage, t time.Time) []byte {
	return []byte(`{"algorithm": "algorithm", "usage": ` + fmt.Sprintf("%d", usage) + `, "certificate": {"key": {"cryptoType": 0, "algorithm": "enc", "keyID": "id", "crypted": "cHJpdmF0ZUtleQ=="}, "expiry": "` + t.Format(time.RFC3339) + `"}}`)
}
