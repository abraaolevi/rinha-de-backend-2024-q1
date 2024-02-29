package http

import (
	"context"
	"net/http"
	"strconv"

	"github.com.br/abraaolevi/rinha-backend-2024/internal"
	"github.com.br/abraaolevi/rinha-backend-2024/internal/account"
	"github.com.br/abraaolevi/rinha-backend-2024/internal/database"
)

var service account.Service

func Configure() {
	service = account.Service{
		Repo: &account.RepositoryPostgres{
			Conn: database.Conn,
			Ctx:  context.Background(),
		},
	}
}

func getPathId(r *http.Request) (uint64, error) {
	strId := r.PathValue("id")
	return strconv.ParseUint(strId, 10, 64)
}

func CreateTransaction(w http.ResponseWriter, r *http.Request) {
	accountId, err := getPathId(r)
	if err != nil {
		JsonError(w, account.ErrAccountNotFound)
		return
	}

	if !service.HasAccount(accountId) {
		JsonError(w, account.ErrAccountNotFound)
		return
	}

	transaction, err := internal.NewTransactionFromRequest(r)
	if err != nil {
		JsonError(w, err)
		return
	}

	acc, err := service.CreateTransaction(accountId, *transaction)
	if err != nil {
		JsonError(w, err)
		return
	}

	JsonResponse(w, acc, http.StatusOK)
}

func GetStatement(w http.ResponseWriter, r *http.Request) {
	accountId, err := getPathId(r)
	if err != nil {
		JsonError(w, account.ErrAccountNotFound)
		return
	}

	acc, err := service.GetAccount(accountId)
	if err != nil {
		JsonError(w, account.ErrAccountNotFound)
		return
	}

	stm, err := service.GetStatement(accountId, acc)
	if err != nil {
		JsonError(w, err)
		return
	}

	JsonResponse(w, stm, http.StatusOK)
}
