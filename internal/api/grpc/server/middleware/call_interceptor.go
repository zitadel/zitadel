package middleware

import (
	"context"
	"database/sql"

	"github.com/zitadel/logging"
	"google.golang.org/grpc"

	"github.com/zitadel/zitadel/internal/api/call"
	"github.com/zitadel/zitadel/internal/api/transaction"
	"github.com/zitadel/zitadel/internal/database"
)

func CallDurationHandler() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ctx = call.WithTimestamp(ctx)
		return handler(ctx, req)
	}
}

func BeginMiddlewareTx(client *database.DB) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		tx, err := client.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted, ReadOnly: true})
		if err != nil {
			return nil, err
		}
		defer func() {
			if tx := transaction.FromContext(ctx); tx != nil {
				tx.Commit()
				ctx = transaction.WithTx(ctx, nil)
			}
		}()
		ctx = transaction.WithTx(ctx, tx)
		// ctx = transaction.WithTx(ctx, tx)
		ctx = call.WithTimestamp(ctx)
		return handler(ctx, req)
	}
}

func CloseMiddlewareTx() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		tx := transaction.FromContext(ctx)
		if tx == nil {
			return handler(ctx, req)
		}
		logging.OnError(tx.Commit()).Debug("commit failed")

		ctx = transaction.WithTx(ctx, nil)
		return handler(ctx, req)
	}
}
