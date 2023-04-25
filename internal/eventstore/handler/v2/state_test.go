package handler

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"testing"
	"time"

	"github.com/zitadel/zitadel/internal/database/mock"
	errs "github.com/zitadel/zitadel/internal/errors"
)

func TestHandler_lockState(t *testing.T) {
	type fields struct {
		projection Projection
		mock       *mock.SQLMock
	}
	type args struct {
		instanceID string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		isErr  func(t *testing.T, err error)
	}{
		{
			name: "tx closed",
			fields: fields{
				projection: &projection{
					name: "projection",
				},
				mock: mock.NewSQLMock(t,
					mock.ExpectBegin(nil),
					mock.ExcpectExec(
						lockStateStmt,
						mock.WithExecArgs(
							"projection",
							"instance",
						),
						mock.WithExecErr(sql.ErrTxDone),
					),
				),
			},
			args: args{
				instanceID: "instance",
			},
			isErr: func(t *testing.T, err error) {
				if !errors.Is(err, sql.ErrTxDone) {
					t.Errorf("unexpected error, want: %v got: %v", sql.ErrTxDone, err)
				}
			},
		},
		{
			name: "no rows affeced",
			fields: fields{
				projection: &projection{
					name: "projection",
				},
				mock: mock.NewSQLMock(t,
					mock.ExpectBegin(nil),
					mock.ExcpectExec(
						lockStateStmt,
						mock.WithExecArgs(
							"projection",
							"instance",
						),
						mock.WithExecNoRowsAffected(),
					),
				),
			},
			args: args{
				instanceID: "instance",
			},
			isErr: func(t *testing.T, err error) {
				if !errors.Is(err, errs.ThrowInternal(nil, "V2-lpiK0", "")) {
					t.Errorf("unexpected error: want internal (V2lpiK0), got: %v", err)
				}
			},
		},
		{
			name: "rows affected",
			fields: fields{
				projection: &projection{
					name: "projection",
				},
				mock: mock.NewSQLMock(t,
					mock.ExpectBegin(nil),
					mock.ExcpectExec(
						lockStateStmt,
						mock.WithExecArgs(
							"projection",
							"instance",
						),
						mock.WithExecRowsAffected(1),
					),
				),
			},
			args: args{
				instanceID: "instance",
			},
		},
	}
	for _, tt := range tests {
		if tt.isErr == nil {
			tt.isErr = func(t *testing.T, err error) {
				if err != nil {
					t.Error("expected no error got:", err)
				}
			}
		}
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				projection: tt.fields.projection,
			}

			tx, err := tt.fields.mock.DB.Begin()
			if err != nil {
				t.Fatalf("unable to begin transaction: %v", err)
			}

			err = h.lockState(context.Background(), tx, tt.args.instanceID)
			tt.isErr(t, err)

			tt.fields.mock.Assert(t)
		})
	}
}

func TestHandler_updateLastUpdated(t *testing.T) {
	type fields struct {
		projection Projection
		mock       *mock.SQLMock
	}
	type args struct {
		updatedState *state
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		isErr  func(t *testing.T, err error)
	}{
		{
			name: "update fails",
			fields: fields{
				projection: &projection{
					name: "instance",
				},
				mock: mock.NewSQLMock(t,
					mock.ExpectBegin(nil),
					mock.ExcpectExec(updateStateStmt,
						mock.WithExecErr(sql.ErrTxDone),
					),
				),
			},
			args: args{
				updatedState: &state{
					instanceID:     "instance",
					aggregateType:  "aggregate type",
					aggregateID:    "aggregate id",
					eventTimestamp: time.Now(),
					eventSequence:  42,
				},
			},
			isErr: func(t *testing.T, err error) {
				if !errors.Is(err, sql.ErrTxDone) {
					t.Errorf("unexpected error, want: %v, got %v", sql.ErrTxDone, err)
				}
			},
		},
		{
			name: "no rows affected",
			fields: fields{
				projection: &projection{
					name: "instance",
				},
				mock: mock.NewSQLMock(t,
					mock.ExpectBegin(nil),
					mock.ExcpectExec(updateStateStmt,
						mock.WithExecNoRowsAffected(),
					),
				),
			},
			args: args{
				updatedState: &state{
					instanceID:     "instance",
					aggregateType:  "aggregate type",
					aggregateID:    "aggregate id",
					eventTimestamp: time.Now(),
					eventSequence:  42,
				},
			},
			isErr: func(t *testing.T, err error) {
				if !errors.Is(err, errs.ThrowInternal(nil, "V2-FGEKi", "")) {
					t.Errorf("unexpected error, want: %v, got %v", sql.ErrTxDone, err)
				}
			},
		},
		{
			name: "success",
			fields: fields{
				projection: &projection{
					name: "projection",
				},
				mock: mock.NewSQLMock(t,
					mock.ExpectBegin(nil),
					mock.ExcpectExec(updateStateStmt,
						mock.WithExecArgs(
							"projection",
							"instance",
							mock.AnyType[time.Time]{},
							"aggregate type",
							"aggregate id",
							uint64(42),
						),
						mock.WithExecRowsAffected(1),
					),
				),
			},
			args: args{
				updatedState: &state{
					instanceID:     "instance",
					aggregateType:  "aggregate type",
					aggregateID:    "aggregate id",
					eventTimestamp: time.Now(),
					eventSequence:  42,
				},
			},
		},
	}
	for _, tt := range tests {
		if tt.isErr == nil {
			tt.isErr = func(t *testing.T, err error) {
				if err != nil {
					t.Error("expected no error got:", err)
				}
			}
		}
		t.Run(tt.name, func(t *testing.T) {
			tx, err := tt.fields.mock.DB.Begin()
			if err != nil {
				t.Fatalf("unable to begin transaction: %v", err)
			}

			h := &Handler{
				projection: tt.fields.projection,
			}
			err = h.setState(context.Background(), tx, tt.args.updatedState)

			tt.isErr(t, err)
			tt.fields.mock.Assert(t)
		})
	}
}
