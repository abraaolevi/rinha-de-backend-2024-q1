package internal

import (
	"encoding/json"
	"net/http"
	"unicode/utf8"
)

var ErrUnmarshal = NewError("invalid body request", http.StatusUnprocessableEntity)
var ErrInvalidTransactionAmount = NewError("invalid Transaction Value. It should be a positive integer > 0", http.StatusUnprocessableEntity)
var ErrInvalidTransactionType = NewError("invalid Transaction Type. It should be \"c\" for Credit or \"d\" for Debit", http.StatusUnprocessableEntity)
var ErrInvalidTransactionDescription = NewError("invalid Transaction Description. It should be between 1 and 10 characters", http.StatusUnprocessableEntity)

type Transaction struct {
	Amount      int64  `json:"valor"`     // deve ser um número inteiro positivo que representa centavos. Por exemplo, R$ 10 são 1000 centavos.
	Type        string `json:"tipo"`      // deve ser apenas c para crédito ou d para débito.
	Description string `json:"descricao"` // deve ser uma string de 1 a 10 caracteres.
}

func NewTransactionFromRequest(r *http.Request) (*Transaction, error) {
	var transaction Transaction
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		return nil, ErrUnmarshal
	}

	if err := transaction.Validate(); err != nil {
		return nil, err
	}

	return &transaction, nil
}

func (t Transaction) Validate() error {
	if t.Amount < 1 {
		return ErrInvalidTransactionAmount
	}

	if t.Type != "c" && t.Type != "d" {
		return ErrInvalidTransactionType
	}

	c := utf8.RuneCountInString(t.Description)
	if c < 1 || c > 10 {
		return ErrInvalidTransactionDescription
	}

	return nil
}

type Account struct {
	Limit   int64 `json:"limite"` // deve ser o limite cadastrado do cliente.
	Balance int64 `json:"saldo"`  // deve ser o novo saldo após a conclusão da transação.
}

type Statement struct {
	Balance      BalanceStatement       `json:"saldo"`
	Transactions []TransactionStatement `json:"ultimas_transacoes"` // lista ordenada por data/hora das transações de forma decrescente contendo até as 10 últimas
}

type BalanceStatement struct {
	Total int64  `json:"total"`        // saldo total atual do cliente (não apenas das últimas transações seguintes exibidas).
	Date  string `json:"data_extrato"` // data/hora da consulta do extrato.
	Limit int64  `json:"limite"`       // limite cadastrado do cliente.
}

type TransactionStatement struct {
	Amount      int64  `json:"valor"`        // o valor da transação.
	Type        string `json:"tipo"`         // deve ser apenas c para crédito ou d para débito.
	Description string `json:"descricao"`    // descrição
	Date        string `json:"realizada_em"` // data/hora da realização da transação.
}
