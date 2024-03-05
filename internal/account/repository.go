package account

import (
	"context"
	"time"

	"github.com.br/abraaolevi/rinha-backend-2024/internal"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	AddStatement(accountId uint64, transaction internal.Transaction, tx *pgx.Tx) error
	BeginTransaction() (*pgx.Tx, error)
	EndTransaction(tx *pgx.Tx) error
	EndTransactionWithError(tx *pgx.Tx, err error) error
	GetAccount(id uint64, tx *pgx.Tx, forUpdate bool) (*internal.Account, error)
	GetStatementList(accountId uint64, count uint, tx *pgx.Tx) ([]internal.TransactionStatement, error)
	UpdateAccountBalance(id uint64, newBalance int64, tx *pgx.Tx) error
}

type RepositoryPostgres struct {
	Conn *pgxpool.Pool
	Ctx  context.Context
}

func (r *RepositoryPostgres) BeginTransaction() (*pgx.Tx, error) {
	tx, err := r.Conn.Begin(r.Ctx)
	if err != nil {
		return nil, err
	}
	return &tx, err
}

func (r *RepositoryPostgres) EndTransaction(tx *pgx.Tx) error {
	if err := (*tx).Commit(r.Ctx); err != nil {
		return err
	}
	return nil
}

func (r *RepositoryPostgres) EndTransactionWithError(tx *pgx.Tx, err error) error {
	if cerr := (*tx).Commit(r.Ctx); cerr != nil {
		return cerr
	}
	return err
}

func (r *RepositoryPostgres) GetAccount(id uint64, tx *pgx.Tx, forUpdate bool) (*internal.Account, error) {
	var account = internal.Account{}

	sql := `SELECT account_limit, balance
		FROM accounts
		WHERE id = $1 `
	if forUpdate {
		sql += "FOR UPDATE"
	}

	err := (*tx).QueryRow(r.Ctx, sql, id).Scan(&account.Limit, &account.Balance)

	if err == pgx.ErrNoRows {
		return nil, ErrAccountNotFound
	}

	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (r *RepositoryPostgres) UpdateAccountBalance(accountId uint64, newBalance int64, tx *pgx.Tx) error {
	_, err := (*tx).Exec(
		r.Ctx,
		`UPDATE accounts 
			SET balance = $1 
			WHERE id = $2`,
		newBalance,
		accountId,
	)
	return err
}

func (r *RepositoryPostgres) AddStatement(accountId uint64, transaction internal.Transaction, tx *pgx.Tx) error {
	_, err := (*tx).Exec(
		r.Ctx,
		`INSERT INTO transactions (account_id, amount, operation, description) 
			VALUES ($1, $2, $3, $4)`,
		accountId,
		transaction.Amount,
		transaction.Type,
		transaction.Description,
	)
	return err
}

func (r *RepositoryPostgres) GetStatementList(accountId uint64, count uint, tx *pgx.Tx) ([]internal.TransactionStatement, error) {
	rows, err := (*tx).Query(
		r.Ctx,
		`SELECT amount, operation, description, created_at 
			FROM transactions 
			WHERE account_id = $1 
			ORDER BY id DESC 
			LIMIT $2`,
		accountId,
		count,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []internal.TransactionStatement

	for rows.Next() {
		var t internal.TransactionStatement
		var c time.Time
		if err := rows.Scan(&t.Amount, &t.Type, &t.Description, &c); err != nil {
			return nil, err
		}
		t.Date = c.Format(time.RFC3339Nano)
		transactions = append(transactions, t)
	}

	return transactions, nil
}
