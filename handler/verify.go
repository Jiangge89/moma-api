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

type verifyReceiptResult struct {
	success bool `json:"success"`
}

func (vh *VerifyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	transactionId := r.URL.Query().Get("transaction_id")

	resData := verifyReceiptResult{}
	code := success
	errMsg := ""

	err := service.VerifyReceipt(transactionId)

	if err != nil {
		resData.success = false
		code = errFailVerifyReceipt
		errMsg = fmt.Sprintf("fail to verify receipt due to: %v\n", err)
	}

	resp := NewResponse(code, resData, errMsg)

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
