package http

import (
	"context"
	"net/http"
	"strconv"

	"github.com.br/abraaolevi/rinha-backend-2024/internal"
	"github.com.br/abraaolevi/rinha-backend-2024/internal/database"
	"github.com.br/abraaolevi/rinha-backend-2024/internal/v1/account"
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
		internal.JsonError(w, internal.ErrAccountNotFound)
		return
	}

	transaction, err := internal.NewTransactionFromRequest(r)
	if err != nil {
		internal.JsonError(w, err)
		return
	}

	acc, err := service.CreateTransactionDeprecated(accountId, *transaction)
	if err != nil {
		internal.JsonError(w, err)
		return
	}

	internal.JsonResponse(w, acc, http.StatusOK)
}

func GetStatement(w http.ResponseWriter, r *http.Request) {
	accountId, err := getPathId(r)
	if err != nil {
		internal.JsonError(w, internal.ErrAccountNotFound)
		return
	}

	stm, err := service.GetStatementDeprecated(accountId)
	if err != nil {
		internal.JsonError(w, err)
		return
	}

	internal.JsonResponse(w, stm, http.StatusOK)
}
