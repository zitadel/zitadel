package database

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/crypto"
	z_db "github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type Database struct {
	client    *z_db.DB
	masterKey string
	encrypt   func(key, masterKey string) (encryptedKey string, err error)
	decrypt   func(encryptedKey, masterKey string) (key string, err error)
}

const (
	EncryptionKeysTable  = "system.encryption_keys"
	encryptionKeysIDCol  = "id"
	encryptionKeysKeyCol = "key"
)

func NewKeyStorage(client *z_db.DB, masterKey string) (*Database, error) {
	if err := checkMasterKeyLength(masterKey); err != nil {
		return nil, err
	}
	return &Database{
		client:    client,
		masterKey: masterKey,
		encrypt:   crypto.EncryptAESString,
		decrypt:   crypto.DecryptAESString,
	}, nil
}

func (d *Database) ReadKeys() (crypto.Keys, error) {
	keys := make(map[string]string)
	stmt, args, err := sq.Select(encryptionKeysIDCol, encryptionKeysKeyCol).
		From(EncryptionKeysTable).
		ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "", "unable to read keys")
	}
	err = d.client.Query(func(rows *sql.Rows) error {
		for rows.Next() {
			var id, encryptionKey string
			err = rows.Scan(&id, &encryptionKey)
			if err != nil {
				return zerrors.ThrowInternal(err, "", "unable to read keys")
			}
			key, err := d.decrypt(encryptionKey, d.masterKey)
			if err != nil {
				return zerrors.ThrowInternal(err, "", "unable to decrypt key")
			}
			keys[id] = key
		}
		return nil
	}, stmt, args...)

	if err != nil {
		return nil, zerrors.ThrowInternal(err, "", "unable to read keys")
	}

	return keys, nil
}

func (d *Database) ReadKey(id string) (_ *crypto.Key, err error) {
	stmt, args, err := sq.Select(encryptionKeysKeyCol).
		From(EncryptionKeysTable).
		Where(sq.Eq{encryptionKeysIDCol: id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "", "unable to read key")
	}
	var key string
	err = d.client.QueryRow(func(row *sql.Row) error {
		var encryptionKey string
		err = row.Scan(&encryptionKey)
		if err != nil {
			return zerrors.ThrowInternal(err, "", "unable to read key")
		}
		key, err = d.decrypt(encryptionKey, d.masterKey)
		if err != nil {
			return zerrors.ThrowInternal(err, "", "unable to decrypt key")
		}
		return nil
	}, stmt, args...)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "", "unable to read key")
	}

	return &crypto.Key{
		ID:    id,
		Value: key,
	}, nil
}

func (d *Database) CreateKeys(ctx context.Context, keys ...*crypto.Key) error {
	insert := sq.Insert(EncryptionKeysTable).
		Columns(encryptionKeysIDCol, encryptionKeysKeyCol).PlaceholderFormat(sq.Dollar)
	for _, key := range keys {
		encryptionKey, err := d.encrypt(key.Value, d.masterKey)
		if err != nil {
			return zerrors.ThrowInternal(err, "", "unable to encrypt key")
		}
		insert = insert.Values(key.ID, encryptionKey)
	}
	stmt, args, err := insert.ToSql()
	if err != nil {
		return zerrors.ThrowInternal(err, "", "unable to insert new keys")
	}
	ctx, spanBeginTx := tracing.NewNamedSpan(ctx, "db.BeginTx")
	tx, err := d.client.BeginTx(ctx, nil)
	spanBeginTx.EndWithError(err)
	if err != nil {
		return zerrors.ThrowInternal(err, "", "unable to insert new keys")
	}
	_, err = tx.Exec(stmt, args...)
	if err != nil {
		tx.Rollback()
		return zerrors.ThrowInternal(err, "", "unable to insert new keys")
	}
	err = tx.Commit()
	if err != nil {
		return zerrors.ThrowInternal(err, "", "unable to insert new keys")
	}
	return nil
}

func checkMasterKeyLength(masterKey string) error {
	if length := len([]byte(masterKey)); length != 32 {
		return zerrors.ThrowInternalf(nil, "", "masterkey must be 32 bytes, but is %d", length)
	}
	return nil
}
