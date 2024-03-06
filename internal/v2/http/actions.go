package http

import (
	"net/http"
	"strconv"

	"github.com.br/abraaolevi/rinha-backend-2024/internal"
	"github.com.br/abraaolevi/rinha-backend-2024/internal/v2/account"
)

var service account.Service

func Configure() {
	service = account.Service{}
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

	if !service.HasAccount(accountId) {
		internal.JsonError(w, internal.ErrAccountNotFound)
		return
	}

	transaction, err := internal.NewTransactionFromRequest(r)
	if err != nil {
		internal.JsonError(w, err)
		return
	}

	acc, err := service.CreateTransaction(accountId, *transaction)
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

	acc, err := service.GetAccount(accountId)
	if err != nil {
		internal.JsonError(w, internal.ErrAccountNotFound)
		return
	}

	stm, err := service.GetStatement(accountId, acc)
	if err != nil {
		internal.JsonError(w, err)
		return
	}

	internal.JsonResponse(w, stm, http.StatusOK)
}
