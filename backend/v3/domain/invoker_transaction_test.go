package domain_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/zitadel/zitadel/backend/v3/domain"
	domainmock "github.com/zitadel/zitadel/backend/v3/domain/mock"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/dbmock"
)

type testTransactionalExecutor struct {
	domain.Transactional
	domain.Executor
}

func Test_transactionInvoker_Invoke(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		executor func(ctrl *gomock.Controller) domain.Executor
		db       func(ctrl *gomock.Controller) database.Pool
		wantErr  error
	}{
		{
			name: "non-transactional executor does not start a transaction",
			executor: func(ctrl *gomock.Controller) domain.Executor {
				mock := domainmock.NewMockExecutor(ctrl)
				mock.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil)
				return mock
			},
			db: func(ctrl *gomock.Controller) database.Pool {
				return dbmock.NewMockPool(ctrl)
			},
			wantErr: nil,
		},
		{
			name: "transactional executor starts a transaction execution successful",
			executor: func(ctrl *gomock.Controller) domain.Executor {
				transactional := domainmock.NewMockTransactional(ctrl)
				executor := domainmock.NewMockExecutor(ctrl)
				executor.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil)
				return &testTransactionalExecutor{
					Transactional: transactional,
					Executor:      executor,
				}
			},
			db: func(ctrl *gomock.Controller) database.Pool {
				pool := dbmock.NewMockPool(ctrl)
				tx := dbmock.NewMockTransaction(ctrl)
				pool.EXPECT().Begin(gomock.Any(), gomock.Any()).Return(tx, nil)
				tx.EXPECT().End(gomock.Any(), nil).Return(nil)
				return pool
			},
			wantErr: nil,
		},
		{
			name: "transactional executor starts a transaction execution failed",
			executor: func(ctrl *gomock.Controller) domain.Executor {
				transactional := domainmock.NewMockTransactional(ctrl)
				executor := domainmock.NewMockExecutor(ctrl)
				executor.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(assert.AnError)
				return &testTransactionalExecutor{
					Transactional: transactional,
					Executor:      executor,
				}
			},
			db: func(ctrl *gomock.Controller) database.Pool {
				pool := dbmock.NewMockPool(ctrl)
				tx := dbmock.NewMockTransaction(ctrl)
				pool.EXPECT().Begin(gomock.Any(), gomock.Any()).Return(tx, nil)
				tx.EXPECT().End(gomock.Any(), assert.AnError).DoAndReturn(func(ctx context.Context, err error) error { return err })
				return pool
			},
			wantErr: assert.AnError,
		},
	}
	for _, tt := range tests {
            t.Parallel()
			ctrl := gomock.NewController(t)
			ctrl := gomock.NewController(t)
			opts := &domain.InvokeOpts{}
			domain.WithQueryExecutor(tt.db(ctrl))(opts)

			invoker := domain.NewTransactionInvoker(nil)
			gotErr := invoker.Invoke(t.Context(), tt.executor(ctrl), opts)
			require.ErrorIs(t, gotErr, tt.wantErr)
		})
	}
}
