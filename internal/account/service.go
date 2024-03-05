package account

import (
	"context"
	"fmt"
	"time"

	"github.com.br/abraaolevi/rinha-backend-2024/internal"
	"github.com.br/abraaolevi/rinha-backend-2024/internal/database"
)

type Service struct {
	Repo Repository
}

func (as *Service) GetAccount(accountId uint64) (*internal.Account, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var account = internal.Account{}

	err := database.Conn.QueryRow(
		ctx,
		`SELECT account_limit, balance 
			FROM accounts 
			WHERE id = $1`,
		accountId,
	).Scan(&account.Limit, &account.Balance)

	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (as *Service) HasAccount(accountId uint64) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var exists bool
	err := database.Conn.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM accounts WHERE id = $1)", accountId).Scan(&exists)

	if err != nil {
		fmt.Println(err)
		return false
	}

	return exists
}

func (as *Service) CreateTransaction(accountId uint64, transaction internal.Transaction) (*internal.Account, error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var account internal.Account
	var err error

	if transaction.Type == "d" {
		err = database.Conn.QueryRow(
			ctx,
			`UPDATE accounts
				SET balance = balance - $1
				WHERE id = $2
					AND account_limit > abs(balance - $3)
				RETURNING account_limit, balance`,
			transaction.Amount,
			accountId,
			transaction.Amount,
		).Scan(&account.Limit, &account.Balance)
	} else {
		err = database.Conn.QueryRow(
			ctx,
			`UPDATE accounts
				SET balance = balance + $1
				WHERE id = $2
				RETURNING account_limit, balance`,
			transaction.Amount,
			accountId,
		).Scan(&account.Limit, &account.Balance)
	}

	if err != nil {
		return nil, ErrInsufficientBalance
	}

	_, err = database.Conn.Exec(
		ctx,
		`INSERT INTO transactions (account_id, amount, operation, description) 
			VALUES ($1, $2, $3, $4)`,
		accountId,
		transaction.Amount,
		transaction.Type,
		transaction.Description,
	)

	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (as *Service) CreateTransactionDeprecated(accountId uint64, transaction internal.Transaction) (*internal.Account, error) {
	tx, err := as.Repo.BeginTransaction()
	if err != nil {
		return nil, err
	}

	account, err := as.Repo.GetAccount(accountId, tx, true)
	if err != nil {
		return nil, as.Repo.EndTransactionWithError(tx, err)
	}

	limit := account.Limit + account.Balance
	if transaction.Type == "d" && transaction.Amount > limit {
		return nil, as.Repo.EndTransactionWithError(tx, ErrInsufficientBalance)
	}

	var amount = transaction.Amount
	if transaction.Type == "d" {
		amount = -transaction.Amount
	}

	account.Balance = account.Balance + amount

	if err := as.Repo.UpdateAccountBalance(accountId, account.Balance, tx); err != nil {
		return nil, as.Repo.EndTransactionWithError(tx, err)
	}

	if err := as.Repo.AddStatement(accountId, transaction, tx); err != nil {
		return nil, as.Repo.EndTransactionWithError(tx, err)
	}

	if err := as.Repo.EndTransaction(tx); err != nil {
		return nil, err
	}

	return account, nil
}

func (as *Service) GetStatement(accountId uint64, account *internal.Account) (*internal.Statement, error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	rows, err := database.Conn.Query(
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
		accountId,
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

	stm := internal.Statement{
		Balance: internal.BalanceStatement{
			Total: account.Balance,
			Date:  currentTime,
			Limit: account.Limit,
		},
		Transactions: transactions,
	}

	return &stm, nil
}

func (as *Service) GetStatementDeprecated(accountId uint64) (*internal.Statement, error) {
	tx, err := as.Repo.BeginTransaction()
	if err != nil {
		return nil, err
	}

	account, err := as.Repo.GetAccount(accountId, tx, false)
	if err != nil {
		return nil, as.Repo.EndTransactionWithError(tx, err)
	}

	transactions, err := as.Repo.GetStatementList(accountId, 10, tx)
	if err != nil {
		return nil, as.Repo.EndTransactionWithError(tx, err)
	}

	if err := as.Repo.EndTransaction(tx); err != nil {
		return nil, err
	}

	currentTime := time.Now().Format(time.RFC3339Nano)

	stm := internal.Statement{
		Balance: internal.BalanceStatement{
			Total: account.Balance,
			Date:  currentTime,
			Limit: account.Limit,
		},
		Transactions: transactions,
	}

	return &stm, nil
}
