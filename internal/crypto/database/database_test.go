package database

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/crypto"
	z_db "github.com/zitadel/zitadel/internal/database"
	db_mock "github.com/zitadel/zitadel/internal/database/mock"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func Test_database_ReadKeys(t *testing.T) {
	type fields struct {
		client    db
		masterKey string
		decrypt   func(encryptedKey, masterKey string) (key string, err error)
	}
	type res struct {
		keys crypto.Keys
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		res    res
	}{
		{
			"query fails, error",
			fields{
				client:    dbMock(t, expectQueryErr("SELECT id, key FROM system.encryption_keys", sql.ErrConnDone)),
				masterKey: "",
				decrypt:   nil,
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, sql.ErrConnDone)
				},
			},
		},
		{
			"decryption error",
			fields{
				client: dbMock(t, expectQueryScanErr(
					"SELECT id, key FROM system.encryption_keys",
					[]string{"id", "key"},
					[][]driver.Value{
						{
							"id1",
							"key1",
						},
					})),
				masterKey: "wrong key",
				decrypt: func(encryptedKey, masterKey string) (key string, err error) {
					return "", fmt.Errorf("wrong masterkey")
				},
			},
			res{
				err: zerrors.IsInternal,
			},
		},
		{
			"single key ok",
			fields{
				client: dbMock(t, expectQuery(
					"SELECT id, key FROM system.encryption_keys",
					[]string{"id", "key"},
					[][]driver.Value{
						{
							"id1",
							"key1",
						},
					})),
				masterKey: "masterKey",
				decrypt: func(encryptedKey, masterKey string) (key string, err error) {
					return encryptedKey, nil
				},
			},
			res{
				keys: crypto.Keys(map[string]string{"id1": "key1"}),
			},
		},
		{
			"multiple keys ok",
			fields{
				client: dbMock(t, expectQuery(
					"SELECT id, key FROM system.encryption_keys",
					[]string{"id", "key"},
					[][]driver.Value{
						{
							"id1",
							"key1",
						},
						{
							"id2",
							"key2",
						},
					})),
				masterKey: "masterKey",
				decrypt: func(encryptedKey, masterKey string) (key string, err error) {
					return encryptedKey, nil
				},
			},
			res{
				keys: crypto.Keys(map[string]string{"id1": "key1", "id2": "key2"}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Database{
				client:    tt.fields.client.db,
				masterKey: tt.fields.masterKey,
				decrypt:   tt.fields.decrypt,
			}
			got, err := d.ReadKeys()
			if tt.res.err == nil {
				assert.NoError(t, err)
			} else if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.keys, got)
			}
			if err := tt.fields.client.mock.ExpectationsWereMet(); err != nil {
				t.Error(err)
			}
		})
	}
}

