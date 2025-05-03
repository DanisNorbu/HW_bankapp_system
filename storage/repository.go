package repository

import (
	"bankapp_system/models"
	"fmt"
	"sync"

	"github.com/shopspring/decimal"
)

type TempRepository struct {
	users        map[string]models.Client
	accounts     map[string]models.Wallet
	cards        map[string]models.VirtualCard
	loans        map[string]models.CreditAgreement
	transactions []models.MoneyTransfer
	userIndex    map[string]string
	emailIndex   map[string]string
	accountIndex map[string][]string
	cardIndex    map[string][]string
	loanIndex    map[string][]string
	mu           sync.RWMutex // Используем RWMutex для синхронизации
}

var repo *TempRepository

// Инициализация хранилища
func InitStorage() {
	repo = &TempRepository{
		users:        make(map[string]models.Client),
		accounts:     make(map[string]models.Wallet),
		cards:        make(map[string]models.VirtualCard),
		loans:        make(map[string]models.CreditAgreement),
		transactions: make([]models.MoneyTransfer, 0),
		userIndex:    make(map[string]string),
		emailIndex:   make(map[string]string),
		accountIndex: make(map[string][]string),
		cardIndex:    make(map[string][]string),
		loanIndex:    make(map[string][]string),
	}
}

// Добавление клиента с проверкой уникальности
func AddUser(client models.Client) error {
	repo.mu.Lock() // Блокировка записи для добавления пользователя
	defer repo.mu.Unlock()

	// Проверки на уникальность
	if _, exists := repo.userIndex[client.Username]; exists {
		return fmt.Errorf("username '%s' already taken", client.Username)
	}
	if _, exists := repo.emailIndex[client.Email]; exists {
		return fmt.Errorf("email '%s' already registered", client.Email)
	}

	// Добавляем клиента в репозиторий
	repo.users[client.ID] = client
	repo.userIndex[client.Username] = client.ID
	repo.emailIndex[client.Email] = client.ID
	return nil
}

// Получение клиента по username
func GetUserByUsername(username string) (models.Client, bool) {
	repo.mu.RLock() // Блокировка чтения для поиска
	defer repo.mu.RUnlock()

	userID, ok := repo.userIndex[username]
	if !ok {
		return models.Client{}, false
	}
	client, ok := repo.users[userID]
	return client, ok
}

// Добавление счета
func AddAccount(wallet models.Wallet) error {
	repo.mu.Lock() // Блокировка записи для добавления счета
	defer repo.mu.Unlock()
	if _, exists := repo.users[wallet.UserID]; !exists {
		return fmt.Errorf("client with ID %s not found", wallet.UserID)
	}
	repo.accounts[wallet.ID] = wallet
	repo.accountIndex[wallet.UserID] = append(repo.accountIndex[wallet.UserID], wallet.ID)
	return nil
}

// Получение счета по accountID
func GetAccount(accountID string) (models.Wallet, bool) {
	repo.mu.RLock() // Блокировка чтения для получения счета
	defer repo.mu.RUnlock()
	acc, ok := repo.accounts[accountID]
	return acc, ok
}

// Получение всех счетов пользователя
func GetUserAccounts(userID string) []models.Wallet {
	repo.mu.RLock() // Блокировка чтения для получения счетов
	defer repo.mu.RUnlock()
	accountIDs := repo.accountIndex[userID]
	accounts := make([]models.Wallet, 0, len(accountIDs))
	for _, id := range accountIDs {
		if acc, ok := repo.accounts[id]; ok {
			accounts = append(accounts, acc)
		}
	}
	return accounts
}

// Обновление баланса счета
func UpdateAccountBalance(accountID string, amount decimal.Decimal) error {
	repo.mu.Lock() // Блокировка записи для обновления баланса
	defer repo.mu.Unlock()

	acc, ok := repo.accounts[accountID]
	if !ok {
		return fmt.Errorf("wallet %s not found", accountID)
	}

	newBalance := acc.Balance.Add(amount)
	if newBalance.IsNegative() {
		return fmt.Errorf("insufficient funds")
	}

	acc.Balance = newBalance
	repo.accounts[accountID] = acc
	return nil
}

// Добавление транзакции
func AddTransaction(tx models.MoneyTransfer) {
	repo.mu.Lock() // Блокировка записи для добавления транзакции
	defer repo.mu.Unlock()
	repo.transactions = append(repo.transactions, tx)
}

// Получение транзакций по счету
func GetAccountTransactions(accountID string) []models.MoneyTransfer {
	repo.mu.RLock() // Блокировка чтения для получения транзакций
	defer repo.mu.RUnlock()
	var accountTxs []models.MoneyTransfer
	for _, tx := range repo.transactions {
		if tx.FromAccountID == accountID || tx.ToAccountID == accountID {
			accountTxs = append(accountTxs, tx)
		}
	}
	return accountTxs
}
