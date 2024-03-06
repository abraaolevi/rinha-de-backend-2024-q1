package v3

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com.br/abraaolevi/rinha-backend-2024/internal"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	service *Service
}

func NewHandler(pool *pgxpool.Pool) *Handler {
	return &Handler{
		service: &Service{
			Conn: pool,
		},
	}
}

func (h *Handler) HandlePostTransactions(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var data internal.Transaction
	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	if err := data.Validate(); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	resp, err := h.service.CreateTransaction(id, data)

	if err != nil {
		if err == internal.ErrAccountNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if err == internal.ErrInsufficientBalance {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	j, err := json.Marshal(internal.Account{
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
}

func (h *Handler) HandleGetStatements(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	stmt, err := h.service.GetStatements(id)
	if err != nil {
		if err == internal.ErrAccountNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	j, err := json.Marshal(stmt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}
