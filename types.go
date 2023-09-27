package main

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"time"
)

const (
	countryCode = "UA"
)

type LoginRequest struct {
	Number   string `json:"number"`
	Password string `json:"password"`
}

type loginResponse struct {
	Number string `json:"number"`
	Token  string `json:"token"`
}

type TransferRequest struct {
	ToAccountNumber string `json:"to_account_number"`
	Amount          int64  `json:"amount"`
}

type CreateAccountRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
}

type Account struct {
	UUID      uuid.UUID `json:"uuid"`
	FirstName string    `json:"name"`
	LastName  string    `json:"last_name"`
	Number    string    `json:"number"`
	Encrypted string    `json:"-"`
	Balance   int64     `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

func NewAccount(firstName, lastName, password string) (*Account, error) {
	passwordEnc, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("can't hash password: %w", err)
	}
	return &Account{
		UUID:      uuid.NewV4(),
		FirstName: firstName,
		LastName:  lastName,
		Encrypted: string(passwordEnc),
		Number:    ibanGenerator(),
		CreatedAt: time.Now().UTC(),
	}, nil
}

func ibanGenerator() string {
	checkDigits := fmt.Sprintf("%02d", rand.Intn(100))
	bankCode := fmt.Sprintf("%04d", rand.Intn(10000))
	accountNumber := fmt.Sprintf("%010d", rand.Intn(10000000000))

	iban := countryCode + checkDigits + bankCode + accountNumber

	return iban
}