func Test_database_ReadKey(t *testing.T) {
	type fields struct {
		client    db
		masterKey string
		decrypt   func(encryptedKey, masterKey string) (key string, err error)
	}
	type args struct {
		id string
	}
	type res struct {
		key *crypto.Key
		err func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"query fails, error",
			fields{
				client:    dbMock(t, expectQueryErr("SELECT key FROM system.encryption_keys WHERE id = $1", sql.ErrConnDone)),
				masterKey: "",
				decrypt:   nil,
			},
			args{
				id: "id1",
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, sql.ErrConnDone)
				},
			},
		},
		{
			"key not found err",
			fields{
				client: dbMock(t, expectQueryScanErr(
					"SELECT key FROM system.encryption_keys WHERE id = $1",
					nil,
					nil,
					"id1")),
				masterKey: "masterKey",
				decrypt: func(encryptedKey, masterKey string) (key string, err error) {
					return encryptedKey, nil
				},
			},
			args{
				id: "id1",
			},
			res{
				err: zerrors.IsInternal,
			},
		},
		{
			"decryption error",
			fields{
				client: dbMock(t, expectQueryScanErr(
					"SELECT key FROM system.encryption_keys WHERE id = $1",
					[]string{"key"},
					[][]driver.Value{
						{
							"key1",
						},
					},
					"id1",
				)),
				masterKey: "wrong key",
				decrypt: func(encryptedKey, masterKey string) (key string, err error) {
					return "", fmt.Errorf("wrong masterkey")
				},
			},
			args{
				id: "id1",
			},
			res{
				err: zerrors.IsInternal,
			},
		},
		{
			"key ok",
			fields{
				client: dbMock(t, expectQuery(
					"SELECT key FROM system.encryption_keys WHERE id = $1",
					[]string{"key"},
					[][]driver.Value{
						{
							"key1",
						},
					},
					"id1",
				)),
				masterKey: "masterKey",
				decrypt: func(encryptedKey, masterKey string) (key string, err error) {
					return encryptedKey, nil
				},
			},
			args{
				id: "id1",
			},
			res{
				key: &crypto.Key{
					ID:    "id1",
					Value: "key1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Database{
				client:    tt.fields.client.db,
				masterKey: tt.fields.masterKey,
				decrypt:   tt.fields.decrypt,
			}
			got, err := d.ReadKey(tt.args.id)
			if tt.res.err == nil {
				assert.NoError(t, err)
			} else if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.key, got)
			}
			if err := tt.fields.client.mock.ExpectationsWereMet(); err != nil {
				t.Error(err)
			}
		})
	}
}

func Test_database_CreateKeys(t *testing.T) {
	type fields struct {
		client    db
		masterKey string
		encrypt   func(key, masterKey string) (encryptedKey string, err error)
	}
	type args struct {
		keys []*crypto.Key
	}
	type res struct {
		err func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"encryption fails, error",
			fields{
				client:    dbMock(t),
				masterKey: "",
				encrypt: func(key, masterKey string) (encryptedKey string, err error) {
					return "", fmt.Errorf("encryption failed")
				},
			},
			args{
				keys: []*crypto.Key{
					{
						"id1",
						"key1",
					},
				},
			},
			res{
				err: zerrors.IsInternal,
			},
		},
		{
			"insert fails, error",
			fields{
				client: dbMock(t,
					expectBegin(nil),
					expectExec("INSERT INTO system.encryption_keys (id,key) VALUES ($1,$2)", sql.ErrTxDone),
					expectRollback(nil),
				),
				masterKey: "masterkey",
				encrypt: func(key, masterKey string) (encryptedKey string, err error) {
					return key, nil
				},
			},
			args{
				keys: []*crypto.Key{
					{
						"id1",
						"key1",
					},
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, sql.ErrTxDone)
				},
			},
		},
		{
			"single insert ok",
			fields{
				client: dbMock(t,
					expectBegin(nil),
					expectExec("INSERT INTO system.encryption_keys (id,key) VALUES ($1,$2)", nil, "id1", "key1"),
					expectCommit(nil),
				),
				masterKey: "masterkey",
				encrypt: func(key, masterKey string) (encryptedKey string, err error) {
					return key, nil
				},
			},
			args{
				keys: []*crypto.Key{
					{
						"id1",
						"key1",
					},
				},
			},
			res{
				err: nil,
			},
		},
		{
			"multiple insert ok",
			fields{
				client: dbMock(t,
					expectBegin(nil),
					expectExec("INSERT INTO system.encryption_keys (id,key) VALUES ($1,$2)", nil, "id1", "key1", "id2", "key2"),
					expectCommit(nil),
				),
				masterKey: "masterkey",
				encrypt: func(key, masterKey string) (encryptedKey string, err error) {
					return key, nil
				},
			},
			args{
				keys: []*crypto.Key{
					{
						"id1",
						"key1",
					},
					{
						"id2",
						"key2",
					},
				},
			},
			res{
				err: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Database{
				client:    tt.fields.client.db,
				masterKey: tt.fields.masterKey,
				encrypt:   tt.fields.encrypt,
			}
			err := d.CreateKeys(context.Background(), tt.args.keys...)
			if tt.res.err == nil {
				assert.NoError(t, err)
			} else if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v", err)
			}
			if err := tt.fields.client.mock.ExpectationsWereMet(); err != nil {
				t.Error(err)
			}
		})
	}
}

