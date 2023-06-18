package handler

import (
	"context"
	"encoding/json"
	"moma-api/db"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type RateHandler struct {
	db db.RateDB
}

func (rh *RateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fromCode := r.URL.Query().Get("fromCode")
	toCode := r.URL.Query().Get("toCode")
	timeoutStr := r.URL.Query().Get("timeout") // for example 10s, only second unit is supported

	var timeout time.Duration
	if timeoutStr == "" {
		timeout = time.Second * 10
	} else {
		timeoutStr = strings.TrimSuffix(timeoutStr, "s")
		intVar, err := strconv.Atoi(timeoutStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		timeout = time.Second * time.Duration(intVar)
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer func() {
		cancel()
	}()

	rate, err := rh.db.GetRate(ctx, fromCode, toCode)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(rate)
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}

func NewRateHandler() *RateHandler {
	return &RateHandler{
		db: db.NewRateDB(db.NewClient()),
	}
}