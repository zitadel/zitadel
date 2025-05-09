package database

// import (
// 	"context"
// 	"fmt"

// 	"github.com/zitadel/zitadel/backend/handler"
// )

// func Begin[In, Out, NextOut any](ctx context.Context, beginner Beginner, opts *TransactionOptions) handler.Defer[In, Out, NextOut] {
// 	// func(ctx context.Context, in *VerifyEmail) (_ *VerifyEmail, _ func(context.Context, error) error, err error) {
// 	return func(handle handler.DeferrableHandle[In, Out], next handler.Handle[Out, NextOut]) handler.Handle[In, NextOut] {
// 		return func(ctx context.Context, in In) (out NextOut, err error) {
// 			tx, err := beginner.Begin(ctx, opts)
// 			if err != nil {
// 				return out, err
// 			}
// 			defer func() {
// 				if err != nil {
// 					rollbackErr := tx.Rollback(ctx)
// 					if rollbackErr != nil {
// 						err = fmt.Errorf("query failed: %w, rollback failed: %v", err, rollbackErr)
// 					}
// 				} else {
// 					err = tx.Commit(ctx)
// 				}
// 			}()
// 			return handle(ctx, in, tx)
// 		}
// 	}

// }

// type QueryExecutorSetter interface {
// 	SetQueryExecutor(QueryExecutor)
// }

// func Begin[In QueryExecutorSetter](ctx context.Context, beginner Beginner, in In) (_ In, _ func(context.Context, error) error, err error) {
// 	tx, err := beginner.Begin(ctx, nil)
// 	if err != nil {
// 		return in, nil, err
// 	}
// 	in.SetQueryExecutor(tx)
// 	return in, func(ctx context.Context, err error) error {
// 		err = tx.End(ctx, err)
// 		if err != nil {
// 			return err
// 		}
// 		in.SetQueryExecutor(beginner)
// 		return nil
// 	}, err
// }
