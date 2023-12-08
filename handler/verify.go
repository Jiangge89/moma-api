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

type VerifyReceiptResponse struct {
	ExpiredDate int64 `json:"expired_date"`
}

func (vh *VerifyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	transactionId := r.URL.Query().Get("transaction_id")

	code := success
	errMsg := ""

	expiredDate, err := service.VerifyReceipt(transactionId)
	if err != nil {
		code = errFailVerifyReceipt
		errMsg = fmt.Sprintf("fail to verify receipt due to: %v\n", err)
	}

	data := VerifyReceiptResponse{ExpiredDate: expiredDate}
	resp := NewResponse(code, data, errMsg)

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
