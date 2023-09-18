package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
	"log/slog"
)

type Storage interface {
	CreateAccount(account *Account) error
	DeleteAccount(accountUUID string) error
	UpdateAccount(account *Account) error
	GetAccounts() ([]*Account, error)
	GetAccountByID(accountUUID uuid.UUID) (*Account, error)
}

type PostgresStorage struct {
	db *sql.DB
}

func (s *PostgresStorage) Init() error {
	return s.createAccountTable()
}

func (s *PostgresStorage) createAccountTable() error {
	query := `create table if not exists account (
		id uuid primary key,
		first_name varchar(50),
		last_name varchar(50),
		number serial,
		balance serial,
		created_at timestamp
	)`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStorage) CreateAccount(account *Account) error {
	queryString := `
	insert into account 
	(id, first_name, last_name, number, balance, created_at)
	values ($1, $2, $3, $4, $5, $6)
	`
	resp, err := s.db.Query(queryString,
		account.UUID,
		account.FirstName,
		account.LastName,
		account.Number,
		account.Balance,
		account.CreatedAt)
	if err != nil {
		return err
	}

	slog.Info("CreateAccount done")
	slog.Info("resp", resp)

	return nil
}

func (s *PostgresStorage) DeleteAccount(accountUUID string) error {
	//TODO implement me
	panic("implement me")
}

func (s *PostgresStorage) UpdateAccount(account *Account) error {
	//TODO implement me
	panic("implement me")
}

func (s *PostgresStorage) GetAccounts() ([]*Account, error) {
	rows, err := s.db.Query("select * from account")
	if err != nil {
		return nil, err
	}

	var accounts []*Account

	for rows.Next() {
		account, err := scanIntoAccount(rows)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (s *PostgresStorage) GetAccountByID(accountUUID uuid.UUID) (*Account, error) {
	rows, err := s.db.Query("select * from account where id = $1", accountUUID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccount(rows)
	}

	return nil, fmt.Errorf("account %s not found", accountUUID)
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

func scanIntoAccount(rows *sql.Rows) (*Account, error) {
	account := &Account{}
	err := rows.Scan(
		&account.UUID,
		&account.FirstName,
		&account.LastName,
		&account.Number,
		&account.Balance,
		&account.CreatedAt,
	)

	return account, err
}
