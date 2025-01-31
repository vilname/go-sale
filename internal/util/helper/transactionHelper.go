package helper

import (
	"context"
	"database/sql"
	"github.com/jackc/pgx/v5"
	"log/slog"
	"telemetry-sale/internal/config/storage"
)

type Transaction struct {
	Transaction pgx.Tx
	Ctx         context.Context
	Stmt        *sql.Stmt
}

func NewTransaction() *Transaction {
	return new(Transaction)
}

func (transaction *Transaction) StartTransaction() {
	ctx := context.Background()
	db := storage.GetDB()

	transactionBegin, err := db.BeginTx(ctx, pgx.TxOptions{})

	if err != nil {
		return
	}

	transaction.Transaction = transactionBegin
	transaction.Ctx = ctx
}

func (transaction *Transaction) SaveAll(query string, args []interface{}, prepareName string) error {
	query = query[0 : len(query)-1]
	stmt, err := transaction.Transaction.Prepare(transaction.Ctx, prepareName, query)

	if err != nil {
		return err
	}

	_, err = transaction.Transaction.Exec(transaction.Ctx, stmt.SQL, args...)

	if err != nil {
		return err
	}
	return nil
}

func (transaction *Transaction) MustRollback() {
	if err := transaction.Transaction.Rollback(transaction.Ctx); err != nil {
		slog.Error(err.Error())
		panic(err.Error())
	}
}
