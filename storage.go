package main

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(account *Account) error
	DeleteAccount(accountUUID string) error
	UpdateAccount(account *Account) error
	GetAccount(accountUUID string) (*Account, error)
}

type PostgresStorage struct {
	db *sql.DB
}

func (p PostgresStorage) CreateAccount(account *Account) error {
	//TODO implement me
	panic("implement me")
}

func (p PostgresStorage) DeleteAccount(accountUUID string) error {
	//TODO implement me
	panic("implement me")
}

func (p PostgresStorage) UpdateAccount(account *Account) error {
	//TODO implement me
	panic("implement me")
}

func (p PostgresStorage) GetAccount(accountUUID string) (*Account, error) {
	//TODO implement me
	panic("implement me")
}

func NewPostgresStorage() (*PostgresStorage, error) {
	conStr := "user=admin dbname=gobank-postgres password=gobank sslmode=verify-full"
	db, err := sql.Open("postgres", conStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStorage{
		db: db,
	}, nil
}
