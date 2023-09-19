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
	DeleteAccount(accountUUID uuid.UUID) error
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
	query := `CREATE TABLE IF NOT EXISTS account (
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
	INSERT INTO  account 
	(id, first_name, last_name, number, balance, created_at)
	values ($1, $2, $3, $4, $5, $6)
	`
	resp, err := s.db.Exec(queryString,
		account.UUID,
		account.FirstName,
		account.LastName,
		account.Number,
		account.Balance,
		account.CreatedAt)
	if err != nil {
		return err
	}

	slog.Info("Account created", slog.Any("response", resp))

	return nil
}

func (s *PostgresStorage) DeleteAccount(accountUUID uuid.UUID) error {
	queryString := `
		DELETE FROM account 
       	WHERE id = $1
       	`

	result, err := s.db.Exec(queryString, accountUUID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("account with ID %s does not exist", accountUUID.String())
	}
	return nil
}

func (s *PostgresStorage) UpdateAccount(account *Account) error {
	//TODO implement me
	panic("implement me")
}

func (s *PostgresStorage) GetAccounts() ([]*Account, error) {
	queryString := `
		SELECT * 
		FROM account
       	`

	rows, err := s.db.Query(queryString)
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
	queryString := `
		SELECT * 
		FROM account 
		WHERE id = $1
       	`

	rows, err := s.db.Query(queryString, accountUUID)
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
