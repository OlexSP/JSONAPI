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
	GetAccountByNumber(accountNumber string) (*Account, error)
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
		number varchar(35),
		balance serial,
		created_at timestamp,
        password_hash varchar(255)
	)`

	_, err := s.db.Exec(query)

	slog.Info(query, slog.Any("error", err))

	return err
}

func (s *PostgresStorage) CreateAccount(account *Account) error {
	queryString := `
	INSERT INTO  account 
	(id, first_name, last_name, number, balance, created_at, password_hash)
	values ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := s.db.Exec(queryString,
		account.UUID,
		account.FirstName,
		account.LastName,
		account.Number,
		account.Balance,
		account.CreatedAt,
		account.Encrypted)
	if err != nil {
		return err
	}

	slog.Info("Account created", slog.Any("accountID", account.UUID))

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
		return fmt.Errorf("invalid ID %s ", accountUUID.String())
	}

	slog.Info(queryString, slog.Any("accountID", accountUUID))

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

	slog.Info(queryString, slog.Any("accounts", accounts))

	return accounts, nil
}

func (s *PostgresStorage) GetAccountByNumber(accountNumber string) (*Account, error) {
	queryString := `
		SELECT * 
		FROM account 
		WHERE number = $1
       	`

	rows, err := s.db.Query(queryString, accountNumber)
	if err != nil {
		return nil, err
	}

	slog.Info(queryString, slog.Any("err", err), slog.Any("Account Number", accountNumber))

	for rows.Next() {
		return scanIntoAccount(rows)
	}

	return nil, fmt.Errorf("invalid Account Number %s", accountNumber)
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

	slog.Info(queryString, slog.Any("err", err), slog.Any("ID", accountUUID))

	for rows.Next() {
		return scanIntoAccount(rows)
	}

	return nil, fmt.Errorf("invalid ID %s", accountUUID)
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
		&account.Encrypted,
	)

	return account, err
}
