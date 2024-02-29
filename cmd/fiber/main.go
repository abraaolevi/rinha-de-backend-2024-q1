package main

import (
	"context"
	"net/http"
	"strconv"
	"time"
	"unicode/utf8"

	"github.com.br/abraaolevi/rinha-backend-2024/config"
	"github.com.br/abraaolevi/rinha-backend-2024/internal"
	"github.com/gofiber/fiber/v2"
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

	app := fiber.New()

	app.Post("/clientes/:id/transacoes", func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.SendStatus(http.StatusNotFound)
		}

		var data struct {
			Amount      int64  `json:"valor"`
			Type        string `json:"tipo"`
			Description string `json:"descricao"`
		}

		if err := c.BodyParser(&data); err != nil {
			return err
		}

		count := utf8.RuneCountInString(data.Description)
		if data.Amount < 1 || (data.Type != "c" && data.Type != "d") || count < 1 || count > 10 {
			return c.SendStatus(http.StatusUnprocessableEntity)
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
			return c.SendStatus(http.StatusInternalServerError)
		}

		if resp.hasErr {
			if resp.message == "not_found" {
				return c.SendStatus(http.StatusNotFound)
			}
			if resp.message == "insufficient_limit" {
				return c.SendStatus(http.StatusUnprocessableEntity)
			}
			return c.SendStatus(http.StatusInternalServerError)
		}

		return c.JSON(struct {
			Limit   int64 `json:"limite"`
			Balance int64 `json:"saldo"`
		}{
			Limit:   resp.account_limit,
			Balance: resp.balance,
		})
	})

	app.Get("/clientes/:id/extrato", func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.SendStatus(http.StatusNotFound)
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
			return c.SendStatus(http.StatusNotFound)
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
			return c.SendStatus(http.StatusInternalServerError)
		}
		defer rows.Close()

		var transactions = make([]internal.TransactionStatement, 0)

		for rows.Next() {
			var t internal.TransactionStatement
			var created time.Time
			if err := rows.Scan(&t.Amount, &t.Type, &t.Description, &created); err != nil {
				return c.SendStatus(http.StatusInternalServerError)
			}
			t.Date = created.Format(time.RFC3339Nano)
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

		return c.JSON(stm)
	})

	app.Listen(":3000")

}
