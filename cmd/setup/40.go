package setup

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"path"
	"strings"
	"text/template"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

// query filenames
const (
	fileInTxOrderType = "00_in_tx_order_type.sql"
	fileType          = "01_type.sql"
	fileFunc          = "02_func.sql"
)

var (
	//go:embed 40/*.sql
	initPushFunc embed.FS
)

type InitPushFunc struct {
	dbClient *database.DB
}

func (mig *InitPushFunc) Execute(ctx context.Context, _ eventstore.Event) (err error) {
	conn, err := mig.dbClient.Conn(ctx)
	if err != nil {
		return err
	}
	defer func() {
		closeErr := conn.Close()
		logging.OnError(closeErr).Debug("failed to release connection")
		// Force the pool to reopen connections to apply the new types
		mig.dbClient.Pool.Reset()
	}()
	statements, err := mig.prepareStatements(ctx)
	if err != nil {
		return err
	}
	for _, stmt := range statements {
		logging.WithFields("file", stmt.file, "migration", mig.String()).Info("execute statement")
		if _, err := conn.ExecContext(ctx, stmt.query); err != nil {
			return fmt.Errorf("%s %s: %w", mig.String(), stmt.file, err)
		}
	}

	return nil
}

func (mig *InitPushFunc) String() string {
	return "40_init_push_func_v4"
}

func (mig *InitPushFunc) prepareStatements(ctx context.Context) ([]statement, error) {
	funcTmpl, err := template.ParseFS(initPushFunc, mig.filePath(fileFunc))
	if err != nil {
		return nil, fmt.Errorf("prepare steps: %w", err)
	}
	typeName, err := mig.inTxOrderType(ctx)
	if err != nil {
		return nil, fmt.Errorf("prepare steps: %w", err)
	}
	var funcStep strings.Builder
	err = funcTmpl.Execute(&funcStep, struct {
		InTxOrderType string
	}{
		InTxOrderType: typeName,
	})
	if err != nil {
		return nil, fmt.Errorf("prepare steps: %w", err)
	}
	typeStatement, err := fs.ReadFile(initPushFunc, mig.filePath(fileType))
	if err != nil {
		return nil, fmt.Errorf("prepare steps: %w", err)
	}
	return []statement{
		{
			file:  fileType,
			query: string(typeStatement),
		},
		{
			file:  fileFunc,
			query: funcStep.String(),
		},
	}, nil
}

func (mig *InitPushFunc) inTxOrderType(ctx context.Context) (typeName string, err error) {
	query, err := fs.ReadFile(initPushFunc, mig.filePath(fileInTxOrderType))
	if err != nil {
		return "", fmt.Errorf("get in_tx_order_type: %w", err)
	}

	err = mig.dbClient.QueryRowContext(ctx, func(row *sql.Row) error {
		return row.Scan(&typeName)
	}, string(query))
	if err != nil {
		return "", fmt.Errorf("get in_tx_order_type: %w", err)
	}
	return typeName, nil
}

func (mig *InitPushFunc) filePath(fileName string) string {
	return path.Join("40", fileName)
}
