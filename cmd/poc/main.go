package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"unicode/utf8"

	"github.com.br/abraaolevi/rinha-backend-2024/config"
	"github.com.br/abraaolevi/rinha-backend-2024/internal"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()
	cfg := config.ReadConfig()

	conn, err := pgxpool.New(ctx, cfg.GetPostgresConnectionString())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("POST /clientes/{id}/transacoes", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		var data struct {
			Amount      int64  `json:"valor"`
			Type        string `json:"tipo"`
			Description string `json:"descricao"`
		}

		err = json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		c := utf8.RuneCountInString(data.Description)
		if data.Amount < 1 || (data.Type != "c" && data.Type != "d") || c < 1 || c > 10 {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		var resp struct {
			message       string
			hasErr        bool
			balance       int64
			account_limit int64
		}

		err = conn.QueryRow(
			context.Background(),
			`select out_message, out_has_error, out_balance, out_account_limit 
				from update_balance($1, $2, $3, $4)`,
			id,
			data.Type,
			data.Amount,
			data.Description,
		).Scan(&resp.message, &resp.hasErr, &resp.balance, &resp.account_limit)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if resp.hasErr {
			if resp.message == "not_found" {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			if resp.message == "insufficient_limit" {
				w.WriteHeader(http.StatusUnprocessableEntity)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		j, err := json.Marshal(struct {
			Limit   int64 `json:"limite"`
			Balance int64 `json:"saldo"`
		}{
			Limit:   resp.account_limit,
			Balance: resp.balance,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(j)
	})

	mux.HandleFunc("GET /clientes/{id}/extrato", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		var acc struct {
			account_limit int64
			balance       int64
		}
		err = conn.QueryRow(
			ctx,
			"select account_limit, balance from accounts where id = $1",
			id,
		).Scan(&acc.account_limit, &acc.balance)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		rows, err := conn.Query(
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
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var transactions = make([]internal.TransactionStatement, 0)

		for rows.Next() {
			var t internal.TransactionStatement
			var c time.Time
			if err := rows.Scan(&t.Amount, &t.Type, &t.Description, &c); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			t.Date = c.Format(time.RFC3339Nano)
			transactions = append(transactions, t)
		}

		currentTime := time.Now().Format(time.RFC3339Nano)

		stm := internal.Statement{
			Balance: internal.BalanceStatement{
				Total: acc.balance,
				Date:  currentTime,
				Limit: acc.account_limit,
			},
			Transactions: transactions,
		}

		j, err := json.Marshal(stm)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(j)

	})

	if err := http.ListenAndServe(":3000", mux); err != nil {
		panic(err)
	}
}
