package utils

import (
	"fmt"
	"math/rand"
	"time"
)

// Генерация случайного ID
func GenerateRandomID() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%d", rand.Int63())
}

// Генерация номера карты с использованием алгоритма Луна
func GenerateCardNumber() string {
	var cardNumber string
	for i := 0; i < 16; i++ {
		cardNumber += fmt.Sprintf("%d", rand.Intn(10))
	}
	return cardNumber
}

// Проверка правильности номера карты по алгоритму Луна
func ValidateCardNumber(cardNumber string) bool {
	var sum int
	var shouldDouble bool

	// Преобразование строки в массив цифр
	for i := len(cardNumber) - 1; i >= 0; i-- {
		num := int(cardNumber[i] - '0')

		// Удваиваем каждую вторую цифру с конца
		if shouldDouble {
			num *= 2
			if num > 9 {
				num -= 9
			}
		}

		// Добавляем к сумме
		sum += num
		shouldDouble = !shouldDouble
	}

	// Если сумма делится на 10, то карта валидна
	return sum%10 == 0
}

// Генерация CVV для карты
func GenerateCVV() string {
	var cvv string
	for i := 0; i < 3; i++ {
		cvv += fmt.Sprintf("%d", rand.Intn(10))
	}
	return cvv
}
