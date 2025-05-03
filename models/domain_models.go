package models

import (
	"time"
)

// Структура для клиента
type Client struct {
	ID       string    `json:"clientID"`
	Email    string    `json:"email"`
	Username string    `json:"username"`
	Password string    `json:"password"`
	Created  time.Time `json:"created"`
}

// Структура для счета
type Wallet struct {
	ID       string    `json:"walletID"`
	ClientID string    `json:"clientID"`
	Balance  float64   `json:"balance"`
	Currency string    `json:"currency"`
	Created  time.Time `json:"created"`
	Updated  time.Time `json:"updated"`
}

// Структура для виртуальной карты
type VirtualCard struct {
	ID             string    `json:"vCardID"`
	WalletID       string    `json:"walletID"`
	CardNumber     string    `json:"cardNumber"`
	ExpirationDate string    `json:"expirationDate"`
	CVV            string    `json:"cvv"`
	Created        time.Time `json:"created"`
}

// Структура для перевода средств
type MoneyTransfer struct {
	ID            string    `json:"transferID"`
	FromAccountID string    `json:"fromWalletID"`
	ToAccountID   string    `json:"toWalletID"`
	Amount        float64   `json:"amount"`
	TransferDate  time.Time `json:"transferDate"`
}

// Структура для кредита
type CreditAgreement struct {
	ID             string    `json:"creditID"`
	ClientID       string    `json:"clientID"`
	Amount         float64   `json:"amount"`
	InterestRate   float64   `json:"interestRate"`
	Term           int       `json:"term"`
	MonthlyPayment float64   `json:"monthlyPayment"`
	Created        time.Time `json:"created"`
}
