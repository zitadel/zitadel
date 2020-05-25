package spooler

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

type dbMock struct {
	db   *sql.DB
	mock sqlmock.Sqlmock
}

func mockDB(t *testing.T) *dbMock {
	mockDB := dbMock{}
	var err error
	mockDB.db, mockDB.mock, err = sqlmock.New()
	if err != nil {
		t.Fatalf("error occured while creating stub db %v", err)
	}

	mockDB.mock.MatchExpectationsInOrder(true)

	return &mockDB
}

func (db *dbMock) expectCommit() *dbMock {
	db.mock.ExpectCommit()

	return db
}

func (db *dbMock) expectRollback() *dbMock {
	db.mock.ExpectRollback()

	return db
}

func (db *dbMock) expectBegin() *dbMock {
	db.mock.ExpectBegin()

	return db
}

func (db *dbMock) expectSavepoint() *dbMock {
	db.mock.ExpectExec("SAVEPOINT").WillReturnResult(sqlmock.NewResult(1, 1))
	return db
}

func (db *dbMock) expectReleaseSavepoint() *dbMock {
	db.mock.ExpectExec("RELEASE SAVEPOINT").WillReturnResult(sqlmock.NewResult(1, 1))

	return db
}

func (db *dbMock) expectRenew(lockerID, view string, affectedRows int64) *dbMock {
	query := db.mock.
		ExpectExec(`INSERT INTO admin_api\.locks \(object_type, locker_id, locked_until\) VALUES \(\$1, \$2, now\(\)\+\$3\) ON CONFLICT \(object_type\) DO UPDATE SET locked_until = now\(\)\+\$4, locker_id = \$5 WHERE \(locks\.locked_until < now\(\) OR locks\.locker_id = \$6\) AND locks\.object_type = \$7`).
		WithArgs(view, lockerID, sqlmock.AnyArg(), sqlmock.AnyArg(), lockerID, lockerID, view).
		WillReturnResult(sqlmock.NewResult(1, 1))

	if affectedRows == 0 {
		query.WillReturnResult(sqlmock.NewResult(0, 0))
	} else {
		query.WillReturnResult(sqlmock.NewResult(1, affectedRows))
	}

	return db
}

func Test_locker_Renew(t *testing.T) {
	type fields struct {
		db *dbMock
	}
	type args struct {
		lockerID  string
		viewModel string
		waitTime  time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "renew succeeded",
			fields: fields{
				db: mockDB(t).
					expectBegin().
					expectSavepoint().
					expectRenew("locker", "view", 1).
					expectReleaseSavepoint().
					expectCommit(),
			},
			args:    args{lockerID: "locker", viewModel: "view", waitTime: 1 * time.Second},
			wantErr: false,
		},
		{
			name: "renew now rows updated",
			fields: fields{
				db: mockDB(t).
					expectBegin().
					expectSavepoint().
					expectRenew("locker", "view", 0).
					expectRollback(),
			},
			args:    args{lockerID: "locker", viewModel: "view", waitTime: 1 * time.Second},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &locker{
				dbClient: tt.fields.db.db,
			}
			if err := l.Renew(tt.args.lockerID, tt.args.viewModel, tt.args.waitTime); (err != nil) != tt.wantErr {
				t.Errorf("locker.Renew() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := tt.fields.db.mock.ExpectationsWereMet(); err != nil {
				t.Errorf("not all database expectations met: %v", err)
			}
		})
	}
}
