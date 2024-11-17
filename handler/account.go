package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"moma-api/db"
	"moma-api/db/model"
	"moma-api/db/sql"
)

// Account represents the structure of the JSON payload
type Account struct {
	UserID    string `json:"user_id"`
	UserName  string `json:"user_name"`
	UserEmail string `json:"user_email"`
}

// RequestPayload represents the structure of the incoming JSON
type RequestPayload struct {
	Account Account `json:"account"`
}

type AccountHandler struct {
	DB db.AccountDB
}

func (rh *AccountHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is supported", http.StatusMethodNotAllowed)
		return
	}

	var payload RequestPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid JSON in request body", http.StatusBadRequest)
		return
	}

	account := payload.Account
	if account.UserID == "" || account.UserName == "" || account.UserEmail == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	_, err := rh.DB.GetAccount(context.Background(), account.UserID)
	if err == nil {
		http.Error(w, "Account already exists", http.StatusConflict)
	}

	// Insert into the database
	err = rh.DB.CreateAccount(context.Background(), &model.Account{
		UserID:    account.UserID,
		UserName:  account.UserName,
		UserEmail: account.UserEmail,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to insert data into database: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte("Account created"))
}

func (rh *AccountHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")

	if len(userID) == 0 {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
	}

	account, err := rh.DB.GetAccount(r.Context(), userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get data from database: %v", err), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	resp := NewResponse(0, account, "")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Printf("fail to encode account due to %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func NewAccountHandler() (*AccountHandler, error) {
	accountDB, err := sql.NewAccountDB()
	if err != nil {
		return nil, err
	}

	return &AccountHandler{
		DB: accountDB,
	}, nil
}
