package main

import (
	"context"
	"net/http"

	"github.com.br/abraaolevi/rinha-backend-2024/config"
	"github.com.br/abraaolevi/rinha-backend-2024/internal/database"
	apphttp "github.com.br/abraaolevi/rinha-backend-2024/internal/v1/http"
)

func main() {
	ctx := context.Background()
	cfg := config.ReadConfig()

	conn, err := database.NewConnection(ctx, cfg.GetPostgresConnectionString())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	apphttp.Configure()

	mux := http.NewServeMux()
	mux.HandleFunc("POST /clientes/{id}/transacoes", apphttp.CreateTransaction)
	mux.HandleFunc("GET /clientes/{id}/extrato", apphttp.GetStatement)

	if err := http.ListenAndServe(":3000", mux); err != nil {
		panic(err)
	}
}
