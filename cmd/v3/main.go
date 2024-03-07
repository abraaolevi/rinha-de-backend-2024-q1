package main

import (
	"context"
	"net"
	"net/http"
	"os"

	v3 "github.com.br/abraaolevi/rinha-backend-2024/internal/v3"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()
	config, err := pgxpool.ParseConfig("user=admin password=123 host=/var/run/postgresql dbname=rinha sslmode=disable")
	if err != nil {
		panic(err)
	}

	conn, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	handler := v3.NewHandler(conn)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /clientes/{id}/transacoes", handler.HandlePostTransactions)
	mux.HandleFunc("GET /clientes/{id}/extrato", handler.HandleGetStatements)

	socketPath := "/app_tmp/rinha.sock"
	os.Remove(socketPath)

	server := &http.Server{
		Handler: mux,
	}

	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	if err := server.Serve(listener); err != nil {
		panic(err)
	}
}
