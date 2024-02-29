package http

import (
	"encoding/json"
	"net/http"

	"github.com.br/abraaolevi/rinha-backend-2024/internal"
)

func JsonResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	j, err := json.Marshal(data)
	if err != nil {
		JsonError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(j)
}

func JsonError(w http.ResponseWriter, err error) {
	var code = http.StatusInternalServerError
	// var message = "interval server error"
	var message = err.Error()

	if ie, ok := err.(*internal.Error); ok {
		code = ie.Status
		message = ie.Error()
	}

	e := struct {
		Message string `json:"message"`
	}{
		Message: message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(e)
}
