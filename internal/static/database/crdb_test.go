package database

import (
	"bytes"
	"context"
	"database/sql"
	"database/sql/driver"
	"io"
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"

	db_mock "github.com/zitadel/zitadel/internal/database/mock"
	"github.com/zitadel/zitadel/internal/static"
)

var (
	testNow = time.Now()
)

const (
	createObjectStmt = "INSERT INTO system.assets" +
		" (instance_id,resource_owner,name,asset_type,content_type,data,updated_at)" +
		" VALUES ($1,$2,$3,$4,$5,$6,$7)" +
		" ON CONFLICT (instance_id, resource_owner, name) DO UPDATE SET" +
		" content_type = $5, data = $6" +
		" RETURNING hash"
	removeObjectStmt = "DELETE FROM system.assets" +
		" WHERE instance_id = $1" +
		" AND name = $2" +
		" AND resource_owner = $3"
	removeObjectsStmt = "DELETE FROM system.assets" +
		" WHERE asset_type = $1" +
		" AND instance_id = $2" +
		" AND resource_owner = $3"
	removeInstanceObjectsStmt = "DELETE FROM system.assets" +
		" WHERE instance_id = $1"
)

func Test_crdbStorage_CreateObject(t *testing.T) {
	type fields struct {
		client db
	}
	type args struct {
		ctx           context.Context
		instanceID    string
		location      string
		resourceOwner string
		name          string
		contentType   string
		objectType    static.ObjectType
		data          io.Reader
		objectSize    int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *static.Asset
		wantErr bool
	}{
		{
			"create ok",
			fields{
				client: prepareDB(t,
					expectQuery(
						createObjectStmt,
						[]string{
							"hash",
							"updated_at",
						},
						[][]driver.Value{
							{
								"md5Hash",
								testNow,
							},
						},
						"instanceID",
						"resourceOwner",
						"name",
						static.ObjectTypeUserAvatar.String(),
						"contentType",
						[]byte("test"),
						"now()",
					)),
			},
			args{
				ctx:           context.Background(),
				instanceID:    "instanceID",
				location:      "location",
				resourceOwner: "resourceOwner",
				name:          "name",
				contentType:   "contentType",
				objectType:    static.ObjectTypeUserAvatar,
				data:          bytes.NewReader([]byte("test")),
				objectSize:    4,
			},
			&static.Asset{
				InstanceID:   "instanceID",
				Name:         "name",
				Hash:         "md5Hash",
				Size:         4,
				LastModified: testNow,
				Location:     "location",
				ContentType:  "contentType",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &crdbStorage{
				client: tt.fields.client.db,
			}
			got, err := c.PutObject(tt.args.ctx, tt.args.instanceID, tt.args.location, tt.args.resourceOwner, tt.args.name, tt.args.contentType, tt.args.objectType, tt.args.data, tt.args.objectSize)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateObject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateObject() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_crdbStorage_RemoveObject(t *testing.T) {
	type fields struct {
		client db
	}
	type args struct {
		ctx           context.Context
		instanceID    string
		resourceOwner string
		name          string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"remove ok",
			fields{
				client: prepareDB(t,
					expectExec(
						removeObjectStmt,
						nil,
						"instanceID",
						"name",
						"resourceOwner",
					)),
			},
			args{
				ctx:           context.Background(),
				instanceID:    "instanceID",
				resourceOwner: "resourceOwner",
				name:          "name",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &crdbStorage{
				client: tt.fields.client.db,
			}
			err := c.RemoveObject(tt.args.ctx, tt.args.instanceID, tt.args.resourceOwner, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("RemoveObject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_crdbStorage_RemoveObjects(t *testing.T) {
	type fields struct {
		client db
	}
	type args struct {
		ctx           context.Context
		instanceID    string
		resourceOwner string
		objectType    static.ObjectType
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"remove ok",
			fields{
				client: prepareDB(t,
					expectExec(
						removeObjectsStmt,
						nil, static.ObjectTypeUserAvatar.String(),
						"instanceID",
						"resourceOwner",
					)),
			},
			args{
				ctx:           context.Background(),
				instanceID:    "instanceID",
				resourceOwner: "resourceOwner",
				objectType:    static.ObjectTypeUserAvatar,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &crdbStorage{
				client: tt.fields.client.db,
			}
			err := c.RemoveObjects(tt.args.ctx, tt.args.instanceID, tt.args.resourceOwner, tt.args.objectType)
			if (err != nil) != tt.wantErr {
				t.Errorf("RemoveObjects() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
func Test_crdbStorage_RemoveInstanceObjects(t *testing.T) {
	type fields struct {
		client db
	}
	type args struct {
		ctx        context.Context
		instanceID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"remove ok",
			fields{
				client: prepareDB(t,
					expectExec(
						removeInstanceObjectsStmt,
						nil,
						"instanceID",
					)),
			},
			args{
				ctx:        context.Background(),
				instanceID: "instanceID",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &crdbStorage{
				client: tt.fields.client.db,
			}
			err := c.RemoveInstanceObjects(tt.args.ctx, tt.args.instanceID)
			if (err != nil) != tt.wantErr {
				t.Errorf("RemoveInstanceObjects() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

type db struct {
	mock sqlmock.Sqlmock
	db   *sql.DB
}

func prepareDB(t *testing.T, expectations ...expectation) db {
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
		db:   client,
	}
}

type expectation func(m sqlmock.Sqlmock)

func expectExists(query string, value bool, args ...driver.Value) expectation {
	return func(m sqlmock.Sqlmock) {
		m.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(args...).WillReturnRows(m.NewRows([]string{"exists"}).AddRow(value))
	}
}

func expectQueryErr(query string, err error, args ...driver.Value) expectation {
	return func(m sqlmock.Sqlmock) {
		m.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(args...).WillReturnError(err)
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

func expectExec(stmt string, err error, args ...driver.Value) expectation {
	return func(m sqlmock.Sqlmock) {
		query := m.ExpectExec(regexp.QuoteMeta(stmt)).WithArgs(args...)
		if err != nil {
			query.WillReturnError(err)
			return
		}
		query.WillReturnResult(sqlmock.NewResult(1, 1))
	}
}

func expectBegin(err error) expectation {
	return func(m sqlmock.Sqlmock) {
		query := m.ExpectBegin()
		if err != nil {
			query.WillReturnError(err)
		}
	}
}

func expectCommit(err error) expectation {
	return func(m sqlmock.Sqlmock) {
		query := m.ExpectCommit()
		if err != nil {
			query.WillReturnError(err)
		}
	}
}

func expectRollback(err error) expectation {
	return func(m sqlmock.Sqlmock) {
		query := m.ExpectRollback()
		if err != nil {
			query.WillReturnError(err)
		}
	}
}
