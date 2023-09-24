package main

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"math/rand"
	"time"
)

const (
	countryCode = "UA"
)

type TransferRequest struct {
	ToAccountNumber string `json:"to_account_number"`
	Amount          int64  `json:"amount"`
}

type CreateAccountRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type Account struct {
	UUID      uuid.UUID `json:"uuid"`
	FirstName string    `json:"name"`
	LastName  string    `json:"last_name"`
	Number    string    `json:"number"`
	Balance   int64     `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

func NewAccount(firstName, lastName string) *Account {
	return &Account{
		UUID:      uuid.NewV4(),
		FirstName: firstName,
		LastName:  lastName,
		Number:    ibanGenerator(),
		CreatedAt: time.Now().UTC(),
	}
}

func ibanGenerator() string {
	checkDigits := fmt.Sprintf("%02d", rand.Intn(100))
	bankCode := fmt.Sprintf("%04d", rand.Intn(10000))
	accountNumber := fmt.Sprintf("%010d", rand.Intn(10000000000))

	iban := countryCode + checkDigits + bankCode + accountNumber

	return iban
}