func Test_checkMasterKeyLength(t *testing.T) {
	type args struct {
		masterKey string
	}
	tests := []struct {
		name string
		args args
		err  func(error) bool
	}{
		{
			"invalid length",
			args{
				masterKey: "",
			},
			zerrors.IsInternal,
		},
		{
			"valid length",
			args{
				masterKey: "!themasterkeywhichis32byteslong!",
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkMasterKeyLength(tt.args.masterKey)
			if tt.err == nil {
				assert.NoError(t, err)
			} else if tt.err != nil && !tt.err(err) {
				t.Errorf("got wrong err: %v", err)
			}
		})
	}
}

type db struct {
	mock sqlmock.Sqlmock
	db   *z_db.DB
}

func dbMock(t *testing.T, expectations ...func(m sqlmock.Sqlmock)) db {
	t.Helper()
	client, mock, err := sqlmock.New(sqlmock.ValueConverterOption(new(db_mock.TypeConverter)))
	if err != nil {
		t.Fatalf("unable to create sql mock: %v", err)
	}
	for _, expectation := range expectations {
		expectation(mock)
	}
	return db{
		mock: mock,
		db:   &z_db.DB{DB: client},
	}
}

func expectQueryErr(query string, err error, args ...driver.Value) func(m sqlmock.Sqlmock) {
	return func(m sqlmock.Sqlmock) {
		m.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(args...).WillReturnError(err)
	}
}

func expectQueryScanErr(stmt string, cols []string, rows [][]driver.Value, args ...driver.Value) func(m sqlmock.Sqlmock) {
	return func(m sqlmock.Sqlmock) {
		q := m.ExpectQuery(regexp.QuoteMeta(stmt)).WithArgs(args...)
		result := m.NewRows(cols)
		count := uint64(len(rows))
		for _, row := range rows {
			if cols[len(cols)-1] == "count" {
				row = append(row, count)
			}
			result.AddRow(row...)
		}
		q.WillReturnRows(result)
		q.RowsWillBeClosed()
	}
}

func expectQuery(stmt string, cols []string, rows [][]driver.Value, args ...driver.Value) func(m sqlmock.Sqlmock) {
	return func(m sqlmock.Sqlmock) {
		q := m.ExpectQuery(regexp.QuoteMeta(stmt)).WithArgs(args...)
		result := m.NewRows(cols)
		count := uint64(len(rows))
		for _, row := range rows {
			if cols[len(cols)-1] == "count" {
				row = append(row, count)
			}
			result.AddRow(row...)
		}
		q.WillReturnRows(result)
		q.RowsWillBeClosed()
	}
}

func expectExec(stmt string, err error, args ...driver.Value) func(m sqlmock.Sqlmock) {
	return func(m sqlmock.Sqlmock) {
		query := m.ExpectExec(regexp.QuoteMeta(stmt)).WithArgs(args...)
		if err != nil {
			query.WillReturnError(err)
			return
		}
		query.WillReturnResult(sqlmock.NewResult(1, 1))
	}
}

func expectBegin(err error) func(m sqlmock.Sqlmock) {
	return func(m sqlmock.Sqlmock) {
		query := m.ExpectBegin()
		if err != nil {
			query.WillReturnError(err)
		}
	}
}

func expectCommit(err error) func(m sqlmock.Sqlmock) {
	return func(m sqlmock.Sqlmock) {
		query := m.ExpectCommit()
		if err != nil {
			query.WillReturnError(err)
		}
	}
}

func expectRollback(err error) func(m sqlmock.Sqlmock) {
	return func(m sqlmock.Sqlmock) {
		query := m.ExpectRollback()
		if err != nil {
			query.WillReturnError(err)
		}
	}
}
