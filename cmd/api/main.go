package main

import (
	"context"
	"net/http"

	"github.com.br/abraaolevi/rinha-backend-2024/config"
	"github.com.br/abraaolevi/rinha-backend-2024/internal/database"
	apphttp "github.com.br/abraaolevi/rinha-backend-2024/internal/http"
)

func main() {
	ctx := context.Background()
	cfg := config.ReadConfig()

	initializeDatabase(ctx, *cfg)

	mux := http.NewServeMux()
	initializeRoutes(mux)

	if err := http.ListenAndServe(":3000", mux); err != nil {
		panic(err)
	}
}

func initializeDatabase(ctx context.Context, cfg config.Config) {
	_, err := database.NewConnection(ctx, cfg.GetPostgresConnectionString())
	if err != nil {
		panic(err)
	}
	// defer conn.Close()
}

func initializeRoutes(mux *http.ServeMux) {
	apphttp.Configure()

	mux.HandleFunc("POST /clientes/{id}/transacoes", apphttp.CreateTransaction)
	mux.HandleFunc("GET /clientes/{id}/extrato", apphttp.GetStatement)
}
