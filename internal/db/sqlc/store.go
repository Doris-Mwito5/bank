package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

type Store struct {
	db *sql.DB
	*Queries
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountID int64   `json:"from_account_id"`
	ToAccountID   int64   `json:"to_account_id"`
	Amount        float64 `json:"amount"`
}

type TransferResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_Entry"`
	ToEntry     Entry    `json:"to_Entry"`
}

type ctxKey string

var txKey = ctxKey("tx")

func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferResult, error) {

	var result TransferResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		txName := ctx.Value(txKey)
		log.Println(txName, "create transfer")
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        int64(arg.Amount),
		})

		if err != nil {
			return err
		}

		log.Println(txName, "create entry 1")
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -int64(arg.Amount),
		})
		if err != nil {
			return err
		}

		log.Println(txName, "create entry 2")
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    +int64(arg.Amount),
		})
		if err != nil {
			return err
		}

		//update account balance
		log.Println(txName, "create transferget account1 for update")
		account1, err := q.GetAccountForUpdate(ctx, arg.FromAccountID)
		if err != nil {
			return err
		}

		log.Println(txName, "update account1 balance")
		result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID:      arg.FromAccountID,
			Balance: account1.Balance - int64(arg.Amount),
		})
		if err != nil {
			return err
		}

		log.Println(txName, "get account 2 for update")
		account2, err := q.GetAccountForUpdate(ctx, arg.ToAccountID)
		if err != nil {
			return err
		}

		log.Println(txName, "update account 2 balance")
		result.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID:      arg.ToAccountID,
			Balance: account2.Balance + int64(arg.Amount),
		})
		if err != nil {
			return err
		}
		return nil

	})

	return result, err
}
