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

func (s *PostgresStorage) Init() error {
	return s.createAccountTable()
}

func (s *PostgresStorage) createAccountTable() error {
	query := `create table if not exists account (
		id serial primary key,
		first_name varchar(50),
		last_name varchar(50),
		number serial,
		balance int,
		created_at timestamp
	)`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStorage) CreateAccount(account *Account) error {
	//TODO implement me
	panic("implement me")
}

func (s *PostgresStorage) DeleteAccount(accountUUID string) error {
	//TODO implement me
	panic("implement me")
}

func (s *PostgresStorage) UpdateAccount(account *Account) error {
	//TODO implement me
	panic("implement me")
}

func (s *PostgresStorage) GetAccount(accountUUID string) (*Account, error) {
	//TODO implement me
	panic("implement me")
}

func NewPostgresStorage() (*PostgresStorage, error) {
	conStr := "user=postgres dbname=postgres password=gobank sslmode=disable"
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
