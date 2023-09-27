package main

import (
	"flag"
	"log/slog"
	"os"
)

func seedAccount(s Storage, fName, lName, password string) *Account {
	acc, err := NewAccount(fName, lName, password)
	if err != nil {
		slog.Error("can't seed accounts (NewAccount)", slog.String("error", err.Error()))
		os.Exit(1)
	}
	err = s.CreateAccount(acc)
	if err != nil {
		slog.Error("can't seed accounts (CreateAccount)", slog.String("error", err.Error()))
		os.Exit(1)
	}
	return acc
}

func seedAccounts(s Storage) {
	seedAccount(s, "Antony", "Gold", "password")
	seedAccount(s, "Robert", "Full", "password2")
	seedAccount(s, "George", "Brown", "password3")

}

func main() {
	seedFlag := flag.Bool("seed", false, "seed accounts to db")
	flag.Parse()

	storage, err := NewPostgresStorage()
	if err != nil {
		slog.Error("can't create storage", slog.String("error", err.Error()))
		os.Exit(1)
	}

	if err := storage.Init(); err != nil {
		slog.Error("can't init storage", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// seed accounts to storage for testing
	if *seedFlag {
		slog.Info("seeding accounts to storage")
		seedAccounts(storage)
	}

	server := NewAPIServer(":3000", storage)

	server.Run()
}
