package handler

import (
	"context"
	"database/sql"
	"database/sql/driver"
	_ "embed"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/shopspring/decimal"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/database/mock"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

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
					eventTimestamp: time.Now(),
					position:       decimal.NewFromInt(42),
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
					eventTimestamp: time.Now(),
					position:       decimal.NewFromInt(42),
				},
			},
			isErr: func(t *testing.T, err error) {
				if !errors.Is(err, zerrors.ThrowInternal(nil, "V2-FGEKi", "")) {
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
							"aggregate id",
							eventstore.AggregateType("aggregate type"),
							uint64(42),
							mock.AnyType[time.Time]{},
							decimal.NewFromInt(42),
							uint32(0),
						),
						mock.WithExecRowsAffected(1),
					),
				),
			},
			args: args{
				updatedState: &state{
					instanceID:     "instance",
					eventTimestamp: time.Now(),
					position:       decimal.NewFromInt(42),
					aggregateType:  "aggregate type",
					aggregateID:    "aggregate id",
					sequence:       42,
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
			tx, err := tt.fields.mock.DB.BeginTx(context.Background(), nil)
			if err != nil {
				t.Fatalf("unable to begin transaction: %v", err)
			}

			h := &Handler{
				projection: tt.fields.projection,
			}
			err = h.setState(tx, tt.args.updatedState)

			tt.isErr(t, err)
			tt.fields.mock.Assert(t)
		})
	}
}

func TestHandler_currentState(t *testing.T) {
	testTime := time.Now()
	type fields struct {
		projection Projection
		mock       *mock.SQLMock
	}
	type args struct {
		ctx context.Context
	}
	type want struct {
		currentState *state
		isErr        func(t *testing.T, err error)
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "connection done",
			fields: fields{
				projection: &projection{
					name: "projection",
				},
				mock: mock.NewSQLMock(t,
					mock.ExpectBegin(nil),
					mock.ExpectQuery(currentStateStmt,
						mock.WithQueryArgs(
							"instance",
							"projection",
						),
						mock.WithQueryErr(sql.ErrConnDone),
					),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance"),
			},
			want: want{
				isErr: func(t *testing.T, err error) {
					if !errors.Is(err, sql.ErrConnDone) {
						t.Errorf("unexpected error, want: %v, got: %v", sql.ErrConnDone, err)
					}
				},
			},
		},
		{
			name: "state locked",
			fields: fields{
				projection: &projection{
					name: "projection",
				},
				mock: mock.NewSQLMock(t,
					mock.ExpectBegin(nil),
					mock.ExpectQuery(currentStateStmt,
						mock.WithQueryArgs(
							"instance",
							"projection",
						),
						mock.WithQueryErr(&pgconn.PgError{Code: "55P03"}),
					),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance"),
			},
			want: want{
				isErr: func(t *testing.T, err error) {
					pgErr := new(pgconn.PgError)
					if !errors.As(err, &pgErr) {
						t.Errorf("error should be PgErr but was %T", err)
						return
					}
					if pgErr.Code != "55P03" {
						t.Errorf("expected code 55P03 got: %s", pgErr.Code)
					}
				},
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
					mock.ExpectQuery(currentStateStmt,
						mock.WithQueryArgs(
							"instance",
							"projection",
						),
						mock.WithQueryResult(
							[]string{"aggregate_id", "aggregate_type", "event_sequence", "event_date", "position", "offset"},
							[][]driver.Value{
								{
									"aggregate id",
									"aggregate type",
									int64(42),
									testTime,
									decimal.NewFromInt(42).String(),
									uint16(10),
								},
							},
						),
					),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance"),
			},
			want: want{
				currentState: &state{
					instanceID:     "instance",
					eventTimestamp: testTime,
					position:       decimal.NewFromInt(42),
					aggregateType:  "aggregate type",
					aggregateID:    "aggregate id",
					sequence:       42,
					offset:         10,
				},
			},
		},
	}
	for _, tt := range tests {
		if tt.want.isErr == nil {
			tt.want.isErr = func(t *testing.T, err error) {
				if err != nil {
					t.Error("expected no error got:", err)
				}
			}
		}
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				projection: tt.fields.projection,
			}

			tx, err := tt.fields.mock.DB.BeginTx(context.Background(), nil)
			if err != nil {
				t.Fatalf("unable to begin transaction: %v", err)
			}

			gotCurrentState, err := h.currentState(tt.args.ctx, tx)

			tt.want.isErr(t, err)
			if !reflect.DeepEqual(gotCurrentState, tt.want.currentState) {
				t.Errorf("Handler.currentState() gotCurrentState = %v, want %v", gotCurrentState, tt.want.currentState)
			}
			tt.fields.mock.Assert(t)
		})
	}
}
