package crdb

// import (
// 	"context"
// 	"database/sql"

// 	"github.com/caos/zitadel/internal/eventstore/handler"
// )

// type StatementHandler struct {
// 	tableName     string
// 	sequenceTable string
// 	lockTable     string
// 	client        *sql.DB
// }

// func NewStatementHandler(
// 	client *sql.DB,
// 	tableName,
// 	sequenceTable,
// 	lockTable string,
// ) StatementHandler {
// 	return StatementHandler{
// 		client:        client,
// 		tableName:     tableName,
// 		sequenceTable: sequenceTable,
// 		lockTable:     lockTable,
// 	}
// }

// func (h *StatementHandler) Lock() error {
// 	query := "INSERT INTO " + h.lockTable + " (table_name) VALUES ($1)"
// 	_, err := h.client.Exec(query, h.tableName)
// 	return err
// }

// func (h *StatementHandler) Unlock() error {
// 	query := "DELETE FROM " + h.lockTable + " WHERE table_name = $1"
// 	_, err := h.client.Exec(query, h.tableName)
// 	return err
// }

// func (h *StatementHandler) Update(ctx context.Context, stmts []handler.Statement) error {
// 	tx, err := h.client.BeginTx(ctx, nil)
// 	if err != nil {
// 		return err
// 	}

// 	for _, stmt := range stmts {
// 		currentSeq, err := stmt.CurrentSequence(tx, h.sequenceTable)
// 		if err != nil {
// 			tx.Rollback()
// 			return err
// 		}
// 		if stmt.PreviousSequence < currentSeq {
// 			continue
// 		}
// 		if stmt.PreviousSequence > currentSeq {
// 			break
// 		}
// 		if err = executeStmt(tx, stmt); err != nil {
// 			return err
// 		}
// 		if err = stmt.UpdateCurrentSequence(tx, h.sequenceTable); err != nil {
// 			tx.Rollback()
// 			return err
// 		}
// 	}

// 	return tx.Commit()
// }

// //executeStmt handles sql statements
// // the transaction is closed properly if an error occurres
// func executeStmt(tx *sql.Tx, stmt handler.Statement) error {
// 	_, err := tx.Query("SAVEPOINT push_stmt")
// 	if err != nil {
// 		tx.Rollback()
// 		return err
// 	}
// 	err = stmt.Execute(tx)
// 	if err != nil {
// 		_, err = tx.Query("ROLLBACK TO SAVEPOINT push_stmt")
// 		if err != nil {
// 			tx.Rollback()
// 			return err
// 		}
// 		err = tx.Commit()
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	_, err = tx.Query("RELEASE push_stmt")
// 	if err != nil {
// 		tx.Rollback()
// 		return err
// 	}
// 	return nil
// }
