package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"moma-api/service"
	"net/http"
)

type VerifyHandler struct {
}

func (vh *VerifyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	transactionId := r.URL.Query().Get("transaction_id")

	code := success
	errMsg := ""

	err := service.VerifyReceipt(transactionId)

	if err != nil {
		code = errFailVerifyReceipt
		errMsg = fmt.Sprintf("fail to verify receipt due to: %v\n", err)
	}

	resp := NewResponse(code, nil, errMsg)

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Printf("fail to encode resp for verifyReceipt due to %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

func NewVerifyHandler() (*VerifyHandler, error) {
	return &VerifyHandler{}, nil
}
