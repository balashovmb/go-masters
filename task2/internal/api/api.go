package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"go-masters/task2/internal/errs"
	"go-masters/task2/internal/rates"
)

type API struct{}

type RateRequest struct {
	From string `json:"from"`
	To   string `json:"to"`
}

func (a *API) Rate(w http.ResponseWriter, r *http.Request) {

	var rReq = RateRequest{}

	err := json.NewDecoder(r.Body).Decode(&rReq)
	if err != nil {
		a.WriteError(w, &errs.ErrBadRequest{})
		return
	}

	if rReq.From == "" || rReq.To == "" {
		a.WriteError(w, &errs.ErrBadRequest{})
		return
	}

	pair := strings.ToUpper(rReq.From) + "-" + strings.ToUpper(rReq.To)
	rate, ok := rates.Rates[pair]
	if !ok {
		a.WriteError(w, &errs.ErrNotFound{Pair: pair})
		return
	}

	response := fmt.Sprintf(" %s: %.2f", pair, rate)
	w.Header().Set("Content-Type", "application/json")

	w.Write([]byte(response))
}

func (a *API) WriteError(w http.ResponseWriter, err error) {
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
