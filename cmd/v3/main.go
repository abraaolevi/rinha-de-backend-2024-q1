package main

import (
	"context"
	"net/http"

	"github.com.br/abraaolevi/rinha-backend-2024/config"
	v3 "github.com.br/abraaolevi/rinha-backend-2024/internal/v3"
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

	handler := v3.NewHandler(conn)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /clientes/{id}/transacoes", handler.HandlePostTransactions)
	mux.HandleFunc("GET /clientes/{id}/extrato", handler.HandleGetStatements)

	if err := http.ListenAndServe(":3000", mux); err != nil {
		panic(err)
	}
}
