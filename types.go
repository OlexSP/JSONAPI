package main

import (
	uuid "github.com/satori/go.uuid"
	"math/rand"
)

type Account struct {
	UUID      string `json:"uuid"`
	FirstName string `json:"name"`
	LastName  string `json:"last_name"`
	Number    int64  `json:"number"`
	Balance   int64  `json:"balance"`
}

func NewAccount(firstName, lastName string) *Account {
	return &Account{
		UUID:      uuid.NewV4().String(),
		FirstName: firstName,
		LastName:  lastName,
		Number:    int64(rand.Intn(100000)),
	}
}
