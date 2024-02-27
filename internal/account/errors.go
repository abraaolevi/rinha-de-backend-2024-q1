package account

import (
	"net/http"

	"github.com.br/abraaolevi/rinha-backend-2024/internal"
)

var ErrAccountNotFound = internal.NewError("account not found", http.StatusNotFound)
var ErrInsufficientBalance = internal.NewError("insufficient account balance", http.StatusUnprocessableEntity)
