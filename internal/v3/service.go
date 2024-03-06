package v3

import (
	"context"
	"errors"
	"time"

	"github.com.br/abraaolevi/rinha-backend-2024/internal"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	Conn *pgxpool.Pool
}

type TransactionResult struct {
	message       string
	hasErr        bool
	balance       int64
	account_limit int64
}

func (s *Service) CreateTransaction(id int, data internal.Transaction) (*TransactionResult, error) {
	var resp TransactionResult

	err := s.Conn.QueryRow(
		context.Background(),
		`select out_message, out_has_error, out_balance, out_account_limit 
			from update_balance($1, $2, $3, $4)`,
		id,
		data.Type,
		data.Amount,
		data.Description,
	).Scan(&resp.message, &resp.hasErr, &resp.balance, &resp.account_limit)

	if err != nil {
		return nil, err
	}

	if resp.hasErr {
		if resp.message == "not_found" {
			return nil, internal.ErrAccountNotFound
		}
		if resp.message == "insufficient_limit" {
			return nil, internal.ErrInsufficientBalance
		}
		return nil, errors.New("unexpected error")
	}

	return &resp, nil
}

func (s *Service) GetStatements(id int) (*internal.Statement, error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var acc struct {
		account_limit int64
		balance       int64
	}
	err := s.Conn.QueryRow(
		ctx,
		"select account_limit, balance from accounts where id = $1",
		id,
	).Scan(&acc.account_limit, &acc.balance)
	if err != nil {
		return nil, internal.ErrAccountNotFound
	}

	rows, err := s.Conn.Query(
		ctx,
		`SELECT 
			amount, 
			operation, 
			description, 
			created_at
		FROM transactions
		WHERE account_id = $1 
		ORDER BY id DESC 
		LIMIT 10`,
		id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions = make([]internal.TransactionStatement, 0)

	for rows.Next() {
		var t internal.TransactionStatement
		var c time.Time
		if err := rows.Scan(&t.Amount, &t.Type, &t.Description, &c); err != nil {
			return nil, err
		}
		t.Date = c.Format(time.RFC3339Nano)
		transactions = append(transactions, t)
	}

	currentTime := time.Now().Format(time.RFC3339Nano)

	return &internal.Statement{
		Balance: internal.BalanceStatement{
			Total: acc.balance,
			Date:  currentTime,
			Limit: acc.account_limit,
		},
		Transactions: transactions,
	}, nil
}
