package services

import (
	"bankapp_system/models"
	"bankapp_system/storage"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Функция для регистрации клиента
func RegisterClient(client models.Client) error {
	// Проверка уникальности email и username
	for _, existingClient := range storage.Clients {
		if existingClient.Email == client.Email {
			return errors.New("email already in use")
		}
		if existingClient.Username == client.Username {
			return errors.New("username already taken")
		}
	}

	// Генерация уникального ID и добавление клиента
	client.ID = generateID()
	client.Created = time.Now()
	storage.Clients = append(storage.Clients, client)

	return nil
}

// Функция для аутентификации клиента и генерации JWT
func AuthenticateClient(username, password string) (string, error) {
	var client models.Client
	var found bool
	for _, c := range storage.Clients {
		if c.Username == username && c.Password == password {
			client = c
			found = true
			break
		}
	}

	if !found {
		return "", errors.New("invalid credentials")
	}

	// Генерация JWT
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["clientID"] = client.ID
	claims["username"] = client.Username
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	// Подпись токена
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Генерация уникального ID для клиента
func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// Создание счета (кошелька)
func CreateWallet(wallet models.Wallet) error {
	for _, existingWallet := range storage.Wallets {
		if existingWallet.ID == wallet.ID {
			return errors.New("wallet already exists")
		}
	}

	// Добавляем счет в репозиторий
	storage.Wallets = append(storage.Wallets, wallet)

	return nil
}

// Пополнение счета
func DepositFunds(accountID string, amount float64) error {
	// Логика пополнения баланса
	for i, wallet := range storage.Wallets {
		if wallet.ID == accountID {
			storage.Wallets[i].Balance += amount
			return nil
		}
	}

	return errors.New("wallet not found")
}

// Перевод средств между счетами
func TransferFunds(fromAccountID, toAccountID string, amount float64) error {
	var fromWallet, toWallet *models.Wallet

	// Находим оба счета
	for i, wallet := range storage.Wallets {
		if wallet.ID == fromAccountID {
			fromWallet = &storage.Wallets[i]
		}
		if wallet.ID == toAccountID {
			toWallet = &storage.Wallets[i]
		}
	}

	if fromWallet == nil {
		return errors.New("source wallet not found")
	}
	if toWallet == nil {
		return errors.New("destination wallet not found")
	}

	// Проверка на достаточность средств
	if fromWallet.Balance < amount {
		return errors.New("insufficient funds")
	}

	// Выполняем перевод
	fromWallet.Balance -= amount
	toWallet.Balance += amount

	return nil
}
