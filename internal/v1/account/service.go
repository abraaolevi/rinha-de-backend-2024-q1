package account

import (
	"time"

	"github.com.br/abraaolevi/rinha-backend-2024/internal"
)

type Service struct {
	Repo Repository
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
		return nil, as.Repo.EndTransactionWithError(tx, internal.ErrInsufficientBalance)
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
