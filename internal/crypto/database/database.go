package database

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"

	caos_errs "github.com/caos/zitadel/internal/errors"

	"github.com/caos/zitadel/internal/crypto"
)

type database struct {
	client    *sql.DB
	masterKey string
	encrypt   func(key, masterKey string) (encryptedKey string, err error)
	decrypt   func(encryptedKey, masterKey string) (key string, err error)
}

const (
	EncryptionKeysTable  = "system.encryption_keys"
	encryptionKeysIDCol  = "id"
	encryptionKeysKeyCol = "key"
)

func NewKeyStorage(client *sql.DB, masterKey string) (*database, error) {
	if err := checkMasterKeyLength(masterKey); err != nil {
		return nil, err
	}
	return &database{
		client:    client,
		masterKey: masterKey,
		encrypt:   crypto.EncryptAESString,
		decrypt:   crypto.DecryptAESString,
	}, nil
}

func (d *database) ReadKeys() (crypto.Keys, error) {
	keys := make(map[string]string)
	stmt, args, err := sq.Select(encryptionKeysIDCol, encryptionKeysKeyCol).
		From(EncryptionKeysTable).
		ToSql()
	if err != nil {
		return nil, caos_errs.ThrowInternal(err, "", "unable to read keys")
	}
	rows, err := d.client.Query(stmt, args...)
	if err != nil {
		return nil, caos_errs.ThrowInternal(err, "", "unable to read keys")
	}
	for rows.Next() {
		var id, encryptionKey string
		err = rows.Scan(&id, &encryptionKey)
		if err != nil {
			return nil, caos_errs.ThrowInternal(err, "", "unable to read keys")
		}
		key, err := d.decrypt(encryptionKey, d.masterKey)
		if err != nil {
			if err := rows.Close(); err != nil {
				return nil, caos_errs.ThrowInternal(err, "", "unable to close rows")
			}
			return nil, caos_errs.ThrowInternal(err, "", "unable to decrypt key")
		}
		keys[id] = key
	}
	if err := rows.Close(); err != nil {
		return nil, caos_errs.ThrowInternal(err, "", "unable to close rows")
	}
	return keys, err
}

func (d *database) ReadKey(id string) (*crypto.Key, error) {
	stmt, args, err := sq.Select(encryptionKeysKeyCol).
		From(EncryptionKeysTable).
		Where(sq.Eq{encryptionKeysIDCol: id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, caos_errs.ThrowInternal(err, "", "unable to read key")
	}
	row := d.client.QueryRow(stmt, args...)
	if err != nil {
		return nil, caos_errs.ThrowInternal(err, "", "unable to read key")
	}
	var encryptionKey string
	err = row.Scan(&encryptionKey)
	if err != nil {
		return nil, caos_errs.ThrowInternal(err, "", "unable to read key")
	}
	key, err := d.decrypt(encryptionKey, d.masterKey)
	if err != nil {
		return nil, caos_errs.ThrowInternal(err, "", "unable to decrypt key")
	}
	return &crypto.Key{
		ID:    id,
		Value: key,
	}, nil
}

func (d *database) CreateKeys(keys ...*crypto.Key) error {
	insert := sq.Insert(EncryptionKeysTable).
		Columns(encryptionKeysIDCol, encryptionKeysKeyCol).PlaceholderFormat(sq.Dollar)
	for _, key := range keys {
		encryptionKey, err := d.encrypt(key.Value, d.masterKey)
		if err != nil {
			return caos_errs.ThrowInternal(err, "", "unable to encrypt key")
		}
		insert = insert.Values(key.ID, encryptionKey)
	}
	stmt, args, err := insert.ToSql()
	if err != nil {
		return caos_errs.ThrowInternal(err, "", "unable to insert new keys")
	}
	tx, err := d.client.Begin()
	if err != nil {
		return caos_errs.ThrowInternal(err, "", "unable to insert new keys")
	}
	_, err = tx.Exec(stmt, args...)
	if err != nil {
		tx.Rollback()
		return caos_errs.ThrowInternal(err, "", "unable to insert new keys")
	}
	err = tx.Commit()
	if err != nil {
		return caos_errs.ThrowInternal(err, "", "unable to insert new keys")
	}
	return nil
}

func checkMasterKeyLength(masterKey string) error {
	if length := len([]byte(masterKey)); length != 32 {
		return caos_errs.ThrowInternalf(nil, "", "masterkey must be 32 bytes, but is %d", length)
	}
	return nil
}
