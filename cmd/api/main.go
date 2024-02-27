package main

import (
	"context"
	"net/http"

	"github.com.br/abraaolevi/rinha-backend-2024/internal/database"
	apphttp "github.com.br/abraaolevi/rinha-backend-2024/internal/http"
)

func main() {
	ctx := context.Background()
	initializeDatabase(ctx)

	mux := http.NewServeMux()
	initializeRoutes(mux)

	if err := http.ListenAndServe(":3000", mux); err != nil {
		panic(err)
	}
}

func initializeDatabase(ctx context.Context) {
	// connectionString := "postgresql://admin:123@db:5432/rinha"
	connectionString := "postgresql://admin:123@localhost:5432/rinha"
	_, err := database.NewConnection(ctx, connectionString)
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
