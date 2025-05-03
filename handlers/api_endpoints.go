package handlers

import (
	"bankapp_system/models"
	"bankapp_system/services"
	"encoding/json"
	"fmt"
	"net/http"
)

// Регистрация клиента
func RegisterClient(w http.ResponseWriter, r *http.Request) {
	var client models.Client
	if err := json.NewDecoder(r.Body).Decode(&client); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Логика регистрации
	if err := services.RegisterClient(client); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Client registered successfully")
}

// Аутентификация клиента
func ClientLogin(w http.ResponseWriter, r *http.Request) {
	var loginData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Логика аутентификации
	token, err := services.AuthenticateClient(loginData.Username, loginData.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Возвращаем токен
	w.Header().Set("Authorization", "Bearer "+token)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Login successful")
}

// Создание счета
func CreateAccount(w http.ResponseWriter, r *http.Request) {
	var wallet models.Wallet
	if err := json.NewDecoder(r.Body).Decode(&wallet); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Логика создания счета
	if err := services.CreateWallet(wallet); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Wallet created successfully")
}

// Пополнение счета
func Deposit(w http.ResponseWriter, r *http.Request) {
	var depositData struct {
		AccountID string  `json:"accountID"`
		Amount    float64 `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&depositData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Логика пополнения счета
	if err := services.DepositFunds(depositData.AccountID, depositData.Amount); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Funds deposited successfully")
}

// Перевод средств между счетами
func Transfer(w http.ResponseWriter, r *http.Request) {
	var transferData struct {
		FromAccountID string  `json:"fromAccountID"`
		ToAccountID   string  `json:"toAccountID"`
		Amount        float64 `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&transferData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Логика перевода средств
	if err := services.TransferFunds(transferData.FromAccountID, transferData.ToAccountID, transferData.Amount); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Transfer successful")
}
