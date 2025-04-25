package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"go-masters/task2/internal/errs"
	"go-masters/task2/internal/rates"
)

type RateRequest struct {
	From string `json:"from"`
	To   string `json:"to"`
}

func Rate(w http.ResponseWriter, r *http.Request) {

	var rReq = RateRequest{}

	err := json.NewDecoder(r.Body).Decode(&rReq)
	if err != nil {
		WriteError(w, errs.NewErrBadRequest(err.Error()))
		return
	}

	if rReq.From == "" || rReq.To == "" {
		WriteError(w, errs.NewErrBadRequest("currency code is empty"))
		return
	}

	pair := strings.ToUpper(rReq.From) + "-" + strings.ToUpper(rReq.To)
	rate, ok := rates.Rates[pair]
	if !ok {
		WriteError(w, errs.NewErrNotFound(pair))
		return
	}

	response := fmt.Sprintf(" %s: %.2f", pair, rate)
	w.Header().Set("Content-Type", "application/json")

	w.Write([]byte(response))
}

func WriteError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	errMassage := err.Error()
	var errStatus int
	switch err.(type) {
	case *errs.ErrBadRequest:
		errStatus = http.StatusBadRequest
	case *errs.ErrNotFound:
		errStatus = http.StatusNotFound
	default:
		errStatus = http.StatusInternalServerError
	}

	http.Error(w, errMassage, errStatus)
}
